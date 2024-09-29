package websocket

import (
	"context"

	"github.com/tusmasoma/go-chat-app/entity"
)

// type HubManager interface{}

type HubManager struct {
	Hub             *entity.Hub
	clientManagers  map[*clientManager]bool
	channelManagers map[*channelManager]bool
	Register        chan *entity.Client
	unregister      chan *entity.Client
	broadcast       chan []byte
}

func NewHubManager(hub *entity.Hub) HubManager {
	return HubManager{
		Hub:             hub,
		clientManagers:  make(map[*clientManager]bool),
		channelManagers: make(map[*channelManager]bool),
		Register:        make(chan *entity.Client),
		unregister:      make(chan *entity.Client),
		broadcast:       make(chan []byte),
	}
}

func (hm *HubManager) Run() {
	for {
		select {
		case client := <-hm.Register:
			hm.registerClient(client)
		case client := <-hm.unregister:
			hm.unregisterClient(client)
		case message := <-hm.broadcast:
			hm.broadcastToClients(message)
		}
	}
}

func (hm *HubManager) registerClient(client *entity.Client) {
	hm.Hub.RegisterClient(client)
}

func (hm *HubManager) unregisterClient(client *entity.Client) {
	for cm := range hm.channelManagers {
		cm.unregister <- client
	}
	hm.Hub.UnRegisterClient(client)
}

func (hm *HubManager) broadcastToClients(message []byte) {
	for cm := range hm.clientManagers {
		cm.send <- message
	}
}

func (hm *HubManager) findChannelManagerByChannelID(channelID string) *channelManager {
	for cm := range hm.channelManagers {
		if cm.channel.ID == channelID {
			return cm
		}
	}
	return nil
}

func (hm *HubManager) RegisterChannel(ctx context.Context, channel *entity.Channel) {
	cm := NewChannelManager(channel, nil)
	go cm.Run(ctx)
	hm.channelManagers[cm] = true
}

func (hm *HubManager) RegisterChannelManager(cm *channelManager) { // 一旦DIのためのメソッドを追加
	hm.channelManagers[cm] = true
}

// HubManagerに登録されているChannelManagerのClientManagerにClientを登録
func (hm *HubManager) RegisterClientManagerInChannelManager(client *clientManager) {
	for cm := range hm.channelManagers {
		cm.clientManagers[client] = true
	}
}
