package nacos

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
	"mond/wind/logger"
	"mond/wind/registry"
	"mond/wind/registry/define"
	"sync"
)

const scheme = "meta"

type ResolverBuilder struct {
	reg registry.Registry
	targetMap     sync.Map
	targetMapLock sync.RWMutex
}

func (m *ResolverBuilder) Scheme() string {
	return scheme
}
func NewMetaResolverBuilder(reg registry.Registry) *ResolverBuilder {
	return &ResolverBuilder{reg: reg}
}
func (m *ResolverBuilder) GetInstance(appId string) []*define.Instance {
	val, ok := m.targetMap.Load(appId)
	if !ok {
		return nil
	}
	r := val.(*nacosResolver)
	return r.instances
}
func (m *ResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	appId := target.URL.Host
	r := &nacosResolver{
		target: target,
		cc:     cc,
		appId:  appId,
		reg:    m.reg,
	}
	r.start()
	logger.GetLogger().Debug(context.TODO(), "Build", zap.String("targetAppId", r.appId))
	m.targetMap.Store(appId, r)
	return r, nil
}

type nacosResolver struct {
	target          resolver.Target
	cc              resolver.ClientConn
	appId           string
	reg             registry.Registry
	instances       []*define.Instance
	subscribeEntity define.SubscribeEntity
}

func (m *nacosResolver) start() {
	instances, err := m.reg.SelectInstances(m.appId)
	//如果刚启动时没有获取到地址则忽略
	if err != nil {
		logger.GetLogger().Error(context.TODO(), "nacosResolver", zap.Any("err", err.Error()), zap.String("targetAppId", m.appId))
	}
	m.instances = instances
	addrs := make([]resolver.Address, 0, len(instances))
	for _, v := range instances {
		attr := attributes.New("hostName", v.Metadata["hostName"])
		if v.Metadata["label"] != "" {
			attr = attr.WithValue("label", v.Metadata["label"])
		}
		addrs = append(addrs, resolver.Address{
			Addr:               fmt.Sprintf("%s:%d", v.Ip, v.Port),
			BalancerAttributes: attr,
		})
	}
	m.cc.UpdateState(resolver.State{Addresses: addrs})
	//订阅地址变化
	subscribeEntity := m.reg.Subscribe(m.appId, func(instances []*define.Instance) {
		//fmt.Println("Subscribe", utils.StructToJson(instances))
		m.instances = instances
		addrs := make([]resolver.Address, 0, len(instances))
		for _, v := range instances {
			attr := attributes.New("hostName", v.Metadata["hostName"])
			if v.Metadata["label"] != "" {
				attr = attr.WithValue("label", v.Metadata["label"])
			}
			addrs = append(addrs, resolver.Address{
				Addr:               fmt.Sprintf("%s:%d", v.Ip, v.Port),
				BalancerAttributes: attr,
			})
		}
		m.cc.UpdateState(resolver.State{Addresses: addrs})
	})
	m.subscribeEntity = subscribeEntity
}

func (*nacosResolver) ResolveNow(o resolver.ResolveNowOptions) {
}

func (m *nacosResolver) Close() {
	m.subscribeEntity.Unsubscribe()
	//m.reg.Unsubscribe(m.appId)
}
