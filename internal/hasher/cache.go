package hasher

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// CacheEntry represents a cached hash result
type CacheEntry struct {
	Path      string    `json:"path"`
	Hash      string    `json:"hash"`
	Size      int64     `json:"size"`
	ModTime   time.Time `json:"mod_time"`
	CachedAt  time.Time `json:"cached_at"`
}

// HashCache provides caching for computed hashes
type HashCache struct {
	cache     map[string]CacheEntry
	maxSize   int
	mu        sync.RWMutex
	filePath  string
	autoSave  bool
}

// NewHashCache creates a new hash cache
func NewHashCache(maxSize int) *HashCache {
	if maxSize <= 0 {
		maxSize = 10000
	}
	
	return &HashCache{
		cache:    make(map[string]CacheEntry, maxSize),
		maxSize:  maxSize,
		autoSave: false,
	}
}

// NewHashCacheWithPersistence creates a cache with file persistence
func NewHashCacheWithPersistence(maxSize int, cacheFile string) (*HashCache, error) {
	cache := NewHashCache(maxSize)
	cache.filePath = cacheFile
	cache.autoSave = true
	
	// Load existing cache
	if err := cache.Load(); err != nil {
		return nil, err
	}
	
	return cache, nil
}

// Get retrieves a cached hash for a file path
func (c *HashCache) Get(path string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	entry, ok := c.cache[path]
	if !ok {
		return "", false
	}
	
	// Verify file hasn't changed
	info, err := os.Stat(path)
	if err != nil {
		return "", false
	}
	
	if info.ModTime() != entry.ModTime || info.Size() != entry.Size {
		return "", false
	}
	
	return entry.Hash, true
}

// Set stores a hash in the cache
func (c *HashCache) Set(path, hash string, size int64, modTime time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Check if we need to evict
	if len(c.cache) >= c.maxSize {
		c.evictOldest()
	}
	
	c.cache[path] = CacheEntry{
		Path:     path,
		Hash:     hash,
		Size:     size,
		ModTime:  modTime,
		CachedAt: time.Now(),
	}
	
	if c.autoSave {
		go c.Save()
	}
}

// GetOrCompute gets a cached hash or computes a new one
func (c *HashCache) GetOrCompute(path string, hasher *Hasher) (string, error) {
	// Try cache first
	if hash, ok := c.Get(path); ok {
		return hash, nil
	}
	
	// Compute new hash
	hash, err := hasher.CalculateHash(path)
	if err != nil {
		return "", err
	}
	
	// Get file info
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	
	// Cache the result
	c.Set(path, hash, info.Size(), info.ModTime())
	
	return hash, nil
}

// Remove removes a path from the cache
func (c *HashCache) Remove(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, path)
	
	if c.autoSave {
		go c.Save()
	}
}

// Clear clears the entire cache
func (c *HashCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = make(map[string]CacheEntry, c.maxSize)
	
	if c.autoSave {
		go c.Save()
	}
}

// Size returns the number of entries in the cache
func (c *HashCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.cache)
}

// evictOldest removes the oldest entry from the cache
func (c *HashCache) evictOldest() {
	var oldestPath string
	var oldestTime time.Time
	
	for path, entry := range c.cache {
		if oldestPath == "" || entry.CachedAt.Before(oldestTime) {
			oldestPath = path
			oldestTime = entry.CachedAt
		}
	}
	
	if oldestPath != "" {
		delete(c.cache, oldestPath)
	}
}

// Save saves the cache to disk
func (c *HashCache) Save() error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	if c.filePath == "" {
		return nil
	}
	
	// Ensure directory exists
	dir := filepath.Dir(c.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	// Convert to slice for JSON
	entries := make([]CacheEntry, 0, len(c.cache))
	for _, entry := range c.cache {
		entries = append(entries, entry)
	}
	
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(c.filePath, data, 0644)
}

// Load loads the cache from disk
func (c *HashCache) Load() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if c.filePath == "" {
		return nil
	}
	
	data, err := os.ReadFile(c.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No cache file yet
		}
		return err
	}
	
	var entries []CacheEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return err
	}
	
	c.cache = make(map[string]CacheEntry, len(entries))
	for _, entry := range entries {
		c.cache[entry.Path] = entry
	}
	
	return nil
}

// Stats returns cache statistics
func (c *HashCache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	return CacheStats{
		Size:      len(c.cache),
		MaxSize:   c.maxSize,
		AutoSave:  c.autoSave,
	}
}

// CacheStats contains cache statistics
type CacheStats struct {
	Size     int  `json:"size"`
	MaxSize  int  `json:"max_size"`
	AutoSave bool `json:"auto_save"`
}
