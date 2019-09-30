package tracing

import (
	"context"
	"github.com/ProtocolONE/go-core/config"
	"github.com/ProtocolONE/go-core/logger"
	"github.com/google/wire"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

// ProviderCfg returns configuration for production jaeger client
func ProviderCfg(cfg config.Configurator) (*jaegerConfig.Configuration, func(), error) {
	c := &jaegerConfig.Configuration{
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans: false,
		},
	}
	e := cfg.UnmarshalKey(UnmarshalKey, c)
	return c, func() {}, e
}

// Provider returns instance implemented of opentracing.Tracer interface with resolved dependencies
func Provider(ctx context.Context, cfg *jaegerConfig.Configuration, log logger.Logger) (Tracer, func(), error) {
	if cfg.Disabled {
		cfg.ServiceName = "disabled"
	}
	t, e := New(ctx, log, cfg, jaegerConfig.Logger(jaeger.StdLogger))
	if !cfg.Disabled {
		opentracing.SetGlobalTracer(t)
	}
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
