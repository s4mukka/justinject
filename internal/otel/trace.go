package otel

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/s4mukka/justinject/domain"
)

type TracerProvider struct {
	handler domain.IOtelTracerProvider
}

func (t TracerProvider) Get() domain.IOtelTracerProvider {
	return t.handler
}

func (t TracerProvider) Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	return t.handler.Tracer(name, opts...)
}

func (t TracerProvider) Shutdown(ctx context.Context) error {
	return t.handler.Shutdown(ctx)
}

var otlptracegrpcNew = otlptracegrpc.New

func InitTracer(ctx domain.IContext) (domain.ITracerProvider, error) {
	environment := ctx.Value(domain.EnvironmentKey).(*domain.Environment)

	var otelEndpoint string

	if otelEndpoint = os.Getenv("OTEL_ENDPOINT_GRPC"); otelEndpoint == "" {
		return nil, fmt.Errorf("OTEL_ENDPOINT_GRPC environment variable is not defined")
	}

	exporter, err := otlptracegrpcNew(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(otelEndpoint),
	)
	if err != nil {
		return nil, fmt.Errorf("Error creating OTLP exporter: %v\n", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(environment.Instance),
		)),
	)

	tracerProvider := TracerProvider{handler: tp}

	otel.SetTracerProvider(tracerProvider.Get())
	return tracerProvider, nil
}
