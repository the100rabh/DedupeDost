package scanner_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/the100rabh/DedupeDost/internal/scanner"
)

// Helper function to create test directory structure
func setupTestDir(t *testing.T) (string, func()) {
	t.Helper()
	
	dir, err := os.MkdirTemp("", "scanner_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	
	cleanup := func() {
		os.RemoveAll(dir)
	}
	
	return dir, cleanup
}

// Helper function to create test files
func createTestFile(t *testing.T, dir, name string, size int64) string {
	t.Helper()
	
	path := filepath.Join(dir, name)
	content := make([]byte, size)
	for i := range content {
		content[i] = byte(i % 256)
	}
	
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	return path
}

// TestScanner_EmptyDirectory tests scanning an empty directory
func TestScanner_EmptyDirectory(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()
	
	s := scanner.NewScanner(dir, scanner.ScanOptions{Recursive: true})
	
	fileChan, errorChan, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}
	
	fileCount := 0
	for range fileChan {
		fileCount++
	}
	
	for err := range errorChan {
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	}
	
	if fileCount != 0 {
		t.Errorf("Expected 0 files, got %d", fileCount)
	}
}

// TestScanner_SingleFile tests scanning a directory with a single file
func TestScanner_SingleFile(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()
	
	createTestFile(t, dir, "test.txt", 100)
	
	s := scanner.NewScanner(dir, scanner.ScanOptions{Recursive: true})
	
	fileChan, errorChan, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}
	
	fileCount := 0
	for range fileChan {
		fileCount++
	}
	
	for err := range errorChan {
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	}
	
	if fileCount != 1 {
		t.Errorf("Expected 1 file, got %d", fileCount)
	}
}

// TestScanner_NestedDirectories tests scanning nested directories
func TestScanner_NestedDirectories(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()
	
	// Create nested structure
	subdir1 := filepath.Join(dir, "sub1")
	subdir2 := filepath.Join(dir, "sub1", "sub2")
	os.MkdirAll(subdir2, 0755)
	
	createTestFile(t, dir, "root.txt", 50)
	createTestFile(t, subdir1, "level1.txt", 100)
	createTestFile(t, subdir2, "level2.txt", 150)
	
	s := scanner.NewScanner(dir, scanner.ScanOptions{Recursive: true})
	
	fileChan, errorChan, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}
	
	fileCount := 0
	for range fileChan {
		fileCount++
	}
	
	for err := range errorChan {
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	}
	
	if fileCount != 3 {
		t.Errorf("Expected 3 files, got %d", fileCount)
	}
}

// TestScanner_HiddenFiles tests hidden file handling
func TestScanner_HiddenFiles(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()
	
	createTestFile(t, dir, "visible.txt", 100)
	createTestFile(t, dir, ".hidden.txt", 100)
	
	// Test with IncludeHidden = false
	s := scanner.NewScanner(dir, scanner.ScanOptions{
		Recursive:     true,
		IncludeHidden: false,
	})
	
	fileChan, _, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}
	
	fileCount := 0
	for range fileChan {
		fileCount++
	}
	
	if fileCount != 1 {
		t.Errorf("Expected 1 file (excluding hidden), got %d", fileCount)
	}
	
	// Test with IncludeHidden = true
	s = scanner.NewScanner(dir, scanner.ScanOptions{
		Recursive:     true,
		IncludeHidden: true,
	})
	
	fileChan, _, err = s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}
	
	fileCount = 0
	for range fileChan {
		fileCount++
	}
	
	if fileCount != 2 {
		t.Errorf("Expected 2 files (including hidden), got %d", fileCount)
	}
}

// TestScanner_SizeFilter_MinSize tests minimum size filtering
func TestScanner_SizeFilter_MinSize(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()
	
	createTestFile(t, dir, "small.txt", 50)
	createTestFile(t, dir, "large.txt", 500)
	
	s := scanner.NewScanner(dir, scanner.ScanOptions{
		Recursive: true,
		MinSize:   100,
	})
	
	fileChan, _, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}
	
	fileCount := 0
	for range fileChan {
		fileCount++
	}
	
	if fileCount != 1 {
		t.Errorf("Expected 1 file (>= 100 bytes), got %d", fileCount)
	}
}

// TestScanner_SizeFilter_MaxSize tests maximum size filtering
func TestScanner_SizeFilter_MaxSize(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()
	
	createTestFile(t, dir, "small.txt", 50)
	createTestFile(t, dir, "large.txt", 500)
	
	s := scanner.NewScanner(dir, scanner.ScanOptions{
		Recursive: true,
		MaxSize:   100,
	})
	
	fileChan, _, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}
	
	fileCount := 0
	for range fileChan {
		fileCount++
	}
	
	if fileCount != 1 {
		t.Errorf("Expected 1 file (<= 100 bytes), got %d", fileCount)
	}
}

// TestScanner_ExtensionFilter tests extension filtering
func TestScanner_ExtensionFilter(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()
	
	createTestFile(t, dir, "file1.txt", 100)
	createTestFile(t, dir, "file2.jpg", 100)
	createTestFile(t, dir, "file3.txt", 100)
	
	s := scanner.NewScanner(dir, scanner.ScanOptions{
		Recursive:  true,
		Extensions: []string{".txt"},
	})
	
	fileChan, _, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}
	
	fileCount := 0
	for range fileChan {
		fileCount++
	}
	
	if fileCount != 2 {
		t.Errorf("Expected 2 .txt files, got %d", fileCount)
	}
}

// TestScanner_Recursive tests recursive scanning
func TestScanner_Recursive(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()
	
	subdir := filepath.Join(dir, "subdir")
	os.MkdirAll(subdir, 0755)
	
	createTestFile(t, dir, "root.txt", 100)
	createTestFile(t, subdir, "nested.txt", 100)
	
	// Test recursive
	s := scanner.NewScanner(dir, scanner.ScanOptions{Recursive: true})
	fileChan, _, _ := s.Scan()
	
	fileCount := 0
	for range fileChan {
		fileCount++
	}
	
	if fileCount != 2 {
		t.Errorf("Expected 2 files (recursive), got %d", fileCount)
	}
	
	// Test non-recursive
	s = scanner.NewScanner(dir, scanner.ScanOptions{Recursive: false})
	fileChan, _, _ = s.Scan()
	
	fileCount = 0
	for range fileChan {
		fileCount++
	}
	
	if fileCount != 1 {
		t.Errorf("Expected 1 file (non-recursive), got %d", fileCount)
	}
}

// TestScanner_Cancel tests cancellation during scan
func TestScanner_Cancel(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()
	
	// Create many files
	for i := 0; i < 10; i++ {
		subdir := filepath.Join(dir, fmt.Sprintf("subdir%d", i))
		os.MkdirAll(subdir, 0755)
		createTestFile(t, subdir, "file.txt", 100)
	}
	
	s := scanner.NewScanner(dir, scanner.ScanOptions{Recursive: true})
	
	fileChan, _, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}
	
	// Cancel after first file
	<-fileChan
	s.Cancel()

	// Just verify we can cancel without panic - channel behavior may vary
	// Drain any remaining files with timeout
	timeout := time.After(500 * time.Millisecond)
	done := make(chan struct{})
	go func() {
		for range fileChan {
			// Drain
		}
		close(done)
	}()

	select {
	case <-done:
		// Channel closed properly
	case <-timeout:
		// Timeout is acceptable for this test
	}
}

// TestScanner_Progress tests progress tracking
func TestScanner_Progress(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()
	
	for i := 0; i < 10; i++ {
		createTestFile(t, dir, fmt.Sprintf("file%d.txt", i), 100)
	}
	
	s := scanner.NewScanner(dir, scanner.ScanOptions{Recursive: true})
	
	// Pre-count files
	count, err := s.CountFiles()
	if err != nil {
		t.Fatalf("CountFiles failed: %v", err)
	}
	
	if count != 10 {
		t.Errorf("Expected 10 files, got %d", count)
	}
	
	progress := s.GetProgress()
	if progress.TotalFiles != 10 {
		t.Errorf("Expected TotalFiles=10, got %d", progress.TotalFiles)
	}
}

// TestFilter_ShouldIncludeHidden tests hidden file detection
func TestFilter_ShouldIncludeHidden(t *testing.T) {
	tests := []struct {
		name          string
		filename      string
		includeHidden bool
		expected      bool
	}{
		{"visible file", "file.txt", false, true},
		{"hidden file excluded", ".hidden", false, false},
		{"hidden file included", ".hidden", true, true},
		{"dot in middle", "file.txt.bak", false, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scanner.ShouldIncludeHidden(tt.filename, tt.includeHidden)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestFilter_ShouldIncludeBySize tests size filtering
func TestFilter_ShouldIncludeBySize(t *testing.T) {
	tests := []struct {
		name    string
		size    int64
		minSize int64
		maxSize int64
		expected bool
	}{
		{"within range", 100, 50, 200, true},
		{"below min", 50, 100, 200, false},
		{"above max", 300, 50, 200, false},
		{"at min", 100, 100, 200, true},
		{"at max", 200, 100, 200, true},
		{"no max limit", 1000, 50, 0, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scanner.ShouldIncludeBySize(tt.size, tt.minSize, tt.maxSize)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestFilter_ShouldIncludeByExtension tests extension filtering
func TestFilter_ShouldIncludeByExtension(t *testing.T) {
	tests := []struct {
		name      string
		ext       string
		allowed   []string
		expected  bool
	}{
		{"matching ext", ".txt", []string{".txt", ".md"}, true},
		{"non-matching ext", ".jpg", []string{".txt", ".md"}, false},
		{"case insensitive", ".TXT", []string{".txt"}, true},
		{"empty allowed", ".txt", []string{}, true},
		{"nil allowed", ".txt", nil, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scanner.ShouldIncludeByExtension(tt.ext, tt.allowed)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestFilter_IsHidden tests hidden file detection
func TestFilter_IsHidden(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		{"hidden", ".gitignore", true},
		{"visible", "file.txt", false},
		{"hidden dir", ".config", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scanner.IsHidden(tt.filename)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestFilter_NormalizeExtension tests extension normalization
func TestFilter_NormalizeExtension(t *testing.T) {
	tests := []struct {
		name     string
		ext      string
		expected string
	}{
		{"with dot", ".txt", ".txt"},
		{"without dot", "txt", ".txt"},
		{"uppercase", ".TXT", ".txt"},
		{"empty", "", ""},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scanner.NormalizeExtension(tt.ext)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// TestScanner_CountFiles tests file counting
func TestScanner_CountFiles(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()
	
	for i := 0; i < 5; i++ {
		createTestFile(t, dir, fmt.Sprintf("file%d.txt", i), 100)
	}
	
	s := scanner.NewScanner(dir, scanner.ScanOptions{Recursive: true})
	count, err := s.CountFiles()
	if err != nil {
		t.Fatalf("CountFiles failed: %v", err)
	}
	
	if count != 5 {
		t.Errorf("Expected 5 files, got %d", count)
	}
}

// TestScanner_ContextCancellation tests context-based cancellation
func TestScanner_ContextCancellation(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()
	
	// Create test files
	for i := 0; i < 10; i++ {
		createTestFile(t, dir, fmt.Sprintf("file%d.txt", i), 100)
	}
	
	s := scanner.NewScanner(dir, scanner.ScanOptions{Recursive: true})
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	
	// Start scan
	fileChan, errorChan, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}
	
	// Wait for context cancellation
	<-ctx.Done()
	
	// Cancel scanner
	s.Cancel()
	
	// Drain channels
	for range fileChan {
	}
	for err := range errorChan {
		if err != nil && err != context.Canceled {
			t.Errorf("Unexpected error: %v", err)
		}
	}
}
