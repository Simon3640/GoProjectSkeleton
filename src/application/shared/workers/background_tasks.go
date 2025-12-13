// Package workers provides background tasks for the application.
package workers

import (
	"context"
	"errors"
	"log"
	"runtime/debug"
	"sync"
)

// BackgroundTask is the unit submitted to the BackgroundExecutor.
type BackgroundTask func(ctx context.Context)

// BackgroundExecutor executes tasks in background with a configurable worker pool.
// It accepts a parent context so tasks can be cancelled when appropriate.
type BackgroundExecutor struct {
	tasks        chan BackgroundTask
	workers      int
	wg           sync.WaitGroup
	startOnce    sync.Once
	stopOnce     sync.Once
	shutdown     chan struct{}
	started      bool
	mu           sync.Mutex
	parentCtx    context.Context
	parentCancel context.CancelFunc
}

// NewBackgroundExecutor creates a new executor. If workers <= 0 it defaults to 4.
func NewBackgroundExecutor(parent context.Context, workers int, queueSize int) *BackgroundExecutor {
	if workers <= 0 {
		workers = 4
	}
	if queueSize <= 0 {
		queueSize = 100
	}
	ctx, cancel := context.WithCancel(parent)
	return &BackgroundExecutor{
		tasks:        make(chan BackgroundTask, queueSize),
		workers:      workers,
		shutdown:     make(chan struct{}),
		parentCtx:    ctx,
		parentCancel: cancel,
	}
}

// Start brings up worker goroutines. Safe to call multiple times.
func (b *BackgroundExecutor) Start() {
	b.startOnce.Do(func() {
		for i := 0; i < b.workers; i++ {
			b.wg.Add(1)
			go func(workerID int) {
				defer b.wg.Done()
				for {
					select {
					case <-b.parentCtx.Done():
						return
					case <-b.shutdown:
						return
					case t, ok := <-b.tasks:
						if !ok {
							return
						}
						// execute with panic recovery
						func() {
							defer func() {
								if r := recover(); r != nil {
									log.Printf("panic recovered in background worker %d: %v\n%s", workerID, r, debug.Stack())
								}
							}()
							t(b.parentCtx)
						}()
					}
				}
			}(i)
		}
		b.mu.Lock()
		b.started = true
		b.mu.Unlock()
	})
}

// Submit enqueues a task; returns error if executor is shutting down.
func (b *BackgroundExecutor) Submit(task BackgroundTask) error {
	b.mu.Lock()
	started := b.started
	b.mu.Unlock()
	if !started {
		// start lazily
		b.Start()
	}
	select {
	case <-b.parentCtx.Done():
		return b.parentCtx.Err()
	case b.tasks <- task:
		return nil
	default:
		// queue full — decide policy. Here we return an error.
		return errors.New("background queue full")
	}
}

// Wait waits until all currently queued tasks are finished.
func (b *BackgroundExecutor) Wait() {
	// close tasks channel to signal workers to exit when queue is drained
	b.stopOnce.Do(func() {
		close(b.shutdown)
	})
	// Wait until channel drained and workers exited
	// Note: we don't close(b.tasks) here because multiple submitters may exist. Instead
	// we rely on parent cancellation or explicit Stop() to stop accepting tasks.
	b.wg.Wait()
}

// Stop cancels parent context and waits for workers to finish their current jobs.
func (b *BackgroundExecutor) Stop() {
	b.stopOnce.Do(func() {
		b.parentCancel()
		close(b.tasks)
		// wait for workers
		b.wg.Wait()
	})
}

var (
	backgroundExecutorSingleton *BackgroundExecutor
	singletonOnce               sync.Once
	singletonMu                 sync.Mutex
)

// InitializeBackgroundExecutor initializes the singleton BackgroundExecutor instance.
// This should be called during application startup in the infrastructure layer.
// This function is thread-safe and can only be called once.
func InitializeBackgroundExecutor(parent context.Context, workers int, queueSize int) {
	singletonOnce.Do(func() {
		backgroundExecutorSingleton = NewBackgroundExecutor(parent, workers, queueSize)
		backgroundExecutorSingleton.Start()
	})
}

// GetBackgroundExecutor returns the singleton instance of BackgroundExecutor.
// The executor must be initialized first by calling InitializeBackgroundExecutor.
// This function is thread-safe.
func GetBackgroundExecutor() *BackgroundExecutor {
	return backgroundExecutorSingleton
}

// ResetBackgroundExecutorSingleton resets the singleton instance.
// This is primarily useful for testing purposes.
// WARNING: Only call this when you're sure no other goroutines are using the executor.
func ResetBackgroundExecutorSingleton() {
	singletonMu.Lock()
	defer singletonMu.Unlock()
	if backgroundExecutorSingleton != nil {
		backgroundExecutorSingleton.Stop()
		backgroundExecutorSingleton = nil
		singletonOnce = sync.Once{}
	}
}
