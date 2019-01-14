// Copyright © 2016-2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

package rxrouter

// Server

import (
	"fmt"
	"github.com/rohanthewiz/rxrouter/bxog"
	"github.com/valyala/fasthttp"
)

type RxRouter struct {
	Mux *bxog.Mux
}

func New() *RxRouter {
	mx := &bxog.Mux{}
	mx.Add("/", func (ctx *fasthttp.RequestCtx, mx *bxog.Mux) {
		_, _ = fmt.Fprintf(ctx, "Hello, world! Requested path is %q", ctx.Path())
	})
	mx.Add("/abc", func (ctx *fasthttp.RequestCtx, mx *bxog.Mux) {
		_, _ = fmt.Fprintf(ctx, "Hello ABC! Requested path is %q", ctx.Path())
	})
	
	return &RxRouter{ Mux: mx }
}

func (rx *RxRouter) Start() {
	rx.Mux.Load() // create new index; compile routes

	reqHandler := func(ctx *fasthttp.RequestCtx) {
		if route := rx.Mux.Index.FindTree(ctx); route != nil {
			route.Handler(ctx, rx.Mux)
		} else {
			rx.Default(ctx)
		}
	}

	fasthttp.ListenAndServe(":3200", reqHandler)
}


// Default Handler
func (rx *RxRouter) Default(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusNotFound)
}
