package scanner_test

import (
	"sync"
	"testing"
	"time"

	"github.com/dupdel/dup-del/internal/scanner"
)

// TestNewProgressTracker tests NewProgressTracker
func TestNewProgressTracker(t *testing.T) {
	pt := scanner.NewProgressTracker()
	if pt == nil {
		t.Fatal("NewProgressTracker returned nil")
	}

	progress := pt.GetProgress()
	if progress.FilesScanned != 0 {
		t.Errorf("Initial FilesScanned = %d, want 0", progress.FilesScanned)
	}
	if progress.TotalFiles != 0 {
		t.Errorf("Initial TotalFiles = %d, want 0", progress.TotalFiles)
	}
	if progress.PercentComplete != 0 {
		t.Errorf("Initial PercentComplete = %f, want 0", progress.PercentComplete)
	}
}

// TestProgressTracker_AddCallback tests AddCallback
func TestProgressTracker_AddCallback(t *testing.T) {
	pt := scanner.NewProgressTracker()

	called := false
	pt.AddCallback(func(progress scanner.ScanProgress) {
		called = true
	})

	pt.Update(1, 10, "/test/file.txt", 1024)

	if !called {
		t.Error("Callback should have been called")
	}
}

// TestProgressTracker_AddCallback_Multiple tests multiple callbacks
func TestProgressTracker_AddCallback_Multiple(t *testing.T) {
	pt := scanner.NewProgressTracker()

	var calls []int
	pt.AddCallback(func(progress scanner.ScanProgress) {
		calls = append(calls, 1)
	})
	pt.AddCallback(func(progress scanner.ScanProgress) {
		calls = append(calls, 2)
	})
	pt.AddCallback(func(progress scanner.ScanProgress) {
		calls = append(calls, 3)
	})

	pt.Update(1, 10, "/test/file.txt", 1024)

	if len(calls) != 3 {
		t.Errorf("Expected 3 callback calls, got %d", len(calls))
	}
}

// TestProgressTracker_Update tests Update
func TestProgressTracker_Update(t *testing.T) {
	pt := scanner.NewProgressTracker()

	pt.Update(5, 10, "/test/file.txt", 5120)

	progress := pt.GetProgress()
	if progress.FilesScanned != 5 {
		t.Errorf("FilesScanned = %d, want 5", progress.FilesScanned)
	}
	if progress.TotalFiles != 10 {
		t.Errorf("TotalFiles = %d, want 10", progress.TotalFiles)
	}
	if progress.CurrentFile != "/test/file.txt" {
		t.Errorf("CurrentFile = %s, want /test/file.txt", progress.CurrentFile)
	}
	if progress.TotalSize != 5120 {
		t.Errorf("TotalSize = %d, want 5120", progress.TotalSize)
	}
	if progress.PercentComplete != 50.0 {
		t.Errorf("PercentComplete = %f, want 50.0", progress.PercentComplete)
	}
}

// TestProgressTracker_Update_PercentComplete tests percent complete calculation
func TestProgressTracker_Update_PercentComplete(t *testing.T) {
	pt := scanner.NewProgressTracker()

	tests := []struct {
		scanned    int64
		total      int64
		expectedPct float64
	}{
		{0, 10, 0},
		{1, 10, 10},
		{5, 10, 50},
		{10, 10, 100},
		{25, 100, 25},
		{50, 200, 25},
	}

	for _, tt := range tests {
		pt.Update(tt.scanned, tt.total, "/test/file.txt", 1024)
		progress := pt.GetProgress()
		if progress.PercentComplete != tt.expectedPct {
			t.Errorf("Update(%d, %d): PercentComplete = %f, want %f",
				tt.scanned, tt.total, progress.PercentComplete, tt.expectedPct)
		}
	}
}

// TestProgressTracker_Update_ZeroTotalFiles tests with zero total files
func TestProgressTracker_Update_ZeroTotalFiles(t *testing.T) {
	pt := scanner.NewProgressTracker()

	pt.Update(5, 0, "/test/file.txt", 1024)

	progress := pt.GetProgress()
	if progress.PercentComplete != 0 {
		t.Errorf("PercentComplete with zero total = %f, want 0", progress.PercentComplete)
	}
}

// TestProgressTracker_GetProgress tests GetProgress
func TestProgressTracker_GetProgress(t *testing.T) {
	pt := scanner.NewProgressTracker()

	pt.Update(10, 100, "/test/file.txt", 10240)

	progress := pt.GetProgress()
	if progress.FilesScanned != 10 {
		t.Errorf("FilesScanned = %d, want 10", progress.FilesScanned)
	}
}

// TestProgressTracker_GetElapsed tests GetElapsed
func TestProgressTracker_GetElapsed(t *testing.T) {
	pt := scanner.NewProgressTracker()

	// Initial elapsed should be very small
	elapsed := pt.GetElapsed()
	if elapsed < 0 {
		t.Errorf("GetElapsed = %v, should be positive", elapsed)
	}

	// Wait a bit and check again
	time.Sleep(100 * time.Millisecond)
	elapsed2 := pt.GetElapsed()
	if elapsed2 < elapsed {
		t.Error("Elapsed time should increase")
	}
}

// TestProgressTracker_GetFilesPerSecond tests GetFilesPerSecond
func TestProgressTracker_GetFilesPerSecond(t *testing.T) {
	pt := scanner.NewProgressTracker()

	// Initially should be 0
	rate := pt.GetFilesPerSecond()
	if rate != 0 {
		t.Errorf("Initial GetFilesPerSecond = %f, want 0", rate)
	}

	// Update with some files
	pt.Update(100, 1000, "/test/file.txt", 102400)

	// Should be positive
	rate = pt.GetFilesPerSecond()
	if rate <= 0 {
		t.Errorf("GetFilesPerSecond after update = %f, should be positive", rate)
	}
}

// TestProgressTracker_GetETA tests GetETA
func TestProgressTracker_GetETA(t *testing.T) {
	pt := scanner.NewProgressTracker()

	// Initially should be 0
	eta := pt.GetETA()
	if eta != 0 {
		t.Errorf("Initial GetETA = %v, want 0", eta)
	}

	// Update with partial progress
	pt.Update(50, 100, "/test/file.txt", 51200)

	// Should estimate some remaining time
	eta = pt.GetETA()
	if eta < 0 {
		t.Errorf("GetETA = %v, should be positive", eta)
	}
}

// TestProgressTracker_GetETA_EdgeCases tests GetETA edge cases
func TestProgressTracker_GetETA_EdgeCases(t *testing.T) {
	pt := scanner.NewProgressTracker()

	// Zero total files
	pt.Update(10, 0, "/test/file.txt", 1024)
	eta := pt.GetETA()
	if eta != 0 {
		t.Errorf("GetETA with zero total = %v, want 0", eta)
	}

	// Zero files scanned
	pt2 := scanner.NewProgressTracker()
	pt2.Update(0, 100, "/test/file.txt", 1024)
	eta = pt2.GetETA()
	if eta != 0 {
		t.Errorf("GetETA with zero scanned = %v, want 0", eta)
	}

	// Complete
	pt3 := scanner.NewProgressTracker()
	pt3.Update(100, 100, "/test/file.txt", 1024)
	// Wait a tiny bit for time to pass
	time.Sleep(10 * time.Millisecond)
	eta = pt3.GetETA()
	if eta > 0 {
		t.Logf("GetETA when complete = %v (may be small positive)", eta)
	}
}

// TestProgressTracker_Reset tests Reset
func TestProgressTracker_Reset(t *testing.T) {
	pt := scanner.NewProgressTracker()

	// Update with some data
	pt.Update(50, 100, "/test/file.txt", 51200)

	// Wait a bit
	time.Sleep(10 * time.Millisecond)

	// Reset
	pt.Reset()

	progress := pt.GetProgress()
	if progress.FilesScanned != 0 {
		t.Errorf("FilesScanned after reset = %d, want 0", progress.FilesScanned)
	}
	if progress.TotalFiles != 0 {
		t.Errorf("TotalFiles after reset = %d, want 0", progress.TotalFiles)
	}
	if progress.CurrentFile != "" {
		t.Errorf("CurrentFile after reset = %s, want empty", progress.CurrentFile)
	}
	if progress.TotalSize != 0 {
		t.Errorf("TotalSize after reset = %d, want 0", progress.TotalSize)
	}
	if progress.PercentComplete != 0 {
		t.Errorf("PercentComplete after reset = %f, want 0", progress.PercentComplete)
	}
}

// TestProgressTracker_ConcurrentAccess tests concurrent access
func TestProgressTracker_ConcurrentAccess(t *testing.T) {
	pt := scanner.NewProgressTracker()

	var wg sync.WaitGroup

	// Multiple goroutines updating
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				pt.Update(int64(id*100+j), 1000, "/test/file.txt", 1024)
			}
		}(i)
	}

	// Multiple goroutines reading
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				pt.GetProgress()
				pt.GetElapsed()
				pt.GetFilesPerSecond()
				pt.GetETA()
			}
		}()
	}

	wg.Wait()
}

// TestProgressTracker_CallbackWithData tests callback receives correct data
func TestProgressTracker_CallbackWithData(t *testing.T) {
	pt := scanner.NewProgressTracker()

	var receivedProgress scanner.ScanProgress
	pt.AddCallback(func(progress scanner.ScanProgress) {
		receivedProgress = progress
	})

	expectedProgress := scanner.ScanProgress{
		FilesScanned:    25,
		TotalFiles:      100,
		CurrentFile:     "/test/specific/file.txt",
		TotalSize:       25600,
		PercentComplete: 25.0,
	}

	pt.Update(expectedProgress.FilesScanned, expectedProgress.TotalFiles,
		expectedProgress.CurrentFile, expectedProgress.TotalSize)

	if receivedProgress.FilesScanned != expectedProgress.FilesScanned {
		t.Errorf("Callback FilesScanned = %d, want %d",
			receivedProgress.FilesScanned, expectedProgress.FilesScanned)
	}
	if receivedProgress.CurrentFile != expectedProgress.CurrentFile {
		t.Errorf("Callback CurrentFile = %s, want %s",
			receivedProgress.CurrentFile, expectedProgress.CurrentFile)
	}
}

// TestProgressTracker_Update_MultipleTimes tests multiple updates
func TestProgressTracker_Update_MultipleTimes(t *testing.T) {
	pt := scanner.NewProgressTracker()

	updateCount := 0
	pt.AddCallback(func(progress scanner.ScanProgress) {
		updateCount++
	})

	for i := 1; i <= 10; i++ {
		pt.Update(int64(i*10), 100, "/test/file.txt", 1024)
	}

	if updateCount != 10 {
		t.Errorf("Expected 10 callback calls, got %d", updateCount)
	}

	progress := pt.GetProgress()
	if progress.FilesScanned != 100 {
		t.Errorf("Final FilesScanned = %d, want 100", progress.FilesScanned)
	}
	if progress.PercentComplete != 100 {
		t.Errorf("Final PercentComplete = %f, want 100", progress.PercentComplete)
	}
}

// TestScanProgress tests ScanProgress struct
func TestScanProgress(t *testing.T) {
	progress := scanner.ScanProgress{
		FilesScanned:    50,
		TotalFiles:      100,
		CurrentFile:     "/current/file.txt",
		TotalSize:       51200,
		PercentComplete: 50.0,
	}

	if progress.FilesScanned != 50 {
		t.Errorf("FilesScanned = %d, want 50", progress.FilesScanned)
	}
	if progress.TotalFiles != 100 {
		t.Errorf("TotalFiles = %d, want 100", progress.TotalFiles)
	}
	if progress.CurrentFile != "/current/file.txt" {
		t.Errorf("CurrentFile = %s, want /current/file.txt", progress.CurrentFile)
	}
	if progress.TotalSize != 51200 {
		t.Errorf("TotalSize = %d, want 51200", progress.TotalSize)
	}
	if progress.PercentComplete != 50.0 {
		t.Errorf("PercentComplete = %f, want 50.0", progress.PercentComplete)
	}
}

// TestProgressTracker_GetETA_Complete tests ETA when complete
func TestProgressTracker_GetETA_Complete(t *testing.T) {
	pt := scanner.NewProgressTracker()
	
	// Simulate complete scan
	pt.Update(100, 100, "/last/file.txt", 102400)
	
	// Wait a tiny bit for time calculation
	time.Sleep(10 * time.Millisecond)
	
	eta := pt.GetETA()
	// When complete, ETA should be 0 or very small
	if eta > time.Second {
		t.Errorf("GetETA when complete = %v, should be 0 or very small", eta)
	}
}
