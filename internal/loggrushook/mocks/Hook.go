package mocks

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

type MockedHook struct {
	mock.Mock
}

func (m *MockedHook) Levels() []logrus.Level {
	args := m.Called()
	return args.Get(0).([]logrus.Level)
}

func (h *MockedHook) Fire(entry *logrus.Entry) error {
	args := h.Called(entry)
	return args.Error(0)
}
