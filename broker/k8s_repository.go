package broker

import (
	"crypto/x509"
	"fmt"

	"github.com/s4mukka/justinject/domain"
	"github.com/s4mukka/justinject/internal/utils"
)

type K8sRepository struct {
	httpClient domain.IHTTPClient

	namespace string
	token     string
	caCert    *x509.CertPool
}

var (
	u domain.IUtils = &utils.Utils{}
)

const (
	apiserver          string = "https://kubernetes.default.svc"
	serviceaccountPath string = "/var/run/secrets/kubernetes.io/serviceaccount"
	tmplPath           string = "/k8s/job.tmpl"
)

func (r *K8sRepository) CreateJob(job domain.IJob) error {
	body, err := u.TemplateToString(tmplPath, job.ParseTemplate())
	if err != nil {
		return err
	}
	_, err = r.httpClient.HTTPRequest(
		fmt.Sprintf("%s/apis/batch/v1/namespaces/%s/jobs", apiserver, r.namespace),
		"POST",
		map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", r.token),
			"Content-Type":  "application/yaml",
		},
		[]byte(body),
		r.caCert,
	)
	return err
}

type K8sRepositoryFactory struct {
}

func (f K8sRepositoryFactory) Create() (domain.IK8sRepository, error) {
	namespace, err := u.ReadFile(fmt.Sprintf("%s/namespace", serviceaccountPath))
	if err != nil {
		return nil, err
	}
	token, err := u.ReadFile(fmt.Sprintf("%s/token", serviceaccountPath))
	if err != nil {
		return nil, err
	}
	caCertFile, err := u.ReadFile(fmt.Sprintf("%s/ca.crt", serviceaccountPath))
	if err != nil {
		return nil, err
	}
	caCert := x509.NewCertPool()
	caCert.AppendCertsFromPEM(caCertFile)
	return &K8sRepository{
		httpClient: nil,
		namespace:  string(namespace),
		token:      string(token),
		caCert:     caCert,
	}, nil
}
