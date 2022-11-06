package sentinel

import (
	"context"
	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/config"
	"go.uber.org/zap"
	"mond/wind/env"
	"mond/wind/logger"
	"mond/wind/utils"
	"sync"
)

var (
	once   sync.Once
	lister = &stateChangeTestListener{}
)

func InitSentinel() {
	once.Do(func() {
		conf := config.NewDefaultConfig()
		conf.Sentinel.Log.Logger = nil
		conf.Sentinel.App.Name = env.GetAppId()
		err := api.InitWithConfig(conf)
		if err != nil {
			return
		}
		utils.MustNil(err)
		circuitbreaker.RegisterStateChangeListeners(lister)
		logger.GetLogger().Info(context.TODO(), "sentinel init")
	})
}

type stateChangeTestListener struct {
	fuseMap sync.Map
}

func (s *stateChangeTestListener) OnTransformToClosed(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	logger.GetLogger().Info(context.TODO(), "OnTransformToClosed", zap.Any("rule", rule), zap.Int32("prev", int32(prev)))
	s.fuseMap.Delete(rule.Resource)
	//fmt.Printf("rule.resource: %s, rule.steategy: %+v, From %s to Closed, time: %d\n", rule.Resource, rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func (s *stateChangeTestListener) OnTransformToOpen(prev circuitbreaker.State, rule circuitbreaker.Rule, snapshot interface{}) {
	logger.GetLogger().Error(context.TODO(), "OnTransformToOpen", zap.Any("rule", rule), zap.Int32("prev", int32(prev)), zap.Any("snapshot", snapshot))
	s.fuseMap.Store(rule.Resource, 1)
	//fmt.Printf("rule.resource: %s, rule.steategy: %+v, From %s to Open, snapshot: %d, time: %d\n", rule.Resource, rule.Strategy, prev.String(), snapshot, util.CurrentTimeMillis())
}

func (s *stateChangeTestListener) OnTransformToHalfOpen(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	s.fuseMap.Store(rule.Resource, 2)
	logger.GetLogger().Info(context.TODO(), "OnTransformToOpen", zap.Any("rule", rule), zap.Int32("prev", int32(prev)))
	//fmt.Printf("rule.resource: %s, rule.steategy: %+v, From %s to Half-Open, time: %d\n", rule.Resource, rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

//检查资源是否处于熔断状态 如果处于半开状态，则迅速修改
func ResourceInFuse(r string) bool {
	v, ok := lister.fuseMap.Load(r)
	//fmt.Println("ResourceInFuse", v, ok)
	if !ok {
		return false
	}
	if v == 2 {
		v = 1
	}
	return true
}
