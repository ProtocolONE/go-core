package invoker

import "context"

const (
	Prefix       = "go-core.invoker"
	UnmarshalKey = "invoker"
)

type Observer interface {
	OnReload(callback func(ctx context.Context))
}

type Reloader interface {
	Reload(ctx context.Context)
}
