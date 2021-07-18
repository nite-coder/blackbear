package log

import (
	"context"
	"fmt"
	stdlog "log"
	"runtime"
	"strings"
	"sync/atomic"
)

// Logger is the default instance of the log package
var (
	atomicLogger = new()

	// ErrorHandler is called whenever handler fails to write an event on its
	// output. If not set, an error is printed on the stderr. This handler must
	// be thread safe and non-blocking.
	ErrorHandler func(err error)

	// AutoStaceTrace add stack trace into entry when use `Error`, `Panic`, `Fatal` level.
	// Default: true
	AutoStaceTrace = true
)

func new() atomic.Value {
	logger := New()

	v := atomic.Value{}
	v.Store(logger)
	return v
}

func SetLogger(logger *logger) *logger {
	if logger == nil {
		return nil
	}
	atomicLogger.Store(logger)
	return logger
}

func Logger() *logger {
	logger, _ := atomicLogger.Load().(*logger)
	return logger
}

// Debug level formatted message
func Debug(msg string) {
	Logger().Debug(msg)
}

// Debugf level formatted message
func Debugf(msg string, v ...interface{}) {
	Logger().Debugf(msg, v...)
}

// Info level formatted message
func Info(msg string) {
	Logger().Info(msg)
}

// Infof level formatted message
func Infof(msg string, v ...interface{}) {
	Logger().Infof(msg, v...)
}

// Warn level formatted message
func Warn(msg string) {
	Logger().Warn(msg)
}

// Warnf level formatted message
func Warnf(msg string, v ...interface{}) {
	Logger().Warnf(msg, v...)
}

// Error level formatted message
func Error(msg string) {
	Logger().Error(msg)
}

// Errorf level formatted message
func Errorf(msg string, v ...interface{}) {
	Logger().Errorf(msg, v...)
}

// Panic level formatted message
func Panic(msg string) {
	Logger().Panic(msg)
}

// Panicf level formatted message
func Panicf(msg string, v ...interface{}) {
	Logger().Panicf(msg, v...)
}

// Fatal level formatted message, followed by an exit.
func Fatal(msg string) {
	Logger().Fatal(msg)
}

// Fatalf level formatted message, followed by an exit.
func Fatalf(msg string, v ...interface{}) {
	Logger().Fatalf(msg, v...)
}

// Str add string field to current context
func Str(key string, val string) Context {
	return Logger().Str(key, val)
}

// Bool add bool field to current context
func Bool(key string, val bool) Context {
	return Logger().Bool(key, val)
}

// Int add Int field to current context
func Int(key string, val int) Context {
	return Logger().Int(key, val)
}

// Int8 add Int8 field to current context
func Int8(key string, val int8) Context {
	return Logger().Int8(key, val)
}

// Int16 add Int16 field to current context
func Int16(key string, val int16) Context {
	return Logger().Int16(key, val)
}

// Int32 add Int32 field to current context
func Int32(key string, val int32) Context {
	return Logger().Int32(key, val)
}

// Int64 add Int64 field to current context
func Int64(key string, val int64) Context {
	return Logger().Int64(key, val)
}

// Uint add Uint field to current context
func Uint(key string, val uint) Context {
	return Logger().Uint(key, val)
}

// Uint8 add Uint8 field to current context
func Uint8(key string, val uint8) Context {
	return Logger().Uint8(key, val)
}

// Uint16 add Uint16 field to current context
func Uint16(key string, val uint16) Context {
	return Logger().Uint16(key, val)
}

// Uint32 add Uint32 field to current context
func Uint32(key string, val uint32) Context {
	return Logger().Uint32(key, val)
}

// Uint64 add Uint64 field to current context
func Uint64(key string, val uint64) Context {
	return Logger().Uint64(key, val)
}

// Float32 add float32 field to current context
func Float32(key string, val float32) Context {
	return Logger().Float32(key, val)
}

// Float64 add Float64 field to current context
func Float64(key string, val float64) Context {
	return Logger().Float64(key, val)
}

// Err add error field to current context
func Err(err error) Context {
	return Logger().Err(err)
}

var (
	ctxKey = &struct {
		name string
	}{
		name: "log",
	}
)

// FromContext return a log context from the standard context
func FromContext(ctx context.Context) Context {
	_logger := Logger()
	v := ctx.Value(ctxKey)
	if v == nil {
		return newContext(_logger)
	}

	return v.(Context)
}

// Flush clear all handler's buffer
func Flush() {
	for _, h := range Logger().handles {
		flusher, ok := h.(Flusher)
		if ok {
			err := flusher.Flush()
			if err != nil {
				stdlog.Printf("log: flush log handler: %v", err)
			}
		}
	}
}

func getStackTrace() string {
	stackBuf := make([]uintptr, 50)
	length := runtime.Callers(3, stackBuf)
	stack := stackBuf[:length]

	var b strings.Builder
	frames := runtime.CallersFrames(stack)

	for {
		frame, more := frames.Next()

		if !strings.Contains(frame.File, "runtime/") {
			_, _ = b.WriteString(fmt.Sprintf("\n\tFile: %s, Line: %d. Function: %s", frame.File, frame.Line, frame.Function))
		}

		if !more {
			break
		}
	}
	return b.String()
}
