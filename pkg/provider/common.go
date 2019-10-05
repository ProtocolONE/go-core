package provider

import (
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/metric"
	"github.com/ProtocolONE/go-core/v2/pkg/tracing"
)

type AwareSet struct {
	Logger logger.Logger
	Metric metric.Scope
	Tracer tracing.Tracer
}

type LMT interface {
	L() logger.Logger
	M() metric.Scope
	T() tracing.Tracer
}

// L returns logger instance implemented of Logger interface with resolved dependencies
func (a *AwareSet) L() logger.Logger {
	return a.Logger
}

// M returns client metric instance implemented of Scope interface with resolved dependencies
func (a *AwareSet) M() metric.Scope {
	return a.Metric
}

// T returns instance implemented of Tracer interface with resolved dependencies
func (a *AwareSet) T() tracing.Tracer {
	return a.Tracer
}
