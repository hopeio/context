/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package httpctx

import (
	"context"
	"github.com/hopeio/context/reqctx"
	"net/http"
)

type RequestCtx struct {
	Request  *http.Request
	Response http.ResponseWriter
}

func (ctx RequestCtx) SetHeaders(md http.Header) {
	header := ctx.Response.Header()
	for k, v := range md {
		header[k] = v
	}
}

func (ctx RequestCtx) SetHeader(k, v string) {
	ctx.Response.Header().Set(k, v)
}

func (ctx RequestCtx) AddHeader(k, v string) {
	ctx.Response.Header().Add(k, v)
}

func (ctx RequestCtx) GetHeader(k string) string {
	return ctx.Request.Header.Get(k)
}

type Context = reqctx.Context[RequestCtx]

func FromContextValue(ctx context.Context) *Context {
	return reqctx.FromContextValue[RequestCtx](ctx)
}

func FromRequest(req RequestCtx) *Context {
	r := req.Request
	var ctx context.Context
	if r != nil {
		ctx = r.Context()
	}
	return reqctx.New[RequestCtx](ctx, req)
}
