package benchmarks

import (
	"io"
	"testing"

	"github.com/nite-coder/blackbear/pkg/log"
)

// We pass Attrs (or zap.Fields) inline because it affects allocations: building
// up a list outside of the benchmarked code and passing it in with "..."
// reduces measured allocations.

func BenchmarkAttrs(b *testing.B) {
	//ctx := context.Background()
	opts := log.HandlerOptions{
		Level:       log.DebugLevel,
		DisableTime: true,
	}

	for _, handler := range []struct {
		name string
		h    log.Handler
	}{
		{"disabled", disabledHandler{}},
		{"JSON discard", log.NewJSONHandler(io.Discard, &opts)},
	} {
		logger := log.New(handler.h)
		b.Run(handler.name, func(b *testing.B) {
			for _, call := range []struct {
				name string
				f    func()
			}{
				{
					// The number should match nAttrsInline in log/record.go.
					// This should exercise the code path where no allocations
					// happen in Record or Attr. If there are allocations, they
					// should only be from Duration.String and Time.String.
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
				// {
				// 	"5 args ctx",
				// 	func() {
				// 		logger.LogAttrs(ctx, log.LevelInfo, TestMessage,
				// 			log.String("string", TestString),
				// 			log.Int("status", TestInt),
				// 			log.Duration("duration", TestDuration),
				// 			log.Time("time", TestTime),
				// 			log.Any("error", TestError),
				// 		)
				// 	},
				// },
				// {
				// 	"10 args",
				// 	func() {
				// 		logger.LogAttrs(nil, log.LevelInfo, TestMessage,
				// 			log.String("string", TestString),
				// 			log.Int("status", TestInt),
				// 			log.Duration("duration", TestDuration),
				// 			log.Time("time", TestTime),
				// 			log.Any("error", TestError),
				// 			log.String("string", TestString),
				// 			log.Int("status", TestInt),
				// 			log.Duration("duration", TestDuration),
				// 			log.Time("time", TestTime),
				// 			log.Any("error", TestError),
				// 		)
				// 	},
				// },
				// {
				// 	"40 args",
				// 	func() {
				// 		logger.LogAttrs(nil, log.LevelInfo, TestMessage,
				// 			log.String("string", TestString),
				// 			log.Int("status", TestInt),
				// 			log.Duration("duration", TestDuration),
				// 			log.Time("time", TestTime),
				// 			log.Any("error", TestError),
				// 			log.String("string", TestString),
				// 			log.Int("status", TestInt),
				// 			log.Duration("duration", TestDuration),
				// 			log.Time("time", TestTime),
				// 			log.Any("error", TestError),
				// 			log.String("string", TestString),
				// 			log.Int("status", TestInt),
				// 			log.Duration("duration", TestDuration),
				// 			log.Time("time", TestTime),
				// 			log.Any("error", TestError),
				// 			log.String("string", TestString),
				// 			log.Int("status", TestInt),
				// 			log.Duration("duration", TestDuration),
				// 			log.Time("time", TestTime),
				// 			log.Any("error", TestError),
				// 			log.String("string", TestString),
				// 			log.Int("status", TestInt),
				// 			log.Duration("duration", TestDuration),
				// 			log.Time("time", TestTime),
				// 			log.Any("error", TestError),
				// 			log.String("string", TestString),
				// 			log.Int("status", TestInt),
				// 			log.Duration("duration", TestDuration),
				// 			log.Time("time", TestTime),
				// 			log.Any("error", TestError),
				// 			log.String("string", TestString),
				// 			log.Int("status", TestInt),
				// 			log.Duration("duration", TestDuration),
				// 			log.Time("time", TestTime),
				// 			log.Any("error", TestError),
				// 			log.String("string", TestString),
				// 			log.Int("status", TestInt),
				// 			log.Duration("duration", TestDuration),
				// 			log.Time("time", TestTime),
				// 			log.Any("error", TestError),
				// 		)
				// 	},
				// },
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

func BenchmarkWithoutFields(b *testing.B) {
	b.Logf("Logging without any structured context.")

	b.Run("blackbear/log", func(b *testing.B) {
		opts := log.HandlerOptions{
			Level:       log.DebugLevel,
			DisableTime: true,
		}
		log.SetDefault(log.New(log.NewJSONHandler(io.Discard, &opts)))

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				log.Info().Msg(getMessage(0))
			}
		})
	})
}
