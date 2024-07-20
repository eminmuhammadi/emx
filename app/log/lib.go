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
package log

import (
	"fmt"
	"sort"
	"strings"

	"github.com/eminmuhammadi/emx/app/request"
	"github.com/eminmuhammadi/emx/app/response"
	"github.com/eminmuhammadi/emx/pkg/sql"
	"github.com/eminmuhammadi/emx/pkg/util"
	"github.com/valyala/fasthttp"
)

type Log struct {
	request.Request
	Response response.Response `json:"response"`
}

func Get(id int64) (Log, error) {
	var reqLog Log

	req, err := request.Get(id)
	if err != nil {
		return reqLog, err
	}

	resp := response.GetBySessionID(req.SessionID)

	reqLog = Log{
		Request:  req,
		Response: resp,
	}

	return reqLog, nil
}

func List(cursor int, limit int, reverse bool, args *fasthttp.Args) ([]Log, error) {
	var reqList []request.Request
	var reqLogList []Log
	var reqDetails request.Request

	// set order
	order := "asc"
	if reverse {
		order = "desc"
	}

	// get args string
	argsStr := args.String()

	// Create db query chain
	db := sql.Database

	// set cursor
	db = db.Where("id > ?", cursor)

	if argsStr != "" {
		splitedArgs := strings.Split(argsStr, "&")
		for _, arg := range splitedArgs {
			splitedArg := strings.Split(arg, "=")
			if len(splitedArg) != 2 {
				continue
			}

			key := strings.TrimSpace(splitedArg[0])
			value := strings.TrimSpace(splitedArg[1])

			if key == "cursor" || key == "limit" || key == "reverse" || key == "id" {
				continue
			}

			value, err := util.NormalizeQueryValue(value)
			if err != nil {
				return []Log{}, err
			}

			reqDetails = reqDetails.Set(key, value)
		}
	}

	if err := db.
		Where(&reqDetails).
		Order(fmt.Sprintf("id %s", order)).
		Limit(limit).
		Find(&reqList).Error; err != nil {
		return []Log{}, err
	}

	sort.Slice(reqList, func(i, j int) bool {
		return reqList[i].ID > reqList[j].ID
	})

	for _, req := range reqList {
		reqLogList = append(reqLogList, Log{
			Request:  req,
			Response: response.GetBySessionID(req.SessionID),
		})
	}

	return reqLogList, nil
}
