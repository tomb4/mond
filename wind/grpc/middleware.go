package mgrpc

import (
	"google.golang.org/grpc"
	"mond/wind/config"
	merr "mond/wind/err"
	"mond/wind/logger"
	"mond/wind/sentinel/breaker"
	"mond/wind/trace"
)

func ClientMiddleware(opt config.GrpcClientOption) []grpc.UnaryClientInterceptor {
	item := []grpc.UnaryClientInterceptor{
		trace.GrpcClientMiddleware(),
	}
	if opt.OpenLog {
		item = append(item, logger.GrpcClientMiddleware)
	}
	item = append(item, merr.GrpcClientMiddleware)
	if config.GetSentinelOption().Breaker.GrpcClientOpen {
		item = append(item, breaker.MakeGRpcClientMiddleware())
	}
	
	return item
}
