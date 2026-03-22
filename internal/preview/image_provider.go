package preview

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/bmp"
	"golang.org/x/image/webp"

	"github.com/the100rabh/DedupeDost/pkg/models"
	"github.com/the100rabh/DedupeDost/pkg/types"
)

// ImageProvider provides preview for image files
type ImageProvider struct {
	zoomLevel    float32
	enablePan    bool
	cache        map[string]image.Image
	maxCacheSize int
}

// NewImageProvider creates a new image preview provider
func NewImageProvider() *ImageProvider {
	return &ImageProvider{
		zoomLevel:    1.0,
		enablePan:    true,
		cache:        make(map[string]image.Image),
		maxCacheSize: 50,
	}
}

// CanPreview checks if this provider can preview the file
func (p *ImageProvider) CanPreview(fileType types.FileType, extension string) bool {
	return fileType == types.FileTypeImage
}

// Preview creates a preview widget for the file
func (p *ImageProvider) Preview(file models.FileEntry) (fyne.CanvasObject, error) {
	img, err := p.loadImage(file.Path)
	if err != nil {
		return p.createErrorWidget(err), nil
	}

	// Convert image.Image to Fyne image
	imgData := p.imageToImage(img)

	// Get image dimensions
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Create info label
	infoLabel := widget.NewLabel(fmt.Sprintf("%dx%d pixels • %s", width, height, file.GetDisplaySize()))

	// Create zoom controls
	zoomSlider := widget.NewSlider(0.25, 4.0)
	zoomSlider.Value = 1.0
	zoomLabel := widget.NewLabel("Zoom: 100%")

	zoomSlider.OnChanged = func(value float64) {
		p.zoomLevel = float32(value)
		zoomLabel.SetText(fmt.Sprintf("Zoom: %.0f%%", value*100))
		imgData.Refresh()
	}

	// Create main image container with scroll
	imageContainer := container.NewScroll(imgData)

	// Create controls
	controls := container.NewHBox(
		widget.NewLabel("Zoom:"),
		zoomSlider,
		zoomLabel,
		layout.NewSpacer(),
		widget.NewButtonWithIcon("Fit", theme.ZoomFitIcon(), func() {
			// Calculate fit zoom
			scrollSize := imageContainer.Size()
			if scrollSize.Width > 0 && scrollSize.Height > 0 {
				zoomW := float32(scrollSize.Width) / float32(width)
				zoomH := float32(scrollSize.Height) / float32(height)
				zoom := zoomW
				if zoomH < zoomW {
					zoom = zoomH
				}
				p.zoomLevel = zoom
				zoomSlider.Value = float64(zoom)
				zoomLabel.SetText(fmt.Sprintf("Zoom: %.0f%%", zoom*100))
			}
		}),
		widget.NewButtonWithIcon("Reset", theme.ViewRefreshIcon(), func() {
			p.zoomLevel = 1.0
			zoomSlider.Value = 1.0
			zoomLabel.SetText("Zoom: 100%")
		}),
	)

	// Create main container
	content := container.NewBorder(
		container.NewVBox(infoLabel, controls),
		nil,
		nil,
		nil,
		imageContainer,
	)

	return content, nil
}

// imageToImage converts image.Image to Fyne Image
func (p *ImageProvider) imageToImage(img image.Image) *canvas.Image {
	return canvas.NewImageFromImage(img)
}

// loadImage loads an image from file with caching
func (p *ImageProvider) loadImage(path string) (image.Image, error) {
	// Check cache
	if img, ok := p.cache[path]; ok {
		return img, nil
	}
	
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	var img image.Image
	ext := strings.ToLower(filepath.Ext(path))
	
	switch ext {
	case ".bmp":
		img, err = bmp.Decode(file)
	case ".webp":
		img, err = webp.Decode(file)
	default:
		img, _, err = image.Decode(file)
	}
	
	if err != nil {
		return nil, err
	}
	
	// Cache the image
	if len(p.cache) >= p.maxCacheSize {
		// Simple cache eviction - clear half the cache
		count := 0
		for k := range p.cache {
			delete(p.cache, k)
			count++
			if count >= p.maxCacheSize/2 {
				break
			}
		}
	}
	p.cache[path] = img
	
	return img, nil
}

// createErrorWidget creates an error display widget
func (p *ImageProvider) createErrorWidget(err error) fyne.CanvasObject {
	return container.NewCenter(
		container.NewVBox(
			widget.NewIcon(theme.BrokenImageIcon()),
			widget.NewLabel("Error loading image"),
			widget.NewLabel(err.Error()),
		),
	)
}

// GetMetadata returns metadata about the image file
func (p *ImageProvider) GetMetadata(file models.FileEntry) (Metadata, error) {
	img, err := p.loadImage(file.Path)
	if err != nil {
		return Metadata{}, err
	}
	
	bounds := img.Bounds()
	
	// Try to get EXIF data
	exifData := p.getExifData(file.Path)
	
	properties := map[string]string{
		"Width":     fmt.Sprintf("%d px", bounds.Dx()),
		"Height":    fmt.Sprintf("%d px", bounds.Dy()),
		"Size":      file.GetDisplaySize(),
		"Format":    strings.TrimPrefix(filepath.Ext(file.Path), "."),
	}
	
	// Add EXIF data if available
	for k, v := range exifData {
		properties[k] = v
	}
	
	return Metadata{
		Properties: properties,
	}, nil
}

// getExifData extracts EXIF metadata from image
func (p *ImageProvider) getExifData(path string) map[string]string {
	data := make(map[string]string)
	
	file, err := os.Open(path)
	if err != nil {
		return data
	}
	defer file.Close()
	
	// Try to read EXIF using standard library
	// Note: For full EXIF support, use github.com/rwcarlsen/goexif
	_, _, err = image.DecodeConfig(file)
	if err != nil {
		return data
	}
	
	// Get file info for date
	info, err := os.Stat(path)
	if err == nil {
		data["Modified"] = info.ModTime().Format("2006-01-02 15:04")
	}

	return data
}
