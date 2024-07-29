package server

import (
	"context"
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/s4mukka/justinject/domain"
	"github.com/s4mukka/justinject/internal/server/mocks"
)

func TestServer_Init_Success(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	mockLogger.On("Infof", "Listening and serving HTTP on :%d", []interface{}{8080})

	mockEnv := &domain.Environment{
		Instance: "test-instance",
		Logger:   mockLogger,
	}
	ctx := context.WithValue(context.Background(), domain.EnvironmentKey, mockEnv)

	server := &Server{Ctx: ctx}

	intializeRoutes := func(r domain.IRouter) {
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, "pong")
		})
	}

	router = new(mocks.MockRouter)
	router.(*mocks.MockRouter).On("Run", []string{":8080"}).Return(nil)
	router.(*mocks.MockRouter).On("Use", mock.Anything).Return(router)
	router.(*mocks.MockRouter).On("GET", mock.Anything, mock.Anything).Return(router)

	err := server.Init(intializeRoutes, 8080)
	assert.NoError(t, err)

	mockLogger.AssertExpectations(t)
	router.(*mocks.MockRouter).AssertExpectations(t)
}

func TestServer_Init_RunError(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	mockLogger.On("Infof", "Listening and serving HTTP on :%d", []interface{}{8080})

	mockEnv := &domain.Environment{
		Instance: "test-instance",
		Logger:   mockLogger,
	}
	ctx := context.WithValue(context.Background(), domain.EnvironmentKey, mockEnv)

	server := &Server{Ctx: ctx}

	intializeRoutes := func(r domain.IRouter) {}

	router = new(mocks.MockRouter)
	router.(*mocks.MockRouter).On("Run", []string{":8080"}).Return(errors.New("server error"))
	router.(*mocks.MockRouter).On("Use", mock.Anything).Return(router)

	err := server.Init(intializeRoutes, 8080)
	assert.Error(t, err)
	assert.Equal(t, "server error", err.Error())

	mockLogger.AssertExpectations(t)
	router.(*mocks.MockRouter).AssertExpectations(t)
}

func TestServerFactory_MakeServer(t *testing.T) {
	ctx := context.Background()
	sf := &ServerFactory{}

	s := sf.MakeServer(ctx)

	assert.NotNil(t, s)

	assert.IsType(t, &Server{}, s)

	serverInstance, ok := s.(*Server)
	assert.True(t, ok)
	assert.Equal(t, ctx, serverInstance.Ctx)
}
