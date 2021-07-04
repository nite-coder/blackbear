package web

import (
	"net/http/httptest"
)

func createTestContext() (*Context, *httptest.ResponseRecorder, *WebServer) {
	s := NewServer()
	w := httptest.NewRecorder()
	c := &Context{
		Writer: newResponseWriter(),
	}
	c.WebServer = s
	c.Writer.reset(w)
	return c, w, s
}
