package main

import (
	"github.com/spf13/viper"

	"typerium/internal/app/profiles_manager/handlers"
	"typerium/internal/app/profiles_manager/store"
	"typerium/internal/pkg/broker"
	"typerium/internal/pkg/broker/proto"
	"typerium/internal/pkg/waiter"

	_ "typerium/internal/pkg/config"
)

const (
	dbURI     = "DB_URI"
	dbVersion = "DB_VERSION"
	grpcAddr  = "GRPC_ADDR"
)

var grpcErrorsWrap = map[error]error{}

func main() {
	viper.SetDefault(dbURI, "postgresql://dev:123456@localhost:5432/profiles_manager?sslmode=disable")
	viper.SetDefault(dbVersion, uint(0))

	db := store.New(viper.GetString(dbURI), viper.GetUint(dbVersion))
	defer db.Close()

	grpcServer := broker.NewGRPCServer(grpcErrorsWrap)
	proto.RegisterProfilesManagerServiceServer(grpcServer.Server, handlers.NewGRPCServer(db))

	viper.SetDefault(grpcAddr, ":50051")
	grpcServer.ServeOnAddress(viper.GetString(grpcAddr))
	defer grpcServer.GracefulStop()

	waiter.Wait()
}
