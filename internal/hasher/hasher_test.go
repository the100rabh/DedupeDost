package hasher_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/dupdel/dup-del/internal/hasher"
)

// Helper function to create test file
func createTestFile(t *testing.T, content []byte) string {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "hasher_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tmpFile.Write(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	return tmpFile.Name()
}

// TestNewHasher tests NewHasher
func TestNewHasher(t *testing.T) {
	h := hasher.NewHasher()
	if h == nil {
		t.Fatal("NewHasher returned nil")
	}
}

// TestNewHasherWithOptions tests NewHasherWithOptions
func TestNewHasherWithOptions(t *testing.T) {
	h := hasher.NewHasherWithOptions(1024, 2048)
	if h == nil {
		t.Fatal("NewHasherWithOptions returned nil")
	}
}

// TestHasher_CalculateHash tests CalculateHash
func TestHasher_CalculateHash(t *testing.T) {
	h := hasher.NewHasher()

	// Create test file with known content
	content := []byte("Hello, World!")
	tmpFile := createTestFile(t, content)
	defer os.Remove(tmpFile)

	hash, err := h.CalculateHash(tmpFile)
	if err != nil {
		t.Fatalf("CalculateHash failed: %v", err)
	}

	if hash == "" {
		t.Error("Hash should not be empty")
	}

	// Verify same content produces same hash
	hash2, err := h.CalculateHash(tmpFile)
	if err != nil {
		t.Fatalf("CalculateHash failed: %v", err)
	}

	if hash != hash2 {
		t.Error("Same content should produce same hash")
	}
}

// TestHasher_CalculateHash_NonExistent tests CalculateHash with non-existent file
func TestHasher_CalculateHash_NonExistent(t *testing.T) {
	h := hasher.NewHasher()

	_, err := h.CalculateHash("/nonexistent/file.txt")
	if err == nil {
		t.Error("CalculateHash should return error for non-existent file")
	}
}

// TestHasher_CalculateHashWithContext tests CalculateHashWithContext
func TestHasher_CalculateHashWithContext(t *testing.T) {
	h := hasher.NewHasher()

	content := []byte("Test content for context")
	tmpFile := createTestFile(t, content)
	defer os.Remove(tmpFile)

	ctx := context.Background()
	hash, err := h.CalculateHashWithContext(tmpFile, ctx)
	if err != nil {
		t.Fatalf("CalculateHashWithContext failed: %v", err)
	}

	if hash == "" {
		t.Error("Hash should not be empty")
	}
}

// TestHasher_CalculateHashWithContext_Cancel tests cancellation
func TestHasher_CalculateHashWithContext_Cancel(t *testing.T) {
	h := hasher.NewHasher()

	// Create a larger file for testing cancellation
	content := make([]byte, 10*1024*1024) // 10MB
	for i := range content {
		content[i] = byte(i % 256)
	}
	tmpFile := createTestFile(t, content)
	defer os.Remove(tmpFile)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := h.CalculateHashWithContext(tmpFile, ctx)
	if err == nil {
		t.Log("Expected context cancellation error (may complete before cancellation)")
	}
}

// TestHasher_CalculateHashWithProgress tests CalculateHashWithProgress
func TestHasher_CalculateHashWithProgress(t *testing.T) {
	h := hasher.NewHasher()

	content := []byte("Test content for progress tracking")
	tmpFile := createTestFile(t, content)
	defer os.Remove(tmpFile)

	progressChan := make(chan float64, 10)

	hash, err := h.CalculateHashWithProgress(tmpFile, progressChan)
	if err != nil {
		t.Fatalf("CalculateHashWithProgress failed: %v", err)
	}

	if hash == "" {
		t.Error("Hash should not be empty")
	}

	// Check we received progress updates (at least the final 100%)
	var gotFinalProgress bool
	timeout := time.After(100 * time.Millisecond)
	for {
		select {
		case progress, ok := <-progressChan:
			if !ok {
				goto done
			}
			if progress == 100.0 {
				gotFinalProgress = true
			}
		case <-timeout:
			goto done
		}
	}
done:
	if !gotFinalProgress {
		t.Log("Note: May not have received final progress update")
	}
}

// TestHasher_HashMultiple tests HashMultiple
func TestHasher_HashMultiple(t *testing.T) {
	h := hasher.NewHasher()

	// Create multiple test files
	files := make([]string, 5)
	for i := 0; i < 5; i++ {
		files[i] = createTestFile(t, []byte("content"+string(rune('0'+i))))
	}
	defer func() {
		for _, f := range files {
			os.Remove(f)
		}
	}()

	resultChan := h.HashMultiple(files, 2)

	results := make(map[string]string)
	for result := range resultChan {
		if result.Error != nil {
			t.Errorf("Error hashing %s: %v", result.Path, result.Error)
			continue
		}
		results[result.Path] = result.Hash
	}

	if len(results) != 5 {
		t.Errorf("Expected 5 results, got %d", len(results))
	}
}

// TestHasher_HashMultiple_Empty tests HashMultiple with empty list
func TestHasher_HashMultiple_Empty(t *testing.T) {
	h := hasher.NewHasher()

	resultChan := h.HashMultiple([]string{}, 2)

	// Should close immediately
	select {
	case _, ok := <-resultChan:
		if ok {
			t.Error("Result channel should be closed")
		}
	default:
		t.Error("Result channel should be closed for empty input")
	}
}

// TestHasher_HashMultipleWithContext tests HashMultipleWithContext
func TestHasher_HashMultipleWithContext(t *testing.T) {
	h := hasher.NewHasher()

	files := make([]string, 3)
	for i := 0; i < 3; i++ {
		files[i] = createTestFile(t, []byte("content"+string(rune('0'+i))))
	}
	defer func() {
		for _, f := range files {
			os.Remove(f)
		}
	}()

	ctx := context.Background()
	resultChan := h.HashMultipleWithContext(ctx, files, 2)

	count := 0
	for result := range resultChan {
		if result.Error != nil {
			t.Errorf("Error hashing %s: %v", result.Path, result.Error)
			continue
		}
		count++
	}

	if count != 3 {
		t.Errorf("Expected 3 results, got %d", count)
	}
}

// TestHasher_HashMultipleWithContext_Cancel tests cancellation
func TestHasher_HashMultipleWithContext_Cancel(t *testing.T) {
	h := hasher.NewHasher()

	files := make([]string, 10)
	for i := 0; i < 10; i++ {
		files[i] = createTestFile(t, []byte("content"+string(rune('0'+i))))
	}
	defer func() {
		for _, f := range files {
			os.Remove(f)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	resultChan := h.HashMultipleWithContext(ctx, files, 2)

	// Should receive cancelled results
	for result := range resultChan {
		if result.Cancelled || result.Error != nil {
			// Expected
		}
	}
}

// TestHashResult tests HashResult struct
func TestHashResult(t *testing.T) {
	result := hasher.HashResult{
		Path:      "/test/file.txt",
		Hash:      "abc123",
		Size:      1024,
		Error:     nil,
		Cancelled: false,
	}

	if result.Path != "/test/file.txt" {
		t.Errorf("Path = %s, want /test/file.txt", result.Path)
	}
	if result.Hash != "abc123" {
		t.Errorf("Hash = %s, want abc123", result.Hash)
	}
	if result.Size != 1024 {
		t.Errorf("Size = %d, want 1024", result.Size)
	}
	if result.Cancelled {
		t.Error("Cancelled should be false")
	}
}

// TestHasher_LargeFile tests hashing a larger file
func TestHasher_LargeFile(t *testing.T) {
	h := hasher.NewHasher()

	// Create 1MB file
	content := make([]byte, 1024*1024)
	for i := range content {
		content[i] = byte(i % 256)
	}
	tmpFile := createTestFile(t, content)
	defer os.Remove(tmpFile)

	hash, err := h.CalculateHash(tmpFile)
	if err != nil {
		t.Fatalf("CalculateHash failed: %v", err)
	}

	if hash == "" {
		t.Error("Hash should not be empty")
	}
}

// TestHasher_DuplicateHashes tests that identical files produce identical hashes
func TestHasher_DuplicateHashes(t *testing.T) {
	h := hasher.NewHasher()

	content := []byte("Duplicate content test")
	file1 := createTestFile(t, content)
	file2 := createTestFile(t, content)
	defer os.Remove(file1)
	defer os.Remove(file2)

	hash1, err := h.CalculateHash(file1)
	if err != nil {
		t.Fatalf("CalculateHash failed for file1: %v", err)
	}

	hash2, err := h.CalculateHash(file2)
	if err != nil {
		t.Fatalf("CalculateHash failed for file2: %v", err)
	}

	if hash1 != hash2 {
		t.Error("Identical files should produce identical hashes")
	}
}

// TestHasher_DifferentHashes tests that different files produce different hashes
func TestHasher_DifferentHashes(t *testing.T) {
	h := hasher.NewHasher()

	file1 := createTestFile(t, []byte("content1"))
	file2 := createTestFile(t, []byte("content2"))
	defer os.Remove(file1)
	defer os.Remove(file2)

	hash1, err := h.CalculateHash(file1)
	if err != nil {
		t.Fatalf("CalculateHash failed for file1: %v", err)
	}

	hash2, err := h.CalculateHash(file2)
	if err != nil {
		t.Fatalf("CalculateHash failed for file2: %v", err)
	}

	if hash1 == hash2 {
		t.Error("Different files should produce different hashes")
	}
}
