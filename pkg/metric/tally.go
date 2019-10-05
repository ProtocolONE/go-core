package metric

import (
	"context"
	"time"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/uber-go/tally"
)

// NewTally returns instance of uber/tally metric client implemented of Scope interface
func NewTally(ctx context.Context, log logger.Logger, options tally.ScopeOptions, interval time.Duration) Scope {
	log = log.WithFields(logger.Fields{"service": Prefix})
	scope, closer := tally.NewRootScope(options, interval)
	go func() {
		<-ctx.Done()
		if e := closer.Close(); e != nil {
			log.Error("%v", logger.Args(e))
		}
	}()
	return scope
}
