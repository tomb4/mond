package breaker

import (
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"mond/wind/config"
	"strings"
)

const (
	
	//熔断器在scope是server的时候的默认配置
	DefaultServerRetryTimeoutMs               = 10000 //10S
	DefaultServerMinRequestAmount             = 100   //
	DefaultServerStatSlidingWindowBucketCount = 10    //
	DefaultServerThreshold                    = 0.7   //阈值 70%
	DefaultServerStatIntervalMs               = 10000 //10S
	
	DefaultClientRetryTimeoutMs               = 10000 //10S
	DefaultClientMinRequestAmount             = 100   //
	DefaultClientStatSlidingWindowBucketCount = 10    //
	DefaultClientThreshold                    = 0.7   //阈值 70%
	DefaultClientStatIntervalMs               = 10000 //10S
	
	//scope是mongo的时候的默认配置
	DefaultMongoRetryTimeoutMs               = 10000 //10S
	DefaultMongoMinRequestAmount             = 20    //
	DefaultMongoStatSlidingWindowBucketCount = 10
	DefaultMongoMaxAllowedRtMs               = 99    //99ms
	DefaultMongoThreshold                    = 0.7   //阈值
	DefaultMongoStatIntervalMs               = 10000 //10S
	DefaultMongoErrorCountThreshold          = 50    //阈值
	
	Type_Server = "server"
	Type_Client = "client"
	Type_Mongo  = "mongo"
)

func getRuleByScope(scopes []string) (res []*circuitbreaker.Rule) {
	//fmt.Println("scopes", scopes)
	scope := strings.Join(scopes, ".")
	rule := config.GetSentinelBreaker(scope)
	if rule == nil {
		scope := strings.Join(scopes[0:2], ".")
		rule = config.GetSentinelBreaker(scope)
	}
	if rule == nil {
		scope := scopes[0]
		rule = config.GetSentinelBreaker(scope)
	}
	var ruleSlowRatio *circuitbreaker.Rule
	var ruleErrorRatio *circuitbreaker.Rule
	var ruleErrorCount *circuitbreaker.Rule
	switch scopes[0] {
	case Type_Server:
		ruleErrorRatio = &circuitbreaker.Rule{
			Resource:                     strings.Join(scopes, "."),
			Strategy:                     circuitbreaker.ErrorRatio,
			RetryTimeoutMs:               DefaultServerRetryTimeoutMs,
			MinRequestAmount:             DefaultServerMinRequestAmount,
			StatIntervalMs:               DefaultServerStatIntervalMs,
			StatSlidingWindowBucketCount: DefaultServerStatSlidingWindowBucketCount,
			Threshold:                    DefaultServerThreshold,
		}
	case Type_Client:
		ruleErrorRatio = &circuitbreaker.Rule{
			Resource:                     strings.Join(scopes, "."),
			Strategy:                     circuitbreaker.ErrorRatio,
			RetryTimeoutMs:               DefaultClientRetryTimeoutMs,
			MinRequestAmount:             DefaultClientMinRequestAmount,
			StatIntervalMs:               DefaultClientStatIntervalMs,
			StatSlidingWindowBucketCount: DefaultClientStatSlidingWindowBucketCount,
			Threshold:                    DefaultClientThreshold,
		}
	case Type_Mongo:
		ruleSlowRatio = &circuitbreaker.Rule{
			Resource:                     strings.Join(scopes, "."),
			Strategy:                     circuitbreaker.SlowRequestRatio,
			RetryTimeoutMs:               DefaultMongoRetryTimeoutMs,
			MinRequestAmount:             DefaultMongoMinRequestAmount,
			StatIntervalMs:               DefaultMongoStatIntervalMs,
			StatSlidingWindowBucketCount: DefaultMongoStatSlidingWindowBucketCount,
			Threshold:                    DefaultMongoThreshold,
			MaxAllowedRtMs:               DefaultMongoMaxAllowedRtMs,
		}
		ruleErrorCount = &circuitbreaker.Rule{
			Resource:                     strings.Join(scopes, "."),
			Strategy:                     circuitbreaker.ErrorCount,
			RetryTimeoutMs:               DefaultMongoRetryTimeoutMs,
			MinRequestAmount:             DefaultMongoMinRequestAmount,
			StatIntervalMs:               DefaultMongoStatIntervalMs,
			StatSlidingWindowBucketCount: DefaultMongoStatSlidingWindowBucketCount,
			Threshold:                    DefaultMongoErrorCountThreshold,
		}
	default:
		panic("未知type")
	}
	if rule != nil && rule.MaxAllowedRtMs != 0 {
		if ruleErrorRatio != nil {
			ruleErrorRatio.MaxAllowedRtMs = rule.MaxAllowedRtMs
		}
		if ruleSlowRatio != nil {
			ruleSlowRatio.MaxAllowedRtMs = rule.MaxAllowedRtMs
		}
		if ruleErrorCount != nil {
			ruleErrorCount.MaxAllowedRtMs = rule.MaxAllowedRtMs
		}
	}
	if rule != nil && rule.RetryTimeoutMs != 0 {
		if ruleErrorRatio != nil {
			ruleErrorRatio.RetryTimeoutMs = rule.RetryTimeoutMs
		}
		if ruleSlowRatio != nil {
			ruleSlowRatio.RetryTimeoutMs = rule.RetryTimeoutMs
		}
		if ruleErrorCount != nil {
			ruleErrorCount.RetryTimeoutMs = rule.RetryTimeoutMs
		}
	}
	if rule != nil && rule.MinRequestAmount != 0 {
		if ruleErrorRatio != nil {
			ruleErrorRatio.MinRequestAmount = rule.MinRequestAmount
		}
		if ruleSlowRatio != nil {
			ruleSlowRatio.MinRequestAmount = rule.MinRequestAmount
		}
		if ruleErrorCount != nil {
			ruleErrorCount.MinRequestAmount = rule.MinRequestAmount
		}
	}
	if rule != nil && rule.StatIntervalMs != 0 {
		if ruleErrorRatio != nil {
			ruleErrorRatio.StatIntervalMs = rule.StatIntervalMs
		}
		if ruleSlowRatio != nil {
			ruleSlowRatio.StatIntervalMs = rule.StatIntervalMs
		}
		if ruleErrorCount != nil {
			ruleErrorCount.StatIntervalMs = rule.StatIntervalMs
		}
	}
	if rule != nil && rule.ErrRatioThreshold != 0 {
		if ruleErrorRatio != nil {
			ruleErrorRatio.Threshold = rule.ErrRatioThreshold
		}
	}
	if rule != nil && rule.ErrCountThreshold != 0 {
		if ruleErrorCount != nil {
			ruleErrorCount.Threshold = rule.ErrCountThreshold
		}
	}
	if rule != nil && rule.SlowRatioThreshold != 0 {
		if ruleSlowRatio != nil {
			ruleSlowRatio.Threshold = rule.SlowRatioThreshold
		}
	}
	if rule != nil && rule.StatSlidingWindowBucketCount != 0 {
		if ruleErrorRatio != nil {
			ruleErrorRatio.StatSlidingWindowBucketCount = rule.StatSlidingWindowBucketCount
		}
		if ruleSlowRatio != nil {
			ruleSlowRatio.StatSlidingWindowBucketCount = rule.StatSlidingWindowBucketCount
		}
		if ruleErrorCount != nil {
			ruleErrorCount.StatSlidingWindowBucketCount = rule.StatSlidingWindowBucketCount
		}
	}
	
	if ruleSlowRatio != nil {
		res = append(res, ruleSlowRatio)
	}
	if ruleErrorRatio != nil {
		res = append(res, ruleErrorRatio)
	}
	if ruleErrorCount != nil {
		res = append(res, ruleErrorCount)
	}
	return
}
