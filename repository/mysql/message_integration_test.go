package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"

	"github.com/tusmasoma/go-chat-app/entity"
)

func Test_MessageRepository(t *testing.T) {
	ctx := context.Background()

	repo := NewMessageRepository(db)

	userID := uuid.New().String()
	workspaceID := uuid.New().String()
	channelID := uuid.New().String()

	msg1, err := entity.NewMessage(
		uuid.New().String(),
		userID,
		workspaceID,
		"Hello, World!",
		entity.CreateMessageAction,
		channelID,
		time.Time{},
	)
	ValidateErr(t, err, nil)

	msg2, err := entity.NewMessage(
		uuid.New().String(),
		userID,
		workspaceID,
		"Hello, World! 2",
		entity.CreateMessageAction,
		channelID,
		time.Time{},
	)
	ValidateErr(t, err, nil)

	// Create
	err = repo.Create(ctx, *msg1)
	ValidateErr(t, err, nil)
	err = repo.Create(ctx, *msg2)
	ValidateErr(t, err, nil)

	// Get
	gotMsg, err := repo.Get(ctx, msg1.ID)
	ValidateErr(t, err, nil)
	if d := cmp.Diff(msg1, gotMsg, cmpopts.IgnoreFields(entity.Message{}, "Action", "CreatedAt")); len(d) != 0 {
		t.Errorf("differs: (-want +got)\n%s", d)
	}

	// List
	msgs, err := repo.List(ctx, channelID)
	ValidateErr(t, err, nil)
	if len(msgs.Messages) != 2 {
		t.Errorf("len(msgs.Messages) got: %d, want: 2", len(msgs.Messages))
	}

	// Update
	msg1.Text = "Hello, World! Updated"
	err = repo.Update(ctx, *msg1)
	ValidateErr(t, err, nil)

	gotMsg, err = repo.Get(ctx, msg1.ID)
	ValidateErr(t, err, nil)
	if d := cmp.Diff(gotMsg.Text, "Hello, World! Updated"); len(d) != 0 {
		t.Errorf("differs: (-want +got)\n%s", d)
	}

	// Delete
	err = repo.Delete(ctx, msg1.ID)
	ValidateErr(t, err, nil)

	_, err = repo.Get(ctx, msg1.ID)
	if err == nil {
		t.Error("want error, but got nil")
	}
}
