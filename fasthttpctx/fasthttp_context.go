/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package fasthttpctx

import (
	"context"
	"github.com/hopeio/context/reqctx"
	stringsi "github.com/hopeio/utils/strings"
	"github.com/valyala/fasthttp"
	"net/http"
)

type RequestCtx struct {
	*fasthttp.RequestCtx
}

func (ctx RequestCtx) SetHeaders(md http.Header) {
	for k, v := range md {
		for _, vv := range v {
			ctx.RequestCtx.Response.Header.Add(k, vv)
		}
	}
}

func (ctx RequestCtx) SetHeader(k, v string) {
	ctx.RequestCtx.Response.Header.Set(k, v)
}

func (ctx RequestCtx) AddHeader(k, v string) {
	ctx.RequestCtx.Response.Header.Add(k, v)
}

func (ctx RequestCtx) GetHeader(k string) string {
	return stringsi.BytesToString(ctx.RequestCtx.Request.Header.Peek(k))
}

type Context = reqctx.Context[RequestCtx]

func FromContextValue(ctx context.Context) *Context {
	return reqctx.FromContextValue[RequestCtx](ctx)
}

func FromRequest(req *fasthttp.RequestCtx) *Context {
	r := &req.Request

	var ctx context.Context
	if r != nil {
		ctx = req
	}

	ctxi := reqctx.New[RequestCtx](ctx, RequestCtx{req})

	return ctxi
}
