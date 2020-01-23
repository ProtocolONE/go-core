package entrypoint

import (
	"context"
	"os"
	"sync"

	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/metric"
	"github.com/ProtocolONE/go-core/v2/pkg/tracing"
	"github.com/pkg/errors"
)

var (
	ErrViperNotInitialized = errors.New("viper not initialized")
	ErrExecutorNotPresent  = errors.New("executor is not set, nothing to execute")
)

// NewEntryPoint returns instance of entry point manager
func NewEntryPoint(set AppSet, initial config.Initial) (Master, error) {
	shutdownCtx, cancelFn := context.WithCancel(context.Background())
	ep := &EntryPoint{
		initial:     &initial,
		shutdownCtx: context.WithValue(shutdownCtx, CtxKeyInitial, &initial),
		cancelFn:    cancelFn,
		invoker:     invoker.NewInvoker(),
	}
	if initial.Viper == nil {
		panic(errors.WithMessage(ErrViperNotInitialized, Prefix))
	}
	if len(initial.WorkDir) > 0 {
		ep.initial.WorkDir = initial.WorkDir
	} else {
		ep.initial.WorkDir, _ = os.Getwd()
	}
	ep.initial.Viper = initial.Viper
	ep.set = set
	return ep, nil
}

// Slave returns instance of entry point for slaves
func (e *EntryPoint) Slave() (Slaver, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e, nil
}

// OnShutdown subscribe on shutdown event for gracefully exit via context.
func (e *EntryPoint) OnShutdown() context.Context {
	return e.shutdownCtx
}

// OnReload subscribe on reload event.
func (e *EntryPoint) OnReload(callback func(ctx context.Context)) {
	e.invoker.OnReload(callback)
}

// EntryPoint manager of single point of application
type EntryPoint struct {
	initial     *config.Initial
	set         AppSet
	builder     func(ctx context.Context) error
	runner      func(ctx context.Context) error
	mu          sync.Mutex
	shutdownCtx context.Context
	cancelFn    context.CancelFunc
	invoker     *invoker.Invoker
}

// Initial returns initial settings
func (e *EntryPoint) Initial() config.Initial {
	return *e.initial
}

// Logger returns logger instance implemented of Logger interface
func (e *EntryPoint) Logger() logger.Logger {
	return e.set.Logger
}

// Metric returns client metric instance implemented of Scope interface
func (e *EntryPoint) Metric() metric.Scope {
	return e.set.Metric
}

// Tracer returns instance implemented of opentracing.Tracer interface
func (e *EntryPoint) Tracer() tracing.Tracer {
	return e.set.Tracer
}

// Executor provide interface for set builder and runner callback functions
func (e *EntryPoint) Executor(builder func(ctx context.Context) error, runner func(ctx context.Context) error) {
	e.builder = builder
	e.runner = runner
}

// Serve execute builder and runner functions with callback for pre run
func (e *EntryPoint) Serve(preRun func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(error); ok {
				err = errors.Wrap(re, "recovered error:")
			} else {
				err = errors.Errorf("unhandled panic, err: %v", r)
			}
		}
	}()
	if e.builder == nil || e.runner == nil {
		return ErrExecutorNotPresent
	}
	if err := e.builder(e.OnShutdown()); err != nil {
		return err
	}
	if preRun != nil {
		if err := preRun(); err != nil {
			return err
		}
	}
	if err := e.runner(e.OnShutdown()); err != nil {
		return err
	}
	return nil
}

// Shutdown raise shutdown event.
func (e *EntryPoint) Shutdown(ctx context.Context, code int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.cancelFn()
	if _, ok := ctx.Deadline(); ok {
		<-ctx.Done()
	}
	os.Exit(code)
}

// Reload raise reload event.
func (e *EntryPoint) Reload() {
	e.invoker.Reload(e.OnShutdown())
}

// WorkDir returns current work directory
func (e *EntryPoint) WorkDir() string {
	return e.initial.WorkDir
}

// CtxExtractInitial returns initial settings from context
func CtxExtractInitial(ctx context.Context) (config.Initial, bool) {
	if initial, ok := ctx.Value(CtxKeyInitial).(*config.Initial); ok {
		return *initial, true
	}
	return config.Initial{}, false
}
