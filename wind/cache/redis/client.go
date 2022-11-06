package mredis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"mond/wind/config"
	"mond/wind/env"
	"sync"
	"time"
)

var (
	addrEmpty = errors.New("redis地址为空")
)

type Client struct {
	*redis.Client
	lockScriptSha string
	rand          *rand.Rand
	lock          sync.Mutex
}

func (m *Client) int64(max int64) int64 {
	m.lock.Lock()
	n := m.rand.Int63n(max)
	m.lock.Unlock()
	return n
}

type DbManager interface {
	GetClient(instance string) (*Client, error)
	Close()
}

type dbManager struct {
	clientMap  map[string]*Client
	clientLock sync.Mutex
}

func (m *dbManager) Close() {
	m.clientLock.Lock()
	defer m.clientLock.Unlock()
	for _, v := range m.clientMap {
		v.Close()
	}
}

func NewDbManager() DbManager {
	if env.GetAppState() != env.Init {
		panic("redis NewDbManager只有frame初始化时可以自加载")
	}
	return &dbManager{clientMap: map[string]*Client{}}
}

func (m *dbManager) GetClient(instance string) (*Client, error) {
	if env.GetAppState() != env.Starting {
		panic("redis GetClient只能在ResourceInit中调用")
	}
	m.clientLock.Lock()
	defer m.clientLock.Unlock()

	if m.clientMap[instance] != nil {
		return m.clientMap[instance], nil
	}

	opt, err := config.GetRedisClientOption(instance)
	if err != nil {
		return nil, err
	}
	if opt.Addr == "" {
		return nil, addrEmpty
	}
	if opt.PoolSize == 0 {
		opt.PoolSize = 20
	}
	if opt.MinIdleConns == 0 {
		opt.MinIdleConns = 5
	}
	rdbCli := redis.NewClient(&redis.Options{
		Addr:         opt.Addr,
		Password:     opt.Password,     // no password set
		DB:           opt.Db,           // use default DB
		PoolSize:     opt.PoolSize,     // 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
		MinIdleConns: opt.MinIdleConns, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。
		//超时
		DialTimeout:  time.Second,            //连接建立超时时间，默认5秒。
		ReadTimeout:  100 * time.Millisecond, //读超时，默认3秒， -1表示取消读超时
		WriteTimeout: 100 * time.Millisecond, //写超时，默认等于读超时
		PoolTimeout:  100 * time.Millisecond, //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。

		//闲置连接检查包括IdleTimeout，MaxConnAge
		IdleCheckFrequency: 60 * time.Second, //闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理。
		IdleTimeout:        5 * time.Minute,  //闲置超时，默认5分钟，-1表示取消闲置超时检查
		MaxConnAge:         0 * time.Second,  //连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接

		//命令执行失败时的重试策略
		MaxRetries:      0,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
		MinRetryBackoff: 8 * time.Millisecond,   //每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
		MaxRetryBackoff: 512 * time.Millisecond, //每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔
	})
	rdbCli.AddHook(&hook{instance: instance, appId: config.GetAppid()})
	res, err := rdbCli.ScriptLoad(context.TODO(), compareAndDeleteLuaScript).Result()
	if err != nil {
		return nil, err
	}
	m.clientMap[instance] = &Client{Client: rdbCli, lockScriptSha: res, rand: rand.New(rand.NewSource(time.Now().UnixNano()))}
	return m.clientMap[instance], nil
}
