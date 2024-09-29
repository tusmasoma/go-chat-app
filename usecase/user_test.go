package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/tusmasoma/go-chat-app/entity"
	"github.com/tusmasoma/go-chat-app/repository/mock"
)

func TestUserUseCase_SignUpAndGenerateToken(t *testing.T) { //nolint:gocognit // The number of lines is acceptable
	t.Helper()
	t.Setenv("WORKSPACE_ID", uuid.New().String())
	t.Setenv("PROFILE_IMAGE_URL", "https://example.com")

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserRepository,
			m1 *mock.MockMembershipRepository,
			m2 *mock.MockTransactionRepository,
			m3 *mock.MockAuthRepository,
		)
		arg struct {
			ctx      context.Context
			email    string
			password string
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(m *mock.MockUserRepository, m1 *mock.MockMembershipRepository, m2 *mock.MockTransactionRepository, m3 *mock.MockAuthRepository) {
				m2.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				m.EXPECT().LockUserByEmail(
					gomock.Any(),
					"test@gmail.com",
				).Return(false, nil)
				m.EXPECT().Create(
					gomock.Any(),
					gomock.Any(),
				).Do(func(_ context.Context, user entity.User) {
					if user.Email != "test@gmail.com" {
						t.Errorf("unexpected Email: got %v, want %v", user.Email, "test@gmail.com")
					}
					if user.Password == "" {
						t.Error("Password is empty")
					}
				}).Return(nil)
				m1.EXPECT().Create(
					gomock.Any(),
					gomock.Any(),
				).Do(func(_ context.Context, membership entity.Membership) {
					if membership.UserID == "" {
						t.Error("UserID is empty")
					}
					if membership.WorkspaceID == "" {
						t.Error("WorkspaceID is empty")
					}
					if membership.Name != "test" {
						t.Errorf("unexpected Name: got %v, want %v", membership.Name, "test")
					}
				}).Return(nil)
				m3.EXPECT().GenerateToken(
					gomock.Any(),
					"test@gmail.com",
				).Return("jwt", "jti")
			},
			arg: struct {
				ctx      context.Context
				email    string
				password string
			}{
				ctx:      context.Background(),
				email:    "test@gmail.com",
				password: "password123",
			},
			wantErr: nil,
		},
		{
			name: "Fail: Username already exists",
			setup: func(m *mock.MockUserRepository, m1 *mock.MockMembershipRepository, m2 *mock.MockTransactionRepository, m3 *mock.MockAuthRepository) {
				m2.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				m.EXPECT().LockUserByEmail(
					gomock.Any(),
					"test@gmail.com",
				).Return(true, nil)
			},
			arg: struct {
				ctx      context.Context
				email    string
				password string
			}{
				ctx:      context.Background(),
				email:    "test@gmail.com",
				password: "password123",
			},
			wantErr: errors.New("user with this email already exists"),
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ur := mock.NewMockUserRepository(ctrl)
			mr := mock.NewMockMembershipRepository(ctrl)
			tr := mock.NewMockTransactionRepository(ctrl)
			ar := mock.NewMockAuthRepository(ctrl)

			if tt.setup != nil {
				tt.setup(ur, mr, tr, ar)
			}

			usecase := NewUserUseCase(ur, mr, tr, ar)
			jwt, err := usecase.SignUpAndGenerateToken(tt.arg.ctx, tt.arg.email, tt.arg.password)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("SignUpAndGenerateToken() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("SignUpAndGenerateToken() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr == nil && jwt == "" {
				t.Error("Failed to generate token")
			}
		})
	}
}

func TestUserUseCase_LoginAndGenerateToken(t *testing.T) {
	t.Parallel()

	userID := uuid.New().String()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserRepository,
			m1 *mock.MockMembershipRepository,
			m2 *mock.MockTransactionRepository,
			m3 *mock.MockAuthRepository,
		)
		arg struct {
			ctx      context.Context
			email    string
			passward string
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(
				m *mock.MockUserRepository,
				m1 *mock.MockMembershipRepository,
				m2 *mock.MockTransactionRepository,
				m3 *mock.MockAuthRepository,
			) {
				hashPassword, _ := entity.PasswordEncrypt("password123")
				m.EXPECT().GetByEmail(
					gomock.Any(),
					"test@gmail.com",
				).Return(
					&entity.User{
						ID:       userID,
						Email:    "test@gmail.com",
						Password: hashPassword,
					}, nil,
				)
				m3.EXPECT().GenerateToken(userID, "test@gmail.com").Return(
					"jwt", "jti",
				)
			},
			arg: struct {
				ctx      context.Context
				email    string
				passward string
			}{
				ctx:      context.Background(),
				email:    "test@gmail.com",
				passward: "password123",
			},
			wantErr: nil,
		},
		{
			name: "Fail: invalid passward",
			setup: func(
				m *mock.MockUserRepository,
				m1 *mock.MockMembershipRepository,
				m2 *mock.MockTransactionRepository,
				m3 *mock.MockAuthRepository,
			) {
				hashPassword, _ := entity.PasswordEncrypt("invalidPassword123")
				m.EXPECT().GetByEmail(
					gomock.Any(),
					"test@gmail.com",
				).Return(
					&entity.User{
						ID:       userID,
						Email:    "test@gmail.com",
						Password: hashPassword,
					}, nil,
				)
			},
			arg: struct {
				ctx      context.Context
				email    string
				passward string
			}{
				ctx:      context.Background(),
				email:    "test@gmail.com",
				passward: "password123",
			},
			wantErr: fmt.Errorf("crypto/bcrypt: hashedPassword is not the hash of the given password"),
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ur := mock.NewMockUserRepository(ctrl)
			mr := mock.NewMockMembershipRepository(ctrl)
			tr := mock.NewMockTransactionRepository(ctrl)
			ar := mock.NewMockAuthRepository(ctrl)

			if tt.setup != nil {
				tt.setup(ur, mr, tr, ar)
			}

			usecase := NewUserUseCase(ur, mr, tr, ar)
			jwt, err := usecase.LoginAndGenerateToken(tt.arg.ctx, tt.arg.email, tt.arg.passward)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("LoginAndGenerateToken() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("LoginAndGenerateToken() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr == nil && jwt == "" {
				t.Error("Failed to generate token")
			}
		})
	}
}
