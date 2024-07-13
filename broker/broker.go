package broker

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/s4mukka/justinject/domain"
	log "github.com/s4mukka/justinject/internal/logging"
	"github.com/s4mukka/justinject/internal/otel"
	"github.com/s4mukka/justinject/internal/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	instance = "broker"
	port     = 8080
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

var (
	logger *logrus.Entry
)

func start() {
	ctx := initializeContext()
	defer shutdownComponents(ctx)

	logger.Info("Starting broker...")
	server := server.Server{Ctx: &ctx}
	if err := server.Init(intializeRoutes, port); err != nil {
		logger.Errorf("Error starting broker: %v\n", err)
	}
}

func initializeContext() context.Context {
	environment := domain.Environment{Instance: instance}
	ctx := context.WithValue(context.Background(), "environment", &environment)

	logger = log.Init(&ctx)
	environment.Logger = logger

	loggerProvider, err := initializeLoggerProvider(ctx)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		return nil
	}
	environment.LoggerProvider = loggerProvider

	log.AddOtelHook(&ctx)

	tracerProvider, err := initializeTracerProvider(ctx)
	if err != nil {
		logger.Warnf("Error initializing tracer: %v", err)
	}
	environment.TracerProvider = tracerProvider

	return ctx
}

func initializeLoggerProvider(ctx context.Context) (*sdklog.LoggerProvider, error) {
	environment := ctx.Value("environment").(*domain.Environment)

	loggerProvider, err := otel.InitLogger(&ctx)
	if err != nil {
		return nil, err
	}
	loggerProvider.Logger(environment.Instance)
	return loggerProvider, nil
}

func initializeTracerProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	tracerProvider, err := otel.InitTracer(&ctx)
	if err != nil {
		return nil, err
	}
	tracerProvider.Tracer(instance)
	return tracerProvider, nil
}

func shutdownComponents(ctx context.Context) {
	if ctx == nil {
		return
	}
	environment := ctx.Value("environment").(*domain.Environment)

	loggerProvider := environment.LoggerProvider
	if loggerProvider != nil {
		if err := loggerProvider.Shutdown(context.Background()); err != nil {
			logger.Warnf("Error closing logger: %v", err)
		}
	}

	tracerProvider := environment.TracerProvider
	if tracerProvider != nil {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			logger.Warnf("Error closing tracer: %v", err)
		}
	}
}

func intializeRoutes(router *gin.Engine) {
	router.GET("/ping", func(ctx *gin.Context) {
		logger.WithContext(ctx.Request.Context()).Warn("pong")
		logger.Warn("pong")
		ctx.JSON(200, "pong")
	})
}
