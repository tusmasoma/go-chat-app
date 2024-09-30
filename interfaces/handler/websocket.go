package handler

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/go-chat-app/config"
	"github.com/tusmasoma/go-chat-app/entity"
	ws "github.com/tusmasoma/go-chat-app/interfaces/websocket"
	"github.com/tusmasoma/go-chat-app/usecase"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  config.BufferSize,
	WriteBufferSize: config.BufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebsocketHandler struct {
	hm  *ws.HubManager // 現状、Workspaceは一つの為、containerにてHubManagerを生成して、DIする
	muc usecase.MessageUseCase
}

func NewWebsocketHandler(hm *ws.HubManager, muc usecase.MessageUseCase) *WebsocketHandler {
	return &WebsocketHandler{
		hm:  hm,
		muc: muc,
	}
}

func (wsh *WebsocketHandler) WebSocket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDValue := ctx.Value(config.ContextUserIDKey)
	userID, ok := userIDValue.(string)
	if !ok {
		log.Error("User ID not found in request context")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil) // conn is *websocket.Conn
	if err != nil {
		log.Error("Failed to upgrade connection", log.Ferror(err))
		return
	}

	client, err := entity.NewClient("", userID, wsh.hm.Hub)
	if err != nil {
		log.Error("Failed to create new client", log.Ferror(err))
		return
	}
	clientManager := ws.NewClientManager(client, conn, wsh.hm, wsh.muc)

	go clientManager.WritePump()
	go clientManager.ReadPump()

	wsh.hm.Register <- client

	// HubManagerに登録さているChannelにClientを登録
	wsh.hm.RegisterClientManagerInChannelManager(clientManager)

	log.Info(
		"Successfully Client connected",
		log.Fstring("userID", userID),
		log.Fstring("workspaceID", wsh.hm.Hub.ID),
	)
}
