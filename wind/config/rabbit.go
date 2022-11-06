package config

import (
	"encoding/json"
	"fmt"
)

type RabbitClientOption struct {
	Url      string `json:"url"`
	VHost    string `json:"vhost"`
	Consumer struct {
		MaxChannel    int32 `json:"maxchannel"`
		PrefetchCount int32 `json:"prefetchcount"`
	} `json:"consumer"`
	Producer struct {
		MaxChannel int32 `json:"maxchannel"`
	} `json:"producer"`
}

func GetRabbitClientOption(instance string) (RabbitClientOption, error) {
	option := _viper.GetStringMap(fmt.Sprintf("rabbit.%s", instance))
	bytes, err := json.Marshal(option)
	if err != nil {
		return RabbitClientOption{}, err
	}
	opt := RabbitClientOption{}
	err = json.Unmarshal(bytes, &opt)
	if err != nil {
		return RabbitClientOption{}, err
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
