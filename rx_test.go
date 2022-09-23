package rxrouter

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strconv"
	"testing"

	"github.com/rohanthewiz/rxrouter/core/constants"
	"github.com/rohanthewiz/rxrouter/test_helpers"
	"github.com/valyala/fasthttp"
)

func TestThatServerWorks(t *testing.T) {
	// t.Parallel() // ?

	req := httptest.NewRequest("GET", "/hello/john/24", nil)

	// FIXUPS

	if req.Header.Get(constants.HeaderContentType) == "" {
		req.Header.Add(constants.HeaderContentType, constants.ContentTypeText)
	}

	if req.Body != http.NoBody && req.Header.Get(constants.HeaderContentLength) == "" {
		req.Header.Add(constants.HeaderContentLength, strconv.FormatInt(req.ContentLength, 10))
	}

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
