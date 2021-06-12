package http

import (
	"net/http/httptest"
)

func createTestContext() (*Context, *httptest.ResponseRecorder, *NapNap) {
	nap := New()
	w := httptest.NewRecorder()
	c := &Context{
		Writer: NewResponseWriter(),
	}
	c.NapNap = nap
	c.Writer.reset(w)
	//c := NewContext(nap, nil, w)
	return c, w, nap
}
