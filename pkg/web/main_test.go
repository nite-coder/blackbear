package web

import (
	"net/http/httptest"
)

func createTestContext() (*Context, *httptest.ResponseRecorder, *WebServer) {
	s := NewServer()
	w := httptest.NewRecorder()
	c := &Context{
		Writer: NewResponseWriter(),
	}
	c.WebServer = s
	c.Writer.reset(w)
	//c := NewContext(nap, nil, w)
	return c, w, s
}
