package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/go-chat-app/entity"
	"github.com/tusmasoma/go-chat-app/repository"
)

type MessageUseCase interface {
	CreateMessage(ctx context.Context, message *entity.Message) error
	UpdateMessage(ctx context.Context, message *entity.Message) error
	DeleteMessage(ctx context.Context, message *entity.Message) error
}

type messageUseCase struct {
	mr repository.MessageRepository
}

func NewMessageUseCase(mr repository.MessageRepository) MessageUseCase {
	return &messageUseCase{
		mr: mr,
	}
}

func (muc *messageUseCase) CreateMessage(ctx context.Context, message *entity.Message) error {
	message.ID = uuid.New().String() // TODO: messageの生成を再度する？
	if err := muc.mr.Create(ctx, *message); err != nil {
		log.Error("Failed to create message", log.Ferror(err))
		return err
	}
	return nil
}

func (muc *messageUseCase) UpdateMessage(ctx context.Context, message *entity.Message) error {
	// TODO: user認証 & messageの所有権確認
	if err := muc.mr.Update(ctx, *message); err != nil {
		log.Error("Failed to update message", log.Ferror(err))
		return err
	}
	return nil
}

func (muc *messageUseCase) DeleteMessage(ctx context.Context, message *entity.Message) error {
	// TODO: user認証 & messageの所有権確認
	if err := muc.mr.Delete(ctx, message.ID); err != nil {
		log.Error("Failed to delete message", log.Ferror(err))
		return err
	}
	return nil
}
