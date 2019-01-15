// Copyright 2019 - Rohan Allison.
// Original Copyright © 2016-2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

package rxrouter

import (
	"fmt"
	"github.com/rohanthewiz/rxrouter/mux"
	"github.com/valyala/fasthttp"
	"log"
)

type RxRouter struct {
	Mux         *mux.Mux
	middlewares []MiddleWare
}

type RouteHandler func(*fasthttp.RequestCtx, *mux.Mux)

type MiddleWareFunc func(ctx *fasthttp.RequestCtx) (ok bool)

type MiddleWare struct {
	MidFunc  MiddleWareFunc
	FailCode int
}

func New() *RxRouter {
	mx := &mux.Mux{}
	return &RxRouter{Mux: mx}
}

func (rx *RxRouter) Use(m MiddleWareFunc, failCode int) {
	rx.middlewares = append(rx.middlewares, MiddleWare{MidFunc: m, FailCode: failCode})
}

func (rx *RxRouter) Start(port string) {
	fmt.Println("Compiling routes...")
	rx.Mux.Load() // create new index; compile routes

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

		if route := rx.Mux.Index.FindTree(ctx); route != nil {
			route.Handler(ctx, rx.Mux)
		} else {
			rx.Default(ctx)
		}
	}

	fmt.Println("RxRouter is listening on port " + port)
	log.Fatal(fasthttp.ListenAndServe(":"+port, reqHandler))
}

// Default Handler
func (rx *RxRouter) Default(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusNotFound)
}

func (rx *RxRouter) Add(path string, handler RouteHandler) {
	rx.Mux.Add(path, handler)
}
