package logging

import (
	"context"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/contrib/bridges/otellogrus"

	"github.com/s4mukka/justinject/domain"
	"github.com/s4mukka/justinject/internal/logging/mocks"
	mocksHook "github.com/s4mukka/justinject/internal/loggrushook/mocks"
)

func TestInit(t *testing.T) {
	environment := &domain.Environment{
		Instance: "test-instance",
	}

	ctx := context.WithValue(context.Background(), domain.EnvironmentKey, environment)

	logger := Init(ctx)

	assert.NotNil(t, logger)
	assert.Equal(t, "test-instance", logger.(*log.Entry).Data["instance"])
}

func TestAddOtelHook(t *testing.T) {
	mockProvider := &mocks.MockedLoggerProvider{}
	mockProvider.On("Get").Return(new(domain.OtelLoggerProvider))

	mockEnv := &domain.Environment{
		Instance:       "test-instance",
		LoggerProvider: mockProvider,
		Logger:         log.New().WithField("mock", "logger"),
	}

	ctx := context.WithValue(context.Background(), domain.EnvironmentKey, mockEnv)

	originalNewHook := newOtelLoggrusHook
	defer func() { newOtelLoggrusHook = originalNewHook }()
	hook := new(mocksHook.MockedHook)
	newOtelLoggrusHook = func(instance string, opts ...otellogrus.Option) domain.IHook {
		return hook
	}

	hook.On("Levels").Return([]log.Level{
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel,
		log.WarnLevel,
		log.InfoLevel,
	})

	logger := Init(ctx)

	AddOtelHook(ctx)

	assert.NotNil(t, logger.(*log.Entry).Logger.Hooks)
}
