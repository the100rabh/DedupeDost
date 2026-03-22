package hasher

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"sync"
)

const (
	// DefaultChunkSize is the default chunk size for reading files (4MB)
	DefaultChunkSize int64 = 4 * 1024 * 1024
	
	// DefaultLargeFileThreshold is the threshold for considering a file "large" (100MB)
	DefaultLargeFileThreshold int64 = 100 * 1024 * 1024
	
	// DefaultAlgorithm is the default hashing algorithm
	DefaultAlgorithm = "sha256"
)

// HashResult contains the result of a hash calculation
type HashResult struct {
	Path      string
	Hash      string
	Size      int64
	Error     error
	Cancelled bool
}

// Hasher calculates file hashes
type Hasher struct {
	algorithm            string
	chunkSize            int64
	largeFileThreshold   int64
	hashPool             sync.Pool
}

// NewHasher creates a new Hasher with default settings
func NewHasher() *Hasher {
	h := &Hasher{
		algorithm:            DefaultAlgorithm,
		chunkSize:            DefaultChunkSize,
		largeFileThreshold:   DefaultLargeFileThreshold,
	}
	
	// Initialize hash pool for reusing hash objects
	h.hashPool = sync.Pool{
		New: func() interface{} {
			return sha256.New()
		},
	}
	
	return h
}

// NewHasherWithOptions creates a new Hasher with custom options
func NewHasherWithOptions(chunkSize, largeFileThreshold int64) *Hasher {
	h := &Hasher{
		algorithm:            DefaultAlgorithm,
		chunkSize:            chunkSize,
		largeFileThreshold:   largeFileThreshold,
	}
	
	h.hashPool = sync.Pool{
		New: func() interface{} {
			return sha256.New()
		},
	}
	
	return h
}

// CalculateHash calculates the SHA-256 hash of a file
func (h *Hasher) CalculateHash(filePath string) (string, error) {
	return h.CalculateHashWithContext(filePath, context.Background())
}

// CalculateHashWithContext calculates the hash with context for cancellation
func (h *Hasher) CalculateHashWithContext(filePath string, ctx context.Context) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	
	// Get file info for size
	info, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to stat file: %w", err)
	}
	
	// Get hash from pool
	hasher := h.hashPool.Get().(hash.Hash)
	defer h.hashPool.Put(hasher)
	
	hasher.Reset()
	
	// Use buffered reading with appropriate buffer size
	var bufSize int
	if info.Size() > h.largeFileThreshold {
		bufSize = int(h.chunkSize)
	} else {
		bufSize = 32 * 1024 // 32KB for smaller files
	}
	
	reader := bufio.NewReaderSize(file, bufSize)
	
	// Create context-aware reader
	limitedReader := &contextReader{
		reader:  reader,
		context: ctx,
	}
	
	_, err = io.Copy(hasher, limitedReader)
	if err != nil {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	
	// Get hex-encoded hash
	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes), nil
}

// CalculateHashWithProgress calculates hash with progress callbacks
func (h *Hasher) CalculateHashWithProgress(filePath string, progressChan chan<- float64) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	
	// Get file info for size
	info, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to stat file: %w", err)
	}
	
	fileSize := info.Size()
	
	// Get hash from pool
	hasher := h.hashPool.Get().(hash.Hash)
	defer h.hashPool.Put(hasher)
	
	hasher.Reset()
	
	// Read file in chunks and report progress
	buf := make([]byte, h.chunkSize)
	var totalRead int64 = 0
	
	for {
		n, err := file.Read(buf)
		if n > 0 {
			hasher.Write(buf[:n])
			totalRead += int64(n)
			
			// Report progress
			if progressChan != nil && fileSize > 0 {
				progress := float64(totalRead) / float64(fileSize) * 100
				select {
				case progressChan <- progress:
				default:
					// Channel not ready, skip
				}
			}
		}
		
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("failed to read file: %w", err)
		}
	}
	
	// Send final progress
	if progressChan != nil {
		select {
		case progressChan <- 100.0:
		default:
		}
	}
	
	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes), nil
}

// HashMultiple calculates hashes for multiple files using a worker pool
func (h *Hasher) HashMultiple(files []string, maxWorkers int) <-chan HashResult {
	resultChan := make(chan HashResult, len(files))
	
	if len(files) == 0 {
		close(resultChan)
		return resultChan
	}
	
	// Limit workers
	if maxWorkers <= 0 {
		maxWorkers = 4
	}
	if maxWorkers > len(files) {
		maxWorkers = len(files)
	}
	
	// Create worker pool
	var wg sync.WaitGroup
	fileChan := make(chan string, len(files))
	
	// Start workers
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filePath := range fileChan {
				hash, err := h.CalculateHash(filePath)
				result := HashResult{
					Path: filePath,
					Hash: hash,
					Error: err,
				}
				
				if err != nil {
					result.Error = err
				}
				
				resultChan <- result
			}
		}()
	}
	
	// Send files to workers
	go func() {
		for _, file := range files {
			fileChan <- file
		}
		close(fileChan)
	}()
	
	// Close result channel when all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	
	return resultChan
}

// HashMultipleWithContext calculates hashes with context support
func (h *Hasher) HashMultipleWithContext(ctx context.Context, files []string, maxWorkers int) <-chan HashResult {
	resultChan := make(chan HashResult, len(files))
	
	if len(files) == 0 {
		close(resultChan)
		return resultChan
	}
	
	if maxWorkers <= 0 {
		maxWorkers = 4
	}
	if maxWorkers > len(files) {
		maxWorkers = len(files)
	}
	
	var wg sync.WaitGroup
	fileChan := make(chan string, len(files))
	
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filePath := range fileChan {
				// Check context before processing
				select {
				case <-ctx.Done():
					resultChan <- HashResult{
						Path:      filePath,
						Cancelled: true,
						Error:     ctx.Err(),
					}
					continue
				default:
				}
				
				hash, err := h.CalculateHashWithContext(filePath, ctx)
				result := HashResult{
					Path: filePath,
					Hash: hash,
					Size: h.getFileSize(filePath),
					Error: err,
				}
				
				if err != nil {
					if ctx.Err() != nil {
						result.Cancelled = true
					}
					result.Error = err
				}
				
				resultChan <- result
			}
		}()
	}
	
	go func() {
		for _, file := range files {
			select {
			case <-ctx.Done():
				break
			case fileChan <- file:
			}
		}
		close(fileChan)
	}()
	
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	
	return resultChan
}

// getFileSize gets the size of a file
func (h *Hasher) getFileSize(filePath string) int64 {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0
	}
	return info.Size()
}

// contextReader wraps an io.Reader to respect context cancellation
type contextReader struct {
	reader  io.Reader
	context context.Context
}

func (r *contextReader) Read(p []byte) (n int, err error) {
	select {
	case <-r.context.Done():
		return 0, r.context.Err()
	default:
		return r.reader.Read(p)
	}
}
