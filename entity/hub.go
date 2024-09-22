package entity

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type Hub struct {
	ID         string
	Name       string
	Clients    map[*Client]bool
	Channels   map[*Channel]bool
	Register   chan *Client
	UnRegister chan *Client
	Broadcast  chan []byte
}

func NewHub(id, name string) (*Hub, error) {
	if id == "" {
		id = uuid.New().String()
	}
	if name == "" {
		log.Error("name is required")
		return nil, errors.New("name is required")
	}
	return &Hub{
		ID:         id,
		Name:       name,
		Clients:    make(map[*Client]bool),
		Channels:   make(map[*Channel]bool),
		Register:   make(chan *Client),
		UnRegister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}, nil
}

func (h *Hub) RegisterClient(client *Client) {
	h.Clients[client] = true
}

func (h *Hub) UnRegisterClient(client *Client) {
	for channel := range client.Channels {
		channel.UnRegister <- client
	}
	delete(h.Clients, client)
}

func (h *Hub) BroadcastToClients(message []byte) {
	for client := range h.Clients {
		client.Send <- message
	}
}

func (h *Hub) FindChannelByID(id string) *Channel {
	for channel := range h.Channels {
		if channel.ID == id {
			return channel
		}
	}
	return nil
}

func (h *Hub) FindChannelByName(name string) *Channel {
	for channel := range h.Channels {
		if channel.Name == name {
			return channel
		}
	}
	return nil
}
