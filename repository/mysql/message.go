package mysql

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/tusmasoma/go-chat-app/entity"
	"github.com/tusmasoma/go-chat-app/repository"
)

type messageModel struct {
	ID          string    `gorm:"type:char(36);primaryKey"`
	UserID      string    `gorm:"column:user_id"`
	WorkspaceID string    `gorm:"column:workspace_id"`
	ChannelID   string    `gorm:"column:channel_id"`
	Text        string    `gorm:"column:text"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

func (messageModel) TableName() string {
	return "Messages"
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) repository.MessageRepository {
	return &messageRepository{
		db: db,
	}
}

func (mr *messageRepository) List(ctx context.Context, channleID string) (*entity.Messages, error) {
	executor := mr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	var mms []messageModel
	if err := executor.WithContext(ctx).Find(&mms, "channel_id = ?", channleID).Error; err != nil {
		return nil, err
	}

	var err error
	msgs := make([]*entity.Message, len(mms))
	for i, mm := range mms {
		msgs[i], err = entity.NewMessage(
			mm.ID,
			mm.UserID,
			mm.WorkspaceID,
			mm.Text,
			entity.NoneAction,
			mm.ChannelID,
			mm.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
	}

	messages, err := entity.NewMessages(msgs, entity.ListMessagesAction, channleID)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (mr *messageRepository) Get(ctx context.Context, id string) (*entity.Message, error) {
	executor := mr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	var mm messageModel
	if err := executor.WithContext(ctx).First(&mm, "id = ?", id).Error; err != nil {
		return nil, err
	}

	msg, err := entity.NewMessage(
		mm.ID,
		mm.UserID,
		mm.WorkspaceID,
		mm.Text,
		entity.GetMessagesAction,
		mm.ChannelID,
		mm.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (mr *messageRepository) Create(ctx context.Context, message entity.Message) error {
	executor := mr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	if err := executor.WithContext(ctx).Create(&messageModel{
		ID:          message.ID,
		UserID:      message.UserID,
		WorkspaceID: message.WorkspaceID,
		ChannelID:   message.TargetID,
		Text:        message.Text,
		CreatedAt:   message.CreatedAt,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (mr *messageRepository) Update(ctx context.Context, message entity.Message) error {
	executor := mr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	if err := executor.WithContext(ctx).Model(&messageModel{}).Where("id = ?", message.ID).Updates(&messageModel{
		Text: message.Text,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (mr *messageRepository) Delete(ctx context.Context, id string) error {
	executor := mr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	if err := executor.WithContext(ctx).Delete(&messageModel{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
