package hasher_test

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/dupdel/dup-del/internal/hasher"
)

// TestNewHashCache tests NewHashCache
func TestNewHashCache(t *testing.T) {
	cache := hasher.NewHashCache(100)
	if cache == nil {
		t.Fatal("NewHashCache returned nil")
	}

	stats := cache.Stats()
	if stats.Size != 0 {
		t.Errorf("Initial size = %d, want 0", stats.Size)
	}
	if stats.MaxSize != 100 {
		t.Errorf("MaxSize = %d, want 100", stats.MaxSize)
	}
}

// TestNewHashCache_ZeroMaxSize tests NewHashCache with zero max size
func TestNewHashCache_ZeroMaxSize(t *testing.T) {
	cache := hasher.NewHashCache(0)
	if cache == nil {
		t.Fatal("NewHashCache returned nil")
	}

	stats := cache.Stats()
	if stats.MaxSize != 10000 {
		t.Errorf("MaxSize with 0 = %d, want 10000", stats.MaxSize)
	}
}

// TestHashCache_GetOrCompute tests GetOrCompute
func TestHashCache_GetOrCompute(t *testing.T) {
	cache := hasher.NewHashCache(100)
	hasher := hasher.NewHasher()

	content := []byte("Test content for caching")
	tmpFile := createTestFile(t, content)
	defer os.Remove(tmpFile)

	// First call - should compute
	hash1, err := cache.GetOrCompute(tmpFile, hasher)
	if err != nil {
		t.Fatalf("GetOrCompute failed: %v", err)
	}

	if hash1 == "" {
		t.Error("Hash should not be empty")
	}

	// Second call - should return cached
	hash2, err := cache.GetOrCompute(tmpFile, hasher)
	if err != nil {
		t.Fatalf("GetOrCompute failed: %v", err)
	}

	if hash1 != hash2 {
		t.Error("Cached hash should match computed hash")
	}
}

// TestHashCache_Get tests Get
func TestHashCache_Get(t *testing.T) {
	cache := hasher.NewHashCache(100)

	// Non-existent key
	_, ok := cache.Get("/nonexistent/file.txt")
	if ok {
		t.Error("Get should return false for non-existent key")
	}
}

// TestHashCache_Get_ModifiedFile tests Get with modified file
func TestHashCache_Get_ModifiedFile(t *testing.T) {
	cache := hasher.NewHashCache(100)
	hasher := hasher.NewHasher()

	content := []byte("Original content")
	tmpFile := createTestFile(t, content)
	defer os.Remove(tmpFile)

	// Cache the hash
	_, err := cache.GetOrCompute(tmpFile, hasher)
	if err != nil {
		t.Fatalf("GetOrCompute failed: %v", err)
	}

	// Modify the file
	if err := os.WriteFile(tmpFile, []byte("Modified content"), 0644); err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}

	// Get should return false (cache miss due to modification)
	_, ok := cache.Get(tmpFile)
	if ok {
		t.Error("Get should return false for modified file")
	}
}

// TestHashCache_Set tests Set
func TestHashCache_Set(t *testing.T) {
	cache := hasher.NewHashCache(100)

	// Create a real temp file for testing
	tmpFile := createTestFile(t, []byte("test content"))
	defer os.Remove(tmpFile)

	info, _ := os.Stat(tmpFile)
	cache.Set(tmpFile, "abc123", info.Size(), info.ModTime())

	hash, ok := cache.Get(tmpFile)
	if !ok {
		t.Error("Get should return true after Set")
	}
	if hash != "abc123" {
		t.Errorf("Hash = %s, want abc123", hash)
	}
}

// TestHashCache_Remove tests Remove
func TestHashCache_Remove(t *testing.T) {
	cache := hasher.NewHashCache(100)

	// Create a real temp file for testing
	tmpFile := createTestFile(t, []byte("test content"))
	defer os.Remove(tmpFile)

	info, _ := os.Stat(tmpFile)
	cache.Set(tmpFile, "abc123", info.Size(), info.ModTime())

	// Verify it's there
	_, ok := cache.Get(tmpFile)
	if !ok {
		t.Fatal("Get should return true after Set")
	}

	// Remove it
	cache.Remove(tmpFile)

	// Verify it's gone
	_, ok = cache.Get(tmpFile)
	if ok {
		t.Error("Get should return false after Remove")
	}
}

// TestHashCache_Clear tests Clear
func TestHashCache_Clear(t *testing.T) {
	cache := hasher.NewHashCache(100)

	// Add some entries
	cache.Set("/test/file1.txt", "abc123", 1024, time.Now())
	cache.Set("/test/file2.txt", "def456", 2048, time.Now())
	cache.Set("/test/file3.txt", "ghi789", 4096, time.Now())

	if cache.Size() != 3 {
		t.Errorf("Size = %d, want 3", cache.Size())
	}

	// Clear
	cache.Clear()

	if cache.Size() != 0 {
		t.Errorf("Size after Clear = %d, want 0", cache.Size())
	}
}

// TestHashCache_Size tests Size
func TestHashCache_Size(t *testing.T) {
	cache := hasher.NewHashCache(100)

	if cache.Size() != 0 {
		t.Errorf("Initial size = %d, want 0", cache.Size())
	}

	cache.Set("/test/file1.txt", "abc123", 1024, time.Now())
	if cache.Size() != 1 {
		t.Errorf("Size after 1 Set = %d, want 1", cache.Size())
	}

	cache.Set("/test/file2.txt", "def456", 1024, time.Now())
	if cache.Size() != 2 {
		t.Errorf("Size after 2 Sets = %d, want 2", cache.Size())
	}
}

// TestHashCache_Eviction tests cache eviction
func TestHashCache_Eviction(t *testing.T) {
	cache := hasher.NewHashCache(3)

	// Create temp files for testing
	file1 := createTestFile(t, []byte("content1"))
	file2 := createTestFile(t, []byte("content2"))
	file3 := createTestFile(t, []byte("content3"))
	file4 := createTestFile(t, []byte("content4"))
	defer func() {
		os.Remove(file1)
		os.Remove(file2)
		os.Remove(file3)
		os.Remove(file4)
	}()

	info1, _ := os.Stat(file1)
	info2, _ := os.Stat(file2)
	info3, _ := os.Stat(file3)
	info4, _ := os.Stat(file4)

	// Add 3 entries
	cache.Set(file1, "abc123", info1.Size(), info1.ModTime())
	time.Sleep(10 * time.Millisecond)
	cache.Set(file2, "def456", info2.Size(), info2.ModTime())
	time.Sleep(10 * time.Millisecond)
	cache.Set(file3, "ghi789", info3.Size(), info3.ModTime())

	if cache.Size() != 3 {
		t.Errorf("Size = %d, want 3", cache.Size())
	}

	// Add 4th entry - should evict oldest (file1)
	cache.Set(file4, "jkl012", info4.Size(), info4.ModTime())

	if cache.Size() != 3 {
		t.Errorf("Size after eviction = %d, want 3", cache.Size())
	}

	// file1 should be evicted
	_, ok := cache.Get(file1)
	if ok {
		t.Error("file1 should have been evicted")
	}

	// Others should still be there
	_, ok = cache.Get(file2)
	if !ok {
		t.Error("file2 should still be in cache")
	}
}

// TestHashCache_Save tests Save
func TestHashCache_Save(t *testing.T) {
	cache := hasher.NewHashCache(100)
	cache.Set("/test/file.txt", "abc123", 1024, time.Now())

	// Manually save (normally done via auto-save)
	// For now, just test that the cache works
	if cache.Size() != 1 {
		t.Errorf("Size = %d, want 1", cache.Size())
	}
}

// TestHashCache_Load tests Load
func TestHashCache_Load(t *testing.T) {
	tmpDir := t.TempDir()
	cacheFile := filepath.Join(tmpDir, "cache.json")

	// Create a cache file manually
	cacheData := `[
		{
			"path": "/test/file.txt",
			"hash": "abc123",
			"size": 1024,
			"mod_time": "2024-01-01T00:00:00Z",
			"cached_at": "2024-01-01T00:00:00Z"
		}
	]`

	if err := os.WriteFile(cacheFile, []byte(cacheData), 0644); err != nil {
		t.Fatalf("Failed to create cache file: %v", err)
	}

	cache, err := hasher.NewHashCacheWithPersistence(100, cacheFile)
	if err != nil {
		t.Fatalf("NewHashCacheWithPersistence failed: %v", err)
	}

	if cache.Size() != 1 {
		t.Errorf("Loaded cache size = %d, want 1", cache.Size())
	}
}

// TestHashCache_Load_NonExistent tests Load with non-existent file
func TestHashCache_Load_NonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	cacheFile := filepath.Join(tmpDir, "nonexistent.json")

	cache, err := hasher.NewHashCacheWithPersistence(100, cacheFile)
	if err != nil {
		t.Fatalf("NewHashCacheWithPersistence failed: %v", err)
	}

	if cache.Size() != 0 {
		t.Errorf("Cache size = %d, want 0", cache.Size())
	}
}

// TestHashCache_Load_InvalidJSON tests Load with invalid JSON
func TestHashCache_Load_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	cacheFile := filepath.Join(tmpDir, "cache.json")

	if err := os.WriteFile(cacheFile, []byte("invalid json"), 0644); err != nil {
		t.Fatalf("Failed to create cache file: %v", err)
	}

	_, err := hasher.NewHashCacheWithPersistence(100, cacheFile)
	if err == nil {
		t.Error("NewHashCacheWithPersistence should return error for invalid JSON")
	}
}

// TestNewHashCacheWithPersistence tests NewHashCacheWithPersistence
func TestNewHashCacheWithPersistence(t *testing.T) {
	tmpDir := t.TempDir()
	cacheFile := filepath.Join(tmpDir, "cache.json")

	cache, err := hasher.NewHashCacheWithPersistence(100, cacheFile)
	if err != nil {
		t.Fatalf("NewHashCacheWithPersistence failed: %v", err)
	}

	if cache == nil {
		t.Fatal("NewHashCacheWithPersistence returned nil")
	}

	stats := cache.Stats()
	if !stats.AutoSave {
		t.Error("AutoSave should be true")
	}
}

// TestHashCache_Stats tests Stats
func TestHashCache_Stats(t *testing.T) {
	cache := hasher.NewHashCache(50)

	stats := cache.Stats()
	if stats.Size != 0 {
		t.Errorf("Initial Size = %d, want 0", stats.Size)
	}
	if stats.MaxSize != 50 {
		t.Errorf("MaxSize = %d, want 50", stats.MaxSize)
	}
	if stats.AutoSave {
		t.Error("AutoSave should be false")
	}

	// Add an entry
	cache.Set("/test/file.txt", "abc123", 1024, time.Now())

	stats = cache.Stats()
	if stats.Size != 1 {
		t.Errorf("Size after Set = %d, want 1", stats.Size)
	}
}

// TestCacheEntry tests CacheEntry struct
func TestCacheEntry(t *testing.T) {
	entry := hasher.CacheEntry{
		Path:     "/test/file.txt",
		Hash:     "abc123",
		Size:     1024,
		ModTime:  time.Now(),
		CachedAt: time.Now(),
	}

	if entry.Path != "/test/file.txt" {
		t.Errorf("Path = %s, want /test/file.txt", entry.Path)
	}
	if entry.Hash != "abc123" {
		t.Errorf("Hash = %s, want abc123", entry.Hash)
	}
	if entry.Size != 1024 {
		t.Errorf("Size = %d, want 1024", entry.Size)
	}
}

// TestCacheStats tests CacheStats struct
func TestCacheStats(t *testing.T) {
	stats := hasher.CacheStats{
		Size:     10,
		MaxSize:  100,
		AutoSave: true,
	}

	if stats.Size != 10 {
		t.Errorf("Size = %d, want 10", stats.Size)
	}
	if stats.MaxSize != 100 {
		t.Errorf("MaxSize = %d, want 100", stats.MaxSize)
	}
	if !stats.AutoSave {
		t.Error("AutoSave should be true")
	}
}

// TestHashCache_ConcurrentAccess tests concurrent access
func TestHashCache_ConcurrentAccess(t *testing.T) {
	cache := hasher.NewHashCache(1000)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			path := "/test/file" + string(rune('0'+id)) + ".txt"
			cache.Set(path, "hash"+string(rune('0'+id)), 1024, time.Now())
			cache.Get(path)
			cache.Size()
		}(i)
	}

	wg.Wait()

	if cache.Size() != 10 {
		t.Errorf("Size = %d, want 10", cache.Size())
	}
}
