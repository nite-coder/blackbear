package log

import (
	"context"
	"io"
	"io/ioutil"
	"net"
	"sync"
	"time"

	stdJSON "encoding/json"

	"github.com/nite-coder/blackbear/internal/buffer"
	"github.com/nite-coder/blackbear/internal/json"
)

var enc = json.Encoder{}

func init() {
	// using closure to reflect the changes at runtime.
	json.JSONMarshalFunc = func(v interface{}) ([]byte, error) {
		return stdJSON.Marshal(v)
	}
}

// Handler is an interface that log handlers need to be implemented
type Handler interface {
	Enabled(context.Context, Level) bool
	Handle(context.Context, *Entry) error
}

type HandlerOptions struct {
	Level        Level
	DisableTime  bool
	DisableColor bool
	// ErrorHandler is called whenever handler fails to write an event on its
	// output. If not set, an error is printed on the stderr. This handler must
	// be thread safe and non-blocking.
	ErrorHandler func(err error)

	// AutoStaceTrace add stack trace into entry when use `Error`, `Panic`, `Fatal` level.
	// Default: false
	DisableAutoStaceTrace bool
}

type JSONHandler struct {
	mu   sync.Mutex
	w    io.Writer
	opts *HandlerOptions
}

func NewJSONHandler(w io.Writer, opts *HandlerOptions) *JSONHandler {
	if w == nil {
		w = ioutil.Discard
	}

	return &JSONHandler{
		w:    w,
		opts: opts,
	}
}

// Enabled reports whether the handler handles records at the given level.
// The handler ignores records whose level is lower.
func (h *JSONHandler) Enabled(_ context.Context, level Level) bool {
	return level >= h.opts.Level
}

// Handle formats its argument Record as a JSON object on a single line.
func (h *JSONHandler) Handle(_ context.Context, e *Entry) error {
	buf := buffer.New()
	defer buf.Free()

	*buf = enc.AppendBeginMarker(*buf)

	// time
	if !h.opts.DisableTime {
		*buf = enc.AppendKey(*buf, "time")
		*buf = enc.AppendTime(*buf, time.Now(), time.RFC3339)
	}

	// level
	*buf = enc.AppendKey(*buf, "level")
	*buf = enc.AppendString(*buf, e.Level.String())

	// msg
	if e.Message != "" {
		*buf = enc.AppendKey(*buf, "msg")
		*buf = enc.AppendString(*buf, e.Message)
	}

	// fields
	for _, field := range e.fields {
		*buf = enc.AppendKey(*buf, field.Key)

		switch val := field.Value.(type) {
		case string:
			*buf = enc.AppendString(*buf, val)
		case []byte:
			*buf = enc.AppendBytes(*buf, val)
		case bool:
			*buf = enc.AppendBool(*buf, val)
		case int:
			*buf = enc.AppendInt(*buf, val)
		case int8:
			*buf = enc.AppendInt8(*buf, val)
		case int16:
			*buf = enc.AppendInt16(*buf, val)
		case int32:
			*buf = enc.AppendInt32(*buf, val)
		case int64:
			*buf = enc.AppendInt64(*buf, val)
		case uint:
			*buf = enc.AppendUint(*buf, val)
		case uint8:
			*buf = enc.AppendUint8(*buf, val)
		case uint16:
			*buf = enc.AppendUint16(*buf, val)
		case uint32:
			*buf = enc.AppendUint32(*buf, val)
		case uint64:
			*buf = enc.AppendUint64(*buf, val)
		case float32:
			*buf = enc.AppendFloat32(*buf, val)
		case float64:
			*buf = enc.AppendFloat64(*buf, val)
		case time.Time:
			*buf = enc.AppendTime(*buf, val, time.RFC3339)
		case time.Duration:
			*buf = enc.AppendDuration(*buf, val, time.Millisecond, false)
		case *string:
			if val != nil {
				*buf = enc.AppendString(*buf, *val)
			} else {
				*buf = enc.AppendNil(*buf)
			}
		case *bool:
			if val != nil {
				*buf = enc.AppendBool(*buf, *val)
			} else {
				*buf = enc.AppendNil(*buf)
			}
		case *int:
			if val != nil {
				*buf = enc.AppendInt(*buf, *val)
			} else {
				*buf = enc.AppendNil(*buf)
			}
		case *int8:
			if val != nil {
				*buf = enc.AppendInt8(*buf, *val)
			} else {
				*buf = enc.AppendNil(*buf)
			}
		case *int16:
			if val != nil {
				*buf = enc.AppendInt16(*buf, *val)
			} else {
				*buf = enc.AppendNil(*buf)
			}
		case *int32:
			if val != nil {
				*buf = enc.AppendInt32(*buf, *val)
			} else {
				*buf = enc.AppendNil(*buf)
			}
		case *int64:
			if val != nil {
				*buf = enc.AppendInt64(*buf, *val)
			} else {
				*buf = enc.AppendNil(*buf)
			}
		case *uint:
			if val != nil {
				*buf = enc.AppendUint(*buf, *val)
			} else {
				*buf = enc.AppendNil(*buf)
			}
		case *uint8:
			if val != nil {
				*buf = enc.AppendUint8(*buf, *val)
			} else {
				*buf = enc.AppendNil(*buf)
			}
		case *uint16:
			if val != nil {
				*buf = enc.AppendUint16(*buf, *val)
			} else {
				*buf = enc.AppendNil(*buf)
			}
		case *uint32:
			if val != nil {
				*buf = enc.AppendUint32(*buf, *val)
			} else {
				*buf = enc.AppendNil(*buf)
			}
		case *uint64:
			if val != nil {
				*buf = enc.AppendUint64(*buf, *val)
			} else {
				*buf = enc.AppendNil(*buf)
			}
		case *float32:
			if val != nil {
				*buf = enc.AppendFloat32(*buf, *val)
			} else {
				*buf = enc.AppendNil(*buf)
			}
		case *float64:
			if val != nil {
				*buf = enc.AppendFloat64(*buf, *val)
			} else {
				*buf = enc.AppendNil(*buf)
			}
		case *time.Time:
			if val != nil {
				*buf = enc.AppendTime(*buf, *val, time.RFC3339)
			} else {
				*buf = enc.AppendNil(*buf)
			}
		case *time.Duration:
			if val != nil {
				*buf = enc.AppendDuration(*buf, *val, time.Millisecond, false)
			} else {
				*buf = enc.AppendNil(*buf)
			}
		case []string:
			*buf = enc.AppendStrings(*buf, val)
		case []bool:
			*buf = enc.AppendBools(*buf, val)
		case []int:
			*buf = enc.AppendInts(*buf, val)
		case []int8:
			*buf = enc.AppendInts8(*buf, val)
		case []int16:
			*buf = enc.AppendInts16(*buf, val)
		case []int32:
			*buf = enc.AppendInts32(*buf, val)
		case []int64:
			*buf = enc.AppendInts64(*buf, val)
		case []uint:
			*buf = enc.AppendUints(*buf, val)
		// case []uint8:
		// 	*buf = enc.AppendUints8(*buf, val)
		case []uint16:
			*buf = enc.AppendUints16(*buf, val)
		case []uint32:
			*buf = enc.AppendUints32(*buf, val)
		case []uint64:
			*buf = enc.AppendUints64(*buf, val)
		case []float32:
			*buf = enc.AppendFloats32(*buf, val)
		case []float64:
			*buf = enc.AppendFloats64(*buf, val)
		case []time.Time:
			*buf = enc.AppendTimes(*buf, val, time.RFC3339)
		case []time.Duration:
			*buf = enc.AppendDurations(*buf, val, time.Millisecond, false)
		case nil:
			*buf = enc.AppendNil(*buf)
		case net.IP:
			*buf = enc.AppendIPAddr(*buf, val)
		case net.IPNet:
			*buf = enc.AppendIPPrefix(*buf, val)
		case net.HardwareAddr:
			*buf = enc.AppendMACAddr(*buf, val)
		case stdJSON.RawMessage:
			*buf = appendJSON(*buf, val)
		default:
			*buf = enc.AppendInterface(*buf, val)
		}
	}

	*buf = enc.AppendEndMarker(*buf)
	*buf = enc.AppendLineBreak(*buf)

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.w.Write(*buf)

	return err
}

func appendJSON(dst []byte, j []byte) []byte {
	return append(dst, j...)
}
