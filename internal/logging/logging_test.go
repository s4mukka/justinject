package logging

import (
	"context"
	"testing"

	"github.com/s4mukka/justinject/domain"
	"github.com/s4mukka/justinject/internal/otellogrusdecorator"
	"github.com/s4mukka/justinject/mock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/contrib/bridges/otellogrus"
)

func TestInit(t *testing.T) {
	environment := &domain.Environment{
		Instance: "test-instance",
	}

	ctx := context.WithValue(context.Background(), "environment", environment)

	logger := Init(&ctx)

	assert.NotNil(t, logger)
	assert.Equal(t, "test-instance", logger.Data["instance"])
}

func TestAddOtelHook(t *testing.T) {
	mockProvider := &mock.MockedLoggerProvider{}
	mockProvider.On("Get").Return(new(domain.OtelLoggerProvider))

	mockEnv := &domain.Environment{
		Instance:       "test-instance",
		LoggerProvider: mockProvider,
		Logger:         log.New().WithField("mock", "logger"),
	}

	ctx := context.WithValue(context.Background(), "environment", mockEnv)

	originalNewHook := newHook
	defer func() { newHook = originalNewHook }()
	mockHook := new(mock.MockedHook)
	newHook = func(instance string, opts ...otellogrus.Option) otellogrusdecorator.Hook {
		return mockHook
	}

	mockHook.On("Levels").Return([]log.Level{
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel,
		log.WarnLevel,
		log.InfoLevel,
	})

	logger := Init(&ctx)

	AddOtelHook(&ctx)

	assert.NotNil(t, logger.Logger.Hooks)
}
