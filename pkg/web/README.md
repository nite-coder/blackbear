# Web



## Examples

### hello world
```go
package main

import (
	"github.com/nite-coder/blackbear/pkg/web"
)

func main() {
	s := web.NewServer()

	s.Get("/", func(c *web.Context) error {
		return c.String(200, "Hello World")
	})

	s.Run(":10080")
}
```

### Using GET, POST, PUT, PATCH, DELETE and OPTIONS
```go
package main

import (
	"github.com/nite-coder/blackbear/pkg/web"
)

func main() {
	s := web.NewServer()

	s.Get("/my-get", myHeadEndpoint)
	s.Post("/my-post", myPostEndpoint)
	s.Put("/my-put", myPutEndpoint)
	s.Delete("/my-delete", myDeleteEndpoint)
	s.Patch("/my-patch", myPatchEndpoint)
	s.Options("/my-options", myOptionsEndpoint)
	s.Head("/my-head", myHeadEndpoint)

    s.Run(":10080")
}
```

### Parameters in path

```go
package main

import (
	"github.com/nite-coder/blackbear/pkg/web"
)

func main() {
	s := web.NewServer()

	s.Get("/users/:name", func(c *web.Context) error {
		name := c.Param("name")
		return c.String(200, "Hello, "+name)
	})

	// /videos/sports/basketball/1.mp4
	// /videos/2.mp4
	// both path will route to the endpoint
	s.Get("/videos/*video_id", func(c *web.Context) error {
		id := c.Param("video_id")
		return c.String(200, "video id is, "+id)
	})

	s.Run(":10080")
}
```

### Get querystring value
```go
package main

import (
	"github.com/nite-coder/blackbear/pkg/web"
)

func main() {
	s := web.NewServer()

	// http://localhost:10080/test?page=3
	s.Get("/test", func(c *web.Context) error {
		page := c.Query("page") //get query string value
		return c.String(200, page)
	})

	err := s.Run(":10080")
	if err != nil {
		panic(err)
	}
}
```

### Get post form value (Multipart/Urlencoded Form)
```go
package main

import (
	"github.com/jasonsoft/napnap"
)

func main() {
	nap := napnap.New()

	nap.Post("/post-form-value", func(c *napnap.Context) error {
		userId := c.Form("user_id") //get post form value
		return c.String(200, userId)
	})

	http.ListenAndServe("127.0.0.1:10080", nap)
}
```

#### JSON binding

```go
package main

import "github.com/jasonsoft/napnap"

func main() {
	nap := napnap.New()

	nap.Post("/json-binding", func(c *napnap.Context) error {
		var person struct {
			Name string `json: name`
			Age  int    `json: age`
		}
        err := c.BindJSON(&person)
        if err != nil {
            return err
        }
		c.String(200, person.Name)
        return nil
	})

	http.ListenAndServe("127.0.0.1:10080", nap)
}
```

#### JSON rendering

```go
package main

import "github.com/jasonsoft/napnap"

func main() {
	nap := napnap.New()

	nap.Get("/json-rendering", func(c *napnap.Context) error {
		var person struct {
			Name string `json: name`
			Age  int    `json: age`
		}

		person.Name = "napnap"
		person.Age = 18

		return c.JSON(200, person)
	})

	http.ListenAndServe("127.0.0.1:10080", nap) 
}
```

#### Http/2 Server

```go
package main

import "github.com/jasonsoft/napnap"

func main() {
	router := napnap.NewRouter()

	router.Get("/hello-world", func(c *napnap.Context) {
		c.String(200, "Hello, World")
	})

	nap := napnap.New()
	nap.Use(router)
	nap.RunTLS(":443", "cert.crt", "key.pem") // nap will use http/2 server as default
}
```

#### combine with autocert

Let's Encrypt disable **tls challenge**, so we can only use **http challenge** with autocert(**dns challenge** not implemented)

notes:

we need to bind http service on 80 port, https service on 443 port. First time, you need to wait a short time for creating and downloading certificate .

```go
package main

import "github.com/jasonsoft/napnap"

func main() {
	router := napnap.NewRouter()

	projRoot, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err.Error())
	}

	config := napnap.Config{
		Domain:        "exmaple.com", // multi domain support ex. "abc.com, 123.com"
		CertCachePath: path.Join(projRoot, ".certCache"),
	}
	server := napnap.NewHttpEngineWithConfig(&config)

	nap := napnap.New()
	nap.Use(router)

	nap.RunAutoTLS(server)
}
```

