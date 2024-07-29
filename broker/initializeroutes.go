package broker

import (
	"github.com/gin-gonic/gin"

	"github.com/s4mukka/justinject/domain"
)

var jobServiceFactory JobServiceFactory = JobServiceFactory{}

func route(fn func(ctx IRestContext)) gin.HandlerFunc {
	return func(ctx *gin.Context) { fn(ctx) }
}

func intializeRoutes(router domain.IRouter) {
	jobService := jobServiceFactory.MakeJobService(nil)

	basePath := "/api/v1"
	v1 := router.Group(basePath)
	{
		v1.POST("/job", route(jobService.CreateJob))
	}
	router.GET("/ping", func(ctx *gin.Context) {
		logger.WithContext(ctx.Request.Context()).Info("pong")
		ctx.JSON(200, "pong")
	})
}
