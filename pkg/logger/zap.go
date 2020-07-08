package logger

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"strings"
)

// Zap is uber/zap logger implemented of Logger interface
type Zap struct {
	ctx     context.Context
	cfg     *Config
	logger  *zap.SugaredLogger
	fields  map[string]interface{}
	tags    []string
	tagsMap []string
}

// Printf is like fmt.Printf, push to log entry with debug level
func (z *Zap) Printf(format string, a ...interface{}) {
	z.Debug(format, Args(a...))
}

// Verbose should return true when verbose logging output is wanted
func (z *Zap) Verbose() bool {
	return z.cfg.Verbose
}

// Emergency push to log entry with emergency level & throw panic
func (z *Zap) Emergency(format string, opts ...Option) {
	z.Log(LevelEmergency, format, opts...)
}

// Alert push to log entry with alert level
func (z *Zap) Alert(format string, opts ...Option) {
	z.Log(LevelAlert, format, opts...)
}

// Critical push to log entry with critical level
func (z *Zap) Critical(format string, opts ...Option) {
	z.Log(LevelCritical, format, opts...)
}

// Error push to log entry with error level
func (z *Zap) Error(format string, opts ...Option) {
	z.Log(LevelError, format, opts...)
}

// Warning push to log entry with warning level
func (z *Zap) Warning(format string, opts ...Option) {
	z.Log(LevelWarning, format, opts...)
}

// Notice push to log entry with notice level
func (z *Zap) Notice(format string, opts ...Option) {
	z.Log(LevelNotice, format, opts...)
}

// Info push to log entry with info level
func (z *Zap) Info(format string, opts ...Option) {
	z.Log(LevelInfo, format, opts...)
}

// Debug push to log entry with debug level
func (z *Zap) Debug(format string, opts ...Option) {
	z.Log(LevelDebug, format, opts...)
}

// Write push to log entry with debug level
func (z *Zap) Write(p []byte) (n int, err error) {
	z.Debug(string(p))
	return len(p), nil
}

// Log push to log with specified level
func (z *Zap) Log(level Level, format string, o ...Option) {
	opts := &opts{}
	for _, option := range o {
		_ = option(opts)
	}
	if !opts.ignoreLevelFilter && level > z.cfg.Level {
		return
	}
	var (
		wargs = []interface{}{"level", level.String()}
		tags  []string
	)
	// tags
	if len(z.tags) > 0 {
		tags = make([]string, len(z.tags))
		copy(tags, z.tags)
	}
	if len(opts.tags) > 0 {
		tags = append(tags, opts.tags...)
	}
	if len(tags) > 0 {
		wargs = append(wargs, "tags", tags)
	}
	// fields
	for k, v := range z.fields {
		wargs = append(wargs, k, v)
	}
	for k, v := range opts.fields {
		wargs = append(wargs, k, v)
	}
	if len(opts.wargs) > 0 {
		wargs = append(wargs, opts.wargs...)
		if len(opts.wargs)%2 != 0 {
			wargs = append(wargs, "%!EXTRA")
		}
	}
	if opts.stack != nil {
		wargs = append(wargs, zap.Stack(*opts.stack))
	}
	if opts.pfields != nil {
		for key, val := range opts.pfields {
			wargs = append(wargs, zap.Any(key, val))
		}
	}
	var logger = z.logger
	if len(wargs) > 0 {
		logger = logger.With(wargs...)
	}
	if !z.pass(level, tags, wargs) {
		return
	}
	if len(opts.args) == 0 {
		var fn func(args ...interface{})
		switch level {
		default:
			fn = logger.Debug
		case LevelInfo, LevelNotice:
			fn = logger.Info
		case LevelWarning:
			fn = logger.Warn
		case LevelError, LevelCritical, LevelAlert:
			fn = logger.Error
		case LevelEmergency:
			fn = logger.Panic
		}
		fn(format)
	} else {
		var fn func(format string, args ...interface{})
		switch level {
		default:
			fn = logger.Debugf
		case LevelInfo, LevelNotice:
			fn = logger.Infof
		case LevelWarning:
			fn = logger.Warnf
		case LevelError, LevelCritical, LevelAlert:
			fn = logger.Errorf
		case LevelEmergency:
			fn = logger.Panicf
		}
		fn(format, opts.args...)
	}
}

func (z *Zap) pass(level Level, tags []string, wargs []interface{}) bool {
	var stop int8
	if len(z.cfg.DebugTags) > 0 {
		for _, tiv := range tags {
			for _, requiredTag := range z.cfg.DebugTags {
				if tiv == requiredTag {
					return true
				}
			}
		}
		stop |= 1
	}
	if len(z.tagsMap) > 0 {
		var k interface{}
		for i, v := range wargs {
			if i%2 == 0 {
				k = v
				continue
			}
			if kk, ok := k.(string); ok {
				if vv, ok := v.(string); ok {
					var mk string
					for mi, mv := range z.tagsMap {
						if mi%2 == 0 {
							mk = mv
							continue
						}
						if kk == mk && vv == mv {
							return true
						}
					}

				}
			}
		}
		stop |= 2
	}
	return stop == 0
}

// WithFields create new instance with fields
func (z *Zap) WithFields(fields Fields) Logger {
	nz := &Zap{}
	copyZap(nz, z, nil, fields)
	return nz
}

// WithTags create new instance with tags
func (z *Zap) WithTags(tags Tags) Logger {
	nz := &Zap{}
	copyZap(nz, z, tags, nil)
	return nz
}

func copyZap(dst, src *Zap, tags []string, fields map[string]interface{}) {
	// tags
	cTags := make([]string, len(src.tags))
	copy(cTags, src.tags)
	dst.tags = cTags
	if tags != nil {
		dst.tags = append(dst.tags, tags...)
	}
	var cFields = map[string]interface{}{}
	// fields
	for k, v := range src.fields {
		cFields[k] = v
	}
	dst.fields = cFields
	for k, v := range fields {
		dst.fields[k] = v
	}
	dst.logger = src.logger
	dst.ctx = src.ctx
	dst.cfg = src.cfg
	dst.tagsMap = src.tagsMap
}

func parseTagsMap(cfg *Config) []string {
	var tagsMap []string
	if len(cfg.DebugTags) > 0 {
		for _, requiredTag := range cfg.DebugTags {
			rt := strings.Split(requiredTag, cfg.MapTagsSplitSep)
			if len(rt) == 2 {
				tagsMap = append(tagsMap, rt...)
			}
		}
	}
	return tagsMap
}

func cfgLevelToZap(lvl Level) zapcore.Level {
	var level zapcore.Level
	switch lvl {
	default:
		level = zap.InfoLevel
	case LevelWarning:
		level = zap.WarnLevel
	case LevelError, LevelCritical, LevelAlert:
		level = zap.ErrorLevel
	case LevelEmergency:
		level = zap.PanicLevel
	}
	return level
}

// NewZap returns uber/zap logger instance implemented of Logger interface
func NewZap(ctx context.Context, cfg *Config) *Zap {
	var (
		logger *zap.Logger
	)
	level := cfgLevelToZap(cfg.Level)
	tagsMap := parseTagsMap(cfg)
	if !cfg.Debug {
		zCfg := zap.Config{
			Level:       zap.NewAtomicLevelAt(level),
			Development: false,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			Encoding:         "json",
			EncoderConfig:    zap.NewProductionEncoderConfig(),
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}
		zCfg.EncoderConfig.LevelKey = ""
		zCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		logger, _ = zCfg.Build(zap.AddCallerSkip(2), zap.AddStacktrace(zap.PanicLevel))
	} else {
		cfg.Level = LevelDebug
		cfg.RedirectLevel = LevelDebug
		zCfg := zap.Config{
			Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
			Development:      true,
			Encoding:         "console",
			EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}
		zCfg.EncoderConfig.LevelKey = ""
		zCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		logger, _ = zCfg.Build(zap.AddCallerSkip(2))
	}
	go func(logger *zap.Logger) {
		<-ctx.Done()
		_ = logger.Sync()
	}(logger)
	copyCfg := *cfg
	z := &Zap{ctx: ctx, cfg: &copyCfg, logger: logger.Sugar(), tagsMap: tagsMap}
	if !copyCfg.DisableRedirectStdLog {
		log.SetOutput(&loggerWriter{
			redirectLevel: &copyCfg.RedirectLevel,
			logFunc:       z.Log,
		})
	}
	cfg.OnReload(func(ctx context.Context) {
		z.cfg.DebugTags = cfg.DebugTags
		z.tagsMap = parseTagsMap(cfg)
	})
	return z
}

//WrapLogger just wraps zap logger without unnecessary actions and return logger
func WrapLogger(ctx context.Context, logger *zap.Logger, cfg *Config) *Zap {
	tagsMap := parseTagsMap(cfg)
	copyCfg := *cfg
	z := &Zap{ctx: ctx, cfg: &copyCfg, logger: logger.Sugar(), tagsMap: tagsMap}

	return z
}