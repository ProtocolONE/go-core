// +build wireinject

package entrypoint

import (
	"context"
	"github.com/ProtocolONE/go-core/config"
	"github.com/ProtocolONE/go-core/invoker"
	"github.com/ProtocolONE/go-core/logger"
	"github.com/ProtocolONE/go-core/metric"
	"github.com/ProtocolONE/go-core/tracing"
	"github.com/google/wire"
)

// Build returns entrypoint instance implemented of Master interface with resolved dependencies
func Build(ctx context.Context, initial config.Initial, observer invoker.Observer) (Master, func(), error) {
	panic(wire.Build(WireSet, logger.WireSet, metric.WireSet, tracing.WireSet, config.WireSet))
}
