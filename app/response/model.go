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
	"github.com/elazarl/goproxy"
	"github.com/eminmuhammadi/emx/pkg/sql"
	"github.com/eminmuhammadi/emx/pkg/util"
	"gorm.io/gorm"
)

type Response struct {
	sql.BaseModel
	SessionID        string `gorm:"index,unique,not null" json:"session_id"`
	Status           string `json:"status" gorm:"text"`
	StatusCode       int    `json:"statusCode" gorm:"numeric"`
	Proto            string `json:"proto" gorm:"text"`
	ProtoMajor       int    `json:"protoMajor" gorm:"numeric"`
	ProtoMinor       int    `json:"protoMinor" gorm:"numeric"`
	Header           string `json:"header" gorm:"text"`
	Body             string `json:"body" gorm:"text"`
	ContentLength    int64  `json:"contentLength" gorm:"numeric"`
	TransferEncoding string `json:"transferEncoding" gorm:"text"`
	Trailer          string `json:"trailer" gorm:"text"`
}

func (r *Response) BeforeCreate(*gorm.DB) (err error) {
	r.CreatedAt = sql.TimeNow()
	r.UpdatedAt = sql.TimeNow()
	r.IsDeleted = 0

	return nil
}

func (r *Response) BeforeUpdate(*gorm.DB) (err error) {
	r.UpdatedAt = sql.TimeNow()

	return nil
}

func (r *Response) BeforeDelete(*gorm.DB) (err error) {
	r.UpdatedAt = sql.TimeNow()
	r.IsDeleted = 1

	return nil
}

func CreateResponseModel(ctx *goproxy.ProxyCtx) Response {
	return Response{
		SessionID:        "", // updated on pkg/proxy/server.go
		Status:           ctx.Resp.Status,
		StatusCode:       ctx.Resp.StatusCode,
		Proto:            ctx.Resp.Proto,
		ProtoMajor:       ctx.Resp.ProtoMajor,
		ProtoMinor:       ctx.Resp.ProtoMinor,
		Header:           util.NormalizeHeaders(ctx.Resp.Header),
		ContentLength:    ctx.Resp.ContentLength,
		TransferEncoding: util.StringifyArray(ctx.Resp.TransferEncoding),
		Trailer:          util.NormalizeHeaders(ctx.Resp.Trailer),
		Body:             "", // updated on pkg/proxy/server.go
	}
}
