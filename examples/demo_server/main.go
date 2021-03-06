package main

import (
	"flag"
	"github.com/RivenZoo/backbone/examples/demo_server/config"
	"github.com/RivenZoo/backbone/examples/demo_server/controllers"
	"github.com/RivenZoo/backbone/examples/demo_server/middlewares"
	"github.com/RivenZoo/backbone/http/handler"
	"github.com/RivenZoo/backbone/logger"
	"github.com/RivenZoo/backbone/objects_container"
	"github.com/RivenZoo/backbone/services"
	"github.com/RivenZoo/backbone/resources"
	"github.com/RivenZoo/backbone/services/httpserver"
	"github.com/gin-gonic/gin"
)

var cfgFile *string

func initMiddleware(g *gin.Engine) {
	g.Use(middlewares.CountURL())
}

func main() {
	cfgFile = flag.String("cfg", "./conf/cfg.json", "config file path")
	flag.Parse()

	// first: config init
	logger.Log("load config")
	config.MustLoadConfig(*cfgFile)

	// second: init all resources and services
	logger.Log("init resources and services")

	resources.Init()
	defer resources.Close()

	objects_container.Init()
	defer objects_container.Close()

	services.Init()

	registerSignal()

	// init controllers
	g := handler.NewGinHandler().GetGin()
	initMiddleware(g)
	controllers.InitRouters(g)
	httpserver.GetHTTPServer().SetHTTPHandler(g)

	// last: run service
	logger.Log("run services")
	services.RunServices()

	logger.Log("exit")
}
