package provider

import (
	"github.com/ProtocolONE/go-core/config"
	"github.com/ProtocolONE/go-core/logger"
	"github.com/ProtocolONE/go-core/metric"
	"github.com/ProtocolONE/go-core/tracing"
	"github.com/google/wire"
)

var Set = wire.NewSet(
	logger.WireSet,
	config.WireSet,
	metric.WireSet,
	tracing.WireSet,
)
