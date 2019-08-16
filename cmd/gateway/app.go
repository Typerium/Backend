package main

import (
	"github.com/spf13/viper"

	"typerium/internal/app/gateway/graphql"
	"typerium/internal/app/gateway/handlers"
	"typerium/internal/pkg/broker"
	_ "typerium/internal/pkg/config"
	"typerium/internal/pkg/waiter"
	"typerium/internal/pkg/web"
)

const (
	httpServerAddr            = "HTTP_SERVER_ADDR"
	profilesManagerServiceURI = "PROFILES_MANAGER_SERVICE_URI"
	authServiceURI            = "AUTH_SERVICE_URI"
)

func main() {
	viper.SetDefault(httpServerAddr, ":10000")
	viper.SetDefault(profilesManagerServiceURI, ":50051")
	viper.SetDefault(authServiceURI, ":50052")

	server := web.NewServer()

	gqlExecutor := graphql.New(
		broker.NewAuthClient(viper.GetString(authServiceURI)),
		broker.NewProfilesManagerClient(viper.GetString(profilesManagerServiceURI)),
	)
	gqlHandler := handlers.Handler(gqlExecutor)
	gqlRoute := "/graphql"
	server.GET(gqlRoute, gqlHandler)
	server.POST(gqlRoute, gqlHandler)

	server.Start(viper.GetString(httpServerAddr))
	defer server.Stop()

	waiter.Wait()
}
