package breaker

import (
	"context"
	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"go.uber.org/zap"
	merr "mond/wind/err"
	"mond/wind/logger"
	"mond/wind/utils/endpoint"
	"strings"
	"sync"
)

const (
	_mongo_max_time_limit_err = "operation exceeded time limit"
)

// scope 作用域，指定熔断器的基础设置，
// scope采用三层，即 type.name.method
// type代表大的种类，比如 server/mongo/mysql/redis/grpc/http
// name代表具体的实例 比如当type=mongo时，name代表collection 如user post等
// method代表具体的方法或表，即熔断的最小单位， 如果是http，则method可以是具体的方法， 如果是mongo，则method是具体的语句 如：find findOne
func Middleware(endpoint endpoint.Endpoint) endpoint.Endpoint {
	ruleLoadMap := sync.Map{}
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_scope := getScopeByCtx(ctx)
		_, exists := ruleLoadMap.Load(_scope)
		//如果该规则不存在，则说明是第一次请求，则需要加载规则,加载规则在并发时可以多次，反正是覆盖，不影响。
		if !exists {
			rules := getRuleByScope(strings.Split(_scope, "."))
			circuitbreaker.LoadRulesOfResource(_scope, rules)
			ruleLoadMap.Store(_scope, "1")
		}
		token, e := api.Entry(_scope, api.WithTrafficType(base.Outbound))
		if e != nil {
			logger.GetLogger().Error(ctx, "外部资源熔断", zap.Any("rule", e.TriggeredRule()), zap.String("_scope", _scope))
			return nil, merr.ResourceErrSentinelBreaker
		}
		defer token.Exit()
		res, err := endpoint(ctx, request)
		if err != nil {
			scopes := strings.Split(_scope, ".")
			//如果是mongo，则对执行耗时超过最大限制的请求标记为错误
			if scopes[0] == Type_Mongo && err.Error() == _mongo_max_time_limit_err {
				api.TraceError(token, e)
			}
			//TODO: 陆续还要支持其它外部资源的熔断
		}
		return res, err
	}
}
