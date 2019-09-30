package logger

import (
	"context"
	"github.com/ProtocolONE/go-core/invoker"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

const (
	UnmarshalKey          = "logger"
	UnmarshalKeyDebug     = UnmarshalKey + ".debug"
	UnmarshalKeyVerbose   = UnmarshalKey + ".verbose"
	UnmarshalKeyDebugTags = UnmarshalKey + ".debugTags"
	UnmarshalKeyLevel     = UnmarshalKey + ".level"
)

type (
	Fields map[string]interface{}
	Tags   []string
)

// Level represent RFC5424 logger severity
type Level int8

// String implements interface Stringer
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelNotice:
		return "notice"
	case LevelWarning:
		return "warning"
	case LevelError:
		return "error"
	case LevelCritical:
		return "critical"
	case LevelAlert:
		return "alert"
	case LevelEmergency:
		return "emergency"
	}
	return "unknown"
}

// FromString set logger level from string representation
func (l *Level) FromString(level string) Level {
	switch level {
	default:
		*l = LevelDebug
	case "info":
		*l = LevelInfo
	case "notice":
		*l = LevelNotice
	case "warning":
		*l = LevelWarning
	case "error":
		*l = LevelError
	case "critical":
		*l = LevelCritical
	case "alert":
		*l = LevelAlert
	case "emergency":
		*l = LevelEmergency
	}
	return *l
}

const (
	// Debug: debug-level messages
	LevelDebug Level = 7
	// Informational: informational messages
	LevelInfo Level = 6
	// Notice: normal but significant condition
	LevelNotice Level = 5
	// Warning: warning conditions
	LevelWarning Level = 4
	// Error: error conditions
	LevelError Level = 3
	// Critical: critical conditions
	LevelCritical Level = 2
	// Alert: action must be taken immediately
	LevelAlert Level = 1
	// Emergency: system is unusable
	LevelEmergency Level = 0
)

type opts struct {
	wargs             []interface{}
	args              []interface{}
	tags              Tags
	fields            Fields
	pfields           Fields
	ignoreLevelFilter bool
	stack             *string
}

// Option is func hook for underling logic call
type Option func(*opts) error

// Stack constructs a field that stores a stacktrace of the current goroutine
// under provided key. Keep in mind that taking a stacktrace is eager and
// expensive (relatively speaking); this function both makes an allocation and
// takes about two microseconds.
func Stack(key string) Option {
	return func(f *opts) error {
		f.stack = &key
		return nil
	}
}

// PairArgs returns func hook a logger for to pair args
func PairArgs(a ...interface{}) Option {
	return func(f *opts) error {
		f.wargs = a
		return nil
	}
}

// Args returns func hook a logger for replace fmt placeholders on represent values
func Args(a ...interface{}) Option {
	return func(f *opts) error {
		f.args = a
		return nil
	}
}

// WithPrettyFields returns func hook a logger for adding fields for call with prettifier
func WithPrettyFields(fields Fields) Option {
	return func(f *opts) error {
		f.pfields = fields
		return nil
	}
}

// WithTags returns func hook a logger for adding tags for call
func WithTags(tags Tags) Option {
	return func(f *opts) error {
		f.tags = tags
		return nil
	}
}

// WithFields returns func hook a logger for adding fields for call
func WithFields(fields Fields) Option {
	return func(f *opts) error {
		f.fields = fields
		return nil
	}
}

// IgnoreLevelFilter returns func hook a logger for ignore filter by logger level
func IgnoreLevelFilter() Option {
	return func(f *opts) error {
		f.ignoreLevelFilter = true
		return nil
	}
}

// Config is a general logger config settings
type Config struct {
	Debug                 bool `fallback:"shared.debug"`
	Verbose               bool
	Level                 Level
	DebugTags             []string
	MapTagsSplitSep       string `default:":"`
	DisableRedirectStdLog bool
	RedirectLevel         Level `default:"6"`
	invoker               *invoker.Invoker
}

// OnReload
func (c *Config) OnReload(callback func(ctx context.Context)) {
	c.invoker.OnReload(callback)
}

// Reload
func (c *Config) Reload(ctx context.Context) {
	c.invoker.Reload(ctx)
}

// Logger is the interface for logger client
type Logger interface {
	// Printf is like fmt.Printf, push to log entry with debug level
	Printf(format string, a ...interface{})
	// Verbose should return true when verbose logging output is wanted
	Verbose() bool
	// Emergency push to log entry with emergency level & throw panic
	Emergency(format string, opts ...Option)
	// Alert push to log entry with alert level
	Alert(format string, opts ...Option)
	// Critical push to log entry with critical level
	Critical(format string, opts ...Option)
	// Error push to log entry with error level
	Error(format string, opts ...Option)
	// Warning push to log entry with warning level
	Warning(format string, opts ...Option)
	// Notice push to log entry with notice level
	Notice(format string, opts ...Option)
	// Info push to log entry with info level
	Info(format string, opts ...Option)
	// Debug push to log entry with debug level
	Debug(format string, opts ...Option)
	// Write push to log entry with debug level
	Write(p []byte) (n int, err error)
	// Log push to log with specified level
	Log(level Level, format string, opts ...Option)
	// WithFields create new instance with fields
	WithFields(fields Fields) Logger
	// WithTags create new instance with tags
	WithTags(tags Tags) Logger
}

// StringToLoggerLevelHookFunc returns decoder func hook for converting string representation to RFC5424 level
func StringToLoggerLevelHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(Level(0)) {
			return data, nil
		}
		var l Level
		return (&l).FromString(data.(string)), nil
	}
}
