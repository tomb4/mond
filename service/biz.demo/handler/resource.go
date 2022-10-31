package handler

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"mond/service/biz.demo/app"
	"mond/service/biz.demo/domain/demo"
	"mond/wind"
	"sync"
	"time"
)

var (
	resourceLock sync.Mutex
	resourceInit bool
)

func (m *BizdemoService) ResourceInit(ctx context.Context) error {
	resourceLock.Lock()
	defer resourceLock.Unlock()
	if resourceInit {
		panic("ResourceInit已经执行过了")
	}
	rdb, err := wind.GetRedisDbManager().GetClient("master")
	if err != nil {
		return err
	}

	dbm := wind.GetMongoDbManager()
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	demoColl, err := dbm.GetCollection("master.meta.demo")
	if err != nil {
		return err
	}
	index := []mongo.IndexModel{
		{
			Keys: bsonx.Doc{
				{Key: "id", Value: bsonx.Int32(1)},
			},
			Options: options.Index().SetBackground(true).SetUnique(true),
		},
	}
	demoColl.Indexes().CreateMany(ctx, index, opts)

	demoSvc := demo.NewService(demoColl, rdb)

	m.app = app.NewApp(rdb, demoSvc)

	resourceInit = true
	return nil
}
