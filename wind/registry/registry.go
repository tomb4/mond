package registry

import (
	"context"
	"go.uber.org/zap"
	"mond/wind/logger"
	"mond/wind/registry/define"
	"mond/wind/registry/nacos"
)


type Registry interface {
	RegistryService() error
	DeregisterInstance() error
	//查询所有健康的实例
	SelectInstances(serviceName string) ([]*define.Instance, error)
	//订阅服务变化
	Subscribe(serviceName string, callback define.Subscribe) define.SubscribeEntity
	//获取所有服务
	GetAllServices() ([]string, error)
}

type RegistryBase struct {
	ErrChan chan error
	reg     Registry
}

func (m *RegistryBase) GetRegistry() Registry {
	return m.reg
}

func (m *RegistryBase) InitRegistry() {
	registry := nacos.Init()
	//err := registry.RegistryService()
	//if err != nil {
	//	m.ErrChan <- err
	//}
	m.reg = registry
}

func (m *RegistryBase) RegistryService() {
	err := m.reg.RegistryService()
	if err != nil {
		m.ErrChan <- err
	}
}
func (m *RegistryBase) GracefulStop() {
	err := m.reg.DeregisterInstance()
	if err != nil {
		logger.GetLogger().Error(context.TODO(), "DeregisterInstance", zap.Any("err", err))
	}
}
