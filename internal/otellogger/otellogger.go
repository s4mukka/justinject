package otellogger

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type OTELLogger struct {
	Formatter log.JSONFormatter
}

func (l OTELLogger) Format(entry *log.Entry) ([]byte, error) {
	entry.Data["AA"] = "AA"
	span := trace.SpanFromContext(entry.Context)
	entry.Data["trace_id"] = span.SpanContext().TraceID().String()
	entry.Data["span_id"] = span.SpanContext().SpanID().String()
	entry.Data["Context"] = span.SpanContext()
	fmt.Printf("ENTRY OTELLOGER %p\n", &entry)
	return l.Formatter.Format(entry)
}
