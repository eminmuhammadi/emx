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
package app

import (
	"github.com/eminmuhammadi/emx/app/log"
	"github.com/eminmuhammadi/emx/app/request"
	"github.com/eminmuhammadi/emx/app/response"
	"github.com/eminmuhammadi/emx/app/ui"
	"github.com/eminmuhammadi/emx/pkg/sql"
	"github.com/gofiber/fiber/v2"
)

func init() {
	sql.Migrate(
		&request.Request{},
		&response.Response{},
	)
}

func RegisterRoutes(app *fiber.App) {
	/*
		/api/v1/log
		Query Params: limit (int), cursor (int), reverse (bool)

		/api/v1/log/:id
		Path Params: id (int)
	*/
	log.LogController(app)

	/*
		/api/v1/request/:id
		Path Params: id (int)
	*/
	request.RequestController(app)

	/*
		/api/v1/response/:id
		Path Params: id (int)
	*/
	response.ResponseController(app)

	/*
		/ui
	*/
	ui.ResponseController(app)
}
