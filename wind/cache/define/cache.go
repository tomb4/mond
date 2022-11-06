package cache

import (
    "context"
    "time"
)

type Cache interface {
    Take(ctx context.Context, key string, v interface{}, query func(v interface{}) error) error
    Delete(ctx context.Context, key ...string) error
    SetCache(ctx context.Context, key string, val string, expiry time.Duration) error
    GetCache(ctx context.Context, key string, v interface{}) error
}
