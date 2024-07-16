/*
Copyright 2024 Emin Muhammadi

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/elazarl/goproxy"
	"github.com/eminmuhammadi/emx/app/request"
	"github.com/eminmuhammadi/emx/app/response"
	"github.com/eminmuhammadi/emx/pkg/logger"
	"github.com/eminmuhammadi/emx/pkg/util"
)

var seedString = func() string {
	return util.GenerateUlid()
}()

func HashSession(ctx *goproxy.ProxyCtx) string {
	sessionID := fmt.Sprintf("%d", ctx.Session)
	seed := seedString

	return util.GenerateHash(sessionID + seed)
}

func CreateProxyServer() *http.Server {
	hostname, port := os.Getenv("PROXY_HOST"), os.Getenv("PROXY_PORT")
	if hostname == "" || port == "" {
		log.Fatalf("PROXY_HOST and PROXY_PORT environment variables must be set")
	}

	proxy := newMitmProxy()
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", hostname, port),
		Handler: proxy,
	}

	return server
}

func StartProxyServer(server *http.Server) {
	logger.Log.Printf("Proxy server started on %s", server.Addr)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Error starting proxy server: %s", err)
	}
}

func RequestHandler(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	if IsMockPatternMatches(ctx) {
		req, _ = MockRequestHandler(req, ctx)
	}

	reqData := request.CreateRequestModel(ctx)
	reqData.SessionID = HashSession(ctx)

	reqBody := util.NormalizeBody(req.Body)
	reqData.Body = reqBody

	// Save into database
	request.Create(reqData)

	// Restore the original request body
	req.Body = io.NopCloser(strings.NewReader(reqBody))

	return req, nil
}

func ResponseHandler(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	if IsMockPatternMatches(ctx) {
		resp = MockResponseHandler(resp, ctx)
	}

	respData := response.CreateResponseModel(ctx)
	respData.SessionID = HashSession(ctx)

	respBody := util.NormalizeBody(resp.Body)
	respData.Body = respBody

	// Save into database
	response.Create(respData)

	// Restore the original response body
	resp.Body = io.NopCloser(strings.NewReader(respBody))

	return resp
}
