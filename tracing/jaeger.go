package tracing

import (
	"context"
	"github.com/ProtocolONE/go-core/logger"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-client-go/config"
)

// New returns instance implemented of opentracing.Tracer interface
func New(ctx context.Context, log logger.Logger, cfg *config.Configuration, option ...config.Option) (Tracer, error) {
	log = log.WithFields(logger.Fields{"service": Prefix})
	tracer, closer, e := cfg.NewTracer(option...)
	if e != nil {
		return tracer, errors.WithMessage(e, Prefix)
	}
	go func() {
		<-ctx.Done()
		if e := closer.Close(); e != nil {
			log.Error("%v", logger.Args(e))
		}
	}()
	return tracer, errors.WithMessage(e, Prefix)
}
