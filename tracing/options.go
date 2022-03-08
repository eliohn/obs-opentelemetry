// Copyright 2022 CloudWeGo Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tracing

import (
	"context"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	instrumentationName = "github.com/kitex-contrib/obs-opentelemetry"
)

// Option opts for opentelemetry tracer provider
type Option interface {
	apply(cfg *config)
}

type option func(cfg *config)

func (fn option) apply(cfg *config) {
	fn(cfg)
}

type config struct {
	tracer trace.Tracer
	meter  metric.Meter

	tracerProvider    trace.TracerProvider
	meterProvider     metric.MeterProvider
	textMapPropagator propagation.TextMapPropagator

	spanNameFormatter func(ctx context.Context) string

	withStackTrace        bool
	recordSourceOperation bool
}

func newConfig(opts []Option) *config {
	cfg := defaultConfig()

	for _, opt := range opts {
		opt.apply(cfg)
	}

	cfg.meter = cfg.meterProvider.Meter(
		instrumentationName,
		metric.WithInstrumentationVersion(SemVersion()),
	)

	return cfg
}

func defaultConfig() *config {
	return &config{
		tracerProvider:        otel.GetTracerProvider(),
		meterProvider:         global.GetMeterProvider(),
		textMapPropagator:     otel.GetTextMapPropagator(),
		recordSourceOperation: false,
		withStackTrace:        true,
		spanNameFormatter: func(ctx context.Context) string {
			endpoint := rpcinfo.GetRPCInfo(ctx).To()
			return endpoint.Method()
		},
		tracer: otel.GetTracerProvider().Tracer(
			instrumentationName,
			trace.WithInstrumentationVersion(SemVersion()),
		),
	}
}

// WithTracer configures tracer
func WithTracer(tracer trace.Tracer) Option {
	return option(func(cfg *config) {
		cfg.tracer = tracer
	})
}

// WithMeter configures meter
func WithMeter(meter metric.Meter) Option {
	return option(func(cfg *config) {
		cfg.meter = meter
	})
}

// WithSpanNameFormatter configures span name formatter
func WithSpanNameFormatter(fn func(ctx context.Context) string) Option {
	return option(func(cfg *config) {
		cfg.spanNameFormatter = fn
	})
}

// WithStackTrace configures stack trace
func WithStackTrace(stackTrace bool) Option {
	return option(func(cfg *config) {
		cfg.withStackTrace = stackTrace
	})
}

// WithRecordSourceOperation configures record source operation dimension
func WithRecordSourceOperation(recordSourceOperation bool) Option {
	return option(func(cfg *config) {
		cfg.recordSourceOperation = recordSourceOperation
	})
}

// WithTextMapPropagator configures propagation
func WithTextMapPropagator(p propagation.TextMapPropagator) Option {
	return option(func(cfg *config) {
		cfg.textMapPropagator = p
	})
}
