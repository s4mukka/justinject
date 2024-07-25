package context

import (
	"context"

	"github.com/s4mukka/justinject/domain"
	log "github.com/s4mukka/justinject/internal/logging"
)

var (
	loggerProviderFactory domain.ILoggerProviderFactory = LoggerProviderFactory{}
	tracerProviderFactory domain.ITracerProviderFactory = TracerProviderFactory{}
	addOtelHook                                         = log.AddOtelHook
	logInit                                             = log.Init
)

func InitializeContext(instance string) domain.IContext {
	environment := domain.Environment{Instance: instance}
	ctx := context.WithValue(context.Background(), "environment", &environment)

	logger := logInit(ctx)
	environment.Logger = logger

	loggerProvider, err := loggerProviderFactory.InitializeLoggerProvider(ctx)
	if err != nil {
		logger.Warnf("Error initializing logger: %v\n", err)
	} else {
		environment.LoggerProvider = loggerProvider
		addOtelHook(ctx)
	}

	tracerProvider, err := tracerProviderFactory.InitializeTracerProvider(ctx)
	if err != nil {
		logger.Warnf("Error initializing tracer: %v\n", err)
	}
	environment.TracerProvider = tracerProvider

	return ctx
}

func ShutdownComponents(ctx domain.IContext) {
	if ctx == nil {
		return
	}
	environment := ctx.Value("environment").(*domain.Environment)
	logger := environment.Logger

	loggerProvider := environment.LoggerProvider
	if loggerProvider != nil {
		if err := loggerProvider.Shutdown(context.Background()); err != nil {
			logger.Warnf("Error closing logger: %v\n", err)
		}
	}

	tracerProvider := environment.TracerProvider
	if tracerProvider != nil {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			logger.Warnf("Error closing tracer: %v\n", err)
		}
	}
}
