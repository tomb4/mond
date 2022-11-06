package config

import (
	"encoding/json"
	"fmt"
	"mond/wind/utils"
	"strings"
	"time"
)

type MethodTimeoutConfig struct {
	timeout map[string]int32
}

func GetMethodTimeoutConfig() *MethodTimeoutConfig {
	option := _viper.GetStringMap(fmt.Sprintf("methodTimeout"))
	bytes, err := json.Marshal(option)
	utils.MustNil(err)
	//fmt.Println(string(bytes))
	opt := MethodTimeoutConfig{timeout: map[string]int32{}}
	err = json.Unmarshal(bytes, &opt.timeout)
	//fmt.Println(opt)
	utils.MustNil(err)
	return &opt
}

var (
	defaultTimeout time.Duration = 5000 * time.Millisecond //ms
)

func (m *MethodTimeoutConfig) GetTimeout(method string) time.Duration {
	v := m.timeout[strings.ToLower(method)]
	if v == 0 {
		return defaultTimeout
	}
	return time.Millisecond * time.Duration(v)
}
