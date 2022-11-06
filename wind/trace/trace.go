package trace

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/rabbitmq/amqp091-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/transport"
	"google.golang.org/grpc/metadata"
	"io"
	config2 "mond/wind/config"
	"mond/wind/utils"
	"mond/wind/utils/constant"
	mctx "mond/wind/utils/ctx"
	"strings"
	"time"
)

type Tracer struct {
	tracer opentracing.Tracer
	closer io.Closer
}

var (
	tracer *Tracer
)

func Init() {
	tc, err := config2.GetTraceConfig()
	utils.MustNil(err)
	env := config2.GetEnv()
	if env == "" {
		env = "local"
	}
	tracer, err = newTracer(fmt.Sprintf("%s_%s", config2.GetAppid(), env), tc.Endpoint, 1, tc.Protocol)
	utils.MustNil(err)
}
func GetTracer() *Tracer {
	if tracer == nil {
		panic("tracer is nil")
	}
	return tracer
}
func newTracer(serviceName string, addr string, ratio float64, protocol string) (*Tracer, error) {
	cfg := config.Configuration{
		ServiceName: serviceName,
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}
	sampler, _ := jaeger.NewProbabilisticSampler(ratio)
	var sender jaeger.Transport
	if protocol == "http" {
		sender = transport.NewHTTPTransport(addr)
	} else if protocol == "udp" {
		sender, _ = jaeger.NewUDPTransport(addr, 0)
		
	}
	
	reporter := jaeger.NewRemoteReporter(sender)
	
	tracer, closer, err := cfg.NewTracer(
		config.Reporter(reporter),
		config.Sampler(sampler),
	)
	if err != nil {
		return nil, err
	}
	t := &Tracer{
		tracer: tracer,
		closer: closer,
	}
	return t, nil
}

func (m *Tracer) StartRedisClientSpanFromContext(ctx context.Context, name string, opts ...opentracing.StartSpanOption) (context.Context, opentracing.Span, error) {
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	}
	sp := m.tracer.StartSpan(name, opts...)
	ext.SpanKindRPCClient.Set(sp)
	ext.DBType.Set(sp, "redis")
	return ctx, sp, nil
}

func (m *Tracer) StartGatewayServerSpanFromContext(ctx context.Context, name string, opts ...opentracing.StartSpanOption) (context.Context, opentracing.Span, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = make(map[string][]string)
	}
	md = md.Copy()
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	} else if spanCtx, err := m.tracer.Extract(opentracing.TextMap, metadataReaderWriter{md}); err == nil {
		opts = append(opts, opentracing.ChildOf(spanCtx))
	}
	
	sp := m.tracer.StartSpan(name, opts...)
	ctx = mctx.WithTraceId(ctx, sp)
	ext.SpanKindRPCServer.Set(sp)
	ext.Component.Set(sp, "gateway")
	ctx = opentracing.ContextWithSpan(ctx, sp)
	return ctx, sp, nil
}

func (m *Tracer) StartGrpcServerSpanFromContext(ctx context.Context, name string, opts ...opentracing.StartSpanOption) (context.Context, opentracing.Span, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = make(map[string][]string)
	}
	md = md.Copy()
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	} else if spanCtx, err := m.tracer.Extract(opentracing.TextMap, metadataReaderWriter{md}); err == nil {
		opts = append(opts, opentracing.ChildOf(spanCtx))
	}
	
	sp := m.tracer.StartSpan(name, opts...)
	ctx = mctx.WithTraceId(ctx, sp)
	ext.SpanKindRPCServer.Set(sp)
	ext.Component.Set(sp, "grpc")
	ctx = opentracing.ContextWithSpan(ctx, sp)
	return ctx, sp, nil
}

func (m *Tracer) StartGRpcClientSpanFromContext(ctx context.Context, name string, opts ...opentracing.StartSpanOption) (context.Context, opentracing.Span, error) {
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	}
	sp := m.tracer.StartSpan(name, opts...)
	//ctx = mctx.WithTraceId(ctx, sp)
	ext.SpanKindRPCClient.Set(sp)
	ext.Component.Set(sp, "grpc")
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	} else {
		md = md.Copy()
	}
	
	mdWriter := metadataReaderWriter{md}
	err := m.tracer.Inject(sp.Context(), opentracing.TextMap, mdWriter)
	if err != nil {
		return nil, nil, err
	}
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx, sp, nil
}

func (m *Tracer) StartConsumerSpanFromContext(ctx context.Context, name string, opts ...opentracing.StartSpanOption) (context.Context, opentracing.Span, error) {
	md := metadata.New(nil)
	headers, ok := ctx.Value(constant.ConsumerMdCtxKey).(amqp091.Table)
	if ok {
		for k, v := range headers {
			if k == constant.UberTraceIdKey {
				md[k] = []string{fmt.Sprintf("%v", v)}
			}
		}
	}
	if spanCtx, err := m.tracer.Extract(opentracing.TextMap, metadataReaderWriter{md}); err == nil {
		opts = append(opts, opentracing.ChildOf(spanCtx))
	}
	sp := m.tracer.StartSpan(name, opts...)
	ctx = mctx.WithTraceId(ctx, sp)
	ext.SpanKindConsumer.Set(sp)
	//ext.Component.Set(sp, "consumer")
	ctx = opentracing.ContextWithSpan(ctx, sp)
	return ctx, sp, nil
}

func (m *Tracer) StartPublishSpanFromContext(ctx context.Context, name string) (context.Context, opentracing.Span, error) {
	var opts []opentracing.StartSpanOption
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	}
	sp := GetTracer().tracer.StartSpan(name, opts...)
	carrier := opentracing.TextMapCarrier{}
	err := GetTracer().tracer.Inject(sp.Context(), opentracing.TextMap, carrier)
	if err != nil {
		return nil, nil, err
	}
	//ext.Component.Set(sp, "publish")
	ext.SpanKindProducer.Set(sp)
	ctx = context.WithValue(ctx, constant.PublishMdCtxKey, carrier)
	return ctx, sp, nil
}

func (m *Tracer) StartHttpServerSpanFromContext(ctx context.Context, name string, opts ...opentracing.StartSpanOption) (context.Context, opentracing.Span, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = make(map[string][]string)
	}
	md = md.Copy()
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	} else if spanCtx, err := m.tracer.Extract(opentracing.TextMap, metadataReaderWriter{md}); err == nil {
		opts = append(opts, opentracing.ChildOf(spanCtx))
	}
	
	sp := m.tracer.StartSpan(name, opts...)
	ctx = mctx.WithTraceId(ctx, sp)
	ext.SpanKindRPCServer.Set(sp)
	ext.Component.Set(sp, "http")
	ctx = opentracing.ContextWithSpan(ctx, sp)
	return ctx, sp, nil
}

func (m *Tracer) StartHttpClientSpanFromContext(ctx context.Context, name string, opts ...opentracing.StartSpanOption) (context.Context, opentracing.Span, error) {
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	}
	sp := m.tracer.StartSpan(name, opts...)
	ext.SpanKindRPCClient.Set(sp)
	ext.Component.Set(sp, "http")
	return ctx, sp, nil
}

func (m *Tracer) StartMongoSpanFromContext(ctx context.Context, name string, opts ...opentracing.StartSpanOption) (context.Context, opentracing.Span, error) {
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	}
	sp := m.tracer.StartSpan(name, opts...)
	ext.SpanKindRPCClient.Set(sp)
	ext.DBType.Set(sp, "mongodb")
	return ctx, sp, nil
}

func (m *Tracer) StartCronSpanFromContext(ctx context.Context, name string, opts ...opentracing.StartSpanOption) (context.Context, opentracing.Span, error) {
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	}
	sp := m.tracer.StartSpan(name, opts...)
	ctx = mctx.WithTraceId(ctx, sp)
	ext.SpanKindRPCServer.Set(sp)
	ext.Component.Set(sp, "cron")
	ctx = opentracing.ContextWithSpan(ctx, sp)
	return ctx, sp, nil
}

func (m *Tracer) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	return m.tracer.Extract(format, carrier)
}
func (m *Tracer) Inject(sm opentracing.SpanContext, format interface{}, carrier interface{}) error {
	return m.tracer.Inject(sm, format, carrier)
}

type metadataReaderWriter struct {
	metadata.MD
}

func (w metadataReaderWriter) Set(key, val string) {
	// The GRPC HPACK implementation rejects any uppercase keys here.
	//
	// As such, since the HTTP_HEADERS format is case-insensitive anyway, we
	// blindly lowercase the key (which is guaranteed to work in the
	// Inject/Extract sense per the OpenTracing spec).
	key = strings.ToLower(key)
	w.MD[key] = append(w.MD[key], val)
}

func (w metadataReaderWriter) ForeachKey(handler func(key, val string) error) error {
	for k, vals := range w.MD {
		for _, v := range vals {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}
	
	return nil
}
