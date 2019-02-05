package main

import (
	"flag"
	"github.com/RivenZoo/backbone/examples/demo_server/config"
	"github.com/RivenZoo/backbone/examples/demo_server/controllers"
	"github.com/RivenZoo/backbone/resources"
	"github.com/RivenZoo/backbone/services"
	"github.com/RivenZoo/backbone/http/logger"
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
	services.GetServiceContainer().Close()

	// init controllers
	controllers.InitRouters()

	// last: run service
	logger.Log("run services")
	services.GetServiceContainer().RunServices()

	logger.Log("exit")
}
