package types_test

import (
	"testing"

	"github.com/dupdel/dup-del/pkg/types"
)

// TestGetFileTypeFromExtension tests GetFileTypeFromExtension function
func TestGetFileTypeFromExtension(t *testing.T) {
	tests := []struct {
		ext      string
		expected types.FileType
	}{
		{".txt", types.FileTypeText},
		{".TXT", types.FileTypeText},
		{".md", types.FileTypeText},
		{".json", types.FileTypeText},
		{".xml", types.FileTypeText},
		{".yaml", types.FileTypeText},
		{".yml", types.FileTypeText},
		{".csv", types.FileTypeText},
		{".log", types.FileTypeText},
		{".html", types.FileTypeText},
		{".css", types.FileTypeText},
		{".js", types.FileTypeText},
		{".ts", types.FileTypeText},
		{".go", types.FileTypeText},
		{".py", types.FileTypeText},
		{".java", types.FileTypeText},
		{".c", types.FileTypeText},
		{".cpp", types.FileTypeText},
		{".h", types.FileTypeText},
		{".rs", types.FileTypeText},
		{".rb", types.FileTypeText},
		{".php", types.FileTypeText},
		{".sh", types.FileTypeText},
		{".sql", types.FileTypeText},
		{".jpg", types.FileTypeImage},
		{".JPG", types.FileTypeImage},
		{".jpeg", types.FileTypeImage},
		{".png", types.FileTypeImage},
		{".gif", types.FileTypeImage},
		{".bmp", types.FileTypeImage},
		{".webp", types.FileTypeImage},
		{".svg", types.FileTypeImage},
		{".ico", types.FileTypeImage},
		{".tiff", types.FileTypeImage},
		{".tif", types.FileTypeImage},
		{".raw", types.FileTypeImage},
		{".heic", types.FileTypeImage},
		{".heif", types.FileTypeImage},
		{".mp4", types.FileTypeVideo},
		{".MP4", types.FileTypeVideo},
		{".avi", types.FileTypeVideo},
		{".mkv", types.FileTypeVideo},
		{".mov", types.FileTypeVideo},
		{".wmv", types.FileTypeVideo},
		{".flv", types.FileTypeVideo},
		{".webm", types.FileTypeVideo},
		{".m4v", types.FileTypeVideo},
		{".mpg", types.FileTypeVideo},
		{".mpeg", types.FileTypeVideo},
		{".3gp", types.FileTypeVideo},
		{".unknown", types.FileTypeOther},
		{".xyz", types.FileTypeOther},
		{"", types.FileTypeOther},
	}

	for _, tt := range tests {
		result := types.GetFileTypeFromExtension(tt.ext)
		if result != tt.expected {
			t.Errorf("GetFileTypeFromExtension(%s) = %d, want %d", tt.ext, result, tt.expected)
		}
	}
}

// TestSupportedTextExtensions tests SupportedTextExtensions
func TestSupportedTextExtensions(t *testing.T) {
	extensions := types.SupportedTextExtensions()
	if len(extensions) == 0 {
		t.Error("SupportedTextExtensions should return non-empty slice")
	}

	// Check for some expected extensions
	expected := []string{".txt", ".md", ".json", ".go", ".py"}
	for _, exp := range expected {
		found := false
		for _, ext := range extensions {
			if ext == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected extension %s not found in SupportedTextExtensions", exp)
		}
	}
}

// TestSupportedImageExtensions tests SupportedImageExtensions
func TestSupportedImageExtensions(t *testing.T) {
	extensions := types.SupportedImageExtensions()
	if len(extensions) == 0 {
		t.Error("SupportedImageExtensions should return non-empty slice")
	}

	// Check for some expected extensions
	expected := []string{".jpg", ".png", ".gif", ".bmp"}
	for _, exp := range expected {
		found := false
		for _, ext := range extensions {
			if ext == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected extension %s not found in SupportedImageExtensions", exp)
		}
	}
}

// TestSupportedVideoExtensions tests SupportedVideoExtensions
func TestSupportedVideoExtensions(t *testing.T) {
	extensions := types.SupportedVideoExtensions()
	if len(extensions) == 0 {
		t.Error("SupportedVideoExtensions should return non-empty slice")
	}

	// Check for some expected extensions
	expected := []string{".mp4", ".avi", ".mkv", ".mov"}
	for _, exp := range expected {
		found := false
		for _, ext := range extensions {
			if ext == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected extension %s not found in SupportedVideoExtensions", exp)
		}
	}
}

// TestIsTextExtension tests IsTextExtension
func TestIsTextExtension(t *testing.T) {
	tests := []struct {
		ext      string
		expected bool
	}{
		{".txt", true},
		{".TXT", true},
		{".go", true},
		{".jpg", false},
		{".mp4", false},
		{".unknown", false},
	}

	for _, tt := range tests {
		result := types.IsTextExtension(tt.ext)
		if result != tt.expected {
			t.Errorf("IsTextExtension(%s) = %v, want %v", tt.ext, result, tt.expected)
		}
	}
}

// TestIsImageExtension tests IsImageExtension
func TestIsImageExtension(t *testing.T) {
	tests := []struct {
		ext      string
		expected bool
	}{
		{".jpg", true},
		{".JPG", true},
		{".png", true},
		{".txt", false},
		{".mp4", false},
		{".unknown", false},
	}

	for _, tt := range tests {
		result := types.IsImageExtension(tt.ext)
		if result != tt.expected {
			t.Errorf("IsImageExtension(%s) = %v, want %v", tt.ext, result, tt.expected)
		}
	}
}

// TestIsVideoExtension tests IsVideoExtension
func TestIsVideoExtension(t *testing.T) {
	tests := []struct {
		ext      string
		expected bool
	}{
		{".mp4", true},
		{".MP4", true},
		{".avi", true},
		{".txt", false},
		{".jpg", false},
		{".unknown", false},
	}

	for _, tt := range tests {
		result := types.IsVideoExtension(tt.ext)
		if result != tt.expected {
			t.Errorf("IsVideoExtension(%s) = %v, want %v", tt.ext, result, tt.expected)
		}
	}
}

// TestFileType_Constants tests FileType constants
func TestFileType_Constants(t *testing.T) {
	if types.FileTypeUnknown != 0 {
		t.Errorf("FileTypeUnknown = %d, want 0", types.FileTypeUnknown)
	}
	if types.FileTypeText != 1 {
		t.Errorf("FileTypeText = %d, want 1", types.FileTypeText)
	}
	if types.FileTypeImage != 2 {
		t.Errorf("FileTypeImage = %d, want 2", types.FileTypeImage)
	}
	if types.FileTypeVideo != 3 {
		t.Errorf("FileTypeVideo = %d, want 3", types.FileTypeVideo)
	}
	if types.FileTypeOther != 4 {
		t.Errorf("FileTypeOther = %d, want 4", types.FileTypeOther)
	}
}
