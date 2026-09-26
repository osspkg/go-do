/*
 *  Copyright (c) 2024-2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package workerpool

import (
	"context"
	"sync"
)

type Task[I comparable, T any] struct {
	ID   I
	Data T
}

type Result[I comparable, T any] struct {
	TaskID I
	Result T
	Err    error
}

type Pool[I comparable, T, R any] struct {
	workers   int
	tasksC    chan Task[I, T]
	resultsC  chan Result[I, R]
	handler   func(ctx context.Context, task Task[I, T]) (R, error)
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	sendWg    sync.WaitGroup
	mu        sync.Mutex
	closeOnce sync.Once
	started   bool
	closed    bool
}

func New[I comparable, T, R any](
	workers int,
	handler func(ctx context.Context, task Task[I, T]) (R, error),
) *Pool[I, T, R] {
	workers = max(workers, 1)

	ctx, cancel := context.WithCancel(context.Background())

	return &Pool[I, T, R]{
		workers:  workers,
		tasksC:   make(chan Task[I, T], workers),
		resultsC: make(chan Result[I, R], workers),
		handler:  handler,
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (p *Pool[I, T, R]) Close() {
	p.closeOnce.Do(func() {
		p.mu.Lock()
		p.closed = true
		p.cancel()
		p.mu.Unlock()

		p.sendWg.Wait()
		p.wg.Wait()
		close(p.resultsC)
	})
}

func (p *Pool[I, T, R]) Start() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.started || p.closed {
		return
	}
	p.started = true

	p.wg.Add(p.workers)
	for i := 0; i < p.workers; i++ {
		go func() {
			defer p.wg.Done()

			for {
				select {
				case task := <-p.tasksC:
					result, err := p.handler(p.ctx, task)
					select {
					case p.resultsC <- Result[I, R]{
						TaskID: task.ID,
						Result: result,
						Err:    err,
					}:
					case <-p.ctx.Done():
						return
					}
				case <-p.ctx.Done():
					return
				}
			}
		}()
	}
}

func (p *Pool[I, T, R]) Send(task Task[I, T]) bool {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return false
	}
	p.sendWg.Add(1)
	p.mu.Unlock()
	defer p.sendWg.Done()

	select {
	case p.tasksC <- task:
		return p.ctx.Err() == nil
	case <-p.ctx.Done():
		return false
	}
}

func (p *Pool[I, T, R]) Receive() <-chan Result[I, R] {
	return p.resultsC
}
