package rabbit

import (
	"context"
	merr "mond/wind/err"
	"mond/wind/utils/pool"
)

//初始化producer
func (m *conn) InitProducer() error {
	maxChannel := 0
	if m.conf.Producer.MaxChannel > 0 {
		maxChannel = int(m.conf.Producer.MaxChannel)
	}
	m.producerPool = pool.NewPool(func(ctx context.Context) (res interface{}, err error) {
		ch, err := m.CreateChannel()
		return ch, err
	}, func(res interface{}) {
		c := res.(*channel)
		c.close()
		m.chs.Delete(c.id)
	}, int32(maxChannel))
	m.producerPool.CreateResource(context.TODO())
	m.producerInit = true
	return nil
}
func (m *conn) Publish(ctx context.Context, msg PublishMessage) error {
	if !m.producerInit {
		return merr.ProducerNotInitErr
	}
	if m.stop {
		return merr.ProducerStopErr
	}
	select {
	case <-ctx.Done():
		return merr.SysErrTimeoutErr.WithMsg("on rabbit publish")
	default:
	}
	//TODO:  暂时不考虑阻塞等情况，后面要优化
	resource, err := m.producerPool.TryAcquire(ctx)
	if err != nil {
		return err
	}
	defer func() {
		//如果发送过程中出现任何错误，则放弃当前channel 选择重建
		if err != nil {
			resource.Destroy()
		} else {
			resource.Release()
		}
	}()
	ch := resource.Value().(*channel)
	err = ch.Publish(ctx, msg)
	if err != nil {
		return err
	}
	return nil
}
