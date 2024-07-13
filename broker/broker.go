package broker

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/s4mukka/justinject/domain"
	log "github.com/s4mukka/justinject/internal/logging"
	"github.com/s4mukka/justinject/internal/otel"
	"github.com/s4mukka/justinject/internal/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

const (
	instance = "broker"
)

var (
	logger *logrus.Entry
)

func Init() *cobra.Command {
	return &cobra.Command{
		Use:   "broker",
		Short: "Starting broker",
		Run: func(cmd *cobra.Command, args []string) {
			start()
		},
	}
}

func start() {
	var lp *sdklog.LoggerProvider
	environment := domain.Environment{Instance: instance}
	ctx := context.WithValue(context.Background(), "environment", &environment)
	logger, lp = log.Init(&ctx)
	defer func() {
		if err := lp.Shutdown(context.Background()); err != nil {
			logger.Warnf("Error closing logger: %v", err)
		}
	}()
	environment.Logger = logger

	tp, err := otel.InitTracer(&ctx)
	if err != nil {
		logger.Warnf("Error initializing tracer: %v", err)
	} else {
		tp.Tracer(environment.Instance)
		defer func() {
			if err := tp.Shutdown(context.Background()); err != nil {
				logger.Warnf("Error closing tracer: %v", err)
			}
		}()
	}

	logger.Info("Starting broker...")
	port := 8080
	server := server.Server{Ctx: &ctx}
	server.Init(intializeRoutes, port)
}

func intializeRoutes(router *gin.Engine) {
	router.GET("/ping", func(ctx *gin.Context) {
		logger.WithContext(ctx.Request.Context()).Warn("pong")
		ctx.JSON(200, "pong")
	})
}
