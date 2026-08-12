package shutdown

import (
	"log"
	"os"
)

var hooks []func()

// Register 注册优雅退出时的清理函数
func Register(f func()) {
	hooks = append(hooks, f)
}

// All 依次执行清理函数后退出进程
func All() {
	log.Println("正在关闭...")
	for i := len(hooks) - 1; i >= 0; i-- {
		hooks[i]()
	}
	log.Println("已退出")
	os.Exit(0)
}
