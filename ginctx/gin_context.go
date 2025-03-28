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
	httpi "github.com/hopeio/utils/net/http"
)

type RequestCtx struct {
	*gin.Context
}

func (ctx RequestCtx) RequestHeader() httpi.Header {
	return httpi.HttpHeader(ctx.Request.Header)
}

func (ctx RequestCtx) ResponseHeader() httpi.Header {
	return httpi.HttpHeader(ctx.Writer.Header())
}

func (ctx RequestCtx) RequestContext() context.Context {
	return ctx.Request.Context()
}

func (ctx RequestCtx) ToHttpReqCtx() httpctx.RequestCtx {
	return httpctx.RequestCtx{Request: ctx.Request, Response: ctx.Writer}
}

type Context = reqctx.Context[RequestCtx]

func FromRequest(req *gin.Context) *Context {
	return reqctx.New[RequestCtx](RequestCtx{req})
}
