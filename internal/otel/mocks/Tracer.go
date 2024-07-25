package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/embedded"
)

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
