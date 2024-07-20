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
package http

import (
	"fmt"
	"os"
	"time"

	"github.com/eminmuhammadi/emx/pkg/logger"
	jsonTranscoder "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
)

func CreateServer() *fiber.App {
	app := fiber.New(fiber.Config{
		EnablePrintRoutes:       false,
		ReduceMemoryUsage:       true,
		DisableStartupMessage:   true,
		Prefork:                 false,
		WriteTimeout:            time.Duration(15 * time.Second),
		IdleTimeout:             time.Duration(30 * time.Second),
		ReadTimeout:             time.Duration(15 * time.Second),
		WriteBufferSize:         4 * 1024,
		ReadBufferSize:          4 * 1024,
		Concurrency:             256 * 1024,
		BodyLimit:               1 * 1024 * 1024, // 1MB
		JSONEncoder:             jsonTranscoder.Marshal,
		JSONDecoder:             jsonTranscoder.Unmarshal,
		ErrorHandler:            errorHandler,
		TrustedProxies:          []string{""},
		ProxyHeader:             "X-Forwarded-For",
		EnableTrustedProxyCheck: true,
	})

	ApplyMiddleware(app)

	return app
}

func ApplyMiddleware(app *fiber.App) {
	logger.Log.Println("Applying middleware")

	app.Use(compress.New(
		compress.Config{
			Level: compress.LevelBestSpeed,
		},
	))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Accept,Authorization,Content-Type,X-CSRF-TOKEN",
		ExposeHeaders:    "Link",
		AllowCredentials: false,
		MaxAge:           300,
	}))

	app.Use(logger.HttpLogger)
	app.Use(etag.New())
	app.Use(helmet.New())

	logger.Log.Println("Middleware applied")
}

func errorHandler(ctx *fiber.Ctx, err error) error {
	if e, ok := err.(*fiber.Error); ok {
		return ctx.Status(e.Code).JSON(ApiMap{
			Success: false,
			Data:    nil,
			Message: e.Message,
		})
	}

	logger.Log.Println(err)

	return ctx.Status(fiber.StatusInternalServerError).JSON(ApiMap{
		Success: false,
		Data:    nil,
		Message: "Internal server error",
	})
}

func StartServer(app *fiber.App) {
	addr := fmt.Sprintf(
		"%s:%s",
		os.Getenv("APP_HOST"),
		os.Getenv("APP_PORT"),
	)

	if addr == ":" {
		logger.Log.Fatalf("APP_HOST and APP_PORT must be set")
	}

	protocol := "http"
	if os.Getenv("TLS_MODE") == "mutual_tls" || os.Getenv("TLS_MODE") == "tls" {
		protocol = "https"
	}

	logger.Log.Printf("Application server started on %s://%s:%s \n",
		protocol,
		os.Getenv("APP_HOST"),
		os.Getenv("APP_PORT"))

	if err := ListenAndServe(app, addr); err != nil {
		logger.Log.Fatalf("error while starting server: %v", err)
	}
}
