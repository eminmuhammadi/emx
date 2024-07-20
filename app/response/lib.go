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
package response

import (
	"github.com/eminmuhammadi/emx/pkg/sql"
)

func Create(response Response) error {
	result := sql.Database.Create(&response)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetBySessionID(sessionID string) Response {
	var response Response
	result := sql.Database.Where("session_id = ?", sessionID).First(&response)
	if result.Error != nil {
		return Response{}
	}

	return response
}

func Get(id int64) (Response, error) {
	var response Response
	result := sql.Database.First(&response, id)
	if result.Error != nil {
		return response, result.Error
	}

	return response, nil
}
