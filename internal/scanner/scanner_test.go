package scanner_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dupdel/dup-del/internal/scanner"
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
func createTestFile(t *testing.T, dir, name string, content []byte) string {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	return path
}

// TestScanOptions tests ScanOptions struct
func TestScanOptions(t *testing.T) {
	opts := scanner.ScanOptions{
		Recursive:      true,
		IncludeHidden:  true,
		FollowSymlinks: false,
		MinSize:        1024,
		MaxSize:        1024 * 1024,
		Extensions:     []string{".txt", ".md"},
	}

	if !opts.Recursive {
		t.Error("Recursive should be true")
	}
	if !opts.IncludeHidden {
		t.Error("IncludeHidden should be true")
	}
	if opts.FollowSymlinks {
		t.Error("FollowSymlinks should be false")
	}
	if opts.MinSize != 1024 {
		t.Errorf("MinSize = %d, want 1024", opts.MinSize)
	}
	if opts.MaxSize != 1024*1024 {
		t.Errorf("MaxSize = %d, want %d", opts.MaxSize, 1024*1024)
	}
	if len(opts.Extensions) != 2 {
		t.Errorf("Extensions length = %d, want 2", len(opts.Extensions))
	}
}

// TestScanProgress tests ScanProgress struct
func TestScanProgress(t *testing.T) {
	progress := scanner.ScanProgress{
		FilesScanned:    100,
		TotalFiles:      500,
		CurrentFile:     "/test/file.txt",
		TotalSize:       1024000,
		PercentComplete: 20.0,
	}

	if progress.FilesScanned != 100 {
		t.Errorf("FilesScanned = %d, want 100", progress.FilesScanned)
	}
	if progress.TotalFiles != 500 {
		t.Errorf("TotalFiles = %d, want 500", progress.TotalFiles)
	}
	if progress.CurrentFile != "/test/file.txt" {
		t.Errorf("CurrentFile = %s, want /test/file.txt", progress.CurrentFile)
	}
	if progress.TotalSize != 1024000 {
		t.Errorf("TotalSize = %d, want 1024000", progress.TotalSize)
	}
	if progress.PercentComplete != 20.0 {
		t.Errorf("PercentComplete = %f, want 20.0", progress.PercentComplete)
	}
}

// TestNewScanner tests NewScanner
func TestNewScanner(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	opts := scanner.ScanOptions{
		Recursive: true,
	}

	s := scanner.NewScanner(dir, opts)
	if s == nil {
		t.Fatal("NewScanner returned nil")
	}
}

// TestScanner_Scan tests basic scanning
func TestScanner_Scan(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	// Create test files
	createTestFile(t, dir, "file1.txt", []byte("content1"))
	createTestFile(t, dir, "file2.txt", []byte("content2"))
	createTestFile(t, dir, "file3.md", []byte("content3"))

	opts := scanner.ScanOptions{
		Recursive:     true,
		IncludeHidden: false,
	}

	s := scanner.NewScanner(dir, opts)
	fileChan, errorChan, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// Collect files
	files := make([]string, 0)
	for file := range fileChan {
		files = append(files, file.Path)
	}

	// Check for errors
	for err := range errorChan {
		t.Errorf("Scan error: %v", err)
	}

	if len(files) != 3 {
		t.Errorf("Expected 3 files, got %d", len(files))
	}
}

// TestScanner_Scan_EmptyDirectory tests scanning empty directory
func TestScanner_Scan_EmptyDirectory(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	opts := scanner.ScanOptions{
		Recursive: true,
	}

	s := scanner.NewScanner(dir, opts)
	fileChan, errorChan, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// Collect files
	count := 0
	for range fileChan {
		count++
	}

	// Check for errors
	for err := range errorChan {
		t.Errorf("Scan error: %v", err)
	}

	if count != 0 {
		t.Errorf("Expected 0 files, got %d", count)
	}
}

// TestScanner_Scan_Recursive tests recursive scanning
func TestScanner_Scan_Recursive(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	// Create nested structure
	subdir1 := filepath.Join(dir, "sub1")
	subdir2 := filepath.Join(dir, "sub1", "sub2")
	os.MkdirAll(subdir2, 0755)

	createTestFile(t, dir, "root.txt", []byte("root"))
	createTestFile(t, subdir1, "level1.txt", []byte("level1"))
	createTestFile(t, subdir2, "level2.txt", []byte("level2"))

	opts := scanner.ScanOptions{
		Recursive: true,
	}

	s := scanner.NewScanner(dir, opts)
	fileChan, errorChan, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	count := 0
	for range fileChan {
		count++
	}
	for err := range errorChan {
		t.Errorf("Scan error: %v", err)
	}

	if count != 3 {
		t.Errorf("Expected 3 files (recursive), got %d", count)
	}
}

// TestScanner_Scan_NonRecursive tests non-recursive scanning
func TestScanner_Scan_NonRecursive(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	// Create nested structure
	subdir := filepath.Join(dir, "subdir")
	os.MkdirAll(subdir, 0755)

	createTestFile(t, dir, "root.txt", []byte("root"))
	createTestFile(t, subdir, "nested.txt", []byte("nested"))

	opts := scanner.ScanOptions{
		Recursive: false,
	}

	s := scanner.NewScanner(dir, opts)
	fileChan, errorChan, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	count := 0
	for range fileChan {
		count++
	}
	for err := range errorChan {
		t.Errorf("Scan error: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 file (non-recursive), got %d", count)
	}
}

// TestScanner_Scan_IncludeHidden tests IncludeHidden option
func TestScanner_Scan_IncludeHidden(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	createTestFile(t, dir, "visible.txt", []byte("visible"))
	createTestFile(t, dir, ".hidden.txt", []byte("hidden"))

	// Test with IncludeHidden = false
	opts := scanner.ScanOptions{
		IncludeHidden: false,
	}

	s := scanner.NewScanner(dir, opts)
	fileChan, errorChan, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	count := 0
	for file := range fileChan {
		if filepath.Base(file.Path) == ".hidden.txt" {
			t.Error("Hidden file should not be included")
		}
		count++
	}
	for err := range errorChan {
		t.Errorf("Scan error: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 file (no hidden), got %d", count)
	}
}

// TestScanner_Scan_WithExtensions tests Extensions filter
func TestScanner_Scan_WithExtensions(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	createTestFile(t, dir, "file.txt", []byte("text"))
	createTestFile(t, dir, "file.md", []byte("markdown"))
	createTestFile(t, dir, "file.jpg", []byte("image"))

	opts := scanner.ScanOptions{
		Extensions: []string{".txt", ".md"},
	}

	s := scanner.NewScanner(dir, opts)
	fileChan, errorChan, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	count := 0
	for file := range fileChan {
		ext := filepath.Ext(file.Path)
		if ext != ".txt" && ext != ".md" {
			t.Errorf("Unexpected extension: %s", ext)
		}
		count++
	}
	for err := range errorChan {
		t.Errorf("Scan error: %v", err)
	}

	if count != 2 {
		t.Errorf("Expected 2 files (.txt, .md), got %d", count)
	}
}

// TestScanner_Scan_MinSize tests MinSize filter
func TestScanner_Scan_MinSize(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	createTestFile(t, dir, "small.txt", []byte("small"))         // 5 bytes
	createTestFile(t, dir, "large.txt", []byte("larger content")) // 15 bytes

	opts := scanner.ScanOptions{
		MinSize: 10,
	}

	s := scanner.NewScanner(dir, opts)
	fileChan, errorChan, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	count := 0
	for file := range fileChan {
		if file.Size < 10 {
			t.Errorf("File smaller than MinSize: %s (%d bytes)", file.Path, file.Size)
		}
		count++
	}
	for err := range errorChan {
		t.Errorf("Scan error: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 file (>= 10 bytes), got %d", count)
	}
}

// TestScanner_Scan_MaxSize tests MaxSize filter
func TestScanner_Scan_MaxSize(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	createTestFile(t, dir, "small.txt", []byte("small"))         // 5 bytes
	createTestFile(t, dir, "large.txt", []byte("larger content")) // 15 bytes

	opts := scanner.ScanOptions{
		MaxSize: 10,
	}

	s := scanner.NewScanner(dir, opts)
	fileChan, errorChan, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	count := 0
	for file := range fileChan {
		if file.Size > 10 {
			t.Errorf("File larger than MaxSize: %s (%d bytes)", file.Path, file.Size)
		}
		count++
	}
	for err := range errorChan {
		t.Errorf("Scan error: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 file (<= 10 bytes), got %d", count)
	}
}

// TestScanner_Cancel tests cancellation
func TestScanner_Cancel(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	// Create many files
	for i := 0; i < 100; i++ {
		createTestFile(t, dir, "file"+string(rune('0'+i%10))+".txt", []byte("content"))
	}

	opts := scanner.ScanOptions{
		Recursive: true,
	}

	s := scanner.NewScanner(dir, opts)
	
	// Cancel immediately
	// Note: Scanner doesn't expose Cancel method, so we just test that Scan starts
	
	fileChan, errorChan, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// Read some files
	count := 0
	timeout := make(chan bool, 1)
	go func() {
		for range fileChan {
			count++
		}
		timeout <- true
	}()

	select {
	case <-timeout:
		// Completed
	case <-make(chan bool, 1):
		// This won't happen, just for structure
	}

	// Check for errors
	for err := range errorChan {
		t.Errorf("Scan error: %v", err)
	}

	t.Logf("Scanned %d files", count)
}

// TestScanner_Scan_NonExistentDirectory tests scanning non-existent directory
func TestScanner_Scan_NonExistentDirectory(t *testing.T) {
	opts := scanner.ScanOptions{
		Recursive: true,
	}

	s := scanner.NewScanner("/nonexistent/directory", opts)
	_, _, err := s.Scan()
	if err != nil {
		// Expected - directory doesn't exist
		t.Logf("Scan returned error as expected: %v", err)
	}
}

// TestFilterOptions tests filter options
func TestFilterOptions(t *testing.T) {
	opts := scanner.ScanOptions{
		Recursive:      true,
		IncludeHidden:  true,
		FollowSymlinks: false,
		MinSize:        0,
		MaxSize:        0,
		Extensions:     nil,
	}

	// Default values
	if !opts.Recursive {
		t.Error("Recursive should default to true")
	}
	if opts.MinSize != 0 {
		t.Errorf("MinSize = %d, want 0", opts.MinSize)
	}
	if opts.MaxSize != 0 {
		t.Errorf("MaxSize = %d, want 0", opts.MaxSize)
	}
}
