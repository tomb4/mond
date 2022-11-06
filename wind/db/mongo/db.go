package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mond/wind/env"
	"sync"
)

type dbManager struct {
	clientMap  map[string]*Client
	clientLock sync.Mutex
	collLock   sync.Mutex
}

type DbManager interface {
	GetClient(instance string) (*Client, error)
	GetCollection(endpoint string, opts ...*options.CollectionOptions) (*Collection, error)
	Close()
}

func NewDbManager() DbManager {
	if env.GetAppState() != env.Init {
		panic("mongo NewDbManager只有frame初始化时可以自加载")
	}
	return &dbManager{clientMap: map[string]*Client{}}
}
func (m *dbManager) Close() {
	m.collLock.Lock()
	defer m.collLock.Unlock()
	for _, v := range m.clientMap {
		v.cli.Disconnect(context.TODO())
	}
}
