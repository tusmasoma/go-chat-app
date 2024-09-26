package mysql

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/tusmasoma/go-chat-app/entity"
	"github.com/tusmasoma/go-chat-app/repository"
)

type userModel struct {
	ID       string `gorm:"type:char(36);primaryKey"`
	Email    string `gorm:"column:email"`
	Password string `gorm:"column:password"`
}

func (userModel) TableName() string {
	return "Users"
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) Get(ctx context.Context, id string) (*entity.User, error) {
	executor := ur.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	var um userModel
	if err := executor.WithContext(ctx).First(&um, "id = ?", id).Error; err != nil {
		return nil, err
	}

	user, err := entity.NewUser(um.ID, um.Email, um.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *userRepository) Create(ctx context.Context, user entity.User) error {
	executor := ur.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	if err := executor.WithContext(ctx).Create(&userModel{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) Update(ctx context.Context, user entity.User) error {
	executor := ur.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	if err := executor.WithContext(ctx).Model(&userModel{}).Where("id = ?", user.ID).Updates(&userModel{
		Password: user.Password,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) Delete(ctx context.Context, id string) error {
	executor := ur.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	if err := executor.WithContext(ctx).Delete(&userModel{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) LockUserByEmail(ctx context.Context, email string) (bool, error) {
	executor := ur.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	if err := executor.WithContext(ctx).Set("gorm:query_option", "FOR UPDATE").Where("email = ?", email).First(&userModel{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
