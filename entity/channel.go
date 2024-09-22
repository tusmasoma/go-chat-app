package entity

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type Channel struct {
	ID         string
	Name       string
	Private    bool
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	// broadcast  chan *WSMessage
}

func NewChannel(id, name string, private bool) (*Channel, error) {
	if id == "" {
		id = uuid.New().String()
	}
	if name == "" {
		log.Error("name is required")
		return nil, errors.New("name is required")
	}
	return &Channel{
		ID:         id,
		Name:       name,
		Private:    private,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		// broadcast:  make(chan *WSMessage),
	}, nil
}

func (c *Channel) broadcastToClientsInChannel(message []byte) {
	for client := range c.clients {
		client.Send <- message
	}
}

func (c *Channel) FindClientByID(id string) *Client {
	for client := range c.clients {
		if client.ID == id {
			return client
		}
	}
	return nil
}

func (c *Channel) FindClientByUserID(userID string) *Client {
	for client := range c.clients {
		if client.UserID == userID {
			return client
		}
	}
	return nil
}
