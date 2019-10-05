package tracing

import (
	"github.com/opentracing/opentracing-go/mocktracer"
)

// NewMock returns mock instance implemented of opentracing.Tracer interface
func NewMock() Tracer {
	return mocktracer.New()
}
