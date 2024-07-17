package broker

import (
	"context"
	"testing"

	"github.com/s4mukka/justinject/domain"
	"github.com/s4mukka/justinject/internal/server"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type MockServer struct{}

func (m *MockServer) Init(routes func(server.IRouter), port int) error {
	return nil
}

type MockServerFactory struct{}

func (sf *MockServerFactory) MakeServer(ctx *context.Context) server.IServer {
	return &MockServer{}
}

func TestInit(t *testing.T) {
	ctx := context.Background()
	environment := &domain.Environment{
		Logger: logrus.NewEntry(logrus.New()),
	}
	ctx = context.WithValue(ctx, "environment", environment)

	// Mock server initialization
	serverFactory = &MockServerFactory{}

	err := Init(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, logger)
}
