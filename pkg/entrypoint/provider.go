package entrypoint

import (
	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/metric"
	"github.com/ProtocolONE/go-core/v2/pkg/tracing"
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
