/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package httpctx

import (
	"context"
	"github.com/hopeio/context/reqctx"
	httpi "github.com/hopeio/utils/net/http"
	"net/http"
)

type RequestCtx struct {
	Request  *http.Request
	Response http.ResponseWriter
}

func (ctx RequestCtx) RequestHeader() httpi.Header {
	return httpi.HttpHeader(ctx.Request.Header)
}

func (ctx RequestCtx) ResponseHeader() httpi.Header {
	return httpi.HttpHeader(ctx.Response.Header())
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
