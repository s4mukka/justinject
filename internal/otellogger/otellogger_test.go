package otellogger

import (
	"context"
	"encoding/json"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func TestOTELLogger_Format(t *testing.T) {
	traceID := trace.TraceID{0x1}
	spanID := trace.SpanID{0x1}
	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
	})

	tp := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(tp)
	tracer := otel.Tracer("test-tracer")

	ctx := trace.ContextWithSpanContext(context.Background(), spanContext)
	_, span := tracer.Start(ctx, "test-span")
	defer span.End()

	entry := &log.Entry{
		Context: ctx,
		Data:    log.Fields{},
		Message: "test message",
		Level:   log.InfoLevel,
	}

	logger := OTELLogger{
		Formatter: log.JSONFormatter{
			FieldMap: log.FieldMap{
				log.FieldKeyTime:  "time",
				log.FieldKeyLevel: "level",
				log.FieldKeyMsg:   "msg",
			},
		},
	}

	formatted, err := logger.Format(entry)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(formatted, &result)
	assert.NoError(t, err)

	assert.Equal(t, traceID.String(), result["trace_id"])
	assert.Equal(t, spanID.String(), result["span_id"])
	assert.Equal(t, spanContext.TraceFlags().String(), result["Context"].(map[string]interface{})["TraceFlags"])
	assert.Equal(t, spanContext.TraceState().String(), result["Context"].(map[string]interface{})["TraceState"])
	assert.Equal(t, "test message", result["msg"])
	assert.Equal(t, "info", result["level"])
}
