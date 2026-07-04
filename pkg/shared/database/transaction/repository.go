package transaction

import (
	"context"

	"gorm.io/gorm"
)

type AppRepository struct{}

func (*AppRepository) GetTx(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	// Attempt to retrieve the transaction from the context
	if tx, ok := ctx.Value(txKey).(*gorm.DB); ok {
		return tx
	}
	// Fallback to the non-transactional DB instance
	return defaultDB.WithContext(ctx)
}
