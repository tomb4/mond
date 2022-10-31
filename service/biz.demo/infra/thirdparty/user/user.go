package user

import (
	"context"
	"mond/service/chimay/user"
	"strconv"
)

type BizUserService struct {
	client user.GrpcUserServiceClient
}

func NewBizUserService() (*BizUserService, error) {
	userClient, err := user.GetGrpcUserServiceClient()
	if err != nil {
		return nil, err
	}
	return &BizUserService{client: userClient}, nil
}

func (m *BizUserService) CheckSSO(ctx context.Context, token string) (int32, error) {
	resp, err := m.client.CheckSso(ctx, &user.CheckSsoRequest{Token: token})
	if err != nil {
		return 0, err
	}
	uid, _ := strconv.ParseInt(resp.Uid, 10, 64)
	return int32(uid), nil
}
