package logging

import (
	"os"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/bridges/otellogrus"

	"github.com/s4mukka/justinject/domain"
	"github.com/s4mukka/justinject/internal/loggrushook"
	"github.com/s4mukka/justinject/internal/otellogger"
)

var newOtelLoggrusHook = loggrushook.NewOtelLoggrusHook

func Init(ctx domain.IContext) domain.ILogger {
	environment := (ctx).Value(domain.EnvironmentKey).(*domain.Environment)
	logger := log.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(
		otellogger.OTELLogger{
			Formatter: log.JSONFormatter{
				FieldMap: log.FieldMap{
					log.FieldKeyLevel: "severity_text",
				},
			},
		},
	)
	logger.WriterLevel(log.DebugLevel)

	return logger.WithField("instance", environment.Instance)
}

func AddOtelHook(ctx domain.IContext) {
	environment := (ctx).Value(domain.EnvironmentKey).(*domain.Environment)

	otelHook := newOtelLoggrusHook(
		environment.Instance,
		otellogrus.WithLoggerProvider(environment.LoggerProvider.Get()),
		otellogrus.WithLevels([]log.Level{
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
			log.WarnLevel,
			log.InfoLevel,
		}),
	)

	environment.Logger.(*log.Entry).Logger.AddHook(otelHook)
}
