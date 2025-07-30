/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package httpctx

import (
	"context"
	"github.com/hopeio/context/reqctx"
	httpi "github.com/hopeio/gox/net/http"
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

func (ctx RequestCtx) RequestContext() context.Context {
	return ctx.Request.Context()
}

type Context = reqctx.Context[RequestCtx]

func FromContext(ctx context.Context) (*Context, bool) {
	return reqctx.FromContext[RequestCtx](ctx)
}

func FromRequest(req RequestCtx) *Context {
	return reqctx.New[RequestCtx](req)
}
