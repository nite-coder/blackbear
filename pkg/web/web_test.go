package web

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultHandlers(t *testing.T) {
	_, w, s := createTestContext()

	m1 := false
	s.UseFunc(func(c *Context, next HandlerFunc) {
		m1 = true
		next(c)
	})

	isError := false
	s.ErrorHandler = func(c *Context, err error) {
		isError = true
		assert.Equal(t, "oops", err.Error())
	}

	isNotFound := false
	s.NotFoundHandler = func(c *Context) error {
		isNotFound = true
		return nil
	}

	s.Get("/error", func(c *Context) error {
		return errors.New("oops")
	})

	req, _ := http.NewRequest("GET", "/error", nil)
	s.ServeHTTP(w, req)

	req, _ = http.NewRequest("GET", "/not_found", nil)
	s.ServeHTTP(w, req)

	assert.Equal(t, true, m1)
	assert.Equal(t, true, isError)
	assert.Equal(t, true, isNotFound)
}

func TestMidderwareOrder(t *testing.T) {
	_, w, s := createTestContext()

	m1 := false
	s.UseFunc(func(c *Context, next HandlerFunc) {
		m1 = true
		next(c)
	})

	m2 := false
	s.UseFunc(func(c *Context, next HandlerFunc) {
		if m1 && m2 == false {
			m2 = true
		}
		next(c)
	})

	m3 := false
	s.Get("/hello", func(c *Context) error {
		if m2 && m3 == false {
			m3 = true
		}
		return nil
	})

	req, _ := http.NewRequest("GET", "/hello", nil)
	s.ServeHTTP(w, req)

	assert.Equal(t, true, m1)
	assert.Equal(t, true, m2)
	assert.Equal(t, true, m3)
}
