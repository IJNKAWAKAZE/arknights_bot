package shutdown

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const (
	// shutdownTimeout 是清理步骤的总预算。
	shutdownTimeout = 30 * time.Second
	// shutdownGrace 是超时之后额外留给资源释放的时间，避免浏览器、数据库连接被跳过。
	shutdownGrace = 10 * time.Second
)

type hook func(context.Context) error

var (
	hooks []hook
	mu    sync.Mutex
	once  sync.Once
)

// Register 按执行顺序注册清理函数：先停止产生任务，再等待任务结束，
// 最后关闭任务使用的资源。
func Register(f func(context.Context) error) { mu.Lock(); defer mu.Unlock(); hooks = append(hooks, f) }

// run 依次执行清理步骤，任何一步失败或超时都不会中断后续步骤：
// 前面的步骤出错时如果直接返回，浏览器、数据库这些资源就会被跳过。
func run(ctx context.Context, steps []hook) error {
	var errs []error
	for i, step := range steps {
		if err := step(ctx); err != nil {
			errs = append(errs, fmt.Errorf("关闭步骤 %d: %w", i+1, err))
		}
	}
	return errors.Join(errs...)
}

// All 必须在受管理的任务之外调用，例如 /kill 使用 go All()，避免等待自身结束。
// 总预算用完后还会多等一小段时间，让最后几步的资源释放有机会跑完。
func All() {
	once.Do(func() {
		log.Println("正在关闭...")
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		mu.Lock()
		steps := append([]hook(nil), hooks...)
		mu.Unlock()
		result := make(chan error, 1)
		go func() { result <- run(ctx, steps) }()
		var err error
		select {
		case err = <-result:
		case <-ctx.Done():
			select {
			case err = <-result:
			case <-time.After(shutdownGrace):
				err = ctx.Err()
			}
		}
		if err != nil {
			log.Println("关闭未完成:", err)
			os.Exit(1)
		}
		log.Println("已退出")
		os.Exit(0)
	})
}
