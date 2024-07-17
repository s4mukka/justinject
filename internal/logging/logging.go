package logging

import (
	"context"
	"os"

	"github.com/s4mukka/justinject/domain"
	"github.com/s4mukka/justinject/internal/otellogger"
	"github.com/s4mukka/justinject/internal/otellogrusdecorator"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/bridges/otellogrus"
)

var (
	newHook = otellogrusdecorator.NewDecoratedHook
)

func Init(ctx *context.Context) *log.Entry {
	environment := (*ctx).Value("environment").(*domain.Environment)
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

func AddOtelHook(ctx *context.Context) {
	environment := (*ctx).Value("environment").(*domain.Environment)

	otelHook := newHook(
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
