package log

import (
	"context"
	"time"
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

	e.CreatedAt = time.Now()

	_ = l.handler.Handle(ctx, e)
	putEntry(e)
}

// Debug level formatted message
func (l *Logger) Debug() *Entry {
	if !l.handler.Enabled(context.TODO(), DebugLevel) {
		return nil
	}
	return newEntry(DebugLevel, l)
}

// Info level formatted message
func (l *Logger) Info() *Entry {
	if !l.handler.Enabled(context.TODO(), InfoLevel) {
		return nil
	}
	return newEntry(InfoLevel, l)
}

// Warn level formatted message
func (l *Logger) Warn() *Entry {
	if !l.handler.Enabled(context.TODO(), WarnLevel) {
		return nil
	}
	return newEntry(WarnLevel, l)
}

// Error level formatted message
func (l *Logger) Error() *Entry {
	if !l.handler.Enabled(context.TODO(), ErrorLevel) {
		return nil
	}
	return newEntry(ErrorLevel, l)
}

// Panic level formatted message
func (l *Logger) Panic() *Entry {
	if !l.handler.Enabled(context.TODO(), PanicLevel) {
		return nil
	}
	return newEntry(PanicLevel, l)
}

// Fatal level formatted message, followed by an exit.
func (l *Logger) Fatal() *Entry {
	if !l.handler.Enabled(context.TODO(), FatalLevel) {
		return nil
	}
	return newEntry(FatalLevel, l)
}

func (l *Logger) With() Context {
	return newContext(l.clone())
}

// WithContext return a new context with a log context value
func (l *Logger) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey, l)
}
