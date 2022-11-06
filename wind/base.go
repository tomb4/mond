package wind

import (
	"context"
	"fmt"
	mredis "mond/wind/cache/redis"
	"mond/wind/config"
	mongodb "mond/wind/db/mongo"
	"mond/wind/env"
	mgrpc "mond/wind/grpc"
	"mond/wind/hook"
	"mond/wind/http/health"
	"mond/wind/logger"
	"mond/wind/mq/rabbit"
	"mond/wind/registry"
	"mond/wind/reload"
	"mond/wind/resolver"
	"mond/wind/sentinel"
	"mond/wind/sentry"
	"mond/wind/trace"
	"mond/wind/utils"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type Base struct {
	configBase   *config.ConfigBase
	grpcBase     *mgrpc.GrpcBase
	registryBase *registry.RegistryBase
	resolverBase *resolver.ResolverBase
	health       *health.Health
	errChan      chan error
	ctx          context.Context
}

var (
	base Base
)

type Option func(f *frameOption)

type frameOption struct {
	testMode bool
}

func (m *frameOption) InTestMode() bool {
	return m.testMode
}

var (
	defaultOption *frameOption = &frameOption{}
)

func WithTestMode() Option {
	return func(f *frameOption) {
		f.testMode = true
	}
}

//初始化框架
func InitFrame(hook hook.FrameStartHook, opts ...Option) {
	ctx := context.Background()
	logger.GetLogger().Info(ctx, "框架开始启动")
	for _, f := range opts {
		logger.GetLogger().Info(ctx, "框架将以测试模式启动")
		f(defaultOption)
	}
	if defaultOption.InTestMode() {
		utils.ChangePosition()
	}
	errChan := make(chan error, 1)
	base = Base{
		configBase:   &config.ConfigBase{},
		grpcBase:     &mgrpc.GrpcBase{ErrChan: errChan},
		registryBase: &registry.RegistryBase{ErrChan: errChan},
		resolverBase: &resolver.ResolverBase{ErrChan: errChan},
		health:       health.NewHealth(),
		errChan:      errChan,
		ctx:          ctx,
	}
	env.StateEnv = env.Init
	logger.GetLogger().Info(ctx, "框架开始加载基础组件")
	//初始化本地配置
	base.configBase.InitLocalConfig(ctx, hook)
	//init sentry
	sentry.Init()
	//初始化trace
	trace.Init()
	//初始化sentinel
	sentinel.InitSentinel()
	//连接注册、配置中心
	base.registryBase.InitRegistry()
	//初始化grpc resolver与balancer
	base.resolverBase.Init(base.registryBase.GetRegistry())
	//框架初始化
	startChan := make(chan int32)
	//初始化mongo dbm
	mongoDbManager = mongodb.NewDbManager()
	//初始化redis dbm
	redisManager = mredis.NewDbManager()
	//初始化rabbit
	rabbitManager = rabbit.NewRabbitManager()
	//进入业务自定义初始化状态
	env.StateEnv = env.Starting
	logger.GetLogger().Info(ctx, "框架开始初始化业务资源")
	err := hook.ResourceInitHook(ctx)
	if err != nil {
		base.errChan <- err
		logger.GetLogger().Error(ctx, "框架初始化业务资源报错，准备退出", zap.Any("err", err))
		goto stop
	}
	if !defaultOption.InTestMode() {
		//启动grpc server
		logger.GetLogger().Info(ctx, "框架开始启动grpc服务")
		go base.grpcBase.InitGrpcServer(ctx, hook, startChan)
		select {
		case e := <-base.errChan:
			base.errChan <- e
		case <-startChan:
			//如果启动成功，则开始注册
			logger.GetLogger().Info(ctx, "框架开始注册服务")
			go base.health.Health()
			base.registryBase.RegistryService()
			reload.Reload(ctx)
			env.StateEnv = env.Running
			logger.GetLogger().Info(ctx, "框架开始平稳运行中")
		}
	} else {
		logger.GetLogger().Info(ctx, "框架以单元测试模式启动成功")
		return
	}
stop:
	sign := make(chan os.Signal)
	signal.Notify(sign, os.Interrupt, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT)
	select {
	case s := <-sign:
		logger.GetLogger().Info(ctx, fmt.Sprintf("receive signal: %v", s))
	case e := <-base.errChan:
		logger.GetLogger().Info(ctx, fmt.Sprintf("receive err chan: %v", e))
	}
	env.StateEnv = env.Stopping

	//优雅退出保证最少有5s
	timer := time.NewTimer(time.Second * 5)
	go func() {
		//先关闭流量入口，在注册中心注销自己
		base.Stop()
		//给业务处理的空间
		hook.AppStopHook()
	}()
	select {
	case <-timer.C:
	}
	//关闭基础资源 如：db
	mongoDbManager.Close()
	redisManager.Close()
	rabbitManager.Close()
	logger.GetLogger().Info(ctx, "应用退出成功")
}

func (m *Base) Stop() {
	m.configBase.GracefulStop()
	m.registryBase.GracefulStop()
	m.grpcBase.GracefulStop()
	m.health.Shutdown(m.ctx)
}
