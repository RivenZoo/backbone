package controllers

import (
	//"github.com/RivenZoo/backbone/http/handler"
	"github.com/RivenZoo/backbone/services/httpserver"
)

func InitRouters() {
	//engine := handler.NewGinHandler()
	httpserver.GetHTTPServer().SetHTTPHandler(nil)
}
