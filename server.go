// Copyright 2019 - Rohan Allison.

package rxrouter

import (
	"bytes"
	"fmt"
	"github.com/rohanthewiz/rxrouter/mux"
	"github.com/valyala/fasthttp"
	"log"
)

type RxRouter struct {
	Options     Options
	mux         *mux.Mux
	middlewares []MiddleWare
}

type Options struct {
	Verbose bool
	AssetPaths []AssetPath
}

type AssetPath struct {
	Prefix []byte // url prefix
	FileSystemRoot string // file locations
	StripSlashes int // how many slash words to strip from the url prefix
}

type RouteHandler func(*fasthttp.RequestCtx, map[string]string)

type MiddleWareFunc func(ctx *fasthttp.RequestCtx) (ok bool)

type MiddleWare struct {
	MidFunc  MiddleWareFunc
	FailCode int
}

func New(opts ...Options) *RxRouter {
	mx := &mux.Mux{}
	r := RxRouter{mux: mx}
	if len(opts) > 0 {
		r.Options = opts[0]
	}
	return &r
}

func (rx *RxRouter) Use(m MiddleWareFunc, failCode int) {
	rx.middlewares = append(rx.middlewares, MiddleWare{MidFunc: m, FailCode: failCode})
}

func (rx *RxRouter) Start(port string) {
	if rx.Options.Verbose {
		fmt.Println("Compiling routes...")
	}
	rx.mux.Load() // create new index; compile routes

	var reqHandler fasthttp.RequestHandler

	// Master handler - select and run a handler on the passed ctx
	reqHandler = func(ctx *fasthttp.RequestCtx) {
		// Middleware - they modify ctx or fail
		var ok bool
		for _, mw := range rx.middlewares {
			if ok = mw.MidFunc(ctx); !ok {
				ctx.SetStatusCode(mw.FailCode) // for now
				return
			}
		}

		// See if we match a file handler
		fileHander, ok := rx.GetFSHandler(ctx)
		if ok {
			fileHander(ctx)
			return
		}

		// Lookup
		if route := rx.mux.Index.FindTree(ctx); route != nil {
			fmt.Printf("route is: %s\n", route.Url())
			route.Handler(ctx, rx.mux.Params(ctx, route.Url()))
		} else {
			if rx.Options.Verbose {
				fmt.Println("Unknown route", string(ctx.Path()))
			}
			rx.Default(ctx)
		}
	}

	if rx.Options.Verbose {
		fmt.Println("RxRouter is listening on port " + port)
	}
	log.Fatal(fasthttp.ListenAndServe(":"+port, reqHandler))
}

// See if we match a file handler - First match is the one we use
func (rx *RxRouter) GetFSHandler(ctx *fasthttp.RequestCtx) (handler fasthttp.RequestHandler, ok bool) {
	path := ctx.Path()
	for _, astPath := range rx.Options.AssetPaths {
		if bytes.HasPrefix(path, astPath.Prefix) {
			return fasthttp.FSHandler(astPath.FileSystemRoot, astPath.StripSlashes), true
		}
	}
	return
}

// Default Handler
func (rx *RxRouter) Default(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusNotFound)
}

func (rx *RxRouter) Add(path string, handler RouteHandler) {
	rx.mux.Add(path, handler)
}
