package log

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
)

// Logger is the default instance of the log package
var (
	defaultLogger = new()


)

func new() atomic.Value {
	opts := HandlerOptions{
		Level: DebugLevel,
	}
	jsonHandler := NewJSONHandler(os.Stderr, &opts)
	logger := New(jsonHandler)

	v := atomic.Value{}
	v.Store(logger)
	return v
}

// SetDefault makes l the default Logger.
func SetDefault(logger *Logger) *Logger {
	if logger == nil {
		return nil
	}
	defaultLogger.Store(logger)
	return logger
}

func Default() *Logger {
	logger, _ := defaultLogger.Load().(*Logger)
	return logger
}

// Debug logs at DebugLevel
func Debug() *Entry {
	return Default().Debug()
}

// DebugCtx logs at LevelDebug with the given context
func DebugCtx(ctx context.Context) *Entry {
	return Default().DebugCtx(ctx)
}

// Info logs at InfoLevel
func Info() *Entry {
	return Default().Info()
}

// InfoCtx logs at InfoLevel with the given context
func InfoCtx(ctx context.Context) *Entry {
	return Default().InfoCtx(ctx)
}

// Warn logs at WarnLevel
func Warn() *Entry {
	return Default().Warn()
}

// WarnCtx logs at WarnLevel with the given context.
func WarnCtx(ctx context.Context) *Entry {
	return Default().WarnCtx(ctx)
}

// Error level formatted message
func Error() *Entry {
	return Default().Error()
}

// ErrorCtx logs at ErrorLevel with the given context.
func ErrorCtx(ctx context.Context) *Entry {
	return Default().ErrorCtx(ctx)
}

// Panic level formatted message
func Panic() *Entry {
	return Default().Panic()
}

// PanicCtx logs at PanicLevel with the given context.
func PanicCtx(ctx context.Context) *Entry {
	return Default().PanicCtx(ctx)
}

// Fatal level formatted message, followed by an exit.
func Fatal() *Entry {
	return Default().Fatal()
}

// FatalCtx logs at FatalLevel with the given context.
func FatalCtx(ctx context.Context) *Entry {
	return Default().FatalCtx(ctx)
}

func With() Context {
	return Default().With()
}

var (
	ctxKey = &struct {
		name string
	}{
		name: "log",
	}
)

// FromContext return a logger from the standard context
func FromContext(ctx context.Context) *Logger {
	v := ctx.Value(ctxKey)
	if v == nil {
		return Default()
	}

	logger, ok := v.(*Logger)
	if !ok {
		return Default()
	}

	return logger
}

// Flush clear all handler's buffer
func Flush() {
	h := Default().handler
	flusher, ok := h.(Flusher)
	if ok {
		_ = flusher.Flush()
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
