// +build wireinject

package entrypoint

import (
	"context"
	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/metric"
	"github.com/ProtocolONE/go-core/v2/pkg/tracing"
	"github.com/google/wire"
)

// Build returns entrypoint instance implemented of Master interface with resolved dependencies
func Build(ctx context.Context, initial config.Initial, observer invoker.Observer) (Master, func(), error) {
	panic(wire.Build(WireSet, logger.WireSet, metric.WireSet, tracing.WireSet, config.WireSet))
}
