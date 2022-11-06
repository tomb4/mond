package hook

import (
	"context"

	"google.golang.org/grpc"
)

type FrameStartHook interface {
	GrpcStartHook(s *grpc.Server) error

	ConfigChangeHook(conf map[string]interface{})

	ResourceInitHook(ctx context.Context) error

	AppStopHook()
}
