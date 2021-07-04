package web

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testRoute(t *testing.T, method string, path string) {
	passed := false
	_, w, s := createTestContext()

	router := newRouter(s)
	router.Add(method, path, func(c *Context) error {
		passed = true
		c.SetStatus(200)
		return nil
	})
	s.Use(router)

	req, _ := http.NewRequest(method, path, nil)
	s.ServeHTTP(w, req)

	assert.True(t, passed)
	assert.Equal(t, 200, w.Code)
}

func TestRouterStaticRoute(t *testing.T) {
	testRoute(t, "GET", "/")
	testRoute(t, "GET", "/hello")
	testRoute(t, "POST", "/hello")
	testRoute(t, "PUT", "/hello/put")
	testRoute(t, "DELETE", "/hello/Delet")
}

func TestRouterParameterRoute(t *testing.T) {
	var name, age string
	_, w, s := createTestContext()

	router := newRouter(s)
	router.Add(GET, "/users/:name", func(c *Context) error {
		name = c.Param("name")
		age = c.Query("age")
		c.SetStatus(200)
		return nil
	})
	s.Use(router)

	req, _ := http.NewRequest("GET", "/users/john?age=18", nil)
	s.ServeHTTP(w, req)

	assert.Equal(t, "john", name)
	assert.Equal(t, "18", age)
	assert.Equal(t, 200, w.Code)
}

func TestRouterMatchAnyRoute(t *testing.T) {
	var action, helo string
	_, w, s := createTestContext()

	router := newRouter(s)
	router.Add(GET, "/video/:action1", func(c *Context) error {
		action = c.Param("action1")
		c.SetStatus(201)
		return nil
	})

	router.Add(GET, "/images/*action2", func(c *Context) error {
		action = c.Param("action2")
		c.SetStatus(200)
		return nil
	})

	router.Add(GET, "/v1/:helo/images/*action2", func(c *Context) error {
		helo = c.Param("helo")
		action = c.Param("action2")
		c.SetStatus(200)

		return nil
	})

	s.Use(router)

	req, _ := http.NewRequest("GET", "/images/play/ball.jpg", nil)
	s.ServeHTTP(w, req)
	assert.Equal(t, "play/ball.jpg", action)
	assert.Equal(t, 200, w.Code)

	req, _ = http.NewRequest("GET", "/v1/aabbc/images/hapyy-ball.jpg", nil)
	s.ServeHTTP(w, req)
	assert.Equal(t, "hapyy-ball.jpg", action)
	assert.Equal(t, "aabbc", helo)
	assert.Equal(t, 200, w.Code)
}
