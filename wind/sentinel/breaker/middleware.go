package breaker

import (
	"context"
	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	merr "mond/wind/err"
	"mond/wind/logger"
	"mond/wind/utils/constant"
	"strings"
	"sync"
)

func MakeGrpcServerInterceptor() func(context.Context, interface{}, *grpc.UnaryServerInfo, grpc.UnaryHandler) (interface{}, error) {
	logger.GetLogger().Info(context.TODO(), "sentinel breaker server is open")
	ruleLoadMap := sync.Map{}
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		_scope := "server.grpc." + strings.Split(info.FullMethod, ".")[1]
		_, exists := ruleLoadMap.Load(_scope)
		//如果该规则不存在，则说明是第一次请求，则需要加载规则,加载规则在并发时可以多次，反正是覆盖，不影响。
		if !exists {
			rules := getRuleByScope(strings.Split(_scope, "."))
			circuitbreaker.LoadRulesOfResource(_scope, rules)
			ruleLoadMap.Store(_scope, "1")
		}
		token, e := api.Entry(_scope, api.WithTrafficType(base.Inbound))
		if e != nil {
			logger.GetLogger().Error(ctx, "grpc接口熔断", zap.Any("rule", e.TriggeredRule()), zap.String("_scope", _scope))
			return nil, merr.SysErrSentinelBreaker
		}
		defer token.Exit()
		res, err := handler(ctx, req)
		
		if err != nil {
			e := merr.ParseErrToMetaErr(err)
			//只有系统错误才算 而且如果是外部资源熔断或者上游熔断返回的错误，是不触发下游节点的熔断的
			if e.Code < merr.MaxSystemErr && e.Code != merr.ResourceErrSentinelBreaker.Code && e.Code != merr.SysErrSentinelBreaker.Code && e.Code != 401 {
				api.TraceError(token, e)
			}
		}
		return res, err
	}
}

func MakeGRpcClientMiddleware() func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	logger.GetLogger().Info(context.TODO(), "sentinel breaker client is open")
	ruleLoadMap := sync.Map{}
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		metadata := map[string]string{}
		ctx = context.WithValue(ctx, constant.GrpcClientAddr, metadata)
		err := invoker(ctx, method, req, reply, cc, opts...)
		addr := metadata["addr"]
		if addr == "" {
			if err != nil {
				return err
			}
			return nil
		}
		addr = strings.ReplaceAll(addr, ".", "_")
		_scope := "client.grpc." + addr
		_, exists := ruleLoadMap.Load(_scope)
		//如果该规则不存在，则说明是第一次请求，则需要加载规则,加载规则在并发时可以多次，反正是覆盖，不影响。
		if !exists {
			rules := getRuleByScope(strings.Split(_scope, "."))
			circuitbreaker.LoadRulesOfResource(_scope, rules)
			ruleLoadMap.Store(_scope, "1")
		}
		token, _ := api.Entry(_scope, api.WithTrafficType(base.Outbound))
		if err != nil {
			if token != nil {
				api.TraceError(token, err)
				token.Exit()
			}
			return err
		}
		return nil
	}
	
}
