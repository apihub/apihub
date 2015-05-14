package main

import (
	"github.com/backstage/backstage/log"
	"github.com/ccding/go-logging/logging"
)

type CustomLogger struct {
	disabled bool
	logger   *logging.Logger
}

func NewCustomLogger() *CustomLogger {
	l := new(CustomLogger)
	l.logger, _ = logging.SimpleLogger("main")
	return l
}

func (l *CustomLogger) Debug(format string, args ...interface{}) {
	if !l.disabled {
		l.logger.Debugf(format, args...)
	}
}

func (l *CustomLogger) Info(format string, args ...interface{}) {
	if !l.disabled {
		l.logger.Infof(format, args...)
	}
}

func (l *CustomLogger) Warn(format string, args ...interface{}) {
	if !l.disabled {
		l.logger.Warnf(format, args...)
	}
}

func (l *CustomLogger) Error(format string, args ...interface{}) {
	if !l.disabled {
		l.logger.Errorf(format, args...)
	}
}

func (l *CustomLogger) Disable() {
	l.disabled = true
}

func (l *CustomLogger) SetLevel(level int32) {
	levelStr := log.GetLevelFlagName(level)
	l.logger.SetLevel(logging.GetLevelValue(levelStr))
}
