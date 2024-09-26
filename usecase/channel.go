package usecase

import (
	"context"

	"github.com/tusmasoma/go-chat-app/entity"
)

type ChannelUseCase interface {
	CreateChannel(ctx context.Context, name string, private bool) (*entity.Channel, error)
}
