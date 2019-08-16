package main

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"typerium/internal/app/auth/handlers"
	"typerium/internal/app/auth/password"
	"typerium/internal/app/auth/signature"
	"typerium/internal/app/auth/store"
	"typerium/internal/pkg/broker"
	"typerium/internal/pkg/broker/proto"
	_ "typerium/internal/pkg/config"
	"typerium/internal/pkg/waiter"

	"github.com/spf13/viper"
)

const (
	dbURI               = "DB_URI"
	dbVersion           = "DB_VERSION"
	grpcAddr            = "GRPC_ADDR"
	hashPasswordCostAlg = "HASH_PASSWORD_COST_ALGORITHM"
	accessTokenTTL      = "ACCESS_TOKEN_TTL"
	refreshTokenTTL     = "REFRESH_TOKEN_TTL"
	rsaSize             = "RSA_SIZE"
)

var grpcErrorsWrap = map[error]error{}

func main() {
	viper.SetDefault(dbURI, "postgresql://dev:123456@localhost:5432/auth?sslmode=disable")
	viper.SetDefault(dbVersion, uint(0))
	viper.SetDefault(grpcAddr, ":50052")
	viper.SetDefault(hashPasswordCostAlg, 12)
	viper.SetDefault(accessTokenTTL, time.Minute*15)
	viper.SetDefault(refreshTokenTTL, time.Hour*24)
	viper.SetDefault(rsaSize, 2048)

	db := store.New(viper.GetString(dbURI), viper.GetUint(dbVersion))
	defer db.Close()

	grpcServer := broker.NewGRPCServer(grpcErrorsWrap)
	proto.RegisterAuthServiceServer(grpcServer.Server,
		handlers.NewGRPCServer(
			db,
			password.NewBcryptProcessor(viper.GetInt(hashPasswordCostAlg)),
			jwt.SigningMethodRS512,
			signature.NewRSACreator(viper.GetInt(rsaSize)),
			viper.GetDuration(accessTokenTTL),
			viper.GetDuration(refreshTokenTTL),
		))

	grpcServer.ServeOnAddress(viper.GetString(grpcAddr))
	defer grpcServer.GracefulStop()

	waiter.Wait()
}
