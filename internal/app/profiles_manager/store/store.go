package store

import (
	"go.uber.org/zap"

	"typerium/internal/app/profiles_manager/store/migrations"
	"typerium/internal/pkg/database"
)

//go:generate go run github.com/kevinburke/go-bindata/go-bindata -o ./migrations/migrations.bindata.go -pkg migrations -ignore=\\*.go ./migrations/...

type Store interface {
	database.Connector
}

func New(uri string, version uint, log *zap.Logger) Store {
	cfg := &database.Config{
		URI:              uri,
		AssetDirFunc:     migrations.AssetDir,
		MigrationsDir:    "migrations",
		AssetFunc:        migrations.Asset,
		MigrationVersion: version,
	}
	conn := &connection{
		database.NewConnection(cfg, log),
	}

	return conn
}

type connection struct {
	*database.Connection
}
