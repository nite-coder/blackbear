package log

import (
	"context"
)

// Flusher is an interface that allow handles have the ability to clear buffer and close connection
type Flusher interface {
	Flush() error
}

type Logger struct {
	handler Handler
	context Context
}

func New(handler Handler) *Logger {
	logger := Logger{
		handler: handler,
	}

	return &logger
}

func (l *Logger) clone() *Logger {
	c := *l
	return &c
}

func (l *Logger) log(ctx context.Context, e *Entry) {
	if ctx == nil {
		ctx = context.Background()
	}

	_ = l.handler.Handle(ctx, e)
	putEntry(e)
}

// Debug logs at DebugLevel
func (l *Logger) Debug() *Entry {
	if !l.handler.Enabled(context.TODO(), DebugLevel) {
		return nil
	}
	return newEntry(context.TODO(), DebugLevel, l)
}

// DebugCtx logs at LevelDebug with the given context
func (l *Logger) DebugCtx(ctx context.Context) *Entry {
	if !l.handler.Enabled(ctx, DebugLevel) {
		return nil
	}
	return newEntry(ctx, DebugLevel, l)
}

// Info logs at InfoLevel
func (l *Logger) Info() *Entry {
	if !l.handler.Enabled(context.TODO(), InfoLevel) {
		return nil
	}
	return newEntry(context.TODO(), InfoLevel, l)
}

// InfoCtx logs at InfoLevel with the given context
func (l *Logger) InfoCtx(ctx context.Context) *Entry {
	if !l.handler.Enabled(ctx, InfoLevel) {
		return nil
	}
	return newEntry(ctx, InfoLevel, l)
}

// Warn logs at WarnLevel
func (l *Logger) Warn() *Entry {
	if !l.handler.Enabled(context.TODO(), WarnLevel) {
		return nil
	}
	return newEntry(context.TODO(), WarnLevel, l)
}

// WarnCtx logs at WarnLevel with the given context.
func (l *Logger) WarnCtx(ctx context.Context) *Entry {
	if !l.handler.Enabled(context.TODO(), WarnLevel) {
		return nil
	}
	return newEntry(ctx, WarnLevel, l)
}

// Error logs at ErrorLevel
func (l *Logger) Error() *Entry {
	if !l.handler.Enabled(context.TODO(), ErrorLevel) {
		return nil
	}
	return newEntry(context.TODO(), ErrorLevel, l)
}

// ErrorCtx logs at ErrorLevel with the given context.
func (l *Logger) ErrorCtx(ctx context.Context) *Entry {
	if !l.handler.Enabled(context.TODO(), ErrorLevel) {
		return nil
	}
	return newEntry(ctx, ErrorLevel, l)
}

// Panic logs at PanicLevel
func (l *Logger) Panic() *Entry {
	if !l.handler.Enabled(context.TODO(), PanicLevel) {
		return nil
	}
	return newEntry(context.TODO(), PanicLevel, l)
}

// PanicCtx logs at PanicLevel with the given context.
func (l *Logger) PanicCtx(ctx context.Context) *Entry {
	if !l.handler.Enabled(context.TODO(), PanicLevel) {
		return nil
	}
	return newEntry(ctx, PanicLevel, l)
}

// Fatal logs at FatalLevel
func (l *Logger) Fatal() *Entry {
	if !l.handler.Enabled(context.TODO(), FatalLevel) {
		return nil
	}
	return newEntry(context.TODO(), FatalLevel, l)
}

// FatalCtx logs at FatalLevel with the given context.
func (l *Logger) FatalCtx(ctx context.Context) *Entry {
	if !l.handler.Enabled(context.TODO(), FatalLevel) {
		return nil
	}
	return newEntry(ctx, FatalLevel, l)
}

func (l *Logger) With() Context {
	return newContext(l.clone())
}

// WithContext return a new context with a log context value
func (l *Logger) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey, l)
}
