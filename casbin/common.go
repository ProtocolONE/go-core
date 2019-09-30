package casbin

import (
	"context"
	"github.com/ProtocolONE/go-core/invoker"
)

const (
	UnmarshalKey = "casbin"
)

// Config is a general casbin config settings
type Config struct {
	ModelConfPath  string
	PolicyConfPath string
	invoker        *invoker.Invoker
}

// OnReload
func (c *Config) OnReload(callback func(ctx context.Context)) {
	c.invoker.OnReload(callback)
}

// Reload
func (c *Config) Reload(ctx context.Context) {
	c.invoker.Reload(ctx)
}

type (
	Model  string
	Policy string
)

func (m Model) String() string {
	return string(m)
}

func (p Policy) String() string {
	return string(p)
}
