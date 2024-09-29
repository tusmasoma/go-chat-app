package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/tusmasoma/go-chat-app/repository"
)

type pubsubRepository struct {
	client *redis.Client
}

func NewPubSubRepository(client *redis.Client) repository.PubSubRepository {
	return &pubsubRepository{
		client,
	}
}

func (r *pubsubRepository) Publish(ctx context.Context, channelID string, message []byte) error {
	return r.client.Publish(ctx, channelID, message).Err()
}

func (r *pubsubRepository) Subscribe(ctx context.Context, channelID string) *redis.PubSub {
	return r.client.Subscribe(ctx, channelID)
}
