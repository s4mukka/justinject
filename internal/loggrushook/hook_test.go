package loggrushook

import (
	"testing"

	"github.com/s4mukka/justinject/internal/loggrushook/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewDecoratedHook(t *testing.T) {
	name := "test-hook"
	hook := NewOtelLoggrusHook(name)
	assert.NotNil(t, hook)
	assert.IsType(t, &OtelLoggrusHook{}, hook)
}

func TestDecoratedHook_Levels(t *testing.T) {
	mockHook := new(mocks.MockedHook)
	mockHook.On("Levels").Return([]logrus.Level{logrus.InfoLevel, logrus.WarnLevel})
	hook := &OtelLoggrusHook{
		hook: mockHook,
	}

	levels := hook.Levels()
	assert.Equal(t, []logrus.Level{logrus.InfoLevel, logrus.WarnLevel}, levels)

	mockHook.AssertCalled(t, "Levels")
}

func TestDecoratedHook_Fire(t *testing.T) {
	mockHook := new(mocks.MockedHook)
	entry := &logrus.Entry{
		Level: logrus.WarnLevel,
		Data:  logrus.Fields{},
	}
	mockHook.On("Fire", entry).Return(nil)

	hook := &OtelLoggrusHook{
		hook: mockHook,
	}

	err := hook.Fire(entry)
	assert.NoError(t, err)
	assert.Equal(t, "warn", entry.Data["level"])

	mockHook.AssertCalled(t, "Fire", entry)
}

func TestConvertLogLevel(t *testing.T) {
	tests := []struct {
		input    logrus.Level
		expected string
	}{
		{logrus.DebugLevel, "debug"},
		{logrus.InfoLevel, "info"},
		{logrus.WarnLevel, "warn"},
		{logrus.ErrorLevel, "error"},
		{logrus.FatalLevel, "fatal"},
		{logrus.PanicLevel, "panic"},
		{logrus.Level(123), "unknown"},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, convertLogLevel(test.input))
	}
}
