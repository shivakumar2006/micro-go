package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

func Init(ctx context.Context) (*sdktrace.TracerProvider, error) {
	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpoint("otel-collector-service:4318"), otlptracehttp.WithInsecure())
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx, resource.WithAttributes(semconv.ServiceNameKey.String("notification-service")))
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter), sdktrace.WithResource(res))

	otel.SetTracerProvider(tp)

	return tp, nil
}
