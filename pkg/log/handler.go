package log

import (
	"context"
)

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
}
