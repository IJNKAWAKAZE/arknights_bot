package shutdown

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
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

func run(ctx context.Context, steps []hook) error {
	for i, step := range steps {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := step(ctx); err != nil {
			return fmt.Errorf("关闭步骤 %d: %w", i+1, err)
		}
	}
	return nil
}

// All 必须在受管理的任务之外调用，例如 /kill 使用 go All()，避免等待自身结束。
// 等待超时时直接退出进程，避免提前关闭运行中任务仍在使用的存储连接。
func All() {
	once.Do(func() {
		log.Println("正在关闭...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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
			err = ctx.Err()
		}
		if err != nil {
			log.Println("关闭未完成:", err)
			os.Exit(1)
		}
		log.Println("已退出")
		os.Exit(0)
	})
}
