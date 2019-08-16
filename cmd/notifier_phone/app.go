package main

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"typerium/internal/app/notifier_phone/handlers"
	"typerium/internal/app/notifier_phone/provider"
	_ "typerium/internal/pkg/config"
	"typerium/internal/pkg/logging"
	"typerium/internal/pkg/waiter"
	"typerium/internal/pkg/web"
)

const (
	twilioUsername = "TWILIO_USERNAME"
	twilioPassword = "TWILIO_PASSWORD"
)

func main() {
	log := logging.New("")
	
	viper.SetDefault(twilioUsername, "user")
	viper.SetDefault(twilioPassword, "123456")

	clientFactory := web.NewFasthttpClientFactory()
	twilioProvider, err := provider.NewTwilioProvider(clientFactory, viper.GetString(twilioUsername),
		viper.GetString(twilioPassword))
	if err != nil {
		log.Fatal("twilio provider can't create", zap.Error(err))
	}

	handlers.NewQueueBroker(twilioProvider)

	waiter.Wait()
}
