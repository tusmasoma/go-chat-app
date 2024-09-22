package entity

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type Client struct {
	ID       string
	UserID   string // TODO: membershipIDに変更すべきか？
	Hub      *Hub
	Channels map[*Channel]bool
	Send     chan []byte
}

func NewClient(id, userID string, hub *Hub) (*Client, error) {
	if id == "" {
		id = uuid.New().String()
	}
	if userID == "" {
		log.Error("userID is required")
		return nil, errors.New("userID is required")
	}
	return &Client{
		ID:       id,
		UserID:   userID,
		Hub:      hub,
		Channels: make(map[*Channel]bool),
		Send:     make(chan []byte),
	}, nil
}

func (c *Client) JoinChannel(channel *Channel) {
	if c.isInChannel(channel) {
		return
	}
	c.Channels[channel] = true
	channel.Register <- c
}

func (c *Client) LeaveChannel(channel *Channel) {
	if !c.isInChannel(channel) {
		return
	}
	channel.UnRegister <- c
	delete(c.Channels, channel)
}

func (c *Client) isInChannel(channel *Channel) bool {
	if _, ok := c.Channels[channel]; ok {
		return true
	}
	return false
}
