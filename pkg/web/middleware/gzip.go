package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/nite-coder/blackbear/pkg/web"
)

// These compression constants are copied from the compress/gzip package.
const (
	encodingGzip = "gzip"

	headerAcceptEncoding  = "Accept-Encoding"
	headerContentEncoding = "Content-Encoding"
	headerContentType     = "Content-Type"
	headerVary            = "Vary"
	headerSecWebSocketKey = "Sec-WebSocket-Key"

	BestCompression    = gzip.BestCompression
	BestSpeed          = gzip.BestSpeed
	DefaultCompression = gzip.DefaultCompression
	NoCompression      = gzip.NoCompression
)

// gzipResponseWriter is the ResponseWriter that http.ResponseWriter is
// wrapped in.
type gzipResponseWriter struct {
	gz        *gzip.Writer
	napWriter web.ResponseWriter
	web.ResponseWriter
}

// Write writes bytes to the gzip.Writer. It will also set the Content-Type
// header using the net/http library content type detection if the Content-Type
// header was not set yet.
func (grw gzipResponseWriter) Write(b []byte) (int, error) {
	if len(grw.Header().Get(headerContentType)) == 0 {
		grw.Header().Set(headerContentType, http.DetectContentType(b))
	}

	if len(grw.Header().Get(headerContentEncoding)) > 0 {
		// compress the content
		return grw.gz.Write(b)
	}
	// no compress
	grw.gz.Reset(io.Discard)
	return grw.napWriter.Write(b)
}

// GzipMiddleware struct is gzip middlware
type GzipMiddleware struct {
	pool sync.Pool
}

// NewGzip returns a middleware which will handle the Gzip compression in Invoke.
// Valid values for level are identical to those in the compress/gzip package.
func NewGzip(level int) *GzipMiddleware {
	h := &GzipMiddleware{}

	h.pool.New = func() interface{} {
		gz, err := gzip.NewWriterLevel(io.Discard, level)
		if err != nil {
			panic(err)
		}

		return gz
	}

	return h
}

// Invoke function is a middleware entry
func (h *GzipMiddleware) Invoke(c *web.Context, next web.HandlerFunc) {
	r := c.Request
	w := c.Writer
	// Skip compression if the client doesn't accept gzip encoding.
	if !strings.Contains(r.Header.Get(headerAcceptEncoding), encodingGzip) {
		_ = next(c)
		return
	}

	// Skip compression if client attempt WebSocket connection
	if len(r.Header.Get(headerSecWebSocketKey)) > 0 {
		_ = next(c)
		return
	}

	// Skip compression if already compressed
	if w.Header().Get(headerContentEncoding) == encodingGzip {
		_ = next(c)
		return
	}

	// Retrieve gzip writer from the pool. Reset it to use the ResponseWriter.
	// This allows us to re-use an already allocated buffer rather than
	// allocating a new buffer for every request.
	// We defer g.pool.Put here so that the gz writer is returned to the
	// pool if any thing after here fails for some reason (functions in
	// next could potentially panic, etc)
	gz, _ := h.pool.Get().(*gzip.Writer)
	defer h.pool.Put(gz)
	gz.Reset(w)

	// Set the appropriate gzip headers.
	headers := w.Header()
	headers.Set(headerContentEncoding, encodingGzip)
	headers.Set(headerVary, headerAcceptEncoding)

	// Wrap the original http.ResponseWriter
	// and create the gzipResponseWriter.
	grw := gzipResponseWriter{
		gz,
		w,
		w,
	}

	// Call the next handler supplying the gzipResponseWriter instead of
	// the original.
	c.Writer = grw
	_ = next(c)

	_ = gz.Close()
}
