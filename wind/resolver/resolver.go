package resolver

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
	"mond/wind/registry"
	"mond/wind/registry/define"
	balancer2 "mond/wind/resolver/balancer"
	"mond/wind/resolver/nacos"
)

type ResolverBase struct {
	ErrChan chan error
	rb      *nacos.ResolverBuilder
}

func (m *ResolverBase) Init(reg registry.Registry) {
	rb := nacos.NewMetaResolverBuilder(reg)
	//注册meta框架的resolver
	resolver.Register(rb)
	//注册负载均衡器
	_balancer, err := balancer2.NewBalancerBuilder()
	if err != nil {
		m.ErrChan <- err
		return
	}
	balancer.Register(base.NewBalancerBuilder("meta", _balancer, base.Config{HealthCheck: true}))
}
func (m *ResolverBase) GetInstance(appId string) []*define.Instance {
	return m.rb.GetInstance(appId)
}
