//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/go-chat-app/entity"
)

type MessageRepository interface {
	List(ctx context.Context, channleID string) (*entity.Messages, error)
	Get(ctx context.Context, id string) (*entity.Message, error)
	Create(ctx context.Context, message entity.Message) error
	Update(ctx context.Context, message entity.Message) error
	Delete(ctx context.Context, id string) error
}
