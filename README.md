RxRouter is a simple and fast HTTP router for Go. RxRouter is a marriage of the Bxog request multiplexer (one of the fastest muxes) to fasthttp server.
All credit goes to [Valyala] (https://github.com/valyala) and [Claygod] (https://github.com/claygod)

[![API documentation](https://godoc.org/github.com/claygod/Bxog?status.svg)](https://godoc.org/github.com/claygod/Bxog)
[![Go Report Card](https://goreportcard.com/badge/github.com/claygod/Bxog)](https://goreportcard.com/report/github.com/claygod/Bxog)

## Warning: Currently this is totally a POC -- do not use in production!

## Usage

An example of using the multiplexer:
See: https://github.com/rohanthewiz/rxrun

# Settings

Necessary changes in the configuration of the multiplexer can be made in the configuration file

# Perfomance

Bxog is the fastest mux, showing the speed of query processing. The benchmark results show how it compares with popular multiplexers: Bone, Httprouter, Gorilla, Zeus. The test is done on a computer with a i3-6320 3.7GHz processor and 8 GB RAM. In short (less time, the better):

- Bxog 163 ns/op
- HttpRouter 183 ns/op
- Zeus 12302 ns/op
- GorillaMux 14928 ns/op
- GorillaPat 618 ns/op
- Bone 47333 ns/op

Detailed benchmark [here](https://github.com/claygod/BxogTest)

# API

Methods:
-  *New* - create a new multiplexer
-  *Add* - add a rule specifying the handler (the default method - GET, ID - as a string to this rule)
-  *Start* - start the server indicating the listening port
-  *Params* - extract parameters from URL
-  *Create* - generate URL of the available options
-  *Shutdown* - graceful stop the server
-  *Stop* - aggressive stop the server
-  *Test* - Start analogue (for testing only)

Example:

```go
package main

import (
	"fmt"
	"github.com/rohanthewiz/rxrouter"
	"github.com/rohanthewiz/rxrouter/mux"
	"github.com/valyala/fasthttp"
	"log"
)

const appEnv = "[DEV]"
const authed = true

func main() {
	rx := rxrouter.New()

	// Rudimentary request logging middleware
	rx.Use(func(ctx *fasthttp.RequestCtx) (retCtx *fasthttp.RequestCtx, ok bool) {
		log.Printf("Requested path: %s", ctx.Path())
		return ctx, true
	}, fasthttp.StatusServiceUnavailable) // 503

	// Auth middleware
	rx.Use(func(ctx *fasthttp.RequestCtx) (retCtx *fasthttp.RequestCtx, ok bool) {
		if !authed { return ctx, false }
		return ctx, true
	}, fasthttp.StatusUnauthorized)

	// Prepend to output middleware
	rx.Use(func(ctx *fasthttp.RequestCtx) (retCtx *fasthttp.RequestCtx, ok bool) {
		_, _ = fmt.Fprintf(ctx, "%s ", appEnv)
		return ctx, true
	}, fasthttp.StatusNotImplemented)


	// Add some routes
	rx.Mux.Add("/", func (ctx *fasthttp.RequestCtx, mx *mux.Mux) {
		_, _ = fmt.Fprintf(ctx, "Hello, world! Requested path is %q", ctx.Path())
	})
	rx.Mux.Add("/abc", handleABC)


	// Let it rip!
	rx.Start("3020")
}
func handleABC(ctx *fasthttp.RequestCtx, mx *mux.Mux) {
	_, _ = fmt.Fprintf(ctx, "Hello ABC! Requested path is %q", ctx.Path())
}
```

# Named parameters

Arguments in the rules designated route colon. Example route: */abc/:param* , where *abc* is a static section and *:param* - the dynamic section(argument).

# Static files

The directory path to the file and its nickname as part of URL specified in the configuration file. This constants *FILE_PREF* and *FILE_PATH*
