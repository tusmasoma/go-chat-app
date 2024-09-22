package handler

import (
	"context"

	"github.com/tusmasoma/go-chat-app/entity"
)

type channelHandler struct {
	channel *entity.Channel
}

func NewChannelHandler(channel *entity.Channel) *channelHandler {
	return &channelHandler{
		channel: channel,
	}
}

func (c *channelHandler) Run(ctx context.Context) {
	for {
		select {
		case client := <-c.channel.Register:
			c.channel.RegisterClientInChannel(client)
		case client := <-c.channel.UnRegister:
			c.channel.UnRegisterClientInChannel(client)
			// case message := <-c.channel.Broadcast:
			// 	c.channel.publishChannelMessage(ctx, message)
		}
	}
}
