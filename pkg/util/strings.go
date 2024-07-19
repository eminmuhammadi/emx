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
package util

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

func StringifyArray(arr []string) string {
	var str string
	for _, item := range arr {
		str += item + " "
	}

	return str
}

func NormalizeHeaders(headers http.Header) string {
	var normalizedHeaders string
	for key, values := range headers {
		for _, value := range values {
			normalizedHeaders += key + ": " + value + "\n"
		}
	}

	return normalizedHeaders
}

func NormalizeBody(body io.ReadCloser) string {
	var bodyCopy bytes.Buffer
	tee := io.TeeReader(body, &bodyCopy)

	bodyBytes, err := io.ReadAll(tee)
	if err != nil {
		return ""
	}

	return string(bodyBytes)
}

func StringToReader(str string) io.ReadCloser {
	return io.NopCloser(strings.NewReader(str))
}

func Contains(arr []string, str string) bool {
	for _, item := range arr {
		if item == str {
			return true
		}
	}

	return false
}
