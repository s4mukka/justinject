package otel

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/s4mukka/justinject/domain"
	"github.com/s4mukka/justinject/internal/otel/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
)

func TestInitLogger_Success(t *testing.T) {
	environment := &domain.Environment{
		Instance: "test-instance",
	}
	ctx := context.WithValue(context.Background(), domain.EnvironmentKey, environment)

	os.Setenv("OTEL_ENDPOINT_HTTP", "localhost:4318")
	defer os.Unsetenv("OTEL_ENDPOINT_HTTP")

	loggerProvider, err := InitLogger(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, loggerProvider)

	lp := loggerProvider.(*LoggerProvider)
	assert.NotNil(t, lp.handler)
}

func TestInitLogger_NoOtelEndpoint(t *testing.T) {
	environment := &domain.Environment{
		Instance: "test-instance",
	}
	ctx := context.WithValue(context.Background(), domain.EnvironmentKey, environment)

	os.Unsetenv("OTEL_ENDPOINT_HTTP")

	loggerProvider, err := InitLogger(ctx)

	assert.Error(t, err)
	assert.Nil(t, loggerProvider)
	assert.Equal(t, fmt.Errorf("OTEL_ENDPOINT_HTTP environment variable is not defined"), err)
}

func TestInitLogger_ExporterError(t *testing.T) {
	environment := &domain.Environment{
		Instance: "test-instance",
	}
	ctx := context.WithValue(context.Background(), domain.EnvironmentKey, environment)

	os.Setenv("OTEL_ENDPOINT_HTTP", "invalid-endpoint")
	defer os.Unsetenv("OTEL_ENDPOINT_HTTP")

	originalNew := otlploghttpNew
	defer func() { otlploghttpNew = originalNew }()
	otlploghttpNew = func(ctx context.Context, opts ...otlploghttp.Option) (*otlploghttp.Exporter, error) {
		return nil, fmt.Errorf("fake error")
	}

	loggerProvider, err := InitLogger(ctx)

	assert.Error(t, err)
	assert.Nil(t, loggerProvider)
	assert.Equal(t, fmt.Errorf("Error creating OTLP exporter: fake error\n"), err)
}

func TestLoggerProvider_Get(t *testing.T) {
	provider := new(domain.OtelLoggerProvider)
	loggerProvider := &LoggerProvider{handler: *provider}

	result := loggerProvider.Get()

	assert.Equal(t, *provider, result)
}

func TestLoggerProvider_Logger(t *testing.T) {
	provider := new(mocks.MockOtelLoggerProvider)
	loggerProvider := &LoggerProvider{handler: provider}
	mockLogger := new(mocks.MockLogger)
	provider.On("Logger", "test", mock.Anything).Return(mockLogger)

	result := loggerProvider.Logger("test")

	assert.Equal(t, mockLogger, result)
	provider.AssertCalled(t, "Logger", "test", mock.Anything)
}

func TestLoggerProvider_Shutdown(t *testing.T) {
	provider := new(mocks.MockOtelLoggerProvider)
	loggerProvider := &LoggerProvider{handler: provider}
	provider.On("Shutdown", mock.Anything).Return(nil)

	err := loggerProvider.Shutdown(context.Background())

	assert.NoError(t, err)
	provider.AssertCalled(t, "Shutdown", mock.Anything)
}
