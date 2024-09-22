package entity

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type Hub struct {
	ID         string
	Name       string
	clients    map[*Client]bool
	channels   map[*Channel]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
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
		clients:    make(map[*Client]bool),
		channels:   make(map[*Channel]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}, nil
}

func (h *Hub) registerClient(client *Client) {
	h.clients[client] = true
}

func (h *Hub) unregisterClient(client *Client) {
	for channel := range client.Channels {
		channel.unregister <- client
	}
	delete(h.clients, client)
}

func (h *Hub) broadcastToClients(message []byte) {
	for client := range h.clients {
		client.Send <- message
	}
}

func (h *Hub) FindChannelByID(id string) *Channel {
	for channel := range h.channels {
		if channel.ID == id {
			return channel
		}
	}
	return nil
}

func (h *Hub) FindChannelByName(name string) *Channel {
	for channel := range h.channels {
		if channel.Name == name {
			return channel
		}
	}
	return nil
}
