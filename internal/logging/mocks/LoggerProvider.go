package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/s4mukka/justinject/domain"
)

type MockedLoggerProvider struct {
	mock.Mock
}

func (m *MockedLoggerProvider) Get() domain.IOtelLoggerProvider {
	return nil
}

func (m *MockedLoggerProvider) Logger(name string, opts ...domain.OtelLoggerOption) domain.OtelLogger {
	args := m.Called(name)
	return args.Get(0).(domain.OtelLogger)
}

func (m *MockedLoggerProvider) Shutdown(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
