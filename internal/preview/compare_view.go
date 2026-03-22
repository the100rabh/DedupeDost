package preview

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/the100rabh/DedupeDost/pkg/models"
)

// CompareView provides side-by-side file comparison
type CompareView struct {
	container     *fyne.Container
	leftCard      *CompareFileCard
	rightCard     *CompareFileCard
	syncScroll    bool
	syncZoom      bool
	registry      *Registry
}

// NewCompareView creates a new comparison view
func NewCompareView(registry *Registry) *CompareView {
	cv := &CompareView{
		registry:   registry,
		syncScroll: true,
		syncZoom:   true,
	}
	
	cv.leftCard = NewCompareFileCard()
	cv.rightCard = NewCompareFileCard()
	
	// Create split container
	split := container.NewHSplit(cv.leftCard.container, cv.rightCard.container)
	split.SetOffset(0.5)

	cv.container = container.NewVBox(split)

	return cv
}

// SetFiles sets the files to compare
func (cv *CompareView) SetFiles(left, right models.FileEntry) {
	cv.leftCard.SetFile(left)
	cv.rightCard.SetFile(right)
	
	// Load previews
	cv.loadPreviews(left, right)
}

// loadPreviews loads previews for both files
func (cv *CompareView) loadPreviews(left, right models.FileEntry) {
	// Load left preview
	if cv.registry != nil {
		leftPreview, err := cv.registry.Preview(left)
		if err == nil {
			cv.leftCard.SetPreview(leftPreview)
		}
		
		rightPreview, err := cv.registry.Preview(right)
		if err == nil {
			cv.rightCard.SetPreview(rightPreview)
		}
	}
}

// SetSyncScroll enables/disables synchronized scrolling
func (cv *CompareView) SetSyncScroll(enabled bool) {
	cv.syncScroll = enabled
}

// SetSyncZoom enables/disables synchronized zoom
func (cv *CompareView) SetSyncZoom(enabled bool) {
	cv.syncZoom = enabled
}

// CompareFileCard displays a file for comparison
type CompareFileCard struct {
	container     *fyne.Container
	pathLabel     *widget.Label
	previewArea   *fyne.Container
	sizeLabel     *widget.Label
	modifiedLabel *widget.Label
	hashLabel     *widget.Label
	keepRadio     *widget.RadioGroup
	fileEntry     models.FileEntry
	preview       fyne.CanvasObject
}

// NewCompareFileCard creates a new comparison file card
func NewCompareFileCard() *CompareFileCard {
	fc := &CompareFileCard{}
	
	// Path label
	fc.pathLabel = widget.NewLabel("")
	fc.pathLabel.TextStyle = fyne.TextStyle{Bold: true}
	fc.pathLabel.Truncation = fyne.TextTruncateEllipsis
	
	// Preview area
	fc.previewArea = container.NewStack(
		widget.NewLabel("No preview"),
	)
	
	// Size label
	fc.sizeLabel = widget.NewLabel("")
	
	// Modified label
	fc.modifiedLabel = widget.NewLabel("")
	
	// Hash label
	fc.hashLabel = widget.NewLabel("")
	fc.hashLabel.TextStyle = fyne.TextStyle{Monospace: true}
	fc.hashLabel.Truncation = fyne.TextTruncateEllipsis
	
	// Keep radio
	fc.keepRadio = widget.NewRadioGroup([]string{"Keep", "Delete"}, func(s string) {})
	fc.keepRadio.Horizontal = true
	
	// Create container
	fc.container = container.NewVBox(
		fc.pathLabel,
		fc.previewArea,
		fc.sizeLabel,
		fc.modifiedLabel,
		fc.hashLabel,
		fc.keepRadio,
	)
	
	return fc
}

// SetFile sets the file to display
func (fc *CompareFileCard) SetFile(file models.FileEntry) {
	fc.fileEntry = file
	
	// Set path
	fc.pathLabel.SetText(file.Path)
	
	// Set size
	fc.sizeLabel.SetText(fmt.Sprintf("Size: %s", file.GetDisplaySize()))
	
	// Set modified time
	fc.modifiedLabel.SetText(fmt.Sprintf("Modified: %s", 
		file.ModTime.Format("2006-01-02 15:04:05")))
	
	// Set hash (truncate for display)
	hash := file.Hash
	if len(hash) > 16 {
		hash = hash[:8] + "..." + hash[len(hash)-8:]
	}
	fc.hashLabel.SetText(fmt.Sprintf("Hash: %s", hash))
}

// SetPreview sets the preview content
func (fc *CompareFileCard) SetPreview(preview fyne.CanvasObject) {
	fc.preview = preview
	fc.previewArea.Objects = []fyne.CanvasObject{preview}
	fc.previewArea.Refresh()
}

// SetKeep sets the keep/delete state
func (fc *CompareFileCard) SetKeep(keep bool) {
	if keep {
		fc.keepRadio.SetSelected("Keep")
	} else {
		fc.keepRadio.SetSelected("Delete")
	}
}

// GetKeep returns whether the file is marked to keep
func (fc *CompareFileCard) GetKeep() bool {
	return fc.keepRadio.Selected == "Keep"
}

// Clear clears the card
func (fc *CompareFileCard) Clear() {
	fc.pathLabel.SetText("")
	fc.sizeLabel.SetText("")
	fc.modifiedLabel.SetText("")
	fc.hashLabel.SetText("")
	fc.keepRadio.SetSelected("")
	fc.previewArea.Objects = []fyne.CanvasObject{widget.NewLabel("No preview")}
	fc.previewArea.Refresh()
}
