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
	Login     string    `db:"login"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (c *connection) CreateUser(ctx context.Context, id uuid.UUID, hashPassword string, logins ...string,
) (out *User, err error) {
	err = c.ExecuteInTransaction(ctx, func(tx *sqlx.Tx) error {
		query := `
INSERT INTO "Users"(%s password) 
VALUES (%s :password)
RETURNING id, password, created_at, updated_at;`
		var idField, idValue string
		if !uuid.Equal(id, uuid.Nil) {
			idField = "id,"
			idValue = ":id,"
		}
		query, args, err := tx.BindNamed(fmt.Sprintf(query, idField, idValue),
			map[string]interface{}{
				"id":       id,
				"password": hashPassword,
			})
		if err != nil {
			return errors.WithStack(err)
		}

		out = new(User)
		err = tx.GetContext(ctx, out, query, args...)
		if err != nil {
			return errors.WithStack(err)
		}

		stmt, err := tx.PrepareNamedContext(ctx, `
INSERT INTO "Logins" (login, user_id)
VALUES (:login, :user_id);`)
		if err != nil {
			return errors.WithStack(err)
		}

		for _, login := range logins {
			_, err = stmt.ExecContext(ctx, map[string]interface{}{
				"user_id": out.ID,
				"login":   login,
			})
			if err != nil {
				return errors.WithStack(err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, c.Wrap(err)
	}

	return
}

func (c *connection) DeleteUser(ctx context.Context, userID uuid.UUID) (err error) {
	err = c.ExecuteInTransaction(ctx, func(tx *sqlx.Tx) error {
		query, args, err := tx.BindNamed(`
DELETE FROM "Users"
WHERE id = :id;`,
			map[string]interface{}{
				"id": userID,
			})

		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			return errors.WithStack(err)
		}

		return nil
	})
	if err != nil {
		return c.Wrap(err)
	}

	return
}

func (c *connection) GetUserByLogin(ctx context.Context, login string) (out *User, err error) {
	query, args, err := c.BindNamed(`
SELECT id,
       password,
       created_at,
       updated_at
FROM "Users"
WHERE id = (
    SELECT user_id 
    FROM "Logins" 
    WHERE login = :login
    LIMIT 1)
LIMIT 1;`,
		map[string]interface{}{
			"login": login,
		})
	if err != nil {
		return nil, c.Wrap(errors.WithStack(err))
	}

	out = &User{
		Login: login,
	}
	err = c.GetContext(ctx, out, query, args...)
	if err != nil {
		return nil, c.Wrap(errors.WithStack(err))
	}

	return
}
