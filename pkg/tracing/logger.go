package tracing

import (
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	jaegerLogger "github.com/uber/jaeger-client-go/log"
)

type logAdapter struct {
	redirectLevel logger.Level
	logFunc       func(level logger.Level, format string, o ...logger.Option)
}

// Error logs a message at error priority
func (l *logAdapter) Error(msg string) {
	l.logFunc(logger.LevelError, msg)
}

// Infof logs a message at info priority
func (l *logAdapter) Infof(format string, v ...interface{}) {
	l.logFunc(logger.LevelInfo, format, logger.Args(v...))
}

// NewLoggerAdapter
func NewLoggerAdapter(log logger.Logger) jaegerLogger.Logger {
	return &logAdapter{
		logFunc: log.Log,
	}
}
