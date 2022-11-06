package trace

import (
    "context"
    "fmt"
    "github.com/go-redis/redis/v8"
    "github.com/opentracing/opentracing-go"
    "github.com/opentracing/opentracing-go/ext"
    "github.com/opentracing/opentracing-go/log"
    "mond/wind/utils/constant"
)

func RedisProcessBeforeMiddleware(ctx context.Context, method string, key string) (context.Context, error) {
    tracer := GetTracer()
    var err error
    ctx, span, _ := tracer.StartRedisClientSpanFromContext(ctx, fmt.Sprintf("redis.%v", method))
    span.SetTag("key", key)
    ctx = context.WithValue(ctx, constant.RedisCmdKey, span)
    return ctx, err

}

func RedisProcessAfterMiddleware(ctx context.Context, err error) error {
    spanInter := ctx.Value(constant.RedisCmdKey)
    if spanInter == nil {
        return nil
    }
    span, ok := spanInter.(opentracing.Span)
    if ok {
        if err != nil && err != redis.Nil {
            ext.Error.Set(span, true)
            span.LogFields(
                log.String("event", "error"),
                log.String("stack", err.Error()),
            )
        }
        span.Finish()
    }
    return nil
}
