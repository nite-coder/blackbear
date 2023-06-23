package log

import "strings"

// Level of the log
type Level uint8

// Log levels.
const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
	TraceLevel
)

var levelNames = []string{
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
	"PANIC",
	"FATAL",
	"TRACE",
}

// String returns the string representation of a logging level.
func (p Level) String() string {
	return levelNames[p]
}

// NewLevel returns Level struct
func NewLevel(level string) Level {
	level = strings.ToLower(level)
	switch level {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "panic":
		return PanicLevel
	case "fatal":
		return FatalLevel
	default:
		return DebugLevel
	}
}
