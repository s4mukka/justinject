package server

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/s4mukka/justinject/domain"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Server struct {
	Ctx *context.Context
}

func (s *Server) Init(intializeRoutes func(router *gin.Engine), port int) {
	environment := (*s.Ctx).Value("environment").(*domain.Environment)
	logger := environment.Logger

	gin.SetMode(os.Getenv("LOG_LEVEL"))
	router := gin.New()

	router.Use(gin.LoggerWithWriter(logger.Writer()))
	router.Use(otelgin.Middleware(environment.Instance))

	intializeRoutes(router)

	logger.Infof("Listening and serving HTTP on :%d", port)

	router.Run(fmt.Sprintf(":%d", port))
}
