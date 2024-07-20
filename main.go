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
package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/eminmuhammadi/emx/app"
	"github.com/eminmuhammadi/emx/pkg/http"
	"github.com/eminmuhammadi/emx/pkg/logger"
	"github.com/eminmuhammadi/emx/pkg/proxy"
	"github.com/eminmuhammadi/emx/pkg/sql"
)

func main() {
	/*
	  Graceful shutdown configuration
	*/
	var wg sync.WaitGroup
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	ctx := context.Background()

	/*
	  Main Application Servers
	*/
	httpServer := http.CreateServer()
	proxyServer := proxy.CreateProxyServer()

	/*
		Starting the servers
		===============================
	*/
	wg.Add(2)

	// (1) Start the application server
	go func() {
		defer wg.Done()

		app.RegisterRoutes(httpServer)
		http.StartServer(httpServer)
	}()

	// (2) Start the proxy server
	go func() {
		defer wg.Done()

		proxy.StartProxyServer(proxyServer)
	}()

	/*
		Shutdown the servers
		===============================
	*/
	<-shutdown // Wait for the shutdown signal

	if err := sql.CloseConnection(sql.Database); err != nil {
		logger.Log.Printf("error while closing database connection: %v", err)
	}

	logger.Log.Println("Shutting down application server...")

	if err := httpServer.Shutdown(); err != nil {
		logger.Log.Printf("error while shutting down server: %v", err)
	}

	logger.Log.Println("Shutting down proxy server...")

	if err := proxyServer.Shutdown(ctx); err != nil {
		logger.Log.Printf("error while shutting down proxy server: %v", err)
	}

	wg.Wait() // Wait for all goroutines to finish

	logger.Log.Println("Service gracefully shutdown")
}
