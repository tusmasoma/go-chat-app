package repository

import "context"

type TransactionRepository interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
