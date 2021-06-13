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
	//c := newContext(nap, nil, w)
	return c, w, s
}
