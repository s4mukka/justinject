package context

import (
	"github.com/s4mukka/justinject/domain"
	"github.com/s4mukka/justinject/internal/otel"
)

var otelInitLogger = otel.InitLogger

type LoggerProviderFactory struct{}

func (l LoggerProviderFactory) InitializeLoggerProvider(ctx domain.IContext) (domain.ILoggerProvider, error) {
	environment := ctx.Value(domain.EnvironmentKey).(*domain.Environment)

	loggerProvider, err := otelInitLogger(ctx)
	if err != nil {
		return nil, err
	}
	loggerProvider.Logger(environment.Instance)
	return loggerProvider, nil
}
