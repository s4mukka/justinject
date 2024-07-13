package domain

import (
	log "github.com/sirupsen/logrus"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Environment struct {
	Instance string

	Logger *log.Entry

	LoggerProvider *sdklog.LoggerProvider
	TracerProvider *sdktrace.TracerProvider
}
