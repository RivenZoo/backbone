package main

import "os"
import (
	"github.com/RivenZoo/backbone/objects_container"
	"github.com/RivenZoo/backbone/signalutils"
)

func registerSignal() {
	signalutils.HandleSignals(func(sig os.Signal) {
		objects_container.Close()
	}, os.Interrupt)
}
