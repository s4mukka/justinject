package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/s4mukka/justinject/domain"
)

type MockedTracerProvider struct {
	mock.Mock
}

func (m *MockedTracerProvider) Get() domain.IOtelTracerProvider {
	return nil
}

func (m *MockedTracerProvider) Tracer(name string, opts ...domain.OtelTracerOption) domain.OtelTracer {
	args := m.Called(name)
	return args.Get(0).(domain.OtelTracer)
}

func (m *MockedTracerProvider) Shutdown(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockedOtelTracer struct {
	domain.OtelTracer
	mock.Mock
}

type MockTracerProviderFactory struct{}

func (m MockTracerProviderFactory) InitializeTracerProvider(ctx domain.IContext) (domain.ITracerProvider, error) {
	return &MockedTracerProvider{}, nil
}

type MockTracerProviderFactoryWithError struct {
	Error error
}

func (m MockTracerProviderFactoryWithError) InitializeTracerProvider(ctx domain.IContext) (domain.ITracerProvider, error) {
	return nil, m.Error
}
