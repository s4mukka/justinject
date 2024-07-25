package loggrushook

import (
	"github.com/s4mukka/justinject/domain"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/bridges/otellogrus"
)

type OtelLoggrusHook struct {
	hook domain.IHook
}

func NewOtelLoggrusHook(name string, options ...otellogrus.Option) domain.IHook {
	return &OtelLoggrusHook{
		hook: otellogrus.NewHook(name, options...),
	}
}

func (h *OtelLoggrusHook) Levels() []logrus.Level {
	return h.hook.Levels()
}

func (h *OtelLoggrusHook) Fire(entry *logrus.Entry) error {
	entry.Data["level"] = convertLogLevel(entry.Level)
	return h.hook.Fire(entry)
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
