package context

import (
	"context"

	"github.com/s4mukka/justinject/domain"
	"github.com/s4mukka/justinject/internal/otel"
)

var (
	otelInitLogger = otel.InitLogger
)

type LoggerProviderFactory struct{}

func (l LoggerProviderFactory) InitializeLoggerProvider(ctx context.Context) (domain.ILoggerProvider, error) {
	environment := ctx.Value("environment").(*domain.Environment)

	loggerProvider, err := otelInitLogger(&ctx)
	if err != nil {
		return nil, err
	}
	loggerProvider.Logger(environment.Instance)
	return loggerProvider, nil
}
