package balancer

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"mond/wind/logger"
	"mond/wind/sentinel"
	"mond/wind/utils/constant"
	"strings"
	"sync/atomic"
)

const (
	Strict   = "Strict"
	Tolerate = "Tolerate"
)

var (
	targetBalanceErr = errors.New("target balance not support")
)

type MetaBalancerBuilder struct {
}

func NewBalancerBuilder() (*MetaBalancerBuilder, error) {
	return &MetaBalancerBuilder{
	}, nil
}

func (m *MetaBalancerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	//fmt.Println("MetaBalancerBuilder", info)
	//fmt.Println(info.ReadySCs)
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	hostTargetSubConn := map[string][]*conn{}
	labelTargetSubConn := map[string][]*conn{}
	for k, v := range info.ReadySCs {
		hostName, _ := v.Address.BalancerAttributes.Value("hostName").(string)
		label, _ := v.Address.BalancerAttributes.Value("label").(string)
		//fmt.Println(info.ReadySCs)
		hostTargetSubConn[hostName] = append(hostTargetSubConn[hostName], &conn{c: k, addr: v.Address.Addr})
		labelTargetSubConn[label] = append(labelTargetSubConn[label], &conn{c: k, addr: v.Address.Addr})
	}
	b := &metaBalancer{
		hostTargetSubConn:  hostTargetSubConn,
		labelTargetSubConn: labelTargetSubConn,
		//labelTargetSubConnLen: int64(len(labelTargetSubConn["stable"])),
	}
	//if b.labelTargetSubConnLen == 0 {
	//	b.labelTargetSubConnLen = 1
	//}
	//fmt.Println(b)
	return b
}

type conn struct {
	c    balancer.SubConn
	addr string
}
type metaBalancer struct {
	//以host为目标选择条件
	hostTargetSubConn map[string][]*conn
	//以label为目标选择条件
	labelTargetSubConn map[string][]*conn
	//labelTargetSubConnLen int64
	next int64
}

func (m *metaBalancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	hostName, ok := info.Ctx.Value(constant.BalanceTargetHostName).(string)
	next := atomic.AddInt64(&m.next, 1)
	//fmt.Println("Pick", hostName)
	if ok && hostName != "" {
		subConn := m.hostTargetSubConn[hostName]
		//fmt.Println(subConn)
		if len(subConn) == 0 {
			return balancer.PickResult{}, errors.New("sub conn is empty")
		}
		conn := subConn[next%int64(len(subConn))]
		return balancer.PickResult{SubConn: conn.c}, nil
	}
	label, ok := info.Ctx.Value(constant.BalanceTargetLabel).(string)
	if !ok || label == "" {
		label = "stable"
	}
	subConn := m.labelTargetSubConn[label]
	//fmt.Println("labelTargetSubConn", subConn)
	if len(subConn) == 0 {
		if label == "stable" {
			return balancer.PickResult{}, errors.New("sub conn is empty")
		}
		subConn = m.labelTargetSubConn["stable"]
		if len(subConn) == 0 {
			return balancer.PickResult{}, errors.New("sub conn is empty")
		}
	}
	conn := subConn[next%int64(len(subConn))]
	//如果超过50%的节点处于熔断状态，则直接随便选一个
	max := len(subConn)/2 + 1
	//max := len(subConn)
	for i := 0; i < max; i++ {
		//如果当前下游pod处于熔断状态，则忽略
		if !sentinel.ResourceInFuse(fmt.Sprintf("client.grpc.%s", strings.ReplaceAll(conn.addr, ".", "_"))) {
			break
		}
		next = next + 1
		conn = subConn[next%int64(len(subConn))]
		if max == i+1 {
			logger.GetLogger().Error(info.Ctx, "所有的下游都处于熔断状态", zap.String("addr", conn.addr))
		}
	}
	metadata := info.Ctx.Value(constant.GrpcClientAddr)
	if metadata != nil {
		m := metadata.(map[string]string)
		m["addr"] = conn.addr
	}
	return balancer.PickResult{SubConn: conn.c}, nil
}
