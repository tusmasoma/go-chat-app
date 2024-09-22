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
	ListMessagesAction        = "LIST_MESSAGES"
	CreateMessageAction       = "CREATE_MESSAGE"
	DeleteMessageAction       = "DELETE_MESSAGE"
	UpdateMessageAction       = "UPDATE_MESSAGE"
	CreatePublicChannelAction = "CREATE_PUBLIC_CHANNEL"
	JoinPublicChannelAction   = "JOIN_PUBLIC_CHANNEL"
	LeavePublicChannelAction  = "LEAVE_PUBLIC_CHANNEL"
)

var validActions = map[string]bool{
	ListMessagesAction:        true,
	CreateMessageAction:       true,
	DeleteMessageAction:       true,
	UpdateMessageAction:       true,
	CreatePublicChannelAction: true,
	JoinPublicChannelAction:   true,
	LeavePublicChannelAction:  true,
}

type Message struct {
	ID        string
	UserID    string
	Text      string
	CreatedAt time.Time
	Action    string
	TargetID  string // TargetID is the ID of the channel or user the message is intended for
	SenderID  string // SenderID is the ID of the user who sent the message
}

type Messages []*Message

func NewMessage(id, userID, text, action, targetID, senderID string, createdAt *time.Time) (*Message, error) {
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
	if senderID == "" {
		log.Error("senderID is required")
		return nil, errors.New("senderID is required")
	}
	if createdAt == nil {
		now := time.Now()
		createdAt = &now
	}
	return &Message{
		ID:        id,
		UserID:    userID,
		Text:      text,
		CreatedAt: *createdAt,
		Action:    action,
		TargetID:  targetID,
		SenderID:  senderID,
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
