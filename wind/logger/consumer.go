package logger

import (
	"context"
	"mond/wind/mq/rabbit/function"
	"mond/wind/utils"
	"mond/wind/utils/constant"
	"time"

	"go.uber.org/zap"
)

func ConsumerMiddleware(f function.ConsumerFunc) function.ConsumerFunc {
	log := GetLogger()
	return func(ctx context.Context, queue string) (err error) {
		start := time.Now()
		defer func() {
			useTime := utils.UseTimeToStr(time.Now().Sub(start))
			fields := []zap.Field{
				zap.String("queue", queue),
				zap.String("type", "consumer"),
				zap.String("useTime", useTime),
			}
			if ctx.Value(constant.AsyncMethodCtxKey) != nil {
				fields = append(fields, zap.Any("method", ctx.Value(constant.AsyncMethodCtxKey)))
			}
			if err != nil {
				fields = append(fields, zap.Any("err", err))
				log.Error(ctx, "", fields...)
			} else {
				log.Info(ctx, "", fields...)
			}
		}()
		err = f(ctx, queue)
		return err
	}
}
