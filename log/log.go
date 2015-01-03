package log

import (
	"fmt"
	"time"
)

const (
	DEBUG int32 = iota
	INFO
	WARN
	ERROR
)

var levelFlags = map[int32]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
}

type Log interface {
	// Debug logs information that is diagnostically helpful to developers.
	Debug(format string, args ...interface{})

	// Info logs useful information to log.
	Info(format string, args ...interface{})

	// Error logs any error which is fatal to the operation.
	Error(format string, args ...interface{})

	// Warn logs message with severity "warn". Anything that can potentially cause error.
	Warn(format string, args ...interface{})

	// Disable will prevent the application to log anything.
	Disable()

	// SetLevel sets the error reporting level.
	SetLevel(level int32)
}

type SimpleLogger struct {
	disabled bool
	level    int32
}

// SimpleLogger is a basic Log implementation which sends the output to stdout.
// It is possible to use a custom logger, just need to implement 'Log' interface
// and call Logger() function.
var Logger Log = new(SimpleLogger)

func (l *SimpleLogger) log(level string, format string, args ...interface{}) {
	fmt.Printf("[%s] %s - ", level, time.Now().In(time.UTC).Format("2006-01-02T15:04:05Z07:00"))
	fmt.Println(fmt.Sprintf(format, args...))
}

// Return the flag level name based on provided iota.
func GetLevelFlagName(level int32) string {
	return levelFlags[level]
}

func (l *SimpleLogger) Debug(format string, args ...interface{}) {
	if !l.disabled && l.level == DEBUG {
		l.log("DEBUG", format, args...)
	}
}

func (l *SimpleLogger) Info(format string, args ...interface{}) {
	if !l.disabled && (l.level >= INFO || l.level == DEBUG) {
		l.log("INFO", format, args...)
	}
}

func (l *SimpleLogger) Warn(format string, args ...interface{}) {
	if !l.disabled && (l.level >= WARN || l.level == DEBUG) {
		l.log("WARN", format, args...)
	}
}

func (l *SimpleLogger) Error(format string, args ...interface{}) {
	if !l.disabled && (l.level >= ERROR || l.level == DEBUG) {
		l.log("ERROR", format, args...)
	}
}

func (l *SimpleLogger) Disable() {
	l.disabled = true
}

func (l *SimpleLogger) SetLevel(level int32) {
	l.level = level
}
