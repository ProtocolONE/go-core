package entrypoint

import (
	"context"
	"github.com/ProtocolONE/go-core/config"
	"github.com/ProtocolONE/go-core/logger"
	"github.com/ProtocolONE/go-core/metric"
	"github.com/ProtocolONE/go-core/tracing"
)

type WithKeyInitial string

const (
	Prefix        = "go-core.entrypoint"
	CtxKeyInitial = WithKeyInitial(Prefix)
)

type Master interface {
	Slaver
	// Shutdown raise event of shutdown for all subscribers
	Shutdown(ctx context.Context, code int)
	// Reload raise event of reload for all subscribers
	Reload()
	// Serve execute builder and runner functions with callback for pre run
	Serve(preRun func() error) error
}

type Slaver interface {
	// WorkDir returns current work directory
	WorkDir() string
	// OnShutdown subscribe on shutdown event for gracefully exit via context.
	OnShutdown() context.Context
	// OnReload subscribe on reload event.
	OnReload(callback func(ctx context.Context))
	// Initial returns initial settings
	Initial() config.Initial
	// Logger returns logger instance implemented of Logger interface
	Logger() logger.Logger
	// Metric returns client metric instance implemented of Scope interface
	Metric() metric.Scope
	// Tracer returns instance implemented of opentracing.Tracer interface
	Tracer() tracing.Tracer
	// Executor provide interface for set builder and runner callback functions
	Executor(builder func(ctx context.Context) error, runner func(ctx context.Context) error)
}
