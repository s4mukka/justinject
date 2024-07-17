package mock

import (
	"context"

	"github.com/s4mukka/justinject/domain"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

type MockedLoggerProvider struct {
	mock.Mock
}

func (m *MockedLoggerProvider) Get() domain.IOtelLoggerProvider {
	return nil
}

func (m *MockedLoggerProvider) Logger(name string, opts ...domain.OtelLoggerOption) domain.OtelLogger {
	args := m.Called(name)
	return args.Get(0).(domain.OtelLogger)
}

func (m *MockedLoggerProvider) Shutdown(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

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
