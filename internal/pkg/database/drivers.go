package database

import (
	"net/url"
	"strings"

	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func getDriver(uri string) (db *sqlx.DB, driverName string, driver database.Driver, err error) {
	parsedURL, err := url.Parse(uri)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	switch strings.ToLower(parsedURL.Scheme) {
	case "postgresql":
		driverName = "postgres"
		db, err = sqlx.Connect(driverName, uri)
		if err != nil {
			err = errors.WithStack(err)
			return
		}
		driver, err = postgres.WithInstance(db.DB, new(postgres.Config))
		if err != nil {
			err = errors.WithStack(err)
			return
		}
	default:
		err = errors.New("not support database")
	}

	return
}
