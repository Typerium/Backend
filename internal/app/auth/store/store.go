package store

import (
	"context"
	"time"

	uuid "github.com/satori/go.uuid"

	"typerium/internal/app/auth/store/migrations"
	"typerium/internal/pkg/database"
)

//go:generate go run github.com/kevinburke/go-bindata/go-bindata -o ./migrations/migrations.bindata.go -pkg migrations -ignore=\\*.go ./migrations/...

// Store data access layer
type Store interface {
	database.Connector

	CreateUser(ctx context.Context, id uuid.UUID, hashPassword string, logins ...string) (out *User, err error)
	DeleteUser(ctx context.Context, userID uuid.UUID) (err error)
	GetUserByLogin(ctx context.Context, login string) (out *User, err error)

	CreateSession(ctx context.Context, userID uuid.UUID, keySignature []byte, ttl time.Duration) (out *Session,
		err error)
	DeleteSession(ctx context.Context, sessionID uuid.UUID) (err error)
	GetSessionByID(ctx context.Context, id uuid.UUID) (out *Session, err error)
	UpdateSession(ctx context.Context, id uuid.UUID, keySignature []byte, ttl time.Duration) (out *Session, err error)
}

// New create database connection
func New(uri string, version uint) Store {
	cfg := &database.Config{
		URI:              uri,
		AssetDirFunc:     migrations.AssetDir,
		MigrationsDir:    "migrations",
		AssetFunc:        migrations.Asset,
		MigrationVersion: version,
	}
	conn := &connection{
		database.NewConnection(cfg),
	}
	return conn
}

type connection struct {
	*database.Connection
}
