package log

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var entryPool = &sync.Pool{
	New: func() interface{} {
		return &Entry{
			fields: make([]*Field, 0, 5),
		}
	},
}

// Entry defines a single log entry
type Entry struct {
	Logger    *Logger
	Level     Level     `json:"level"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"time"`
	fields    []*Field
}

func newEntry(level Level, l *Logger) *Entry {
	e, _ := entryPool.Get().(*Entry)
	e.Level = level
	e.Logger = l

	if len(l.context.fields) > 0 {
		e.fields = append(e.fields, l.context.fields...)
	}

	return e
}

func putEntry(e *Entry) {
	e.fields = e.fields[:0]
	entryPool.Put(e)
}

// Msg print the message.
func (e *Entry) Msg(msg string) {
	if e == nil {
		return
	}
	e.Message = msg
	e.Logger.log(context.TODO(), e)
}

// Msgf print the formatted message.
func (e *Entry) Msgf(msg string, v ...any) {
	if e == nil {
		return
	}
	e.Message = fmt.Sprintf(msg, v...)
	e.Logger.log(context.TODO(), e)
}

// Str add string field to current entry
func (e *Entry) Str(key string, val string) *Entry {
	if e == nil {
		return e
	}

	f := Field{
		Key:   key,
		Value: val,
	}

	e.fields = append(e.fields, &f)
	return e
}

// Strs add string field to current entry
func (e *Entry) Strs(key string, val []string) *Entry {
	if e == nil {
		return e
	}

	f := Field{
		Key:   key,
		Value: val,
	}

	e.fields = append(e.fields, &f)
	return e
}

// Bool add bool field to current entry
func (e *Entry) Bool(key string, val bool) *Entry {
	if e == nil {
		return e
	}

	f := Field{
		Key:   key,
		Value: val,
	}

	e.fields = append(e.fields, &f)
	return e
}

// Int adds Int field to current entry
func (e *Entry) Int(key string, val int) *Entry {
	if e == nil {
		return e
	}

	f := Field{
		Key:   key,
		Value: val,
	}

	e.fields = append(e.fields, &f)
	return e
}

// Ints adds Int field to current entry
func (e *Entry) Ints(key string, val []int) *Entry {
	if e == nil {
		return e
	}

	f := Field{
		Key:   key,
		Value: val,
	}

	e.fields = append(e.fields, &f)
	return e
}

// Int8 add Int8 field to current entry
func (e *Entry) Int8(key string, val int8) *Entry {
	if e == nil {
		return e
	}

	f := Field{
		Key:   key,
		Value: val,
	}

	e.fields = append(e.fields, &f)
	return e
}

// Int16 add Int16 field to current entry
func (e *Entry) Int16(key string, val int16) *Entry {
	if e == nil {
		return e
	}

	f := Field{
		Key:   key,
		Value: val,
	}

	e.fields = append(e.fields, &f)
	return e
}

// Int32 adds Int32 field to current entry
func (e *Entry) Int32(key string, val int32) *Entry {
	if e == nil {
		return e
	}

	f := Field{
		Key:   key,
		Value: val,
	}

	e.fields = append(e.fields, &f)
	return e
}

// Int64 add Int64 field to current entry
func (e *Entry) Int64(key string, val int64) *Entry {
	if e == nil {
		return e
	}

	f := Field{
		Key:   key,
		Value: val,
	}

	e.fields = append(e.fields, &f)
	return e
}

// Uint add Uint field to current entry
func (e *Entry) Uint(key string, val uint) *Entry {
	if e == nil {
		return e
	}

	f := Field{
		Key:   key,
		Value: val,
	}

	e.fields = append(e.fields, &f)
	return e
}

// Uint8 add Uint8 field to current entry
func (e *Entry) Uint8(key string, val uint8) *Entry {
	if e == nil {
		return e
	}

	f := Field{
		Key:   key,
		Value: val,
	}

	e.fields = append(e.fields, &f)
	return e
}

// Uint16 add Uint16 field to current entry
func (e *Entry) Uint16(key string, val uint16) *Entry {
	if e == nil {
		return e
	}

	f := Field{
		Key:   key,
		Value: val,
	}

	e.fields = append(e.fields, &f)
	return e
}

// Uint32 add Uint32 field to current entry
func (e *Entry) Uint32(key string, val uint32) *Entry {
	if e == nil {
		return e
	}

	f := Field{
		Key:   key,
		Value: val,
	}

	e.fields = append(e.fields, &f)
	return e
}

// Uint64 add Uint64 field to current entry
func (e *Entry) Uint64(key string, val uint64) *Entry {
	if e == nil {
		return e
	}

	f := Field{
		Key:   key,
		Value: val,
	}

	e.fields = append(e.fields, &f)
	return e
}

// Float32 add Float32 field to current entry
func (e *Entry) Float32(key string, val float32) *Entry {
	if e == nil {
		return e
	}

	f := Field{
		Key:   key,
		Value: val,
	}

	e.fields = append(e.fields, &f)
	return e
}

// Float64 adds Float64 field to current entry
func (e *Entry) Float64(key string, val float64) *Entry {
	if e == nil {
		return e
	}

	f := Field{
		Key:   key,
		Value: val,
	}

	e.fields = append(e.fields, &f)
	return e
}

// Time adds Time field to current entry
func (e *Entry) Time(key string, val time.Time) *Entry {
	if e == nil {
		return e
	}

	f := Field{
		Key:   key,
		Value: val,
	}

	e.fields = append(e.fields, &f)
	return e
}

// Times adds Time field to current entry
func (e *Entry) Times(key string, val []time.Time) *Entry {
	if e == nil {
		return e
	}

	f := Field{
		Key:   key,
		Value: val,
	}

	e.fields = append(e.fields, &f)
	return e
}

// Duration adds Duration field to current entry
func (e *Entry) Duration(key string, val time.Duration) *Entry {
	if e == nil {
		return e
	}

	f := Field{
		Key:   key,
		Value: val,
	}

	e.fields = append(e.fields, &f)
	return e
}

// Any adds the field key with i marshaled using reflection.
func (e *Entry) Any(key string, val interface{}) *Entry {
	if e == nil {
		return e
	}

	f := Field{
		Key:   key,
		Value: val,
	}

	e.fields = append(e.fields, &f)

	return e
}

// Err add error field to current context
func (e *Entry) Err(err error) *Entry {
	f := Field{
		Key:   "error",
		Value: err.Error(),
	}

	e.fields = append(e.fields, &f)
	return e
}

// StackTrace adds stack_trace field to the current context
func (e *Entry) StackTrace() *Entry {
	if e == nil {
		return e
	}

	f := Field{
		Key:   "stack_trace",
		Value: getStackTrace(),
	}

	e.fields = append(e.fields, &f)
	return e
}
