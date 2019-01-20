RxRouter is a simple and fast HTTP router for Go. RxRouter is a marriage of the Bxog request multiplexer (one of the fastest muxes) to fasthttp server.
Credit goes to [Valyala] (https://github.com/valyala) for fasthttp and [Claygod] (https://github.com/claygod) for Bxog

[![API documentation](https://godoc.org/github.com/rohanthewiz/rxrouter?status.svg)](https://godoc.org/github.com/rohanthewiz/rxrouter)

## Warning: Currently this is totally a POC -- do not use in production!

## Usage
Please see this example: https://github.com/rohanthewiz/rxrun

# Settings
Settings are passed into the rxrouter.New() function

# Performance
Since we are based on perhaps the fastest http package (fasthttp) and one of the fastest route multiplexers available (Bxog),
we should be hitting some high marks. A benchmark is on the todo.


Example:

```go
package main

import (
	"fmt"
	"github.com/rohanthewiz/rxrouter"
	"github.com/valyala/fasthttp"
	"log"
)

func main() {
	rx := rxrouter.New(
		rxrouter.Options{Verbose: false, Port: "3026"}, // the Options argument here is optional
	)

	// Logging middleware
	rx.Use(
		func(ctx *fasthttp.RequestCtx) (ok bool) {
			log.Printf("Requested path: %s", ctx.Path())
			return true
		},
		fasthttp.StatusServiceUnavailable, // 503
	)

	// Auth middleware
	rx.Use(
		func(ctx *fasthttp.RequestCtx) (ok bool) {
			authed := true // pretend we got a good response from our auth check
			if !authed {
				return false
			}
			return true
		},
		fasthttp.StatusUnauthorized,
	)

	// Add some routes
	rx.AddRoute("/", func(ctx *fasthttp.RequestCtx, params map[string]string) {
		fmt.Fprintf(ctx, "Hello, world! Requested path is %q", ctx.Path())
	})
	rx.AddRoute("/hello/:name/:age", handleHello)
	rx.AddRoute("/store/:number/:location", handleStore)
	// Routes for static files
	rx.AddStaticFilesRoute("/images/", "./assets/images", 1)
	rx.AddStaticFilesRoute("/css/", "./assets/css", 1)

	// Let it rip!
	rx.Start()
}

func handleHello(ctx *fasthttp.RequestCtx, params map[string]string) {
	ctx.WriteString(fmt.Sprintf("Hello %s", params["name"]))
}

func handleStore(ctx *fasthttp.RequestCtx, params map[string]string) {
	fmt.Fprintf(ctx, "Store: %s  location: %s", params["number"], params["location"])
}
```

# Named parameters

Arguments in the rules designated route colon. Example route: */abc/:param* , where *abc* is a static section and *:param* - the dynamic section(argument).

# Static files
Use the AddStaticFilesRoute method. Example: `rx.AddStaticFilesRoute("/images/", "./assets/images", 1)`

# API

Methods:
-  *New* - create a new router passing in options (including port)
-  *AddRoute* - add a rule specifying the handler
-  *AddStaticFilesRoute* - add a route for serving static files
-  *Start* - start the server
