package mredis

import (
	"context"
	merr "mond/wind/err"
	utils2 "mond/wind/utils"
	"time"
)

const (
	compareAndDeleteLuaScript = `
		local val = redis.call("get",KEYS[1])
		if (val == ARGV[1]) then
			redis.call("del",KEYS[1])
			return "ok"
		else
			return "fail"
		end
`
	minLockExpire = time.Millisecond * 100
)

var (
//h          = sha1.New()
//_, _       = io.WriteString(h, compareAndDeleteLuaScript)
//hashScript = hex.EncodeToString(h.Sum(nil))
)

type Lock struct {
	uuid string
	rdb  *Client
	key  string
}

func (m *Client) Lock(ctx context.Context, key string, expire time.Duration) (*Lock, error) {
	if expire < minLockExpire {
		expire = minLockExpire
	}
	uuid := utils2.GetNoDashUUIDStr()
	lock := &Lock{
		uuid: uuid,
		rdb:  m,
		key:  key,
	}
	max := 0
	for {
		select {
		case <-ctx.Done():
			return nil, merr.SysErrTimeoutErr
		default:
		}
		ok, err := m.SetNX(ctx, key, uuid, expire).Result()
		if err != nil {
			return nil, err
		}
		if ok {
			break
		}
		max++
		if max > 10 {
			return nil, merr.RedisLockTimeoutErr
		}
		//随机sleep 20-29ms
		time.Sleep(time.Millisecond * (time.Duration(m.int64(10)) + 20))
	}
	return lock, nil
}

func (m *Lock) UnLock(ctx context.Context) error {
	_, err := m.rdb.EvalSha(ctx, m.rdb.lockScriptSha, []string{m.key}, []string{m.uuid}).Result()
	if err != nil {
		return err
	}
	return nil
}
