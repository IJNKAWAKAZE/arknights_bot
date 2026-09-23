package bot

import (
	"arknights_bot/config"
	"context"
	"errors"
	"log"
	"runtime/debug"
	"sync"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

var (
	errQueueFull         = errors.New("消息处理队列已满，请稍后重试")
	errDispatcherStopped = errors.New("机器人正在关闭")
	commands             = newDispatcher(6, 32, 256)
	moderation           = newDispatcher(2, 32, 128)
	notices              = newDispatcher(1, 1, 16)
)

type chatQueue struct {
	tasks       []func()
	outstanding int
}

// 使用同一把锁保护队列的查找、入队和回收，避免任务加入已脱离管理的队列。
// 容量限制包含正在执行的任务。
type dispatcher struct {
	mu                               sync.Mutex
	queues                           map[int64]*chatQueue
	sem                              chan struct{}
	perChat, totalLimit, outstanding int
	stopped                          bool
	stopCh                           chan struct{}
	wg                               sync.WaitGroup
	done                             chan struct{}
}

func newDispatcher(workers, perChat, total int) *dispatcher {
	return &dispatcher{queues: make(map[int64]*chatQueue), sem: make(chan struct{}, workers),
		perChat: perChat, totalLimit: total, stopCh: make(chan struct{}), done: make(chan struct{})}
}

func (d *dispatcher) enqueue(id int64, task func()) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.stopped {
		return errDispatcherStopped
	}
	q := d.queues[id]
	if d.outstanding >= d.totalLimit || (q != nil && q.outstanding >= d.perChat) {
		return errQueueFull
	}
	start := q == nil
	if start {
		q = &chatQueue{}
		d.queues[id] = q
	}
	q.tasks = append(q.tasks, task)
	q.outstanding++
	d.outstanding++
	if start {
		d.wg.Add(1)
		go d.run(id, q)
	}
	return nil
}

func (d *dispatcher) run(id int64, q *chatQueue) {
	defer d.wg.Done()
	for {
		d.mu.Lock()
		if d.stopped || len(q.tasks) == 0 {
			d.outstanding -= len(q.tasks)
			delete(d.queues, id)
			d.mu.Unlock()
			return
		}
		task := q.tasks[0]
		q.tasks[0] = nil
		q.tasks = q.tasks[1:]
		d.mu.Unlock()

		select {
		case d.sem <- struct{}{}:
			d.mu.Lock()
			stopped := d.stopped
			d.mu.Unlock()
			if !stopped {
				runTask(task)
			}
			<-d.sem
		case <-d.stopCh:
		}
		d.mu.Lock()
		q.outstanding--
		d.outstanding--
		d.mu.Unlock()
	}
}

func runTask(task func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("处理消息异常: %v\n%s", r, debug.Stack())
		}
	}()
	task()
}

// stopAdmission 停止接收新任务，丢弃排队任务，让已开始的处理继续执行。
func (d *dispatcher) stopAdmission() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.stopped {
		return
	}
	d.stopped = true
	close(d.stopCh)
	go func() { d.wg.Wait(); close(d.done) }()
}

func (d *dispatcher) wait(ctx context.Context) error {
	select {
	case <-d.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func StopAdmission() {
	commands.stopAdmission()
	moderation.stopAdmission()
	notices.stopAdmission()
}
func Drain(ctx context.Context) error {
	if err := commands.wait(ctx); err != nil {
		return err
	}
	if err := moderation.wait(ctx); err != nil {
		return err
	}
	return notices.wait(ctx)
}

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

func async(handler func(tgbotapi.Update) error) func(tgbotapi.Update) error {
	return commands.wrap(handler)
}
func asyncControl(handler func(tgbotapi.Update) error) func(tgbotapi.Update) error {
	return moderation.wrap(handler)
}
func (d *dispatcher) wrap(handler func(tgbotapi.Update) error) func(tgbotapi.Update) error {
	return func(update tgbotapi.Update) error {
		err := d.enqueue(chatIDOf(update), func() {
			if err := handler(update); err != nil {
				log.Println("Plugin Error", err)
			}
		})
		if errors.Is(err, errQueueFull) {
			// 繁忙提示也使用有容量上限的队列，避免阻塞 Telegram 消息接收循环。
			_ = notices.enqueue(chatIDOf(update), func() {
				if update.CallbackQuery != nil {
					update.CallbackQuery.Answer(false, "请求较多，请稍后重试。")
				} else if update.Message != nil && update.Message.IsCommand() {
					_, sendErr := config.Arknights.ReplyText(update.Message.Chat.ID, update.Message.MessageID, "请求较多，请稍后重试。")
					if sendErr != nil {
						log.Println("发送繁忙提示失败:", sendErr)
					}
				}
			})
		}
		return err
	}
}
