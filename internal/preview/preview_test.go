//go:build !gui

package preview

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dupdel/dup-del/pkg/models"
	"github.com/dupdel/dup-del/pkg/types"
)

// TestMetadata tests Metadata struct
func TestMetadata(t *testing.T) {
	metadata := Metadata{
		Properties: map[string]string{
			"Name": "test.txt",
			"Size": "1 KB",
		},
		VideoMetadata: nil,
	}

	if len(metadata.Properties) != 2 {
		t.Errorf("Properties length = %d, want 2", len(metadata.Properties))
	}
	if metadata.Properties["Name"] != "test.txt" {
		t.Errorf("Name = %s, want test.txt", metadata.Properties["Name"])
	}
	if metadata.VideoMetadata != nil {
		t.Error("VideoMetadata should be nil")
	}
}

// TestNewTextProvider tests NewTextProvider
func TestNewTextProvider(t *testing.T) {
	provider := NewTextProvider()
	if provider == nil {
		t.Fatal("NewTextProvider returned nil")
	}
	if !provider.wordWrap {
		t.Error("wordWrap should be true by default")
	}
}

// TestTextProvider_CanPreview tests TextProvider.CanPreview
func TestTextProvider_CanPreview(t *testing.T) {
	provider := NewTextProvider()

	tests := []struct {
		fileType types.FileType
		ext      string
		want     bool
	}{
		{types.FileTypeText, ".txt", true},
		{types.FileTypeText, ".md", true},
		{types.FileTypeText, ".json", true},
		{types.FileTypeImage, ".jpg", false},
		{types.FileTypeVideo, ".mp4", false},
		{types.FileTypeOther, ".xyz", false},
	}

	for _, tt := range tests {
		got := provider.CanPreview(tt.fileType, tt.ext)
		if got != tt.want {
			t.Errorf("CanPreview(%v, %s) = %v, want %v", tt.fileType, tt.ext, got, tt.want)
		}
	}
}

// TestTextProvider_GetMetadata tests TextProvider.GetMetadata
func TestTextProvider_GetMetadata(t *testing.T) {
	provider := NewTextProvider()

	// Create temp file
	tmpFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write test content
	content := "Hello World\nLine 2\nLine 3"
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	file := models.FileEntry{
		Path: tmpFile.Name(),
		Name: filepath.Base(tmpFile.Name()),
	}

	metadata, err := provider.GetMetadata(file)
	if err != nil {
		t.Fatalf("GetMetadata failed: %v", err)
	}

	if metadata.Properties["Lines"] != "3" {
		t.Errorf("Lines = %s, want 3", metadata.Properties["Lines"])
	}
	// "Hello World\nLine 2\nLine 3" has 6 words
	if metadata.Properties["Words"] != "6" {
		t.Errorf("Words = %s, want 6", metadata.Properties["Words"])
	}
	// 25 chars
	if metadata.Properties["Chars"] != "25" {
		t.Errorf("Chars = %s, want 25", metadata.Properties["Chars"])
	}
}

// TestTextProvider_GetMetadata_NonExistent tests GetMetadata with non-existent file
func TestTextProvider_GetMetadata_NonExistent(t *testing.T) {
	provider := NewTextProvider()

	file := models.FileEntry{
		Path: "/nonexistent/file.txt",
	}

	_, err := provider.GetMetadata(file)
	if err == nil {
		t.Error("GetMetadata should return error for non-existent file")
	}
}

// TestTextProvider_GetMetadata_EmptyFile tests GetMetadata with empty file
func TestTextProvider_GetMetadata_EmptyFile(t *testing.T) {
	provider := NewTextProvider()

	tmpFile, err := os.CreateTemp("", "empty_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	file := models.FileEntry{
		Path: tmpFile.Name(),
		Name: filepath.Base(tmpFile.Name()),
	}

	metadata, err := provider.GetMetadata(file)
	if err != nil {
		t.Fatalf("GetMetadata failed: %v", err)
	}

	if metadata.Properties["Lines"] != "1" {
		t.Errorf("Lines for empty file = %s, want 1", metadata.Properties["Lines"])
	}
	if metadata.Properties["Words"] != "0" {
		t.Errorf("Words for empty file = %s, want 0", metadata.Properties["Words"])
	}
	if metadata.Properties["Chars"] != "0" {
		t.Errorf("Chars for empty file = %s, want 0", metadata.Properties["Chars"])
	}
}

// TestNewFallbackProvider tests NewFallbackProvider
func TestNewFallbackProvider(t *testing.T) {
	provider := NewFallbackProvider()
	if provider == nil {
		t.Fatal("NewFallbackProvider returned nil")
	}
}

// TestFallbackProvider_CanPreview tests FallbackProvider.CanPreview
func TestFallbackProvider_CanPreview(t *testing.T) {
	provider := NewFallbackProvider()

	// Should always return true
	tests := []struct {
		fileType types.FileType
		ext      string
	}{
		{types.FileTypeText, ".txt"},
		{types.FileTypeImage, ".jpg"},
		{types.FileTypeVideo, ".mp4"},
		{types.FileTypeOther, ".xyz"},
	}

	for _, tt := range tests {
		if !provider.CanPreview(tt.fileType, tt.ext) {
			t.Errorf("CanPreview(%v, %s) should return true", tt.fileType, tt.ext)
		}
	}
}

// TestFallbackProvider_GetMetadata tests FallbackProvider.GetMetadata
func TestFallbackProvider_GetMetadata(t *testing.T) {
	provider := NewFallbackProvider()

	file := models.FileEntry{
		Name:      "test.txt",
		Size:      1024,
		Extension: ".txt",
		Path:      "/test/path/test.txt",
	}

	metadata, err := provider.GetMetadata(file)
	if err != nil {
		t.Fatalf("GetMetadata failed: %v", err)
	}

	if metadata.Properties["Name"] != "test.txt" {
		t.Errorf("Name = %s, want test.txt", metadata.Properties["Name"])
	}
	if metadata.Properties["Size"] == "" {
		t.Error("Size should not be empty")
	}
	if metadata.Properties["Type"] != ".txt" {
		t.Errorf("Type = %s, want .txt", metadata.Properties["Type"])
	}
	if metadata.Properties["Path"] != "/test/path/test.txt" {
		t.Errorf("Path = %s, want /test/path/test.txt", metadata.Properties["Path"])
	}
	if metadata.Properties["Modified"] == "" {
		t.Error("Modified should not be empty")
	}
}

// TestNewRegistry tests NewRegistry
func TestNewRegistry(t *testing.T) {
	registry := NewRegistry()
	if registry == nil {
		t.Fatal("NewRegistry returned nil")
	}
	if len(registry.providers) == 0 {
		t.Error("Registry should have providers")
	}
}

// TestRegistry_Register tests Registry.Register
func TestRegistry_Register(t *testing.T) {
	registry := NewRegistry()
	initialCount := len(registry.providers)

	// Register a new provider
	provider := NewFallbackProvider()
	registry.Register(provider)

	if len(registry.providers) != initialCount+1 {
		t.Errorf("Provider count = %d, want %d", len(registry.providers), initialCount+1)
	}
}

// TestRegistry_ListProviders tests Registry.ListProviders
func TestRegistry_ListProviders(t *testing.T) {
	registry := NewRegistry()

	providers := registry.ListProviders()
	if len(providers) == 0 {
		t.Error("ListProviders should return non-empty list")
	}

	// Check for expected providers
	foundText := false
	foundImage := false
	foundVideo := false
	foundFallback := false

	for _, p := range providers {
		switch p {
		case "*preview.TextProvider":
			foundText = true
		case "*preview.ImageProvider":
			foundImage = true
		case "*preview.VideoProvider":
			foundVideo = true
		case "*preview.FallbackProvider":
			foundFallback = true
		}
	}

	if !foundText {
		t.Error("Should have TextProvider")
	}
	if !foundImage {
		t.Error("Should have ImageProvider")
	}
	if !foundVideo {
		t.Error("Should have VideoProvider")
	}
	if !foundFallback {
		t.Error("Should have FallbackProvider")
	}
}

// TestRegistry_GetProvider tests Registry.GetProvider for text files
func TestRegistry_GetProvider(t *testing.T) {
	registry := NewRegistry()

	textFile := models.FileEntry{
		FileType:   types.FileTypeText,
		Extension:  ".txt",
		Name:       "test.txt",
	}

	provider := registry.GetProvider(textFile)
	if provider == nil {
		t.Fatal("GetProvider returned nil for text file")
	}

	if _, ok := provider.(*TextProvider); !ok {
		t.Error("GetProvider should return TextProvider for text files")
	}
}

// TestRegistry_GetProvider_Fallback tests fallback provider
func TestRegistry_GetProvider_Fallback(t *testing.T) {
	registry := NewRegistry()

	unknownFile := models.FileEntry{
		FileType:   types.FileTypeOther,
		Extension:  ".xyz",
		Name:       "test.xyz",
	}

	provider := registry.GetProvider(unknownFile)
	if provider == nil {
		t.Fatal("GetProvider returned nil for unknown file")
	}

	if _, ok := provider.(*FallbackProvider); !ok {
		t.Error("GetProvider should return FallbackProvider for unknown files")
	}
}

// TestRegistry_GetProvider_Image tests image provider
func TestRegistry_GetProvider_Image(t *testing.T) {
	registry := NewRegistry()

	imageFile := models.FileEntry{
		FileType:   types.FileTypeImage,
		Extension:  ".jpg",
		Name:       "test.jpg",
	}

	provider := registry.GetProvider(imageFile)
	if provider == nil {
		t.Fatal("GetProvider returned nil for image file")
	}

	if _, ok := provider.(*ImageProvider); !ok {
		t.Error("GetProvider should return ImageProvider for image files")
	}
}

// TestRegistry_GetProvider_Video tests video provider
func TestRegistry_GetProvider_Video(t *testing.T) {
	registry := NewRegistry()

	videoFile := models.FileEntry{
		FileType:   types.FileTypeVideo,
		Extension:  ".mp4",
		Name:       "test.mp4",
	}

	provider := registry.GetProvider(videoFile)
	if provider == nil {
		t.Fatal("GetProvider returned nil for video file")
	}

	if _, ok := provider.(*VideoProvider); !ok {
		t.Error("GetProvider should return VideoProvider for video files")
	}
}

// TestNewImageProvider tests NewImageProvider
func TestNewImageProvider(t *testing.T) {
	provider := NewImageProvider()
	if provider == nil {
		t.Fatal("NewImageProvider returned nil")
	}
	if provider.zoomLevel != 1.0 {
		t.Errorf("zoomLevel = %f, want 1.0", provider.zoomLevel)
	}
	if !provider.enablePan {
		t.Error("enablePan should be true")
	}
	if provider.maxCacheSize != 50 {
		t.Errorf("maxCacheSize = %d, want 50", provider.maxCacheSize)
	}
	if provider.cache == nil {
		t.Error("cache should not be nil")
	}
}

// TestImageProvider_CanPreview tests ImageProvider.CanPreview
func TestImageProvider_CanPreview(t *testing.T) {
	provider := NewImageProvider()

	tests := []struct {
		fileType types.FileType
		ext      string
		want     bool
	}{
		{types.FileTypeImage, ".jpg", true},
		{types.FileTypeImage, ".png", true},
		{types.FileTypeImage, ".gif", true},
		{types.FileTypeText, ".txt", false},
		{types.FileTypeVideo, ".mp4", false},
	}

	for _, tt := range tests {
		got := provider.CanPreview(tt.fileType, tt.ext)
		if got != tt.want {
			t.Errorf("CanPreview(%v, %s) = %v, want %v", tt.fileType, tt.ext, got, tt.want)
		}
	}
}

// TestRegistry_GetProvider_UnknownType tests unknown file type
func TestRegistry_GetProvider_UnknownType(t *testing.T) {
	registry := NewRegistry()

	unknownFile := models.FileEntry{
		FileType:   types.FileType(999), // Unknown type
		Extension:  ".unknown",
		Name:       "test.unknown",
	}

	provider := registry.GetProvider(unknownFile)
	if provider == nil {
		t.Fatal("GetProvider should return a provider (fallback)")
	}
}

// TestTextProvider_Preview_NonExistent tests Preview with non-existent file
func TestTextProvider_Preview_NonExistent(t *testing.T) {
	provider := NewTextProvider()

	file := models.FileEntry{
		Path: "/nonexistent/file.txt",
	}

	_, err := provider.Preview(file)
	if err == nil {
		t.Error("Preview should return error for non-existent file")
	}
}

// TestImageProvider_CacheBehavior tests image provider cache
func TestImageProvider_CacheBehavior(t *testing.T) {
	provider := NewImageProvider()

	if len(provider.cache) != 0 {
		t.Errorf("Initial cache should be empty, got %d entries", len(provider.cache))
	}

	if provider.maxCacheSize <= 0 {
		t.Errorf("maxCacheSize should be positive, got %d", provider.maxCacheSize)
	}
}

// TestFallbackProvider_Preview tests FallbackProvider.Preview
func TestFallbackProvider_Preview(t *testing.T) {
	provider := NewFallbackProvider()

	file := models.FileEntry{
		Name: "test.txt",
		Size: 1024,
	}

	// Preview should not return error (creates widget even without file)
	_, err := provider.Preview(file)
	if err != nil {
		t.Errorf("Preview should not return error: %v", err)
	}
}

// TestRegistry_EmptyProviders tests registry with no providers
func TestRegistry_EmptyProviders(t *testing.T) {
	registry := &Registry{
		providers: []PreviewProvider{},
	}

	file := models.FileEntry{
		FileType:  types.FileTypeText,
		Extension: ".txt",
	}

	// Should return nil when no providers registered
	provider := registry.GetProvider(file)
	if provider != nil {
		t.Errorf("GetProvider should return nil when no providers registered, got %T", provider)
	}
}

// TestVideoMetadata_GetDisplayDuration tests VideoMetadata.GetDisplayDuration
func TestVideoMetadata_GetDisplayDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{"zero", 0, "Unknown"},
		{"seconds_only", 30 * time.Second, "0:30"},
		{"minutes_seconds", 2*time.Minute + 30*time.Second, "2:30"},
		{"hours_minutes_seconds", 1*time.Hour + 30*time.Minute + 45*time.Second, "1:30:45"},
		{"long_duration", 5*time.Hour + 0*time.Minute + 0*time.Second, "5:00:00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &VideoMetadata{Duration: tt.duration}
			if got := m.GetDisplayDuration(); got != tt.want {
				t.Errorf("GetDisplayDuration() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestVideoMetadata_GetDisplayResolution tests VideoMetadata.GetDisplayResolution
func TestVideoMetadata_GetDisplayResolution(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
		want   string
	}{
		{"zero", 0, 0, "Unknown"},
		{"480p", 640, 480, "640x480 (480p)"},
		{"720p", 1280, 720, "1280x720 (720p)"},
		{"1080p", 1920, 1080, "1920x1080 (1080p)"},
		{"2k", 2560, 1440, "2560x1440 (2K)"},
		{"4k", 3840, 2160, "3840x2160 (4K)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &VideoMetadata{Width: tt.width, Height: tt.height}
			if got := m.GetDisplayResolution(); got != tt.want {
				t.Errorf("GetDisplayResolution() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestVideoMetadata_GetDisplayBitrate tests VideoMetadata.GetDisplayBitrate
func TestVideoMetadata_GetDisplayBitrate(t *testing.T) {
	tests := []struct {
		name    string
		bitrate int64
		want    string
	}{
		{"zero", 0, "Unknown"},
		{"bps", 500, "500 bps"},
		{"kbps", 5000, "5.0 Kbps"},
		{"mbps", 5000000, "5.0 Mbps"},
		{"high_bitrate", 50000000, "50.0 Mbps"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &VideoMetadata{Bitrate: tt.bitrate}
			if got := m.GetDisplayBitrate(); got != tt.want {
				t.Errorf("GetDisplayBitrate() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestVideoMetadata_GetDisplayFrameRate tests VideoMetadata.GetDisplayFrameRate
func TestVideoMetadata_GetDisplayFrameRate(t *testing.T) {
	tests := []struct {
		name      string
		frameRate float64
		want      string
	}{
		{"zero", 0, "Unknown"},
		{"24fps", 24.0, "24.0 fps"},
		{"30fps", 30.0, "30.0 fps"},
		{"60fps", 60.0, "60.0 fps"},
		{"fractional", 29.97, "30.0 fps"}, // %.1f rounds 29.97 to 30.0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &VideoMetadata{FrameRate: tt.frameRate}
			if got := m.GetDisplayFrameRate(); got != tt.want {
				t.Errorf("GetDisplayFrameRate() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestFindFFmpeg tests findFFmpeg function
func TestFindFFmpeg(t *testing.T) {
	// This test checks if ffmpeg is found (may or may not be installed)
	path := findFFmpeg()
	t.Logf("findFFmpeg returned: %q", path)
	// Don't fail if not found - ffmpeg may not be installed
}

// TestFindFFprobe tests findFFprobe function
func TestFindFFprobe(t *testing.T) {
	// This test checks if ffprobe is found (may or may not be installed)
	path := findFFprobe()
	t.Logf("findFFprobe returned: %q", path)
	// Don't fail if not found - ffprobe may not be installed
}

// TestNewVideoProvider tests NewVideoProvider
func TestNewVideoProvider(t *testing.T) {
	provider := NewVideoProvider()
	if provider == nil {
		t.Fatal("NewVideoProvider returned nil")
	}
	if provider.thumbnailCache == nil {
		t.Error("thumbnailCache should not be nil")
	}
}

// TestVideoProvider_CanPreview tests VideoProvider.CanPreview
func TestVideoProvider_CanPreview(t *testing.T) {
	provider := NewVideoProvider()

	tests := []struct {
		fileType types.FileType
		ext      string
		want     bool
	}{
		{types.FileTypeVideo, ".mp4", true},
		{types.FileTypeVideo, ".avi", true},
		{types.FileTypeVideo, ".mkv", true},
		{types.FileTypeVideo, ".mov", true},
		{types.FileTypeText, ".txt", false},
		{types.FileTypeImage, ".jpg", false},
	}

	for _, tt := range tests {
		got := provider.CanPreview(tt.fileType, tt.ext)
		if got != tt.want {
			t.Errorf("CanPreview(%v, %s) = %v, want %v", tt.fileType, tt.ext, got, tt.want)
		}
	}
}

// TestVideoProvider_createPlaceholderThumbnail tests createPlaceholderThumbnail
func TestVideoProvider_createPlaceholderThumbnail(t *testing.T) {
	provider := NewVideoProvider()

	img := provider.createPlaceholderThumbnail()
	if img == nil {
		t.Fatal("createPlaceholderThumbnail returned nil")
	}

	bounds := img.Bounds()
	if bounds.Dx() != 320 {
		t.Errorf("Thumbnail width = %d, want 320", bounds.Dx())
	}
	if bounds.Dy() != 180 {
		t.Errorf("Thumbnail height = %d, want 180", bounds.Dy())
	}
}

// TestVideoProvider_OpenInPlayer tests OpenInPlayer
func TestVideoProvider_OpenInPlayer(t *testing.T) {
	provider := NewVideoProvider()

	// Test with non-existent file
	err := provider.OpenInPlayer("/nonexistent/video.mp4")
	// May or may not error depending on platform
	t.Logf("OpenInPlayer returned: %v", err)
}

// TestVideoProvider_GetMetadata tests GetMetadata
func TestVideoProvider_GetMetadata(t *testing.T) {
	provider := NewVideoProvider()

	file := models.FileEntry{
		Name:      "test.mp4",
		Size:      1024 * 1024,
		Extension: ".mp4",
		FileType:  types.FileTypeVideo,
	}

	metadata, err := provider.GetMetadata(file)
	if err != nil {
		t.Fatalf("GetMetadata failed: %v", err)
	}

	if len(metadata.Properties) == 0 {
		t.Error("Properties should not be empty")
	}
}

// TestVideoProvider_GetMetadata_NonExistent tests GetMetadata with non-existent file
func TestVideoProvider_GetMetadata_NonExistent(t *testing.T) {
	provider := NewVideoProvider()

	file := models.FileEntry{
		Name:      "nonexistent.mp4",
		Size:      0,
		Extension: ".mp4",
		FileType:  types.FileTypeVideo,
	}

	metadata, err := provider.GetMetadata(file)
	if err != nil {
		t.Logf("GetMetadata returned error (expected): %v", err)
	}
	// Should still return some metadata
	if len(metadata.Properties) == 0 {
		t.Error("Should return basic metadata even on error")
	}
}

// TestVideoProvider_ExtractThumbnail_NonExistent tests ExtractThumbnail with non-existent file
func TestVideoProvider_ExtractThumbnail_NonExistent(t *testing.T) {
	provider := NewVideoProvider()

	_, err := provider.ExtractThumbnail("/nonexistent/video.mp4")
	if err == nil {
		t.Error("ExtractThumbnail should return error for non-existent file")
	}
}

// TestVideoProvider_ExtractThumbnail_EmptyPath tests ExtractThumbnail with empty path
func TestVideoProvider_ExtractThumbnail_EmptyPath(t *testing.T) {
	provider := NewVideoProvider()

	_, err := provider.ExtractThumbnail("")
	if err == nil {
		t.Error("ExtractThumbnail should return error for empty path")
	}
}

// TestVideoMetadata_JSON tests VideoMetadata struct fields
func TestVideoMetadata_JSON(t *testing.T) {
	m := &VideoMetadata{
		Duration:     2 * time.Minute,
		Width:        1920,
		Height:       1080,
		VideoCodec:   "h264",
		AudioCodec:   "aac",
		Bitrate:      5000000,
		FrameRate:    30.0,
		HasAudio:     true,
		HasVideo:     true,
		Format:       "mp4",
		FileSize:     1024 * 1024 * 100,
		CreationTime: time.Now(),
	}

	if m.Duration != 2*time.Minute {
		t.Errorf("Duration = %v, want 2m", m.Duration)
	}
	if m.Width != 1920 {
		t.Errorf("Width = %d, want 1920", m.Width)
	}
	if m.Height != 1080 {
		t.Errorf("Height = %d, want 1080", m.Height)
	}
	if !m.HasAudio {
		t.Error("HasAudio should be true")
	}
	if !m.HasVideo {
		t.Error("HasVideo should be true")
	}
}
