package config

import (
	"github.com/spf13/viper"
	"strings"
)

type Viper struct {
	*viper.Viper
	envKeyReplacer *strings.Replacer
	envPrefix      string
	settings       []CfgItem
}

// SetEnvPrefix defines a prefix that ENVIRONMENT variables will use.
// E.g. if your prefix is "spf", the env registry will look for env
// variables that start with "SPF_".
func (v *Viper) SetEnvPrefix(in string) {
	if in != "" {
		v.Viper.SetEnvPrefix(in)
	}
}

// SetEnvKeyReplacer sets the strings.Replacer on the viper object
// Useful for mapping an environmental variable to a key that does
// not match it.
func (v *Viper) SetEnvKeyReplacer(r *strings.Replacer) {
	v.envKeyReplacer = r
	v.Viper.SetEnvKeyReplacer(r)
}

// EnvPrefix returns env prefix
func (v *Viper) EnvPrefix() string {
	return v.envPrefix
}

// EnvPrefix returns strings.Replacer
func (v *Viper) EnvKeyReplacer() *strings.Replacer {
	return v.envKeyReplacer
}

// AllEnrichedSettings returns all settings with enriched info
func (v *Viper) AllEnrichedSettings() []CfgItem {
	return v.settings
}

// NewViper
func NewViper() *Viper {
	return &Viper{
		Viper: viper.New(),
	}
}
