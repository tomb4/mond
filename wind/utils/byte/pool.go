package byte

import (
	"bytes"
	"sync"
)

type BufferPool struct {
	pool sync.Pool
}

type BufferPoolOption struct {
	InitBufferSize int
}

func NewBufferPool(opt ...BufferPoolOption) *BufferPool {
	p := &BufferPool{pool: sync.Pool{}}
	initBufferSize := 4096
	for _, v := range opt {
		if v.InitBufferSize != 0 {
			initBufferSize = v.InitBufferSize
		}
	}
	p.pool.New = func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, initBufferSize))
	}
	return p
}

func (m *BufferPool) Get() *bytes.Buffer {
	return m.pool.Get().(*bytes.Buffer)
}

func (m *BufferPool) Put(v *bytes.Buffer) {
	v.Reset()
	m.pool.Put(v)
}
