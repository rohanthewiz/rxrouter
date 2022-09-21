// Copyright 2019 - Rohan Allison.

package rxrouter

import (
	"fmt"
	"log"

	"github.com/rohanthewiz/rxrouter/mux"
	"github.com/valyala/fasthttp"
)

const defaultPort = "3020"
const defaultTLSPort = "443"
const ipAny = "0.0.0.0"

type RxRouter struct {
	Options     Options
	mux         *mux.Mux
	middlewares []MiddleWare
}

type Options struct {
	Verbose             bool
	Port                string
	TLS                 RxTLS
	assetPaths          []AssetPath
	CustomMasterHandler *fasthttp.RequestHandler
}

// Specify whether to use TLS and
// CertFile and KeyFile or CertData and KeyData (for embedded certs)
type RxTLS struct {
	UseTLS   bool
	CertFile string
	KeyFile  string
	CertData []byte
	KeyData  []byte
}

type RouteHandler func(*fasthttp.RequestCtx, map[string]string)

// New creates a new router instance
// Afterwards you will want to add some routes then call the instance's Start()
func New(opts ...Options) *RxRouter {
	r := RxRouter{mux: &mux.Mux{}}
	if len(opts) > 0 {
		r.Options = opts[0]
	}
	return &r
}

// Start serving routes
func (rx *RxRouter) Start() {
	if rx.Options.Verbose {
		fmt.Println("Compiling routes...")
	}
	rx.LoadRoutes()

	var reqHandler fasthttp.RequestHandler

	// Master handler
	if rx.Options.CustomMasterHandler != nil {
		reqHandler = *rx.Options.CustomMasterHandler
	} else {
		reqHandler = InitStdMasterHandler(rx)
	}

	if rx.Options.Port == "" {
		if rx.Options.TLS.UseTLS {
			rx.Options.Port = defaultTLSPort
		} else {
			rx.Options.Port = defaultPort
		}
	}

	if rx.Options.Verbose {
		fmt.Println("RxRouter is listening on port " + rx.Options.Port)
	}

	if rx.Options.TLS.UseTLS && rx.Options.TLS.CertFile != "" {
		log.Fatal(fasthttp.ListenAndServeTLS(ipAny+":"+rx.Options.Port, rx.Options.TLS.CertFile,
			rx.Options.TLS.KeyFile, reqHandler))
	} else if rx.Options.TLS.UseTLS && len(rx.Options.TLS.CertData) > 0 {
		log.Fatal(fasthttp.ListenAndServeTLSEmbed(ipAny+":"+rx.Options.Port, rx.Options.TLS.CertData,
			rx.Options.TLS.KeyData, reqHandler))
	} else {
		log.Fatal(fasthttp.ListenAndServe(":"+rx.Options.Port, reqHandler))
	}
}

// LoadRoutes creates new index and compiles routes
func (rx *RxRouter) LoadRoutes() {
	rx.mux.Load()
}

func InitStdMasterHandler(rx *RxRouter) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		// Middlewares - they modify ctx or fail with the provided code
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

		// Route Lookup
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
}

func (rx *RxRouter) AddCustomMasterHandler(mh *fasthttp.RequestHandler) {
	rx.Options.CustomMasterHandler = mh
}

// Default Handler
func (rx *RxRouter) Default(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusNotFound)
}

// AddRoute maps a handler to a route
func (rx *RxRouter) AddRoute(path string, handler RouteHandler) {
	rx.mux.Add(path, handler)
}
