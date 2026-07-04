package transaction

import (
	"context"

	"gorm.io/gorm"
)

const GormServiceName = "transaction_manager"

type gormManager struct {
	db *gorm.DB
}

func NewGormManager(db *gorm.DB) Manager {
	return &gormManager{db: db}
}

type contextKey string

const txKey contextKey = "gorm_transaction_context_key"

func (m *gormManager) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	// .Transaction() automatically handles Rollback/Commit based on error return
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Inject the transaction object into a NEW context
		txCtx := context.WithValue(ctx, txKey, tx)

		// 2. Execute the business logic provided by the caller
		// We pass txCtx so that repositories can find the 'tx'
		return fn(txCtx)
	})
}
