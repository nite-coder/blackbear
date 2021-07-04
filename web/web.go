package web

import (
	"html/template"
	"net/http"
	"path"
	"sync"
	"time"
)

var (
	_logger *logger
)

func init() {
	_logger = &logger{
		mode: off,
	}
}

// HandlerFunc defines a function to server HTTP requests
type HandlerFunc func(c *Context) error

// ErrorHandler defines a function to handle HTTP errors
type ErrorHandler func(c *Context, err error)

// MiddlewareHandler is an interface that objects can implement to be registered to serve as middleware
// in the WebServer middleware stack.
type MiddlewareHandler interface {
	Invoke(c *Context, next HandlerFunc)
}

// MiddlewareFunc is an adapter to allow the use of ordinary functions as WebServer handlers.
type MiddlewareFunc func(c *Context, next HandlerFunc)

// Invoke function is a middleware entry
func (m MiddlewareFunc) Invoke(c *Context, next HandlerFunc) {
	m(c, next)
}

type middleware struct {
	handler MiddlewareHandler
	next    *middleware
}

func (m middleware) Execute(c *Context) error {
	m.handler.Invoke(c, m.next.Execute)
	return nil
}

// WrapHandler wraps `http.Handler` into `web.HandlerFunc`.
func WrapHandler(h http.Handler) HandlerFunc {
	return func(c *Context) error {
		h.ServeHTTP(c.Writer, c.Request)
		return nil
	}
}

// WebServer is root level of framework instance
type WebServer struct {
	pool             sync.Pool
	handlers         []MiddlewareHandler
	middleware       middleware
	template         *template.Template
	templateRootPath string
	router           *router

	MaxRequestBodySize int64
	ErrorHandler       ErrorHandler
	NotFoundHandler    HandlerFunc
}

// NewServer returns a new WebServer instance
func NewServer(mHandlers ...MiddlewareHandler) *WebServer {
	s := &WebServer{
		handlers:           mHandlers,
		middleware:         build(mHandlers),
		MaxRequestBodySize: 10485760, // default 10MB for request body size
	}

	s.pool.New = func() interface{} {
		rw := newResponseWriter()
		return newContext(s, nil, rw)
	}

	s.router = newRouter(s)
	s.Use(s.router)

	s.NotFoundHandler = func(c *Context) error {
		c.SetStatus(404)
		return nil
	}

	return s
}

// UseFunc adds an anonymous function onto middleware stack.
func (s *WebServer) UseFunc(aFunc func(c *Context, next HandlerFunc)) {
	s.Use(MiddlewareFunc(aFunc))
}

// Use adds a Handler onto the middleware stack. Handlers are invoked in the order they are added to a WebServer.
func (s *WebServer) Use(mHandler MiddlewareHandler) {
	if len(s.handlers) == 0 {
		s.handlers = append(s.handlers, mHandler)
	} else {
		end := len(s.handlers) - 1
		s.handlers = append(s.handlers[:end], mHandler, s.router)
	}

	s.middleware = build(s.handlers)
}

// All is a shortcut for adding all methods
func (s *WebServer) All(path string, handler HandlerFunc) {
	s.router.Add(GET, path, handler)
	s.router.Add(POST, path, handler)
	s.router.Add(PUT, path, handler)
	s.router.Add(DELETE, path, handler)
	s.router.Add(PATCH, path, handler)
	s.router.Add(OPTIONS, path, handler)
	s.router.Add(HEAD, path, handler)
}

// Get is a shortcut for router.Add("GET", path, handle)
func (s *WebServer) Get(path string, handler HandlerFunc) {
	s.router.Add(GET, path, handler)
}

// Post is a shortcut for router.Add("POST", path, handle)
func (s *WebServer) Post(path string, handler HandlerFunc) {
	s.router.Add(POST, path, handler)
}

// Put is a shortcut for router.Add("PUT", path, handle)
func (s *WebServer) Put(path string, handler HandlerFunc) {
	s.router.Add(PUT, path, handler)
}

// Delete is a shortcut for router.Add("DELETE", path, handle)
func (s *WebServer) Delete(path string, handler HandlerFunc) {
	s.router.Add(DELETE, path, handler)
}

// Patch is a shortcut for router.Add("PATCH", path, handle)
func (s *WebServer) Patch(path string, handler HandlerFunc) {
	s.router.Add(PATCH, path, handler)
}

// Options is a shortcut for router.Add("OPTIONS", path, handle)
func (s *WebServer) Options(path string, handler HandlerFunc) {
	s.router.Add(OPTIONS, path, handler)
}

// Head is a shortcut for router.Add("HEAD", path, handle)
func (s *WebServer) Head(path string, handler HandlerFunc) {
	s.router.Add(HEAD, path, handler)
}

// SetTemplate function allows user to set their own template instance.
func (s *WebServer) SetTemplate(t *template.Template) {
	s.template = t
}

// SetRender function allows user to set template location.
func (s *WebServer) SetRender(templateRootPath string) {
	sharedTemplatePath := path.Join(templateRootPath, "shares/*")
	tmpl, err := template.ParseGlob(sharedTemplatePath)
	template := template.Must(tmpl, err)

	if template == nil {
		_logger.debug("no template")
		template = template.New("")
	}

	s.template = template
	s.templateRootPath = templateRootPath
}

// Run will start to run a http server
// TODO: allow multiple ports and addrs
func (s *WebServer) Run(addr string) error {
	serv := new(http.Server)
	serv.Addr = addr
	serv.Handler = s
	return serv.ListenAndServe()
}

// RunAll will listen on multiple port
// func (s *WebServer) RunAll(addrs []string) error {
// 	if len(addrs) == 0 {
// 		return errors.New("addrs can't be empty")
// 	}

// 	wg := &sync.WaitGroup{}

// 	for _, addr := range addrs {
// 		wg.Add(1)
// 		go func(newAddr string) {
// 			err := http.ListenAndServe(newAddr, s)
// 			if err != nil {
// 				panic(err)
// 			}
// 			wg.Done()
// 		}(addr)
// 	}

// 	wg.Wait()
// 	return nil
// }

type ServerOptions struct {
	Addr          string
	Domain        string // abc123.com, abc456.com
	CertCachePath string
	TLSCertFile   string
	TLSKeyFile    string
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
}

// RunTLS will run http/2 server
func (s *WebServer) RunTLS(addr, cert, key string) error {
	serv := new(http.Server)
	serv.Addr = addr
	serv.Handler = s
	return serv.ListenAndServeTLS(cert, key)
}

// Conforms to the http.Handler interface.
func (s *WebServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	req.Body = http.MaxBytesReader(w, req.Body, s.MaxRequestBodySize)
	c := s.pool.Get().(*Context)
	c.reset(w, req)
	_ = s.middleware.Execute(c)
	s.pool.Put(c)
}

func build(handlers []MiddlewareHandler) middleware {
	var next middleware

	if len(handlers) == 0 {
		return voidMiddleware()
	} else if len(handlers) > 1 {
		next = build(handlers[1:])
	} else {
		next = voidMiddleware()
	}

	return middleware{handlers[0], &next}
}

func voidMiddleware() middleware {
	return middleware{
		MiddlewareFunc(func(c *Context, next HandlerFunc) {}),
		&middleware{},
	}
}
