package rgroup

// Package rgroup provides synchronization, error propagation, and Context
// rgroup 包为一组子任务的 goroutine 提供了 goroutine 同步,错误取消功能.
//
//rgroup 包含三种常用方式
//
//1、直接使用 此时不会因为一个任务失败导致所有任务被 cancel:
//		g := &rgroup.Group{}
//		g.Go(func(ctx context.Context) {
//			// NOTE: 此时 ctx 为 context.Background()
//			// do something
//		})
//
//2、WithContext 使用 WithContext 时不会因为一个任务失败导致所有任务被 cancel:
//		g := rgroup.WithContext(ctx)
//		g.Go(func(ctx context.Context) {
//			// NOTE: 此时 ctx 为 rgroup.WithContext 传递的 ctx
//			// do something
//		})
//
//3、WithCancel 使用 WithCancel 时如果有一个人任务失败会导致所有*未进行或进行中*的任务被 cancel:
//		g := rgroup.WithCancel(ctx)
//		g.Go(func(ctx context.Context) {
//			// NOTE: 此时 ctx 是从 rgroup.WithContext 传递的 ctx 派生出的 ctx
//			// do something
//		})
//
//设置最大并行数 GOMAXPROCS 对以上三种使用方式均起效
//NOTE: 由于 rgroup 实现问题,设定 GOMAXPROCS 的 rgroup 需要立即调用 Wait() 例如:
//
//		g := rgroup.WithCancel(ctx)
//		g.GOMAXPROCS(2)
//		// task1
//		g.Go(func(ctx context.Context) {
//			fmt.Println("task1")
//		})
//		// task2
//		g.Go(func(ctx context.Context) {
//			fmt.Println("task2")
//		})
//		// task3
//		g.Go(func(ctx context.Context) {
//			fmt.Println("task3")
//		})
//		// NOTE: 此时设置的 GOMAXPROCS 为2, 添加了三个任务 task1, task2, task3 此时 task3 是不会运行的!
//		// 只有调用了 Wait task3 才有运行的机会
//		g.Wait() // task3 运行
