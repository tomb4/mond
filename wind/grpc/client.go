package mgrpc

import "google.golang.org/grpc"

var (
	DefaultDialOptions []grpc.DialOption = []grpc.DialOption{
		grpc.WithInitialConnWindowSize(1024 * 1024 * 1024),
		grpc.WithInitialWindowSize(1024 * 1024 * 1024),
	}
)
