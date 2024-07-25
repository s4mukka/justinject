package broker

import (
	"github.com/gin-gonic/gin"
	"github.com/s4mukka/justinject/domain"
)

func intializeRoutes(router domain.IRouter) {
	router.GET("/ping", func(ctx *gin.Context) {
		logger.WithContext(ctx.Request.Context()).Info("pong")
		ctx.JSON(200, "pong")
	})
}
