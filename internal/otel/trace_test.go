package otel

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/s4mukka/justinject/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/embedded"
)

func TestInitTracer_Success(t *testing.T) {
	environment := &domain.Environment{
		Instance: "test-instance",
	}
	ctx := context.WithValue(context.Background(), "environment", environment)

	os.Setenv("OTEL_ENDPOINT_GRPC", "localhost:4318")
	defer os.Unsetenv("OTEL_ENDPOINT_GRPC")

	tracerProvider, err := InitTracer(&ctx)

	assert.NoError(t, err)
	assert.NotNil(t, tracerProvider)

	tp := tracerProvider.(TracerProvider)
	assert.NotNil(t, tp.handler)
}

func TestInitTracer_NoOtelEndpoint(t *testing.T) {
	environment := &domain.Environment{
		Instance: "test-instance",
	}
	ctx := context.WithValue(context.Background(), "environment", environment)

	os.Unsetenv("OTEL_ENDPOINT_GRPC")

	tracerProvider, err := InitTracer(&ctx)

	assert.Error(t, err)
	assert.Nil(t, tracerProvider)
	assert.Equal(t, fmt.Errorf("OTEL_ENDPOINT_GRPC environment variable is not defined"), err)
}

func TestInitTracer_ExporterError(t *testing.T) {
	environment := &domain.Environment{
		Instance: "test-instance",
	}
	ctx := context.WithValue(context.Background(), "environment", environment)

	os.Setenv("OTEL_ENDPOINT_GRPC", "invalid-endpoint")
	defer os.Unsetenv("OTEL_ENDPOINT_GRPC")

	originalNew := otlptracegrpcNew
	defer func() { otlptracegrpcNew = originalNew }()
	otlptracegrpcNew = func(ctx context.Context, opts ...otlptracegrpc.Option) (*otlptrace.Exporter, error) {
		return nil, fmt.Errorf("mock exporter creation error")
	}

	tracerProvider, err := InitTracer(&ctx)

	assert.Error(t, err)
	assert.Nil(t, tracerProvider)
	assert.Equal(t, fmt.Errorf("Error creating OTLP exporter: mock exporter creation error\n"), err)
}

func TestTracerProvider_Get(t *testing.T) {
	mockProvider := new(domain.OtelTracerProvider)
	tracerProvider := &TracerProvider{handler: *mockProvider}

	result := tracerProvider.Get()

	assert.Equal(t, *mockProvider, result)
}

type MockOtelTracerProvider struct {
	embedded.TracerProvider
	mock.Mock
}

func (m *MockOtelTracerProvider) Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	args := m.Called(name, opts)
	return args.Get(0).(trace.Tracer)
}

func (m *MockOtelTracerProvider) ForceFlush(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockOtelTracerProvider) Shutdown(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockOtelTracerProvider) RegisterSpanProcessor(sp sdktrace.SpanProcessor) {
	m.Called(sp)
}

func (m *MockOtelTracerProvider) UnregisterSpanProcessor(sp sdktrace.SpanProcessor) {
	m.Called(sp)
}

type MockTracer struct {
	embedded.Tracer
	mock.Mock
}

func (m *MockTracer) Emit(ctx context.Context, record log.Record) {
	m.Called(ctx, record)
}

func (m *MockTracer) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	args := m.Called(ctx, spanName, opts)
	return args.Get(0).(context.Context), args.Get(1).(trace.Span)
}

func (m *MockTracer) Enabled(ctx context.Context, record log.Record) bool {
	args := m.Called(ctx, record)
	return args.Bool(0)
}

func TestTracerProvider_Tracer(t *testing.T) {
	mockProvider := new(MockOtelTracerProvider)
	tracerProvider := &TracerProvider{handler: mockProvider}
	mockTracer := new(MockTracer)
	mockProvider.On("Tracer", "test", mock.Anything).Return(mockTracer)

	result := tracerProvider.Tracer("test")

	assert.Equal(t, mockTracer, result)
	mockProvider.AssertCalled(t, "Tracer", "test", mock.Anything)
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

func TestTracerProvider_RegisterSpanProcessor(t *testing.T) {
	mockProcessor := new(MockSpanProcessor)
	mockProvider := new(MockOtelTracerProvider)
	tracerProvider := &TracerProvider{handler: mockProvider}

	mockProvider.On("RegisterSpanProcessor", mockProcessor)

	tracerProvider.RegisterSpanProcessor(mockProcessor)

	mockProvider.AssertCalled(t, "RegisterSpanProcessor", mockProcessor)
}

func TestTracerProvider_UnregisterSpanProcessor(t *testing.T) {
	mockProcessor := new(MockSpanProcessor)
	mockProvider := new(MockOtelTracerProvider)
	tracerProvider := &TracerProvider{handler: mockProvider}

	mockProvider.On("UnregisterSpanProcessor", mockProcessor)

	tracerProvider.UnregisterSpanProcessor(mockProcessor)

	mockProvider.AssertCalled(t, "UnregisterSpanProcessor", mockProcessor)
}

func TestTracerProvider_ForceFlush(t *testing.T) {
	mockProvider := new(MockOtelTracerProvider)
	tracerProvider := &TracerProvider{handler: mockProvider}
	mockProvider.On("ForceFlush", mock.Anything).Return(nil)

	err := tracerProvider.ForceFlush(context.Background())

	assert.NoError(t, err)
	mockProvider.AssertCalled(t, "ForceFlush", mock.Anything)
}

func TestTracerProvider_Shutdown(t *testing.T) {
	mockProvider := new(MockOtelTracerProvider)
	tracerProvider := &TracerProvider{handler: mockProvider}
	mockProvider.On("Shutdown", mock.Anything).Return(nil)

	err := tracerProvider.Shutdown(context.Background())

	assert.NoError(t, err)
	mockProvider.AssertCalled(t, "Shutdown", mock.Anything)
}
