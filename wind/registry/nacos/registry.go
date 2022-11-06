package nacos

import (
	"context"
	"errors"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"go.uber.org/zap"
	"mond/wind/config"
	"mond/wind/env"
	"mond/wind/logger"
	"mond/wind/registry/define"
	"mond/wind/utils"
	"mond/wind/utils/constant"
	"time"
)

func (n *nacosRegistry) RegistryService() error {
	ok, err := n.namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          env.GetPodIp(),
		Port:        uint64(env.GetPort()),
		ServiceName: config.GetAppid(),
		Weight:      1,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata: map[string]string{
			"hostName":  config.GetHostName(),
			"label":     "stable",
			"createdAt": utils.FormatTime(time.Now(), "MM/DD HH:mm:ss"),
		},
		ClusterName: env.GetCluster(),
		GroupName:   constant.Group,
	})
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("注册节点失败")
	}
	return nil
}

func (n *nacosRegistry) DeregisterInstance() error {
	ok, err := n.namingClient.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          env.GetPodIp(),
		Port:        uint64(env.GetPort()),
		ServiceName: config.GetAppid(),
		Ephemeral:   true,
		Cluster:     env.GetCluster(),
		GroupName:   constant.Group,
	})
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("取消注册节点失败")
	}
	return nil
}

func (n *nacosRegistry) SelectInstances(serviceName string) ([]*define.Instance, error) {
	instances, err := n.namingClient.SelectInstances(vo.SelectInstancesParam{
		ServiceName: serviceName,
		GroupName:   constant.Group,             // 默认值DEFAULT_GROUP
		Clusters:    []string{env.GetCluster()}, // 默认值DEFAULT
		HealthyOnly: true,
	})
	if err != nil {
		return nil, err
	}
	items := make([]*define.Instance, 0, len(instances))
	for _, v := range instances {
		items = append(items, &define.Instance{
			ServiceName: v.ServiceName,
			ClusterName: v.ClusterName,
			Ip:          v.Ip,
			Port:        v.Port,
			Metadata:    v.Metadata,
			Weight:      v.Weight,
			Enable:      v.Enable,
			Healthy:     v.Healthy,
			Ephemeral:   v.Ephemeral,
		})
	}
	return items, nil
}

type nacosSubscribeEntity struct {
	param *vo.SubscribeParam
	m     *nacosRegistry
}

func (m *nacosSubscribeEntity) Unsubscribe() {
	err := m.m.namingClient.Unsubscribe(m.param)
	if err != nil {
		logger.GetLogger().Error(context.TODO(), "Unsubscribe err", zap.Error(err))
	}
}

func (m *nacosRegistry) Subscribe(serviceName string, callback define.Subscribe) define.SubscribeEntity {
	param := &vo.SubscribeParam{
		ServiceName: serviceName,
		GroupName:   constant.Group,
		Clusters:    []string{env.GetCluster()},
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			if err != nil {
				logger.GetLogger().Error(context.TODO(), "Subscribe", zap.Any("err", err))
				return
			}
			items := make([]*define.Instance, 0, len(services))
			for _, v := range services {
				items = append(items, &define.Instance{
					ServiceName: v.ServiceName,
					ClusterName: v.ClusterName,
					Ip:          v.Ip,
					Port:        v.Port,
					Metadata:    v.Metadata,
					Weight:      v.Weight,
					Enable:      v.Enable,
					Healthy:     v.Healthy,
				})
			}
			callback(items)
		},
	}
	m.namingClient.Subscribe(param)
	return &nacosSubscribeEntity{param: param, m: m}
}

//func (m *nacosRegistry) Unsubscribe(serviceName string) {
//	err := m.namingClient.Unsubscribe(&vo.SubscribeParam{
//		ServiceName: serviceName,
//		GroupName:   constant.Group,
//		Clusters:    []string{constant.Cluster},
//	})
//	//fmt.Println("Unsubscribe", serviceName)
//	if err != nil {
//		logger.GetLogger().Error(context.TODO(), "Unsubscribe", zap.Any("err", err))
//	}
//}

func (m *nacosRegistry) GetAllServices() ([]string, error) {
	svcList, err := m.namingClient.GetAllServicesInfo(vo.GetAllServiceInfoParam{
		NameSpace: env.GetNamespaceId(),
		GroupName: constant.Group,
		PageNo:    1,
		PageSize:  500,
	})
	if err != nil {
		return nil, err
	}
	return svcList.Doms, nil
}
