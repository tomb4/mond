package mongodb

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mond/wind/env"
	"mond/wind/sentinel/breaker"
	"mond/wind/trace"
	constant2 "mond/wind/utils/constant"
	mctx "mond/wind/utils/ctx"
	"mond/wind/utils/endpoint"
	"strings"
	"time"
)

// endpoint: main.uki.user
func (m *dbManager) GetCollection(endpoint string, opts ...*options.CollectionOptions) (*Collection, error) {
	//if env.GetAppState() != env.Starting {
	//	panic("mongo GetCollection只能在ResourceInit中调用")
	//}
	m.collLock.Lock()
	defer m.collLock.Unlock()
	addr := strings.Split(endpoint, ".")
	if len(addr) != 3 {
		return nil, errors.New("endpoint格式不正确")
	}
	var err error
	client, ok := m.clientMap[addr[0]]
	if !ok {
		client, err = m.GetClient(addr[0])
		if err != nil {
			return nil, err
		}
	}
	coll := client.cli.Database(addr[1]).Collection(addr[2], opts...)
	c := &Collection{coll: coll, dbName: addr[1], collectionName: addr[2]}
	c.makeEndpoint()
	//设置中间件
	c.attachMiddleware(trace.MongoMiddleware)
	//统一的ctx检查
	c.attachMiddleware(ctxDead)
	//数据库的熔断
	c.attachMiddleware(breaker.Middleware)
	return c, nil
}

type Collection struct {
	coll                      *mongo.Collection
	insertOneEndpoint         endpoint.Endpoint
	insertManyEndpoint        endpoint.Endpoint
	deleteOneEndpoint         endpoint.Endpoint
	deleteManyEndpoint        endpoint.Endpoint
	updateByIDEndpoint        endpoint.Endpoint
	updateOneEndpoint         endpoint.Endpoint
	updateManyEndpoint        endpoint.Endpoint
	replaceOneEndpoint        endpoint.Endpoint
	aggregateEndpoint         endpoint.Endpoint
	countDocumentsEndpoint    endpoint.Endpoint
	findEndpoint              endpoint.Endpoint
	findOneEndpoint           endpoint.Endpoint
	findOneAndReplaceEndpoint endpoint.Endpoint
	findOneAndDeleteEndpoint  endpoint.Endpoint
	findOneAndUpdateEndpoint  endpoint.Endpoint
	bulkWriteEndpoint         endpoint.Endpoint
	dbName                    string
	collectionName            string
}

type mongoReq struct {
	document []interface{}
	filter   []interface{}
	opts     []interface{}
	withTFunc interface{}
}

func (m *Collection) makeEndpoint() {
	m.insertOneEndpoint = m.makeInsertOneEndpoint()
	m.insertManyEndpoint = m.makeInsertManyEndpoint()
	m.deleteOneEndpoint = m.makeDeleteOneEndpoint()
	m.deleteManyEndpoint = m.makeDeleteManyEndpoint()
	m.updateByIDEndpoint = m.makeUpdateByIDEndpoint()
	m.updateOneEndpoint = m.makeUpdateOneEndpoint()
	m.updateManyEndpoint = m.makeUpdateManyEndpoint()
	m.replaceOneEndpoint = m.makeReplaceOneEndpoint()
	m.aggregateEndpoint = m.makeAggregateEndpoint()
	m.countDocumentsEndpoint = m.makeCountDocumentsEndpoint()
	m.findEndpoint = m.makeFindEndpoint()
	m.findOneEndpoint = m.makeFindOneEndpoint()
	m.findOneAndReplaceEndpoint = m.makeFindOneAndReplaceEndpoint()
	m.findOneAndDeleteEndpoint = m.makeFindOneAndDeleteEndpoint()
	m.findOneAndUpdateEndpoint = m.makeFindOneAndUpdateEndpoint()
	m.bulkWriteEndpoint = m.makeBulkWriteEndpoint()
	//TODO: 所有接口
}

//封装endpoint
func (m *Collection) attachMiddleware(mw endpoint.MiddleWare) {
	m.insertOneEndpoint = mw(m.insertOneEndpoint)
	m.insertManyEndpoint = mw(m.insertManyEndpoint)
	m.deleteOneEndpoint = mw(m.deleteOneEndpoint)
	m.deleteManyEndpoint = mw(m.deleteManyEndpoint)
	m.updateByIDEndpoint = mw(m.updateByIDEndpoint)
	m.updateOneEndpoint = mw(m.updateOneEndpoint)
	m.updateManyEndpoint = mw(m.updateManyEndpoint)
	m.replaceOneEndpoint = mw(m.replaceOneEndpoint)
	m.aggregateEndpoint = mw(m.aggregateEndpoint)
	m.countDocumentsEndpoint = mw(m.countDocumentsEndpoint)
	m.findEndpoint = mw(m.findEndpoint)
	m.findOneEndpoint = mw(m.findOneEndpoint)
	m.findOneAndReplaceEndpoint = mw(m.findOneAndReplaceEndpoint)
	m.findOneAndDeleteEndpoint = mw(m.findOneAndDeleteEndpoint)
	m.findOneAndUpdateEndpoint = mw(m.findOneAndUpdateEndpoint)
	m.bulkWriteEndpoint = mw(m.bulkWriteEndpoint)
}

func (m *Collection) injectContext(ctx context.Context, operation string) context.Context {
	ctx = context.WithValue(ctx, constant2.MongoDbKey, m.dbName)
	ctx = context.WithValue(ctx, constant2.MongoCollectionKey, m.collectionName)
	ctx = context.WithValue(ctx, constant2.MongoOperationKey, operation)
	ctx = context.WithValue(ctx, constant2.SentinelBreaker, fmt.Sprintf("mongo.%s.%s", m.collectionName, operation))
	return ctx
}

func (m *Collection) BulkWrite(ctx context.Context, models []mongo.WriteModel,
	opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	_opts := make([]interface{}, 0, len(opts))
	for _, v := range opts {
		_opts = append(_opts, v)
	}
	_models := make([]interface{}, 0, len(models))
	for _, v := range models {
		_models = append(_models, v)
	}
	req := &mongoReq{opts: _opts, document: _models}
	ctx = m.injectContext(ctx, "BulkWrite")
	resp, err := m.bulkWriteEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*mongo.BulkWriteResult), nil
	//return m.coll.BulkWrite(ctx, models, opts...)
}

func (m *Collection) InsertOne(ctx context.Context, document interface{},
	opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	
	_opts := make([]interface{}, 0, len(opts))
	for _, v := range opts {
		_opts = append(_opts, v)
	}
	req := &mongoReq{document: []interface{}{document}, opts: _opts}
	ctx = m.injectContext(ctx, "InsertOne")
	resp, err := m.insertOneEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*mongo.InsertOneResult), nil
}

func (m *Collection) InsertMany(ctx context.Context, documents []interface{},
	opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	_opts := make([]interface{}, 0, len(opts))
	for _, v := range opts {
		_opts = append(_opts, v)
	}
	req := &mongoReq{document: documents, opts: _opts}
	ctx = m.injectContext(ctx, "InsertMany")
	resp, err := m.insertManyEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*mongo.InsertManyResult), nil
	//return m.coll.InsertMany(ctx, documents, opts...)
}

func (m *Collection) DeleteOne(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	_opts := make([]interface{}, 0, len(opts))
	for _, v := range opts {
		_opts = append(_opts, v)
	}
	req := &mongoReq{filter: []interface{}{filter}, opts: _opts}
	ctx = m.injectContext(ctx, "DeleteOne")
	resp, err := m.deleteOneEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*mongo.DeleteResult), nil
	//return m.coll.DeleteOne(ctx, filter, opts...)
}

func (m *Collection) DeleteMany(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	_opts := make([]interface{}, 0, len(opts))
	for _, v := range opts {
		_opts = append(_opts, v)
	}
	req := &mongoReq{filter: []interface{}{filter}, opts: _opts}
	ctx = m.injectContext(ctx, "DeleteMany")
	resp, err := m.deleteManyEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*mongo.DeleteResult), nil
	//return m.coll.DeleteMany(ctx, filter, opts...)
}

func (m *Collection) UpdateByID(ctx context.Context, id interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	_opts := make([]interface{}, 0, len(opts))
	for _, v := range opts {
		_opts = append(_opts, v)
	}
	req := &mongoReq{filter: []interface{}{id}, document: []interface{}{update}, opts: _opts}
	ctx = m.injectContext(ctx, "UpdateByID")
	resp, err := m.updateByIDEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*mongo.UpdateResult), nil
	//return m.coll.UpdateByID(ctx, id, update, opts...)
}

func (m *Collection) UpdateOne(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	_opts := make([]interface{}, 0, len(opts))
	for _, v := range opts {
		_opts = append(_opts, v)
	}
	req := &mongoReq{filter: []interface{}{filter}, document: []interface{}{update}, opts: _opts}
	ctx = m.injectContext(ctx, "UpdateOne")
	resp, err := m.updateOneEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*mongo.UpdateResult), nil
	//return m.coll.UpdateOne(ctx, filter, update, opts...)
}

func (m *Collection) UpdateMany(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	_opts := make([]interface{}, 0, len(opts))
	for _, v := range opts {
		_opts = append(_opts, v)
	}
	req := &mongoReq{filter: []interface{}{filter}, document: []interface{}{update}, opts: _opts}
	ctx = m.injectContext(ctx, "UpdateMany")
	resp, err := m.updateManyEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*mongo.UpdateResult), nil
	//return m.coll.UpdateMany(ctx, filter, update, opts...)
}

func (m *Collection) ReplaceOne(ctx context.Context, filter interface{},
	replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	_opts := make([]interface{}, 0, len(opts))
	for _, v := range opts {
		_opts = append(_opts, v)
	}
	req := &mongoReq{filter: []interface{}{filter}, document: []interface{}{replacement}, opts: _opts}
	ctx = m.injectContext(ctx, "ReplaceOne")
	resp, err := m.replaceOneEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*mongo.UpdateResult), nil
	//return m.coll.ReplaceOne(ctx, filter, replacement, opts...)
}
func getComment(ctx context.Context) string {
	return fmt.Sprintf("%s_%s", env.GetAppId(), mctx.GetTraceId(ctx))
}
func (m *Collection) Aggregate(ctx context.Context, pipeline interface{},
	opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	_opts := make([]interface{}, 0, len(opts))
	_opts = append(_opts, options.Aggregate().SetMaxTime(time.Second).SetMaxAwaitTime(time.Second).SetComment(getComment(ctx)))
	for _, v := range opts {
		_opts = append(_opts, v)
	}
	req := &mongoReq{filter: []interface{}{pipeline}, opts: _opts}
	ctx = m.injectContext(ctx, "Aggregate")
	resp, err := m.aggregateEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*mongo.Cursor), nil
}

func (m *Collection) CountDocuments(ctx context.Context, filter interface{},
	opts ...*options.CountOptions) (int64, error) {
	_opts := make([]interface{}, 0, len(opts))
	_opts = append(_opts, options.Count().SetMaxTime(time.Second))
	for _, v := range opts {
		_opts = append(_opts, v)
	}
	req := &mongoReq{filter: []interface{}{filter}, opts: _opts}
	ctx = m.injectContext(ctx, "CountDocuments")
	resp, err := m.countDocumentsEndpoint(ctx, req)
	if err != nil {
		return 0, err
	}
	return resp.(int64), nil
}

func (m *Collection) Find(ctx context.Context, filter interface{},
	opts ...*options.FindOptions) (cur *mongo.Cursor, err error) {
	_opts := make([]interface{}, 0, len(opts))
	//默认返回100个
	_opts = append(_opts, options.Find().SetMaxTime(time.Second).SetMaxAwaitTime(time.Second).SetComment(getComment(ctx)).SetLimit(100))
	for _, v := range opts {
		_opts = append(_opts, v)
	}
	req := &mongoReq{filter: []interface{}{filter}, opts: _opts}
	ctx = m.injectContext(ctx, "Find")
	resp, err := m.findEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*mongo.Cursor), nil
	//return m.coll.Find(ctx, filter, opts...)
}

func (m *Collection) FindOne(ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions) (*mongo.SingleResult, error) {
	_opts := make([]interface{}, 0, len(opts))
	_opts = append(_opts, options.FindOne().SetMaxTime(time.Second).SetComment(getComment(ctx)))
	for _, v := range opts {
		_opts = append(_opts, v)
	}
	req := &mongoReq{filter: []interface{}{filter}, opts: _opts}
	ctx = m.injectContext(ctx, "FindOne")
	resp, err := m.findOneEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*mongo.SingleResult), nil
	//return m.coll.FindOne(ctx, filter, opts...)
}

func (m *Collection) FindOneAndDelete(ctx context.Context, filter interface{},
	opts ...*options.FindOneAndDeleteOptions) (*mongo.SingleResult, error) {
	_opts := make([]interface{}, 0, len(opts))
	
	_opts = append(_opts, options.FindOneAndDelete().SetMaxTime(time.Second))
	for _, v := range opts {
		_opts = append(_opts, v)
	}
	req := &mongoReq{filter: []interface{}{filter}, opts: _opts}
	ctx = m.injectContext(ctx, "FindOneAndDelete")
	resp, err := m.findOneAndDeleteEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*mongo.SingleResult), nil
	//return m.coll.FindOneAndDelete(ctx, filter, opts...)
}

func (m *Collection) FindOneAndReplace(ctx context.Context, filter interface{},
	replacement interface{}, opts ...*options.FindOneAndReplaceOptions) (*mongo.SingleResult, error) {
	_opts := make([]interface{}, 0, len(opts))
	
	_opts = append(_opts, options.FindOneAndReplace().SetMaxTime(time.Second))
	for _, v := range opts {
		_opts = append(_opts, v)
	}
	req := &mongoReq{filter: []interface{}{filter}, document: []interface{}{replacement}, opts: _opts}
	ctx = m.injectContext(ctx, "FindOneAndReplace")
	resp, err := m.findOneAndReplaceEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*mongo.SingleResult), nil
	//return m.coll.FindOneAndReplace(ctx, filter, replacement, opts...)
}

func (m *Collection) FindOneAndUpdate(ctx context.Context, filter interface{},
	update interface{}, opts ...*options.FindOneAndUpdateOptions) (*mongo.SingleResult, error) {
	_opts := make([]interface{}, 0, len(opts))
	_opts = append(_opts, options.FindOneAndUpdate().SetMaxTime(time.Second))
	for _, v := range opts {
		_opts = append(_opts, v)
	}
	req := &mongoReq{filter: []interface{}{filter}, document: []interface{}{update}, opts: _opts}
	ctx = m.injectContext(ctx, "FindOneAndUpdate")
	resp, err := m.findOneAndUpdateEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*mongo.SingleResult), nil
	//return m.coll.FindOneAndUpdate(ctx, filter, update, opts...)
}

func (m *Collection) Indexes() mongo.IndexView {
	return m.coll.Indexes()
}
