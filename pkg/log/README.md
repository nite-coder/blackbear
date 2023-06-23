# log 

It is a simple structured logging package for Go. 
## Features
* fast, easy to use, and pretty logging for development
* low to zero allocation
* JSON encoding format
* colored text for text handler
* `context.Context` integration

## Handlers
* Text (development use)
* JSON (default)

## Installation
Use go get 

```go
go get -u github.com/nite-coder/blackbear
```

## Get Started

```go
package main

import (
	"os"

	"github.com/nite-coder/blackbear/pkg/log"
	"github.com/nite-coder/blackbear/pkg/log/handler/text"
)

func main() {
	// json handler
	log.Debug().Msg("Hello World") // {"time":"2023-06-23T06:17:43Z","level":"DEBUG","msg":"Hello World"}

	// text handler
	opts := log.HandlerOptions{
		Level:       log.DebugLevel,
		DisableTime: true,
	}
	logger := log.New(text.New(os.Stderr, &opts))
	log.SetDefault(logger)
	log.Debug().Msg("Hello World") // 06:17:43.991 DEBUG  Hello World
}
```
### Fields
```go
package main

import (
	"github.com/nite-coder/blackbear/pkg/log"
)

func main() {
    // example1
	logger := log.With().Str("app_id", "blackbear").Logger()
	logger.Debug().Msg("Hello World")

    // example2
    log.Debug().Str("request_id", "abc").Msg("cool")
    
}
```
### Pass Context
```go
package main

import (
	"github.com/nite-coder/blackbear/pkg/log"
)

func main() {
    ctx := context.Background()
    log.DebugCtx(ctx).Str("request_id", "abc").Msg("cool")
}
```