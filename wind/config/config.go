package config

import (
	"bytes"
	"context"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	env2 "mond/wind/env"
	"mond/wind/hook"
	"mond/wind/logger"
	"mond/wind/utils"
	constant2 "mond/wind/utils/constant"
	"net"
	"os"
	"strings"
)

type ConfigBase struct {
	client config_client.IConfigClient
}

var (
	baseConfig = []byte(`
base:
  NamespaceId: "f456ea81-f4ff-423b-bb7d-0134bd2520f4"
  ServerIp: "dev-nacos.neoclub.cn"
base_develop:
  NamespaceId: "f456ea81-f4ff-423b-bb7d-0134bd2520f4"
  ServerIp: "dev-nacos.neoclub.cn"
base_test:
  NamespaceId: "499fbff2-970b-457a-9dfe-2a65f3917c61"
  ServerIp: "mse-7ea41ab0-nacos-ans.mse.aliyuncs.com"
`)
	_viper *viper.Viper
)

func (m *ConfigBase) InitLocalConfig(ctx context.Context, frameHook hook.FrameStartHook) {
	_viper = viper.New()
	env := os.Getenv("ENV")
	env2.SetEnv(env)
	hostName := os.Getenv("HostName")
	if hostName == "" {
		hostName = os.Getenv("HOSTNAME")
		if hostName == "" {
			hostName = "unknown"
		}
	}
	env2.SetHostName(hostName)
	cluster := os.Getenv("CLUSTER")
	if cluster != "" {
		env2.SetCluster(cluster)
	}
	podIp := os.Getenv("PodIp")
	if podIp == "" {
		podIp = getIp()
	}
	if env == "develop" {
		podIp = "192.168.155.10"
	}
	env2.SetPodIp(podIp)
	
	switch env {
	case "":
	case "develop":
	case "test":
	case "prod":
	default:
		panic(fmt.Sprintf("不支持的ENV:%s", env))
	}
	_viper.SetConfigName("config_yaml")
	_viper.SetConfigType("yaml")
	_viper.AddConfigPath("./conf")
	err := _viper.ReadInConfig()
	utils.MustNil(err)
	
	appId := _viper.GetString("appId")
	if appId == "" {
		panic("appId不能为空")
	}
	env2.SetAppId(appId)
	port := _viper.GetInt32("port")
	if port == 0 {
		panic("port不能为空")
	}
	env2.SetPort(port)
	
	//_viper.SetConfigName("config_json")
	//_viper.SetConfigType("json")
	//err = _viper.MergeInConfig()
	//utils.MustNil(err)
	_viper.SetConfigType("yaml")
	err = _viper.MergeConfig(bytes.NewBuffer(baseConfig))
	utils.MustNil(err)
	
	SetNamespaceId()
	SetServerIp()
	
	//加载环境配置文件
	if GetEnv() != "" {
		_viper.SetConfigName(fmt.Sprintf("config_yaml_%s", GetEnv()))
		_viper.SetConfigType("yaml")
		_viper.AddConfigPath("./conf")
		err = _viper.MergeInConfig()
		utils.MustNil(err)
		//_viper.SetConfigName(fmt.Sprintf("config_json_%s", GetEnv()))
		//_viper.SetConfigType("json")
		//err = _viper.MergeInConfig()
		//utils.MustNil(err)
	}
	//加载全局配置
	sc := []constant.ServerConfig{
		{
			IpAddr: env2.GetServerIp(),
			Port:   8848,
		},
	}
	cc := constant.ClientConfig{
		NamespaceId:         env2.GetNamespaceId(), //namespace id
		TimeoutMs:           3000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "error",
	}
	// a more graceful way to create config client
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	utils.MustNil(err)
	m.client = client
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: "MetaBase",
		Group:  constant2.Group,
	})
	utils.MustNil(err)
	if content == "" {
		panic("没有拉到基础配置")
	}
	mm := map[string]interface{}{}
	//err = json.Unmarshal([]byte(content), &mm)
	//if err != nil {
	err = yaml.Unmarshal([]byte(content), &mm)
	if err != nil {
		panic("基础配置必须是yaml")
	}
	//}
	err = _viper.MergeConfigMap(mm)
	utils.MustNil(err)
	//加载服务自身配置
	content, err = client.GetConfig(vo.ConfigParam{
		DataId: GetAppid(),
		Group:  constant2.Group,
	})
	utils.MustNil(err)
	//允许有服务没有设置自己的配置
	if content == "" {
		return
	}
	mm = map[string]interface{}{}
	//err = json.Unmarshal([]byte(content), &mm)
	//if err != nil {
	err = yaml.Unmarshal([]byte(content), &mm)
	if err != nil {
		panic("服务本身的配置必须是json或yaml")
	}
	err = _viper.MergeConfigMap(mm)
	utils.MustNil(err)
	//}
	//err = _viper.MergeConfigMap(mm)
	//utils.MustNil(err)
	//监听config变化
	err = client.ListenConfig(vo.ConfigParam{
		DataId: GetAppid(),
		Group:  constant2.Group,
		OnChange: func(namespace, group, dataId, data string) {
			mm := map[string]interface{}{}
			//err = json.Unmarshal([]byte(data), &mm)
			//if err != nil {
			err = yaml.Unmarshal([]byte(data), &mm)
			if err != nil {
				logger.GetLogger().Error(ctx, "动态改变的配置必须是json或yaml格式", zap.Any("err", err))
				return
			}
			
			//}
			if mm["dynamic"] == nil {
				return
			}
			err = _viper.MergeConfigMap(map[string]interface{}{"dynamic": mm["dynamic"]})
			if err != nil {
				logger.GetLogger().Error(ctx, "动态改变时merge出错", zap.Any("err", err))
				return
			}
			frameHook.ConfigChangeHook(_viper.GetStringMap("dynamic"))
		},
	})
	utils.MustNil(err)
	//第一次拉取完配置也需要配置一下动态配置
	if _viper.GetStringMap("dynamic") != nil {
		frameHook.ConfigChangeHook(_viper.GetStringMap("dynamic"))
	}
}

func (m *ConfigBase) GracefulStop() {
	m.client.CancelListenConfig(vo.ConfigParam{
		DataId: GetAppid(),
		Group:  constant2.Group,
	})
}

func GetString(key string) string {
	return _viper.GetString(key)
}

func GetInt32(key string) int32 {
	return _viper.GetInt32(key)
}

func getIp() string {
	conn, err := net.Dial("udp", "baidu.com:80")
	utils.MustNil(err)
	defer conn.Close()
	ipAddr := strings.Split(conn.LocalAddr().String(), ":")[0]
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		panic("ip must not be nil")
	}
	return ip.String()
}

func GetSentry() string {
	return _viper.GetString("sentry")
}
