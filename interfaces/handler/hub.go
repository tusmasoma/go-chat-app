package handler

import (
	"context"

	"github.com/tusmasoma/go-chat-app/entity"
)

type hubHandler struct {
	hub *entity.Hub
}

func NewHubHandler(hub *entity.Hub) *hubHandler {
	return &hubHandler{
		hub: hub,
	}
}

func (h *hubHandler) Run(ctx context.Context) {
	for {
		select {
		case client := <-h.hub.Register:
			h.hub.RegisterClient(client)
		case client := <-h.hub.UnRegister:
			h.hub.UnRegisterClient(client)
		case message := <-h.hub.Broadcast:
			h.hub.BroadcastToClients(message)
		}
	}
}
