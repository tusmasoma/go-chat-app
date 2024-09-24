package handler

import (
	"context"

	"github.com/tusmasoma/go-chat-app/entity"
)

type HubManager interface{}

type hubManager struct {
	hub             *entity.Hub
	clientManagers  map[*clientManager]bool
	channelManagers map[*channelManager]bool
	register        chan *entity.Client
	unregister      chan *entity.Client
	broadcast       chan []byte
}

func NewHubHandler(hub *entity.Hub) HubManager {
	return &hubManager{
		hub:        hub,
		register:   make(chan *entity.Client),
		unregister: make(chan *entity.Client),
		broadcast:  make(chan []byte),
	}
}

func (hm *hubManager) Run(ctx context.Context) {
	for {
		select {
		case client := <-hm.register:
			hm.registerClient(client)
		case client := <-hm.unregister:
			hm.hub.UnRegisterClient(client)
		case message := <-hm.broadcast:
			hm.broadcastToClients(message)
		}
	}
}

func (hm *hubManager) registerClient(client *entity.Client) {
	hm.hub.RegisterClient(client)
}

func (hm *hubManager) unregisterClient(client *entity.Client) {
	for cm := range hm.channelManagers {
		cm.unregister <- client
	}
	hm.hub.UnRegisterClient(client)
}

func (hm *hubManager) broadcastToClients(message []byte) {
	for cm := range hm.clientManagers {
		cm.send <- message
	}
}

func (hm *hubManager) findChannelManagerByChannelID(channelID string) *channelManager {
	for cm := range hm.channelManagers {
		if cm.channel.ID == channelID {
			return cm
		}
	}
	return nil
}
