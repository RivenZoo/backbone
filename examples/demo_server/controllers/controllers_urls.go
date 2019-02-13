package controllers

import "github.com/RivenZoo/backbone/http/handler"
import "github.com/RivenZoo/backbone/services/httpserver"

func InitRouters() {
	engine := handler.NewGinHandler().GetGin()
	defer httpserver.GetHTTPServer().SetHTTPHandler(engine)

	engine.POST("/url/abbr", abbreviateURLReqHandler)

}
