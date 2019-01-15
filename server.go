// Copyright 2019 - Rohan Allison.

package rxrouter

import (
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
	reqHandler = func(ctx *fasthttp.RequestCtx) {

		// Run middlewares - they modify ctx or fail
		var ok bool
		for _, mw := range rx.middlewares {
			if ok = mw.MidFunc(ctx); !ok {
				ctx.SetStatusCode(mw.FailCode) // for now
				return
			}
		}

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

// Default Handler
func (rx *RxRouter) Default(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusNotFound)
}

func (rx *RxRouter) Add(path string, handler RouteHandler) {
	rx.mux.Add(path, handler)
}
