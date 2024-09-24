package handler

import (
	"context"

	"github.com/tusmasoma/go-chat-app/entity"
	"github.com/tusmasoma/go-chat-app/repository"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type channelHandler struct {
	channel *entity.Channel
	psr     repository.PubSubRepository
}

func NewChannelHandler(channel *entity.Channel, psr repository.PubSubRepository) *channelHandler {
	return &channelHandler{
		channel: channel,
		psr:     psr,
	}
}

func (ch *channelHandler) Run(ctx context.Context) {
	go ch.subscribeToChannelMessages(ctx)

	for {
		select {
		case client := <-ch.channel.Register:
			ch.channel.RegisterClientInChannel(client)
		case client := <-ch.channel.UnRegister:
			ch.channel.UnRegisterClientInChannel(client)
		case message := <-ch.channel.Broadcast:
			ch.publishChannelMessage(ctx, message)
		}
	}
}

func (ch *channelHandler) publishChannelMessage(ctx context.Context, message *entity.Message) {
	msg, err := message.Encode()
	if err != nil {
		log.Error("Failed to encode message", log.Ferror(err))
		return
	}
	if err := ch.psr.Publish(ctx, ch.channel.ID, msg); err != nil {
		log.Error("Failed to publish message", log.Ferror(err))
	}
}

func (ch *channelHandler) subscribeToChannelMessages(ctx context.Context) {
	pubsub := ch.psr.Subscribe(ctx, ch.channel.ID)
	defer pubsub.Close()

	msgs := pubsub.Channel()

	for msg := range msgs {
		ch.channel.BroadcastToClientsInChannel([]byte(msg.Payload))
	}
}
