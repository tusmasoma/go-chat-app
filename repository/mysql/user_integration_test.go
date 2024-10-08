package mysql

import (
	"context"
	"reflect"
	"testing"

	"github.com/tusmasoma/go-chat-app/entity"
)

func Test_UserRepository(t *testing.T) {
	ctx := context.Background()

	repo := NewUserRepository(db)

	hashPassword, _ := entity.PasswordEncrypt("password")
	user, err := entity.NewUser(
		"",
		"test@gmail.com",
		hashPassword,
	)
	ValidateErr(t, err, nil)

	// Create
	err = repo.Create(ctx, *user)
	ValidateErr(t, err, nil)

	// Get
	gotUser, err := repo.Get(ctx, user.ID)
	ValidateErr(t, err, nil)
	if !reflect.DeepEqual(user, gotUser) {
		t.Errorf("want: %v, got: %v", user, gotUser)
	}

	// GetByEmail
	gotUserByEmail, err := repo.GetByEmail(ctx, "test@gmail.com")
	ValidateErr(t, err, nil)
	if !reflect.DeepEqual(user, gotUserByEmail) {
		t.Errorf("want: %v, got: %v", user, gotUserByEmail)
	}

	// LockUserByEmail
	exists, err := repo.LockByEmail(ctx, "test@gmail.com")
	ValidateErr(t, err, nil)
	if !exists {
		t.Fatalf("Failed to get user by email")
	}

	// Update
	gotUser.Password, _ = entity.PasswordEncrypt("newpassword")
	err = repo.Update(ctx, *gotUser)
	ValidateErr(t, err, nil)

	updatedUser, err := repo.Get(ctx, user.ID)
	ValidateErr(t, err, nil)
	if !reflect.DeepEqual(gotUser, updatedUser) {
		t.Errorf("want: %v, got: %v", gotUser, updatedUser)
	}

	// Delete
	err = repo.Delete(ctx, user.ID)
	ValidateErr(t, err, nil)

	_, err = repo.Get(ctx, user.ID)
	if err == nil {
		t.Errorf("want: %v, got: %v", nil, err)
	}
}
