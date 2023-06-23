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

	// ErrorHandler is called whenever handler fails to write an event on its
	// output. If not set, an error is printed on the stderr. This handler must
	// be thread safe and non-blocking.
	ErrorHandler func(err error)

	// AutoStaceTrace add stack trace into entry when use `Error`, `Panic`, `Fatal` level.
	// Default: true
	AutoStaceTrace = true
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

// Debug level formatted message
func Debug() *Entry {
	return Default().Debug()
}

// Info level formatted message
func Info() *Entry {
	return Default().Info()
}

// Warn level formatted message
func Warn() *Entry {
	return Default().Warn()
}

// Error level formatted message
func Error() *Entry {
	return Default().Error()
}

// Panic level formatted message
func Panic() *Entry {
	return Default().Panic()
}

// Fatal level formatted message, followed by an exit.
func Fatal() *Entry {
	return Default().Fatal()
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
	// for _, h := range Logger().handles {
	// 	flusher, ok := h.(Flusher)
	// 	if ok {
	// 		err := flusher.Flush()
	// 		if err != nil {
	// 			stdlog.Printf("log: flush log handler: %v", err)
	// 		}
	// 	}
	// }
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
