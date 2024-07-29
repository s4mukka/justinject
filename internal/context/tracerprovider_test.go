package context

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/s4mukka/justinject/domain"
	"github.com/s4mukka/justinject/internal/context/mocks"
)

func TestInitializeTracerProviderSuccess(t *testing.T) {
	environment := &domain.Environment{
		Instance: "test-instance",
	}
	ctx := context.WithValue(context.Background(), domain.EnvironmentKey, environment)

	mockedTracerProvider := &mocks.MockedTracerProvider{}
	mockedTracerProvider.On("Tracer", "test-instance").Return(&mocks.MockedOtelTracer{})

	otelInitTracer = func(ctx domain.IContext) (domain.ITracerProvider, error) {
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
	ctx := context.WithValue(context.Background(), domain.EnvironmentKey, environment)

	otelInitTracer = func(ctx domain.IContext) (domain.ITracerProvider, error) {
		return nil, errors.New("simulated error")
	}

	factory := TracerProviderFactory{}
	tracerProvider, err := factory.InitializeTracerProvider(ctx)

	assert.Error(t, err)
	assert.Nil(t, tracerProvider)
	assert.Equal(t, "simulated error", err.Error())
}
