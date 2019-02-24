package main

import (
	"flag"
	"github.com/RivenZoo/backbone/examples/demo_server/config"
	"github.com/RivenZoo/backbone/examples/demo_server/controllers"
	"github.com/RivenZoo/backbone/http/handler"
	"github.com/RivenZoo/backbone/logger"
	"github.com/RivenZoo/backbone/services"
	"github.com/RivenZoo/backbone/objects_container"
	"github.com/RivenZoo/backbone/services/httpserver"
)

var cfgFile *string

func main() {
	cfgFile = flag.String("cfg", "./conf/cfg.json", "config file path")
	flag.Parse()

	// first: config init
	logger.Log("load config")
	config.MustLoadConfig(*cfgFile)

	// second: init all resources and services
	logger.Log("init resources and services")
	objects_container.Init()

	registerSignal()

	// init controllers
	g := handler.NewGinHandler().GetGin()
	controllers.InitRouters(g)
	httpserver.GetHTTPServer().SetHTTPHandler(g)

	// last: run service
	logger.Log("run services")
	services.GetServiceContainer().RunServices()

	logger.Log("exit")
}
