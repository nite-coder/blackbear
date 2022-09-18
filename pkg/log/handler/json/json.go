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
	//Out []byte
}

// New handler.
func New() *Handler {
	return &Handler{
		//Out: make([]byte, 500),
	}
}

// BeforeWriting implements log.Handler.
func (h *Handler) BeforeWriting(e *log.Entry) error {
	e.Str("level", e.Level.String())

	return nil
}

// Write implements log.Handler.
func (h *Handler) Write(bytes []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	fmt.Fprintf(os.Stdout, "%s\n", bytes)
	//h.Out = bytes
	return nil
}

// Flush clear all buffer
func (h *Handler) Flush() error {
	return nil
}
