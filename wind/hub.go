package wind

import (
	"errors"
	"mond/wind/async"
	mredis "mond/wind/cache/redis"
	"mond/wind/config"
	mongodb "mond/wind/db/mongo"
	"mond/wind/mq/rabbit"
	"mond/wind/registry/define"
)

var (
	mongoDbManager mongodb.DbManager
	redisManager   mredis.DbManager
	rabbitManager  rabbit.RabbitManager
)

func GetMongoDbManager() mongodb.DbManager {
	return mongoDbManager
}
func GetRedisDbManager() mredis.DbManager {
	return redisManager
}

func GetRabbitManager() rabbit.RabbitManager {
	return rabbitManager
}

func GetAsync() (async.Async, error) {
	conf, err := config.GetAsyncOption()
	if err != nil {
		return nil, err
	}
	if conf.Type != "rabbit" {
		return nil, errors.New("只支持rabbit")
	}
	client, err := GetRabbitManager().GetAsyncClient()
	if err != nil {
		return nil, err
	}
	return async.GetRabbitAsyncClient(client)
}

//获取所有服务
func GetAllService() ([]string, error) {
	return base.registryBase.GetRegistry().GetAllServices()
}

//获取服务的所有实例 从grpc client中拿，如果没有该grpc client，则永远拿不到
func GetInstanceFromGrpcClient(appId string) []*define.Instance {
	return base.resolverBase.GetInstance(appId)
}

//获取服务的所有实例 每次都是从注册中心中拿  不适合非常高频的操作
func GetInstancesFromRegistry(appId string) ([]*define.Instance, error) {
	return base.registryBase.GetRegistry().SelectInstances(appId)
}
