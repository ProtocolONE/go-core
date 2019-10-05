package provider

import (
	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/metric"
	"github.com/ProtocolONE/go-core/v2/pkg/tracing"
	"github.com/google/wire"
)

var Set = wire.NewSet(
	logger.WireSet,
	config.WireSet,
	metric.WireSet,
	tracing.WireSet,
)
