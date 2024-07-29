package context

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/s4mukka/justinject/domain"
	"github.com/s4mukka/justinject/internal/context/mocks"
)

type MockedOtelLogger struct {
	domain.OtelLogger
	mock.Mock
}

func TestInitializeLoggerProviderSuccess(t *testing.T) {
	environment := &domain.Environment{
		Instance: "test-instance",
	}
	ctx := context.WithValue(context.Background(), domain.EnvironmentKey, environment)

	mockedLoggerProvider := &mocks.MockedLoggerProvider{}
	mockedLoggerProvider.On("Logger", "test-instance").Return(&MockedOtelLogger{})

	otelInitLogger = func(ctx domain.IContext) (domain.ILoggerProvider, error) {
		return mockedLoggerProvider, nil
	}

	factory := LoggerProviderFactory{}
	loggerProvider, err := factory.InitializeLoggerProvider(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, loggerProvider)

	mockedLoggerProvider.AssertCalled(t, "Logger", "test-instance")
}

func TestInitializeLoggerProviderFailure(t *testing.T) {
	environment := &domain.Environment{
		Instance: "test-instance",
	}
	ctx := context.WithValue(context.Background(), domain.EnvironmentKey, environment)

	otelInitLogger = func(ctx domain.IContext) (domain.ILoggerProvider, error) {
		return nil, errors.New("simulated error")
	}

	factory := LoggerProviderFactory{}
	loggerProvider, err := factory.InitializeLoggerProvider(ctx)

	assert.Error(t, err)
	assert.Nil(t, loggerProvider)
	assert.Equal(t, "simulated error", err.Error())
}
