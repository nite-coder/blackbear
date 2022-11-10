package json

import (
	"fmt"
	"os"
	"sync"

	"github.com/nite-coder/blackbear/pkg/log"
)

// Handler implementation.
type Handler struct {
	mu sync.Mutex
}

// New handler.
func New() *Handler {
	return &Handler{}
}

// BeforeWriting implements log.Handler.
func (h *Handler) BeforeWriting(e *log.Entry) error {
	e.Str("level", e.Level.String())

	if !e.Logger.DisableTimeField {
		e.Str("time", e.CreatedAt.Format("2006-01-02 15:04:05.000Z"))
	}

	return nil
}

// Write implements log.Handler.
func (h *Handler) Write(bytes []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	fmt.Fprintf(os.Stdout, "%s", bytes)
	return nil
}

// Flush clear all buffer
func (h *Handler) Flush() error {
	return nil
}
