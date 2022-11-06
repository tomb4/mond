package trace

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"mond/wind/config"
	"mond/wind/logger"
	"mond/wind/utils"
	"time"
)

func GrpcServerMiddleware() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	t := GetTracer()
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx, span, _ := t.StartGrpcServerSpanFromContext(ctx, info.FullMethod)
		defer span.Finish()
		resp, err = handler(ctx, req)
		if err != nil {
			ext.Error.Set(span, true)
			span.LogFields(
				log.String("event", "error"),
				log.String("stack", err.Error()),
				log.Object("metaErr", err),
				log.Object("req", utils.StructToJson(req)),
				log.Object("resp", utils.StructToJson(resp)),
			)
		}
		return resp, err
	}
}

func GrpcClientMiddleware() func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	tracer := GetTracer()
	timeoutConf := config.GetMethodTimeoutConfig()
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var cancel context.CancelFunc
		//fmt.Println(fmt.Sprintf("client.%s", method))
		timeout := timeoutConf.GetTimeout(fmt.Sprintf("client.%s", method))
		ctx, cancel = context.WithTimeout(ctx, time.Millisecond*timeout)
		defer cancel()
		ctx, span, err := tracer.StartGRpcClientSpanFromContext(ctx, method)
		if err != nil {
			logger.GetLogger().Error(ctx, "StartGRpcClientSpanFromContext", zap.Any("err", err))
		} else {
			defer span.Finish()
		}
		err = invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			ext.Error.Set(span, true)
			span.LogFields(
				log.String("event", "error"),
				log.String("stack", err.Error()),
				log.Object("req", utils.StructToJson(req)),
				log.Object("resp", utils.StructToJson(reply)),
			)
			return err
		}
		return nil
		
	}
}
