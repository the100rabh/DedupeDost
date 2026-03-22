package hasher

import (
	"context"
	"sync"
)

// WorkerPool manages a pool of workers for parallel hash calculation
type WorkerPool struct {
	workers    int
	hasher     *Hasher
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	jobChan    chan string
	resultChan chan HashResult
	mu         sync.RWMutex
	running    bool
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workers int, hasher *Hasher) *WorkerPool {
	if workers <= 0 {
		workers = 4
	}
	
	ctx, cancel := context.WithCancel(context.Background())
	
	return &WorkerPool{
		workers:    workers,
		hasher:     hasher,
		ctx:        ctx,
		cancel:     cancel,
		jobChan:    make(chan string, 100),
		resultChan: make(chan HashResult, 100),
	}
}

// Start starts the worker pool
func (wp *WorkerPool) Start() {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	
	if wp.running {
		return
	}
	
	wp.running = true
	
	// Start workers
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
	
	// Start result forwarder
	go wp.forwardResults()
}

// worker is a goroutine that processes hash jobs
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	
	for {
		select {
		case <-wp.ctx.Done():
			return
		case filePath, ok := <-wp.jobChan:
			if !ok {
				return
			}
			
			hash, err := wp.hasher.CalculateHashWithContext(filePath, wp.ctx)
			
			result := HashResult{
				Path: filePath,
				Hash: hash,
				Error: err,
			}
			
			if err != nil {
				if wp.ctx.Err() != nil {
					result.Cancelled = true
				}
			}
			
			select {
			case wp.resultChan <- result:
			case <-wp.ctx.Done():
				return
			}
		}
	}
}

// forwardResults forwards results from internal channel to external
func (wp *WorkerPool) forwardResults() {
	// This is a placeholder for more complex result handling
	// In a real implementation, you might want to add buffering,
	// batching, or other processing
}

// Submit submits a file for hashing
func (wp *WorkerPool) Submit(filePath string) bool {
	wp.mu.RLock()
	running := wp.running
	wp.mu.RUnlock()
	
	if !running {
		return false
	}
	
	select {
	case wp.jobChan <- filePath:
		return true
	case <-wp.ctx.Done():
		return false
	default:
		// Channel full, block until space available
		wp.jobChan <- filePath
		return true
	}
}

// SubmitBatch submits multiple files for hashing
func (wp *WorkerPool) SubmitBatch(files []string) int {
	submitted := 0
	for _, file := range files {
		if wp.Submit(file) {
			submitted++
		}
	}
	return submitted
}

// Results returns the result channel
func (wp *WorkerPool) Results() <-chan HashResult {
	return wp.resultChan
}

// Stop stops the worker pool
func (wp *WorkerPool) Stop() {
	wp.mu.Lock()
	if !wp.running {
		wp.mu.Unlock()
		return
	}
	wp.running = false
	wp.mu.Unlock()
	
	wp.cancel()
	close(wp.jobChan)
	wp.wg.Wait()
	close(wp.resultChan)
}

// IsRunning returns whether the pool is running
func (wp *WorkerPool) IsRunning() bool {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return wp.running
}

// WorkerCount returns the number of workers
func (wp *WorkerPool) WorkerCount() int {
	return wp.workers
}
