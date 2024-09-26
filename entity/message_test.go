package entity

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestEntity_NewMessage(t *testing.T) {
	t.Parallel()

	msgID := uuid.New().String()
	userID := uuid.New().String()
	channelID := uuid.New().String()

	patterns := []struct {
		name string
		arg  struct {
			id        string
			userID    string
			text      string
			action    string
			targetID  string
			createdAt time.Time
		}
		want struct {
			message *Message
			err     error
		}
	}{
		{
			name: "Success: id is not empty",
			arg: struct {
				id        string
				userID    string
				text      string
				action    string
				targetID  string
				createdAt time.Time
			}{
				id:        msgID,
				userID:    userID,
				text:      "text",
				action:    CreateMessageAction,
				targetID:  channelID,
				createdAt: time.Time{},
			},
			want: struct {
				message *Message
				err     error
			}{
				message: &Message{
					ID:       msgID,
					UserID:   userID,
					Text:     "text",
					Action:   CreateMessageAction,
					TargetID: channelID,
				},
				err: nil,
			},
		},
		{
			name: "Success: id is empty",
			arg: struct {
				id        string
				userID    string
				text      string
				action    string
				targetID  string
				createdAt time.Time
			}{
				id:        "",
				userID:    userID,
				text:      "text",
				action:    CreateMessageAction,
				targetID:  channelID,
				createdAt: time.Time{},
			},
			want: struct {
				message *Message
				err     error
			}{
				message: &Message{
					UserID:   userID,
					Text:     "text",
					Action:   CreateMessageAction,
					TargetID: channelID,
				},
				err: nil,
			},
		},
		{
			name: "Fail: userID is empty",
			arg: struct {
				id        string
				userID    string
				text      string
				action    string
				targetID  string
				createdAt time.Time
			}{
				id:        "",
				userID:    "",
				text:      "text",
				action:    CreateMessageAction,
				targetID:  channelID,
				createdAt: time.Time{},
			},
			want: struct {
				message *Message
				err     error
			}{
				message: nil,
				err:     errors.New("userID is required"),
			},
		},
		{
			name: "Fail: text is empty",
			arg: struct {
				id        string
				userID    string
				text      string
				action    string
				targetID  string
				createdAt time.Time
			}{
				id:        msgID,
				userID:    userID,
				text:      "",
				action:    CreateMessageAction,
				targetID:  channelID,
				createdAt: time.Time{},
			},
			want: struct {
				message *Message
				err     error
			}{
				message: nil,
				err:     errors.New("text is required"),
			},
		},
		{
			name: "Fail: action is invalid",
			arg: struct {
				id        string
				userID    string
				text      string
				action    string
				targetID  string
				createdAt time.Time
			}{
				id:        msgID,
				userID:    userID,
				text:      "text",
				action:    "InvalidAction",
				targetID:  channelID,
				createdAt: time.Time{},
			},
			want: struct {
				message *Message
				err     error
			}{
				message: nil,
				err:     errors.New("invalid action"),
			},
		},
		{
			name: "Fail: targetID is empty",
			arg: struct {
				id        string
				userID    string
				text      string
				action    string
				targetID  string
				createdAt time.Time
			}{
				id:        msgID,
				userID:    userID,
				text:      "text",
				action:    CreateMessageAction,
				targetID:  "",
				createdAt: time.Time{},
			},
			want: struct {
				message *Message
				err     error
			}{
				message: nil,
				err:     errors.New("targetID is required"),
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			message, err := NewMessage(tt.arg.id, tt.arg.userID, tt.arg.text, tt.arg.action, tt.arg.targetID, tt.arg.createdAt)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("NewMessage() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("NewMessage() error = %v, wantErr %v", err, tt.want.err)
			}

			if d := cmp.Diff(message, tt.want.message, cmpopts.IgnoreFields(Message{}, "ID", "CreatedAt")); len(d) != 0 {
				t.Errorf("NewMessage() mismatch (-got +want):\n%s", d)
			}
		})
	}
}

func TestEntity_Message_Encode(t *testing.T) {
	t.Parallel()

	msgID := uuid.New().String()
	userID := uuid.New().String()
	channelID := uuid.New().String()

	message, _ := NewMessage(
		msgID,
		userID,
		"text",
		CreateMessageAction,
		channelID,
		time.Time{},
	)

	encoded, err := message.Encode()
	if err != nil {
		t.Errorf("Encode() error = %v, want nil", err)
	}

	if encoded == nil {
		t.Errorf("Encode() got nil, want not nil")
	}
}

func TestEntity_NewMessages(t *testing.T) {
	t.Parallel()

	msgID := uuid.New().String()
	userID := uuid.New().String()
	channelID := uuid.New().String()

	message, _ := NewMessage(
		msgID,
		userID,
		"text",
		CreateMessageAction,
		channelID,
		time.Time{},
	)

	patterns := []struct {
		name string
		arg  struct {
			messages []*Message
			action   string
			targetID string
		}
		want struct {
			messages *Messages
			err      error
		}
	}{
		{
			name: "Success",
			arg: struct {
				messages []*Message
				action   string
				targetID string
			}{
				messages: []*Message{message},
				action:   CreateMessageAction,
				targetID: channelID,
			},
			want: struct {
				messages *Messages
				err      error
			}{
				messages: &Messages{
					Messages: []*Message{message},
					Action:   CreateMessageAction,
					TargetID: channelID,
				},
				err: nil,
			},
		},
		{
			name: "Fail: messages is empty",
			arg: struct {
				messages []*Message
				action   string
				targetID string
			}{
				messages: []*Message{},
				action:   CreateMessageAction,
				targetID: channelID,
			},
			want: struct {
				messages *Messages
				err      error
			}{
				messages: nil,
				err:      errors.New("messages is required"),
			},
		},
		{
			name: "Fail: action is invalid",
			arg: struct {
				messages []*Message
				action   string
				targetID string
			}{
				messages: []*Message{message},
				action:   "InvalidAction",
				targetID: channelID,
			},
			want: struct {
				messages *Messages
				err      error
			}{
				messages: nil,
				err:      errors.New("invalid action"),
			},
		},
		{
			name: "Fail: targetID is empty",
			arg: struct {
				messages []*Message
				action   string
				targetID string
			}{
				messages: []*Message{message},
				action:   CreateMessageAction,
				targetID: "",
			},
			want: struct {
				messages *Messages
				err      error
			}{
				messages: nil,
				err:      errors.New("targetID is required"),
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			messages, err := NewMessages(tt.arg.messages, tt.arg.action, tt.arg.targetID)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("NewMessages() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("NewMessages() error = %v, wantErr %v", err, tt.want.err)
			}

			if d := cmp.Diff(messages, tt.want.messages, cmpopts.IgnoreFields(Messages{}, "Messages")); len(d) != 0 {
				t.Errorf("NewMessages() mismatch (-got +want):\n%s", d)
			}
		})
	}
}

func TestEntity_Messages_Encode(t *testing.T) {
	t.Parallel()

	msgID := uuid.New().String()
	userID := uuid.New().String()
	channelID := uuid.New().String()

	message, _ := NewMessage(
		msgID,
		userID,
		"text",
		CreateMessageAction,
		channelID,
		time.Time{},
	)

	messages, _ := NewMessages([]*Message{message}, CreateMessageAction, channelID)

	encoded, err := messages.Encode()
	if err != nil {
		t.Errorf("Encode() error = %v, want nil", err)
	}

	if encoded == nil {
		t.Errorf("Encode() got nil, want not nil")
	}
}
