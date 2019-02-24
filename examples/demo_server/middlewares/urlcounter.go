package middlewares

import (
	"github.com/RivenZoo/backbone/examples/demo_server/model"
	"github.com/RivenZoo/backbone/logger"
	"github.com/gin-gonic/gin"
)

func CountURL() func(c *gin.Context) {
	return func(c *gin.Context) {
		counter := model.GetRedisCounter()
		n, e := counter.Count(c.Request.URL.Path)
		if e != nil {
			logger.Errorf("count url path %s error %v", c.Request.URL.Path, e)
		} else {
			logger.Infof("url path %s count %d", c.Request.URL.Path, n)
		}
		c.Next()
	}
}
