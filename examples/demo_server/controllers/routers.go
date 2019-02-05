package controllers

import (
	"github.com/RivenZoo/backbone/http/handler"
	"github.com/RivenZoo/backbone/services/httpserver"
)

func InitRouters() {
	r := handler.NewGinHandler().GetGin()
	defer httpserver.GetHTTPServer().SetHTTPHandler(r)

	// set url handlers
	r.POST("/url/abbr", abbreviateURLProcessor)
}
