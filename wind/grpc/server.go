package mgrpc

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"mond/wind/config"
	merr "mond/wind/err"
	"mond/wind/hook"
	"mond/wind/logger"
	"mond/wind/sentinel/breaker"
	"mond/wind/sentry"
	"mond/wind/trace"
	"net"
	"time"
)

type GrpcBase struct {
	server  *grpc.Server
	ErrChan chan error
}

func (m *GrpcBase) InitGrpcServer(ctx context.Context, frameHook hook.FrameStartHook, startChan chan int32) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GetInt32("port")))
	if err != nil {
		logger.GetLogger().Error(ctx, fmt.Sprintf("grpc Listen err"), zap.Any("err", err))
		m.ErrChan <- err
		return
	}
	ms := make([]grpc.UnaryServerInterceptor, 0)
	ms = append(ms,
		trace.GrpcServerMiddleware(),
		logger.GrpcServerMiddleware(),
		merr.GrpcServerMiddleware(),
		sentry.GrpcServerMiddleware(),
	)
	conf := config.GetSentinelOption()
	if conf.Breaker.GrpcServerOpen {
		ms = append(ms, breaker.MakeGrpcServerInterceptor())
	}
	// 实例化grpc服务端
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(ms...)),
		grpc.MaxConcurrentStreams(1000000),
		grpc.InitialWindowSize(1034*1024*1024),
		grpc.InitialConnWindowSize(1034*1024*1024),
	)
	m.server = s
	frameHook.GrpcStartHook(s)
	reflection.Register(s)
	logger.GetLogger().Info(ctx, fmt.Sprintf("grpc服务启动成功, port:%d", config.GetInt32("port")))
	go func() {
		time.Sleep(time.Millisecond * 10)
		startChan <- 1
	}()
	
	err = s.Serve(lis)
	if err != nil {
		logger.GetLogger().Error(ctx, fmt.Sprintf("grpc server err"), zap.Any("err", err))
		m.ErrChan <- err
		return
	}
	lis.Close()
	logger.GetLogger().Info(ctx, fmt.Sprintf("grpc服务停止成功"))
}

func (m *GrpcBase) GracefulStop() {
	if m.server != nil {
		logger.GetLogger().Info(context.TODO(), fmt.Sprintf("grpc优雅退出开始"))
		m.server.GracefulStop()
		logger.GetLogger().Info(context.TODO(), fmt.Sprintf("grpc优雅退出成功"))
	}
}
