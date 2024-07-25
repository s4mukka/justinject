package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/embedded"
)

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
