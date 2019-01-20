// Copyright 2019 - Rohan Allison.

package rxrouter

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/rohanthewiz/rxrouter/mux"
	"github.com/valyala/fasthttp"
	"log"
	"net"
)

const defaultPort = "3020"
const defaultTLSPort = "443"

type RxRouter struct {
	Options     Options
	mux         *mux.Mux
	middlewares []MiddleWare
}

type Options struct {
	Verbose    bool
	Port string
	TLSCfg *tls.Config
	assetPaths []AssetPath
}

type AssetPath struct {
	Prefix         []byte // url prefix
	FileSystemRoot string // file locations
	StripSlashes   int    // how many slash words to strip from the url prefix
}

type RouteHandler func(*fasthttp.RequestCtx, map[string]string)

type MiddleWareFunc func(ctx *fasthttp.RequestCtx) (ok bool)

type MiddleWare struct {
	MidFunc  MiddleWareFunc
	FailCode int
}

// Create a new router instance
// Afterwards you will want to add some routes then call the instance's Start()
func New(opts ...Options) *RxRouter {
	mx := &mux.Mux{}
	r := RxRouter{mux: mx}
	if len(opts) > 0 {
		r.Options = opts[0]
	}
	return &r
}

// Add a middleware function before regular routes
// failCode is the htttp response code to return if we are immediately stopping the request
func (rx *RxRouter) Use(m MiddleWareFunc, failCode int) {
	rx.middlewares = append(rx.middlewares, MiddleWare{MidFunc: m, FailCode: failCode})
}

// Start serving routes
func (rx *RxRouter) Start() {
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
			if rx.Options.Verbose {
				fmt.Printf("route is: %s\n", route.Url())
			}
			route.Handler(ctx, rx.mux.Params(ctx, route.Url()))
		} else {
			if rx.Options.Verbose {
				fmt.Println("Unknown route", string(ctx.Path()))
			}
			rx.Default(ctx)
		}
	}

	if rx.Options.Verbose {
		fmt.Println("RxRouter is listening on port " + rx.Options.Port)
	}
	if rx.Options.Port == "" { rx.Options.Port = defaultPort }

	if rx.Options.TLSCfg != nil {
		if rx.Options.Port == "" { rx.Options.Port = defaultTLSPort }
		ln, err := net.Listen("tcp4", "0.0.0.0." + rx.Options.Port)
		if err != nil {
			panic(err)
		} // todo - better handling of err
		lnTls := tls.NewListener(ln, rx.Options.TLSCfg)
		if err := fasthttp.Serve(lnTls, reqHandler); err != nil {
			panic(err) // todo - better handling here too
		}
	} else {
		log.Fatal(fasthttp.ListenAndServe(":"+rx.Options.Port, reqHandler))
	}
}

// See if we match a file handler - First match is the one we use
func (rx *RxRouter) GetFSHandler(ctx *fasthttp.RequestCtx) (handler fasthttp.RequestHandler, ok bool) {
	path := ctx.Path()
	for _, astPath := range rx.Options.assetPaths {
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

// Map a handler to a route
func (rx *RxRouter) AddRoute(path string, handler RouteHandler) {
	rx.mux.Add(path, handler)
}

// Add a route to static files
// Prefix is a starting portion of the URL delimited by slashes
// fsRoot is the path to the top-level folder to serve files from
// StripSlashes is the number of slash delimited tokens to remove from the URL
// before appending it to the fsRoot to form the full file path
// Todo - example
func (rx *RxRouter) AddStaticFilesRoute(prefix, fsRoot string, slashesToStrip int) {
	ap := AssetPath{Prefix: []byte(prefix), FileSystemRoot: fsRoot, StripSlashes: slashesToStrip}
	rx.Options.assetPaths = append(rx.Options.assetPaths, ap)
}
