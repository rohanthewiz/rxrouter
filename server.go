// Copyright © 2016-2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

package rxrouter

// Server

import (
	"net/http"
	"github.com/rohanthewiz/rxrouter/bxog
)

type RxRouter struct {
	BxRtr *bxog.Router
}

func New() *RxRouter {
	return &RxRouter{}
}

func (rx *RxRouter) Start() {
	rx.BxRtr.Start() // create new index; compile routes
}

// ServeHTTP looks for a matching route
func (rx *RxRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if route := rx.BxRtr.Index.FindTree(req); route != nil {
		route.Handler(w, req, rx.BxRtr)
	} else {
		rx.Default(w, req)
	}
}

// Default Handler
func (rx *RxRouter) Default(w http.ResponseWriter, req *http.Request) {
	// w.WriteHeader(404)
	http.Error(w, "Page not found", 404)
}
