package tracing

import (
	"github.com/opentracing/opentracing-go"
)

const (
	Prefix       = "go-core.tracing"
	UnmarshalKey = "tracing"
)

type (
	Tracer = opentracing.Tracer
)
