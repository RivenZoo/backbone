package main

import (
	"flag"
	"github.com/RivenZoo/backbone/examples/demo_server/config"
	"github.com/RivenZoo/backbone/examples/demo_server/controllers"
	"github.com/RivenZoo/backbone/http/logger"
	"github.com/RivenZoo/backbone/resources"
	"github.com/RivenZoo/backbone/services"
	"github.com/RivenZoo/backbone/signalutils"
	"os"
	"syscall"
)

var cfgFile *string

func main() {
	cfgFile = flag.String("cfg", "./conf/cfg.json", "config file path")
	flag.Parse()

	// first: config init
	logger.Log("load config")
	config.MustLoadConfig(*cfgFile)

	// second: init all resource
	logger.Log("init resources")
	resources.GetResourceContainer().Init()
	defer resources.GetResourceContainer().Close()

	// third: init all service
	logger.Log("init services")
	services.GetServiceContainer().Init()

	signalutils.HandleSignals(func(sig os.Signal) {
		services.GetServiceContainer().Close()
	}, syscall.SIGINT, syscall.SIGTERM)

	// init controllers
	controllers.InitRouters()

	// last: run service
	logger.Log("run services")
	services.GetServiceContainer().RunServices()

	logger.Log("exit")
}
