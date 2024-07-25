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
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestInitTracer_Success(t *testing.T) {
	environment := &domain.Environment{
		Instance: "test-instance",
	}
	ctx := context.WithValue(context.Background(), domain.EnvironmentKey, environment)

	os.Setenv("OTEL_ENDPOINT_GRPC", "localhost:4318")
	defer os.Unsetenv("OTEL_ENDPOINT_GRPC")

	tracerProvider, err := InitTracer(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, tracerProvider)

	tp := tracerProvider.(TracerProvider)
	assert.NotNil(t, tp.handler)
}

func TestInitTracer_NoOtelEndpoint(t *testing.T) {
	environment := &domain.Environment{
		Instance: "test-instance",
	}
	ctx := context.WithValue(context.Background(), domain.EnvironmentKey, environment)

	os.Unsetenv("OTEL_ENDPOINT_GRPC")

	tracerProvider, err := InitTracer(ctx)

	assert.Error(t, err)
	assert.Nil(t, tracerProvider)
	assert.Equal(t, fmt.Errorf("OTEL_ENDPOINT_GRPC environment variable is not defined"), err)
}

func TestInitTracer_ExporterError(t *testing.T) {
	environment := &domain.Environment{
		Instance: "test-instance",
	}
	ctx := context.WithValue(context.Background(), domain.EnvironmentKey, environment)

	os.Setenv("OTEL_ENDPOINT_GRPC", "invalid-endpoint")
	defer os.Unsetenv("OTEL_ENDPOINT_GRPC")

	originalNew := otlptracegrpcNew
	defer func() { otlptracegrpcNew = originalNew }()
	otlptracegrpcNew = func(ctx context.Context, opts ...otlptracegrpc.Option) (*otlptrace.Exporter, error) {
		return nil, fmt.Errorf("mock exporter creation error")
	}

	tracerProvider, err := InitTracer(ctx)

	assert.Error(t, err)
	assert.Nil(t, tracerProvider)
	assert.Equal(t, fmt.Errorf("Error creating OTLP exporter: mock exporter creation error\n"), err)
}

func TestTracerProvider_Get(t *testing.T) {
	provider := new(domain.OtelTracerProvider)
	tracerProvider := &TracerProvider{handler: *provider}

	result := tracerProvider.Get()

	assert.Equal(t, *provider, result)
}

func TestTracerProvider_Tracer(t *testing.T) {
	provider := new(mocks.MockOtelTracerProvider)
	tracerProvider := &TracerProvider{handler: provider}
	tracer := new(mocks.MockTracer)
	provider.On("Tracer", "test", mock.Anything).Return(tracer)

	result := tracerProvider.Tracer("test")

	assert.Equal(t, tracer, result)
	provider.AssertCalled(t, "Tracer", "test", mock.Anything)
}

type MockSpanProcessor struct {
	mock.Mock
}

func (m *MockSpanProcessor) OnStart(parent context.Context, s sdktrace.ReadWriteSpan) {
	m.Called(parent, s)
}

func (m *MockSpanProcessor) OnEnd(s sdktrace.ReadOnlySpan) {
	m.Called(s)
}

func (m *MockSpanProcessor) Shutdown(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockSpanProcessor) ForceFlush(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestTracerProvider_Shutdown(t *testing.T) {
	provider := new(mocks.MockOtelTracerProvider)
	tracerProvider := &TracerProvider{handler: provider}
	provider.On("Shutdown", mock.Anything).Return(nil)

	err := tracerProvider.Shutdown(context.Background())

	assert.NoError(t, err)
	provider.AssertCalled(t, "Shutdown", mock.Anything)
}
