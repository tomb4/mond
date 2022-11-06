// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package Bizdemo

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// BizdemoServiceClient is the client API for BizdemoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BizdemoServiceClient interface {
	Ping(ctx context.Context, in *PingReq, opts ...grpc.CallOption) (*PingResp, error)
}

type bizdemoServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBizdemoServiceClient(cc grpc.ClientConnInterface) BizdemoServiceClient {
	return &bizdemoServiceClient{cc}
}

func (c *bizdemoServiceClient) Ping(ctx context.Context, in *PingReq, opts ...grpc.CallOption) (*PingResp, error) {
	out := new(PingResp)
	err := c.cc.Invoke(ctx, "/Bizdemo.BizdemoService/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BizdemoServiceServer is the server API for BizdemoService service.
// All implementations must embed UnimplementedBizdemoServiceServer
// for forward compatibility
type BizdemoServiceServer interface {
	Ping(context.Context, *PingReq) (*PingResp, error)
	mustEmbedUnimplementedBizdemoServiceServer()
}

// UnimplementedBizdemoServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBizdemoServiceServer struct {
}

func (UnimplementedBizdemoServiceServer) Ping(context.Context, *PingReq) (*PingResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedBizdemoServiceServer) mustEmbedUnimplementedBizdemoServiceServer() {}

// UnsafeBizdemoServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BizdemoServiceServer will
// result in compilation errors.
type UnsafeBizdemoServiceServer interface {
	mustEmbedUnimplementedBizdemoServiceServer()
}

func RegisterBizdemoServiceServer(s grpc.ServiceRegistrar, srv BizdemoServiceServer) {
	s.RegisterService(&BizdemoService_ServiceDesc, srv)
}

func _BizdemoService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BizdemoServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Bizdemo.BizdemoService/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BizdemoServiceServer).Ping(ctx, req.(*PingReq))
	}
	return interceptor(ctx, in, info, handler)
}

// BizdemoService_ServiceDesc is the grpc.ServiceDesc for BizdemoService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BizdemoService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Bizdemo.BizdemoService",
	HandlerType: (*BizdemoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _BizdemoService_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/Bizdemo.proto",
}