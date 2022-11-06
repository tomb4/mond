package utils

import (
	"context"
	"mond/wind/utils/constant"
)

//注入目标下游的target
func InjectCtxTargetHostName(ctx context.Context, targetHostName string) context.Context {
	return context.WithValue(ctx, constant.BalanceTargetHostName, targetHostName)
}

func ExtractCtxTargetHostName(ctx context.Context) string {
	host, _ := ctx.Value(constant.BalanceTargetHostName).(string)
	return host
}
