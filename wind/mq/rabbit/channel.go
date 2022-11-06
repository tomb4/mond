package rabbit

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"mond/wind/config"
	"mond/wind/logger"
	"mond/wind/sentry"
	"mond/wind/trace"
	"mond/wind/utils"
	"mond/wind/utils/constant"
	"time"
)

type channel struct {
	ch      *amqp091.Channel
	id      string
	_conn   *conn
	errChan chan *amqp091.Error
	stop    bool
}

func newChannel(c *conn) (*channel, error) {
	m := &channel{_conn: c, id: utils.GetNoDashUUIDStr()}
	err := m.createChannel()
	if err != nil {
		return nil, err
	}
	go m.loop()
	return m, nil
}

func (m *channel) close() {
	m.stop = true
	m.ch.Close()
}

func (m *channel) createChannel() error {
	ch, err := m._conn.c.Channel()
	if err != nil {
		return err
	}
	m.errChan = make(chan *amqp091.Error)
	ch.NotifyClose(m.errChan)
	m.ch = ch
	err = ch.Qos(int(m._conn.conf.Consumer.PrefetchCount), 0, false)
	//err = ch.Qos(100, 0, false)
	if err != nil {
		return err
	}
	return nil
}

func (m *channel) loop() {
	for {
		e := <-m.errChan
		//如果已经停了，则直接返回
		if m.stop {
			logger.GetLogger().Info(context.TODO(), "channel守护者收到错误，发现channel已经关闭，准备退出",
				zap.String("instance", m._conn.instance),
				zap.Any("err", e))
			return
		}
		logger.GetLogger().Info(context.TODO(), "channel守护者收到错误，准备开始重建channel",
			zap.String("instance", m._conn.instance),
			zap.Any("err", e))
		//logger.GetLogger().Info(context.TODO(), "重建channel", zap.String("instance", m._conn.instance))
		for {
			err := m.createChannel()
			if err != nil {
				if m.stop {
					logger.GetLogger().Info(context.TODO(), "channel守护者重建channel失败，发现channel已经关闭，准备退出",
						zap.String("instance", m._conn.instance),
						zap.Any("err", e))
					return
				}
				logger.GetLogger().Error(context.TODO(), "重建channel发生错误，准备1s后重试",
					zap.String("instance", m._conn.instance),
					zap.Any("err", err))
				//等待1s后重连
				time.Sleep(time.Second)
				continue
			}
			break
		}
	}
}

type Consume func(ctx context.Context, msg Message) error

type consumeMsg struct {
	*amqp091.Delivery
	ack bool
}

func (m *consumeMsg) Exchange() string {
	return m.Delivery.Exchange
}

func (m *consumeMsg) RoutingKey() string {
	return m.Delivery.RoutingKey
}

func (m *consumeMsg) Body() []byte {
	return m.Delivery.Body
}

func (m *consumeMsg) Nack(multiple, requeue bool) error {
	if m.ack {
		return nil
	}
	m.ack = true
	return m.Delivery.Nack(multiple, requeue)
}

func (m *consumeMsg) Ack(multiple bool) error {
	if m.ack {
		return nil
	}
	m.ack = true
	return m.Delivery.Ack(multiple)
}

func (m *channel) Consume(queue string, autoAck bool, handler Consume) {
	for {
	Loop:
		errChan := make(chan *amqp091.Error)
		msgs, err := m.ch.Consume(queue, fmt.Sprintf("%s-%s", config.GetAppid(), config.GetHostName()), autoAck, false, false, false, nil)
		if err != nil {
			close(errChan)
			if m.stop {
				logger.GetLogger().Info(context.TODO(), "consume创建失败后发现channel已经停止", zap.String("queue", queue), zap.Any("err", err))
				return
			}
			logger.GetLogger().Error(context.TODO(), "consume失败,等待1s后重试", zap.Any("err", err))
			time.Sleep(time.Second * 1)
			continue
		}
		for {
			select {
			case msg, ok := <-msgs:
				if !ok {
					if m.stop {
						logger.GetLogger().Info(context.TODO(), "consume接收msg失败了，channel已关停。准备退出", zap.Any("queue", queue))
						return
					}
					logger.GetLogger().Info(context.TODO(), "consume接收msg失败了，准备重建consume", zap.Any("queue", queue))
					goto Loop
				}
				ctx := context.TODO()
				ctx = context.WithValue(ctx, constant.ConsumerMdCtxKey, msg.Headers)
				if msg.Headers[constant.AsyncMethodCtxKey] != nil {
					ctx = context.WithValue(ctx, constant.AsyncMethodCtxKey, msg.Headers[constant.AsyncMethodCtxKey])
				}
				ctx = context.WithValue(ctx, constant.MqExchangeKey, msg.Exchange)
				_ = trace.ConsumerMiddleware(logger.ConsumerMiddleware(sentry.Consumer(func(ctx context.Context, queue string) error {
					mm := &consumeMsg{Delivery: &msg}
					err := handler(ctx, mm)
					//如果有错误且没有主动答复，则nack
					if !mm.ack {
						if err != nil {
							mm.Delivery.Nack(false, true)
						} else {
							mm.Delivery.Ack(false)
						}
					}
					if err != nil {
						logger.GetLogger().Error(ctx, "consumer返回错误", zap.String("queue", queue), zap.Any("err", err))
					}
					return err
				})))(ctx, queue)
			
			case e := <-errChan:
				//如果已经停了则直接返回
				if m.stop {
					logger.GetLogger().Info(context.TODO(), "consume接收到err channel后停止", zap.String("queue", queue), zap.Any("err", e))
					return
				}
				logger.GetLogger().Info(context.TODO(), "consume接收到err channel后准备重建监听", zap.String("queue", queue), zap.Any("err", e))
				goto Loop
			}
			
		}
		
	}
}

type PublishMessage interface {
	Exchange() string
	RoutingKey() string
	Body() []byte
}

type Message interface {
	Exchange() string
	RoutingKey() string
	Body() []byte
	Nack(multiple, requeue bool) error
	Ack(multiple bool) error
}

func (m *channel) Publish(ctx context.Context, msg PublishMessage) error {
	header := amqp091.Table{}
	name := fmt.Sprintf("producer.%s", msg.RoutingKey())
	if ctx.Value(constant.AsyncMethodCtxKey) != nil {
		name = fmt.Sprintf("producer.async.%s", ctx.Value(constant.AsyncMethodCtxKey))
		header[constant.AsyncMethodCtxKey] = ctx.Value(constant.AsyncMethodCtxKey)
	}
	ctx, span, err := trace.PublishMiddleware(ctx, name)
	if err == nil {
		defer span.Finish()
	}
	if ctx.Value(constant.AsyncMethodCtxKey) != nil {
		span.SetTag("async", "true")
	}
	span.SetTag("exchange", msg.Exchange())
	span.SetTag("routingKey", msg.RoutingKey())
	md := ctx.Value(constant.PublishMdCtxKey)
	textMap, ok := md.(opentracing.TextMapCarrier)
	if ok {
		for k, v := range textMap {
			header[k] = v
		}
	}
	err = m.ch.Publish(msg.Exchange(), msg.RoutingKey(), false, false, amqp091.Publishing{
		Headers: header,
		Body:    msg.Body(),
	})
	if err != nil {
		return err
	}
	return nil
}
