package mocks

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type MockRouter struct {
	mock.Mock
	gin.IRouter
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
