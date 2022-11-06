package mredis

import (
    "context"
    "fmt"
    "github.com/go-redis/redis/v8"
    "mond/wind/trace"
)

type hook struct {
    instance string
    appId    string
}

func (h *hook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
    key := ""
    if len(cmd.Args()) >= 2 {
        key = fmt.Sprintf("%s", cmd.Args()[1])
    }
    ctx, err := trace.RedisProcessBeforeMiddleware(ctx, cmd.Name(), key)
    return ctx, err
}

func (h *hook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
    err := trace.RedisProcessAfterMiddleware(ctx, cmd.Err())
    return err
}

func (h *hook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
    return ctx, nil
}

func (h *hook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
    return nil
}
