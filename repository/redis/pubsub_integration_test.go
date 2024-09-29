package redis

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test_PubSubRepository(t *testing.T) {
	repo := NewPubSubRepository(client)
	ctx := context.Background()

	channelID := uuid.New().String()
	message := []byte("testMessage")

	pubsub := repo.Subscribe(ctx, channelID)
	defer pubsub.Close()

	time.Sleep(5 * time.Second) // PublishとSubscribeの間に少し遅延を入れます

	err := repo.Publish(ctx, channelID, message)
	ValidateErr(t, err, nil)

	ch := pubsub.Channel()
	select {
	case msg := <-ch:
		if msg.Payload != string(message) {
			t.Errorf("Subscribe() \n got = %v,\n want = %v", msg.Payload, message)
		}
	case <-time.After(10 * time.Second):
		t.Error("Timeout waiting for message")
	}
}
