package tracing

import (
	"github.com/opentracing/opentracing-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

const (
	Prefix       = "go-core.tracing"
	UnmarshalKey = "tracing"
)

type (
	Tracer = opentracing.Tracer
)

type Config struct {
	Enabled bool
	Jaeger  jaegerConfig.Configuration
}
