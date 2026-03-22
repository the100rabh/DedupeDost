package preview_test

import (
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/image/bmp"

	"github.com/the100rabh/DedupeDost/internal/preview"
	"github.com/the100rabh/DedupeDost/pkg/models"
	"github.com/the100rabh/DedupeDost/pkg/types"
)

// Helper function to create test directory
func setupTestDir(t *testing.T) (string, func()) {
	t.Helper()

	dir, err := os.MkdirTemp("", "preview_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(dir)
	}

	return dir, cleanup
}

// Helper function to create a test image
func createTestImage(t *testing.T, dir, name string, format string) string {
	t.Helper()

	// Create a simple test image (100x100 pixels with red color)
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	path := filepath.Join(dir, name)
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	switch format {
	case "png":
		err = png.Encode(file, img)
	case "jpeg", "jpg":
		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 80})
	case "gif":
		palette := []color.Color{color.RGBA{255, 0, 0, 255}}
		err = gif.Encode(file, img, &gif.Options{NumColors: len(palette)})
	case "bmp":
		err = bmp.Encode(file, img)
	case "webp":
		// For WebP, create a simple valid WebP file (RIFF header + VP8)
		// Using a minimal valid WebP file structure
		webpData := []byte{
			// RIFF header
			0x52, 0x49, 0x46, 0x46, // "RIFF"
			0x1a, 0x00, 0x00, 0x00, // File size - 8
			0x57, 0x45, 0x42, 0x50, // "WEBP"
			// VP8 chunk
			0x56, 0x50, 0x38, 0x20, // "VP8 "
			0x0e, 0x00, 0x00, 0x00, // Chunk size
			0x30, 0x01, 0x00, 0x9d, 0x01, 0x2a, 0x01, 0x00,
			0x01, 0x00, 0xfe, 0xfb, 0x94, 0x00, 0x00,
		}
		_, err = file.Write(webpData)
	default:
		t.Fatalf("Unknown format: %s", format)
	}

	if err != nil {
		t.Fatalf("Failed to encode image: %v", err)
	}

	return path
}

// TestImageProvider_CanPreview tests format detection
func TestImageProvider_CanPreview(t *testing.T) {
	p := preview.NewImageProvider()

	tests := []struct {
		name      string
		fileType  types.FileType
		extension string
		expected  bool
	}{
		{"PNG file", types.FileTypeImage, ".png", true},
		{"JPEG file", types.FileTypeImage, ".jpg", true},
		{"GIF file", types.FileTypeImage, ".gif", true},
		{"BMP file", types.FileTypeImage, ".bmp", true},
		{"WebP file", types.FileTypeImage, ".webp", true},
		{"Text file", types.FileTypeText, ".txt", false},
		{"Unknown file", types.FileTypeUnknown, ".xyz", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := p.CanPreview(tt.fileType, tt.extension)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestImageProvider_Preview_PNG tests PNG preview
func TestImageProvider_Preview_PNG(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	path := createTestImage(t, dir, "test.png", "png")
	file, _ := models.NewFileEntry(path)

	p := preview.NewImageProvider()
	result, err := p.Preview(*file)

	if err != nil {
		t.Fatalf("Preview failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil preview widget")
	}
}

// TestImageProvider_Preview_JPEG tests JPEG preview
func TestImageProvider_Preview_JPEG(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	path := createTestImage(t, dir, "test.jpg", "jpeg")
	file, _ := models.NewFileEntry(path)

	p := preview.NewImageProvider()
	result, err := p.Preview(*file)

	if err != nil {
		t.Fatalf("Preview failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil preview widget")
	}
}

// TestImageProvider_Preview_GIF tests GIF preview
func TestImageProvider_Preview_GIF(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	path := createTestImage(t, dir, "test.gif", "gif")
	file, _ := models.NewFileEntry(path)

	p := preview.NewImageProvider()
	result, err := p.Preview(*file)

	if err != nil {
		t.Fatalf("Preview failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil preview widget")
	}
}

// TestImageProvider_Preview_BMP tests BMP preview
func TestImageProvider_Preview_BMP(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	path := createTestImage(t, dir, "test.bmp", "bmp")
	file, _ := models.NewFileEntry(path)

	p := preview.NewImageProvider()
	result, err := p.Preview(*file)

	if err != nil {
		t.Fatalf("Preview failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil preview widget")
	}
}

// TestImageProvider_Preview_WebP tests WebP preview
func TestImageProvider_Preview_WebP(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	path := createTestImage(t, dir, "test.webp", "webp")
	file, _ := models.NewFileEntry(path)

	p := preview.NewImageProvider()
	result, err := p.Preview(*file)

	if err != nil {
		t.Fatalf("Preview failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil preview widget")
	}
}

// TestImageProvider_Preview_Error tests error handling for non-existent file
func TestImageProvider_Preview_Error(t *testing.T) {
	p := preview.NewImageProvider()

	file := models.FileEntry{
		Path: "/nonexistent/file.png",
		Name: "missing.png",
	}

	result, err := p.Preview(file)

	if err != nil {
		t.Fatalf("Preview should not return error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected error widget for missing file")
	}
}

// TestImageProvider_Preview_CorruptFile tests error handling for corrupt images
func TestImageProvider_Preview_CorruptFile(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	// Create a file with invalid image data
	path := filepath.Join(dir, "corrupt.png")
	err := os.WriteFile(path, []byte("not a valid image"), 0644)
	if err != nil {
		t.Fatalf("Failed to create corrupt file: %v", err)
	}

	file, _ := models.NewFileEntry(path)

	p := preview.NewImageProvider()
	result, err := p.Preview(*file)

	if err != nil {
		t.Fatalf("Preview should not return error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected error widget for corrupt file")
	}
}

// TestImageProvider_GetMetadata tests metadata extraction
func TestImageProvider_GetMetadata(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	path := createTestImage(t, dir, "test.png", "png")
	file, _ := models.NewFileEntry(path)

	p := preview.NewImageProvider()
	metadata, err := p.GetMetadata(*file)

	if err != nil {
		t.Fatalf("GetMetadata failed: %v", err)
	}

	props := metadata.Properties

	// Check width
	if props["Width"] != "100 px" {
		t.Errorf("Expected width 100 px, got %s", props["Width"])
	}

	// Check height
	if props["Height"] != "100 px" {
		t.Errorf("Expected height 100 px, got %s", props["Height"])
	}

	// Check format
	if props["Format"] != "png" {
		t.Errorf("Expected format png, got %s", props["Format"])
	}

	// Check size exists
	if props["Size"] == "" {
		t.Error("Expected size to be set")
	}
}

// TestImageProvider_Cache tests image caching
func TestImageProvider_Cache(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	path := createTestImage(t, dir, "test.png", "png")
	file, _ := models.NewFileEntry(path)

	p := preview.NewImageProvider()

	// First load - should cache
	_, err := p.Preview(*file)
	if err != nil {
		t.Fatalf("First preview failed: %v", err)
	}

	// Second load - should use cache
	_, err = p.Preview(*file)
	if err != nil {
		t.Fatalf("Cached preview failed: %v", err)
	}
}

// TestImageProvider_DifferentSizes tests images of different sizes
func TestImageProvider_DifferentSizes(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	// Create small image
	smallImg := image.NewRGBA(image.Rect(0, 0, 10, 10))
	smallPath := filepath.Join(dir, "small.png")
	smallFile, _ := os.Create(smallPath)
	png.Encode(smallFile, smallImg)
	smallFile.Close()

	// Create large image
	largeImg := image.NewRGBA(image.Rect(0, 0, 1000, 1000))
	largePath := filepath.Join(dir, "large.png")
	largeFile, _ := os.Create(largePath)
	png.Encode(largeFile, largeImg)
	largeFile.Close()

	p := preview.NewImageProvider()

	// Preview small image
	smallEntry, _ := models.NewFileEntry(smallPath)
	smallResult, err := p.Preview(*smallEntry)
	if err != nil {
		t.Fatalf("Small image preview failed: %v", err)
	}
	if smallResult == nil {
		t.Error("Expected non-nil preview for small image")
	}

	// Preview large image
	largeEntry, _ := models.NewFileEntry(largePath)
	largeResult, err := p.Preview(*largeEntry)
	if err != nil {
		t.Fatalf("Large image preview failed: %v", err)
	}
	if largeResult == nil {
		t.Error("Expected non-nil preview for large image")
	}
}
