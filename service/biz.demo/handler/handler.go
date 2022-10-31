package handler

import (
	"context"
	"mond/service/biz.demo/app"
	"mond/service/biz.demo/proto"
	"mond/wind/logger"
)

type BizdemoService struct {
	_log logger.Logger
	app  *app.App
}

func (m *BizdemoService) Ping(ctx context.Context, n *Bizdemo.PingReq) (*Bizdemo.PingResp, error) {
	panic("implement me")
}

func NewBizdemoService() *BizdemoService {
	return &BizdemoService{_log: logger.GetLogger()}
}
