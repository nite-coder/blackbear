package log

import (
	"time"
)

// Context use for meta data
type Context struct {
	logger *Logger
	fields []*Field
}

func newContext(l *Logger) Context {
	l.context = Context{
		logger: l,
		fields: make([]*Field, 0, 5),
	}

	return l.context
}

// Logger returns the logger with the context previously set.
func (c Context) Logger() *Logger {
	return c.logger
}

// Str add string field to current context
func (c Context) Str(key string, val string) Context {
	f := Field{
		Key:   key,
		Value: val,
	}

	c.fields = append(c.fields, &f)
	c.logger.context = c
	return c
}

// Strs add string field to current context
func (c Context) Strs(key string, val []string) Context {
	f := Field{
		Key:   key,
		Value: val,
	}

	c.fields = append(c.fields, &f)
	c.logger.context = c
	return c
}

// Bool add bool field to current context
func (c Context) Bool(key string, val bool) Context {
	f := Field{
		Key:   key,
		Value: val,
	}

	c.fields = append(c.fields, &f)
	c.logger.context = c
	return c
}

// Int add Int field to current context
func (c Context) Int(key string, val int) Context {
	f := Field{
		Key:   key,
		Value: val,
	}

	c.fields = append(c.fields, &f)
	c.logger.context = c
	return c
}

// Ints add Int field to current context
func (c Context) Ints(key string, val []int) Context {
	f := Field{
		Key:   key,
		Value: val,
	}

	c.fields = append(c.fields, &f)
	c.logger.context = c
	return c
}

// Int8 add Int8 field to current context
func (c Context) Int8(key string, val int8) Context {
	f := Field{
		Key:   key,
		Value: val,
	}

	c.fields = append(c.fields, &f)
	c.logger.context = c
	return c
}

// Int16 add Int16 field to current context
func (c Context) Int16(key string, val int16) Context {
	f := Field{
		Key:   key,
		Value: val,
	}

	c.fields = append(c.fields, &f)
	c.logger.context = c
	return c
}

// Int32 add Int32 field to current context
func (c Context) Int32(key string, val int32) Context {
	f := Field{
		Key:   key,
		Value: val,
	}

	c.fields = append(c.fields, &f)
	c.logger.context = c
	return c
}

// Int64 add Int64 field to current context
func (c Context) Int64(key string, val int64) Context {
	f := Field{
		Key:   key,
		Value: val,
	}

	c.fields = append(c.fields, &f)
	c.logger.context = c
	return c
}

// Uint add Uint field to current context
func (c Context) Uint(key string, val uint) Context {
	f := Field{
		Key:   key,
		Value: val,
	}

	c.fields = append(c.fields, &f)
	c.logger.context = c
	return c
}

// Uint8 add Uint8 field to current context
func (c Context) Uint8(key string, val uint8) Context {
	f := Field{
		Key:   key,
		Value: val,
	}

	c.fields = append(c.fields, &f)
	c.logger.context = c
	return c
}

// Uint16 add Uint16 field to current context
func (c Context) Uint16(key string, val uint16) Context {
	f := Field{
		Key:   key,
		Value: val,
	}

	c.fields = append(c.fields, &f)
	c.logger.context = c
	return c
}

// Uint32 add Uint32 field to current context
func (c Context) Uint32(key string, val uint32) Context {
	f := Field{
		Key:   key,
		Value: val,
	}

	c.fields = append(c.fields, &f)
	c.logger.context = c
	return c
}

// Uint64 add Uint64 field to current context
func (c Context) Uint64(key string, val uint64) Context {
	f := Field{
		Key:   key,
		Value: val,
	}

	c.fields = append(c.fields, &f)
	c.logger.context = c
	return c
}

// Float32 add float32 field to current context
func (c Context) Float32(key string, val float32) Context {
	f := Field{
		Key:   key,
		Value: val,
	}

	c.fields = append(c.fields, &f)
	c.logger.context = c
	return c
}

// Float64 add Float64 field to current context
func (c Context) Float64(key string, val float64) Context {
	f := Field{
		Key:   key,
		Value: val,
	}

	c.fields = append(c.fields, &f)
	c.logger.context = c
	return c
}

// Err add error field to current context
func (c Context) Err(err error) Context {
	f := Field{
		Key:   "error",
		Value: err.Error(),
	}

	c.fields = append(c.fields, &f)
	c.logger.context = c
	return c
}

// StackTrace adds stack_trace field to the current context
func (c Context) StackTrace() Context {
	f := Field{
		Key:   "stack_trace",
		Value: getStackTrace(),
	}

	c.fields = append(c.fields, &f)
	c.logger.context = c
	return c
}

// Time adds Time field to current context
func (c Context) Time(key string, val time.Time) Context {
	f := Field{
		Key:   key,
		Value: val,
	}

	c.fields = append(c.fields, &f)
	c.logger.context = c
	return c
}

// Times adds Time field to current context
func (c Context) Times(key string, val []time.Time) Context {
	f := Field{
		Key:   key,
		Value: val,
	}

	c.fields = append(c.fields, &f)
	c.logger.context = c
	return c
}

// Any adds the field key with i marshaled using reflection.
func (c Context) Any(key string, val interface{}) Context {
	f := Field{
		Key:   key,
		Value: val,
	}

	c.fields = append(c.fields, &f)
	c.logger.context = c
	return c
}
