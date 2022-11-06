package merr

import (
	"context"
	"google.golang.org/grpc"
)

func GrpcServerMiddleware() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		if err != nil {
			//设置业务错误可以grpc透传
			e := ParseErrToMetaErr(err)
			err = ParseMetaErrToStatusErr(e)
		}
		return resp, err
	}
}

func GrpcClientMiddleware(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		//设置业务错误可以grpc透传
		err = ParseErrToMetaErr(err)
	}
	return err
}
