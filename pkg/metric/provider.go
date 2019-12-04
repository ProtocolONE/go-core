package metric

import (
	"context"
	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/cactus/go-statsd-client/statsd"
	"github.com/google/wire"
	promreporter "github.com/uber-go/tally/prometheus"
	tallystatsd "github.com/uber-go/tally/statsd"
	"net/http"
	"net/url"
	"sync"
)

var (
	m  Scope
	mu sync.Mutex
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

// ProviderPrometheus returns prometheus connector metric instance implemented of Scope interface with resolved dependencies
func ProviderPrometheus(ctx context.Context, log logger.Logger, cfg *Config) (Scope, func(), error) {
	defer mu.Unlock()
	mu.Lock()
	if m != nil {
		return m, func() {}, nil
	}
	if !cfg.Enabled {
		return ProviderTest()
	}
	u, e := url.Parse(cfg.Prometheus.Address)
	if e != nil {
		return nil, nil, e
	}
	cfgCopy := cfg
	r := promreporter.NewReporter(cfgCopy.Prometheus.Options)
	cfgCopy.Scope.Tags = map[string]string{}
	cfgCopy.Scope.CachedReporter = r
	if cfgCopy.Scope.Separator == "" {
		cfgCopy.Scope.Separator = promreporter.DefaultSeparator
	}
	//
	http.Handle(u.Path, r.HTTPHandler())
	go func() {
		err := http.ListenAndServe(u.Host, nil)
		if err != nil {
			panic(err)
		}
	}()
	m = NewTally(ctx, log, cfgCopy.Scope, cfgCopy.Interval)
	return m, func() {}, nil
}

// ProviderTest returns stub/mock client metric instance implemented of Scope interface with resolved dependencies
func ProviderTest() (Scope, func(), error) {
	m := NewMock()
	return m, func() {}, nil
}

var (
	WireSet     = wire.NewSet(ProviderPrometheus, ProviderCfg)
	WireTestSet = wire.NewSet(ProviderTest)
)
