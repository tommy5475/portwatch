package daemon

import (
	"context"
	"sync"
)

// workerPool runs a fixed number of concurrent tasks using a semaphore
// to cap parallelism. It is intended for parallel port-scan batches.
type workerPool struct {
	sem *semaphore
}

// newWorkerPool creates a pool limited to concurrency parallel workers.
// If concurrency < 1 it defaults to 1.
func newWorkerPool(concurrency int) *workerPool {
	return &workerPool{sem: newSemaphore(concurrency)}
}

// RunAll executes fn for every item in items, blocking until all
// invocations complete or ctx is cancelled. Errors from individual
// workers are collected and returned as a slice; a nil slice means all
// tasks succeeded.
func (p *workerPool) RunAll(ctx context.Context, items []int, fn func(ctx context.Context, item int) error) []error {
	var (
		mu   sync.Mutex
		errs []error
		wg   sync.WaitGroup
	)

	for _, item := range items {
		item := item
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := p.sem.Acquire(ctx); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				return
			}
			defer p.sem.Release()

			if err := fn(ctx, item); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	return errs
}
