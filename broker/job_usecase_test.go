package broker

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/s4mukka/justinject/domain"
)

type MockExtractorRepository struct {
	mock.Mock
}

func (m *MockExtractorRepository) GetById(extractorId string) (*domain.Extractor, error) {
	args := m.Called(extractorId)
	return args.Get(0).(*domain.Extractor), args.Error(1)
}

type MockJobRepository struct {
	mock.Mock
}

func (m *MockJobRepository) Create(job domain.IJob) error {
	return m.Called(job).Error(0)
}

type MockK8sRepository struct {
	mock.Mock
}

func (m *MockK8sRepository) CreateJob(job domain.IJob) error {
	return m.Called(job).Error(0)
}

type MockExtractorRepositoryFactory struct {
	mock.Mock
}

func (m *MockExtractorRepositoryFactory) Create() (domain.IExtractorRepository, error) {
	args := m.Called()
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockExtractorRepository), args.Error(1)
}

type MockJobRepositoryFactory struct {
	mock.Mock
}

func (m *MockJobRepositoryFactory) Create() (domain.IJobRepository, error) {
	args := m.Called()
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockJobRepository), args.Error(1)
}

type MockK8sRepositoryFactory struct {
	mock.Mock
}

func (m *MockK8sRepositoryFactory) Create() (domain.IK8sRepository, error) {
	args := m.Called()
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockK8sRepository), args.Error(1)
}

func TestJobUseCase_CreateJob(t *testing.T) {
	mockJobRepository := MockJobRepository{}
	mockExtractorRepository := MockExtractorRepository{}
	mockK8sRepository := MockK8sRepository{}
	type fields struct {
		JobRepository       domain.IJobRepository
		ExtractorRepository domain.IExtractorRepository
		K8sRepository       domain.IK8sRepository
	}
	type args struct {
		request domain.CreateJobRequest
	}
	type mockArgs struct {
		obj     string
		method  string
		args    []interface{}
		returns []interface{}
	}
	request := domain.CreateJobRequest{
		ExtractorId:   "any",
		Query:         "any",
		UpperBound:    100,
		LowerBound:    0,
		NumPartitions: 10,
	}
	job := domain.Job{
		CreateJobRequest: request,
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     domain.IJob
		wantErr  bool
		mockArgs []mockArgs
	}{
		{
			name: "Should returns an error when ExtractorRepository.GetById returns an error",
			fields: fields{
				JobRepository:       &mockJobRepository,
				ExtractorRepository: &mockExtractorRepository,
				K8sRepository:       &mockK8sRepository,
			},
			args: args{
				request: request,
			},
			want:    nil,
			wantErr: true,
			mockArgs: []mockArgs{
				{
					obj:     "ExtractorRepository",
					method:  "GetById",
					args:    []interface{}{mock.Anything},
					returns: []interface{}{&domain.Extractor{Id: request.ExtractorId}, fmt.Errorf("any")},
				},
			},
		},
		{
			name: "Should returns an error when JobRepository.CreateJob returns an error",
			fields: fields{
				JobRepository:       &mockJobRepository,
				ExtractorRepository: &mockExtractorRepository,
				K8sRepository:       &mockK8sRepository,
			},
			args: args{
				request: request,
			},
			want:    nil,
			wantErr: true,
			mockArgs: []mockArgs{
				{
					obj:     "ExtractorRepository",
					method:  "GetById",
					args:    []interface{}{mock.Anything},
					returns: []interface{}{&domain.Extractor{Id: request.ExtractorId}, nil},
				},
				{
					obj:     "JobRepository",
					method:  "Create",
					args:    []interface{}{&job},
					returns: []interface{}{fmt.Errorf("any")},
				},
			},
		},
		{
			name: "Should returns an error when K8sRepository.CreateJob returns an error",
			fields: fields{
				JobRepository:       &mockJobRepository,
				ExtractorRepository: &mockExtractorRepository,
				K8sRepository:       &mockK8sRepository,
			},
			args: args{
				request: request,
			},
			want:    nil,
			wantErr: true,
			mockArgs: []mockArgs{
				{
					obj:     "ExtractorRepository",
					method:  "GetById",
					args:    []interface{}{mock.Anything},
					returns: []interface{}{&domain.Extractor{Id: request.ExtractorId}, nil},
				},
				{
					obj:     "JobRepository",
					method:  "Create",
					args:    []interface{}{&job},
					returns: []interface{}{nil},
				},
				{
					obj:     "K8sRepository",
					method:  "CreateJob",
					args:    []interface{}{&job},
					returns: []interface{}{fmt.Errorf("any")},
				},
			},
		},
		{
			name: "Should returns nil on successful",
			fields: fields{
				JobRepository:       &mockJobRepository,
				ExtractorRepository: &mockExtractorRepository,
				K8sRepository:       &mockK8sRepository,
			},
			args: args{
				request: request,
			},
			want:    &job,
			wantErr: false,
			mockArgs: []mockArgs{
				{
					obj:     "ExtractorRepository",
					method:  "GetById",
					args:    []interface{}{mock.Anything},
					returns: []interface{}{&domain.Extractor{Id: request.ExtractorId}, nil},
				},
				{
					obj:     "JobRepository",
					method:  "Create",
					args:    []interface{}{&job},
					returns: []interface{}{nil},
				},
				{
					obj:     "K8sRepository",
					method:  "CreateJob",
					args:    []interface{}{&job},
					returns: []interface{}{nil},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &JobUseCase{
				JobRepository:       tt.fields.JobRepository,
				ExtractorRepository: tt.fields.ExtractorRepository,
				K8sRepository:       tt.fields.K8sRepository,
			}
			for _, mockArgs := range tt.mockArgs {
				if mockArgs.obj == "JobRepository" {
					uc.JobRepository.(*MockJobRepository).
						On(mockArgs.method, mockArgs.args...).
						Return(mockArgs.returns...).
						Once()
				} else if mockArgs.obj == "K8sRepository" {
					uc.K8sRepository.(*MockK8sRepository).
						On(mockArgs.method, mockArgs.args...).
						Return(mockArgs.returns...).
						Once()
				} else if mockArgs.obj == "ExtractorRepository" {
					uc.ExtractorRepository.(*MockExtractorRepository).
						On(mockArgs.method, mockArgs.args...).
						Return(mockArgs.returns...).
						Once()
				}
			}
			got, err := uc.CreateJob(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("JobUseCase.CreateJob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JobUseCase.CreateJob() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJobUseCaseFactory_Create(t *testing.T) {
	extractorRepoFactory := &MockExtractorRepositoryFactory{}
	jobRepoFactory := &MockJobRepositoryFactory{}
	k8sRepoFactory := &MockK8sRepositoryFactory{}
	extractorRepo := &MockExtractorRepository{}
	jobRepo := &MockJobRepository{}
	k8sRepo := &MockK8sRepository{}
	useCase := &JobUseCase{
		ExtractorRepository: extractorRepo,
		JobRepository:       jobRepo,
		K8sRepository:       k8sRepo,
	}
	type fields struct {
		extractorRepositoryFactory domain.IFactory[domain.IExtractorRepository]
		jobRepositoryFactory       domain.IFactory[domain.IJobRepository]
		k8sRepositoryFactory       domain.IFactory[domain.IK8sRepository]
	}
	type args struct{}
	type mockArgs struct {
		obj     string
		method  string
		args    []interface{}
		returns []interface{}
	}
	testCases := []struct {
		name     string
		fields   fields
		args     args
		want     domain.IJobUseCase
		wantErr  bool
		mockArgs []mockArgs
	}{
		{
			name: "Should return a valid JobUseCase when factories successfully",
			fields: fields{
				extractorRepositoryFactory: extractorRepoFactory,
				jobRepositoryFactory:       jobRepoFactory,
				k8sRepositoryFactory:       k8sRepoFactory,
			},
			mockArgs: []mockArgs{
				{
					obj:     "extractorRepositoryFactory",
					method:  "Create",
					args:    []interface{}{},
					returns: []interface{}{extractorRepo, nil},
				},
				{
					obj:     "jobRepositoryFactory",
					method:  "Create",
					args:    []interface{}{},
					returns: []interface{}{jobRepo, nil},
				},
				{
					obj:     "k8sRepositoryFactory",
					method:  "Create",
					args:    []interface{}{},
					returns: []interface{}{k8sRepo, nil},
				},
			},
			want:    useCase,
			wantErr: false,
		},
		{
			name: "Should return an error if extractorRepositoryFactory fails",
			fields: fields{
				extractorRepositoryFactory: extractorRepoFactory,
				jobRepositoryFactory:       jobRepoFactory,
				k8sRepositoryFactory:       k8sRepoFactory,
			},
			mockArgs: []mockArgs{
				{
					obj:     "extractorRepositoryFactory",
					method:  "Create",
					args:    []interface{}{},
					returns: []interface{}{nil, errors.New("extractor repository error")},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should return an error if jobRepositoryFactory fails",
			fields: fields{
				extractorRepositoryFactory: extractorRepoFactory,
				jobRepositoryFactory:       jobRepoFactory,
				k8sRepositoryFactory:       k8sRepoFactory,
			},
			mockArgs: []mockArgs{
				{
					obj:     "extractorRepositoryFactory",
					method:  "Create",
					args:    []interface{}{},
					returns: []interface{}{extractorRepo, nil},
				},
				{
					obj:     "jobRepositoryFactory",
					method:  "Create",
					args:    []interface{}{},
					returns: []interface{}{nil, errors.New("extractor repository error")},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should return an error if k8sRepositoryFactory fails",
			fields: fields{
				extractorRepositoryFactory: extractorRepoFactory,
				jobRepositoryFactory:       jobRepoFactory,
				k8sRepositoryFactory:       k8sRepoFactory,
			},
			mockArgs: []mockArgs{
				{
					obj:     "extractorRepositoryFactory",
					method:  "Create",
					args:    []interface{}{},
					returns: []interface{}{extractorRepo, nil},
				},
				{
					obj:     "jobRepositoryFactory",
					method:  "Create",
					args:    []interface{}{},
					returns: []interface{}{jobRepo, nil},
				},
				{
					obj:     "k8sRepositoryFactory",
					method:  "Create",
					args:    []interface{}{},
					returns: []interface{}{nil, errors.New("extractor repository error")},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			factory := &JobUseCaseFactory{
				extractorRepositoryFactory: tt.fields.extractorRepositoryFactory,
				jobRepositoryFactory:       tt.fields.jobRepositoryFactory,
				k8sRepositoryFactory:       tt.fields.k8sRepositoryFactory,
			}
			for _, mockArgs := range tt.mockArgs {
				switch mockArgs.obj {
				case "extractorRepositoryFactory":
					tt.fields.extractorRepositoryFactory.(*MockExtractorRepositoryFactory).On(mockArgs.method, mockArgs.args...).Return(mockArgs.returns...).Once()
				case "jobRepositoryFactory":
					tt.fields.jobRepositoryFactory.(*MockJobRepositoryFactory).On(mockArgs.method, mockArgs.args...).Return(mockArgs.returns...).Once()
				case "k8sRepositoryFactory":
					tt.fields.k8sRepositoryFactory.(*MockK8sRepositoryFactory).On(mockArgs.method, mockArgs.args...).Return(mockArgs.returns...).Once()
				}
			}
			got, err := factory.Create()
			if (err != nil) != tt.wantErr {
				t.Errorf("JobUseCaseFactory.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JobUseCaseFactory.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
