package screenshot

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2"
)

// Capture takes a screenshot of a canvas
func Capture(canvas fyne.Canvas) (image.Image, error) {
	if canvas == nil {
		return nil, nil
	}
	
	// Capture the canvas
	return canvas.Capture(), nil
}

// SaveImage saves an image to a file
func SaveImage(img image.Image, path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	// Create file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	
	// Encode as PNG
	return png.Encode(file, img)
}

// LoadImage loads an image from a file
func LoadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	return png.Decode(file)
}

// CompareImages compares two images and returns the difference percentage
func CompareImages(img1, img2 image.Image) (float64, error) {
	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()
	
	// Check if sizes match
	if bounds1.Dx() != bounds2.Dx() || bounds1.Dy() != bounds2.Dy() {
		return 100.0, nil // Completely different
	}
	
	width := bounds1.Dx()
	height := bounds1.Dy()
	
	diffPixels := 0
	totalPixels := width * height
	
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			p1 := img1.At(x, y)
			p2 := img2.At(x, y)
			
			r1, g1, b1, _ := p1.RGBA()
			r2, g2, b2, _ := p2.RGBA()
			
			// Check if pixels are different
			if r1 != r2 || g1 != g2 || b1 != b2 {
				diffPixels++
			}
		}
	}
	
	return float64(diffPixels) / float64(totalPixels) * 100.0, nil
}

// AssertMatchesGolden compares a screenshot against a golden file
func AssertMatchesGolden(t *testing.T, img image.Image, goldenPath string, update bool) {
	t.Helper()
	
	// If update mode, save the new golden image
	if update {
		if err := SaveImage(img, goldenPath); err != nil {
			t.Fatalf("Failed to save golden image: %v", err)
		}
		t.Logf("Updated golden image: %s", goldenPath)
		return
	}
	
	// Load golden image
	golden, err := LoadImage(goldenPath)
	if err != nil {
		if os.IsNotExist(err) {
			t.Fatalf("Golden image not found: %s. Run with -update to create.", goldenPath)
		}
		t.Fatalf("Failed to load golden image: %v", err)
	}
	
	// Compare images
	diff, err := CompareImages(img, golden)
	if err != nil {
		t.Fatalf("Failed to compare images: %v", err)
	}
	
	// Allow 2% difference for anti-aliasing and rendering variations
	if diff > 2.0 {
		// Save diff image
		diffPath := goldenPath + ".diff.png"
		CreateDiffImage(img, golden, diffPath)
		
		t.Errorf(
			"Screenshot differs from golden image by %.2f%%\nGolden: %s\nDiff saved to: %s",
			diff, goldenPath, diffPath,
		)
	}
}

// CreateDiffImage creates a visual diff of two images
func CreateDiffImage(img1, img2 image.Image, path string) error {
	bounds := img1.Bounds()
	diff := image.NewRGBA(bounds)
	
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			p1 := img1.At(x, y)
			p2 := img2.At(x, y)
			
			r1, g1, b1, _ := p1.RGBA()
			r2, g2, b2, _ := p2.RGBA()
			
			if r1 == r2 && g1 == g2 && b1 == b2 {
				// Same pixel - show original
				diff.Set(x, y, p1)
			} else {
				// Different pixel - highlight in red
				diff.Set(x, y, red)
			}
		}
	}
	
	return SaveImage(diff, path)
}

// GoldenPath returns the path to a golden file
func GoldenPath(name string) string {
	return filepath.Join("testdata", "golden", name+".png")
}

// DiffPath returns the path to a diff file
func DiffPath(name string) string {
	return filepath.Join("testdata", "diff", name+".diff.png")
}
