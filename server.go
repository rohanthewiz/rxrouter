// Original Copyright © 2016-2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

package rxrouter

// Server

import (
	"github.com/rohanthewiz/rxrouter/mux"
	"github.com/valyala/fasthttp"
	"log"
)

type RxRouter struct {
	Mux *mux.Mux
}

type Handler func(*fasthttp.RequestCtx, *mux.Mux)

func New() *RxRouter {
	mx := &mux.Mux{}
	return &RxRouter{Mux: mx}
}

func (rx *RxRouter) Start(port string) {
	rx.Mux.Load() // create new index; compile routes

	reqHandler := func(ctx *fasthttp.RequestCtx) {
		if route := rx.Mux.Index.FindTree(ctx); route != nil {
			route.Handler(ctx, rx.Mux)
		} else {
			rx.Default(ctx)
		}
	}

	log.Fatal(fasthttp.ListenAndServe(":" + port, reqHandler))
}

// Default Handler
func (rx *RxRouter) Default(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusNotFound)
}

func (rx *RxRouter) Add(path string, handler Handler) {
	rx.Mux.Add(path, handler)
}
