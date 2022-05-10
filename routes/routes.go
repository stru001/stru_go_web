package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stru_web/logger"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(),logger.GinRecovery(true))
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK,"welcome to stru_web with gin")
	})
	return r
}
