package config

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"sync"
)

type MockConfigurator struct {
	viper    *Viper
	settings map[string]interface{}
	initial  Initial
	observer invoker.Observer
	mu       sync.Mutex
}

// WorkDir returns current work directory
func (p *MockConfigurator) WorkDir() string {
	return p.initial.WorkDir
}

// UnmarshalKeyOnReload
func (p *MockConfigurator) UnmarshalKeyOnReload(key string, reloader invoker.Reloader, hook ...DecodeHookFunc) error {
	if p.observer != nil {
		p.observer.OnReload(func(ctx context.Context) {
			_ = p.UnmarshalKey(key, reloader, hook...)
			reloader.Reload(ctx)
		})
	}
	return p.UnmarshalKey(key, reloader, hook...)
}

// UnmarshalKey
func (p *MockConfigurator) UnmarshalKey(key string, rawVal interface{}, hook ...DecodeHookFunc) error {
	p.mu.Lock()
	defer p.mu.Unlock()
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

// NewMockConfigurator
func NewMockConfigurator(initial Initial, observer invoker.Observer, settings Settings) (Configurator, error) {
	v := initial.Viper
	if v == nil {
		v = NewViper()
	}
	b, e := json.Marshal(settings)
	if e != nil {
		return nil, e
	}
	v.SetConfigType("json")
	e = v.ReadConfig(bytes.NewBuffer(b))
	if e != nil {
		return nil, e
	}
	return &MockConfigurator{
		viper:    v,
		settings: settings,
		initial:  initial,
		observer: observer,
	}, nil
}
