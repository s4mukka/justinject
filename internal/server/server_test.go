package server

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/s4mukka/justinject/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServer_Init_Success(t *testing.T) {
	mockLogger := new(MockLogger)
	mockLogger.On("Infof", "Listening and serving HTTP on :%d", []interface{}{8080})

	mockEnv := &domain.Environment{
		Instance: "test-instance",
		Logger:   mockLogger,
	}
	ctx := context.WithValue(context.Background(), "environment", mockEnv)

	server := &Server{Ctx: &ctx}

	intializeRoutes := func(r IRouter) {
		fmt.Printf("R=%v\n", r)
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, "pong")
		})
	}

	router = new(MockRouter)
	router.(*MockRouter).On("Run", []string{":8080"}).Return(nil)
	router.(*MockRouter).On("Use", mock.Anything).Return(router)
	router.(*MockRouter).On("GET", mock.Anything, mock.Anything).Return(router)

	err := server.Init(intializeRoutes, 8080)
	assert.NoError(t, err)

	mockLogger.AssertExpectations(t)
	router.(*MockRouter).AssertExpectations(t)
}

func TestServer_Init_RunError(t *testing.T) {
	mockLogger := new(MockLogger)
	mockLogger.On("Infof", "Listening and serving HTTP on :%d", []interface{}{8080})

	mockEnv := &domain.Environment{
		Instance: "test-instance",
		Logger:   mockLogger,
	}
	ctx := context.WithValue(context.Background(), "environment", mockEnv)

	server := &Server{Ctx: &ctx}

	intializeRoutes := func(r IRouter) {}

	router = new(MockRouter)
	router.(*MockRouter).On("Run", []string{":8080"}).Return(errors.New("server error"))
	router.(*MockRouter).On("Use", mock.Anything).Return(router)

	err := server.Init(intializeRoutes, 8080)
	assert.Error(t, err)
	assert.Equal(t, "server error", err.Error())

	mockLogger.AssertExpectations(t)
	router.(*MockRouter).AssertExpectations(t)
}

func TestServerFactory_MakeServer(t *testing.T) {
	ctx := context.Background()
	sf := &ServerFactory{}

	s := sf.MakeServer(&ctx)

	assert.NotNil(t, s)

	assert.IsType(t, &Server{}, s)

	serverInstance, ok := s.(*Server)
	assert.True(t, ok)
	assert.Equal(t, &ctx, serverInstance.Ctx)
}
