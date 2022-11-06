package rabbit

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"mond/wind/config"
	"mond/wind/logger"
	"mond/wind/utils/pool"
	"sync"
	"time"
)

type conn struct {
	c        *amqp091.Connection
	instance string
	errChan  chan *amqp091.Error
	stop     bool
	chs      sync.Map
	//producerPool chan *channel
	producerPool *pool.Pool
	producerInit bool
	conf         config.RabbitClientOption
}

func newConn(instance string) (*conn, error) {
	conf, err := config.GetRabbitClientOption(instance)
	if err != nil {
		return nil, err
	}
	c := conn{
		instance: instance,
		conf:     conf,
	}
	err = c.createConn()
	if err != nil {
		return nil, err
	}
	go c.loop()
	return &c, nil
}
func newAsyncConn() (*conn, error) {
	asyncConfig, err := config.GetAsyncOption()
	if err != nil {
		return nil, err
	}
	conf, err := config.GetRabbitClientOption(asyncConfig.Instance)
	if err != nil {
		return nil, err
	}
	conf.Consumer = asyncConfig.Consumer
	conf.Producer = asyncConfig.Producer
	c := conn{
		instance: "async",
		conf:     conf,
	}
	err = c.createConn()
	if err != nil {
		return nil, err
	}
	go c.loop()
	return &c, nil
}

func (m *conn) createConn() error {
	cc, err := amqp091.DialConfig(m.conf.Url, amqp091.Config{
		Heartbeat: defaultHeartbeat,
		Vhost:     m.conf.VHost,
	})
	if err != nil {
		return err
	}
	m.c = cc
	m.errChan = make(chan *amqp091.Error)
	m.c.NotifyClose(m.errChan)
	return nil
}
func (m *conn) loop() {
	for {
		e := <-m.errChan
		//如果已经停了，则直接返回
		if m.stop {
			return
		}
		if e != nil {
			logger.GetLogger().Error(context.TODO(), "rabbit loop chan error",
				zap.String("instance", m.instance),
				zap.Any("err", e))
		}
		logger.GetLogger().Info(context.TODO(), "重建连接", zap.String("instance", m.instance))
		for {
			err := m.createConn()
			if err != nil {
				logger.GetLogger().Error(context.TODO(), "rabbit loop recreate error",
					zap.String("instance", m.instance),
					zap.Any("err", err))
				//等待1s后重连
				time.Sleep(time.Second)
				continue
			}
			break
		}
	}
}

func (m *conn) Close() {
	m.stop = true
	m.chs.Range(func(key, value interface{}) bool {
		ch := value.(*channel)
		ch.close()
		return true
	})
	//for _, ch := range m.chs {
	//}
	m.c.Close()
}

func (m *conn) CreateExchange(name, kind string, durable, autoDelete bool) error {
	ch, err := m.c.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	err = ch.ExchangeDeclare(name, kind, durable, autoDelete, false, false, nil)
	return err
}

func (m *conn) CreateQueue(name string, durable, autoDelete bool) error {
	ch, err := m.c.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	_, err = ch.QueueDeclare(name, durable, autoDelete, false, false, nil)
	if err != nil {
		return err
	}
	return nil
}

func (m *conn) QueueBind(name, key, exchange string) error {
	ch, err := m.c.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	err = ch.QueueBind(name, key, exchange, false, nil)
	if err != nil {
		return err
	}
	return nil
}

func (m *conn) CreateChannel() (*channel, error) {
	c, err := newChannel(m)
	if err != nil {
		return nil, err
	}
	m.chs.Store(c.id, c)
	return c, nil
}
