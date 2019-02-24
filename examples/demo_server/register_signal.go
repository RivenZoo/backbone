// +build !windows

package main

import (
	"github.com/RivenZoo/backbone/objects_container"
	"github.com/RivenZoo/backbone/signalutils"

	"os"
	"syscall"
)

func registerSignal() {
	signalutils.HandleSignals(func(sig os.Signal) {
		objects_container.Close()
	}, syscall.SIGINT, syscall.SIGTERM)
}
