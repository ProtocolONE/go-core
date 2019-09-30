package casbin

import (
	"github.com/ProtocolONE/go-core/config"
	"github.com/ProtocolONE/go-core/invoker"
	"github.com/casbin/casbin"
	"github.com/google/wire"
	"github.com/qiangmzsx/string-adapter"
)

// ProviderCfg returns configuration for production GORM
func ProviderCfg(cfg config.Configurator) (*Config, func(), error) {
	c := &Config{
		ModelConfPath:  cfg.WorkDir() + "/casbin/model.conf",
		PolicyConfPath: cfg.WorkDir() + "/casbin/policy.csv",
		invoker: invoker.NewInvoker(),
	}
	e := cfg.UnmarshalKeyOnReload(UnmarshalKey, c)
	return c, func() {}, e
}

// Provider returns casbin.Enforcer instance with resolved dependencies
func Provider(cfg *Config) (*casbin.Enforcer, func(), error) {
	enf := casbin.NewEnforcer(cfg.ModelConfPath, cfg.PolicyConfPath)
	return enf, func() {}, nil
}

// ProviderTest returns stub/mock casbin.Enforcer instance with resolved dependencies
func ProviderTest(model Model, policy Policy) (*casbin.Enforcer, func(), error) {
	enf := casbin.NewEnforcer(casbin.NewModel(model.String()), string_adapter.NewAdapter(policy.String()))
	return enf, func() {}, nil
}

var (
	WireSet     = wire.NewSet(Provider, ProviderCfg)
	WireTestSet = wire.NewSet(ProviderTest)
)
