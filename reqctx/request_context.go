/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package reqctx

import (
	"context"
	"github.com/google/uuid"
	context2 "github.com/hopeio/context"
	httpi "github.com/hopeio/utils/net/http"
	"github.com/hopeio/utils/net/http/consts"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
	"sync"
)

func GetPool[REQ ReqCtx]() sync.Pool {
	return sync.Pool{New: func() any {
		return new(Context[REQ])
	}}
}

type ReqValue struct {
	Token string
	Auth
	device *DeviceInfo
	grpc.ServerTransportStream
	Internal string
	RequestAt
}

type ReqCtx interface {
	RequestHeader() httpi.Header
	ResponseHeader() httpi.Header
	RequestContext() context.Context
}

type Context[REQ ReqCtx] struct {
	context2.Context
	ReqValue
	ReqCtx REQ
}

func methodFamily(m string) string {
	m = strings.TrimPrefix(m, "/") // remove leading slash
	if i := strings.Index(m, "/"); i >= 0 {
		m = m[:i] // remove everything from second slash
	}
	return m
}

func (c *Context[REQ]) Wrapper() context.Context {
	return context.WithValue(c.Context.Base(), context2.WrapperKey(), c)
}

func (c *Context[REQ]) StartSpanX(name string, o ...trace.SpanStartOption) (*Context[REQ], trace.Span) {
	span := c.Context.StartSpan(name, o...)
	return c, span
}

func FromContextValue[REQ ReqCtx](ctx context.Context) (*Context[REQ], bool) {
	if ctx == nil {
		return nil, false
	}
	ctxi := ctx.Value(context2.WrapperKey())
	c, ok := ctxi.(*Context[REQ])
	if !ok {
		return nil, false
	}
	if c.ServerTransportStream == nil {
		c.ServerTransportStream = grpc.ServerTransportStreamFromContext(ctx)
	}
	c.SetBase(ctx)
	return c, ok
}

func New[REQ ReqCtx](req REQ) *Context[REQ] {
	ctx := req.RequestContext()
	c, ok := FromContextValue[REQ](ctx)
	if ok {
		return c
	}
	return &Context[REQ]{
		Context: *context2.New(ctx),
		ReqValue: ReqValue{
			RequestAt:             NewRequestAt(),
			ServerTransportStream: grpc.ServerTransportStreamFromContext(ctx),
			Internal:              req.RequestHeader().Get(consts.HeaderGrpcInternal),
			Token:                 GetToken(req),
		},
		ReqCtx: req,
	}
}

func (c *Context[REQ]) reset(ctx context.Context) *Context[REQ] {
	span := trace.SpanFromContext(ctx)
	traceId := span.SpanContext().TraceID().String()
	if traceId == "" {
		traceId = uuid.New().String()
	}
	c.SetBase(ctx)
	c.RequestAt = NewRequestAt()
	return c
}

func (c *Context[REQ]) Device() *DeviceInfo {
	if c.device == nil {
		header := c.ReqCtx.RequestHeader()
		c.device = Device(header.Get(consts.HeaderDeviceInfo),
			header.Get(consts.HeaderArea), header.Get(consts.HeaderLocation),
			header.Get(consts.HeaderUserAgent), header.Get(consts.HeaderXForwardedFor))
	}
	return c.device
}

func (c *Context[REQ]) Method() string {
	if c.ServerTransportStream != nil {
		return c.ServerTransportStream.Method()
	}
	return ""
}

func (c *Context[REQ]) SetHeader(md metadata.MD) error {
	header := c.ReqCtx.ResponseHeader()
	for k, v := range md {
		header.Set(k, v[0])
	}
	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SetHeader(md)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Context[REQ]) SendHeader(md metadata.MD) error {
	header := c.ReqCtx.ResponseHeader()
	for k, v := range md {
		header.Set(k, v[0])
	}
	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SendHeader(md)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Context[REQ]) SetCookie(v string) error {
	c.ReqCtx.ResponseHeader().Set(consts.HeaderSetCookie, v)
	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SendHeader(metadata.MD{consts.HeaderSetCookie: []string{v}})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Context[REQ]) SetTrailer(md metadata.MD) error {
	header := c.ReqCtx.ResponseHeader()
	for k, v := range md {
		header.Set(k, v[0])
	}
	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SetTrailer(md)
		if err != nil {
			return err
		}
	}
	return nil
}
