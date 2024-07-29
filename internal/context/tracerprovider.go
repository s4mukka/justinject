package context

import (
	"github.com/s4mukka/justinject/domain"
	"github.com/s4mukka/justinject/internal/otel"
)

var otelInitTracer = otel.InitTracer

type TracerProviderFactory struct{}

func (l TracerProviderFactory) InitializeTracerProvider(ctx domain.IContext) (domain.ITracerProvider, error) {
	environment := ctx.Value(domain.EnvironmentKey).(*domain.Environment)
	tracerProvider, err := otelInitTracer(ctx)
	if err != nil {
		return nil, err
	}
	tracerProvider.Tracer(environment.Instance)
	return tracerProvider, nil
}
