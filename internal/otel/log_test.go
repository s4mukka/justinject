package otel

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/s4mukka/justinject/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
)

func TestInitLogger_Success(t *testing.T) {
	environment := &domain.Environment{
		Instance: "test-instance",
	}
	ctx := context.WithValue(context.Background(), "environment", environment)

	os.Setenv("OTEL_ENDPOINT_HTTP", "localhost:4318")
	defer os.Unsetenv("OTEL_ENDPOINT_HTTP")

	loggerProvider, err := InitLogger(&ctx)

	assert.NoError(t, err)
	assert.NotNil(t, loggerProvider)

	lp := loggerProvider.(*LoggerProvider)
	assert.NotNil(t, lp.handler)
}

func TestInitLogger_NoOtelEndpoint(t *testing.T) {
	environment := &domain.Environment{
		Instance: "test-instance",
	}
	ctx := context.WithValue(context.Background(), "environment", environment)

	os.Unsetenv("OTEL_ENDPOINT_HTTP")

	loggerProvider, err := InitLogger(&ctx)

	assert.Error(t, err)
	assert.Nil(t, loggerProvider)
	assert.Equal(t, fmt.Errorf("OTEL_ENDPOINT_HTTP environment variable is not defined"), err)
}

func TestInitLogger_ExporterError(t *testing.T) {
	environment := &domain.Environment{
		Instance: "test-instance",
	}
	ctx := context.WithValue(context.Background(), "environment", environment)

	os.Setenv("OTEL_ENDPOINT_HTTP", "invalid-endpoint")
	defer os.Unsetenv("OTEL_ENDPOINT_HTTP")

	originalNew := otlploghttpNew
	defer func() { otlploghttpNew = originalNew }()
	otlploghttpNew = func(ctx context.Context, opts ...otlploghttp.Option) (*otlploghttp.Exporter, error) {
		return nil, fmt.Errorf("mock exporter creation error")
	}

	loggerProvider, err := InitLogger(&ctx)

	assert.Error(t, err)
	assert.Nil(t, loggerProvider)
	assert.Equal(t, fmt.Errorf("Error creating OTLP exporter: mock exporter creation error\n"), err)
}

func TestLoggerProvider_Get(t *testing.T) {
	mockProvider := new(domain.OtelLoggerProvider)
	loggerProvider := &LoggerProvider{handler: *mockProvider}

	result := loggerProvider.Get()

	assert.Equal(t, *mockProvider, result)
}

type MockOtelLoggerProvider struct {
	embedded.LoggerProvider
	mock.Mock
}

func (m *MockOtelLoggerProvider) Logger(name string, opts ...log.LoggerOption) log.Logger {
	args := m.Called(name, opts)
	return args.Get(0).(log.Logger)
}

func (m *MockOtelLoggerProvider) ForceFlush(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockOtelLoggerProvider) Shutdown(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockLogger struct {
	embedded.Logger
	mock.Mock
}

func (m *MockLogger) Emit(ctx context.Context, record log.Record) {
	m.Called(ctx, record)
}

func (m *MockLogger) Enabled(ctx context.Context, record log.Record) bool {
	args := m.Called(ctx, record)
	return args.Bool(0)
}

func TestLoggerProvider_Logger(t *testing.T) {
	mockProvider := new(MockOtelLoggerProvider)
	loggerProvider := &LoggerProvider{handler: mockProvider}
	mockLogger := new(MockLogger)
	mockProvider.On("Logger", "test", mock.Anything).Return(mockLogger)

	result := loggerProvider.Logger("test")

	assert.Equal(t, mockLogger, result)
	mockProvider.AssertCalled(t, "Logger", "test", mock.Anything)
}

func TestLoggerProvider_Shutdown(t *testing.T) {
	mockProvider := new(MockOtelLoggerProvider)
	loggerProvider := &LoggerProvider{handler: mockProvider}
	mockProvider.On("Shutdown", mock.Anything).Return(nil)

	err := loggerProvider.Shutdown(context.Background())

	assert.NoError(t, err)
	mockProvider.AssertCalled(t, "Shutdown", mock.Anything)
}
