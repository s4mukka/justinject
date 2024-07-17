package context

import (
	"context"
	"errors"
	"testing"

	"github.com/s4mukka/justinject/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestInitializeTracerProviderSuccess(t *testing.T) {
	environment := &domain.Environment{
		Instance: "test-instance",
	}
	ctx := context.WithValue(context.Background(), "environment", environment)

	mockedTracerProvider := &MockedTracerProvider{}
	mockedTracerProvider.On("Tracer", "test-instance").Return(&MockedOtelTracer{})

	otelInitTracer = func(ctx *context.Context) (domain.ITracerProvider, error) {
		return mockedTracerProvider, nil
	}

	factory := TracerProviderFactory{}
	tracerProvider, err := factory.InitializeTracerProvider(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, tracerProvider)

	mockedTracerProvider.AssertCalled(t, "Tracer", "test-instance")
}

func TestInitializeTracerProviderFailure(t *testing.T) {
	environment := &domain.Environment{
		Instance: "test-instance",
	}
	ctx := context.WithValue(context.Background(), "environment", environment)

	otelInitTracer = func(ctx *context.Context) (domain.ITracerProvider, error) {
		return nil, errors.New("simulated error")
	}

	factory := TracerProviderFactory{}
	tracerProvider, err := factory.InitializeTracerProvider(ctx)

	assert.Error(t, err)
	assert.Nil(t, tracerProvider)
	assert.Equal(t, "simulated error", err.Error())
}
