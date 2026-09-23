package gatekeeper

import (
	"context"
	"log"
	"sync"
	"time"
)

type delayedTasks struct {
	mu      sync.Mutex
	stopped bool
	cancel  chan struct{}
	done    chan struct{}
	wg      sync.WaitGroup
}

var verificationTasks = &delayedTasks{cancel: make(chan struct{}), done: make(chan struct{})}

func (d *delayedTasks) schedule(delay time.Duration, task func()) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.stopped {
		return
	}
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		defer func() {
			if r := recover(); r != nil {
				log.Println("延迟验证任务异常:", r)
			}
		}()
		timer := time.NewTimer(delay)
		defer timer.Stop()
		select {
		case <-d.cancel:
			return
		case <-timer.C:
		}
		select {
		case <-d.cancel:
			return
		default:
		}
		task()
	}()
}
func (d *delayedTasks) stop(ctx context.Context) error {
	d.mu.Lock()
	if !d.stopped {
		d.stopped = true
		close(d.cancel)
		go func() { d.wg.Wait(); close(d.done) }()
	}
	d.mu.Unlock()
	select {
	case <-d.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
func Stop(ctx context.Context) error { return verificationTasks.stop(ctx) }
