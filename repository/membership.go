//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/go-chat-app/entity"
)

type MembershipRepository interface {
	Get(ctx context.Context, userID, workspaceID string) (*entity.Membership, error)
	Create(ctx context.Context, membership entity.Membership) error
	Update(ctx context.Context, membership entity.Membership) error
	Delete(ctx context.Context, userID, workspaceID string) error
}
