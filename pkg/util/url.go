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
package util

import (
	"net/url"
	"strings"
)

var specialChars = []string{"\"", "'", "`"}

func NormalizeQueryValue(value string) (string, error) {
	value, err := url.QueryUnescape(value)
	if err != nil {
		return "", err
	}

	for _, char := range specialChars {
		value = strings.TrimPrefix(value, char)
	}

	for _, char := range specialChars {
		value = strings.TrimSuffix(value, char)
	}

	return value, nil
}
