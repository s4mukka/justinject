package broker

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/s4mukka/justinject/domain"
)

type MockJobUseCase struct {
	mock.Mock
}

func (m *MockJobUseCase) CreateJob(request domain.CreateJobRequest) (*domain.Job, error) {
	args := m.Called(request)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Job), args.Error(1)
}

type MockGinContext struct {
	mock.Mock
}

func (m *MockGinContext) JSON(code int, obj any) {
	m.Called(code, obj)
}

func (m *MockGinContext) ShouldBindJSON(obj any) error {
	return m.Called(obj).Error(0)
}

func TestCreateJobBadRequestError(t *testing.T) {
	useCase := &MockJobUseCase{}
	jf := &JobServiceFactory{}
	j := jf.MakeJobService(useCase)
	mockError := fmt.Errorf("any error")
	ctx := &MockGinContext{}
	ctx.On("ShouldBindJSON", mock.Anything).Return(mockError)
	ctx.On("JSON", http.StatusBadRequest, mockError.Error())
	j.CreateJob(ctx)
	ctx.AssertExpectations(t)
}

func TestCreateJobInternalServerError(t *testing.T) {
	useCase := &MockJobUseCase{}
	mockError := fmt.Errorf("any error")
	useCase.On("CreateJob", mock.Anything).Return(nil, mockError)
	jf := &JobServiceFactory{}
	j := jf.MakeJobService(useCase)
	ctx := &MockGinContext{}
	ctx.On("ShouldBindJSON", mock.Anything).Return(nil)
	ctx.On("JSON", http.StatusInternalServerError, mockError.Error())
	j.CreateJob(ctx)
	ctx.AssertExpectations(t)
}

func TestCreateJobSuccess(t *testing.T) {
	useCase := &MockJobUseCase{}
	job := domain.Job{}
	useCase.On("CreateJob", mock.Anything).Return(&job, nil)
	jf := &JobServiceFactory{}
	j := jf.MakeJobService(useCase)
	ctx := &MockGinContext{}
	ctx.On("ShouldBindJSON", mock.Anything).Return(nil)
	ctx.On("JSON", http.StatusOK, job)
	j.CreateJob(ctx)
	ctx.AssertExpectations(t)
}

func TestJobServiceFactory_MakeJobService(t *testing.T) {
	useCase := &MockJobUseCase{}
	jf := &JobServiceFactory{}
	j := jf.MakeJobService(useCase)

	assert.NotNil(t, j)

	assert.IsType(t, &JobService{}, j)

	assert.Equal(t, useCase, j.jobUseCase)
}
