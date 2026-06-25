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
	workers  int
	tasksC   chan Task[I, T]
	resultsC chan Result[I, R]
	handler  func(ctx context.Context, task Task[I, T]) (R, error)
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
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
	p.cancel()
	p.wg.Wait()
	close(p.tasksC)
	close(p.resultsC)
}

func (p *Pool[I, T, R]) Start() {
	p.wg.Add(p.workers)
	for i := 0; i < p.workers; i++ {
		go func() {
			defer p.wg.Done()

			for {
				select {
				case task := <-p.tasksC:
					result, err := p.handler(p.ctx, task)
					p.resultsC <- Result[I, R]{
						TaskID: task.ID,
						Result: result,
						Err:    err,
					}
				case <-p.ctx.Done():
					return
				}
			}
		}()
	}
}

func (p *Pool[I, T, R]) Send(task Task[I, T]) bool {
	select {
	case <-p.ctx.Done():
		return false
	default:
	}

	select {
	case p.tasksC <- task:
		return true
	case <-p.ctx.Done():
		return false
	}
}

func (p *Pool[I, T, R]) Receive() <-chan Result[I, R] {
	return p.resultsC
}
