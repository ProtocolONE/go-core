package config

import (
	"fmt"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/gurukami/typ/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/shurcooL/graphql/ident"
	"os"
	"reflect"
	"strings"
)

// CfgItem
type CfgItem struct {
	Key      string
	ENV      []string
	Usage    string
	Default  string
	Type     string
	Required bool
	Fallback string
	Value    interface{}
}

const (
	Prefix                    = "go-core.config"
	LookupTag                 = "mapstructure"
	LookupDefaultTag          = "default"
	LookupRequiredTag         = "required"
	LookupFallbackTag         = "fallback"
	LookupEnvConfigTag        = "envconfig"
	StringToSliceSep          = ","
	BindEnvSep                = "."
	EnvSep                    = "_"
	ErrDefaultPlaceholder     = "default value for option %v can't set, error occurred: %v"
	ErrFallbackPlaceholder    = "fallback value for option %v can't set, error occurred: %v"
	UnmarshalKeyDebug         = "shared.debug"          // Do not change, usage as fallback
	UnmarshalKeyConfigFile    = "shared.path"           // Do not change, usage as fallback
	UnmarshalKeyGracefulDelay = "shared.graceful.delay" // Do not change, usage as fallback
)

var ErrUnmarshalNotStruct = errors.New("given value under interface not a struct")

type (
	DecodeHookFunc = mapstructure.DecodeHookFunc
	Settings       map[string]interface{}
)

type Initial struct {
	Viper   *Viper
	WorkDir string
	// Disable to bind mixed caps keys name to humanize ENV someEnvHere -> SOME_ENV_HERE
	DisableBindMixedCapsEnv bool
}

type Configurator interface {
	WorkDir() string
	UnmarshalKey(key string, rawVal interface{}, hook ...DecodeHookFunc) error
	UnmarshalKeyOnReload(key string, reloader invoker.Reloader, hook ...DecodeHookFunc) error
}

func setValue(v *Viper, tpl string, key string, rv reflect.Value, defaultValue interface{}) error {
	tv := typ.Of(defaultValue)
	var errString string
	switch rv.Kind() {
	case reflect.Bool:
		nv := tv.Bool()
		if nv.Err() != nil {
			errString = nv.Err().Error()
		} else {
			v.Set(key, nv.V())
		}
	case reflect.Int:
		nv := tv.Int()
		if nv.Err() != nil {
			errString = nv.Err().Error()
		} else {
			v.Set(key, nv.V())
		}
	case reflect.Int8:
		nv := tv.Int8()
		if nv.Err() != nil {
			errString = nv.Err().Error()
		} else {
			v.Set(key, nv.V())
		}
	case reflect.Int16:
		nv := tv.Int16()
		if nv.Err() != nil {
			errString = nv.Err().Error()
		} else {
			v.Set(key, nv.V())
		}
	case reflect.Int32:
		nv := tv.Int32()
		if nv.Err() != nil {
			errString = nv.Err().Error()
		} else {
			v.Set(key, nv.V())
		}
	case reflect.Int64:
		nv := tv.Int64()
		if nv.Err() != nil {
			errString = nv.Err().Error()
		} else {
			v.Set(key, nv.V())
		}
	case reflect.Uint:
		nv := tv.Uint()
		if nv.Err() != nil {
			errString = nv.Err().Error()
		} else {
			v.Set(key, nv.V())
		}
	case reflect.Uint8:
		nv := tv.Uint8()
		if nv.Err() != nil {
			errString = nv.Err().Error()
		} else {
			v.Set(key, nv.V())
		}
	case reflect.Uint16:
		nv := tv.Uint16()
		if nv.Err() != nil {
			errString = nv.Err().Error()
		} else {
			v.Set(key, nv.V())
		}
	case reflect.Uint32:
		nv := tv.Uint32()
		if nv.Err() != nil {
			errString = nv.Err().Error()
		} else {
			v.Set(key, nv.V())
		}
	case reflect.Uint64:
		nv := tv.Uint64()
		if nv.Err() != nil {
			errString = nv.Err().Error()
		} else {
			v.Set(key, nv.V())
		}
	case reflect.Float32:
		nv := tv.Float32()
		if nv.Err() != nil {
			errString = nv.Err().Error()
		} else {
			v.Set(key, nv.V())
		}
	case reflect.Float64:
		nv := tv.Float()
		if nv.Err() != nil {
			errString = nv.Err().Error()
		} else {
			v.Set(key, nv.V())
		}
	case reflect.Complex64:
		nv := tv.Complex64()
		if nv.Err() != nil {
			errString = nv.Err().Error()
		} else {
			v.Set(key, nv.V())
		}
	case reflect.Complex128:
		nv := tv.Complex()
		if nv.Err() != nil {
			errString = nv.Err().Error()
		} else {
			v.Set(key, nv.V())
		}
	case reflect.String:
		nv := tv.String()
		if nv.Err() != nil {
			errString = nv.Err().Error()
		} else {
			v.Set(key, nv.V())
		}
	default:
		errString = "type " + rv.Kind().String() + " not supported"
	}
	if errString != "" {
		e := fmt.Errorf(tpl, key, errString)
		return errors.WithMessage(e, Prefix)
	}
	return nil
}

func valueFromHumanizeEnvPath(v *Viper, path []string) (envKey string, val interface{}) {
	envPath := make([]string, len(path))
	for i, key := range path {
		name := ident.ParseMixedCaps(key)
		envPath[i] = strings.Join(name, EnvSep)
	}
	prefix := ""
	if p := v.EnvPrefix(); len(p) > 0 {
		prefix = p + EnvSep
	}
	envKey = strings.ToUpper(prefix + strings.Join(envPath, EnvSep))
	if v.EnvKeyReplacer() != nil {
		envKey = v.EnvKeyReplacer().Replace(envKey)
	}
	if val, ok := os.LookupEnv(envKey); ok {
		return envKey, val
	}
	return
}

func bindValues(v *Viper, disableBindMixedCapsEnv bool, iface interface{}, parts ...string) error {
	ifv := reflect.Indirect(reflect.ValueOf(iface))
	ift := ifv.Type()
	if ift.Kind() != reflect.Struct {
		return ErrUnmarshalNotStruct
	}
	for i := 0; i < ift.NumField(); i++ {
		fieldv := ifv.Field(i)
		t := ift.Field(i)
		name := t.Name
		if name[:1] == strings.ToLower(name[:1]) {
			continue
		}
		tag, ok := t.Tag.Lookup(LookupTag)
		if ok {
			name = tag
		}
		path := append(parts, name)
		switch fieldv.Kind() {
		case reflect.Interface:
			continue
		case reflect.Ptr:
			if e := bindValues(v, disableBindMixedCapsEnv, reflect.Zero(fieldv.Type().Elem()).Interface(), path...); e != nil {
				return e
			}
		case reflect.Struct:
			if e := bindValues(v, disableBindMixedCapsEnv, fieldv.Interface(), path...); e != nil {
				return e
			}
		default:
			var (
				item   CfgItem
				envKey string
			)
			item.Type = fieldv.Type().String()
			item.Key = strings.Join(path, BindEnvSep)
			// skip binding if already bind
			if v.findCfgItemByName(item.Key) != nil {
				return nil
			}
			//
			if err := v.BindEnv(item.Key); err != nil {
				return err
			}
			// bind mixed caps keys name to humanize ENV
			if !disableBindMixedCapsEnv {
				envKey, item.Value = valueFromHumanizeEnvPath(v, path)
				item.ENV = append(item.ENV, envKey)
			}
			// bind to exact ENV name
			if envConfigValue, testEnvConfig := t.Tag.Lookup(LookupEnvConfigTag); testEnvConfig {
				item.ENV = append(item.ENV, envConfigValue)
				if v, ok := os.LookupEnv(envConfigValue); ok {
					item.Value = v
				}
			}
			if item.Value == nil {
				item.Value = v.Get(item.Key)
			}
			//
			if item.Value != nil {
				v.Set(item.Key, item.Value)
			}
			vv := typ.Of(item.Value)
			//
			var (
				requiredValue                           string
				testDefault, testRequired, testFallback bool
			)
			item.Default, testDefault = t.Tag.Lookup(LookupDefaultTag)
			requiredValue, testRequired = t.Tag.Lookup(LookupRequiredTag)
			item.Fallback, testFallback = t.Tag.Lookup(LookupFallbackTag)
			item.Required = testRequired && typ.StringBoolHumanize(requiredValue).V()
			if testDefault && item.Required {
				e := fmt.Errorf("ambiguous usage in %v, only one of 'default' or 'required' should specified", item.Key)
				return errors.WithMessage(e, Prefix)
			}
			if vv.Empty().V() {
				if item.Required {
					e := fmt.Errorf("option %v is required", item.Key)
					return errors.WithMessage(e, Prefix)
				}
				if testFallback {
					if v.IsSet(item.Fallback) {
						if e := setValue(v, ErrFallbackPlaceholder, item.Key, fieldv, v.Get(item.Fallback)); e != nil {
							return e
						}
					}
				} else if testDefault {
					if e := setValue(v, ErrDefaultPlaceholder, item.Key, fieldv, item.Default); e != nil {
						return e
					}
				}
			}
			if v, ok := t.Tag.Lookup("usage"); ok {
				item.Usage = v
			}
			item.Value = v.Get(item.Key)
			v.settings = append(v.settings, item)
		}
	}
	return nil
}
