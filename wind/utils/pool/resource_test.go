package pool

import (
	"context"
	"fmt"
	"testing"
	"time"
)

var (
	ctx = context.TODO()
)

func TestNewPool(t *testing.T) {
	p := NewPool(func(ctx context.Context) (res interface{}, err error) {
		return 1, nil
	}, func(res interface{}) {
		fmt.Println("close")
	}, 1)
	ctx,_ = context.WithTimeout(ctx, time.Second)
	fmt.Println(p.TryAcquire(ctx))
	fmt.Println(p.Stat())
	fmt.Println(p.Acquire(ctx))
}
