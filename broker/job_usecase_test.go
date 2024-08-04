package broker

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/s4mukka/justinject/domain"
)

type MockJobRepository struct {
	mock.Mock
}

func (m *MockJobRepository) GetExtractorById(extractorId string) (*domain.Extractor, error) {
	args := m.Called(extractorId)
	return args.Get(0).(*domain.Extractor), args.Error(1)
}

func (m *MockJobRepository) CreateJob(job *domain.Job) error {
	return m.Called(job).Error(0)
}

type MockK8sRepository struct {
	mock.Mock
}

func (m *MockK8sRepository) CreateJob(job *domain.Job) error {
	return m.Called(job).Error(0)
}

func TestJobUseCase_CreateJob(t *testing.T) {
	mockJobRepository := MockJobRepository{}
	mockK8sRepository := MockK8sRepository{}
	type fields struct {
		JobRepository domain.IJobRepository
		K8sRepository domain.IK8sRepository
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
		want     *domain.Job
		wantErr  bool
		mockArgs []mockArgs
	}{
		{
			name: "Should returns an error when JobRepository.GetExtractorById returns an error",
			fields: fields{
				JobRepository: &mockJobRepository,
				K8sRepository: &mockK8sRepository,
			},
			args: args{
				request: request,
			},
			want:    nil,
			wantErr: true,
			mockArgs: []mockArgs{
				{
					obj:     "JobRepository",
					method:  "GetExtractorById",
					args:    []interface{}{mock.Anything},
					returns: []interface{}{&domain.Extractor{Id: request.ExtractorId}, fmt.Errorf("any")},
				},
			},
		},
		{
			name: "Should returns an error when JobRepository.CreateJob returns an error",
			fields: fields{
				JobRepository: &mockJobRepository,
				K8sRepository: &mockK8sRepository,
			},
			args: args{
				request: request,
			},
			want:    nil,
			wantErr: true,
			mockArgs: []mockArgs{
				{
					obj:     "JobRepository",
					method:  "GetExtractorById",
					args:    []interface{}{mock.Anything},
					returns: []interface{}{&domain.Extractor{Id: request.ExtractorId}, nil},
				},
				{
					obj:     "JobRepository",
					method:  "CreateJob",
					args:    []interface{}{&job},
					returns: []interface{}{fmt.Errorf("any")},
				},
			},
		},
		{
			name: "Should returns an error when K8sRepository.CreateJob returns an error",
			fields: fields{
				JobRepository: &mockJobRepository,
				K8sRepository: &mockK8sRepository,
			},
			args: args{
				request: request,
			},
			want:    nil,
			wantErr: true,
			mockArgs: []mockArgs{
				{
					obj:     "JobRepository",
					method:  "GetExtractorById",
					args:    []interface{}{mock.Anything},
					returns: []interface{}{&domain.Extractor{Id: request.ExtractorId}, nil},
				},
				{
					obj:     "JobRepository",
					method:  "CreateJob",
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
				JobRepository: &mockJobRepository,
				K8sRepository: &mockK8sRepository,
			},
			args: args{
				request: request,
			},
			want:    &job,
			wantErr: false,
			mockArgs: []mockArgs{
				{
					obj:     "JobRepository",
					method:  "GetExtractorById",
					args:    []interface{}{mock.Anything},
					returns: []interface{}{&domain.Extractor{Id: request.ExtractorId}, nil},
				},
				{
					obj:     "JobRepository",
					method:  "CreateJob",
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
				JobRepository: tt.fields.JobRepository,
				K8sRepository: tt.fields.K8sRepository,
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
