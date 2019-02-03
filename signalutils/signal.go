package signalutils

import (
	"os"
	"os/signal"
)

type SignalHandler func(sig os.Signal)

func HandleSignals(h SignalHandler, sigs ...os.Signal) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, sigs...)
	go func() {
		sig := <-sigChan
		h(sig)
		close(sigChan)
	}()
}
