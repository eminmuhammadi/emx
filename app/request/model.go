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
package request

import (
	"github.com/elazarl/goproxy"
	"github.com/eminmuhammadi/emx/pkg/sql"
	"github.com/eminmuhammadi/emx/pkg/util"
	"gorm.io/gorm"
)

type Request struct {
	sql.BaseModel
	SessionID        string `gorm:"index,unique,not null" json:"session_id"`
	Method           string `json:"method" gorm:"text"`
	URL              URL    `json:"url" gorm:"embedded;embeddedPrefix:url_"`
	Proto            string `json:"proto" gorm:"text"`
	ProtoMajor       int    `json:"proto_major" gorm:"numeric"`
	ProtoMinor       int    `json:"proto_minor" gorm:"numeric"`
	Header           string `json:"header" gorm:"text"`
	Body             string `json:"body" gorm:"text"`
	ContentLength    int64  `json:"content_length" gorm:"numeric"`
	TransferEncoding string `json:"transfer_encoding" gorm:"text"`
	Host             string `json:"host" gorm:"text"`
	Trailer          string `json:"trailer" gorm:"text"`
	RemoteAddr       string `json:"remote_addr" gorm:"text"`
	RequestURI       string `json:"request_uri" gorm:"text"`
}

type URL struct {
	Scheme      string    `json:"scheme" gorm:"text"`
	Opaque      string    `json:"opaque" gorm:"text"`
	User        *Userinfo `json:"user" gorm:"embedded;embeddedPrefix:user_"`
	Host        string    `json:"host" gorm:"text"`
	Path        string    `json:"path" gorm:"text"`
	RawPath     string    `json:"raw_path" gorm:"text"`
	OmitHost    bool      `json:"omit_host" gorm:"bool"`
	ForceQuery  bool      `json:"force_query" gorm:"bool"`
	RawQuery    string    `json:"raw_query" gorm:"text"`
	Fragment    string    `json:"fragment" gorm:"text"`
	RawFragment string    `json:"raw_fragment" gorm:"text"`
}

type Userinfo struct {
	Username string `json:"username" gorm:"text"`
	Password string `json:"password" gorm:"text"`
}

func (r *Request) BeforeCreate(*gorm.DB) (err error) {
	r.CreatedAt = sql.TimeNow()
	r.UpdatedAt = sql.TimeNow()
	r.IsDeleted = 0

	return nil
}

func (r *Request) BeforeUpdate(*gorm.DB) (err error) {
	r.UpdatedAt = sql.TimeNow()

	return nil
}

func (r *Request) BeforeDelete(*gorm.DB) (err error) {
	r.UpdatedAt = sql.TimeNow()
	r.IsDeleted = 1

	return nil
}

func CreateRequestModel(ctx *goproxy.ProxyCtx) Request {
	return Request{
		SessionID: "", // updated on pkg/proxy/server.go
		Method:    ctx.Req.Method,
		URL: URL{
			Scheme: ctx.Req.URL.Scheme,
			Opaque: ctx.Req.URL.Opaque,
			Host:   ctx.Req.URL.Host,
			User: &Userinfo{
				Username: ctx.Req.URL.User.Username(),
				Password: func() string {
					if password, ok := ctx.Req.URL.User.Password(); ok {
						return password
					}

					return ""
				}(),
			},
			Path:        ctx.Req.URL.Path,
			RawPath:     ctx.Req.URL.RawPath,
			OmitHost:    ctx.Req.URL.OmitHost,
			ForceQuery:  ctx.Req.URL.ForceQuery,
			RawQuery:    ctx.Req.URL.RawQuery,
			Fragment:    ctx.Req.URL.Fragment,
			RawFragment: ctx.Req.URL.RawFragment,
		},
		Proto:            ctx.Req.Proto,
		ProtoMajor:       ctx.Req.ProtoMajor,
		ProtoMinor:       ctx.Req.ProtoMinor,
		Header:           util.NormalizeHeaders(ctx.Req.Header),
		ContentLength:    ctx.Req.ContentLength,
		TransferEncoding: util.StringifyArray(ctx.Req.TransferEncoding),
		Host:             ctx.Req.Host,
		Trailer:          util.NormalizeHeaders(ctx.Req.Trailer),
		RemoteAddr:       ctx.Req.RemoteAddr,
		RequestURI:       ctx.Req.RequestURI,
		Body:             "", // updated on pkg/proxy/server.go
	}
}

func (r *Request) Set(key, value string) Request {
	switch key {
	case "session_id":
		r.SessionID = value
	case "method":
		r.Method = value
	case "proto":
		r.Proto = value
	case "proto_major":
		r.ProtoMajor, _ = util.StringToInt(value)
	case "proto_minor":
		r.ProtoMinor, _ = util.StringToInt(value)
	case "header":
		r.Header = value
	case "body":
		r.Body = value
	case "content_length":
		r.ContentLength, _ = util.StringToInt64(value)
	case "transfer_encoding":
		r.TransferEncoding = value
	case "host":
		r.Host = value
	case "trailer":
		r.Trailer = value
	case "remote_addr":
		r.RemoteAddr = value
	case "request_uri":
		r.RequestURI = value
	case "url_scheme":
		r.URL.Scheme = value
	case "url_opaque":
		r.URL.Opaque = value
	case "url_host":
		r.URL.Host = value
	case "url_path":
		r.URL.Path = value
	case "url_raw_path":
		r.URL.RawPath = value
	case "url_omit_host":
		r.URL.OmitHost, _ = util.StringToBool(value)
	case "url_force_query":
		r.URL.ForceQuery, _ = util.StringToBool(value)
	case "url_raw_query":
		r.URL.RawQuery = value
	case "url_fragment":
		r.URL.Fragment = value
	case "url_raw_fragment":
		r.URL.RawFragment = value
	case "url_user_username":
		r.URL.User.Username = value
	case "url_user_password":
		r.URL.User.Password = value
	}

	return *r
}
