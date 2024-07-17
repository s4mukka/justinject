package otellogrusdecorator

import (
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/bridges/otellogrus"
)

type Hook interface {
	Levels() []logrus.Level
	Fire(entry *logrus.Entry) error
}

type DecoratedHook struct {
	hook Hook
}

func NewDecoratedHook(name string, options ...otellogrus.Option) Hook {
	return &DecoratedHook{
		hook: otellogrus.NewHook(name, options...),
	}
}

func (h *DecoratedHook) Levels() []logrus.Level {
	return h.hook.Levels()
}

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
