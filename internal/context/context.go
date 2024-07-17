package context

import (
	"context"
	"fmt"
	"reflect"

	"github.com/s4mukka/justinject/domain"
	log "github.com/s4mukka/justinject/internal/logging"
)

var (
	loggerProviderFactory domain.ILoggerProviderFactory = LoggerProviderFactory{}
	tracerProviderFactory domain.ITracerProviderFactory = TracerProviderFactory{}
	addOtelHook                                         = log.AddOtelHook
)

func InitializeContext(instance string) context.Context {
	environment := domain.Environment{Instance: instance}
	ctx := context.WithValue(context.Background(), "environment", &environment)

	logger := log.Init(&ctx)
	environment.Logger = logger

	loggerProvider, err := loggerProviderFactory.InitializeLoggerProvider(ctx)
	if err != nil {
		logger.Warnf("Error initializing logger: %v\n", err)
	} else {
		environment.LoggerProvider = loggerProvider
		addOtelHook(&ctx)
	}

	tracerProvider, err := tracerProviderFactory.InitializeTracerProvider(ctx)
	if err != nil {
		logger.Warnf("Error initializing tracer: %v\n", err)
	}
	environment.TracerProvider = tracerProvider

	return ctx
}

func ShutdownComponents(ctx context.Context) {
	if ctx == nil {
		return
	}
	environment := ctx.Value("environment").(*domain.Environment)
	logger := environment.Logger

	loggerProvider := environment.LoggerProvider
	if loggerProvider != nil {
		if err := loggerProvider.Shutdown(context.Background()); err != nil {
			fmt.Printf("SHUT %p %v", &environment.Logger, reflect.TypeOf(logger))
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
