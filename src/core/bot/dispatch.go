package bot

import (
	"log"
	"runtime/debug"
	"sync"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

// maxConcurrentHandlers 限制同时执行的指令数量，避免大量截图任务同时占用浏览器和内存。
const maxConcurrentHandlers = 6

var (
	queueMu    sync.Mutex
	chatQueues = make(map[int64]*chatQueue)
	workerSem  = make(chan struct{}, maxConcurrentHandlers)
)

// chatQueue 保证同一个会话内的指令按顺序串行执行，不同会话之间互不影响。
// 消息库的 Run 循环本身是串行的，一个慢指令会堵住所有用户，这里把它改成异步分发。
type chatQueue struct {
	mu      sync.Mutex
	tasks   []func()
	running bool
}

// chatIDOf 取出更新所属的会话/用户，用于隔离不同会话的执行顺序。
func chatIDOf(update tgbotapi.Update) int64 {
	switch {
	case update.Message != nil && update.Message.Chat != nil:
		return update.Message.Chat.ID
	case update.GuestMessage != nil && update.GuestMessage.Chat != nil:
		return update.GuestMessage.Chat.ID
	case update.CallbackQuery != nil:
		if update.CallbackQuery.Message != nil && update.CallbackQuery.Message.Chat != nil {
			return update.CallbackQuery.Message.Chat.ID
		}
		if update.CallbackQuery.From != nil {
			return update.CallbackQuery.From.ID
		}
	case update.InlineQuery != nil && update.InlineQuery.From != nil:
		return update.InlineQuery.From.ID
	case update.ChatJoinRequest != nil:
		return update.ChatJoinRequest.Chat.ID
	case update.MyChatMember != nil:
		return update.MyChatMember.Chat.ID
	}
	return 0
}

// async 把处理器放到后台执行，消息循环可以立刻继续处理下一个更新。
func async(handler func(tgbotapi.Update) error) func(tgbotapi.Update) error {
	return func(update tgbotapi.Update) error {
		enqueue(chatIDOf(update), func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Recovered err=%v, stack=%s\n", r, string(debug.Stack()))
				}
			}()
			workerSem <- struct{}{}
			defer func() { <-workerSem }()
			if err := handler(update); err != nil {
				log.Println("Plugin Error", err.Error())
			}
		})
		return nil
	}
}

func enqueue(id int64, task func()) {
	queueMu.Lock()
	q := chatQueues[id]
	if q == nil {
		q = &chatQueue{}
		chatQueues[id] = q
	}
	queueMu.Unlock()

	q.mu.Lock()
	q.tasks = append(q.tasks, task)
	start := !q.running
	if start {
		q.running = true
	}
	q.mu.Unlock()

	if start {
		go q.run()
	}
}

func (q *chatQueue) run() {
	for {
		q.mu.Lock()
		if len(q.tasks) == 0 {
			q.running = false
			q.mu.Unlock()
			return
		}
		task := q.tasks[0]
		q.tasks[0] = nil
		q.tasks = q.tasks[1:]
		q.mu.Unlock()

		task()
	}
}
