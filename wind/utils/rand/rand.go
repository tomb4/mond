package mrand

import (
    "math/rand"
    "sync"
    "time"
)

type Rand struct {
    r    *rand.Rand
    lock sync.Mutex
}

func NewRand() *Rand {
    m := Rand{
        r: rand.New(rand.NewSource(time.Now().UnixNano())),
    }
    return &m
}

func (m *Rand) Int63n(n int64) int64 {
    m.lock.Lock()
    defer m.lock.Unlock()
    return m.r.Int63n(n)
}

func (m *Rand) Int31n(n int32) int32 {
    m.lock.Lock()
    defer m.lock.Unlock()
    return m.r.Int31n(n)
}
