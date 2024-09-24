package entity

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type Channel struct {
	ID      string
	Name    string
	Private bool
	Clients map[*Client]bool
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
		ID:      id,
		Name:    name,
		Private: private,
		Clients: make(map[*Client]bool),
	}, nil
}

func (c *Channel) RegisterClientInChannel(client *Client) {
	c.Clients[client] = true
}

func (c *Channel) UnRegisterClientInChannel(client *Client) {
	delete(c.Clients, client)
}

func (c *Channel) FindClientByID(id string) *Client {
	for client := range c.Clients {
		if client.ID == id {
			return client
		}
	}
	return nil
}

func (c *Channel) FindClientByUserID(userID string) *Client {
	for client := range c.Clients {
		if client.UserID == userID {
			return client
		}
	}
	return nil
}
