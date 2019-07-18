package main

import (
	"github.com/spf13/viper"

	"typerium/internal/app/gateway/graphql"
	"typerium/internal/app/gateway/handlers"
	_ "typerium/internal/pkg/config"
	"typerium/internal/pkg/logging"
	"typerium/internal/pkg/waiter"
	"typerium/internal/pkg/web"
)

const (
	httpServerAddr = "HTTP_SERVER_ADDR"
)

func main() {
	log := logging.New()

	server := web.NewServer(log)

	gqlExecutor := graphql.New(log)
	gqlHandler := handlers.Handler(gqlExecutor)
	gqlRoute := "/graphql"
	server.GET(gqlRoute, gqlHandler)
	server.POST(gqlRoute, gqlHandler)

	viper.SetDefault(httpServerAddr, ":10000")

	server.Start(viper.GetString(httpServerAddr))
	defer server.Stop()

	waiter.Wait(log)
}
