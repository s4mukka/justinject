package domain

import "crypto/x509"

type HTTPResponse struct {
	Data       []byte
	StatusCode int
}

type IHTTPClient interface {
	HTTPRequest(url string, method string, header map[string]string, body []byte, cacerts ...*x509.CertPool) (HTTPResponse, error)
}

type IRestContext interface {
	JSON(code int, obj any)
	ShouldBindJSON(obj any) error
}
