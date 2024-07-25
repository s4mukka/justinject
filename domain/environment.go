package domain

import (
	"context"

	log "github.com/sirupsen/logrus"
)

type ILogger interface {
	log.FieldLogger
	WithContext(ctx context.Context) *log.Entry
}

type Environment struct {
	Instance string

	Logger ILogger

	LoggerProvider ILoggerProvider
	TracerProvider ITracerProvider
}

const EnvironmentKey ContextKey = "environment"
