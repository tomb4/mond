package mctx

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"mond/wind/utils/constant"
)

func WithTraceId(ctx context.Context, span opentracing.Span) context.Context {
	if ctx.Value(constant.TraceIdKey) != nil {
		return ctx
	}
	s, ok := span.Context().(jaeger.SpanContext)
	if ok {
		ctx = context.WithValue(ctx, constant.TraceIdKey, s.TraceID().String())
	}
	return ctx
}

func GetTraceId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	traceId := ctx.Value(constant.TraceIdKey)
	if traceId == nil {
		return ""
	} else {
		return traceId.(string)
	}
}
