package config

import (
	"encoding/json"
	"fmt"
)

type MongoClientOption struct {
	Uri               string `json:"uri"`
	HeartbeatInterval int32
	MaxConnIdleTime   int32
	MaxPoolSize       int32
	MinPoolSize       int32
	ReadPreference    string
	ReadConcern       string
	WriteConcern      string
	MaxStaleness      int32
}

func GetMongoClientOption(instance string) (MongoClientOption, error) {
	option := _viper.GetStringMap(fmt.Sprintf("mongo.%s", instance))
	bytes, err := json.Marshal(option)
	if err != nil {
		return MongoClientOption{}, err
	}
	opt := MongoClientOption{}
	err = json.Unmarshal(bytes, &opt)
	if err != nil {
		return MongoClientOption{}, err
	}
	return opt, nil
}
