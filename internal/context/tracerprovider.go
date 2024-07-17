package context

import (
	"context"

	"github.com/s4mukka/justinject/domain"
	"github.com/s4mukka/justinject/internal/otel"
)

var (
	otelInitTracer = otel.InitTracer
)

type TracerProviderFactory struct{}

func (l TracerProviderFactory) InitializeTracerProvider(ctx context.Context) (domain.ITracerProvider, error) {
	environment := ctx.Value("environment").(*domain.Environment)
	tracerProvider, err := otelInitTracer(&ctx)
	if err != nil {
		return nil, err
	}
	tracerProvider.Tracer(environment.Instance)
	return tracerProvider, nil
}
