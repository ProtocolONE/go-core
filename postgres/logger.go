package postgres

import (
	"bytes"
	"fmt"
	"github.com/ProtocolONE/go-core/logger"
)

type loggerAdapter struct {
	level  logger.Level
	logger logger.Logger
}

type AdapterLogger interface {
	// Print format & print log
	Print(v ...interface{})
}

// Print write string to output
func (l *loggerAdapter) Print(v ...interface{}) {
	var str bytes.Buffer
	for _, s := range v {
		if str.Len() > 0 {
			str.WriteString(" ")
		}
		str.WriteString(fmt.Sprint(s))
	}
	l.logger.Log(l.level, "%v", logger.Args(str.String()))
}

// NewLoggerAdapter returns instance adapter for services/logger
func NewLoggerAdapter(logger logger.Logger, level logger.Level) AdapterLogger {
	return &loggerAdapter{level: level, logger: logger}
}
