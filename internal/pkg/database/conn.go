package database

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"typerium/internal/pkg/logging"
)

// AssetDirFunc return names files in migrations directory
type AssetDirFunc func(name string) ([]string, error)

// Config configuration database
type Config struct {
	URI string
	// return names files in directory with migrations
	AssetDirFunc AssetDirFunc
	// path to directory with migrations
	MigrationsDir string
	// return migration by name
	AssetFunc bindata.AssetFunc
	// current version migration (default 0 is updated database to latest version)
	MigrationVersion uint
}

// Connector interface
type Connector interface {
	Ping() error
	PingContext(ctx context.Context) error
	Close()
}

// Connection instance of database access
type Connection struct {
	*sqlx.DB
	Logger   *zap.Logger
	wrappers []WrapperFunc
}

// NewConnection constructor for Connection
func NewConnection(cfg *Config) (conn *Connection) {
	conn = &Connection{
		Logger: logging.New("database"),
	}

	if cfg == nil {
		conn.Logger.Fatal("configuration isn't set")
	}

	var err error
	var driverName string
	var driver database.Driver

	conn.DB, driverName, driver, err = getDriver(cfg.URI)
	if err != nil {
		conn.Logger.Fatal("failed detect driver", zap.Error(err))
	}

	conn.Logger = conn.Logger.With(zap.String("driver", driverName))

	names, err := cfg.AssetDirFunc(cfg.MigrationsDir)
	if err != nil {
		conn.Logger.Fatal("can't get names migrations", zap.Error(err))
	}
	assetFunc := createAssetFunc(cfg.MigrationsDir, cfg.AssetFunc)
	sourceInstance, err := bindata.WithInstance(bindata.Resource(names, assetFunc))
	if err != nil {
		conn.Logger.Fatal("parsing migrations is failed", zap.Error(err))
	}

	m, err := migrate.NewWithInstance("go-bindata", sourceInstance, driverName, driver)
	if err != nil {
		conn.Logger.Fatal("failed create migrate instance", zap.Error(err))
	}
	if cfg.MigrationVersion == 0 {
		err = m.Up()
	} else {
		err = m.Migrate(cfg.MigrationVersion)
	}
	if err == migrate.ErrNoChange {
		err = nil
	}
	if err != nil {
		conn.Logger.Fatal("database migrating is failed", zap.Error(err))
	}
	version, dirty, err := m.Version()
	if err != nil {
		conn.Logger.Error("can't get migrations version", zap.Error(err))
		return
	}
	if dirty {
		conn.Logger.Warn("migrations is dirty", zap.Error(err))
		return
	}
	conn.Logger.Info(fmt.Sprintf("migrations is applied: current version %d", version))

	return
}

// Close db connection
func (c *Connection) Close() {
	err := c.DB.Close()
	if err != nil {
		c.Logger.Error("failed to close connection")
	}
}

func createAssetFunc(migrationsDir string, assetFunc bindata.AssetFunc) bindata.AssetFunc {
	return func(name string) (bytes []byte, e error) {
		data, err := assetFunc(fmt.Sprintf("%s/%s", migrationsDir, name))
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return data, nil
	}
}
