package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/s4mukka/justinject/domain"
)

type MockLogger struct {
	mock.Mock
	domain.ILogger
}

func (m *MockLogger) Warnf(format string, args ...interface{}) {
	m.Called(format, args)
}
