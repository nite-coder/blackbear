package web

import "strings"

type tree struct {
	rootNode *node
}

type kind uint8

// example: /user/:jason/count
// params: ["user", "jason", count]
// pNames: current node ["jason"]

type node struct {
	parent    *node
	children  []*node
	kind      kind
	name      string
	pNames    []string
	params    []string
	sortOrder int
	handler   *methodHandler
}

type methodHandler struct {
	connect HandlerFunc
	delete  HandlerFunc
	get     HandlerFunc
	head    HandlerFunc
	options HandlerFunc
	patch   HandlerFunc
	post    HandlerFunc
	put     HandlerFunc
	trace   HandlerFunc
}

const (
	// CONNECT HTTP method
	CONNECT = "CONNECT"
	// DELETE HTTP method
	DELETE = "DELETE"
	// GET HTTP method
	GET = "GET"
	// HEAD HTTP method
	HEAD = "HEAD"
	// OPTIONS HTTP method
	OPTIONS = "OPTIONS"
	// PATCH HTTP method
	PATCH = "PATCH"
	// POST HTTP method
	POST = "POST"
	// PUT HTTP method
	PUT = "PUT"
	// TRACE HTTP method
	TRACE = "TRACE"
)

const (
	skind kind = iota
	pkind
	akind
)

type router struct {
	webServer *WebServer
	tree      *tree
}

// NewRouter function will create a new router instance
func newRouter(s *WebServer) *router {
	return &router{
		webServer: s,
		tree: &tree{
			rootNode: &node{
				parent:    nil,
				children:  []*node{},
				kind:      0,
				name:      "/",
				sortOrder: 0,
				handler:   &methodHandler{},
			},
		},
	}
}

// Invoke function is a middleware entry
func (r *router) Invoke(c *Context, next HandlerFunc) {
	h := r.Find(c.Request.Method, c.Request.URL.Path, c)

	var err error

	if h == nil {
		if r.webServer.NotFoundHandler != nil {
			err = r.webServer.NotFoundHandler(c)
		}
	} else {
		err = h(c)
	}

	if err != nil && r.webServer.ErrorHandler != nil {
		r.webServer.ErrorHandler(c, err)
	}
}

// All is a shortcut for adding all methods
func (r *router) All(path string, handler HandlerFunc) {
	r.Add(GET, path, handler)
	r.Add(POST, path, handler)
	r.Add(PUT, path, handler)
	r.Add(DELETE, path, handler)
	r.Add(PATCH, path, handler)
	r.Add(OPTIONS, path, handler)
	r.Add(HEAD, path, handler)
}

// Get is a shortcut for router.Add("GET", path, handle)
func (r *router) Get(path string, handler HandlerFunc) {
	r.Add(GET, path, handler)
}

// Post is a shortcut for router.Add("POST", path, handle)
func (r *router) Post(path string, handler HandlerFunc) {
	r.Add(POST, path, handler)
}

// Put is a shortcut for router.Add("PUT", path, handle)
func (r *router) Put(path string, handler HandlerFunc) {
	r.Add(PUT, path, handler)
}

// Delete is a shortcut for router.Add("DELETE", path, handle)
func (r *router) Delete(path string, handler HandlerFunc) {
	r.Add(DELETE, path, handler)
}

// Patch is a shortcut for router.Add("PATCH", path, handle)
func (r *router) Patch(path string, handler HandlerFunc) {
	r.Add(PATCH, path, handler)
}

// Options is a shortcut for router.Add("OPTIONS", path, handle)
func (r *router) Options(path string, handler HandlerFunc) {
	r.Add(OPTIONS, path, handler)
}

// Head is a shortcut for router.Add("HEAD", path, handle)
func (r *router) Head(path string, handler HandlerFunc) {
	r.Add(HEAD, path, handler)
}

// Add function which adding path and handler to router
func (r *router) Add(method string, path string, handler HandlerFunc) {
	_logger.debug("===Add")

	if len(path) == 0 {
		panic("router: path couldn't be empty")
	}

	if path[0] != '/' {
		panic("router: path was invalid")
	}

	if len(path) > 1 {
		path = path[1:]
	}

	_logger.debug("path:" + path)

	currentNode := r.tree.rootNode
	if path == "/" {
		currentNode.addHandler(method, handler)
		return
	}

	pathArray := strings.Split(path, "/")
	count := len(pathArray)
	pathParams := []string{}

	for index, element := range pathArray {
		if len(element) == 0 {
			continue
		}

		var childNode *node

		firstSymbol := element[0]

		switch firstSymbol {
		case ':':
			// this is parameter node
			pName := element[1:]
			_logger.debug("parameter_node_pname:" + pName)
			childNode = currentNode.findChildByKind(pkind)

			if childNode == nil {
				childNode = newNode(pName, pkind)
				currentNode.addChild(childNode)
			}

			isFound := false

			for _, p := range childNode.pNames {
				if p == pName {
					isFound = true
				}
			}

			if !isFound {
				childNode.pNames = append(childNode.pNames, pName)
				_logger.debug("add_parameter_name:" + pName)
			}

			pathParams = append(pathParams, pName)
		case '*':
			// this is match any node.  We should allow one match any node only.
			pName := element[1:]
			_logger.debug("match_node_pname:" + pName)
			childNode = currentNode.findChildByKind(akind)
			if childNode == nil {
				childNode = newNode(pName, akind)
				currentNode.addChild(childNode)
				childNode.pNames = append(childNode.pNames, pName)
			}

			pathParams = append(pathParams, pName)
		default:
			// this is static node
			childNode = currentNode.findChildByName(element)
			if childNode == nil {
				childNode = newNode(element, skind)
				currentNode.addChild(childNode)
			}
		}

		// last node in the path
		if count == index+1 {
			childNode.params = pathParams
			childNode.addHandler(method, handler)
		}

		currentNode = childNode
	}
}

// Find returns http handler for specific path
func (r *router) Find(method string, path string, c *Context) HandlerFunc {
	_logger.debug("===Find")
	_logger.debug("method:" + method)
	_logger.debug("path:" + path)

	path = sanitizeUrl(path)

	currentNode := r.tree.rootNode
	if path == "/" {
		return currentNode.findHandler(method)
	}

	pathArray := strings.Split(path, "/")
	count := len(pathArray)

	pathParams := make(map[int][]Param)

	var paramsNum int

	for index, element := range pathArray {
		if len(element) == 0 {
			continue
		}

		// find static node first
		childNode := currentNode.findChildByName(element)

		if childNode == nil {
			// looking for parameter node
			childNode = currentNode.findChildByKind(pkind)

			if childNode != nil {
				_logger.debugf("parameter node: %s", element)

				var newParams []Param

				for _, pName := range childNode.pNames {
					param := Param{Key: pName, Value: element}
					newParams = append(newParams, param)
				}

				pathParams[paramsNum] = newParams
				paramsNum++
			}
		}

		if childNode == nil {
			// looking for match any node
			childNode = currentNode.findChildByKind(akind)

			if childNode != nil {
				_logger.debugf("match node: %s", element)
				start := 0

				for i := 0; i < index; i++ {
					start += 1 + len(pathArray[i])
				}

				_logger.debugf("start: %d", start)
				_logger.debugf("pname count: %d", len(childNode.pNames))

				var newParams []Param

				for _, pName := range childNode.pNames {
					val := path[start:]
					_logger.debugf("val: %s", val)
					param := Param{Key: pName, Value: val}
					newParams = append(newParams, param)
				}

				pathParams[paramsNum] = newParams
				paramsNum++

				index = count - 1
			}
		}

		if childNode == nil {
			return nil
		}

		// last node in the path
		if count == index+1 {
			myHandler := childNode.findHandler(method)
			if myHandler == nil {
				_logger.debug("handler was not found")
				return nil
			}

			paramsNum = 0
			// println("params_count:", len(pathParams))
			_logger.debug("lastNode_params_count:", len(childNode.params))

			for _, validParam := range childNode.params {
				for _, p := range pathParams[paramsNum] {
					if validParam == p.Key {
						_logger.debug("matched: " + validParam + "," + p.Value)
						c.params = append(c.params, p)
					}
				}
				paramsNum++
			}

			return myHandler
		}

		currentNode = childNode
	}
	return nil
}

func newNode(name string, t kind) *node {
	return &node{
		kind:      t,
		name:      name,
		sortOrder: 0,
		handler:   &methodHandler{},
	}
}

func (n *node) addChild(node *node) {
	node.parent = n
	n.children = append(n.children, node)
}

func (n *node) findChildByName(name string) *node {
	var result *node

	for _, element := range n.children {
		if strings.EqualFold(element.name, name) && element.kind == skind {
			result = element
			break
		}
	}

	return result
}

func (n *node) findChildByKind(t kind) *node {
	for _, c := range n.children {
		if c.kind == t {
			return c
		}
	}

	return nil
}

func (n *node) addHandler(method string, h HandlerFunc) {
	switch method {
	case GET:
		n.handler.get = h
	case POST:
		n.handler.post = h
	case PUT:
		n.handler.put = h
	case DELETE:
		n.handler.delete = h
	case PATCH:
		n.handler.patch = h
	case OPTIONS:
		n.handler.options = h
	case HEAD:
		n.handler.head = h
	case CONNECT:
		n.handler.connect = h
	case TRACE:
		n.handler.trace = h
	default:
		panic("method was invalid")
	}
}

func (n *node) findHandler(method string) HandlerFunc {
	switch method {
	case GET:
		return n.handler.get
	case POST:
		return n.handler.post
	case PUT:
		return n.handler.put
	case DELETE:
		return n.handler.delete
	case PATCH:
		return n.handler.patch
	case OPTIONS:
		return n.handler.options
	case HEAD:
		return n.handler.head
	case CONNECT:
		return n.handler.connect
	case TRACE:
		return n.handler.trace
	default:
		panic("method was invalid")
	}
}

func sanitizeUrl(redir string) string {
	if len(redir) > 1 && redir[0] == '/' && redir[1] != '/' && redir[1] != '\\' {
		return redir
	}
	return "/"
}
