package rabbit

type Options struct {
	pool int
}

type Option func(opts *Options)

func (m *conn) Consume(queue string, autoAck bool, handler Consume, opts ...Option) error {
	opt := &Options{pool: 4}
	if m.conf.Consumer.MaxChannel != 0 {
		opt.pool = int(m.conf.Consumer.MaxChannel)
	}
	for _, v := range opts {
		v(opt)
	}
	for i := 0; i < opt.pool; i++ {
		ch, err := m.CreateChannel()
		if err != nil {
			return err
		}
		go ch.Consume(queue, autoAck, handler)
	}
	return nil
}
