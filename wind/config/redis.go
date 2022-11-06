package config

import (
	"encoding/json"
	"fmt"
)

type RedisClientOption struct {
	Addr         string `json:"addr"`
	Password     string `json:"password"`
	Db           int    `json:"db"`
	PoolSize     int    `json:"poolsize"`
	MinIdleConns int    `json:"minidleconns"`
}

func GetRedisClientOption(instance string) (RedisClientOption, error) {
	option := _viper.GetStringMap(fmt.Sprintf("redis.%s", instance))
	bytes, err := json.Marshal(option)
	if err != nil {
		return RedisClientOption{}, err
	}
	opt := RedisClientOption{}
	err = json.Unmarshal(bytes, &opt)
	if err != nil {
		return RedisClientOption{}, err
	}
	return opt, nil
}
