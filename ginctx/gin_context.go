/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package ginctx

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/hopeio/context/httpctx"
	"github.com/hopeio/context/reqctx"
	"net/http"
)

type RequestCtx struct {
	*gin.Context
}

func (ctx RequestCtx) SetHeaders(md http.Header) {
	header := ctx.Writer.Header()
	for k, v := range md {
		header[k] = v
	}
}

func (ctx RequestCtx) SetHeader(k, v string) {
	ctx.Writer.Header().Set(k, v)
}

func (ctx RequestCtx) AddHeader(k, v string) {
	ctx.Writer.Header().Add(k, v)
}

func (ctx RequestCtx) GetHeader(k string) string {
	return ctx.Request.Header.Get(k)
}

func (ctx RequestCtx) ToHttpReqCtx() httpctx.RequestCtx {
	return httpctx.RequestCtx{Request: ctx.Request, Response: ctx.Writer}
}

type Context = reqctx.Context[RequestCtx]

func FromContextValue(ctx context.Context) *Context {
	return reqctx.FromContextValue[RequestCtx](ctx)
}

func FromRequest(req *gin.Context) *Context {
	r := req.Request

	var ctx context.Context
	if r != nil {
		ctx = r.Context()
	}

	ctxi := reqctx.New[RequestCtx](ctx, RequestCtx{req})
	return ctxi
}
