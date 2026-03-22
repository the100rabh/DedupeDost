package models

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/the100rabh/DedupeDost/pkg/types"
)

// FileEntry represents a file in the file system
type FileEntry struct {
	Path        string         // Absolute file path
	Name        string         // File name with extension
	Size        int64          // File size in bytes
	Extension   string         // Lowercase extension
	ModTime     time.Time      // Last modified time
	CreatedTime time.Time      // Creation time
	Hash        string         // SHA-256 hash
	FileType    types.FileType // Type of file
	IsHidden    bool           // Whether file is hidden
	IsSymlink   bool           // Whether file is a symbolic link
}

// NewFileEntry creates a new FileEntry from a file path
func NewFileEntry(path string) (*FileEntry, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file %s: %w", path, err)
	}

	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Check if it's a symlink
	isSymlink := info.Mode()&os.ModeSymlink != 0

	// Get file info (follow symlink if needed)
	if isSymlink {
		info, err = os.Stat(path)
		if err != nil {
			// Symlink points to non-existent file, use lstat info
			info, _ = os.Lstat(path)
		}
	}

	// Get extension
	ext := strings.ToLower(filepath.Ext(info.Name()))

	// Check if hidden (starts with . on Unix, or has hidden attribute on Windows)
	isHidden := strings.HasPrefix(info.Name(), ".")

	// Get creation time (use modtime as fallback)
	createdTime := info.ModTime()

	return &FileEntry{
		Path:        absPath,
		Name:        info.Name(),
		Size:        info.Size(),
		Extension:   ext,
		ModTime:     info.ModTime(),
		CreatedTime: createdTime,
		FileType:    types.GetFileTypeFromExtension(ext),
		IsHidden:    isHidden,
		IsSymlink:   isSymlink,
	}, nil
}

// GetDisplaySize returns a human-readable file size
func (f *FileEntry) GetDisplaySize() string {
	return FormatSize(f.Size)
}

// FormatSize converts bytes to human-readable format
func FormatSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	switch {
	case size >= TB:
		return fmt.Sprintf("%.2f TB", float64(size)/TB)
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	default:
		return fmt.Sprintf("%d B", size)
	}
}

// GetRelativePath returns the path relative to a base directory
func (f *FileEntry) GetRelativePath(baseDir string) (string, error) {
	rel, err := filepath.Rel(baseDir, f.Path)
	if err != nil {
		return f.Path, err
	}
	return rel, nil
}

// MatchesExtension checks if the file matches any of the given extensions
func (f *FileEntry) MatchesExtension(extensions []string) bool {
	if len(extensions) == 0 {
		return true
	}
	for _, ext := range extensions {
		if strings.ToLower(f.Extension) == strings.ToLower(ext) {
			return true
		}
	}
	return false
}
