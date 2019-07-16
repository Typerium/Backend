package database

import (
	"context"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// TransactionFunc decorator executing in transaction
type TransactionFunc func(tx *sqlx.Tx) error

// ExecuteInTransaction execute decorator in transaction with commit (success executing) or rollback (failed executing)
func (c *Connection) ExecuteInTransaction(ctx context.Context, f TransactionFunc) error {
	tx, err := c.BeginTxx(ctx, nil)
	if err != nil {
		c.Logger.Error("failed begin transaction", zap.Error(err))
		return ErrInternal
	}
	err = f(tx)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			c.Logger.Error("failed rollback transaction", zap.Error(errRollback))
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		c.Logger.Error("failed commit transaction", zap.Error(err))
		return ErrInternal
	}

	return nil
}
