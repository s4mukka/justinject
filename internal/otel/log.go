package otel

import (
	"context"
	"fmt"
	"os"

	"github.com/s4mukka/justinject/domain"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func InitLogger(ctx *context.Context) (*sdklog.LoggerProvider, error) {
	environment := (*ctx).Value("environment").(*domain.Environment)

	var otelEndpoint string

	if otelEndpoint = os.Getenv("OTEL_ENDPOINT_HTTP"); otelEndpoint == "" {
		return nil, fmt.Errorf("OTEL_ENDPOINT_HTTP environment variable is not defined")
	}

	exporter, err := otlploghttp.New(*ctx,
		otlploghttp.WithInsecure(),
		otlploghttp.WithEndpoint(otelEndpoint),
	)
	if err != nil {
		return nil, fmt.Errorf("Error creating OTLP exporter: %v\n", err)
	}

	stdoutExporter, _ := stdoutlog.New(stdoutlog.WithPrettyPrint())

	lp := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(stdoutExporter)),
		sdklog.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(environment.Instance),
		)),
	)

	global.SetLoggerProvider(lp)
	return lp, nil
}
