package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"mond/wind/config"
	env2 "mond/wind/env"
	"mond/wind/utils"
)

type nacosRegistry struct {
	namingClient naming_client.INamingClient
}

func Init() *nacosRegistry {
	clientConfig := constant.ClientConfig{
		NamespaceId:         env2.GetNamespaceId(), // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           3000,
		ListenInterval:      10000, //config的监听每10s一次
		BeatInterval:        1000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "error",
		AppName:             config.GetAppid(),
	}
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      env2.GetServerIp(),
			ContextPath: "/nacos",
			Port:        8848,
			Scheme:      "http",
		},
	}
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	utils.MustNil(err)
	reg := nacosRegistry{}
	reg.namingClient = namingClient
	
	//FIXME :  调试用的
	//fmt.Println(namingClient.DeregisterInstance(vo.DeregisterInstanceParam{
	//	Ip:          "192.192.203.138",
	//	Port:        28002,
	//	ServiceName: "MetaLogic",
	//	Ephemeral:   false,
	//	Cluster:     Cluster,
	//	GroupName:   Group,
	//}))
	//fmt.Println("111")
	return &reg
}
