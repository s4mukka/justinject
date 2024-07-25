package mocks

import (
	"context"

	"github.com/s4mukka/justinject/domain"
	"github.com/stretchr/testify/mock"
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

type MockLoggerProviderFactory struct{}

func (m MockLoggerProviderFactory) InitializeLoggerProvider(ctx domain.IContext) (domain.ILoggerProvider, error) {
	return &MockedLoggerProvider{}, nil
}

type MockLoggerProviderFactoryWithError struct {
	Error error
}

func (m MockLoggerProviderFactoryWithError) InitializeLoggerProvider(ctx domain.IContext) (domain.ILoggerProvider, error) {
	return nil, m.Error
}
