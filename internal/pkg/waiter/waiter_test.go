package waiter

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestWait(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := zaptest.NewLogger(t)

	exit := make(chan struct{}, 1)
	go func() {
		select {
		case <-time.NewTimer(time.Second * 2).C:
			assert.Fail(t, "the stream wasn't released")
			os.Exit(1)
		case <-exit:
			close(exit)
			return

		}
	}()

	go time.AfterFunc(time.Second, func() {
		err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		if err != nil {
			assert.Fail(t, err.Error())
		}
	})

	Wait(log)
	exit <- struct{}{}
}
