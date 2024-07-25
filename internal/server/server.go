package server

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/s4mukka/justinject/domain"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var (
	router domain.IRouter = gin.New()
)

type Server struct {
	Ctx domain.IContext
}

func (s *Server) Init(intializeRoutes func(router domain.IRouter), port int) error {
	environment := s.Ctx.Value("environment").(*domain.Environment)
	logger := environment.Logger

	gin.SetMode(os.Getenv("LOG_LEVEL"))

	router.Use(otelgin.Middleware(environment.Instance))

	intializeRoutes(router)

	logger.Infof("Listening and serving HTTP on :%d", port)

	return router.Run(fmt.Sprintf(":%d", port))
}

type ServerFactory struct{}

func (sf *ServerFactory) MakeServer(ctx domain.IContext) domain.IServer {
	return &Server{Ctx: ctx}
}
