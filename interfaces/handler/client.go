package handler

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/tusmasoma/go-chat-app/entity"
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

type clientHandler struct {
	client *entity.Client
	conn   *websocket.Conn
}

func NewClientHandler(client *entity.Client, conn *websocket.Conn) *clientHandler {
	return &clientHandler{
		client: client,
		conn:   conn,
	}
}

func (ch *clientHandler) ReadPump() {
	// defer func() {
	// 	client.disconnect()
	// }()

	ch.conn.SetReadLimit(MaxMessageSize)
	if err := ch.conn.SetReadDeadline(time.Now().Add(PongWait)); err != nil {
		log.Error("Failed to set read deadline", log.Ferror(err))
	}
	ch.conn.SetPongHandler(func(string) error {
		err := ch.conn.SetReadDeadline(time.Now().Add(PongWait))
		if err != nil {
			log.Error("Error setting read deadline", log.Ferror(err))
			return err
		}
		return nil
	})
	// Start endless read loop, waiting for messages from client
	for {
		_, jsonMessage, err := ch.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Warn("Unexpected close error", log.Ferror(err))
			} else {
				log.Info("Client disconnected", log.Ferror(err))
			}
			break
		}

		ch.handleNewMessage(jsonMessage)
	}
}

func (ch *clientHandler) WritePump() { //nolint: gocognit
	ticker := time.NewTicker(PingPeriod)
	defer func() {
		ticker.Stop()
		ch.conn.Close()
	}()

	for {
		select {
		case message, ok := <-ch.client.Send:
			if err := ch.conn.SetWriteDeadline(time.Now().Add(WriteWait)); err != nil {
				log.Error("Failed to set write deadline", log.Ferror(err))
				return
			}
			if !ok {
				// The Hub closed the channel.
				if err := ch.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Warn("Failed to write close message", log.Ferror(err))
				}
				return
			}

			w, err := ch.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Error("Failed to get next writer", log.Ferror(err))
				return
			}

			if _, err = w.Write(message); err != nil {
				log.Error("Failed to write message", log.Ferror(err))
				return
			}

			// Attach queued chat messages to the current websocket message.
			n := len(ch.client.Send)
			for i := 0; i < n; i++ {
				if _, err = w.Write(newline); err != nil {
					log.Error("Failed to write newline", log.Ferror(err))
					return
				}
				if _, err = w.Write(<-ch.client.Send); err != nil {
					log.Error("Failed to write queued message", log.Ferror(err))
					return
				}
			}

			if err = w.Close(); err != nil {
				log.Error("Failed to close writer", log.Ferror(err))
				return
			}
		case <-ticker.C:
			if err := ch.conn.SetWriteDeadline(time.Now().Add(WriteWait)); err != nil {
				log.Error("Failed to set write deadline", log.Ferror(err))
				return
			}
			if err := ch.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Error("Failed to write ping message", log.Ferror(err))
				return
			}
		}
	}
}

func (ch *clientHandler) handleNewMessage(jsonMessage []byte) {
	// ctx := context.Background()

	// var message entity.WSMessage
	// if err := json.Unmarshal(jsonMessage, &message); err != nil {
	// 	log.Error("Error unmarshalling JSON message", log.Ferror(err))
	// 	return
	// }

	// message.SenderID = client.ID

	// switch message.Action {
	// case entity.ListMessagesAction:
	// 	client.handleListMessages(ctx, message)
	// case entity.CreateMessageAction:
	// 	client.handleCreateMessage(ctx, message)
	// case entity.DeleteMessageAction:
	// 	client.handleDeleteMessage(ctx, message)
	// case entity.UpdateMessageAction:
	// 	client.handleUpdateMessage(ctx, message)
	// case entity.CreatePublicChannelAction:
	// 	client.handleCreatePublicChannel(ctx, message)
	// case entity.JoinPublicChannelAction:
	// 	client.handleJoinPublicChannel(ctx, message)
	// case entity.LeavePublicChannelAction:
	// 	client.handleLeavePublicChannel(ctx, message)
	// default:
	// 	log.Warn("Unknown message action", log.Fstring("action", message.Action))
	// }
}
