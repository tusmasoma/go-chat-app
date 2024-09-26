package entity

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestEntity_NewUser(t *testing.T) {
	t.Parallel()

	userID := uuid.New().String()
	hashPassword, _ := PasswordEncrypt("password123")
	patterns := []struct {
		name string
		arg  struct {
			id       string
			email    string
			password string
		}
		wantErr error
	}{
		{
			name: "Success: id is not empty",
			arg: struct {
				id       string
				email    string
				password string
			}{
				id:       userID,
				email:    "test@gmail.com",
				password: hashPassword,
			},
			wantErr: nil,
		},
		{
			name: "Success: id is empty",
			arg: struct {
				id       string
				email    string
				password string
			}{
				id:       "",
				email:    "test@gmail.com",
				password: hashPassword,
			},
			wantErr: nil,
		},
		{
			name: "Fail: email is required",
			arg: struct {
				id       string
				email    string
				password string
			}{
				id:       userID,
				email:    "",
				password: hashPassword,
			},
			wantErr: fmt.Errorf("email is required"),
		},
		{
			name: "Fail: password is required",
			arg: struct {
				id       string
				email    string
				password string
			}{
				id:       userID,
				email:    "test@gmail.com",
				password: "",
			},
			wantErr: fmt.Errorf("password is required"),
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewUser(tt.arg.id, tt.arg.email, tt.arg.password)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEntity_User_passwordEncrypt(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		name string
		arg  struct {
			password string
		}
		expectError error
	}{
		{
			name:        "empty password",
			arg:         struct{ password string }{password: ""},
			expectError: nil,
		},
		{
			name:        "normal password",
			arg:         struct{ password string }{password: "password123"},
			expectError: nil,
		},
		{
			name:        "complex password",
			arg:         struct{ password string }{password: "p@$$w0rd!23"},
			expectError: nil,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			hash, err := PasswordEncrypt(tt.arg.password)
			require.ErrorIs(t, err, tt.expectError)
			if err == nil && hash == "" {
				t.Errorf("Expected non-empty hash, got empty")
			}
		})
	}
}

func TestEntity_User_CompareHashAndPassword(t *testing.T) {
	t.Parallel()

	userID := uuid.New().String()
	hashPassword, _ := PasswordEncrypt("password123")
	user, _ := NewUser(userID, "test@gmail.com", hashPassword)

	patterns := []struct {
		name string
		arg  struct {
			password string
		}
		expectError error
	}{
		{
			name:        "valid hash and password",
			arg:         struct{ password string }{password: "password123"},
			expectError: nil,
		},
		{
			name:        "invalid password",
			arg:         struct{ password string }{password: "invalidpassword"},
			expectError: bcrypt.ErrMismatchedHashAndPassword,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := user.CompareHashAndPassword(tt.arg.password)
			require.ErrorIs(t, err, tt.expectError)
		})
	}
}
