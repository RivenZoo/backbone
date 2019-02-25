package main

import "os"
import (
	"github.com/RivenZoo/backbone/services"
	"github.com/RivenZoo/backbone/signalutils"
)

func registerSignal() {
	signalutils.HandleSignals(func(sig os.Signal) {
		services.Close()
	}, os.Interrupt)
}
