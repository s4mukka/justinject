package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
)

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
