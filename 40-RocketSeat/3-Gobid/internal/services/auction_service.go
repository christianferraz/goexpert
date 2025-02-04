package services

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type MessageKind int

const (
	PlaceBid MessageKind = iota
	FailedPlaceBid
	SuccessFullyPlacedBid
	NewBidPlaced
	AuctionFinished
	InvalidJSON
)

type Message struct {
	Message string      `json:"message,omitempty"`
	Kind    MessageKind `json:"kind"`
	UserID  uuid.UUID   `json:"user_id,omitempty"`
	Amount  float64     `json:"amount,omitempty"`
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
			slog.Info("Auction room has been closed", "RoomID", r.Id)
			for _, client := range r.Clients {
				client.Send <- Message{Kind: AuctionFinished, Message: "Auction has been closed"}
			}
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
					client.Send <- Message{Kind: FailedPlaceBid, Message: ErrBidIsTooLow.Error(), UserID: msg.UserID}
				}
				return
			}
		}
		if client, ok := r.Clients[msg.UserID]; ok {
			client.Send <- Message{Kind: SuccessFullyPlacedBid, Message: "Bid placed with success", UserID: msg.UserID, Amount: bid.BidAmount}
		}
		for id, client := range r.Clients {
			newBidMessage := Message{Kind: NewBidPlaced, Message: "New bid placed", Amount: bid.BidAmount, UserID: msg.UserID}
			if id == msg.UserID {
				continue
			}
			client.Send <- newBidMessage
		}
	case InvalidJSON:
		client, ok := r.Clients[msg.UserID]
		if !ok {
			slog.Info("Client not found", "UserID", msg.UserID)
			return
		}
		client.Send <- msg
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

const (
	maxMessageSize = 512
	readDeadLine   = 60 * time.Second
	writeWait      = 10 * time.Second
	pingPeriod     = (readDeadLine * 9) / 10
)

func (c *Client) ReadEventLoop() {
	defer func() {
		c.Rooms.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(readDeadLine))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(readDeadLine))
		return nil
	})
	for {
		var m Message
		m.UserID = c.UserId
		err := c.Conn.ReadJSON(&m)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("Unexpected error", "Error", err)
				return
			}
			c.Rooms.Broadcast <- Message{
				Kind:    InvalidJSON,
				Message: "Invalid JSON",
				UserID:  m.UserID,
			}
			continue
		}
		c.Rooms.Broadcast <- m
	}
}

func (c *Client) WriteEventLoop() {
	// mandar ping message
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case <-ticker.C:
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteJSON(Message{Kind: websocket.CloseMessage, Message: "Closing websocket connection"})
				return
			}
			if message.Kind == AuctionFinished {
				close(c.Send)
				return
			}
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.Conn.WriteJSON(message)
			if err != nil {
				slog.Error("Error writing message", "Error", err)
				c.Rooms.Unregister <- c
				return
			}
		}
		c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
		if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			slog.Error("Error sending ping message", "Error", err)
			return
		}
	}
}
