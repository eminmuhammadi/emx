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

import "github.com/gofiber/fiber/v2"

type ApiMap struct {
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
}

func ErrorResponse(ctx *fiber.Ctx, code int, data interface{}, message string) error {
	return ctx.Status(code).JSON(
		ApiMap{
			Success: false,
			Message: message,
			Data:    data,
		},
	)
}

func SuccessResponse(ctx *fiber.Ctx, code int, data interface{}) error {
	return ctx.Status(code).JSON(
		ApiMap{
			Success: true,
			Message: "Success",
			Data:    data,
		},
	)
}
