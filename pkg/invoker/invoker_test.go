package invoker

import (
	"context"
	"testing"
)

func TestClose(t *testing.T) {
	inv := NewInvoker()
	inv.Close()
}

func TestReload(t *testing.T) {
	inv := NewInvoker()
	inv.OnReload(func(_ context.Context) {
		// do nothing
	})
	inv.Reload(context.Background())
}
