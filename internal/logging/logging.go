package logging

import (
	"context"

	"github.com/s4mukka/justinject/domain"
	"github.com/s4mukka/justinject/internal/otel"
	"github.com/s4mukka/justinject/internal/otellogger"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/bridges/otellogrus"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

func Init(ctx *context.Context) (*log.Entry, *sdklog.LoggerProvider) {
	environment := (*ctx).Value("environment").(*domain.Environment)
	logger := log.New()
	logger.SetFormatter(
		otellogger.OTELLogger{
			Formatter: log.JSONFormatter{},
		},
	)
	logger.SetLevel(log.DebugLevel)

	lp, err := otel.InitLogger(ctx)
	if err != nil {
		logger.Warnf("Error initializing logger: %v", err)
	} else {
		lp.Logger(environment.Instance)
	}

	otelHook := otellogrus.NewHook(
		environment.Instance,
		otellogrus.WithLoggerProvider(lp),
		otellogrus.WithLevels([]log.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
			logrus.InfoLevel,
		}),
	)

	logger.AddHook(otelHook)

	return logger.WithField("instance", environment.Instance), lp
}
