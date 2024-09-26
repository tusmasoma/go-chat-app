package mysql

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"

	"github.com/tusmasoma/go-chat-app/entity"
)

func Test_MembershipRepository(t *testing.T) {
	ctx := context.Background()
	repo := NewMembershipRepository(db)

	workspaceID := uuid.New().String()
	userID := uuid.New().String()

	membership, _ := entity.NewMembership(
		userID,
		workspaceID,
		"test",
		"https://www.hoge.com/avatar.jpg",
		false,
	)

	// Create
	err := repo.Create(ctx, *membership)
	ValidateErr(t, err, nil)

	// Get
	gotMembership, err := repo.Get(ctx, userID, workspaceID)
	ValidateErr(t, err, nil)
	if !reflect.DeepEqual(gotMembership, membership) {
		t.Errorf("Get() got = %v, want %v", gotMembership, membership)
	}

	// Update
	membership.Name = "updated"
	err = repo.Update(ctx, *membership)
	ValidateErr(t, err, nil)

	gotMembership, err = repo.Get(ctx, userID, workspaceID)
	ValidateErr(t, err, nil)
	if gotMembership.Name != "updated" {
		t.Errorf("Expected membership name 'updated', got %s", gotMembership.Name)
	}

	// Delete
	err = repo.Delete(ctx, userID, workspaceID)
	ValidateErr(t, err, nil)

	_, err = repo.Get(ctx, userID, workspaceID)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
