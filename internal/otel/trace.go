package otel

import (
	"context"
	"fmt"
	"os"

	"github.com/s4mukka/justinject/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func InitTracer(ctx *context.Context) (*sdktrace.TracerProvider, error) {
	environment := (*ctx).Value("environment").(*domain.Environment)
	// logger := environment.Logger

	var otelEndpoint string

	if otelEndpoint = os.Getenv("OTEL_ENDPOINT_GRPC"); otelEndpoint == "" {
		return nil, fmt.Errorf("OTEL_ENDPOINT_GRPC environment variable is not defined")
	}

	exporter, err := otlptracegrpc.New(*ctx,
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

	otel.SetTracerProvider(tp)
	return tp, nil
}
