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
	Register        chan *clientManager
	unregister      chan *clientManager
	broadcast       chan []byte
}

func NewHubManager(hub *entity.Hub) HubManager {
	return HubManager{
		Hub:             hub,
		clientManagers:  make(map[*clientManager]bool),
		channelManagers: make(map[*channelManager]bool),
		Register:        make(chan *clientManager),
		unregister:      make(chan *clientManager),
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

func (hm *HubManager) registerClient(clientM *clientManager) {
	hm.Hub.RegisterClient(clientM.client)
	hm.clientManagers[clientM] = true
}

func (hm *HubManager) unregisterClient(clientM *clientManager) {
	for cm := range hm.channelManagers {
		cm.unregister <- clientM
	}
	hm.Hub.UnRegisterClient(clientM.client)
	delete(hm.clientManagers, clientM)
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
func (hm *HubManager) RegisterClientManagerInChannelManager(clientManager *clientManager) {
	for cm := range hm.channelManagers {
		if !cm.isInChannel(clientManager) {
			cm.register <- clientManager
		}
	}
}

// channelManagerから該当するclientManagerの登録を削除する
func (hm *HubManager) UnRegisterClientManagerInChannelManager(clientManager *clientManager) {
	for cm := range hm.channelManagers {
		if cm.isInChannel(clientManager) {
			cm.unregister <- clientManager
		}
	}
}
