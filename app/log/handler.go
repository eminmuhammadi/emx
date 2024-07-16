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
package log

import (
	"strconv"

	"github.com/eminmuhammadi/emx/pkg/http"
	"github.com/gofiber/fiber/v2"
)

func GetHandler(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return http.ErrorResponse(ctx, fiber.StatusBadRequest, nil, "Invalid ID")
	}

	reqLog, err := Get(id)
	if err != nil {
		return http.ErrorResponse(ctx, fiber.StatusNotFound, nil, "Not Found")
	}

	return http.SuccessResponse(ctx, fiber.StatusOK, reqLog)
}

func ListHandler(ctx *fiber.Ctx) error {
	limit := ctx.Query("limit", "10")
	cursor := ctx.Query("cursor", "0")
	reverse := ctx.Query("reverse", "false")

	cursorInt, err := strconv.Atoi(cursor)
	if err != nil {
		return http.ErrorResponse(ctx, fiber.StatusBadRequest, nil, "Invalid cursor")
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return http.ErrorResponse(ctx, fiber.StatusBadRequest, nil, "Invalid limit")
	}

	reverseBool, err := strconv.ParseBool(reverse)
	if err != nil {
		return http.ErrorResponse(ctx, fiber.StatusBadRequest, nil, "Invalid reverse")
	}

	args := ctx.Context().QueryArgs()

	reqList, err := List(cursorInt, limitInt, reverseBool, args)
	if err != nil {
		return http.ErrorResponse(ctx, fiber.StatusNotFound, nil, "Not Found")
	}

	return http.SuccessResponse(ctx, fiber.StatusOK, reqList)
}
