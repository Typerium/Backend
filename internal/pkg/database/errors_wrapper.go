package database

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// WrapperFunc wrap error and return new error with ok=true or skip wrapping with ok=false
type WrapperFunc func(input error) (output error, ok bool)

// common errors
var (
	ErrUnknown  = errors.New("unprocessed error")
	ErrInternal = errors.New("internal error")
)

// Register add a new wrap function
func (c *Connection) Register(f WrapperFunc) {
	c.wrappers = append(c.wrappers, f)
}

// Wrap processing input error and swap error to another error if wrapping is success by WrapperFunc
func (c *Connection) Wrap(input error) error {
	cleanErr := errors.Cause(input)
	if cleanErr == nil || cleanErr == sql.ErrNoRows || cleanErr == ErrInternal {
		return cleanErr
	}

	for _, wrapper := range c.wrappers {
		out, ok := wrapper(cleanErr)
		if ok {
			return out
		}
	}

	c.Logger.Warn("database error isn't processed", zap.Error(input))
	return ErrUnknown
}

// PostgresConstraintErrorMapper mapper for wrapping error postgres database (key - constraint, value - swapping error)
type PostgresConstraintErrorMapper map[string]error

func postgresWrapperFunc(mapper PostgresConstraintErrorMapper, code pq.ErrorCode) WrapperFunc {
	return func(input error) (output error, ok bool) {
		pgErr, ok := input.(*pq.Error)
		if !ok || pgErr.Code != code {
			return input, false
		}

		output, ok = mapper[pgErr.Constraint]
		if !ok {
			return input, false
		}

		return
	}
}

// PostgresUniqueIndex create wrapping function for code 23505 (unique violation)
func PostgresUniqueIndex(mapper PostgresConstraintErrorMapper) WrapperFunc {
	return postgresWrapperFunc(mapper, "23505")
}

// PostgresFKViolation create wrapping function for code 23503 (foreign key violation)
func PostgresFKViolation(mapper PostgresConstraintErrorMapper) WrapperFunc {
	return postgresWrapperFunc(mapper, "23503")
}
