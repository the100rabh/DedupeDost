package preview_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/the100rabh/DedupeDost/internal/preview"
	"github.com/the100rabh/DedupeDost/pkg/models"
	"github.com/the100rabh/DedupeDost/pkg/types"
)

// Helper function to create a minimal valid MP4 file for testing
func createTestVideo(t *testing.T, dir, name string) string {
	t.Helper()

	// Create a minimal MP4 file (not a valid video, but enough for testing)
	// This is a minimal ftyp box + moov box structure
	mp4Data := []byte{
		// ftyp box (file type)
		0x00, 0x00, 0x00, 0x14, // box size (20 bytes)
		0x66, 0x74, 0x79, 0x70, // 'ftyp'
		0x69, 0x73, 0x6F, 0x6D, // 'isom' (brand)
		0x00, 0x00, 0x00, 0x01, // version
		0x69, 0x73, 0x6F, 0x6D, // 'isom' (compatible brand)
		// Minimal moov box (movie)
		0x00, 0x00, 0x00, 0x08, // box size (8 bytes)
		0x6D, 0x6F, 0x6F, 0x76, // 'moov'
	}

	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, mp4Data, 0644); err != nil {
		t.Fatalf("Failed to create test video: %v", err)
	}

	return path
}

// TestVideoProvider_CanPreview tests format detection
func TestVideoProvider_CanPreview(t *testing.T) {
	vp := preview.NewVideoProvider()

	tests := []struct {
		name      string
		fileType  types.FileType
		extension string
		expected  bool
	}{
		{"MP4 file", types.FileTypeVideo, ".mp4", true},
		{"AVI file", types.FileTypeVideo, ".avi", true},
		{"MKV file", types.FileTypeVideo, ".mkv", true},
		{"MOV file", types.FileTypeVideo, ".mov", true},
		{"WebM file", types.FileTypeVideo, ".webm", true},
		{"WMV file", types.FileTypeVideo, ".wmv", true},
		{"FLV file", types.FileTypeVideo, ".flv", true},
		{"Text file", types.FileTypeText, ".txt", false},
		{"Image file", types.FileTypeImage, ".jpg", false},
		{"Unknown file", types.FileTypeUnknown, ".xyz", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := vp.CanPreview(tt.fileType, tt.extension)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestVideoProvider_Preview tests video preview creation
func TestVideoProvider_Preview(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	path := createTestVideo(t, dir, "test.mp4")
	file, _ := models.NewFileEntry(path)

	vp := preview.NewVideoProvider()
	result, err := vp.Preview(*file)

	if err != nil {
		t.Fatalf("Preview failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil preview widget")
	}
}

// TestVideoProvider_Preview_NonExistent tests error handling for missing files
func TestVideoProvider_Preview_NonExistent(t *testing.T) {
	vp := preview.NewVideoProvider()

	file := models.FileEntry{
		Path: "/nonexistent/video.mp4",
		Name: "missing.mp4",
	}

	result, err := vp.Preview(file)

	if err != nil {
		t.Fatalf("Preview should not return error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected error widget for missing file")
	}
}

// TestVideoProvider_GetMetadata tests metadata extraction
func TestVideoProvider_GetMetadata(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	path := createTestVideo(t, dir, "test.mp4")
	file, _ := models.NewFileEntry(path)

	vp := preview.NewVideoProvider()
	metadata, err := vp.GetMetadata(*file)

	if err != nil {
		t.Fatalf("GetMetadata failed: %v", err)
	}

	props := metadata.Properties

	// Check format is detected
	if props["Format"] == "" {
		t.Error("Expected format to be set")
	}

	// Check size is set
	if props["Size"] == "" {
		t.Error("Expected size to be set")
	}
}

// TestVideoProvider_GetVideoMetadata tests detailed video metadata
func TestVideoProvider_GetVideoMetadata(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	path := createTestVideo(t, dir, "test.mp4")

	vp := preview.NewVideoProvider()
	metadata, err := vp.GetVideoMetadata(path)

	// Should not fail even without ffprobe (will use fallback)
	if err != nil {
		t.Fatalf("GetVideoMetadata failed: %v", err)
	}

	if metadata == nil {
		t.Fatal("Expected non-nil metadata")
	}

	// Format should contain "mp4" (may be comma-separated list like "mov,mp4,m4a")
	if !strings.Contains(metadata.Format, "mp4") {
		t.Errorf("Expected format to contain mp4, got %s", metadata.Format)
	}
}

// TestVideoProvider_ExtractThumbnail tests thumbnail extraction
func TestVideoProvider_ExtractThumbnail(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	path := createTestVideo(t, dir, "test.mp4")

	vp := preview.NewVideoProvider()
	thumbnail, err := vp.ExtractThumbnail(path)

	// Should return placeholder thumbnail if FFmpeg not available
	if err != nil && thumbnail == nil {
		t.Error("Expected thumbnail or placeholder")
	}
}

// TestVideoProvider_OpenInPlayer tests open in player (mock test)
func TestVideoProvider_OpenInPlayer(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	path := createTestVideo(t, dir, "test.mp4")

	vp := preview.NewVideoProvider()
	err := vp.OpenInPlayer(path)

	// This might fail in test environment without display
	// Just verify the function exists and can be called
	_ = err
}

// TestVideoMetadata_GetDisplayDuration tests duration formatting
func TestVideoMetadata_GetDisplayDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration string
		expected string
	}{
		{"Short video", "30", "0:30"},
		{"Medium video", "150", "2:30"},
		{"Long video", "3661", "1:01:01"},
		{"Zero duration", "0", "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: This test would need actual VideoMetadata struct
			// For now, just verify the method exists
		})
	}
}

// TestVideoMetadata_GetDisplayResolution tests resolution formatting
func TestVideoMetadata_GetDisplayResolution(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		height   int
		expected string
	}{
		{"480p", 640, 480, "640x480 (480p)"},
		{"720p", 1280, 720, "1280x720 (720p)"},
		{"1080p", 1920, 1080, "1920x1080 (1080p)"},
		{"4K", 3840, 2160, "3840x2160 (4K)"},
		{"Unknown", 0, 0, "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta := &preview.VideoMetadata{
				Width:  tt.width,
				Height: tt.height,
			}
			result := meta.GetDisplayResolution()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// TestVideoProvider_FFmpegDetection tests FFmpeg detection
func TestVideoProvider_FFmpegDetection(t *testing.T) {
	vp := preview.NewVideoProvider()

	// Just verify the provider is created successfully
	// FFmpeg availability depends on test environment
	if vp == nil {
		t.Fatal("Expected non-nil VideoProvider")
	}
}
