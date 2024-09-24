package handler

import (
	"context"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/go-chat-app/entity"
	"github.com/tusmasoma/go-chat-app/repository"
)

type ChannelManager interface{}

type channelManager struct {
	channel        *entity.Channel
	clientManagers map[*clientManager]bool
	register       chan *entity.Client
	unregister     chan *entity.Client
	broadcast      chan *entity.Message
	psr            repository.PubSubRepository
}

func NewChannelManager(channel *entity.Channel, psr repository.PubSubRepository) ChannelManager {
	return &channelManager{
		channel:    channel,
		register:   make(chan *entity.Client),
		unregister: make(chan *entity.Client),
		broadcast:  make(chan *entity.Message),
		psr:        psr,
	}
}

func (cm *channelManager) Run(ctx context.Context) {
	go cm.subscribeToChannelMessages(ctx)

	for {
		select {
		case client := <-cm.register:
			cm.channel.RegisterClientInChannel(client)
		case client := <-cm.unregister:
			cm.channel.UnRegisterClientInChannel(client)
		case message := <-cm.broadcast:
			cm.publishChannelMessage(ctx, message)
		}
	}
}

func (cm *channelManager) broadcastToClientsInChannel(message []byte) {
	for cm := range cm.clientManagers {
		cm.send <- message
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
}

func (cm *channelManager) subscribeToChannelMessages(ctx context.Context) {
	pubsub := cm.psr.Subscribe(ctx, cm.channel.ID)
	defer pubsub.Close()

	msgs := pubsub.Channel()

	for msg := range msgs {
		cm.broadcastToClientsInChannel([]byte(msg.Payload))
	}
}
