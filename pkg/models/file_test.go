package models_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/the100rabh/DedupeDost/pkg/models"
)

// TestFormatSize tests the FormatSize function
func TestFormatSize(t *testing.T) {
	tests := []struct {
		size     int64
		expected string
	}{
		{0, "0 B"},
		{1, "1 B"},
		{1023, "1023 B"},
		{1024, "1.00 KB"},
		{1536, "1.50 KB"},
		{1048576, "1.00 MB"},
		{1572864, "1.50 MB"},
		{1073741824, "1.00 GB"},
		{1610612736, "1.50 GB"},
		{1099511627776, "1.00 TB"},
		{1649267441664, "1.50 TB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := models.FormatSize(tt.size)
			if result != tt.expected {
				t.Errorf("FormatSize(%d) = %s, want %s", tt.size, result, tt.expected)
			}
		})
	}
}

// TestNewFileEntry tests creating a FileEntry from a file path
func TestNewFileEntry(t *testing.T) {
	// Create temp file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	content := []byte("test content")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Wait a bit to ensure modtime is set
	time.Sleep(10 * time.Millisecond)

	entry, err := models.NewFileEntry(tmpFile)
	if err != nil {
		t.Fatalf("NewFileEntry failed: %v", err)
	}

	if entry.Name != "test.txt" {
		t.Errorf("Name = %s, want test.txt", entry.Name)
	}

	if entry.Size != int64(len(content)) {
		t.Errorf("Size = %d, want %d", entry.Size, len(content))
	}

	if entry.Extension != ".txt" {
		t.Errorf("Extension = %s, want .txt", entry.Extension)
	}

	if entry.IsHidden {
		t.Error("IsHidden = true, want false")
	}

	if entry.IsSymlink {
		t.Error("IsSymlink = true, want false")
	}
}

// TestNewFileEntry_NonExistent tests error handling for non-existent file
func TestNewFileEntry_NonExistent(t *testing.T) {
	_, err := models.NewFileEntry("/non/existent/file.txt")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

// TestNewFileEntry_HiddenFile tests detection of hidden files
func TestNewFileEntry_HiddenFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, ".hidden")
	if err := os.WriteFile(tmpFile, []byte("hidden"), 0644); err != nil {
		t.Fatalf("Failed to create hidden file: %v", err)
	}

	entry, err := models.NewFileEntry(tmpFile)
	if err != nil {
		t.Fatalf("NewFileEntry failed: %v", err)
	}

	if !entry.IsHidden {
		t.Error("IsHidden = false, want true for hidden file")
	}
}

// TestNewFileEntry_Symlink tests symlink detection
func TestNewFileEntry_Symlink(t *testing.T) {
	tmpDir := t.TempDir()
	original := filepath.Join(tmpDir, "original.txt")
	link := filepath.Join(tmpDir, "link.txt")

	if err := os.WriteFile(original, []byte("original"), 0644); err != nil {
		t.Fatalf("Failed to create original file: %v", err)
	}

	if err := os.Symlink(original, link); err != nil {
		t.Fatalf("Failed to create symlink: %v", err)
	}

	entry, err := models.NewFileEntry(link)
	if err != nil {
		t.Fatalf("NewFileEntry failed: %v", err)
	}

	if !entry.IsSymlink {
		t.Error("IsSymlink = false, want true for symlink")
	}
}

// TestFileEntry_GetDisplaySize tests GetDisplaySize method
func TestFileEntry_GetDisplaySize(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	content := []byte("test content 12345")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	entry, err := models.NewFileEntry(tmpFile)
	if err != nil {
		t.Fatalf("NewFileEntry failed: %v", err)
	}

	size := entry.GetDisplaySize()
	expected := models.FormatSize(entry.Size)
	if size != expected {
		t.Errorf("GetDisplaySize() = %s, want %s", size, expected)
	}
}

// TestFileEntry_GetRelativePath tests GetRelativePath method
func TestFileEntry_GetRelativePath(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "subdir", "test.txt")
	if err := os.MkdirAll(filepath.Dir(tmpFile), 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	if err := os.WriteFile(tmpFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	entry, err := models.NewFileEntry(tmpFile)
	if err != nil {
		t.Fatalf("NewFileEntry failed: %v", err)
	}

	relPath, err := entry.GetRelativePath(tmpDir)
	if err != nil {
		t.Fatalf("GetRelativePath failed: %v", err)
	}

	expected := filepath.Join("subdir", "test.txt")
	if relPath != expected {
		t.Errorf("GetRelativePath() = %s, want %s", relPath, expected)
	}
}

// TestFileEntry_MatchesExtension tests MatchesExtension method
func TestFileEntry_MatchesExtension(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.TXT")
	if err := os.WriteFile(tmpFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	entry, err := models.NewFileEntry(tmpFile)
	if err != nil {
		t.Fatalf("NewFileEntry failed: %v", err)
	}

	tests := []struct {
		extensions []string
		expected   bool
	}{
		{[]string{}, true},                        // Empty = match all
		{[]string{".txt"}, true},                  // Exact match
		{[]string{".TXT"}, true},                  // Case insensitive
		{[]string{".txt", ".pdf"}, true},          // Match in list
		{[]string{".pdf", ".doc"}, false},         // No match
		{[]string{".TXT"}, true},                  // Case match
	}

	for _, tt := range tests {
		result := entry.MatchesExtension(tt.extensions)
		if result != tt.expected {
			t.Errorf("MatchesExtension(%v) = %v, want %v", tt.extensions, result, tt.expected)
		}
	}
}
