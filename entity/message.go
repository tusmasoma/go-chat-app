package entity

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

// type MessageAction string

const (
	GetMessagesAction         = "GET_MESSAGES"
	ListMessagesAction        = "LIST_MESSAGES"
	CreateMessageAction       = "CREATE_MESSAGE"
	DeleteMessageAction       = "DELETE_MESSAGE"
	UpdateMessageAction       = "UPDATE_MESSAGE"
	CreatePublicChannelAction = "CREATE_PUBLIC_CHANNEL"
	JoinPublicChannelAction   = "JOIN_PUBLIC_CHANNEL"
	LeavePublicChannelAction  = "LEAVE_PUBLIC_CHANNEL"
	NoneAction                = "NONE"
)

var validActions = map[string]bool{
	GetMessagesAction:         true,
	ListMessagesAction:        true,
	CreateMessageAction:       true,
	DeleteMessageAction:       true,
	UpdateMessageAction:       true,
	CreatePublicChannelAction: true,
	JoinPublicChannelAction:   true,
	LeavePublicChannelAction:  true,
	NoneAction:                true,
}

type Message struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	Action    string    `json:"action"`
	TargetID  string    `json:"target_id"` // TargetID is the ID of the channel or user the message is intended for
	// SenderID  string    `json:"sender_id"` // SenderID is the ID of the user who sent the message
}

func NewMessage(id, userID, text, action, targetID string, createdAt time.Time) (*Message, error) {
	if id == "" {
		id = uuid.New().String()
	}
	if userID == "" {
		log.Error("userID is required")
		return nil, errors.New("userID is required")
	}
	if text == "" {
		// TODO: textの長さ制限を設ける
		log.Error("text is required")
		return nil, errors.New("text is required")
	}
	if !validActions[action] {
		log.Error("invalid action")
		return nil, errors.New("invalid action")
	}
	if targetID == "" {
		log.Error("targetID is required")
		return nil, errors.New("targetID is required")
	}
	if createdAt.IsZero() {
		now := time.Now()
		createdAt = now
	}
	return &Message{
		ID:        id,
		UserID:    userID,
		Text:      text,
		CreatedAt: createdAt,
		Action:    action,
		TargetID:  targetID,
	}, nil
}

func (m *Message) Encode() ([]byte, error) {
	json, err := json.Marshal(m)
	if err != nil {
		log.Error("Failed to encode message", log.Ferror(err))
		return nil, err
	}
	return json, nil
}

type Messages struct {
	Messages []*Message `json:"messages"`
	Action   string     `json:"action"`
	TargetID string     `json:"target_id"` // TargetID is the ID of the channel or user the message is intended for
	// SenderID  string    `json:"sender_id"` // SenderID is the ID of the user who sent the message
}

func NewMessages(messages []*Message, action, targetID string) (*Messages, error) {
	if len(messages) == 0 {
		log.Error("messages is required")
		return nil, errors.New("messages is required")
	}
	if !validActions[action] {
		log.Error("invalid action")
		return nil, errors.New("invalid action")
	}
	if targetID == "" {
		log.Error("targetID is required")
		return nil, errors.New("targetID is required")
	}
	return &Messages{
		Messages: messages,
		Action:   action,
		TargetID: targetID,
	}, nil
}

func (ms *Messages) Encode() ([]byte, error) {
	json, err := json.Marshal(ms)
	if err != nil {
		log.Error("Failed to encode messages", log.Ferror(err))
		return nil, err
	}
	return json, nil
}
