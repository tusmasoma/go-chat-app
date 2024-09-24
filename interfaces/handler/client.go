package handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tusmasoma/go-chat-app/entity"
	"github.com/tusmasoma/go-chat-app/usecase"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

const (
	// Max wait time when writing a message to the peer.
	WriteWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	PongWait       = 60 * time.Second
	PingMultiplier = 9

	// Send pings to peer with this period. Must be less than pongWait.
	PingPeriod = (PongWait * PingMultiplier) / 10

	// Max message size allowed from peer.
	MaxMessageSize = 10000

	// Max buffer size for messages.
	BufferSize = 4096

	// ChannelBufferSize is the buffer size for the channel.
	ChannelBufferSize = 256

	// PubSubGeneralChannel is the general channel for pubsub.
	PubSubGeneralChannel = "general"

	// PubSubChannelPrefix is the prefix for the channel channel.
	WelcomeMessage = "%s joined the channel"
	GoodbyeMessage = "%s left the channel"
)

var newline = []byte{'\n'}

type ClientManager interface{}

type clientManager struct {
	client *entity.Client
	conn   *websocket.Conn
	hm     *hubManager
	send   chan []byte
	muc    usecase.MessageUseCase
}

func NewclientManager(client *entity.Client, conn *websocket.Conn, muc usecase.MessageUseCase) ClientManager {
	return &clientManager{
		client: client,
		conn:   conn,
		send:   make(chan []byte, BufferSize),
		muc:    muc,
	}
}

func (cm *clientManager) ReadPump() {
	defer func() {
		cm.disconnect()
	}()

	cm.conn.SetReadLimit(MaxMessageSize)
	if err := cm.conn.SetReadDeadline(time.Now().Add(PongWait)); err != nil {
		log.Error("Failed to set read deadline", log.Ferror(err))
	}
	cm.conn.SetPongHandler(func(string) error {
		err := cm.conn.SetReadDeadline(time.Now().Add(PongWait))
		if err != nil {
			log.Error("Error setting read deadline", log.Ferror(err))
			return err
		}
		return nil
	})
	// Start endless read loop, waiting for messages from client
	for {
		_, jsonMessage, err := cm.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Warn("Unexpected close error", log.Ferror(err))
			} else {
				log.Info("Client disconnected", log.Ferror(err))
			}
			break
		}

		cm.handleNewMessage(jsonMessage)
	}
}

func (cm *clientManager) WritePump() { //nolint: gocognit
	ticker := time.NewTicker(PingPeriod)
	defer func() {
		ticker.Stop()
		cm.conn.Close()
	}()

	for {
		select {
		case message, ok := <-cm.send:
			if err := cm.conn.SetWriteDeadline(time.Now().Add(WriteWait)); err != nil {
				log.Error("Failed to set write deadline", log.Ferror(err))
				return
			}
			if !ok {
				// The Hub closed the channel.
				if err := cm.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Warn("Failed to write close message", log.Ferror(err))
				}
				return
			}

			w, err := cm.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Error("Failed to get next writer", log.Ferror(err))
				return
			}

			if _, err = w.Write(message); err != nil {
				log.Error("Failed to write message", log.Ferror(err))
				return
			}

			// Attach queued chat messages to the current websocket message.
			n := len(cm.send)
			for i := 0; i < n; i++ {
				if _, err = w.Write(newline); err != nil {
					log.Error("Failed to write newline", log.Ferror(err))
					return
				}
				if _, err = w.Write(<-cm.send); err != nil {
					log.Error("Failed to write queued message", log.Ferror(err))
					return
				}
			}

			if err = w.Close(); err != nil {
				log.Error("Failed to close writer", log.Ferror(err))
				return
			}
		case <-ticker.C:
			if err := cm.conn.SetWriteDeadline(time.Now().Add(WriteWait)); err != nil {
				log.Error("Failed to set write deadline", log.Ferror(err))
				return
			}
			if err := cm.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Error("Failed to write ping message", log.Ferror(err))
				return
			}
		}
	}
}

func (cm *clientManager) disconnect() {
	cm.hm.unregister <- cm.client
	close(cm.send)
	if err := cm.conn.Close(); err != nil {
		log.Warn("Failed to close connection", log.Ferror(err))
	} else {
		log.Info("Client disconnected successfully", log.Fstring("clientID", cm.client.ID))
	}
}

func (cm *clientManager) handleNewMessage(jsonMessage []byte) {
	ctx := context.Background()

	var message entity.Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Error("Error unmarshalling JSON message", log.Ferror(err))
		return
	}

	message.SenderID = cm.client.UserID

	cm.routeMessageAction(ctx, message)
}

func (cm *clientManager) routeMessageAction(ctx context.Context, message entity.Message) {
	switch message.Action {
	case entity.CreateMessageAction:
		cm.muc.CreateMessage(ctx, &message)
		cm.broadcastMessage(message.TargetID, &message)
	case entity.UpdateMessageAction:
		cm.muc.UpdateMessage(ctx, &message)
		cm.broadcastMessage(message.TargetID, &message)
	case entity.DeleteMessageAction:
		cm.muc.DeleteMessage(ctx, &message)
		cm.broadcastMessage(message.TargetID, &message)
	default:
		log.Warn("Unknown message action", log.Fstring("action", message.Action))
	}
}

func (cm *clientManager) broadcastMessage(channelID string, message *entity.Message) {
	if channel := cm.hm.findChannelManagerByChannelID(channelID); channel != nil {
		log.Info("Broadcasting message", log.Fstring("channelID", channelID), log.Fstring("messageID", message.ID))
		channel.broadcast <- message
	} else {
		log.Warn("Channel not found", log.Fstring("channelID", channelID))
	}
}
