package tracing

import (
	"context"
	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/google/wire"
	"github.com/opentracing/opentracing-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	"sync"
)

var (
	t  Tracer
	mu sync.Mutex
)

// ProviderCfg returns configuration for production jaeger client
func ProviderCfg(cfg config.Configurator) (*Config, func(), error) {
	jaeger := jaegerConfig.Configuration{
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans: false,
		},
	}
	c := &Config{
		Jaeger: jaeger,
	}
	e := cfg.UnmarshalKey(UnmarshalKey, c)
	return c, func() {}, e
}

// Provider returns instance implemented of opentracing.Tracer interface with resolved dependencies
func Provider(ctx context.Context, cfg *Config, log logger.Logger) (Tracer, func(), error) {
	defer mu.Unlock()
	mu.Lock()
	if t != nil {
		return t, func() {}, nil
	}
	if !cfg.Enabled {
		return ProviderTest()
	}
	var e error
	t, e = New(ctx, log, cfg, jaegerConfig.Logger(NewLoggerAdapter(log)))
	opentracing.SetGlobalTracer(t)
	return t, func() {}, e
}

// ProviderTest returns stub/mock instance implemented of opentracing.Tracer interface with resolved dependencies
func ProviderTest() (Tracer, func(), error) {
	m := NewMock()
	return m, func() {}, nil
}

var (
	WireSet     = wire.NewSet(Provider, ProviderCfg)
	WireTestSet = wire.NewSet(ProviderTest)
)
