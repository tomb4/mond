package rabbit

import (
	"context"
	"mond/wind/env"
	"sync"
	"time"
)

type client struct {
}

var (
	defaultHeartbeat         = 10 * time.Second
	defaultConnectionTimeout = 30 * time.Second
)

type rabbitManager struct {
	connStore map[string]*conn
	connLock  sync.Mutex
	ctx       context.Context
	cancel    context.CancelFunc
}

type RabbitManager interface {
	Close()
	GetClient(instance string) (Client, error)
	GetAsyncClient() (Client, error)
}

func NewRabbitManager() RabbitManager {
	if env.GetAppState() != env.Init {
		panic("rabbit NewRabbitManager只能在frame初始化时可以自加载")
	}
	rm := rabbitManager{connStore: map[string]*conn{}}
	rm.ctx, rm.cancel = context.WithCancel(context.Background())
	return &rm
}

type Client interface {
	Close()
	Consume(queue string, autoAck bool, handler Consume, opts ...Option) error
	InitProducer() error
	Publish(ctx context.Context, msg PublishMessage) error
	CreateExchange(name, kind string, durable, autoDelete bool) error
	CreateQueue(name string, durable, autoDelete bool) error
	QueueBind(name, key, exchange string) error
}

func (m *rabbitManager) GetClient(instance string) (Client, error) {
	if env.GetAppState() != env.Starting {
		panic("rabbit GetClient只能在ResourceInit中调用")
	}
	m.connLock.Lock()
	defer m.connLock.Unlock()
	if c, ok := m.connStore[instance]; ok {
		return c, nil
	}
	cc, err := newConn(instance)
	if err != nil {
		return nil, err
	}
	m.connStore[instance] = cc
	return m.connStore[instance], nil
}

func (m *rabbitManager) GetAsyncClient() (Client, error) {
	if env.GetAppState() != env.Starting {
		panic("rabbit GetClient只能在ResourceInit中调用")
	}
	instance := "async"
	m.connLock.Lock()
	defer m.connLock.Unlock()
	if c, ok := m.connStore[instance]; ok {
		return c, nil
	}
	cc, err := newAsyncConn()
	if err != nil {
		return nil, err
	}
	m.connStore[instance] = cc
	return m.connStore[instance], nil
}

func (m *rabbitManager) Close() {
	for _, v := range m.connStore {
		v.Close()
	}
}
