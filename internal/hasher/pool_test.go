package hasher_test

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/dupdel/dup-del/internal/hasher"
)

// TestNewWorkerPool tests NewWorkerPool
func TestNewWorkerPool(t *testing.T) {
	h := hasher.NewHasher()
	pool := hasher.NewWorkerPool(4, h)

	if pool == nil {
		t.Fatal("NewWorkerPool returned nil")
	}

	if pool.WorkerCount() != 4 {
		t.Errorf("WorkerCount = %d, want 4", pool.WorkerCount())
	}

	if pool.IsRunning() {
		t.Error("Pool should not be running before Start")
	}
}

// TestNewWorkerPool_ZeroWorkers tests NewWorkerPool with zero workers
func TestNewWorkerPool_ZeroWorkers(t *testing.T) {
	h := hasher.NewHasher()
	pool := hasher.NewWorkerPool(0, h)

	if pool.WorkerCount() != 4 {
		t.Errorf("WorkerCount with 0 = %d, want 4", pool.WorkerCount())
	}
}

// TestWorkerPool_Start tests Start
func TestWorkerPool_Start(t *testing.T) {
	h := hasher.NewHasher()
	pool := hasher.NewWorkerPool(2, h)

	if pool.IsRunning() {
		t.Error("Pool should not be running before Start")
	}

	pool.Start()

	if !pool.IsRunning() {
		t.Error("Pool should be running after Start")
	}

	pool.Stop()
}

// TestWorkerPool_Start_Twice tests starting twice
func TestWorkerPool_Start_Twice(t *testing.T) {
	h := hasher.NewHasher()
	pool := hasher.NewWorkerPool(2, h)

	pool.Start()
	pool.Start() // Should be no-op

	if !pool.IsRunning() {
		t.Error("Pool should be running")
	}

	pool.Stop()
}

// TestWorkerPool_Stop tests Stop
func TestWorkerPool_Stop(t *testing.T) {
	h := hasher.NewHasher()
	pool := hasher.NewWorkerPool(2, h)

	pool.Start()
	if !pool.IsRunning() {
		t.Fatal("Pool should be running")
	}

	pool.Stop()

	if pool.IsRunning() {
		t.Error("Pool should not be running after Stop")
	}
}

// TestWorkerPool_Stop_Twice tests stopping twice
func TestWorkerPool_Stop_Twice(t *testing.T) {
	h := hasher.NewHasher()
	pool := hasher.NewWorkerPool(2, h)

	pool.Start()
	pool.Stop()
	pool.Stop() // Should be no-op

	if pool.IsRunning() {
		t.Error("Pool should not be running")
	}
}

// TestWorkerPool_Submit tests Submit
func TestWorkerPool_Submit(t *testing.T) {
	h := hasher.NewHasher()
	pool := hasher.NewWorkerPool(2, h)

	pool.Start()
	defer pool.Stop()

	content := []byte("Test content")
	tmpFile := createTestFile(t, content)
	defer os.Remove(tmpFile)

	success := pool.Submit(tmpFile)
	if !success {
		t.Error("Submit should return true when pool is running")
	}

	// Wait for result
	timeout := time.After(5 * time.Second)
	select {
	case result, ok := <-pool.Results():
		if !ok {
			t.Error("Result channel closed unexpectedly")
		}
		if result.Error != nil {
			t.Errorf("Hash error: %v", result.Error)
		}
		if result.Hash == "" {
			t.Error("Hash should not be empty")
		}
	case <-timeout:
		t.Error("Timeout waiting for result")
	}
}

// TestWorkerPool_Submit_NotRunning tests Submit when not running
func TestWorkerPool_Submit_NotRunning(t *testing.T) {
	h := hasher.NewHasher()
	pool := hasher.NewWorkerPool(2, h)

	// Don't start the pool
	success := pool.Submit("/test/file.txt")
	if success {
		t.Error("Submit should return false when pool is not running")
	}
}

// TestWorkerPool_SubmitBatch tests SubmitBatch
func TestWorkerPool_SubmitBatch(t *testing.T) {
	h := hasher.NewHasher()
	pool := hasher.NewWorkerPool(2, h)

	pool.Start()
	defer pool.Stop()

	// Create test files
	files := make([]string, 5)
	for i := 0; i < 5; i++ {
		files[i] = createTestFile(t, []byte("content"+string(rune('0'+i))))
	}
	defer func() {
		for _, f := range files {
			os.Remove(f)
		}
	}()

	submitted := pool.SubmitBatch(files)
	if submitted != 5 {
		t.Errorf("SubmitBatch returned %d, want 5", submitted)
	}

	// Wait for all results
	results := make(map[string]string)
	timeout := time.After(5 * time.Second)
	for i := 0; i < 5; i++ {
		select {
		case result, ok := <-pool.Results():
			if !ok {
				t.Fatal("Result channel closed unexpectedly")
			}
			if result.Error != nil {
				t.Errorf("Hash error for %s: %v", result.Path, result.Error)
				continue
			}
			results[result.Path] = result.Hash
		case <-timeout:
			t.Fatalf("Timeout waiting for result %d", i)
		}
	}

	if len(results) != 5 {
		t.Errorf("Got %d results, want 5", len(results))
	}
}

// TestWorkerPool_Results tests Results channel
func TestWorkerPool_Results(t *testing.T) {
	h := hasher.NewHasher()
	pool := hasher.NewWorkerPool(2, h)

	pool.Start()
	defer pool.Stop()

	resultChan := pool.Results()
	if resultChan == nil {
		t.Fatal("Results returned nil channel")
	}

	content := []byte("Test content")
	tmpFile := createTestFile(t, content)
	defer os.Remove(tmpFile)

	pool.Submit(tmpFile)

	select {
	case result, ok := <-resultChan:
		if !ok {
			t.Error("Result channel closed unexpectedly")
		}
		if result.Path != tmpFile {
			t.Errorf("Result.Path = %s, want %s", result.Path, tmpFile)
		}
	case <-time.After(5 * time.Second):
		t.Error("Timeout waiting for result")
	}
}

// TestWorkerPool_ConcurrentSubmit tests concurrent submissions
func TestWorkerPool_ConcurrentSubmit(t *testing.T) {
	h := hasher.NewHasher()
	pool := hasher.NewWorkerPool(4, h)

	pool.Start()
	defer pool.Stop()

	// Create test files
	files := make([]string, 20)
	for i := 0; i < 20; i++ {
		files[i] = createTestFile(t, []byte("content"+string(rune('0'+i%10))))
	}
	defer func() {
		for _, f := range files {
			os.Remove(f)
		}
	}()

	// Submit from multiple goroutines
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			for j := start; j < start+4; j++ {
				pool.Submit(files[j])
			}
		}(i * 4)
	}

	wg.Wait()

	// Collect all results
	results := make(map[string]string)
	timeout := time.After(10 * time.Second)
	for i := 0; i < 20; i++ {
		select {
		case result, ok := <-pool.Results():
			if !ok {
				t.Fatal("Result channel closed unexpectedly")
			}
			if result.Error == nil {
				results[result.Path] = result.Hash
			}
		case <-timeout:
			t.Fatalf("Timeout waiting for result %d", i)
		}
	}

	if len(results) != 20 {
		t.Errorf("Got %d results, want 20", len(results))
	}
}

// TestWorkerPool_Stop_WithPendingJobs tests Stop with pending jobs
func TestWorkerPool_Stop_WithPendingJobs(t *testing.T) {
	h := hasher.NewHasher()
	pool := hasher.NewWorkerPool(1, h) // Single worker to slow down processing

	pool.Start()

	// Submit many files
	files := make([]string, 10)
	for i := 0; i < 10; i++ {
		files[i] = createTestFile(t, []byte("content"+string(rune('0'+i))))
	}
	defer func() {
		for _, f := range files {
			os.Remove(f)
		}
	}()

	for _, f := range files {
		pool.Submit(f)
	}

	// Stop immediately
	pool.Stop()

	// Should not panic or hang
}

// TestWorkerPool_WorkerCount tests WorkerCount
func TestWorkerPool_WorkerCount(t *testing.T) {
	h := hasher.NewHasher()

	tests := []struct {
		workers int
		want    int
	}{
		{0, 4},   // Default
		{1, 1},
		{2, 2},
		{8, 8},
		{-1, 4},  // Negative should default
	}

	for _, tt := range tests {
		pool := hasher.NewWorkerPool(tt.workers, h)
		if pool.WorkerCount() != tt.want {
			t.Errorf("WorkerCount(%d) = %d, want %d", tt.workers, pool.WorkerCount(), tt.want)
		}
	}
}

// TestWorkerPool_IsRunning tests IsRunning
func TestWorkerPool_IsRunning(t *testing.T) {
	h := hasher.NewHasher()
	pool := hasher.NewWorkerPool(2, h)

	if pool.IsRunning() {
		t.Error("IsRunning should return false before Start")
	}

	pool.Start()
	if !pool.IsRunning() {
		t.Error("IsRunning should return true after Start")
	}

	pool.Stop()
	if pool.IsRunning() {
		t.Error("IsRunning should return false after Stop")
	}
}

// TestHashResult_Cancelled tests HashResult with Cancelled flag
func TestHashResult_Cancelled(t *testing.T) {
	result := hasher.HashResult{
		Path:      "/test/file.txt",
		Hash:      "",
		Size:      0,
		Error:     context.Canceled,
		Cancelled: true,
	}

	if !result.Cancelled {
		t.Error("Cancelled should be true")
	}
	if result.Error == nil {
		t.Error("Error should not be nil for cancelled result")
	}
}
