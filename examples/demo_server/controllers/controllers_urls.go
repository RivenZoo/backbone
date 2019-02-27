package controllers

import "github.com/gin-gonic/gin"

func InitRouters(engine *gin.Engine) {

	engine.POST("/url/abbr", ginAbbreviateURLReqHandler)

}
