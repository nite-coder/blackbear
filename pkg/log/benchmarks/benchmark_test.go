package benchmarks

import (
	"context"
	"io"
	"testing"

	"github.com/nite-coder/blackbear/pkg/log"
)

// We pass Attrs (or zap.Fields) inline because it affects allocations: building
// up a list outside of the benchmarked code and passing it in with "..."
// reduces measured allocations.

func BenchmarkHandles(b *testing.B) {
	ctx := context.Background()
	opts := log.HandlerOptions{
		Level:       log.DebugLevel,
		DisableTime: true,
	}

	for _, handler := range []struct {
		name string
		h    log.Handler
	}{
		{"disabled", disabledHandler{}},
		{"Text discard", log.NewTextHandler(io.Discard, &opts)},
		{"JSON discard", log.NewJSONHandler(io.Discard, &opts)},
	} {
		logger := log.New(handler.h)
		b.Run(handler.name, func(b *testing.B) {
			for _, call := range []struct {
				name string
				f    func()
			}{
				{
					"0 args",
					func() {
						logger.Info().
							Msg(TestMessage)

					},
				},
				{
					"5 args",
					func() {
						logger.Info().
							Str("string", TestString).
							Int("int", TestInt).
							Duration("duration", TestDuration).
							Time("time", TestTime).
							Any("error", TestError).
							Msg(TestMessage)

					},
				},
				{
					"5 args ctx",
					func() {
						logger.InfoCtx(ctx).
							Str("string", TestString).
							Int("int", TestInt).
							Duration("duration", TestDuration).
							Time("time", TestTime).
							Any("error", TestError).
							Msg(TestMessage)
					},
				},
				{
					"10 args",
					func() {
						logger.Info().
							Str("string", TestString).
							Int("int", TestInt).
							Duration("duration", TestDuration).
							Time("time", TestTime).
							Any("error", TestError).
							Str("string", TestString).
							Int("int", TestInt).
							Duration("duration", TestDuration).
							Time("time", TestTime).
							Any("error", TestError).
							Msg(TestMessage)
					},
				},
				{
					"40 args",
					func() {
						logger.Info().
							Str("string", TestString).
							Int("int", TestInt).
							Duration("duration", TestDuration).
							Time("time", TestTime).
							Any("error", TestError).
							Str("string", TestString).
							Int("int", TestInt).
							Duration("duration", TestDuration).
							Time("time", TestTime).
							Any("error", TestError).
							Str("string", TestString).
							Int("int", TestInt).
							Duration("duration", TestDuration).
							Time("time", TestTime).
							Any("error", TestError).
							Str("string", TestString).
							Int("int", TestInt).
							Duration("duration", TestDuration).
							Time("time", TestTime).
							Any("error", TestError).
							Str("string", TestString).
							Int("int", TestInt).
							Duration("duration", TestDuration).
							Time("time", TestTime).
							Any("error", TestError).
							Str("string", TestString).
							Int("int", TestInt).
							Duration("duration", TestDuration).
							Time("time", TestTime).
							Any("error", TestError).
							Str("string", TestString).
							Int("int", TestInt).
							Duration("duration", TestDuration).
							Time("time", TestTime).
							Any("error", TestError).
							Str("string", TestString).
							Int("int", TestInt).
							Duration("duration", TestDuration).
							Time("time", TestTime).
							Any("error", TestError).
							Msg(TestMessage)
					},
				},
			} {
				b.Run(call.name, func(b *testing.B) {
					b.ReportAllocs()
					b.RunParallel(func(pb *testing.PB) {
						for pb.Next() {
							call.f()
						}
					})
				})
			}
		})
	}
}
