package middleware

import (
	"strings"

	bearHTTP "github.com/nite-coder/blackbear/pkg/http"
)

// Health is health middleware struct
type Health struct {
}

// NewHealth returns Health middlware instance
func NewHealth() *Health {
	return &Health{}
}

// Invoke function is a middleware entry
func (h *Health) Invoke(c *bearHTTP.Context, next bearHTTP.HandlerFunc) {
	if strings.EqualFold(c.Request.URL.Path, "/health") {
		_ = c.String(200, "OK")
	} else {
		_ = next(c)
	}
}
