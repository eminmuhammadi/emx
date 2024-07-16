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
	"net/http"
	"os"
	"strings"

	"github.com/elazarl/goproxy"
	"github.com/eminmuhammadi/emx/pkg/logger"
	"github.com/eminmuhammadi/emx/pkg/util"
	"gopkg.in/yaml.v3"
)

type MockPattern struct {
	Method   string       `yaml:"method"`
	Host     string       `yaml:"host"`
	Path     string       `yaml:"path"`
	Response ResponseMock `yaml:"response"`
}

type ResponseMock struct {
	StatusCode int    `yaml:"status_code"`
	Headers    string `yaml:"headers"`
	Body       string `yaml:"body"`
}

type PatternConfig struct {
	Patterns []MockPattern `yaml:"patterns"`
}

func LoadMockPatterns() []MockPattern {
	filePath := os.Getenv("MOCK_FILE")
	if filePath == "" {
		return []MockPattern{}
	}

	file, err := os.Open(filePath)
	if err != nil {
		logger.Log.Fatalf("Error opening mock responses file: %s", err)
	}
	defer file.Close()

	var config PatternConfig
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		logger.Log.Fatalf("Error decoding mock responses file: %s", err)
	}

	return config.Patterns
}

var MockPatterns = func() []MockPattern {
	patterns := LoadMockPatterns()
	logger.Log.Println("Loaded mock patterns: ", patterns)

	return patterns
}()

func FindPattern(ctx *goproxy.ProxyCtx) *MockPattern {
	for _, pattern := range MockPatterns {
		if pattern.Method == ctx.Req.Method && pattern.Host == ctx.Req.Host && pattern.Path == ctx.Req.URL.Path {
			return &pattern
		}
	}
	return nil
}

func IsMockPatternMatches(ctx *goproxy.ProxyCtx) bool {
	for _, pattern := range MockPatterns {
		if pattern.Method == ctx.Req.Method && pattern.Host == ctx.Req.Host && pattern.Path == ctx.Req.URL.Path {
			return true
		}
	}
	return false
}

func MockRequestHandler(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	return ctx.Req, nil
}

func MockResponseHandler(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	pattern := FindPattern(ctx)

	// Update the response with the mock pattern dynamically
	if pattern != nil {
		headers := make(http.Header)

		for _, line := range strings.Split(pattern.Response.Headers, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			headers.Add(key, value)
		}

		resp.StatusCode = pattern.Response.StatusCode
		resp.Header = headers
		resp.Body = util.StringToReader(pattern.Response.Body)
		resp.ContentLength = int64(len(pattern.Response.Body))
	}

	return resp
}
