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
	"fmt"
	"sort"
	"strings"

	"github.com/eminmuhammadi/emx/app/request"
	"github.com/eminmuhammadi/emx/app/response"
	"github.com/eminmuhammadi/emx/pkg/sql"
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

	// set order
	order := "asc"
	if reverse {
		order = "desc"
	}

	// convert argsStr to urlQuery
	argsStr := args.String()
	query := fmt.Sprintf("id > %d", cursor)
	if argsStr != "" {
		splitedArgs := strings.Split(argsStr, "&")
		for _, arg := range splitedArgs {
			splitedArg := strings.Split(arg, "=")
			if splitedArg[0] == "cursor" || splitedArg[0] == "limit" || splitedArg[0] == "reverse" {
				continue
			}

			if len(splitedArg) == 2 {
				query += " AND " + splitedArg[0] + " = '" + splitedArg[1] + "'"
			}
		}
	}

	if err := sql.Database.
		Where(query).
		Order(fmt.Sprintf("id %s", order)).
		Limit(limit).
		Find(&reqList).Error; err != nil {
		return nil, err
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
