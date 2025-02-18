/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package fiberctx

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"github.com/hopeio/context/fasthttpctx"
	"github.com/hopeio/context/reqctx"
	"net/http"
)

type RequestCtx struct {
	fiber.Ctx
}

func (ctx RequestCtx) SetHeaders(md http.Header) {
	for k, v := range md {
		for _, vv := range v {
			ctx.Set(k, vv)
		}
	}
}

func (ctx RequestCtx) SetHeader(k, v string) {
	ctx.Set(k, v)
}

func (ctx RequestCtx) AddHeader(k, v string) {
	ctx.Response().Header.Add(k, v)
}

func (ctx RequestCtx) GetHeader(k string) string {
	return ctx.Get(k)
}

type Context = reqctx.Context[RequestCtx]

func FromContextValue(ctx context.Context) *Context {
	return reqctx.FromContextValue[RequestCtx](ctx)
}

func (ctx RequestCtx) ConvertToFastHttpCtx() *fasthttpctx.RequestCtx {
	return &fasthttpctx.RequestCtx{}
}

func FromRequest(req fiber.Ctx) *Context {
	r := req.Request

	var ctx context.Context
	if r != nil {
		ctx = req.Context()
	}
	ctxi := reqctx.New[RequestCtx](ctx, RequestCtx{req})
	return ctxi
}
