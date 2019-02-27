package controllers

import "github.com/RivenZoo/backbone/http/handler"
import "github.com/gin-gonic/gin"

var ginAbbreviateURLReqHandler = handler.NewRequestHandleFunc(&handler.RequestProcessor{
	NewReqFunc: func() interface{} {
		return &abbreviateURLReq{}
	},
	ProcessFunc: func(c *gin.Context, req interface{}) (resp interface{}, err error) {
		concreteReq := req.(*abbreviateURLReq)
		return handleAbbreviateURLReq(c, concreteReq)
	},
})
