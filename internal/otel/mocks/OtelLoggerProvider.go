package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
)

type MockOtelLoggerProvider struct {
	embedded.LoggerProvider
	mock.Mock
}

func (m *MockOtelLoggerProvider) Logger(name string, opts ...log.LoggerOption) log.Logger {
	args := m.Called(name, opts)
	return args.Get(0).(log.Logger)
}

func (m *MockOtelLoggerProvider) Shutdown(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
