package mysql

import (
	"context"

	"gorm.io/gorm"

	"github.com/tusmasoma/go-chat-app/entity"
	"github.com/tusmasoma/go-chat-app/repository"
)

type membershipModel struct {
	UserID          string `gorm:"column:user_id"`
	WorkspaceID     string `gorm:"column:workspace_id"`
	Name            string `gorm:"column:name"`
	ProfileImageURL string `gorm:"column:profile_image_url"`
	IsAdmin         bool   `gorm:"column:is_admin"`
}

func (membershipModel) TableName() string {
	return "Memberships"
}

type membershipRepository struct {
	db *gorm.DB
}

func NewMembershipRepository(db *gorm.DB) repository.MembershipRepository {
	return &membershipRepository{
		db: db,
	}
}

func (mr *membershipRepository) Get(ctx context.Context, userID, workspaceID string) (*entity.Membership, error) {
	executor := mr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	var mm membershipModel
	if err := executor.WithContext(ctx).First(&mm, "user_id = ? AND workspace_id = ?", userID, workspaceID).Error; err != nil {
		return nil, err
	}

	membership, err := entity.NewMembership(
		mm.UserID,
		mm.WorkspaceID,
		mm.Name,
		mm.ProfileImageURL,
		mm.IsAdmin,
	)
	if err != nil {
		return nil, err
	}
	return membership, nil
}

func (mr *membershipRepository) Create(ctx context.Context, membership entity.Membership) error {
	executor := mr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	if err := executor.WithContext(ctx).Create(&membershipModel{
		UserID:          membership.UserID,
		WorkspaceID:     membership.WorkspaceID,
		Name:            membership.Name,
		ProfileImageURL: membership.ProfileImageURL,
		IsAdmin:         membership.IsAdmin,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (mr *membershipRepository) Update(ctx context.Context, membership entity.Membership) error {
	executor := mr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	if err := executor.WithContext(ctx).Model(
		&membershipModel{},
	).Where(
		"user_id = ? AND workspace_id = ?",
		membership.UserID,
		membership.WorkspaceID,
	).Updates(&membershipModel{
		Name:            membership.Name,
		ProfileImageURL: membership.ProfileImageURL,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (mr *membershipRepository) Delete(ctx context.Context, userID, workspaceID string) error {
	executor := mr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	if err := executor.WithContext(ctx).Delete(&membershipModel{}, "user_id = ? AND workspace_id = ?", userID, workspaceID).Error; err != nil {
		return err
	}
	return nil
}
