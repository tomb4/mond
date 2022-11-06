package static

import (
	"google.golang.org/grpc/resolver"
)

const scheme = "static"

type staticResolverBuilder struct {
	addrMap map[string][]string
}

func (m *staticResolverBuilder) Scheme() string {
	return scheme
}
func NewStaticResolverBuilder(addrMap map[string][]string) *staticResolverBuilder {
	return &staticResolverBuilder{addrMap: addrMap}
}

func (m *staticResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &staticResolver{
		target:  target,
		cc:      cc,
		addrMap: m.addrMap,
	}
	r.start()
	return r, nil
}

type staticResolver struct {
	target  resolver.Target
	cc      resolver.ClientConn
	addrMap map[string][]string
}

func (m *staticResolver) start() {
	addrArr := m.addrMap[m.target.URL.Host]
	addrs := make([]resolver.Address, 0, len(addrArr))
	for _, v := range addrArr {
		addrs = append(addrs, resolver.Address{Addr: v})
	}
	m.cc.UpdateState(resolver.State{Addresses: addrs})
}

func (*staticResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*staticResolver) Close()                                  {}
