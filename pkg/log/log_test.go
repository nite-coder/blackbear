package log_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/nite-coder/blackbear/internal/buffer"
	"github.com/nite-coder/blackbear/pkg/log"
	"github.com/stretchr/testify/assert"
)

type Person struct {
	Name string
	Age  int
}

// type ErrHandler struct {
// }

// func (h *ErrHandler) BeforeWriting(e *log.Entry) error {
// 	return nil
// }

// func (h *ErrHandler) Write(bytes []byte) error {
// 	return errors.New("oops")
// }

// func TestErrorHandler(t *testing.T) {
// 	logger := log.New()

// 	h1 := &ErrHandler{}
// 	logger.AddHandler(h1, log.AllLevels...)

// 	log.SetLogger(logger)

// 	isErr := false
// 	log.ErrorHandler = func(err error) {
// 		isErr = true
// 	}

// 	log.Debug("aaa")
// 	assert.Equal(t, true, isErr)
// }

func TestNoHandler(t *testing.T) {
	log.Info().Msg("no handler 1")
	log.InfoCtx(context.Background()).Msg("no handler 2")
}

func TestDisableLevel(t *testing.T) {
	b := buffer.New()
	defer b.Free()

	opts := log.HandlerOptions{
		Level:       log.InfoLevel,
		DisableTime: true,
	}
	log.SetDefault(log.New(log.NewJSONHandler(b, &opts)))

	log.Debug().Msg("no handler 1")
	assert.Equal(t, "", b.String())
	b.Reset()
}

func TestLog(t *testing.T) {
	b := buffer.New()
	defer b.Free()

	opts := log.HandlerOptions{
		Level:       log.DebugLevel,
		DisableTime: true,
	}
	log.SetDefault(log.New(log.NewJSONHandler(b, &opts)))

	log.Debug().Msg("debug")
	assert.Equal(t, `{"level":"DEBUG","msg":"debug"}`+"\n", b.String())
	b.Reset()

	log.Debug().Msgf("debug %s", "hello")
	assert.Equal(t, `{"level":"DEBUG","msg":"debug hello"}`+"\n", b.String())
	b.Reset()

	log.Info().Msg("info")
	assert.Equal(t, `{"level":"INFO","msg":"info"}`+"\n", b.String())
	b.Reset()

	log.Info().Msgf("info %s", "hello")
	assert.Equal(t, `{"level":"INFO","msg":"info hello"}`+"\n", b.String())
	b.Reset()

	log.Warn().Msg("warn")
	assert.Equal(t, `{"level":"WARN","msg":"warn"}`+"\n", b.String())
	b.Reset()

	log.Warn().Msgf("warn %s", "hello")
	assert.Equal(t, `{"level":"WARN","msg":"warn hello"}`+"\n", b.String())
	b.Reset()

	log.Error().Msg("error")
	assert.Equal(t, `{"level":"ERROR","msg":"error"}`+"\n", b.String())
	b.Reset()

	log.Error().Msgf("error %s", "hello")
	assert.Equal(t, `{"level":"ERROR","msg":"error hello"}`+"\n", b.String())
	b.Reset()

	t.Run("test panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if ok {
					assert.Equal(t, "panic", err.Error())
				}
			}
			assert.Equal(t, `{"level":"PANIC","msg":"panic"}`+"\n", b.String())
		}()
		log.Panic().Msg("panic")
	})

	b.Reset()
	t.Run("test panicf", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if ok {
					assert.Equal(t, "panic hello", err.Error())
				}
			}
			assert.Equal(t, `{"level":"PANIC","msg":"panic hello"}`+"\n", b.String())
		}()
		log.Panic().Msgf("panic %s", "hello")
	})
}

func TestFields(t *testing.T) {
	b := buffer.New()
	defer b.Free()

	opts := log.HandlerOptions{
		Level:       log.DebugLevel,
		DisableTime: true,
	}
	log.SetDefault(log.New(log.NewJSONHandler(b, &opts)))

	time1, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	time2, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+08:00")

	logger := log.With().
		Str("string", "hello").
		Strs("strs", []string{"str1", "str2"}).
		Bool("bool", true).
		Int("int", 1).
		Int8("int8", int8(2)).
		Int16("int16", int16(3)).
		Int32("int32", int32(4)).
		Int64("int64", int64(5)).
		Uint("uint", uint(6)).
		Uint8("uint8", uint8(7)).
		Uint16("uint16", uint16(8)).
		Uint32("uint32", uint32(9)).
		Uint64("uint64", uint64(10)).
		Float32("float32", float32(11.123)).
		Float64("float64", float64(12.123)).
		Time("time", time1).
		Times("times", []time.Time{time1, time2}).
		Any("person", Person{Name: "Doge", Age: 18}).
		Logger()

	logger.Info().Msg("test field")
	assert.Equal(t, `{"level":"INFO","msg":"test field","string":"hello","strs":["str1","str2"],"bool":true,"int":1,"int8":2,"int16":3,"int32":4,"int64":5,"uint":6,"uint8":7,"uint16":8,"uint32":9,"uint64":10,"float32":11.123,"float64":12.123,"time":"2012-11-01T22:08:41Z","times":["2012-11-01T22:08:41Z","2012-11-01T22:08:41+08:00"],"person":{"Name":"Doge","Age":18}}`+"\n", b.String())
	b.Reset()

	log.Debug().
		Str("string", "hello").
		Strs("strs", []string{"str1", "str2"}).
		Bool("bool", true).
		Int("int", 1).
		Int8("int8", int8(2)).
		Int16("int16", int16(3)).
		Int32("int32", int32(4)).
		Int64("int64", int64(5)).
		Uint("uint", uint(6)).
		Uint8("uint8", uint8(7)).
		Uint16("uint16", uint16(8)).
		Uint32("uint32", uint32(9)).
		Uint64("uint64", uint64(10)).
		Float32("float32", float32(11.123)).
		Float64("float64", float64(12.123)).
		Time("time", time1).
		Times("times", []time.Time{time1, time2}).
		Any("person", Person{Name: "Doge", Age: 18}).
		Msg("test field")

	assert.Equal(t, `{"level":"DEBUG","msg":"test field","string":"hello","strs":["str1","str2"],"bool":true,"int":1,"int8":2,"int16":3,"int32":4,"int64":5,"uint":6,"uint8":7,"uint16":8,"uint32":9,"uint64":10,"float32":11.123,"float64":12.123,"time":"2012-11-01T22:08:41Z","times":["2012-11-01T22:08:41Z","2012-11-01T22:08:41+08:00"],"person":{"Name":"Doge","Age":18}}`+"\n", b.String())
}

func TestFlush(t *testing.T) {
	b := buffer.New()
	defer b.Free()

	opts := log.HandlerOptions{
		Level:       log.DebugLevel,
		DisableTime: true,
	}
	log.SetDefault(log.New(log.NewJSONHandler(b, &opts)))

	log.Debug().Msg("flush")
	log.Flush()
}

func TestStdContext(t *testing.T) {
	b := buffer.New()
	defer b.Free()

	opts := log.HandlerOptions{
		Level:       log.DebugLevel,
		DisableTime: true,
	}
	log.SetDefault(log.New(log.NewJSONHandler(b, &opts)))

	for i := 0; i < 2; i++ {
		logger1 := log.With().Int("num", i).Logger()
		logger1.Debug().Msg("aaa")
	}
	assert.Equal(t, `{"level":"DEBUG","msg":"aaa","num":0}`+"\n"+`{"level":"DEBUG","msg":"aaa","num":1}`+"\n", b.String())
	b.Reset()

	t.Run("create new context", func(t *testing.T) {
		ctx := context.Background()
		ctx = log.With().Str("request_id", "abc").Logger().WithContext(ctx)

		logger := log.FromContext(ctx)
		logger.Debug().Msg("debug")
		assert.Equal(t, `{"level":"DEBUG","msg":"debug","request_id":"abc"}`+"\n", b.String())
		b.Reset()

		logger.Debug().Str("app", "santa").Msgf("debug %s", "hello")
		assert.Equal(t, `{"level":"DEBUG","msg":"debug hello","request_id":"abc","app":"santa"}`+"\n", b.String())
		b.Reset()
	})

	t.Run("create from background context", func(t *testing.T) {
		ctx := context.Background()
		logger := log.FromContext(ctx)
		logger.Info().Msg("test")
		assert.Equal(t, `{"level":"INFO","msg":"test"}`+"\n", b.String())
	})
}

func TestAdvancedFields(t *testing.T) {
	b := buffer.New()
	defer b.Free()

	opts := log.HandlerOptions{
		Level:       log.DebugLevel,
		DisableTime: true,
	}
	log.SetDefault(log.New(log.NewJSONHandler(b, &opts)))

	err := errors.New("something bad happened")
	log.Error().Err(err).Msg("too bad")

	assert.Equal(t, `{"level":"ERROR","msg":"too bad","error":"something bad happened"}`+"\n", b.String())
}

func TestGoroutineSafe(t *testing.T) {
	b := buffer.New()
	defer b.Free()

	opts := log.HandlerOptions{
		Level:       log.DebugLevel,
		DisableTime: true,
	}
	log.SetDefault(log.New(log.NewJSONHandler(b, &opts)))

	logger := log.With().Str("request_id", "abc").Logger()

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		logger.Info().Str("name", "abc").Msg("test")
	}()

	go func() {
		defer wg.Done()
		logger.Info().Str("name", "xyz").Msg("test")
	}()

	go func() {
		defer wg.Done()

		b1 := buffer.New()
		defer b1.Free()

		opts := log.HandlerOptions{
			Level:       log.DebugLevel,
			DisableTime: true,
		}
		log.SetDefault(log.New(log.NewJSONHandler(b1, &opts)))
		log.Info().Msg("test")
	}()

	wg.Wait()
}
