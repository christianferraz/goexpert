package services

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type MessageKind int

const (
	PlaceBid MessageKind = iota
	FailedPlaceBid
	SuccessFullyPlacedBid
	NewBidPlaced
)

type Message struct {
	Message string
	Kind    MessageKind
	UserID  uuid.UUID
	Amount  float64
}

type AuctionLobby struct {
	sync.Mutex
	Rooms map[uuid.UUID]*AuctionRoom
}

type AuctionRoom struct {
	Id          uuid.UUID
	Context     context.Context
	Broadcast   chan Message
	Register    chan *Client
	Unregister  chan *Client
	Clients     map[uuid.UUID]*Client
	BidsService *BidsService
}

func NewAuctionRoom(ctx context.Context, id uuid.UUID, bs *BidsService) *AuctionRoom {
	return &AuctionRoom{
		Id:          id,
		Context:     ctx,
		Broadcast:   make(chan Message),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Clients:     make(map[uuid.UUID]*Client),
		BidsService: bs,
	}
}

type Client struct {
	Rooms  *AuctionRoom
	Conn   *websocket.Conn
	Send   chan Message
	UserId uuid.UUID
}

func NewClient(room *AuctionRoom, conn *websocket.Conn, userId uuid.UUID) *Client {
	return &Client{
		Rooms:  room,
		Conn:   conn,
		Send:   make(chan Message, 512),
		UserId: userId,
	}
}

func (r *AuctionRoom) Run() {
	defer func() {
		close(r.Broadcast)
		close(r.Register)
		close(r.Unregister)
	}()
	for {
		select {
		case client := <-r.Register:
			r.registerClient(client)
		case client := <-r.Unregister:
			r.unregisterClient(client)

		case <-r.Context.Done():
			return
		}
	}
}

func (r *AuctionRoom) registerClient(c *Client) {
	slog.Info("Registering client", "Client", c)
	r.Clients[c.UserId] = c
}

func (r *AuctionRoom) unregisterClient(c *Client) {
	slog.Info("Unregistering client", "Client", c)
	if _, ok := r.Clients[c.UserId]; ok {
		delete(r.Clients, c.UserId)
		close(c.Send)
	}
}

func (r *AuctionRoom) broadcastMessage(msg Message) {
	slog.Info("Broadcasting message", "RoomID", r.Id, "Message", msg, "Clients", len(r.Clients))
	switch msg.Kind {
	case PlaceBid:
		bid, err := r.BidsService.Placebid(r.Context, r.Id, msg.UserID, msg.Amount)
		if err != nil {
			if errors.Is(err, ErrBidIsTooLow) {
				if client, ok := r.Clients[msg.UserID]; ok {
					client.Send <- Message{Kind: FailedPlaceBid, Message: ErrBidIsTooLow.Error()}
				}
			}
		}
		if client, ok := r.Clients[msg.UserID]; ok {
			client.Send <- Message{Kind: SuccessFullyPlacedBid, Message: "Bid placed with success", UserID: msg.UserID, Amount: bid.Amount}
		}
		for id, client := range r.Clients {
			newBidMessage := Message{Kind: NewBidPlaced, Message: "New bid placed", UserID: msg.UserID, Amount: bid.Amount}
			if id == msg.UserID {
				continue
			}
			client.Send <- newBidMessage
		}
	}
}
