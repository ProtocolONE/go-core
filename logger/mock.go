package logger

import (
	"context"
)

// Entity represent log struct with all assets
type Entity struct {
	Level  Level
	Fields Fields
	Args   []interface{}
	Tags   Tags
	Format string
}

// Mock is the logger with stubbed methods
type Mock struct {
	ctx     context.Context
	ch      chan Entity
	discard bool
	cfg     *Config
	tags    []string
	fields  map[string]interface{}
}

// Printf is like fmt.Printf, push to log entry with debug level
func (m *Mock) Printf(format string, a ...interface{}) {
	m.Debug(format, Args(a...))
}

// Verbose should return true when verbose logging output is wanted
func (m *Mock) Verbose() bool {
	return m.cfg.Verbose
}

// Emergency push to log entry with emergency level & throw panic
func (m *Mock) Emergency(format string, opts ...Option) {
	m.Log(LevelEmergency, format, opts...)
}

// Alert push to log entry with alert level
func (m *Mock) Alert(format string, opts ...Option) {
	m.Log(LevelAlert, format, opts...)
}

// Critical push to log entry with critical level
func (m *Mock) Critical(format string, opts ...Option) {
	m.Log(LevelCritical, format, opts...)
}

// Error push to log entry with error level
func (m *Mock) Error(format string, opts ...Option) {
	m.Log(LevelError, format, opts...)
}

// Warning push to log entry with warning level
func (m *Mock) Warning(format string, opts ...Option) {
	m.Log(LevelWarning, format, opts...)
}

// Notice push to log entry with notice level
func (m *Mock) Notice(format string, opts ...Option) {
	m.Log(LevelNotice, format, opts...)
}

// Info push to log entry with info level
func (m *Mock) Info(format string, opts ...Option) {
	m.Log(LevelInfo, format, opts...)
}

// Debug push to log entry with debug level
func (m *Mock) Debug(format string, opts ...Option) {
	m.Log(LevelDebug, format, opts...)
}

// Write push to log entry with debug level
func (m *Mock) Write(p []byte) (n int, err error) {
	m.Debug(string(p))
	return len(p), nil
}

// Log push to log with specified level
func (m *Mock) Log(level Level, format string, o ...Option) {
	if !m.discard {
		return
	}
	opts := &opts{}
	for _, option := range o {
		_ = option(opts)
	}
	var (
		fields = map[string]interface{}{}
		tags   []string
	)
	// tags
	if len(m.tags) > 0 {
		tags = make([]string, len(m.tags))
		copy(tags, m.tags)
	}
	if len(opts.tags) > 0 {
		tags = append(tags, opts.tags...)
	}
	// fields
	for k, v := range m.fields {
		fields[k] = v
	}
	for k, v := range opts.fields {
		fields[k] = v
	}
	m.ch <- Entity{
		Level:  level,
		Format: format,
		Args:   opts.args,
		Fields: fields,
		Tags:   tags,
	}
}

// WithFields create new instance with fields
func (m *Mock) WithFields(fields Fields) Logger {
	nm := &Mock{}
	copyMock(nm, m, nil, fields)
	return nm
}

// WithTags create new instance with tags
func (m *Mock) WithTags(tags Tags) Logger {
	nm := &Mock{}
	copyMock(nm, m, tags, nil)
	return nm
}

func copyMock(dst, src *Mock, tags []string, fields map[string]interface{}) {
	// tags
	cTags := make([]string, len(src.tags))
	copy(cTags, src.tags)
	dst.tags = cTags
	if tags != nil {
		dst.tags = append(dst.tags, tags...)
	}
	var cFields = map[string]interface{}{}
	// fields
	for k, v := range src.fields {
		cFields[k] = v
	}
	dst.fields = cFields
	for k, v := range fields {
		dst.fields[k] = v
	}
}

// Catch returns channel of entity structure for testing event content
func (m *Mock) Catch() <-chan Entity {
	return m.ch
}

// NewMock returns mock instance implemented of Logger interface
func NewMock(ctx context.Context, cfg *Config, discard bool) *Mock {
	copyCfg := *cfg
	return &Mock{
		ctx:     ctx,
		cfg:     &copyCfg,
		ch:      make(chan Entity),
		discard: discard,
	}
}
