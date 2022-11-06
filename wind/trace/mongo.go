package trace

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"go.mongodb.org/mongo-driver/mongo"
	"mond/wind/utils/constant"
	"mond/wind/utils/endpoint"
)

func MongoMiddleware(endpoint endpoint.Endpoint) endpoint.Endpoint {
	tracer := GetTracer()
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		dbName := ctx.Value(constant.MongoDbKey)
		collectionName := ctx.Value(constant.MongoCollectionKey)
		operation := ctx.Value(constant.MongoOperationKey)
		ctx, span, _ := tracer.StartMongoSpanFromContext(ctx, fmt.Sprintf("mongo.%v.%v.%v", dbName, collectionName, operation))
		defer span.Finish()
		resp, err := endpoint(ctx, req)
		//查不到数据这种数据正常业务报错，不记录
		if err != nil && err != mongo.ErrNoDocuments {
			ext.Error.Set(span, true)
			span.LogFields(
				log.String("event", "error"),
				log.String("stack", err.Error()),
			)
		}
		return resp, err
	}
}
