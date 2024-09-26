//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type PubSubRepository interface {
	Publish(ctx context.Context, channelID string, message []byte) error
	Subscribe(ctx context.Context, channelID string) *redis.PubSub // TODO: *redis.PubSub を汎用化する
}
