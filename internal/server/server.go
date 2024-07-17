package server

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/s4mukka/justinject/domain"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type IServer interface {
	Init(intializeRoutes func(router IRouter), port int) error
}
type Server struct {
	Ctx *context.Context
}

type IRouter interface {
	gin.IRoutes
	Use(middleware ...gin.HandlerFunc) gin.IRoutes
	Run(addr ...string) (err error)
}

var (
	router IRouter = gin.New()
)

func (s *Server) Init(intializeRoutes func(router IRouter), port int) error {
	environment := (*s.Ctx).Value("environment").(*domain.Environment)
	logger := environment.Logger

	gin.SetMode(os.Getenv("LOG_LEVEL"))

	router.Use(otelgin.Middleware(environment.Instance))

	intializeRoutes(router)

	logger.Infof("Listening and serving HTTP on :%d", port)

	return router.Run(fmt.Sprintf(":%d", port))
}

type IServerFactory interface {
	MakeServer(ctx *context.Context) IServer
}

type ServerFactory struct{}

func (sf *ServerFactory) MakeServer(ctx *context.Context) IServer {
	return &Server{Ctx: ctx}
}
