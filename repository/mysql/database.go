package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/go-chat-app/config"
	"github.com/tusmasoma/go-chat-app/repository"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) repository.TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (tr *transactionRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	var err error

	tx := tr.db.WithContext(ctx).Begin(&sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if tx.Error != nil {
		return tx.Error
	}

	ctx = context.WithValue(ctx, CtxTxKey(), tx)

	defer func() {
		if p := recover(); p != nil || err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error("Failed to rollback transaction: %v", rollbackErr)
			}
		}
	}()

	if err = fn(ctx); err != nil {
		return err
	}

	if err = tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

type TxKey string

func CtxTxKey() TxKey {
	return "tx"
}

func TxFromCtx(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(CtxTxKey()).(*gorm.DB)
	if !ok {
		return nil
	}
	return tx
}

func NewMySQLDB(ctx context.Context) (*gorm.DB, error) {
	conf, err := config.NewDBConfig(ctx)
	if err != nil {
		log.Error("Failed to load database config", log.Ferror(err))
		return nil, err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		conf.User, conf.Password, conf.Host, conf.Port, conf.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}) // ping is automatically called
	if err != nil {
		return nil, err
	}

	return db, nil
}
