package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/s4mukka/justinject/domain"
)

type MockLogger struct {
	mock.Mock
	domain.ILogger
}

func (m *MockLogger) Infof(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Errorf(format string, args ...interface{}) {
	m.Called(format, args)
}
