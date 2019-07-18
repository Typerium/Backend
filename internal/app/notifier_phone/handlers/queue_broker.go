package handlers

import (
	"typerium/internal/app/notifier_phone/provider"
)

type queueHandlers struct {
	provider provider.Provider
}

func NewQueueBroker(provider provider.Provider) *queueHandlers {
	return &queueHandlers{provider: provider}
} 
