package config

import (
	"encoding/json"
	"fmt"
)

type GrpcClientOption struct {
	Scheme   string `json:"scheme"`
	Endpoint string `json:"endpoint"`
	Timeout  int32  `json:"timeout"` //TODO:  客户端按实例设置超时
	OpenLog  bool   `json:"openlog"`
}

func GetGrpcClientOption(appId string) GrpcClientOption {
	option := _viper.GetStringMap(fmt.Sprintf("grpcClient.%s", appId))
	opt := GrpcClientOption{}
	bytes, _ := json.Marshal(option)
	json.Unmarshal(bytes, &opt)
	if option["openlog"] == nil {
		opt.OpenLog = true
	}
	return opt
}
