/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package fasthttpctx

import (
	"context"
	"github.com/hopeio/context/reqctx"
	httpi "github.com/hopeio/utils/net/http"
	fiberi "github.com/hopeio/utils/net/http/fiber"
	stringsi "github.com/hopeio/utils/strings"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc/metadata"
)

type Context = reqctx.Context[*fasthttp.RequestCtx]

func FromContextValue(ctx context.Context) *Context {
	return reqctx.FromContextValue[*fasthttp.RequestCtx](ctx)
}

func FromRequest(req *fasthttp.RequestCtx) *Context {
	r := &req.Request

	var ctx context.Context
	if r != nil {
		ctx = req
	}

	ctxi := reqctx.New[*fasthttp.RequestCtx](ctx, req)
	setWithReq(ctxi, r)
	return ctxi
}

func setWithReq(c *Context, r *fasthttp.Request) {
	c.Token = fiberi.GetToken(r)
	c.DeviceInfo = Device(&r.Header)
	c.Internal = stringsi.BytesToString(r.Header.Peek(httpi.HeaderGrpcInternal))
}

func Device(r *fasthttp.RequestHeader) *reqctx.DeviceInfo {
	return reqctx.Device(stringsi.BytesToString(r.Peek(httpi.HeaderDeviceInfo)),
		stringsi.BytesToString(r.Peek(httpi.HeaderArea)),
		stringsi.BytesToString(r.Peek(httpi.HeaderLocation)),
		stringsi.BytesToString(r.Peek(httpi.HeaderUserAgent)),
		stringsi.BytesToString(r.Peek(httpi.HeaderXForwardedFor)),
	)
}

type FastHttpContext Context

func (c *FastHttpContext) SetHeader(md metadata.MD) error {
	header := &c.ReqCtx.Response.Header
	for k, v := range md {
		if len(v) > 0 {
			header.Set(k, v[0])
		}
	}
	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SetHeader(md)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *FastHttpContext) SendHeader(md metadata.MD) error {
	header := &c.ReqCtx.Response.Header
	for k, v := range md {
		if len(v) > 0 {
			header.Set(k, v[0])
		}
	}
	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SendHeader(md)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *FastHttpContext) WriteHeader(k, v string) error {
	c.ReqCtx.Response.Header.Set(k, v)
	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SendHeader(metadata.MD{k: []string{v}})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *FastHttpContext) SetCookie(v string) error {
	c.ReqCtx.Response.Header.Set(httpi.HeaderSetCookie, v)
	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SendHeader(metadata.MD{httpi.HeaderSetCookie: []string{v}})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *FastHttpContext) SetTrailer(md metadata.MD) error {
	for k, v := range md {
		if len(v) > 0 {
			c.ReqCtx.Response.Header.Set(k, v[0])
		}
	}
	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SetTrailer(md)
		if err != nil {
			return err
		}
	}
	return nil
}
