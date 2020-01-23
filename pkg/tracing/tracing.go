package tracing

import (
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	reporterhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

const EndpointURL = "http://localhost:9411/api/v2/spans"

func InitTracing(cfg TracingConfig) {
	tracer, err := NewTracer(cfg)
	if err != nil {
		return
	}
	opentracing.SetGlobalTracer(tracer)
}

func NewTracer(cfg TracingConfig) (opentracing.Tracer, error) {
	// The reporter sends traces to zipkin server
	reporter := reporterhttp.NewReporter(cfg.ReporterUrl)

	// Local endpoint represent the local service information
	localEndpoint := &model.Endpoint{ServiceName: cfg.ReporterUrl, Port: cfg.ReporterPort}

	// Sampler tells you which traces are going to be sampled or not. In this case we will record 100% (1.00) of traces.
	sampler, err := zipkin.NewCountingSampler(1)
	if err != nil {
		return nil, err
	}

	t, err := zipkin.NewTracer(
		reporter,
		zipkin.WithSampler(sampler),
		zipkin.WithLocalEndpoint(localEndpoint),
	)
	if err != nil {
		return nil, err
	}

	return t, err
}
