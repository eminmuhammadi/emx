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
package http

import (
	"os"
	"path"
	"strings"

	"github.com/eminmuhammadi/emx/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func ListenAndServe(app *fiber.App, addr string) error {
	tlsMode := strings.ToLower(os.Getenv("TLS_MODE"))
	certFile := path.Join(os.Getenv("TLS_CERT_FILE"))
	keyFile := path.Join(os.Getenv("TLS_KEY_FILE"))
	caFile := path.Join(os.Getenv("TLS_CA_FILE"))

	if tlsMode == "tls" || tlsMode == "mutual_tls" {
		if certFile == "" || keyFile == "" || caFile == "" {
			logger.Log.Fatalf("TLS_CERT_FILE, TLS_KEY_FILE, TLS_CA_FILE environment variables must be set for TLS_MODE=%s", tlsMode)
		}
	}

	switch tlsMode {
	case "mutual_tls":
		return StartMutualTlsServer(app, addr, certFile, keyFile, caFile)
	case "tls":
		return StartTlsServer(app, addr, certFile, keyFile)
	default:
		return StartInsecureServer(app, addr)
	}
}

func StartInsecureServer(app *fiber.App, addr string) error {
	return app.Listen(addr)
}

func StartMutualTlsServer(app *fiber.App, addr, certFile, keyFile string, caFile string) error {
	return app.ListenMutualTLS(addr, certFile, keyFile, caFile)
}

func StartTlsServer(app *fiber.App, addr, certFile, keyFile string) error {
	return app.ListenTLS(addr, certFile, keyFile)
}
