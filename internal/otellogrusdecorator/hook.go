package otellogrusdecorator

import (
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/bridges/otellogrus"
)

type DecoratedHook struct {
	hook *otellogrus.Hook
}

func NewDecoratedHook(name string, options ...otellogrus.Option) *DecoratedHook {
	return &DecoratedHook{
		hook: otellogrus.NewHook(name, options...),
	}
}

// Levels returns the list of log levels we want to be sent to OpenTelemetry.
func (h *DecoratedHook) Levels() []logrus.Level {
	return h.hook.Levels()
}

// Fire handles the passed record, and sends it to OpenTelemetry.
func (h *DecoratedHook) Fire(entry *logrus.Entry) error {
	entry.Data["level"] = convertLogLevel(entry.Level)
	h.hook.Fire(entry)
	return nil
}

func convertLogLevel(level logrus.Level) string {
	switch level {
	case logrus.DebugLevel:
		return "debug"
	case logrus.InfoLevel:
		return "info"
	case logrus.WarnLevel:
		return "warn"
	case logrus.ErrorLevel:
		return "error"
	case logrus.FatalLevel:
		return "fatal"
	case logrus.PanicLevel:
		return "panic"
	default:
		return "unknown"
	}
}
