package broker

import (
	"errors"
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

func (m *MockJobUseCase) CreateJob(request domain.CreateJobRequest) (domain.IJob, error) {
	args := m.Called(request)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(domain.IJob), args.Error(1)
}

type MockJobUseCaseFactory struct {
	mock.Mock
}

func (m *MockJobUseCaseFactory) Create() (domain.IJobUseCase, error) {
	args := m.Called()
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockJobUseCase), args.Error(1)
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
	jucf := &MockJobUseCaseFactory{}
	jucf.On("Create").Return(useCase, nil)
	jsf := &JobServiceFactory{jucf}
	js, _ := jsf.Create()
	mockError := fmt.Errorf("any error")
	ctx := &MockGinContext{}
	ctx.On("ShouldBindJSON", mock.Anything).Return(mockError)
	ctx.On("JSON", http.StatusBadRequest, mockError.Error())
	js.CreateJob(ctx)
	ctx.AssertExpectations(t)
}

func TestCreateJobInternalServerError(t *testing.T) {
	useCase := &MockJobUseCase{}
	mockError := fmt.Errorf("any error")
	useCase.On("CreateJob", mock.Anything).Return(nil, mockError)
	jucf := &MockJobUseCaseFactory{}
	jucf.On("Create").Return(useCase, nil)
	jsf := &JobServiceFactory{jucf}
	js, _ := jsf.Create()
	ctx := &MockGinContext{}
	ctx.On("ShouldBindJSON", mock.Anything).Return(nil)
	ctx.On("JSON", http.StatusInternalServerError, mockError.Error())
	js.CreateJob(ctx)
	ctx.AssertExpectations(t)
}

func TestCreateJobSuccess(t *testing.T) {
	useCase := &MockJobUseCase{}
	job := domain.Job{}
	useCase.On("CreateJob", mock.Anything).Return(job, nil)
	jucf := &MockJobUseCaseFactory{}
	jucf.On("Create").Return(useCase, nil)
	jsf := &JobServiceFactory{jucf}
	js, _ := jsf.Create()
	ctx := &MockGinContext{}
	ctx.On("ShouldBindJSON", mock.Anything).Return(nil)
	ctx.On("JSON", http.StatusOK, job)
	js.CreateJob(ctx)
	ctx.AssertExpectations(t)
}

func TestJobServiceFactory_Create(t *testing.T) {
	useCase := &MockJobUseCase{}
	jucf := &MockJobUseCaseFactory{}
	jucf.On("Create").Return(useCase, nil)
	jsf := &JobServiceFactory{jucf}
	js, _ := jsf.Create()

	assert.NotNil(t, js)

	assert.IsType(t, &JobService{}, js)
}

func TestJobServiceFactory_CreateError(t *testing.T) {
	jucf := &MockJobUseCaseFactory{}
	jucf.On("Create").Return(nil, errors.New("job use case error"))
	jsf := &JobServiceFactory{jucf}
	js, _ := jsf.Create()

	assert.Nil(t, js)

	jucf.AssertExpectations(t)
}
