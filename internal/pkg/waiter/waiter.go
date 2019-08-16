package waiter

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"typerium/internal/pkg/logging"
)

// Wait function for wait signal and continue execute thread
func Wait() {
	sig := make(chan os.Signal, 1)
	defer close(sig)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	receivedSignal := <-sig

	logging.New("waiter").Info(fmt.Sprintf("received signal '%s'", receivedSignal.String()))
}
