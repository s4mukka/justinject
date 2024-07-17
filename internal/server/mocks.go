package server

import (
	"github.com/gin-gonic/gin"
	"github.com/s4mukka/justinject/domain"
	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
	domain.ILogger
}

func (m *MockLogger) Infof(format string, args ...interface{}) {
	m.Called(format, args)
}

type MockRouter struct {
	mock.Mock
	IRouter
}

func (m *MockRouter) Run(addr ...string) (err error) {
	args := m.Called(addr)
	return args.Error(0)
}

func (m *MockRouter) Use(middleware ...gin.HandlerFunc) gin.IRoutes {
	args := m.Called(middleware)
	return args.Get(0).(gin.IRoutes)
}

func (m *MockRouter) GET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	args := m.Called(relativePath, handlers)
	return args.Get(0).(gin.IRoutes)
}
