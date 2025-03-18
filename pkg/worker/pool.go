package worker

import (
	"context"
	"sync"
	"time"
)

// Task represents a function to be executed by the worker pool
type Task func() error

// Pool represents a worker pool that can execute tasks concurrently
type Pool struct {
	tasks     chan Task
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
	semaphore chan struct{} // Limits the number of concurrent tasks
}

// NewPool creates a new worker pool with the specified number of workers
func NewPool(size int) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	
	pool := &Pool{
		tasks:     make(chan Task, size*10), // Buffer tasks to avoid blocking
		ctx:       ctx,
		cancel:    cancel,
		semaphore: make(chan struct{}, size),
	}
	
	// Start the dispatcher
	go pool.dispatcher()
	
	return pool
}

// dispatcher listens for tasks and dispatches them to workers
func (p *Pool) dispatcher() {
	for {
		select {
		case <-p.ctx.Done():
			return
		case task, ok := <-p.tasks:
			if !ok {
				return
			}
			
			// Acquire semaphore slot
			select {
			case <-p.ctx.Done():
				return
			case p.semaphore <- struct{}{}:
				// Start a worker
				p.wg.Add(1)
				go func(t Task) {
					defer p.wg.Done()
					defer func() { <-p.semaphore }() // Release semaphore slot
					
					// Execute the task with timeout
					taskCtx, cancel := context.WithTimeout(p.ctx, 30*time.Second)
					defer cancel()
					
					done := make(chan error, 1)
					go func() {
						done <- t()
					}()
					
					select {
					case <-taskCtx.Done():
						// Task timed out or pool was shut down
						return
					case <-done:
						// Task completed
						return
					}
				}(task)
			}
		}
	}
}

// Submit adds a task to the pool
func (p *Pool) Submit(task Task) {
	select {
	case <-p.ctx.Done():
		return // Pool is shutting down
	case p.tasks <- task:
		// Task submitted successfully
	}
}

// Shutdown gracefully shuts down the pool, waiting for all tasks to complete
func (p *Pool) Shutdown() {
	p.cancel() // Signal all workers to stop
	close(p.tasks) // Stop accepting new tasks
	p.wg.Wait() // Wait for all workers to finish
}

// ShutdownWithTimeout attempts to gracefully shut down the pool within the timeout
func (p *Pool) ShutdownWithTimeout(timeout time.Duration) bool {
	p.cancel() // Signal all workers to stop
	close(p.tasks) // Stop accepting new tasks
	
	// Wait with timeout
	c := make(chan struct{})
	go func() {
		defer close(c)
		p.wg.Wait()
	}()
	
	select {
	case <-c:
		return true // Completed successfully
	case <-time.After(timeout):
		return false // Timed out
	}
}
