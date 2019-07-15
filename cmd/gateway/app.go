package main

import (
	"github.com/spf13/viper"

	"typerium/internal/app/gateway/graphql"
	"typerium/internal/app/gateway/handlers"
	"typerium/internal/pkg/httpserver"
	"typerium/internal/pkg/logging"
	"typerium/internal/pkg/waiter"
)

const (
	httpServerAddr = "HTTP_SERVER_ADDR"
)

func main() {
	log := logging.New()

	server := httpserver.NewServer(log)

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
