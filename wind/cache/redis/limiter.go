package mredis

import (
	"context"
	"fmt"
	goCache "github.com/patrickmn/go-cache"
    "mond/wind/utils"
    "time"
)

type Limiter struct {
	rdb   *Client
	slot  int64 //时间片  ms  即每多少ms为1个时间片
	limit int64 //单位时间内的 限额
	cache *goCache.Cache
}

func NewLimiter(rdb *Client, slot int64, limit int64) *Limiter {
	limiter := Limiter{
		rdb:   rdb,
		slot:  slot,
		limit: limit,
		cache: goCache.New(time.Second*60, time.Second*10),
	}
	return &limiter
}

func (m *Limiter) Allow(ctx context.Context, key string) (pass bool, err error) {
	return m.AllowN(ctx, key, 1)
}

// pass 是否通过 true-没有触发限流  false-触发限流
// err  错误
// 对已经触发的限流设置内存缓存，可以减少对redis的压力
func (m *Limiter) AllowN(ctx context.Context, key string, n int64) (pass bool, err error) {
	current := utils.CurrentMillis()
	limitKey := fmt.Sprintf("p:l:%s:%d", key, current/m.slot)
	if _, exists := m.cache.Get(fmt.Sprintf("result:cache:%s", limitKey)); exists {
		return false, nil
	}
	var count int64
	if n != 1 {
		count, err = m.rdb.IncrBy(ctx, limitKey, n).Result()
	} else {
		count, err = m.rdb.Incr(ctx, limitKey).Result()
	}
	if err != nil {
		return true, err
	}
	_, exists := m.cache.Get(limitKey)
	if !exists {
		expire := time.Millisecond * time.Duration(m.slot) * 10
		if expire < time.Second {
			expire = time.Second
		}
		//键的过期时间设置为时间片的10倍
		_, err := m.rdb.Expire(ctx, limitKey, expire).Result()
		if err == nil {
			m.cache.SetDefault(limitKey, 1)
		}
	}
	//如果当前的数值已经超过限制
	if count > m.limit {
		//加入缓存
		m.cache.Set(fmt.Sprintf("result:cache:%s", limitKey), false, time.Millisecond*time.Duration(m.slot))
		return false, nil
	}

	return true, nil
}
