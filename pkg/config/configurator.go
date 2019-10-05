package config

import (
	"context"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"sync"
)

var mu sync.Mutex

type ProductionConfigurator struct {
	viper    *Viper
	initial  Initial
	observer invoker.Observer
}

// WorkDir returns current work directory
func (p *ProductionConfigurator) WorkDir() string {
	return p.initial.WorkDir
}

// UnmarshalKeyOnReload
func (p *ProductionConfigurator) UnmarshalKeyOnReload(key string, reloader invoker.Reloader, hook ...DecodeHookFunc) error {
	if p.observer != nil {
		p.observer.OnReload(func(ctx context.Context) {
			_ = p.UnmarshalKey(key, reloader, hook...)
			reloader.Reload(ctx)
		})
	}
	return p.UnmarshalKey(key, reloader, hook...)
}

// UnmarshalKey
func (p *ProductionConfigurator) UnmarshalKey(key string, rawVal interface{}, hook ...DecodeHookFunc) error {
	mu.Lock()
	defer mu.Unlock()
	if e := bindValues(p.viper, p.initial.DisableBindMixedCapsEnv, rawVal, key); e != nil {
		return e
	}
	hook = append(hook,
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(StringToSliceSep),
	)
	return p.viper.UnmarshalKey(key, rawVal, viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
			hook...
		),
	))
}

// NewProductionConfigurator
func NewProductionConfigurator(initial Initial, observer invoker.Observer) (Configurator, error) {
	v := initial.Viper
	if v == nil {
		v = NewViper()
	}
	for _, key := range v.AllKeys() {
		val := v.Get(key)
		v.Set(key, val)
	}
	return &ProductionConfigurator{
		viper:    v,
		initial:  initial,
		observer: observer,
	}, nil
}
