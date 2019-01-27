package rxrouter

import "github.com/valyala/fasthttp"

type MiddleWareFunc func(ctx *fasthttp.RequestCtx) (ok bool)

type MiddleWare struct {
	MidFunc  MiddleWareFunc
	FailCode int
}

// Add a middleware function before regular routes
// failCode is the htttp response code to return if we are immediately stopping the request
func (rx *RxRouter) Use(m MiddleWareFunc, failCode int) {
	rx.middlewares = append(rx.middlewares, MiddleWare{MidFunc: m, FailCode: failCode})
}
