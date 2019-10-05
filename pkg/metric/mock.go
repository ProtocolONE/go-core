package metric

import "time"

// Counter is the emitter type metrics.
type counter struct {
}

// Inc increments the counter by a delta.
func (*counter) Inc(delta int64) {
}

// Gauge is the the emitter gauge metrics.
type gauge struct {
}

// Update sets the gauges absolute value.
func (*gauge) Update(value float64) {
}

// Timer is the emitter timer metrics.
type timer struct {
}

// Record a specific duration directly.
func (*timer) Record(value time.Duration) {
}

// Start gives you back a specific point in time to report via Stop.
func (*timer) Start() Stopwatch {
	return Stopwatch{}
}

// Histogram is the emitter histogram metrics
type histogram struct {
}

// RecordValue records a specific value directly.
// Will use the configured value buckets for the histogram.
func (*histogram) RecordValue(value float64) {
}

// RecordDuration records a specific duration directly.
// Will use the configured duration buckets for the histogram.
func (*histogram) RecordDuration(value time.Duration) {
}

// Start gives you a specific point in time to then record a duration.
// Will use the configured duration buckets for the histogram.
func (*histogram) Start() Stopwatch {
	return Stopwatch{}
}

// Capabilities is a description of metrics reporting capabilities.
type capabilities struct {
}

// Reporting returns whether the reporter has the ability to actively report.
func (*capabilities) Reporting() bool {
	return false
}

// Tagging returns whether the reporter has the capability for tagged metrics.
func (*capabilities) Tagging() bool {
	return false
}

// Mock is the client metric with stubbed methods
type Mock struct {
	tags map[string]string
	name string
}

// Counter returns the Counter object corresponding to the name.
func (t *Mock) Counter(name string) Counter {
	return &counter{}
}

// Gauge returns the Gauge object corresponding to the name.
func (t *Mock) Gauge(name string) Gauge {
	return &gauge{}
}

// Timer returns the Timer object corresponding to the name.
func (t *Mock) Timer(name string) Timer {
	return &timer{}
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
func (t *Mock) Histogram(name string, buckets Buckets) Histogram {
	return &histogram{}
}

// Tagged returns a new child scope with the given tags and current tags.
func (t *Mock) Tagged(tags map[string]string) Scope {
	return &Mock{tags: tags}
}

// SubScope returns a new child scope appending a further name prefix.
func (t *Mock) SubScope(name string) Scope {
	return &Mock{name: name}
}

// Capabilities returns a description of metrics reporting capabilities.
func (t *Mock) Capabilities() Capabilities {
	return &capabilities{}
}

// NewMock returns mock instance implemented of Scope interface
func NewMock() Scope {
	return &Mock{}
}
