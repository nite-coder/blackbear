package log

import (
	"sync"
)

// Handler is an interface that log handlers need to be implemented
type Handler interface {
	BeforeWriting(*Entry) error
	Write([]byte) error
}

// Flusher is an interface that allow handles have the ability to clear buffer and close connection
type Flusher interface {
	Flush() error
}

// Hookfunc is an func that allow us to do something before writing
type Hookfunc func(*Entry) error

type logger struct {
	handles              []Handler
	hooks                []Hookfunc
	leveledHandlers      map[Level][]Handler
	cacheLeveledHandlers func(level Level) []Handler
	rwMutex              sync.RWMutex
	buf                  []byte

	DisableTimeField bool
}

func New() *logger {
	logger := logger{
		leveledHandlers: map[Level][]Handler{},
	}

	logger.cacheLeveledHandlers = logger.getLeveledHandlers()
	return &logger
}

// Debug level formatted message
func (l *logger) Debug(msg string) {
	e := newEntry(l, l.buf)
	e.Debug(msg)
}

// Debugf level formatted message
func (l *logger) Debugf(msg string, v ...interface{}) {
	e := newEntry(l, l.buf)
	e.Debugf(msg, v...)
}

// Info level formatted message
func (l *logger) Info(msg string) {
	e := newEntry(l, l.buf)
	e.Info(msg)
}

// Infof level formatted message
func (l *logger) Infof(msg string, v ...interface{}) {
	e := newEntry(l, l.buf)
	e.Infof(msg, v...)
}

// Warn level formatted message
func (l *logger) Warn(msg string) {
	e := newEntry(l, l.buf)
	e.Warn(msg)
}

// Warnf level formatted message
func (l *logger) Warnf(msg string, v ...interface{}) {
	e := newEntry(l, l.buf)
	e.Warnf(msg, v...)
}

// Error level formatted message
func (l *logger) Error(msg string) {
	e := newEntry(l, l.buf)
	e.Error(msg)
}

// Errorf level formatted message
func (l *logger) Errorf(msg string, v ...interface{}) {
	e := newEntry(l, l.buf)
	e.Errorf(msg, v...)
}

// Panic level formatted message
func (l *logger) Panic(msg string) {
	e := newEntry(l, l.buf)
	e.Panic(msg)
}

// Panicf level formatted message
func (l *logger) Panicf(msg string, v ...interface{}) {
	e := newEntry(l, l.buf)
	e.Panicf(msg, v...)
}

// Fatal level formatted message, followed by an exit.
func (l *logger) Fatal(msg string) {
	e := newEntry(l, l.buf)
	e.Fatal(msg)
}

// Fatalf level formatted message, followed by an exit.
func (l *logger) Fatalf(msg string, v ...interface{}) {
	e := newEntry(l, l.buf)
	e.Fatalf(msg, v...)
}

// Trace returns a new entry with a Stop method to fire off
// a corresponding completion log, useful with defer.
func (l *logger) Trace(msg string) *Entry {
	e := newEntry(l, l.buf)
	return e.Trace(msg)
}

// Str add string field to current context
func (l *logger) Str(key string, val string) Context {
	c := newContext(l)
	return c.Str(key, val)
}

// Bool add bool field to current context
func (l *logger) Bool(key string, val bool) Context {
	c := newContext(l)
	return c.Bool(key, val)
}

// Int add Int field to current context
func (l *logger) Int(key string, val int) Context {
	c := newContext(l)
	return c.Int(key, val)
}

// Int8 add Int8 field to current context
func (l *logger) Int8(key string, val int8) Context {
	c := newContext(l)
	return c.Int8(key, val)
}

// Int16 add Int16 field to current context
func (l *logger) Int16(key string, val int16) Context {
	c := newContext(l)
	return c.Int16(key, val)
}

// Int32 add Int32 field to current context
func (l *logger) Int32(key string, val int32) Context {
	c := newContext(l)
	return c.Int32(key, val)
}

// Int64 add Int64 field to current context
func (l *logger) Int64(key string, val int64) Context {
	c := newContext(l)
	return c.Int64(key, val)
}

// Uint add Uint field to current context
func (l *logger) Uint(key string, val uint) Context {
	c := newContext(l)
	return c.Uint(key, val)
}

// Uint8 add Uint8 field to current context
func (l *logger) Uint8(key string, val uint8) Context {
	c := newContext(l)
	return c.Uint8(key, val)
}

// Uint16 add Uint16 field to current context
func (l *logger) Uint16(key string, val uint16) Context {
	c := newContext(l)
	return c.Uint16(key, val)
}

// Uint32 add Uint32 field to current context
func (l *logger) Uint32(key string, val uint32) Context {
	c := newContext(l)
	return c.Uint32(key, val)
}

// Uint64 add Uint64 field to current context
func (l *logger) Uint64(key string, val uint64) Context {
	c := newContext(l)
	return c.Uint64(key, val)
}

// Float32 add float32 field to current context
func (l *logger) Float32(key string, val float32) Context {
	c := newContext(l)
	return c.Float32(key, val)
}

// Float64 add Float64 field to current context
func (l *logger) Float64(key string, val float64) Context {
	c := newContext(l)
	return c.Float64(key, val)
}

// Any add val field to current context
func (l *logger) Any(key string, val interface{}) Context {
	c := newContext(l)
	return c.Any(key, val)
}

// Err add error field to current context
func (l *logger) Err(err error) Context {
	c := newContext(l)
	return c.Err(err)
}

// AddHandler adds a new Log Handler and specifies what log levels
// the handler will be passed log entries for
func (l *logger) AddHandler(handler Handler, levels ...Level) *logger {
	l.rwMutex.Lock()
	defer l.rwMutex.Unlock()

	for _, level := range levels {
		l.leveledHandlers[level] = append(l.leveledHandlers[level], handler)
	}

	l.handles = append(l.handles, handler)
	l.cacheLeveledHandlers = l.getLeveledHandlers()

	return l
}

// RemoveAllHandlers removes all handlers
func (l *logger) RemoveAllHandlers() {
	l.rwMutex.Lock()
	defer l.rwMutex.Unlock()

	l.leveledHandlers = map[Level][]Handler{}
	l.handles = []Handler{}
	l.hooks = []Hookfunc{}
	l.cacheLeveledHandlers = l.getLeveledHandlers()
}

// AddHook adds a new Hook to log entry
func (l *logger) AddHook(hook Hookfunc) error {
	l.rwMutex.Lock()
	defer l.rwMutex.Unlock()

	l.hooks = append(l.hooks, hook)
	return nil
}

func (l *logger) getLeveledHandlers() func(level Level) []Handler {
	debugHandlers := l.leveledHandlers[DebugLevel]
	infoHandlers := l.leveledHandlers[InfoLevel]
	warnHandlers := l.leveledHandlers[WarnLevel]
	errorHandlers := l.leveledHandlers[ErrorLevel]
	panicHandlers := l.leveledHandlers[PanicLevel]
	fatalHandlers := l.leveledHandlers[FatalLevel]
	traceHandlers := l.leveledHandlers[TraceLevel]

	return func(level Level) []Handler {
		switch level {
		case DebugLevel:
			return debugHandlers
		case InfoLevel:
			return infoHandlers
		case WarnLevel:
			return warnHandlers
		case ErrorLevel:
			return errorHandlers
		case PanicLevel:
			return panicHandlers
		case FatalLevel:
			return fatalHandlers
		case TraceLevel:
			return traceHandlers
		}

		return []Handler{}
	}
}
