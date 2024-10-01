package websocket

import (
	"context"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/go-chat-app/entity"
	"github.com/tusmasoma/go-chat-app/repository"
)

// type ChannelManager interface{}

type channelManager struct {
	channel        *entity.Channel
	clientManagers map[*clientManager]bool
	register       chan *clientManager
	unregister     chan *clientManager
	broadcast      chan *entity.Message
	psr            repository.PubSubRepository
}

func NewChannelManager(channel *entity.Channel, psr repository.PubSubRepository) *channelManager { //nolint:revive // This function is used in other packages
	return &channelManager{
		channel:        channel,
		clientManagers: make(map[*clientManager]bool),
		register:       make(chan *clientManager),
		unregister:     make(chan *clientManager),
		broadcast:      make(chan *entity.Message),
		psr:            psr,
	}
}

func (cm *channelManager) Run(ctx context.Context) {
	go cm.subscribeToChannelMessages(ctx)

	for {
		select {
		case clientM := <-cm.register:
			cm.registerClientInChannel(clientM)
		case clientM := <-cm.unregister:
			cm.unregisterClientInChannel(clientM)
		case message := <-cm.broadcast:
			cm.publishChannelMessage(ctx, message)
		}
	}
}

func (cm *channelManager) registerClientInChannel(clientM *clientManager) {
	cm.channel.RegisterClientInChannel(clientM.client)
	cm.clientManagers[clientM] = true
}

func (cm *channelManager) unregisterClientInChannel(clientM *clientManager) {
	cm.channel.UnRegisterClientInChannel(clientM.client)
	delete(cm.clientManagers, clientM)
}

func (cm *channelManager) broadcastToClientsInChannel(message []byte) {
	for clientManger := range cm.clientManagers {
		log.Info("Broadcasting message to clients in channel", log.Fstring("message", string(message)))
		clientManger.send <- message
	}
}

func (cm *channelManager) publishChannelMessage(ctx context.Context, message *entity.Message) {
	msg, err := message.Encode()
	if err != nil {
		log.Error("Failed to encode message", log.Ferror(err))
		return
	}
	if err := cm.psr.Publish(ctx, cm.channel.ID, msg); err != nil {
		log.Error("Failed to publish message", log.Ferror(err))
	}
	log.Info("Successfully published message", log.Fstring("channelID", cm.channel.ID))
}

func (cm *channelManager) subscribeToChannelMessages(ctx context.Context) {
	pubsub := cm.psr.Subscribe(ctx, cm.channel.ID)
	defer pubsub.Close()

	msgs := pubsub.Channel()

	for msg := range msgs {
		cm.broadcastToClientsInChannel([]byte(msg.Payload))
	}
}

func (cm *channelManager) isInChannel(client *clientManager) bool {
	_, ok := cm.clientManagers[client]
	return ok
}
