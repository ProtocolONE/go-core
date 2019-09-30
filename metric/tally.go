package metric

import (
	"context"
	"github.com/ProtocolONE/go-core/logger"
	"github.com/uber-go/tally"
	"time"
)

// Tally is uber/tally metric client
type Tally struct {
	client Scope
}

// Counter returns the Counter object corresponding to the name.
func (t *Tally) Counter(name string) Counter {
	return t.client.Counter(name)
}

// Gauge returns the Gauge object corresponding to the name.
func (t *Tally) Gauge(name string) Gauge {
	return t.client.Gauge(name)
}

// Timer returns the Timer object corresponding to the name.
func (t *Tally) Timer(name string) Timer {
	return t.client.Timer(name)
}

// Histogram returns the Histogram object corresponding to the name.
// To use default value and duration buckets configured for the scope
// simply pass tally.DefaultBuckets or nil.
// You can use tally.ValueBuckets{x, y, ...} for value buckets.
// You can use tally.DurationBuckets{x, y, ...} for duration buckets.
// You can use tally.MustMakeLinearValueBuckets(start, width, count) for linear values.
// You can use tally.MustMakeLinearDurationBuckets(start, width, count) for linear durations.
// You can use tally.MustMakeExponentialValueBuckets(start, factor, count) for exponential values.
// You can use tally.MustMakeExponentialDurationBuckets(start, factor, count) for exponential durations.
func (t *Tally) Histogram(name string, buckets Buckets) Histogram {
	return t.client.Histogram(name, buckets)
}

// Tagged returns a new child scope with the given tags and current tags.
func (t *Tally) Tagged(tags map[string]string) Scope {
	return t.client.Tagged(tags)
}

// SubScope returns a new child scope appending a further name prefix.
func (t *Tally) SubScope(name string) Scope {
	return t.client.SubScope(name)
}

// Capabilities returns a description of metrics reporting capabilities.
func (t *Tally) Capabilities() Capabilities {
	return t.client.Capabilities()
}

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
	return &Tally{client: scope}
}
