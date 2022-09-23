package rxrouter

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strconv"

	"github.com/rohanthewiz/rerr"
	"github.com/rohanthewiz/rxrouter/core/constants"
	"github.com/rohanthewiz/rxrouter/test_helpers"
	"github.com/valyala/fasthttp"
)

// RunServerTest will test an endpoint on the passed in rxrouter
// An rxrouter with added routes is the first arg;
// the request containing the endpoint to test is the second arg
// See github.com/rohanthewiz/rxrun for example usage
//
// rx := New(Options{Verbose: true})
// rx.AddRoute("/hello/:name/:age",
// 		func(ctx *fasthttp.RequestCtx, params map[string]string) {
// 		  _, _ = ctx.WriteString(fmt.Sprintf("Hello %s. You are %s!", params["name"], params["age"]))
//		}
// })
func RunServerTest(rx *RxRouter, req *http.Request) (resp *fasthttp.Response, err error) {
	if req.Header.Get(constants.HeaderContentType) == "" {
		req.Header.Add(constants.HeaderContentType, constants.ContentTypeText)
	}

	if req.Body != http.NoBody && req.Header.Get(constants.HeaderContentLength) == "" {
		req.Header.Add(constants.HeaderContentLength, strconv.FormatInt(req.ContentLength, 10))
	}

	reqRaw, err := httputil.DumpRequest(req, true)
	if err != nil {
		return resp, rerr.Wrap(err, "Error obtaining raw HTTP request from req",
			"request", fmt.Sprintf("%v", req))
	}

	rx.LoadRoutes() // compile routes

	s := &fasthttp.Server{
		Handler: InitStdMasterHandler(rx), // Standard route table handling
	}

	cw := &test_helpers.ConnWrap{} // connection obj for testing
	cw.R.Write(reqRaw)

	if err := s.ServeConn(cw); err != nil {
		return resp, rerr.Wrap(err, "Unexpected error from serveConn",
			"request", fmt.Sprintf("%v", req))
	}

	body, err := ioutil.ReadAll(&cw.W)
	if err != nil {
		return resp, rerr.Wrap(err, "Unexpected error from ReadAll",
			"request", fmt.Sprintf("%v", req))
	}

	resp.SetBodyRaw(body)
	// fmt.Printf("**-> StatusCode: %d\n", resp.StatusCode())

	return
}
