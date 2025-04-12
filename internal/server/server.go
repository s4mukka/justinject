package server

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/s4mukka/justinject/domain"
)

var router domain.IRouter = gin.New()

type Server struct {
	Ctx domain.IContext
}

func (s *Server) Init(intializeRoutes func(router domain.IRouter) error, port int) error {
	environment := s.Ctx.Value(domain.EnvironmentKey).(*domain.Environment)
	logger := environment.Logger

	gin.SetMode(os.Getenv("LOG_LEVEL"))

	router.Use(otelgin.Middleware(environment.Instance))

	if err := intializeRoutes(router); err != nil {
		logger.Errorf("Routes initialization failed: %s", err.Error())
		return err
	}

	logger.Infof("Listening and serving HTTP on :%d", port)

	return router.Run(fmt.Sprintf(":%d", port))
}

type ServerFactory struct{}

func (sf *ServerFactory) MakeServer(ctx domain.IContext) domain.IServer {
	return &Server{Ctx: ctx}
}
