package mongodb

import (
	"context"
	merr "mond/wind/err"
	"mond/wind/utils/endpoint"
)

func ctxDead(endpoint endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		select {
		case <-ctx.Done():
			return nil, merr.MongoContextTimeoutErr
		default:
			return endpoint(ctx, req)
		}
	}
}
