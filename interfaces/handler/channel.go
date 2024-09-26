package handler

import (
	"encoding/json"
	"io"
	"net/http"

	ws "github.com/tusmasoma/go-chat-app/interfaces/websocket"
	"github.com/tusmasoma/go-chat-app/usecase"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type ChannelHandler interface {
	CreateChannel(w http.ResponseWriter, r *http.Request)
}

type channelHandler struct {
	hm  *ws.HubManager // 現状、Workspaceは一つの為、containerにてHubManagerを生成して、DIする
	cuc usecase.ChannelUseCase
}

func NewChannelHandler() ChannelHandler {
	return &channelHandler{}
}

type CreateChannelRequest struct {
	Name    string `json:"name"`
	Private bool   `json:"private"`
}

func (ch *channelHandler) CreateChannel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var requestBody CreateChannelRequest
	defer r.Body.Close()
	if !ch.isValidCreateChannelRequest(r.Body, &requestBody) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	channel, err := ch.cuc.CreateChannel(ctx, requestBody.Name, requestBody.Private)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ch.hm.RegisterChannel(ctx, channel)

	w.WriteHeader(http.StatusOK)
}

func (ch *channelHandler) isValidCreateChannelRequest(body io.ReadCloser, requestBody *CreateChannelRequest) bool {
	if err := json.NewDecoder(body).Decode(requestBody); err != nil {
		log.Error("Failed to decode request body: %v", err)
		return false
	}
	if requestBody.Name == "" {
		log.Warn("Invalid request body: %v", requestBody)
		return false
	}
	return true
}
