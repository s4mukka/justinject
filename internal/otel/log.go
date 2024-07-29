package otel

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/s4mukka/justinject/domain"
)

type LoggerProvider struct {
	handler domain.IOtelLoggerProvider
}

func (l LoggerProvider) Get() domain.IOtelLoggerProvider {
	return l.handler
}

func (l LoggerProvider) Logger(name string, opts ...log.LoggerOption) log.Logger {
	return l.handler.Logger(name, opts...)
}

func (l LoggerProvider) Shutdown(ctx context.Context) error {
	return l.handler.Shutdown(ctx)
}

var otlploghttpNew = otlploghttp.New

func InitLogger(ctx domain.IContext) (domain.ILoggerProvider, error) {
	environment := ctx.Value(domain.EnvironmentKey).(*domain.Environment)

	var otelEndpoint string

	if otelEndpoint = os.Getenv("OTEL_ENDPOINT_HTTP"); otelEndpoint == "" {
		return nil, fmt.Errorf("OTEL_ENDPOINT_HTTP environment variable is not defined")
	}

	exporter, err := otlploghttpNew(ctx,
		otlploghttp.WithInsecure(),
		otlploghttp.WithEndpoint(otelEndpoint),
	)
	if err != nil {
		return nil, fmt.Errorf("Error creating OTLP exporter: %v\n", err)
	}

	lp := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
		sdklog.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(environment.Instance),
		)),
	)

	loggerProvider := &LoggerProvider{handler: lp}

	global.SetLoggerProvider(loggerProvider.Get())
	return loggerProvider, nil
}
