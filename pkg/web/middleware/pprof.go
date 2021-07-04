package middleware

import (
	"net/http/pprof"

	"github.com/nite-coder/blackbear/pkg/web"
)

// PPROF is middleware struct
type PPROF struct {
}

// NewPPROF returns a mddleware instance
func NewPPROF() *PPROF {
	return &PPROF{}
}

// Invoke function is a middleware entry
func (p *PPROF) Invoke(c *web.Context, next web.HandlerFunc) {
	pprof.Index(c.Writer, c.Request)
	pprof.Cmdline(c.Writer, c.Request)
	pprof.Profile(c.Writer, c.Request)
	pprof.Symbol(c.Writer, c.Request)
	pprof.Trace(c.Writer, c.Request)
	_ = next(c)
}
