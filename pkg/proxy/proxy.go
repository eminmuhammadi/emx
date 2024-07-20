/*
Copyright (c) 2024 Emin Muhammadi

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
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/elazarl/goproxy"
	"github.com/eminmuhammadi/emx/pkg/logger"
)

func loadTLSKeyPair() tls.Certificate {
	certFile, keyFile := os.Getenv("PROXY_DECRYPT_CERT_FILE"), os.Getenv("PROXY_DECRYPT_KEY_FILE")
	if certFile == "" || keyFile == "" {
		logger.Log.Fatal("PROXY_DECRYPT_CERT_FILE and PROXY_DECRYPT_KEY_FILE must be set")
	}

	rootX509KeyPair, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		logger.Log.Fatalf("Failed to load proxy root certificate: %v", err)
	}

	rootX509KeyPair.Leaf, err = x509.ParseCertificate(rootX509KeyPair.Certificate[0])
	if err != nil {
		logger.Log.Fatalf("Failed to parse proxy root certificate: %v", err)
	}

	return rootX509KeyPair
}

func newProxyHttpServer() *goproxy.ProxyHttpServer {
	caKeyPair := loadTLSKeyPair()
	goproxy.GoproxyCa = caKeyPair

	logger.Log.Println("Proxy CA loaded")

	goproxy.OkConnect = &goproxy.ConnectAction{
		Action:    goproxy.ConnectAccept,
		TLSConfig: goproxy.TLSConfigFromCA(&caKeyPair),
	}

	goproxy.MitmConnect = &goproxy.ConnectAction{
		Action:    goproxy.ConnectMitm,
		TLSConfig: goproxy.TLSConfigFromCA(&caKeyPair),
	}

	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{
		Action:    goproxy.ConnectHTTPMitm,
		TLSConfig: goproxy.TLSConfigFromCA(&caKeyPair),
	}

	goproxy.RejectConnect = &goproxy.ConnectAction{
		Action:    goproxy.ConnectReject,
		TLSConfig: goproxy.TLSConfigFromCA(&caKeyPair),
	}

	return goproxy.NewProxyHttpServer()
}

func newMitmProxy() *goproxy.ProxyHttpServer {
	proxy := newProxyHttpServer()
	proxy.Logger = logger.Log
	proxy.Verbose = false

	if os.Getenv("PROXY_VERBOSE") == "true" {
		proxy.Verbose = true
	}

	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest().DoFunc(RequestHandler)
	proxy.OnResponse().DoFunc(ResponseHandler)

	return proxy
}
