package store

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

type Session struct {
	ID           uuid.UUID `db:"id"`
	UserID       uuid.UUID `db:"user_id"`
	KeySignature []byte    `db:"key_signature"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func (c *connection) CreateSession(ctx context.Context, userID uuid.UUID, keySignature []byte, ttl time.Duration,
) (out *Session, err error) {
	err = c.ExecuteInTransaction(ctx, func(tx *sqlx.Tx) error {
		query, args, err := tx.BindNamed(`
INSERT INTO "Sessions"(user_id, key_signature, time_exp) 
VALUES (:user_id, :key_signature, :ttl)
RETURNING id, user_id, key_signature, created_at, updated_at;`,
			map[string]interface{}{
				"user_id":       userID,
				"key_signature": keySignature,
				"ttl":           time.Now().UTC().Add(ttl),
			})
		if err != nil {
			return errors.WithStack(err)
		}

		out = new(Session)
		err = tx.GetContext(ctx, out, query, args...)
		if err != nil {
			return errors.WithStack(err)
		}

		return nil
	})
	if err != nil {
		return nil, c.Wrap(err)
	}

	return
}

func (c *connection) DeleteSession(ctx context.Context, sessionID uuid.UUID) (err error) {
	err = c.ExecuteInTransaction(ctx, func(tx *sqlx.Tx) error {
		query, args, err := tx.BindNamed(`
DELETE FROM "Sessions"
WHERE id = :id;`,
			map[string]interface{}{
				"id": sessionID,
			})
		if err != nil {
			return errors.WithStack(err)
		}
		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
	return c.Wrap(err)
}

func (c *connection) GetSessionByID(ctx context.Context, id uuid.UUID) (out *Session, err error) {
	query, args, err := c.BindNamed(`
SELECT id,
       user_id,
       key_signature,
       created_at,
       updated_at
FROM "Sessions"
WHERE id = :id
LIMIT 1;`,
		map[string]interface{}{
			"id": id,
		})

	out = new(Session)
	err = c.GetContext(ctx, out, query, args...)
	if err != nil {
		return nil, c.Wrap(errors.WithStack(err))
	}

	return
}

func (c *connection) UpdateSession(ctx context.Context, id uuid.UUID, keySignature []byte, ttl time.Duration,
) (out *Session, err error) {
	err = c.ExecuteInTransaction(ctx, func(tx *sqlx.Tx) error {
		query, args, err := tx.BindNamed(`
UPDATE "Sessions" SET key_signature = :key_signature, time_exp = :ttl, updated_at = now()
WHERE id = :id
RETURNING id, user_id, key_signature, created_at, updated_at;`,
			map[string]interface{}{
				"key_signature": keySignature,
				"ttl":           ttl.Nanoseconds(),
			})
		if err != nil {
			return errors.WithStack(err)
		}

		out = new(Session)
		err = tx.GetContext(ctx, out, query, args...)
		if err != nil {
			return errors.WithStack(err)
		}

		return nil
	})
	if err != nil {
		return nil, c.Wrap(err)
	}

	return
}
