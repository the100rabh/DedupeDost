package scanner

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dupdel/dup-del/pkg/models"
)

// ScanOptions contains options for scanning directories
type ScanOptions struct {
	Recursive      bool     // Scan subdirectories
	IncludeHidden  bool     // Include hidden files
	FollowSymlinks bool     // Follow symbolic links
	MinSize        int64    // Minimum file size in bytes
	MaxSize        int64    // Maximum file size in bytes (0 = unlimited)
	Extensions     []string // File extensions to include (empty = all)
}

// ScanProgress represents the current scan progress
type ScanProgress struct {
	FilesScanned    int64   // Number of files scanned so far
	TotalFiles      int64   // Estimated total files (0 = unknown)
	CurrentFile     string  // Currently processing file
	TotalSize       int64   // Total size of scanned files
	PercentComplete float64 // Percentage complete (0-100)
}

// Scanner scans directories for files
type Scanner struct {
	rootDir string
	options ScanOptions
	progress *ScanProgress
	ctx      context.Context
	cancel   context.CancelFunc
	mu       sync.RWMutex
}

// NewScanner creates a new Scanner for the given directory
func NewScanner(rootDir string, options ScanOptions) *Scanner {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &Scanner{
		rootDir: rootDir,
		options: options,
		progress: &ScanProgress{},
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Scan starts scanning the directory and returns channels for files and errors
func (s *Scanner) Scan() (<-chan models.FileEntry, <-chan error, error) {
	fileChan := make(chan models.FileEntry, 100)
	errorChan := make(chan error, 10)

	// Get absolute path
	absDir, err := filepath.Abs(s.rootDir)
	if err != nil {
		return nil, nil, err
	}
	s.rootDir = absDir

	go func() {
		defer close(fileChan)
		defer close(errorChan)

		err := s.walkDirectory(fileChan, errorChan)
		if err != nil && err != context.Canceled {
			select {
			case errorChan <- err:
			default:
			}
		}
	}()

	return fileChan, errorChan, nil
}

// walkDirectory walks the directory tree and sends files to the channel
func (s *Scanner) walkDirectory(fileChan chan<- models.FileEntry, errorChan chan<- error) error {
	walkFn := func(path string, d fs.DirEntry, err error) error {
		// Check for cancellation
		select {
		case <-s.ctx.Done():
			return context.Canceled
		default:
		}

		if err != nil {
			// Log error but continue scanning
			select {
			case errorChan <- err:
			default:
			}
			return nil // Continue scanning
		}

		// Skip directories (we only want files)
		if d.IsDir() {
			// Skip hidden directories if not including hidden
			if !s.options.IncludeHidden && strings.HasPrefix(d.Name(), ".") {
				if s.options.Recursive {
					return filepath.SkipDir
				}
				return nil
			}
			return nil
		}

		// Check if file should be included
		if !s.shouldIncludeFile(d) {
			return nil
		}

		// Create FileEntry
		entry, err := models.NewFileEntry(path)
		if err != nil {
			select {
			case errorChan <- err:
			default:
			}
			return nil
		}

		// Update progress
		s.updateProgress(path, entry.Size)

		// Send file entry
		select {
		case <-s.ctx.Done():
			return context.Canceled
		case fileChan <- *entry:
		}

		return nil
	}

	if s.options.Recursive {
		return filepath.WalkDir(s.rootDir, walkFn)
	}

	// Non-recursive: only scan immediate children
	entries, err := os.ReadDir(s.rootDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := filepath.Join(s.rootDir, entry.Name())
		if err := walkFn(path, entry, nil); err != nil {
			return err
		}
	}

	return nil
}

// shouldIncludeFile checks if a file should be included based on options
func (s *Scanner) shouldIncludeFile(d fs.DirEntry) bool {
	name := d.Name()

	// Check hidden
	if !s.options.IncludeHidden && strings.HasPrefix(name, ".") {
		return false
	}

	// Get file info for size check
	info, err := d.Info()
	if err != nil {
		return false
	}

	// Check size filters
	if !FilterBySize(info.Size(), s.options.MinSize, s.options.MaxSize) {
		return false
	}

	// Check extension filter
	if len(s.options.Extensions) > 0 {
		ext := strings.ToLower(filepath.Ext(name))
		if !FilterByExtension(ext, s.options.Extensions) {
			return false
		}
	}

	return true
}

// updateProgress updates the scan progress
func (s *Scanner) updateProgress(currentFile string, size int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.progress.FilesScanned++
	s.progress.CurrentFile = currentFile
	s.progress.TotalSize += size
	
	// Calculate percentage if we know total
	if s.progress.TotalFiles > 0 {
		s.progress.PercentComplete = float64(s.progress.FilesScanned) / float64(s.progress.TotalFiles) * 100
	}
}

// GetProgress returns the current scan progress
func (s *Scanner) GetProgress() ScanProgress {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	return *s.progress
}

// Cancel stops the scan
func (s *Scanner) Cancel() {
	s.cancel()
}

// CountFiles performs a quick count of files to be scanned
func (s *Scanner) CountFiles() (int64, error) {
	var count int64
	var totalSize int64

	walkFn := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Continue on errors
		}

		if d.IsDir() {
			if !s.options.IncludeHidden && strings.HasPrefix(d.Name(), ".") {
				if s.options.Recursive {
					return filepath.SkipDir
				}
			}
			return nil
		}

		if !s.shouldIncludeFile(d) {
			return nil
		}

		info, err := d.Info()
		if err == nil {
			count++
			totalSize += info.Size()
		}

		return nil
	}

	if s.options.Recursive {
		if err := filepath.WalkDir(s.rootDir, walkFn); err != nil {
			return 0, err
		}
	} else {
		entries, err := os.ReadDir(s.rootDir)
		if err != nil {
			return 0, err
		}
		for _, entry := range entries {
			if !entry.IsDir() && s.shouldIncludeFile(entry) {
				count++
			}
		}
	}

	s.mu.Lock()
	s.progress.TotalFiles = count
	s.progress.TotalSize = totalSize
	s.mu.Unlock()

	return count, nil
}

// FilterBySize checks if a file size matches the filter criteria
func FilterBySize(size, minSize, maxSize int64) bool {
	if size < minSize {
		return false
	}
	if maxSize > 0 && size > maxSize {
		return false
	}
	return true
}

// FilterByExtension checks if an extension matches the filter criteria
func FilterByExtension(ext string, allowed []string) bool {
	if len(allowed) == 0 {
		return true
	}
	ext = strings.ToLower(ext)
	for _, allowedExt := range allowed {
		if ext == strings.ToLower(allowedExt) {
			return true
		}
	}
	return false
}
