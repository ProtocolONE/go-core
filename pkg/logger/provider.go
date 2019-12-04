package logger

import (
	"context"
	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/google/wire"
	"sync"
)

var (
	z  *Zap
	mu sync.Mutex
)

// ProviderCfg returns configuration for production logger
func ProviderCfg(cfg config.Configurator) (*Config, func(), error) {
	c := &Config{
		invoker: invoker.NewInvoker(),
	}
	e := cfg.UnmarshalKeyOnReload(UnmarshalKey, c, StringToLoggerLevelHookFunc())
	return c, func() {}, e
}

// Provider returns logger instance implemented of Logger interface with resolved dependencies
func Provider(ctx context.Context, cfg *Config) (*Zap, func(), error) {
	defer mu.Unlock()
	mu.Lock()
	if z != nil {
		return z, func() {}, nil
	}
	z = NewZap(ctx, cfg)
	return z, func() {}, nil
}

// ProviderTest returns stub/mock logger instance implemented of Logger interface with resolved dependencies
func ProviderTest(ctx context.Context, cfg *Config) (*Mock, func(), error) {
	return NewMock(ctx, cfg, true), func() {}, nil
}

var (
	WireSet     = wire.NewSet(Provider, ProviderCfg, wire.Bind(new(Logger), new(*Zap)))
	WireTestSet = wire.NewSet(ProviderTest, ProviderCfg, wire.Bind(new(Logger), new(*Mock)))
)
