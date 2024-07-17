package context

import (
	"context"
	"errors"
	"testing"

	"github.com/s4mukka/justinject/domain"
	m "github.com/s4mukka/justinject/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedOtelLogger struct {
	domain.OtelLogger
	mock.Mock
}

func TestInitializeLoggerProviderSuccess(t *testing.T) {
	environment := &domain.Environment{
		Instance: "test-instance",
	}
	ctx := context.WithValue(context.Background(), "environment", environment)

	mockedLoggerProvider := &m.MockedLoggerProvider{}
	mockedLoggerProvider.On("Logger", "test-instance").Return(&MockedOtelLogger{})

	otelInitLogger = func(ctx *context.Context) (domain.ILoggerProvider, error) {
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
	ctx := context.WithValue(context.Background(), "environment", environment)

	otelInitLogger = func(ctx *context.Context) (domain.ILoggerProvider, error) {
		return nil, errors.New("simulated error")
	}

	factory := LoggerProviderFactory{}
	loggerProvider, err := factory.InitializeLoggerProvider(ctx)

	assert.Error(t, err)
	assert.Nil(t, loggerProvider)
	assert.Equal(t, "simulated error", err.Error())
}
