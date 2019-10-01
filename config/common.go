package config

import (
	"fmt"
	"github.com/ProtocolONE/go-core/invoker"
	"github.com/gurukami/typ/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/shurcooL/graphql/ident"
	"os"
	"reflect"
	"strings"
)

var usage = map[string]string{}

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

func valueFromHumanizeEnvPath(v *Viper, key string, path []string) interface{} {
	envPath := make([]string, len(path))
	for i, key := range path {
		name := ident.ParseMixedCaps(key)
		envPath[i] = strings.Join(name, EnvSep)
	}
	prefix := ""
	if p := v.EnvPrefix(); len(p) > 0 {
		prefix = p + EnvSep
	}
	if val, ok := os.LookupEnv(strings.ToUpper(prefix+strings.Join(envPath, EnvSep))); ok {
		return val
	}
	return nil
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
		case reflect.Struct:
			if e := bindValues(v, disableBindMixedCapsEnv, fieldv.Interface(), path...); e != nil {
				return e
			}
		default:
			var val interface{}
			key := strings.Join(path, BindEnvSep)
			//
			if err := v.BindEnv(key); err != nil {
				return err
			}
			// bind mixed caps keys name to humanize ENV
			if !disableBindMixedCapsEnv {
				val = valueFromHumanizeEnvPath(v, key, path)
			}
			// bind to exact ENV name
			if envConfigValue, testEnvConfig := t.Tag.Lookup(LookupEnvConfigTag); testEnvConfig {
				if v, ok := os.LookupEnv(envConfigValue); ok {
					val = v
				}
			}
			if val == nil {
				val = v.Get(key)
			}
			//
			if val != nil {
				v.Set(key, val)
			}
			vv := typ.Of(val)
			//
			var (
				defaultValue, requiredValue, fallbackValue string
				testDefault, testRequired, testFallback    bool
			)
			defaultValue, testDefault = t.Tag.Lookup(LookupDefaultTag)
			requiredValue, testRequired = t.Tag.Lookup(LookupRequiredTag)
			fallbackValue, testFallback = t.Tag.Lookup(LookupFallbackTag)
			required := testRequired && typ.StringBoolHumanize(requiredValue).V()
			if testDefault && required {
				e := fmt.Errorf("ambiguous usage in %v, only one of 'default' or 'required' should specified", key)
				return errors.WithMessage(e, Prefix)
			}
			if vv.Empty().V() {
				if required {
					e := fmt.Errorf("option %v is required", key)
					return errors.WithMessage(e, Prefix)
				}
				if testFallback {
					if v.IsSet(fallbackValue) {
						if e := setValue(v, ErrFallbackPlaceholder, key, fieldv, v.Get(fallbackValue)); e != nil {
							return e
						}
					}
				} else if testDefault {
					if e := setValue(v, ErrDefaultPlaceholder, key, fieldv, defaultValue); e != nil {
						return e
					}
				}
			}
			if v, ok := t.Tag.Lookup("usage"); ok {
				usage[key] = v
			}
		}
	}
	return nil
}