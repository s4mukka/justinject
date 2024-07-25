package domain

import (
	"context"

	"go.opentelemetry.io/otel/log"
	logembedded "go.opentelemetry.io/otel/log/embedded"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	traceembedded "go.opentelemetry.io/otel/trace/embedded"
)

type OtelLoggerProvider = *sdklog.LoggerProvider
type OtelLoggerOption = log.LoggerOption
type OtelLogger = log.Logger
type OtelRecord = log.Record

type OtelTracerProvider = *sdktrace.TracerProvider
type OtelTracerOption = trace.TracerOption
type OtelTracer = trace.Tracer

type IOtelLoggerProvider interface {
	logembedded.LoggerProvider
	Logger(name string, opts ...OtelLoggerOption) OtelLogger
	Shutdown(ctx context.Context) error
}

type ILoggerProvider interface {
	Get() IOtelLoggerProvider
	Logger(name string, opts ...OtelLoggerOption) OtelLogger
	Shutdown(ctx context.Context) error
}

type ILoggerProviderFactory interface {
	InitializeLoggerProvider(ctx IContext) (ILoggerProvider, error)
}

type IOtelTracerProvider interface {
	traceembedded.TracerProvider
	Tracer(name string, opts ...trace.TracerOption) trace.Tracer
	Shutdown(ctx context.Context) error
}

type ITracerProvider interface {
	Get() IOtelTracerProvider
	Tracer(name string, opts ...trace.TracerOption) trace.Tracer
	Shutdown(ctx context.Context) error
}

type ITracerProviderFactory interface {
	InitializeTracerProvider(ctx IContext) (ITracerProvider, error)
}
