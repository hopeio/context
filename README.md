# context

一个轻量却强大的上下文管理器,一个请求会生成一个context，贯穿整个请求，context记录原始请求上下文，请求时间，客户端信息，权限校验信息，及负责判断是否内部调用，
及附带唯一traceId的日志记录器

支持http及fasthttp,并支持自定义的请求类型
![context](_assets/context.webp)