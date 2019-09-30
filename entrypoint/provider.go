package entrypoint

import (
	"github.com/ProtocolONE/go-core/config"
	"github.com/ProtocolONE/go-core/logger"
	"github.com/ProtocolONE/go-core/metric"
	"github.com/ProtocolONE/go-core/tracing"
	"github.com/google/wire"
)

type AppSet struct {
	Logger logger.Logger
	Metric metric.Scope
	Tracer tracing.Tracer
}

// Provider returns entrypoint instance implemented of Master interface with resolved dependencies
func Provider(set AppSet, initial config.Initial) (Master, func(), error) {
	e, r := NewEntryPoint(set, initial)
	return e, func() {}, r
}

var (
	WireSet = wire.NewSet(Provider, wire.Struct(new(AppSet), "*"))
)
