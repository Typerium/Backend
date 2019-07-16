package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Phone     *string   `db:"phone"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (c *connection) CreateUser(ctx context.Context, in *User) (out *User, err error) {
	err = c.ExecuteInTransaction(ctx, func(tx *sqlx.Tx) error {
		query := `
INSERT INTO "Profiles" (%s username, email, phone)
VALUES (%s :username, :email, :phone)
RETURNING id, created_at, updated_at;`
		var idField, idValue string
		if !uuid.Equal(in.ID, uuid.Nil) {
			idField = "id,"
			idValue = ":id,"
		}

		query, args, err := tx.BindNamed(fmt.Sprintf(query, idField, idValue), in)
		if err != nil {
			return errors.WithStack(err)
		}

		err = tx.GetContext(ctx, in, query, args...)
		return errors.WithStack(err)
	})
	if err != nil {
		return nil, c.Wrap(err)
	}

	out = in
	return
}

func (c *connection) DeleteUser(ctx context.Context, id uuid.UUID) (err error) {
	err = c.ExecuteInTransaction(ctx, func(tx *sqlx.Tx) error {
		query, args, err := tx.BindNamed(`
DELETE FROM "Profiles"
WHERE id = :id;`, map[string]interface{}{
			"id": id,
		})
		if err != nil {
			return errors.WithStack(err)
		}

		_, err = tx.ExecContext(ctx, query, args...)
		return errors.WithStack(err)
	})
	if err != nil {
		return c.Wrap(err)
	}

	return
}

func (c *connection) GetUserByID(ctx context.Context, id uuid.UUID) (out *User, err error) {
	out = &User{
		ID: id,
	}

	query, args, err := c.BindNamed(`
SELECT id, username, email, phone, created_at, updated_at
FROM "Profiles"
WHERE id = :id
LIMIT 1;`, out)
	if err != nil {
		return nil, c.Wrap(errors.WithStack(err))
	}

	err = c.GetContext(ctx, out, query, args...)
	if err != nil {
		return nil, c.Wrap(errors.WithStack(err))
	}

	return
}
