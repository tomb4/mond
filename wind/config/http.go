package config

import (
    "encoding/json"
    "fmt"
)

type HttpClientOption struct {
    Addr    string `json:"addr"`
    MaxConn int    `json:"maxconn"`
    Timeout int32  `json:"timeout"`
}

func GetHttpClientOption(instance string) (HttpClientOption, error) {
    option := _viper.GetStringMap(fmt.Sprintf("http.%s", instance))
    bytes, err := json.Marshal(option)
    if err != nil {
        return HttpClientOption{}, err
    }
    opt := HttpClientOption{}
    err = json.Unmarshal(bytes, &opt)
    if err != nil {
        return HttpClientOption{}, err
    }
    return opt, nil
}
