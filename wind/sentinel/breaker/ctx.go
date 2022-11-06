package breaker

import (
	"context"
	"mond/wind/utils/constant"
)

func getScopeByCtx(ctx context.Context) string {
	v := ctx.Value(constant.SentinelBreaker)
	return v.(string)
}
