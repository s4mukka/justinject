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

type IOtelLoggerProvider interface {
	logembedded.LoggerProvider
	ForceFlush(ctx context.Context) error
	Logger(name string, opts ...OtelLoggerOption) OtelLogger
	Shutdown(ctx context.Context) error
}

type OtelLoggerProvider = *sdklog.LoggerProvider
type OtelLoggerOption = log.LoggerOption
type OtelLogger = log.Logger
type OtelRecord = log.Record

type ILoggerProvider interface {
	Get() IOtelLoggerProvider
	Logger(name string, opts ...OtelLoggerOption) OtelLogger
	Shutdown(ctx context.Context) error
}

type ILoggerProviderFactory interface {
	InitializeLoggerProvider(ctx context.Context) (ILoggerProvider, error)
}

type IOtelTracerProvider interface {
	traceembedded.TracerProvider
	RegisterSpanProcessor(sp sdktrace.SpanProcessor)
	UnregisterSpanProcessor(sp sdktrace.SpanProcessor)
	ForceFlush(ctx context.Context) error
	Tracer(name string, opts ...trace.TracerOption) trace.Tracer
	Shutdown(ctx context.Context) error
}

type OtelTracerProvider = *sdktrace.TracerProvider
type OtelTracerOption = trace.TracerOption
type OtelTracer = trace.Tracer

type ITracerProvider interface {
	Get() IOtelTracerProvider
	Tracer(name string, opts ...trace.TracerOption) trace.Tracer
	Shutdown(ctx context.Context) error
}

type ITracerProviderFactory interface {
	InitializeTracerProvider(ctx context.Context) (ITracerProvider, error)
}
