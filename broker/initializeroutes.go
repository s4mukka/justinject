package broker

import (
	"github.com/gin-gonic/gin"

	"github.com/s4mukka/justinject/domain"
)

var jobServiceFactory domain.IFactory[domain.IJobService] = JobServiceFactory{
	jobUseCaseFactory: &JobUseCaseFactory{
		extractorRepositoryFactory: nil,
		jobRepositoryFactory:       nil,
		k8sRepositoryFactory:       &K8sRepositoryFactory{},
	},
}

func route(fn func(ctx domain.IRestContext)) gin.HandlerFunc {
	return func(ctx *gin.Context) { fn(ctx) }
}

func intializeRoutes(router domain.IRouter) error {
	jobService, err := jobServiceFactory.Create()
	if err != nil {
		return err
	}

	basePath := "/api/v1"
	v1 := router.Group(basePath)
	{
		v1.POST("/job", route(jobService.CreateJob))
	}
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, "pong")
	})
	return nil
}
