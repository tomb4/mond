package handler

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	dynamicConfig "mond/service/biz.demo/infra/config"
	"mond/service/biz.demo/proto"
	hook2 "mond/wind/hook"
	"mond/wind/logger"
)

type hook struct {
	h *BizdemoService
}

var (
	_handler *BizdemoService
)

func NewHook() hook2.FrameStartHook {
	_handler = NewBizdemoService()
	return &hook{
		h: _handler,
	}
}

// 框架基础组件初始化完成后，加载应用资源
func (h *hook) ResourceInitHook(ctx context.Context) error {
	err := h.h.ResourceInit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (h *hook) GrpcStartHook(s *grpc.Server) error {
	Bizdemo.RegisterBizdemoServiceServer(s, Bizdemo.NewBizdemoServiceServerImpl(h.h))

	return nil
}

func (h *hook) ConfigChangeHook(conf map[string]interface{}) {
	c := &dynamicConfig.DynamicConfig{}
	bytes, err := json.Marshal(conf)
	if err != nil {
		logger.GetLogger().Error(context.TODO(), "ConfigChangeHook失败", zap.Any("err", err))
		return
	}
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		logger.GetLogger().Error(context.TODO(), "ConfigChangeHook失败2", zap.Any("err", err))
		return
	}
	dynamicConfig.SetDynamicConfig(c)
}
func (h *hook) AppStopHook() {

}
