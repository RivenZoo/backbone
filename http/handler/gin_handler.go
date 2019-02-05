package handler

import "github.com/gin-gonic/gin"

type GinHandler struct {
	engine *gin.Engine
}

func NewGinHandler() *GinHandler {
	return &GinHandler{
		engine: gin.New(),
	}
}

func (h *GinHandler) GetGin() *gin.Engine {
	return h.engine
}
