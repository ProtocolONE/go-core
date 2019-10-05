package invoker

import (
	"context"
	"runtime"
	"sync"
)

type Dispatcher struct {
	reloadCh chan struct{}
	ctx      context.Context
}

type Invoker struct {
	mu       sync.Mutex
	dispatch *Dispatcher
	closeCh  chan struct{}
}

func (i *Invoker) init() {
	i.mu.Lock()
	defer i.mu.Unlock()
	if i.closeCh == nil {
		i.closeCh = make(chan struct{})
	}
	if i.dispatch == nil {
		i.dispatch = &Dispatcher{
			reloadCh: make(chan struct{}),
		}
	}
}

// OnReload
func (i *Invoker) OnReload(callback func(ctx context.Context)) {
	i.init()
	go func() {
		// fmt.Println("run on reload")
		// defer fmt.Println("stop on reload")
		for {
			i.mu.Lock()
			dispatch := i.dispatch
			closeCh := i.closeCh
			reloadCh := dispatch.reloadCh
			i.mu.Unlock()
			select {
			case <-closeCh:
				return
			case <-reloadCh:
				i.mu.Lock()
				ctx := dispatch.ctx
				i.mu.Unlock()
				callback(ctx)
			}
		}
	}()
}

// Reload
func (i *Invoker) Reload(ctx context.Context) {
	i.init()
	i.mu.Lock()
	defer i.mu.Unlock()
	i.dispatch.ctx = ctx
	ch := i.dispatch.reloadCh
	i.dispatch = &Dispatcher{
		reloadCh: make(chan struct{}),
	}
	close(ch)
}

// Close unsubscribe all listeners
func (i *Invoker) Close() error {
	i.mu.Lock()
	if i.closeCh != nil {
		close(i.closeCh)
	}
	i.mu.Unlock()
	return nil
}

// OnClose
func (i *Invoker) OnClose() <-chan struct{} {
	i.init()
	return i.closeCh
}

// Observer
func (i *Invoker) Observer() Observer {
	return i
}

// NewInvoker
func NewInvoker() *Invoker {
	inv := &Invoker{}
	runtime.SetFinalizer(inv, func(i *Invoker) {
		i.Close()
	})
	return inv
}
