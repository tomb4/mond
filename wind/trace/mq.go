package trace

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"mond/wind/mq/rabbit/function"
	"mond/wind/utils/constant"
)

func ConsumerMiddleware(f function.ConsumerFunc) function.ConsumerFunc {
	t := GetTracer()
	return func(ctx context.Context, queue string) error {
		name := fmt.Sprintf("consumer.%s", queue)
		if ctx.Value(constant.AsyncMethodCtxKey) != nil {
			name = fmt.Sprintf("consumer.async.%s", ctx.Value(constant.AsyncMethodCtxKey))
		}
		ctx, span, _ := t.StartConsumerSpanFromContext(ctx, name)
		defer span.Finish()
		if ctx.Value(constant.AsyncMethodCtxKey) != nil {
			span.SetTag("async", "true")
		}
		span.SetTag("exchange", ctx.Value(constant.MqExchangeKey))
		span.SetTag("queue", queue)
		err := f(ctx, queue)
		if err != nil {
			ext.Error.Set(span, true)
			span.LogFields(
				log.String("event", "error"),
				log.String("stack", err.Error()),
				log.Object("metaErr", err),
			)
		}
		return err
	}
}

func PublishMiddleware(ctx context.Context, name string) (context.Context, opentracing.Span, error) {
	return GetTracer().StartPublishSpanFromContext(ctx, name)
}
