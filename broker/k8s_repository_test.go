package broker

import (
	"crypto/x509"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/s4mukka/justinject/domain"
	"github.com/s4mukka/justinject/internal/utils"
	"github.com/stretchr/testify/mock"
)

type MockHttpClient struct {
	mock.Mock
}

func (m *MockHttpClient) HTTPRequest(url, method string, headers map[string]string, body []byte, caCert ...*x509.CertPool) (domain.HTTPResponse, error) {
	listArgs := make([]interface{}, len(caCert)+4)
	listArgs[0] = url
	listArgs[1] = method
	listArgs[2] = headers
	listArgs[3] = body
	for i, v := range caCert {
		listArgs[i+4] = v
	}
	args := m.Called(listArgs...)
	return args.Get(0).(domain.HTTPResponse), args.Error(1)
}

type MockJob struct {
	mock.Mock
}

func (m *MockJob) ParseTemplate() *domain.JobTemplate {
	return m.Called().Get(0).(*domain.JobTemplate)
}

type MockUtils struct {
	mock.Mock
	utils.Utils
}

func (m *MockUtils) TemplateToString(path string, data any) (string, error) {
	args := m.Called(path, data)
	return args.String(0), args.Error(1)
}

func (m *MockUtils) ReadFile(name string) ([]byte, error) {
	args := m.Called(name)
	return args.Get(0).([]byte), args.Error(1)
}

func TestK8sRepository_CreateJob(t *testing.T) {
	jobTemplate := &domain.JobTemplate{}
	caCert := x509.NewCertPool()
	type fields struct {
		httpClient domain.IHTTPClient

		namespace string
		token     string
		caCert    *x509.CertPool
	}
	type args struct {
		job domain.IJob
	}
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
		want     domain.IJob
		wantErr  bool
		mockArgs []mockArgs
	}{
		{
			name: "Should returns nil on successful",
			fields: fields{
				httpClient: &MockHttpClient{},
				namespace:  "default",
				token:      "mockToken",
				caCert:     caCert,
			},
			args: args{
				job: &MockJob{},
			},
			wantErr: false,
			mockArgs: []mockArgs{
				{
					obj:     "Job",
					method:  "ParseTemplate",
					args:    []interface{}{},
					returns: []interface{}{jobTemplate},
				},
				{
					obj:     "u",
					method:  "TemplateToString",
					args:    []interface{}{tmplPath, jobTemplate},
					returns: []interface{}{"templateData", nil},
				},
				{
					obj:    "HttpClient",
					method: "HTTPRequest",
					args: []interface{}{
						"https://kubernetes.default.svc/apis/batch/v1/namespaces/default/jobs",
						"POST",
						map[string]string{
							"Authorization": "Bearer mockToken",
							"Content-Type":  "application/yaml",
						},
						[]byte("templateData"),
						caCert,
					},
					returns: []interface{}{
						domain.HTTPResponse{},
						nil,
					},
				},
			},
		},
		{
			name: "Shoulds return an error when TemplateToString returns an error",
			fields: fields{
				httpClient: &MockHttpClient{},
				namespace:  "default",
				token:      "mockToken",
				caCert:     x509.NewCertPool(),
			},
			args: args{
				job: &MockJob{},
			},
			wantErr: true,
			mockArgs: []mockArgs{
				{
					obj:     "Job",
					method:  "ParseTemplate",
					args:    []interface{}{},
					returns: []interface{}{jobTemplate},
				},
				{
					obj:     "u",
					method:  "TemplateToString",
					args:    []interface{}{tmplPath, jobTemplate},
					returns: []interface{}{"", errors.New("template to string error")},
				},
			},
		},
		{
			name: "Shoulds return an error when HTTPRequest returns an error",
			fields: fields{
				httpClient: &MockHttpClient{},
				namespace:  "default",
				token:      "mockToken",
				caCert:     x509.NewCertPool(),
			},
			args: args{
				job: &MockJob{},
			},
			wantErr: true,
			mockArgs: []mockArgs{
				{
					obj:     "Job",
					method:  "ParseTemplate",
					args:    []interface{}{},
					returns: []interface{}{jobTemplate},
				},
				{
					obj:    "HttpClient",
					method: "HTTPRequest",
					args: []interface{}{
						"https://kubernetes.default.svc/apis/batch/v1/namespaces/default/jobs",
						"POST",
						mock.Anything,
						[]byte("templateData"),
						x509.NewCertPool(),
					},
					returns: []interface{}{
						nil,
						errors.New("HTTP request failed"),
					},
				},
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			repo := &K8sRepository{
				httpClient: tt.fields.httpClient,
				namespace:  tt.fields.namespace,
				token:      tt.fields.token,
				caCert:     tt.fields.caCert,
			}
			for _, mockArgs := range tt.mockArgs {
				switch mockArgs.obj {
				case "Job":
					tt.args.job.(*MockJob).On(mockArgs.method, mockArgs.args...).Return(mockArgs.returns...)
				case "HttpClient":
					repo.httpClient.(*MockHttpClient).On(mockArgs.method, mockArgs.args...).Return(mockArgs.returns...)
				case "u":
					u = &MockUtils{}
					u.(*MockUtils).On(mockArgs.method, mockArgs.args...).Return(mockArgs.returns...)
				}
			}
			err := repo.CreateJob(tt.args.job)
			if (err != nil) != tt.wantErr {
				t.Errorf("K8sRepository.CreateJob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestK8sRepositoryFactory_Create(t *testing.T) {
	mockNamespace := []byte("mockNamespace")
	mockToken := []byte("mockToken")
	mockCacert := []byte("mockCert")
	caCert := x509.NewCertPool()
	caCert.AppendCertsFromPEM(mockCacert)
	repo := &K8sRepository{
		httpClient: nil,
		namespace:  string(mockNamespace),
		token:      string(mockToken),
		caCert:     caCert,
	}
	type fields struct{}
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
		want     domain.IK8sRepository
		wantErr  bool
		mockArgs []mockArgs
	}{
		{
			name: "Should return a valid K8sRepository when files are read successfully",
			mockArgs: []mockArgs{
				{
					obj:     "u",
					method:  "ReadFile",
					args:    []interface{}{fmt.Sprintf("%s/namespace", serviceaccountPath)},
					returns: []interface{}{mockNamespace, nil},
				},
				{
					obj:     "u",
					method:  "ReadFile",
					args:    []interface{}{fmt.Sprintf("%s/token", serviceaccountPath)},
					returns: []interface{}{mockToken, nil},
				},
				{
					obj:     "u",
					method:  "ReadFile",
					args:    []interface{}{fmt.Sprintf("%s/ca.crt", serviceaccountPath)},
					returns: []interface{}{mockCacert, nil},
				},
			},
			want:    repo,
			wantErr: false,
		},
		{
			name: "Should return an error if reading the namespace file fails",
			mockArgs: []mockArgs{
				{
					obj:     "u",
					method:  "ReadFile",
					args:    []interface{}{fmt.Sprintf("%s/namespace", serviceaccountPath)},
					returns: []interface{}{[]byte{}, errors.New("namespace error")},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Should return an error if reading the token file fails",
			mockArgs: []mockArgs{
				{
					obj:     "u",
					method:  "ReadFile",
					args:    []interface{}{fmt.Sprintf("%s/namespace", serviceaccountPath)},
					returns: []interface{}{mockNamespace, nil},
				},
				{
					obj:     "u",
					method:  "ReadFile",
					args:    []interface{}{fmt.Sprintf("%s/token", serviceaccountPath)},
					returns: []interface{}{[]byte{}, errors.New("token error")},
				},
			},
			wantErr: true,
		},
		{
			name: "Should return an error if reading the ca.crt file fails",
			mockArgs: []mockArgs{
				{
					obj:     "u",
					method:  "ReadFile",
					args:    []interface{}{fmt.Sprintf("%s/namespace", serviceaccountPath)},
					returns: []interface{}{mockNamespace, nil},
				},
				{
					obj:     "u",
					method:  "ReadFile",
					args:    []interface{}{fmt.Sprintf("%s/token", serviceaccountPath)},
					returns: []interface{}{mockToken, nil},
				},
				{
					obj:     "u",
					method:  "ReadFile",
					args:    []interface{}{fmt.Sprintf("%s/ca.crt", serviceaccountPath)},
					returns: []interface{}{[]byte{}, errors.New("cacert error")},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			factory := &K8sRepositoryFactory{}
			u = &MockUtils{}
			for _, mockArgs := range tt.mockArgs {
				switch mockArgs.obj {
				case "u":
					u.(*MockUtils).On(mockArgs.method, mockArgs.args...).Return(mockArgs.returns...).Once()
				}
			}
			got, err := factory.Create()
			if (err != nil) != tt.wantErr {
				t.Errorf("K8sRepositoryFactory.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("K8sRepositoryFactory.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
