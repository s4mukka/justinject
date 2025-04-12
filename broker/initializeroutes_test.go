package broker

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/s4mukka/justinject/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockJobService struct{}

func (m MockJobService) CreateJob(ctx domain.IRestContext) {
	ctx.JSON(200, "ok")
}

type MockRouterFactory struct {
	mock.Mock
}

func (m *MockRouterFactory) Create() (domain.IJobService, error) {
	args := m.Called()
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockJobService), args.Error(1)
}

func TestIntializeRoutes(t *testing.T) {
	router := gin.Default()
	jobServiceFactory = &MockRouterFactory{}
	jobServiceFactory.(*MockRouterFactory).On("Create").Return(&MockJobService{}, nil)
	err := intializeRoutes(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `"pong"`, w.Body.String())
	assert.Nil(t, err)
}

func TestIntializeRoutesError(t *testing.T) {
	router := gin.Default()
	jobServiceFactory = &MockRouterFactory{}
	jobServiceFactory.(*MockRouterFactory).On("Create").Return(nil, errors.New("job service factory error"))
	err := intializeRoutes(router)

	assert.Equal(t, "job service factory error", err.Error())
	jobServiceFactory.(*MockRouterFactory).AssertExpectations(t)
}

func exampleHandler(ctx domain.IRestContext) {
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func TestRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/test", route(exampleHandler))

	req, err := http.NewRequest(http.MethodGet, "/test", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expectedBody := `{"message":"success"}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}
