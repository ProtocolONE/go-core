package metric

import (
	"context"
	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/cactus/go-statsd-client/statsd"
	"github.com/google/wire"
	tallystatsd "github.com/uber-go/tally/statsd"
)

// ProviderCfg returns configuration for production jaeger client
func ProviderCfg(cfg config.Configurator) (*Config, func(), error) {
	c := &Config{
		invoker: invoker.NewInvoker(),
	}
	e := cfg.UnmarshalKey(UnmarshalKey, c)
	return c, func() {}, e
}

// Provider returns client metric instance implemented of Scope interface with resolved dependencies
func Provider(ctx context.Context, log logger.Logger, cfg *Config) (Scope, func(), error) {
	if !cfg.Enabled {
		return ProviderTest()
	}
	statter, e := statsd.NewBufferedClient(
		cfg.StatsD.Addr,
		cfg.StatsD.Prefix,
		cfg.StatsD.FlushInterval,
		cfg.StatsD.FlushBytes,
	)
	if e != nil {
		return nil, nil, e
	}
	cfg.Scope.Reporter = tallystatsd.NewReporter(statter, cfg.StatsD.Options)
	m := NewTally(ctx, log, cfg.Scope, cfg.Interval)
	return m, func() {}, nil
}

// ProviderTest returns stub/mock client metric instance implemented of Scope interface with resolved dependencies
func ProviderTest() (Scope, func(), error) {
	m := NewMock()
	return m, func() {}, nil
}

var (
	WireSet     = wire.NewSet(Provider, ProviderCfg)
	WireTestSet = wire.NewSet(ProviderTest)
)
