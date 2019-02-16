// +build !windows

package main

import (
	"github.com/RivenZoo/backbone/services"
	"github.com/RivenZoo/backbone/signalutils"
	"os"
	"syscall"
)

func registerSignal() {
	signalutils.HandleSignals(func(sig os.Signal) {
		services.GetServiceContainer().Close()
	}, syscall.SIGINT, syscall.SIGTERM)
}
