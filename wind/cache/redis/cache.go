package mredis

import (
    "context"
    "encoding/json"
    "errors"
    "github.com/go-redis/redis/v8"
    "go.uber.org/zap"
    "golang.org/x/sync/singleflight"
    cache "mond/wind/cache/define"
    "mond/wind/logger"
    mrand "mond/wind/utils/rand"
    "time"
)

const (
    defaultExpiry         = time.Minute * 10 //默认缓存时间10分钟
    defaultNotFoundExpiry = time.Minute      //默认未找到缓存时间 1分钟

    minExpiry = time.Second * 5

    deviationSeed = 5 //偏移种子 5%
)

var (
    notCacheErr = errors.New("no cache")
)

type redisCache struct {
    rds            *Client
    expiry         time.Duration //ms
    notFoundExpiry time.Duration //not found缓存时间
    errNotFound    error         //not found对应的错误
    s              singleflight.Group
    rand           *mrand.Rand
    _log           logger.Logger
}

type Option func(c *redisCache)

func WithExpiry(v time.Duration) Option {
    if v < minExpiry {
        v = minExpiry
    }
    return func(c *redisCache) {
        c.expiry = v
    }
}
func WithNotFoundExpiry(v time.Duration) Option {
    if v < minExpiry {
        v = minExpiry
    }
    return func(c *redisCache) {
        c.notFoundExpiry = v
    }
}

func NewRedisCache(rds *Client, errNotFound error, opts ...Option) cache.Cache {
    r := &redisCache{
        rds:            rds,
        errNotFound:    errNotFound,
        expiry:         defaultExpiry,
        notFoundExpiry: defaultNotFoundExpiry,
        rand:           mrand.NewRand(),
        _log:           logger.GetLogger(),
    }
    for _, o := range opts {
        o(r)
    }
    return r
}

func (m *redisCache) Take(ctx context.Context, key string, v interface{}, query func(v interface{}) error) error {
    result, err, _ := m.s.Do(key, func() (interface{}, error) {
        val, err := m.get(ctx, key)
        if err == nil {
            return val, nil
        }
        if err != notCacheErr {
            return nil, err
        }
        err = query(v)
        if err != nil && err != m.errNotFound {
            return "", err
        }
        if err == m.errNotFound {
            e := m.set(ctx, key, "*", m.randomDeviation(m.notFoundExpiry))
            if e != nil {
                m._log.Error(ctx, "redis cache not found写入失败", zap.Error(e))
            }
            return "", err
        }
        valStr := ""
        if err == nil {
            bytes, err := json.Marshal(v)
            if err != nil {
                return nil, err
            }
            valStr = string(bytes)
        }
        //如果查询结果是空数组，则说明没找到，则按照未找到的缓存时间去缓存
        if valStr == "[]" {
            e := m.set(ctx, key, valStr, m.randomDeviation(m.notFoundExpiry))
            if e != nil {
                m._log.Error(ctx, "redis cache 查询到空数组写入失败", zap.Error(e))
            }
        } else {
            e := m.set(ctx, key, valStr, m.randomDeviation(m.expiry))
            if e != nil {
                m._log.Error(ctx, "redis cache 查询成功，写入失败", zap.Error(e))
            }
        }
        return valStr, nil
    })
    if err != nil {
        return err
    }
    err = json.Unmarshal([]byte(result.(string)), v)
    if err != nil {
        m._log.Error(ctx, "redis cache反序列化失败，删除cache", zap.Error(err))
        m.Delete(ctx, key)
        return err
    }
    return nil
}

func (m *redisCache) get(ctx context.Context, key string) (string, error) {
    val, err := m.rds.Get(ctx, key).Result()
    if err != redis.Nil && err != nil {
        return "", err
    }
    if val == "*" {
        return "", m.errNotFound
    }
    if val == "" {
        return "", notCacheErr
    }

    return val, nil
}
func (m *redisCache) set(ctx context.Context, key string, val string, expiry time.Duration) error {
    e := m.rds.Set(ctx, key, val, expiry).Err()
    return e
}

func (m *redisCache) SetCache(ctx context.Context, key string, val string, expiry time.Duration) error {
    return m.set(ctx, key, val, m.randomDeviation(expiry))
}

func (m *redisCache) GetCache(ctx context.Context, key string, v interface{}) error {
    val, err := m.get(ctx, key)
    if err != nil {
        return err
    }
    err = json.Unmarshal([]byte(val), v)
    if err != nil {
        return err
    }
    return nil
}

func (m *redisCache) Delete(ctx context.Context, key ...string) error {
    return m.rds.Del(ctx, key...).Err()
}
func (m *redisCache) randomDeviation(v time.Duration) time.Duration {
    rd := v * time.Duration(m.rand.Int63n(deviationSeed)/100)
    //随机偏移量最小要有10S
    if rd < time.Second*10 {
        rd = time.Second * 10
    }
    return v + rd
}
