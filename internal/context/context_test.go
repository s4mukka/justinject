package context

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/s4mukka/justinject/domain"
	"github.com/s4mukka/justinject/mock"
	m "github.com/s4mukka/justinject/mock"
	"github.com/stretchr/testify/assert"
)

type MockLoggerProviderFactory struct{}

func (l MockLoggerProviderFactory) InitializeLoggerProvider(ctx context.Context) (domain.ILoggerProvider, error) {
	return &m.MockedLoggerProvider{}, nil
}

type MockLoggerProviderFactoryWithError struct{}

func (l MockLoggerProviderFactoryWithError) InitializeLoggerProvider(ctx context.Context) (domain.ILoggerProvider, error) {
	return nil, fmt.Errorf("simulated error")
}

type MockTracerProviderFactory struct{}

func (l MockTracerProviderFactory) InitializeTracerProvider(ctx context.Context) (domain.ITracerProvider, error) {
	return &MockedTracerProvider{}, nil
}

type MockTracerProviderFactoryWithError struct{}

func (l MockTracerProviderFactoryWithError) InitializeTracerProvider(ctx context.Context) (domain.ITracerProvider, error) {
	return nil, fmt.Errorf("simulated error")
}

func TestInitializeContext(t *testing.T) {
	instance := "anything"

	loggerProviderFactory = MockLoggerProviderFactory{}
	tracerProviderFactory = MockTracerProviderFactory{}

	oldAddOtelHook := addOtelHook
	defer func() { addOtelHook = oldAddOtelHook }()
	addOtelHook = func(ctx *context.Context) {}

	ctx := InitializeContext(instance)

	environment := ctx.Value("environment").(*domain.Environment)
	assert.NotNil(t, environment)
	assert.Equal(t, instance, environment.Instance)
	assert.NotNil(t, environment.Logger)
	assert.NotNil(t, environment.LoggerProvider)
	assert.NotNil(t, environment.TracerProvider)
}

func TestInitializeContextWithError(t *testing.T) {
	instance := "anything"

	loggerProviderFactory = MockLoggerProviderFactoryWithError{}
	tracerProviderFactory = MockTracerProviderFactoryWithError{}

	oldAddOtelHook := addOtelHook
	defer func() { addOtelHook = oldAddOtelHook }()
	addOtelHook = func(ctx *context.Context) {}

	out := mock.CaptureOutput(func() { InitializeContext(instance) })

	var logMessages []map[string]interface{}
	for _, line := range bytes.Split([]byte(out), []byte{'\n'}) {
		if len(line) > 0 {
			var msg map[string]interface{}
			err := json.Unmarshal(line, &msg)
			if err != nil {
				t.Fatalf("Failed to unmarshal log line: %v", err)
			}
			logMessages = append(logMessages, msg)
		}
	}

	assert.Len(t, logMessages, 2)

	expectedMessages := []string{
		"Error initializing logger: simulated error",
		"Error initializing tracer: simulated error",
	}

	for i, expected := range expectedMessages {
		assert.Contains(t, logMessages[i]["msg"], expected)
	}
}

func TestShutdownComponents(t *testing.T) {
	instance := "anything"

	loggerProviderFactory = MockLoggerProviderFactory{}
	tracerProviderFactory = MockTracerProviderFactory{}

	oldAddOtelHook := addOtelHook
	defer func() { addOtelHook = oldAddOtelHook }()
	addOtelHook = func(ctx *context.Context) {}

	ctx := InitializeContext(instance)
	environment := ctx.Value("environment").(*domain.Environment)

	mockLoggerProvider := environment.LoggerProvider.(*m.MockedLoggerProvider)
	mockTracerProvider := environment.TracerProvider.(*MockedTracerProvider)

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

	loggerProviderFactory = MockLoggerProviderFactory{}
	tracerProviderFactory = MockTracerProviderFactory{}

	oldAddOtelHook := addOtelHook
	defer func() { addOtelHook = oldAddOtelHook }()
	addOtelHook = func(ctx *context.Context) {}

	ctx := InitializeContext(instance)
	environment := ctx.Value("environment").(*domain.Environment)

	mockLoggerProvider := environment.LoggerProvider.(*m.MockedLoggerProvider)
	mockTracerProvider := environment.TracerProvider.(*MockedTracerProvider)

	mockLoggerProvider.On("Shutdown", context.Background()).Return(fmt.Errorf("simulated shutdown error"))
	mockTracerProvider.On("Shutdown", context.Background()).Return(fmt.Errorf("simulated shutdown error"))

	output := mock.CaptureLoggerOutput(environment.Logger, func() { ShutdownComponents(ctx) })

	var logMessages []map[string]interface{}
	for _, line := range bytes.Split([]byte(output), []byte{'\n'}) {
		if len(line) > 0 {
			var msg map[string]interface{}
			err := json.Unmarshal(line, &msg)
			if err != nil {
				t.Fatalf("Failed to unmarshal log line: %v", err)
			}
			logMessages = append(logMessages, msg)
		}
	}

	assert.Len(t, logMessages, 2)

	expectedMessages := []string{
		"Error closing logger: simulated shutdown error",
		"Error closing tracer: simulated shutdown error",
	}

	for i, expected := range expectedMessages {
		assert.Contains(t, logMessages[i]["msg"], expected)
	}

	mockLoggerProvider.AssertExpectations(t)
	mockTracerProvider.AssertExpectations(t)
}
