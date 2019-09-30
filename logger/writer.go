package logger

import (
	"bytes"
	"io"
)

type loggerWriter struct {
	redirectLevel *Level
	logFunc       func(level Level, format string, o ...Option)
}

// Write
func (l *loggerWriter) Write(p []byte) (int, error) {
	p = bytes.TrimSpace(p)
	l.logFunc(*l.redirectLevel, string(p))
	return len(p), nil
}

// NewLevelWriter provide adapter from Logger instance to io.Writer with custom level
func NewLevelWriter(logger Logger, lvl Level) io.Writer {
	return &loggerWriter{
		redirectLevel: &lvl,
		logFunc:       logger.Log,
	}
}
