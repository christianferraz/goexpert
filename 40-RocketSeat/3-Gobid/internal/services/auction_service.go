package services

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type MessageKind int

const (
	PlaceBid MessageKind = iota
)

type Message struct {
	Message string
	Kind    MessageKind
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
			r.Clients[client.UserId] = client
		case client := <-r.Unregister:
			if _, ok := r.Clients[client.UserId]; ok {
				delete(r.Clients, client.UserId)
				close(client.Send)
			}
		case message := <-r.Broadcast:
			for _, client := range r.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(r.Clients, client.UserId)
				}
			}
		case <-r.Context.Done():
			return
		}
	}
}
