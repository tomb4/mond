package config

import (
	"encoding/json"
	"errors"
)

type TraceConfig struct {
	Endpoint string `json:"endpoint"`
	Protocol string `json:"protocol"`
}

func GetTraceConfig() (TraceConfig, error) {
	option := _viper.GetStringMap("trace")
	bytes, err := json.Marshal(option)
	if err != nil {
		return TraceConfig{}, err
	}
	opt := TraceConfig{}
	err = json.Unmarshal(bytes, &opt)
	if err != nil {
		return TraceConfig{}, err
	}
	if opt.Endpoint == "" {
		return TraceConfig{}, errors.New("endpoint not found")
	}
	return opt, nil
}
