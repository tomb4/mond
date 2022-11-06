package config

import (
	"encoding/json"
	"fmt"
	"mond/wind/utils"
)

type SentinelOption struct {
	Breaker struct {
		MongoOpen      bool `json:"mongoOpen"`
		GrpcServerOpen bool `json:"grpcServerOpen"`
		GrpcClientOpen bool `json:"grpcClientOpen"`
	} `json:"breaker"`
}

func GetSentinelOption() SentinelOption {
	option := _viper.GetStringMap(fmt.Sprintf("sentinel"))
	bytes, err := json.Marshal(option)
	utils.MustNil(err)
	opt := SentinelOption{}
	err = json.Unmarshal(bytes, &opt)
	utils.MustNil(err)
	return opt
}

type Rule struct {
	RetryTimeoutMs               uint32  `json:"retryTimeoutMs"`
	MinRequestAmount             uint64  `json:"minRequestAmount"`
	StatIntervalMs               uint32  `json:"statIntervalMs"`
	StatSlidingWindowBucketCount uint32  `json:"statSlidingWindowBucketCount"`
	MaxAllowedRtMs               uint64  `json:"maxAllowedRtMs"`
	ErrRatioThreshold            float64 `json:"errRatioThreshold"`  //这是给SlowRequestRatio ErrorRatio用的
	SlowRatioThreshold           float64 `json:"slowRatioThreshold"` //这是给SlowRequestRatio ErrorRatio用的
	ErrCountThreshold            float64 `json:"errCountThreshold"`  //给ErrorCount用的
}

func GetSentinelBreaker(scope string) *Rule {
	opt := _viper.GetStringMap(fmt.Sprintf("sentinel.breaker.%s", scope))
	if len(opt) == 0 {
		return nil
	}
	rule := &Rule{}
	bytes, _ := json.Marshal(opt)
	json.Unmarshal(bytes, rule)
	return rule
}
