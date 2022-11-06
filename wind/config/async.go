package config

import (
	"encoding/json"
	"fmt"
)

type AsyncOption struct {
	Type     string `json:"type"`
	Instance string `json:"instance"`
	Consumer struct {
		MaxChannel    int32 `json:"maxchannel"`
		PrefetchCount int32 `json:"prefetchcount"`
	} `json:"consumer"`
	Producer struct {
		MaxChannel int32 `json:"maxchannel"`
	} `json:"producer"`
}

func GetAsyncOption() (AsyncOption, error) {
	option := _viper.GetStringMap(fmt.Sprintf("async"))
	bytes, err := json.Marshal(option)
	if err != nil {
		return AsyncOption{}, err
	}
	opt := AsyncOption{}
	err = json.Unmarshal(bytes, &opt)
	if err != nil {
		return AsyncOption{}, err
	}
	if opt.Consumer.MaxChannel == 0 {
		opt.Consumer.MaxChannel = 1
	}
	if opt.Consumer.PrefetchCount == 0 {
		opt.Consumer.PrefetchCount = 10
	}
	if opt.Producer.MaxChannel == 0 {
		opt.Producer.MaxChannel = 1
	}
	return opt, nil
}
