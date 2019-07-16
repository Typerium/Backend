package main

import (
	"github.com/spf13/viper"

	"typerium/internal/app/profiles_manager/store"
	"typerium/internal/pkg/logging"
	"typerium/internal/pkg/waiter"

	_ "typerium/internal/pkg/config"
)

const (
	dbURI     = "db_uri"
	dbVersion = "db_version"
)

func main() {
	log := logging.New()

	viper.SetDefault(dbURI, "postgresql://dev:123456@localhost:5432/profiles_manager?sslmode=disable")
	viper.SetDefault(dbVersion, uint(0))

	db := store.New(viper.GetString(dbURI), viper.GetUint(dbVersion), log)
	defer db.Close()

	waiter.Wait(log)
}
