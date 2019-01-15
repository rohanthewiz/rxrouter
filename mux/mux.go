// Original Copyright © 2016-2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

package mux

// mux

import (
	"github.com/valyala/fasthttp"
)

// mux Bxog is a simple and fast HTTP router for Go (HTTP request multiplexer).
type Mux struct {
	routes []*route
	Index  *index
	url    string
}

// New - create a new multiplexer
func New() *Mux {
	return &Mux{}
}

// Add - add a rule specifying the handler (the default method - GET, ID - as a string to this rule)
func (m *Mux) Add(url string, handler func(*fasthttp.RequestCtx, map[string]string)) *route {
	if len(url) > HTTP_PATTERN_COUNT {
		panic("URL is too long")
	} else {
		return m.newRoute(url, handler, HTTP_METHOD_DEFAULT)
	}
}

// Start the mux
func (m *Mux) Load() {
	m.Index = newIndex()
	m.Index.compile(m.routes)
	// The actual server will use fasthttp
	//server := &http.Server{
	//	Addr:           port,
	//	Handler:        nil,
	//	ReadTimeout:    READ_TIME_OUT * time.Second,
	//	WriteTimeout:   WRITE_TIME_OUT * time.Second,
	//	MaxHeaderBytes: 1 << 20,
	//}
	//m.server = server
	//http.Handle(DELIMITER_STRING, m)
	//http.Handle(FILE_PREF, http.StripPrefix(FILE_PREF, http.FileServer(http.Dir(FILE_PATH))))
	//log.Fatal(fasthttp.ListenAndServe())
}

// Shutdown - graceful stop the server
//func (m *mux) Shutdown() error {
//	if m.server == nil {
//		return fmt.Errorf("The server is not running and therefore cannot be stopped.")
//	}
//	return m.server.Shutdown(nil)
//}
//
//// Stop - aggressive stop the server
//func (m *mux) Stop() error {
//	if m.server == nil {
//		return fmt.Errorf("The server is not running and therefore cannot be stopped.")
//	}
//	return m.server.Close()
//}

// Params - extract parameters from URL
func (m *Mux) Params(ctx *fasthttp.RequestCtx, id string) map[string]string {
	out := make(map[string]string)
	if cRoute, ok := m.Index.index[m.Index.genUint(id, 0)]; ok {
		query := cRoute.genSplit(string(ctx.Path())[1:])
		for u := len(cRoute.sections) - 1; u >= 0; u-- {
			if cRoute.sections[u].typeSec == TYPE_ARG {
				out[cRoute.sections[u].id] = query[u]
			}
		}
	}
	return out
}

// GenerateURL - generate URL of the available options
// Hmm TODO - let'server rename this according to what it does
//func (r *mux) GenerateURL(id string, param map[string]string) string {
//	out := ""
//	if route, ok := r.Index.index[r.Index.genUint(id, 0)]; ok {
//		for _, section := range route.sections {
//			if section.typeSec == TYPE_STAT {
//				out = out + DELIMITER_STRING + section.id
//			} else if par, ok := param[section.id]; section.typeSec == TYPE_ARG && ok {
//				out = out + DELIMITER_STRING + par
//			}
//		}
//	}
//	return out
//}

// Test - Start analogue (for testing only)
func (m *Mux) Test() {
	m.Index = newIndex()
	m.Index.compile(m.routes)
}
