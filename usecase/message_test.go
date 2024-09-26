package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/tusmasoma/go-chat-app/entity"
	"github.com/tusmasoma/go-chat-app/repository/mock"
)

func TestMessageUseCase_CreateMessage(t *testing.T) { //nolint:gocognit // ignore
	t.Parallel()

	date := time.Now()
	channelID := uuid.New().String()
	userID := uuid.New().String()
	message := entity.Message{
		UserID:    userID,
		Text:      "test message",
		CreatedAt: date,
		Action:    entity.CreateMessageAction,
		TargetID:  channelID,
	}

	patterns := []struct {
		name  string
		setup func(
			mmr *mock.MockMessageRepository,
		)
		arg struct {
			ctx     context.Context
			message *entity.Message
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(
				mmr *mock.MockMessageRepository,
			) {
				mmr.EXPECT().Create(
					gomock.Any(),
					gomock.Any(),
				).Do(func(_ context.Context, msg entity.Message) {
					if msg.UserID != userID {
						t.Errorf("unexpected UserID: got %v, want %v", msg.UserID, userID)
					}
					if msg.Text != "test message" {
						t.Errorf("unexpected Text: got %v, want %v", msg.Text, "test message")
					}
					if msg.CreatedAt != date {
						t.Errorf("unexpected CreatedAt: got %v, want %v", msg.CreatedAt, date)
					}
					if msg.Action != entity.CreateMessageAction {
						t.Errorf("unexpected Action: got %v, want %v", msg.Action, entity.CreateMessageAction)
					}
					if msg.TargetID != channelID {
						t.Errorf("unexpected TargetID: got %v, want %v", msg.TargetID, channelID)
					}
				}).Return(nil)
			},
			arg: struct {
				ctx     context.Context
				message *entity.Message
			}{
				ctx:     context.Background(),
				message: &message,
			},
			wantErr: nil,
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			mr := mock.NewMockMessageRepository(ctrl)

			if tt.setup != nil {
				tt.setup(mr)
			}

			usecase := NewMessageUseCase(mr)

			err := usecase.CreateMessage(
				tt.arg.ctx,
				tt.arg.message,
			)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("MessageCreate() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("MessageCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMessageUseCase_UpdateMessage(t *testing.T) {
	t.Parallel()

	msgID := uuid.New().String()
	channelID := uuid.New().String()
	userID := uuid.New().String()
	message := entity.Message{
		ID:       msgID,
		UserID:   userID,
		Text:     "updated message",
		Action:   entity.UpdateMessageAction,
		TargetID: channelID,
	}

	patterns := []struct {
		name  string
		setup func(
			mmr *mock.MockMessageRepository,
		)
		arg struct {
			ctx     context.Context
			message *entity.Message
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(
				mmr *mock.MockMessageRepository,
			) {
				mmr.EXPECT().Update(
					gomock.Any(),
					gomock.Any(),
				).Do(func(_ context.Context, msg entity.Message) {
					if msg.Text != "updated message" {
						t.Errorf("unexpected Text: got %v, want %v", msg.Text, "updated message")
					}
				}).Return(nil)
			},
			arg: struct {
				ctx     context.Context
				message *entity.Message
			}{
				ctx:     context.Background(),
				message: &message,
			},
			wantErr: nil,
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			mr := mock.NewMockMessageRepository(ctrl)

			if tt.setup != nil {
				tt.setup(mr)
			}

			usecase := NewMessageUseCase(mr)

			err := usecase.UpdateMessage(
				tt.arg.ctx,
				tt.arg.message,
			)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("MessageUpdate() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("MessageUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMessageUseCase_DeleteMessage(t *testing.T) {
	t.Parallel()

	msgID := uuid.New().String()
	channelID := uuid.New().String()
	userID := uuid.New().String()
	message := entity.Message{
		ID:       msgID,
		UserID:   userID,
		Text:     "test message",
		Action:   entity.DeleteMessageAction,
		TargetID: channelID,
	}

	patterns := []struct {
		name  string
		setup func(
			mmr *mock.MockMessageRepository,
		)
		arg struct {
			ctx     context.Context
			message *entity.Message
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(
				mmr *mock.MockMessageRepository,
			) {
				mmr.EXPECT().Delete(gomock.Any(), msgID).Return(nil)
			},
			arg: struct {
				ctx     context.Context
				message *entity.Message
			}{
				ctx:     context.Background(),
				message: &message,
			},
			wantErr: nil,
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			mr := mock.NewMockMessageRepository(ctrl)

			if tt.setup != nil {
				tt.setup(mr)
			}

			usecase := NewMessageUseCase(mr)

			err := usecase.DeleteMessage(
				tt.arg.ctx,
				tt.arg.message,
			)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("MessageDelete() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("MessageDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
