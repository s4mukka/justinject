package logging

import (
	"context"

	"github.com/s4mukka/justinject/domain"
	"github.com/s4mukka/justinject/internal/otellogger"
	"github.com/s4mukka/justinject/internal/otellogrus"
	log "github.com/sirupsen/logrus"
)

func Init(ctx *context.Context) *log.Entry {
	environment := (*ctx).Value("environment").(*domain.Environment)
	logger := log.New()
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

	return logger.WithField("instance", environment.Instance).WithField("teste", "testee")
}

func AddOtelHook(ctx *context.Context) {
	environment := (*ctx).Value("environment").(*domain.Environment)

	otelHook := otellogrus.NewHook(
		environment.Instance,
		otellogrus.WithLoggerProvider(environment.LoggerProvider),
		otellogrus.WithLevels([]log.Level{
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
			log.WarnLevel,
			log.InfoLevel,
		}),
	)

	environment.Logger.Logger.AddHook(otelHook)
}
