package cmd

import "text/template"

var infraConfigTemplate, _ = template.New("").Parse(`
package dynamicConfig

var (
	dynamicConfigInstance *DynamicConfig = &DynamicConfig{}
)

func GetDynamicConfig() *DynamicConfig {
	return dynamicConfigInstance
}

func SetDynamicConfig(instance *DynamicConfig) {
	dynamicConfigInstance = instance
}

type DynamicConfig struct {
	Id        string ` + "`" + `json:"id"` + "`" + `
}
`)

var infraErrTemplate, _ = template.New("").Parse(`
package xerr

import merr "mond/wind/err"

var (
	LoginErr             = merr.NewError({{.Port}}001, "登录失败，token验证不合法", merr.Abnormal)
)

`)

var infraThirdPartyTemplate, _ = template.New("").Parse(`
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

`)
