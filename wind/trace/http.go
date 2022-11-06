package trace

import (
    "context"
    "fmt"
    "github.com/opentracing/opentracing-go/ext"
    "github.com/opentracing/opentracing-go/log"
    "mond/wind/utils/constant"
    "mond/wind/utils/endpoint"
    "net/http"
)

func HttpClientMiddleware(endpoint endpoint.Endpoint) endpoint.Endpoint {
    tracer := GetTracer()
    return func(ctx context.Context, req interface{}) (interface{}, error) {
        host := ctx.Value(constant.HttpClientHost)
        path := ctx.Value(constant.HttpClientPath)
        method := ctx.Value(constant.HttpClientMethod)
        ctx, span, _ := tracer.StartHttpClientSpanFromContext(ctx, fmt.Sprintf("http.%v_%v", path, method))
        defer span.Finish()
        span.SetTag("host", host)
        resp, err := endpoint(ctx, req)
        //查不到数据这种数据正常业务报错，不记录
        statusCode := 0
        if resp != nil {
            r, ok := resp.(*http.Response)
            if ok && r != nil {
                span.SetTag("status_code", r.StatusCode)
                statusCode = r.StatusCode
            }
        }
        if err != nil || statusCode >= 400 {
            ext.Error.Set(span, true)
            span.LogFields(
                log.String("event", "error"),
            )
            if err != nil {
                span.LogFields(
                    log.String("stack", err.Error()),
                )
            }
        }
        return resp, err
    }
}
