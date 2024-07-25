package context

import (
	"context"
	"fmt"
	"testing"

	"github.com/s4mukka/justinject/domain"
	"github.com/s4mukka/justinject/internal/context/mocks"
	"github.com/stretchr/testify/assert"
)

func TestInitializeContext(t *testing.T) {
	instance := "anything"

	loggerProviderFactory = mocks.MockLoggerProviderFactory{}
	tracerProviderFactory = mocks.MockTracerProviderFactory{}

	oldAddOtelHook := addOtelHook
	defer func() { addOtelHook = oldAddOtelHook }()
	addOtelHook = func(ctx domain.IContext) {}

	ctx := InitializeContext(instance)

	environment := ctx.Value(domain.EnvironmentKey).(*domain.Environment)
	assert.NotNil(t, environment)
	assert.Equal(t, instance, environment.Instance)
	assert.NotNil(t, environment.Logger)
	assert.NotNil(t, environment.LoggerProvider)
	assert.NotNil(t, environment.TracerProvider)
}

func TestInitializeContextWithError(t *testing.T) {
	instance := "anything"

	loggerProviderFactory = mocks.MockLoggerProviderFactoryWithError{
		Error: fmt.Errorf("any error"),
	}
	tracerProviderFactory = mocks.MockTracerProviderFactoryWithError{
		Error: fmt.Errorf("any error"),
	}

	oldAddOtelHook := addOtelHook
	defer func() { addOtelHook = oldAddOtelHook }()
	addOtelHook = func(ctx domain.IContext) {}

	mockLogger := new(mocks.MockLogger)
	logInit = func(ctx domain.IContext) domain.ILogger { return mockLogger }
	mockLogger.On(
		"Warnf",
		"Error initializing logger: %v\n",
		[]interface{}{loggerProviderFactory.(mocks.MockLoggerProviderFactoryWithError).Error},
	)
	mockLogger.On(
		"Warnf",
		"Error initializing tracer: %v\n",
		[]interface{}{tracerProviderFactory.(mocks.MockTracerProviderFactoryWithError).Error},
	)

	InitializeContext(instance)

	mockLogger.AssertExpectations(t)
}

func TestShutdownComponents(t *testing.T) {
	instance := "anything"

	loggerProviderFactory = mocks.MockLoggerProviderFactory{}
	tracerProviderFactory = mocks.MockTracerProviderFactory{}

	oldAddOtelHook := addOtelHook
	defer func() { addOtelHook = oldAddOtelHook }()
	addOtelHook = func(ctx domain.IContext) {}

	ctx := InitializeContext(instance)
	environment := ctx.Value(domain.EnvironmentKey).(*domain.Environment)

	mockLoggerProvider := environment.LoggerProvider.(*mocks.MockedLoggerProvider)
	mockTracerProvider := environment.TracerProvider.(*mocks.MockedTracerProvider)

	mockLoggerProvider.On("Shutdown", context.Background()).Return(nil)
	mockTracerProvider.On("Shutdown", context.Background()).Return(nil)

	ShutdownComponents(ctx)

	// Assert expectations
	mockLoggerProvider.AssertExpectations(t)
	mockTracerProvider.AssertExpectations(t)
}

func TestShutdownComponentsNilContext(t *testing.T) {
	assert.NotPanics(t, func() {
		ShutdownComponents(nil)
	})
}

func TestShutdownComponentsWithError(t *testing.T) {
	instance := "anything"

	mockError := fmt.Errorf("any error")

	loggerProviderFactory = mocks.MockLoggerProviderFactory{}
	tracerProviderFactory = mocks.MockTracerProviderFactory{}

	oldAddOtelHook := addOtelHook
	defer func() { addOtelHook = oldAddOtelHook }()
	addOtelHook = func(ctx domain.IContext) {}

	mockLogger := new(mocks.MockLogger)
	logInit = func(ctx domain.IContext) domain.ILogger { return mockLogger }
	mockLogger.On(
		"Warnf",
		"Error closing logger: %v\n",
		[]interface{}{mockError},
	)
	mockLogger.On(
		"Warnf",
		"Error closing tracer: %v\n",
		[]interface{}{mockError},
	)
	ctx := InitializeContext(instance)
	environment := ctx.Value(domain.EnvironmentKey).(*domain.Environment)

	mockLoggerProvider := environment.LoggerProvider.(*mocks.MockedLoggerProvider)
	mockTracerProvider := environment.TracerProvider.(*mocks.MockedTracerProvider)

	mockLoggerProvider.On("Shutdown", context.Background()).Return(mockError)
	mockTracerProvider.On("Shutdown", context.Background()).Return(mockError)

	ShutdownComponents(ctx)

	mockLoggerProvider.AssertExpectations(t)
	mockTracerProvider.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
