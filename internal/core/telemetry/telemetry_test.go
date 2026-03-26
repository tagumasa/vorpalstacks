package telemetry

import (
	"context"
	"testing"
)

func TestInit_NoExporter(t *testing.T) {
	t.Setenv("OTEL_TRACES_EXPORTER", "none")

	shutdown, err := Init("test-service", "1.0.0")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer shutdown(context.Background())

	tracer := Tracer()
	if tracer == nil {
		t.Error("tracer should not be nil")
	}
}

func TestInit_OffExporter(t *testing.T) {
	t.Setenv("OTEL_TRACES_EXPORTER", "off")

	shutdown, err := Init("test-service", "1.0.0")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer shutdown(context.Background())

	tracer := Tracer()
	if tracer == nil {
		t.Error("tracer should not be nil")
	}
}

func TestInit_EmptyExporter(t *testing.T) {
	t.Setenv("OTEL_TRACES_EXPORTER", "")

	shutdown, err := Init("test-service", "1.0.0")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer shutdown(context.Background())

	tracer := Tracer()
	if tracer == nil {
		t.Error("tracer should not be nil")
	}
}

func TestTracer_Default(t *testing.T) {
	tracer := Tracer()
	if tracer == nil {
		t.Error("tracer should not be nil even without init")
	}
}

func TestStartSpan(t *testing.T) {
	ctx := context.Background()
	ctx, span := StartSpan(ctx, "test-span")
	defer span.End()

	if span == nil {
		t.Error("span should not be nil")
	}
}

func TestSpanFromContext(t *testing.T) {
	ctx := context.Background()
	span := SpanFromContext(ctx)
	if span == nil {
		t.Error("span should not be nil")
	}
}

func TestContextWithSpan(t *testing.T) {
	ctx := context.Background()
	_, span := StartSpan(ctx, "test-span")
	ctxWithSpan := ContextWithSpan(ctx, span)

	if ctxWithSpan == nil {
		t.Error("context with span should not be nil")
	}
}
