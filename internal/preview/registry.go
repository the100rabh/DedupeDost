package preview

import (
	"fmt"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/dupdel/dup-del/pkg/models"
	"github.com/dupdel/dup-del/pkg/types"
)

// PreviewProvider defines the interface for file preview providers
type PreviewProvider interface {
	// CanPreview checks if this provider can preview the given file type
	CanPreview(fileType types.FileType, extension string) bool
	
	// Preview creates a preview widget for the file
	Preview(file models.FileEntry) (fyne.CanvasObject, error)
	
	// GetMetadata returns metadata about the file
	GetMetadata(file models.FileEntry) (Metadata, error)
}

// Registry manages preview providers
type Registry struct {
	providers []PreviewProvider
	mu        sync.RWMutex
}

// NewRegistry creates a new preview provider registry
func NewRegistry() *Registry {
	r := &Registry{
		providers: make([]PreviewProvider, 0),
	}
	
	// Register default providers
	r.Register(NewTextProvider())
	r.Register(NewImageProvider())
	r.Register(NewVideoProvider())
	r.Register(NewFallbackProvider())
	
	return r
}

// Register registers a preview provider
func (r *Registry) Register(provider PreviewProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers = append(r.providers, provider)
}

// GetProvider returns the best provider for the given file
func (r *Registry) GetProvider(file models.FileEntry) PreviewProvider {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	for _, provider := range r.providers {
		if provider.CanPreview(file.FileType, file.Extension) {
			return provider
		}
	}
	
	// Fallback should always be available
	return r.providers[len(r.providers)-1]
}

// Preview creates a preview for the given file
func (r *Registry) Preview(file models.FileEntry) (fyne.CanvasObject, error) {
	provider := r.GetProvider(file)
	return provider.Preview(file)
}

// GetMetadata returns metadata for the given file
func (r *Registry) GetMetadata(file models.FileEntry) (Metadata, error) {
	provider := r.GetProvider(file)
	return provider.GetMetadata(file)
}

// ListProviders returns all registered providers
func (r *Registry) ListProviders() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	names := make([]string, len(r.providers))
	for i, p := range r.providers {
		names[i] = fmt.Sprintf("%T", p)
	}
	return names
}

// FallbackProvider provides a fallback preview for unsupported files
type FallbackProvider struct{}

// NewFallbackProvider creates a new fallback provider
func NewFallbackProvider() *FallbackProvider {
	return &FallbackProvider{}
}

// CanPreview always returns true for fallback
func (p *FallbackProvider) CanPreview(fileType types.FileType, extension string) bool {
	return true
}

// Preview creates a fallback preview widget
func (p *FallbackProvider) Preview(file models.FileEntry) (fyne.CanvasObject, error) {
	content := container.NewCenter(
		container.NewVBox(
			widget.NewIcon(theme.FileIcon()),
			widget.NewLabel(file.Name),
			widget.NewLabel(file.GetDisplaySize()),
			widget.NewLabel("No preview available"),
		),
	)
	
	return content, nil
}

// GetMetadata returns basic metadata
func (p *FallbackProvider) GetMetadata(file models.FileEntry) (Metadata, error) {
	return Metadata{
		Properties: map[string]string{
			"Name":      file.Name,
			"Size":      file.GetDisplaySize(),
			"Type":      file.Extension,
			"Path":      file.Path,
			"Modified":  file.ModTime.Format("2006-01-02 15:04:05"),
		},
	}, nil
}
