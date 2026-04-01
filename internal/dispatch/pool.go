package dispatch

import (
	"context"
	"runtime"
	"sync"
)

// Pool is an in-process goroutine pool dispatcher that limits concurrency
// using a semaphore pattern.
type Pool struct {
	workers int
	sem     chan struct{}
}

// NewPool creates a Pool with the given worker count. If workers is zero or
// negative, runtime.NumCPU() is used.
func NewPool(workers int) *Pool {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}
	return &Pool{
		workers: workers,
		sem:     make(chan struct{}, workers),
	}
}

// Dispatch runs tasks concurrently up to the worker limit and returns a
// channel that receives one Result per task. The channel is closed when all
// tasks complete.
func (p *Pool) Dispatch(ctx context.Context, tasks []Task) <-chan Result {
	results := make(chan Result, len(tasks))
	var wg sync.WaitGroup

	for _, task := range tasks {
		wg.Add(1)
		go func(t Task) {
			defer wg.Done()

			select {
			case p.sem <- struct{}{}:
				defer func() { <-p.sem }()
			case <-ctx.Done():
				results <- Result{TaskID: t.ID, Err: ctx.Err()}
				return
			}

			var err error
			if t.Func != nil {
				err = t.Func(ctx)
			}
			results <- Result{TaskID: t.ID, Err: err}
		}(task)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

// Shutdown is a no-op for Pool as there is no persistent state to clean up.
func (p *Pool) Shutdown(_ context.Context) error {
	return nil
}
