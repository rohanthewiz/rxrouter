package rxrouter

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"net/http/httputil"
	"testing"

	"github.com/rohanthewiz/rxrouter/test_helpers"
	"github.com/valyala/fasthttp"
)

func TestThatServerWorks(t *testing.T) {
	// t.Parallel() // ?

	req := httptest.NewRequest("GET", "/hello/there", nil)
	// req := httptest.NewRequest("GET", "/hello/there/foe", nil)
	// req := httptest.NewRequest("GET", "/hello/john/24", nil)
	reqRaw, err := httputil.DumpRequest(req, true)
	if err != nil {
		t.Errorf("error %s - obtaining raw HTTP request from req %v", err.Error(), req)
		t.Fail()
	}

	// FIXUPS

	rx := New(Options{Verbose: true})

	// rx.AddRoute("/hello/:name/:age", func(ctx *fasthttp.RequestCtx, params map[string]string) {
	// 	_, _ = ctx.WriteString(fmt.Sprintf("Hello %s. You are %s!", params["name"], params["age"]))
	// })
	rx.AddRoute("/hello/there/friend/sue", func(ctx *fasthttp.RequestCtx, params map[string]string) {
		_, _ = ctx.WriteString("Hello Sue")
	})
	rx.AddRoute("/hello/there", func(ctx *fasthttp.RequestCtx, params map[string]string) {
		_, _ = ctx.WriteString("Hello there")
	})
	rx.AddRoute("/hello/there/foe/fah", func(ctx *fasthttp.RequestCtx, params map[string]string) {
		_, _ = ctx.WriteString("Hello foe")
	})
	rx.AddRoute("/hello/harry", func(ctx *fasthttp.RequestCtx, params map[string]string) {
		_, _ = ctx.WriteString("Hello Harry")
	})

	rx.LoadRoutes()

	s := &fasthttp.Server{
		Handler: InitStdMasterHandler(rx),
	}

	cw := &test_helpers.ConnWrap{}
	cw.R.Write(reqRaw)
	// Example: cw.R.WriteString("GET /hello/john/24 HTTP/1.1\r\nHost: localhost\r\n\r\n")

	if err := s.ServeConn(cw); err != nil {
		t.Fatalf("Unexpected error from serveConn: %v", err)
	}

	body, err := ioutil.ReadAll(&cw.W)
	if err != nil {
		t.Fatalf("Unexpected error from ReadAll: %v", err)
	}

	/*	fresp := fasthttp.Response{}
		fresp.SetBodyRaw(body)
		fmt.Printf("**-> StatusCode: %d\n", fresp.StatusCode())
	*/
	fmt.Println("**-> Server resp", string(body))
}
