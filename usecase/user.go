//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/go-chat-app/entity"

	"github.com/tusmasoma/go-chat-app/repository"
)

type UserUseCase interface {
	CreateUserAndToken(ctx context.Context, email string, passward string) (string, error)
}

type userUseCase struct {
	ur repository.UserRepository
	mr repository.MembershipRepository
	tr repository.TransactionRepository
	ar repository.AuthRepository
}

func NewUserUseCase(
	ur repository.UserRepository,
	mr repository.MembershipRepository,
	tr repository.TransactionRepository,
	ar repository.AuthRepository,
) UserUseCase {
	return &userUseCase{
		ur: ur,
		mr: mr,
		tr: tr,
		ar: ar,
	}
}

// 現時点では、User作成時にMembershipも作成する(workspaceIDは固定されています)

func (uuc *userUseCase) CreateUserAndToken(ctx context.Context, email string, password string) (string, error) {
	var user *entity.User
	if err := uuc.tr.Transaction(ctx, func(ctx context.Context) error {
		exists, err := uuc.ur.LockUserByEmail(ctx, email)
		if err != nil {
			log.Error("Error retrieving user by email", log.Fstring("email", email))
			return err
		}
		if exists {
			log.Info("User with this email already exists", log.Fstring("email", email))
			return errors.New("user with this email already exists")
		}

		hashedPassword, err := entity.PasswordEncrypt(password)
		if err != nil {
			log.Error("Error encrypting password", log.Fstring("email", email))
			return err
		}
		user, err = entity.NewUser("", email, hashedPassword)
		if err != nil {
			log.Error("Error creating new user", log.Fstring("email", email))
			return err
		}

		if err = uuc.ur.Create(ctx, *user); err != nil {
			log.Error("Error creating new user", log.Fstring("email", email))
			return err
		}

		// Membership作成
		parts := strings.Split(email, "@")
		name := parts[0]
		membership, err := entity.NewMembership(
			user.ID,
			os.Getenv("WORKSPACE_ID"),
			name,
			os.Getenv("PROFILE_IMAGE_URL"),
			false,
		)
		if err != nil {
			log.Error("Error creating new membership", log.Fstring("email", email))
			return err
		}
		if err = uuc.mr.Create(ctx, *membership); err != nil {
			log.Error("Error creating new membership", log.Fstring("email", email))
			return err
		}

		return nil
	}); err != nil {
		return "", err
	}

	jwt, _ := uuc.ar.GenerateToken(user.ID, user.Email)
	return jwt, nil
}
