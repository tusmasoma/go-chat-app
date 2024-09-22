package repository

import (
	"context"

	"github.com/tusmasoma/go-chat-app/entity"
)

type ChannelRepository interface {
	List(ctx context.Context) ([]entity.Channel, error)
	Get(ctx context.Context, id string) (*entity.Channel, error)
	Create(ctx context.Context, channel entity.Channel) error
	Update(ctx context.Context, channel entity.Channel) error
	Delete(ctx context.Context, id string) error
}
