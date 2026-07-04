package transaction

import (
	"context"
)

type Manager interface {
	// Execute runs a function within a transaction block
	Execute(ctx context.Context, fn func(ctx context.Context) error) error
}
