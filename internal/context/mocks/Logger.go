package mocks

import (
	"github.com/s4mukka/justinject/domain"
	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
	domain.ILogger
}

func (m *MockLogger) Warnf(format string, args ...interface{}) {
	m.Called(format, args)
}
