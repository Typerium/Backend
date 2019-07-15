package waiter

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

// Wait function for wait signal and continue execute thread
func Wait(log *zap.Logger) {
	sig := make(chan os.Signal, 1)
	defer close(sig)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	receivedSignal := <-sig

	log.Info(fmt.Sprintf("received signal '%s'", receivedSignal.String()))
}
