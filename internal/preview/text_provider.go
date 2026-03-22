package preview

import (
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/dupdel/dup-del/pkg/models"
	"github.com/dupdel/dup-del/pkg/types"
)

// TextProvider provides preview for text files
type TextProvider struct {
	wordWrap bool
}

// NewTextProvider creates a new text preview provider
func NewTextProvider() *TextProvider {
	return &TextProvider{
		wordWrap: true,
	}
}

// CanPreview checks if this provider can preview the file
func (p *TextProvider) CanPreview(fileType types.FileType, extension string) bool {
	return fileType == types.FileTypeText
}

// Preview creates a preview widget for the file
func (p *TextProvider) Preview(file models.FileEntry) (fyne.CanvasObject, error) {
	content, err := os.ReadFile(file.Path)
	if err != nil {
		return nil, err
	}

	text := string(content)

	// Create simple label with monospace font
	textLabel := widget.NewLabel(text)
	textLabel.TextStyle = fyne.TextStyle{Monospace: true}

	if p.wordWrap {
		textLabel.Wrapping = fyne.TextWrapWord
	}

	// Wrap in scroll for overflow
	scroll := container.NewScroll(textLabel)

	return scroll, nil
}

// GetMetadata returns metadata about the text file
func (p *TextProvider) GetMetadata(file models.FileEntry) (Metadata, error) {
	content, err := os.ReadFile(file.Path)
	if err != nil {
		return Metadata{}, err
	}

	lines := strings.Count(string(content), "\n") + 1
	words := len(strings.Fields(string(content)))

	return Metadata{
		Properties: map[string]string{
			"Lines": fmt.Sprintf("%d", lines),
			"Words": fmt.Sprintf("%d", words),
			"Chars": fmt.Sprintf("%d", len(content)),
		},
	}, nil
}

// Metadata contains preview metadata
type Metadata struct {
	Properties    map[string]string
	VideoMetadata interface{} // *VideoMetadata for video files
}
