package rxrouter

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strconv"
	"testing"

	"github.com/rohanthewiz/rxrouter/test_helpers"
	"github.com/valyala/fasthttp"
)

const HeaderContentLength = "Content-Length"

func TestThatServerWorks(t *testing.T) {
	// t.Parallel() // ?

	req := httptest.NewRequest("GET", "/hello/john/24", nil)
	req.Header.Set("content-type", "text/html") // "application/json"

	// Add Content-Length if not provided with header
	if req.Body != http.NoBody && req.Header.Get(HeaderContentLength) == "" {
		req.Header.Add(HeaderContentLength, strconv.FormatInt(req.ContentLength, 10))
	}

	// Dump raw http request
	reqRaw, err := httputil.DumpRequest(req, true)
	if err != nil {
		t.Fatalf("Error obtaining raw HTTP request: %v", reqRaw)
	}

	rx := New(Options{Verbose: true})

	rx.AddRoute("/hello/:name/:age", func(ctx *fasthttp.RequestCtx, params map[string]string) {
		_, _ = ctx.WriteString(fmt.Sprintf("Hello %s. You are %s!", params["name"], params["age"]))
	})
	rx.LoadRoutes()

	s := &fasthttp.Server{
		Handler: InitStdMasterHandler(rx),
	}

	getResp := func() []byte {
		cw := &test_helpers.ConnWrap{}
		cw.R.Write(reqRaw)
		// cw.R.WriteString("GET /hello/john/24 HTTP/1.1\r\nHost: localhost\r\n\r\n")

		if err := s.ServeConn(cw); err != nil {
			t.Fatalf("Unexpected error from serveConn: %v", err)
		}

		resp, err := ioutil.ReadAll(&cw.W)
		if err != nil {
			t.Fatalf("Unexpected error from ReadAll: %v", err)
		}
		return resp
	}

	resp := getResp()

	fmt.Println("**-> Server resp", string(resp))
}
