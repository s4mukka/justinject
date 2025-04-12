package mocks

import (
	"github.com/s4mukka/justinject/domain"
)

type MockServer struct{}

func (m *MockServer) Init(routes func(domain.IRouter) error, port int) error {
	return nil
}

type MockServerFactory struct{}

func (sf *MockServerFactory) MakeServer(ctx domain.IContext) domain.IServer {
	return &MockServer{}
}
