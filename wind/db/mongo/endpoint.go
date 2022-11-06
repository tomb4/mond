package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	merr "mond/wind/err"
	"mond/wind/utils/endpoint"
)

func (m *Collection) makeInsertOneEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		_mongoReq := req.(*mongoReq)
		document := _mongoReq.document[0]
		opts := _mongoReq.opts
		_opts := make([]*options.InsertOneOptions, 0, len(opts))
		for _, v := range opts {
			_opts = append(_opts, v.(*options.InsertOneOptions))
		}
		result, err := m.coll.InsertOne(ctx, document, _opts...)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}

func (m *Collection) makeInsertManyEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		_mongoReq := req.(*mongoReq)
		opts := _mongoReq.opts
		_opts := make([]*options.InsertManyOptions, 0, len(opts))
		for _, v := range opts {
			_opts = append(_opts, v.(*options.InsertManyOptions))
		}
		result, err := m.coll.InsertMany(ctx, _mongoReq.document, _opts...)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}

func (m *Collection) makeDeleteOneEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		_mongoReq := req.(*mongoReq)
		opts := _mongoReq.opts
		filter := _mongoReq.filter[0]
		_opts := make([]*options.DeleteOptions, 0, len(opts))
		for _, v := range opts {
			_opts = append(_opts, v.(*options.DeleteOptions))
		}
		result, err := m.coll.DeleteOne(ctx, filter, _opts...)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}

func (m *Collection) makeDeleteManyEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		_mongoReq := req.(*mongoReq)
		opts := _mongoReq.opts
		filter := _mongoReq.filter[0]
		_opts := make([]*options.DeleteOptions, 0, len(opts))
		for _, v := range opts {
			_opts = append(_opts, v.(*options.DeleteOptions))
		}
		result, err := m.coll.DeleteMany(ctx, filter, _opts...)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}

func (m *Collection) makeUpdateByIDEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		_mongoReq := req.(*mongoReq)
		opts := _mongoReq.opts
		filter := _mongoReq.filter[0]
		update := _mongoReq.document[0]
		_opts := make([]*options.UpdateOptions, 0, len(opts))
		for _, v := range opts {
			_opts = append(_opts, v.(*options.UpdateOptions))
		}
		result, err := m.coll.UpdateByID(ctx, filter, update, _opts...)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}

func (m *Collection) makeUpdateOneEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		_mongoReq := req.(*mongoReq)
		opts := _mongoReq.opts
		filter := _mongoReq.filter[0]
		update := _mongoReq.document[0]
		_opts := make([]*options.UpdateOptions, 0, len(opts))
		for _, v := range opts {
			_opts = append(_opts, v.(*options.UpdateOptions))
		}
		result, err := m.coll.UpdateOne(ctx, filter, update, _opts...)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}

func (m *Collection) makeUpdateManyEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		_mongoReq := req.(*mongoReq)
		opts := _mongoReq.opts
		filter := _mongoReq.filter[0]
		update := _mongoReq.document[0]
		_opts := make([]*options.UpdateOptions, 0, len(opts))
		for _, v := range opts {
			_opts = append(_opts, v.(*options.UpdateOptions))
		}
		result, err := m.coll.UpdateMany(ctx, filter, update, _opts...)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}

func (m *Collection) makeReplaceOneEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		_mongoReq := req.(*mongoReq)
		opts := _mongoReq.opts
		filter := _mongoReq.filter[0]
		update := _mongoReq.document[0]
		_opts := make([]*options.ReplaceOptions, 0, len(opts))
		for _, v := range opts {
			_opts = append(_opts, v.(*options.ReplaceOptions))
		}
		result, err := m.coll.ReplaceOne(ctx, filter, update, _opts...)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}

func (m *Collection) makeAggregateEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		_mongoReq := req.(*mongoReq)
		opts := _mongoReq.opts
		pipeline := _mongoReq.filter[0]
		_opts := make([]*options.AggregateOptions, 0, len(opts))
		for _, v := range opts {
			_opts = append(_opts, v.(*options.AggregateOptions))
		}
		result, err := m.coll.Aggregate(ctx, pipeline, _opts...)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}

func (m *Collection) makeCountDocumentsEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		_mongoReq := req.(*mongoReq)
		opts := _mongoReq.opts
		filter := _mongoReq.filter[0]
		_opts := make([]*options.CountOptions, 0, len(opts))
		for _, v := range opts {
			_opts = append(_opts, v.(*options.CountOptions))
		}
		result, err := m.coll.CountDocuments(ctx, filter, _opts...)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}

func (m *Collection) makeFindEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		_mongoReq := req.(*mongoReq)
		opts := _mongoReq.opts
		filter := _mongoReq.filter[0]
		_opts := make([]*options.FindOptions, 0, len(opts))
		for _, v := range opts {
			_opts = append(_opts, v.(*options.FindOptions))
		}
		result, err := m.coll.Find(ctx, filter, _opts...)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}

func (m *Collection) makeFindOneEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		_mongoReq := req.(*mongoReq)
		opts := _mongoReq.opts
		filter := _mongoReq.filter[0]
		_opts := make([]*options.FindOneOptions, 0, len(opts))
		for _, v := range opts {
			_opts = append(_opts, v.(*options.FindOneOptions))
		}
		result := m.coll.FindOne(ctx, filter, _opts...)
		if result == nil {
			return nil, merr.MongoFindOneResultNilErr
		}
		err := result.Err()
		if err != nil {
			//if err == mongo.ErrNoDocuments {
			//	return nil, merr.NoDocumentError
			//}
			return nil, err
		}
		return result, nil
	}
}

func (m *Collection) makeFindOneAndReplaceEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		_mongoReq := req.(*mongoReq)
		opts := _mongoReq.opts
		filter := _mongoReq.filter[0]
		replacement := _mongoReq.document[0]
		_opts := make([]*options.FindOneAndReplaceOptions, 0, len(opts))
		for _, v := range opts {
			_opts = append(_opts, v.(*options.FindOneAndReplaceOptions))
		}
		result := m.coll.FindOneAndReplace(ctx, filter, replacement, _opts...)
		if result == nil {
			return nil, merr.MongoFindOneResultNilErr
		}
		err := result.Err()
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}
func (m *Collection) makeFindOneAndDeleteEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		_mongoReq := req.(*mongoReq)
		opts := _mongoReq.opts
		filter := _mongoReq.filter[0]
		_opts := make([]*options.FindOneAndDeleteOptions, 0, len(opts))
		for _, v := range opts {
			_opts = append(_opts, v.(*options.FindOneAndDeleteOptions))
		}
		result := m.coll.FindOneAndDelete(ctx, filter, _opts...)
		if result == nil {
			return nil, merr.MongoFindOneResultNilErr
		}
		err := result.Err()
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}

func (m *Collection) makeFindOneAndUpdateEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		_mongoReq := req.(*mongoReq)
		opts := _mongoReq.opts
		filter := _mongoReq.filter[0]
		update := _mongoReq.document[0]
		_opts := make([]*options.FindOneAndUpdateOptions, 0, len(opts))
		for _, v := range opts {
			_opts = append(_opts, v.(*options.FindOneAndUpdateOptions))
		}
		result := m.coll.FindOneAndUpdate(ctx, filter, update, _opts...)
		if result == nil {
			return nil, merr.MongoFindOneResultNilErr
		}
		err := result.Err()
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}
func (m *Collection) makeBulkWriteEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		_mongoReq := req.(*mongoReq)
		opts := _mongoReq.opts
		_opts := make([]*options.BulkWriteOptions, 0, len(opts))
		for _, v := range opts {
			_opts = append(_opts, v.(*options.BulkWriteOptions))
		}
		models := make([]mongo.WriteModel, 0, len(_mongoReq.document))
		for _, v := range _mongoReq.document {
			models = append(models, v.(mongo.WriteModel))
		}
		result, err := m.coll.BulkWrite(ctx, models, _opts...)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}
