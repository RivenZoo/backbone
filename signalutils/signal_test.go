package signalutils

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestHandleSignals(t *testing.T) {
	ctx, fn := context.WithCancel(context.Background())
	h := SignalHandler(func(sig os.Signal) {
		t.Log("recv ", sig)
		switch sig {
		case syscall.SIGINT:
		case syscall.SIGHUP:
		}
		fn()
	})
	HandleSignals(h, syscall.SIGINT, syscall.SIGHUP)
	go func() {
		time.Sleep(time.Millisecond*10)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		t.Log("sent signal")
	}()

	t.Log("wait context")
	<-ctx.Done()
	t.Log("context done")
}
