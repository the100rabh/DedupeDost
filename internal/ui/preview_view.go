package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/the100rabh/DedupeDost/internal/preview"
	"github.com/the100rabh/DedupeDost/pkg/models"
)

// PreviewView represents the file preview view
type PreviewView struct {
	container     *fyne.Container
	groupLabel    *widget.Label
	prevButton    *widget.Button
	nextButton    *widget.Button
	fileCard1     *FilePreviewCard
	fileCard2     *FilePreviewCard
	fileIndex     int        // Current file index being compared (0 = first file)
	actionsBox    *fyne.Container
	applyToAll    *widget.Check
	saveBtn       *widget.Button
	skipBtn       *widget.Button
	currentGroup  *models.DuplicateGroup // Pointer to group in filteredGroups
	groupIndex    int
	totalGroups   int
	dupDelApp     *DedupeDostApp
	registry      *preview.Registry
}

// NewPreviewView creates a new preview view
func NewPreviewView(dupDelApp *DedupeDostApp) *PreviewView {
	pv := &PreviewView{
		dupDelApp:   dupDelApp,
		groupIndex:  -1,
		totalGroups: 0,
		registry:    preview.NewRegistry(),
	}
	
	// Group label
	pv.groupLabel = widget.NewLabel("No group selected")
	pv.groupLabel.TextStyle = fyne.TextStyle{Bold: true}
	
	// Previous button
	pv.prevButton = widget.NewButtonWithIcon("Previous", theme.NavigateBackIcon(), func() {
		pv.previousGroup()
	})
	pv.prevButton.Disable()
	
	// Next button
	pv.nextButton = widget.NewButtonWithIcon("Next", theme.NavigateNextIcon(), func() {
		pv.nextGroup()
	})
	pv.nextButton.Disable()

	// File cards
	pv.fileCard1 = NewFilePreviewCard()
	pv.fileCard2 = NewFilePreviewCard()

	// File navigation within group
	prevFileBtn := widget.NewButtonWithIcon("Previous File", theme.NavigateBackIcon(), func() {
		pv.previousFile()
	})
	nextFileBtn := widget.NewButtonWithIcon("Next File", theme.NavigateNextIcon(), func() {
		pv.nextFile()
	})
	pv.fileIndex = 0

	// File index label
	fileIndexLabel := widget.NewLabel("")
	fileIndexLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Apply to all check
	pv.applyToAll = widget.NewCheck("Apply to all remaining groups", nil)

	// Save button
	pv.saveBtn = widget.NewButtonWithIcon("Save & Continue", theme.ConfirmIcon(), func() {
		pv.saveDecision()
	})
	pv.saveBtn.Importance = widget.HighImportance

	// Skip button
	pv.skipBtn = widget.NewButtonWithIcon("Skip", theme.ContentClearIcon(), func() {
		pv.skipGroup()
	})

	// Navigation box
	navBox := container.NewHBox(
		pv.prevButton,
		pv.nextButton,
		layout.NewSpacer(),
		pv.groupLabel,
	)

	// File navigation box
	fileNavBox := container.NewHBox(
		prevFileBtn,
		nextFileBtn,
		layout.NewSpacer(),
		fileIndexLabel,
	)

	// Split preview - expands to fill available space
	splitPreview := container.NewHSplit(
		pv.fileCard1.container,
		pv.fileCard2.container,
	)
	splitPreview.SetOffset(0.5)

	// Actions box with dynamic buttons for multiple files
	pv.actionsBox = container.NewVBox(
		widget.NewSeparator(),
		fileNavBox,
		widget.NewSeparator(),
	)

	// Create container - NO outer scroll, let split expand naturally
	content := container.NewBorder(
		navBox,
		pv.actionsBox,
		nil,
		nil,
		splitPreview,
	)

	pv.container = content

	return pv
}

// Container returns the container
func (pv *PreviewView) Container() *fyne.Container {
	return pv.container
}

// setGroup sets the group to preview
func (pv *PreviewView) setGroup(group *models.DuplicateGroup) {
	pv.currentGroup = group
	pv.fileIndex = 0

	// Get total groups from results view
	if pv.dupDelApp != nil && pv.dupDelApp.mainView != nil {
		pv.totalGroups = len(pv.dupDelApp.mainView.resultsView.filteredGroups)
		pv.groupIndex = pv.dupDelApp.mainView.resultsView.selectedGroup
	}

	// Truncate filename for display if too long
	filename := group.Files[0].Name
	if len(filename) > 30 {
		filename = filename[:27] + "..."
	}
	pv.groupLabel.SetText(fmt.Sprintf("Group %d of %d: %s",
		pv.groupIndex+1, pv.totalGroups, filename))
	pv.groupLabel.Truncation = fyne.TextTruncateEllipsis

	// Set files based on current index
	pv.updateFileCards()

	// Update buttons
	pv.prevButton.Enable()
	pv.nextButton.Enable()
	if pv.groupIndex <= 0 {
		pv.prevButton.Disable()
	}
	if pv.groupIndex >= pv.totalGroups-1 {
		pv.nextButton.Disable()
	}

	// Refresh containers
	pv.fileCard1.container.Refresh()
	pv.fileCard2.container.Refresh()
	pv.updateActionButtons()
}

// updateFileCards updates the file cards based on current fileIndex
func (pv *PreviewView) updateFileCards() {
	if len(pv.currentGroup.Files) == 0 {
		pv.fileCard1.Clear()
		pv.fileCard2.Clear()
		return
	}

	// Always show first file on left
	pv.fileCard1.SetFile(pv.currentGroup.Files[0])
	pv.fileCard1.SetKeep(isFileKept(0, pv.currentGroup))

	// Show next file on right (or clear if only one file)
	if len(pv.currentGroup.Files) > 1 {
		rightIndex := pv.fileIndex + 1
		if rightIndex >= len(pv.currentGroup.Files) {
			pv.fileIndex = 0
			rightIndex = 1
		}
		pv.fileCard2.SetFile(pv.currentGroup.Files[rightIndex])
		pv.fileCard2.SetKeep(isFileKept(rightIndex, pv.currentGroup))
	} else {
		pv.fileCard2.Clear()
	}

	pv.updateActionButtons()
}

// updateActionButtons updates the action buttons based on current state
func (pv *PreviewView) updateActionButtons() {
	// Clear existing action buttons
	pv.actionsBox.Objects = []fyne.CanvasObject{
		widget.NewSeparator(),
	}

	if len(pv.currentGroup.Files) <= 1 {
		// Only one file, show simple keep button
		keepBtn := widget.NewButtonWithIcon("Keep This File", theme.ConfirmIcon(), func() {
			pv.keepFile(0)
		})
		pv.actionsBox.Add(container.NewHBox(keepBtn, layout.NewSpacer()))
	} else {
		// Multiple files - show comparison navigation and keep buttons
		fileNavBox := container.NewHBox(
			widget.NewButtonWithIcon("Previous File", theme.NavigateBackIcon(), func() {
				pv.previousFile()
			}),
			widget.NewLabel(fmt.Sprintf("Comparing: %d of %d", pv.fileIndex+2, len(pv.currentGroup.Files))),
			widget.NewButtonWithIcon("Next File", theme.NavigateNextIcon(), func() {
				pv.nextFile()
			}),
			layout.NewSpacer(),
		)
		pv.actionsBox.Add(fileNavBox)
		pv.actionsBox.Add(widget.NewSeparator())

		// Keep buttons
		rightIndex := pv.fileIndex + 1
		if rightIndex < len(pv.currentGroup.Files) {
			keepLeftBtn := widget.NewButtonWithIcon("Keep Left", theme.ConfirmIcon(), func() {
				pv.keepFile(0)
			})
			keepRightBtn := widget.NewButtonWithIcon("Keep Right", theme.ConfirmIcon(), func() {
				pv.keepFile(rightIndex)
			})
			keepBothBtn := widget.NewButtonWithIcon("Keep Both", theme.ContentAddIcon(), func() {
				pv.keepBothFiles(0, rightIndex)
			})
			pv.actionsBox.Add(container.NewHBox(
				keepLeftBtn,
				keepRightBtn,
				keepBothBtn,
				layout.NewSpacer(),
			))
		}

		pv.actionsBox.Add(pv.applyToAll)
		pv.actionsBox.Add(widget.NewSeparator())
		pv.actionsBox.Add(container.NewHBox(pv.saveBtn, pv.skipBtn))
	}

	pv.actionsBox.Refresh()
}

// isFileKept checks if a file at the given index is kept
func isFileKept(index int, group *models.DuplicateGroup) bool {
	for _, delPath := range group.DeletePaths {
		// Check if this index is in delete paths
		if index < len(group.Files) && group.Files[index].Path == delPath {
			return false
		}
	}
	return true
}

// previousFile navigates to the previous file in the group
func (pv *PreviewView) previousFile() {
	if pv.fileIndex > 0 {
		pv.fileIndex--
		pv.updateFileCards()
	}
}

// nextFile navigates to the next file in the group
func (pv *PreviewView) nextFile() {
	if pv.fileIndex+1 < len(pv.currentGroup.Files)-1 {
		pv.fileIndex++
		pv.updateFileCards()
	}
}

// keepFile marks a file to keep
func (pv *PreviewView) keepFile(index int) {
	if pv.currentGroup.Files == nil || index >= len(pv.currentGroup.Files) {
		return
	}

	pv.currentGroup.MarkForDeletion([]int{index})
	pv.fileCard1.SetKeep(isFileKept(0, pv.currentGroup))
	if pv.fileIndex+1 < len(pv.currentGroup.Files) {
		pv.fileCard2.SetKeep(isFileKept(pv.fileIndex+1, pv.currentGroup))
	}

	// Refresh results view to show updated decisions
	pv.refreshResultsView()
}

// keepBothFiles marks both files to keep
func (pv *PreviewView) keepBothFiles(index1, index2 int) {
	pv.currentGroup.ClearDecisions()
	pv.fileCard1.SetKeep(true)
	pv.fileCard2.SetKeep(true)

	// Refresh results view to show updated decisions
	pv.refreshResultsView()
}

// keepBoth marks both files to keep (alias for keepBothFiles)
func (pv *PreviewView) keepBoth() {
	pv.keepBothFiles(0, 1)
}

// refreshResultsView refreshes the results view to show updated decisions
func (pv *PreviewView) refreshResultsView() {
	if pv.dupDelApp != nil && pv.dupDelApp.mainView != nil {
		pv.dupDelApp.mainView.resultsView.refreshCards()
	}
}

// SetGroups sets all groups for navigation
func (pv *PreviewView) SetGroups(groups []*models.DuplicateGroup) {
	pv.totalGroups = len(groups)
	if len(groups) > 0 {
		pv.setGroup(groups[0])
		pv.groupIndex = 0
	}
}

// previousGroup navigates to the previous group
func (pv *PreviewView) previousGroup() {
	if pv.groupIndex > 0 {
		pv.groupIndex--
		if pv.dupDelApp != nil && pv.dupDelApp.mainView != nil {
			groups := pv.dupDelApp.mainView.resultsView.filteredGroups
			if pv.groupIndex < len(groups) {
				pv.setGroup(groups[pv.groupIndex])
			}
		}
	}
}

// nextGroup navigates to the next group
func (pv *PreviewView) nextGroup() {
	if pv.groupIndex < pv.totalGroups-1 {
		pv.groupIndex++
		if pv.dupDelApp != nil && pv.dupDelApp.mainView != nil {
			groups := pv.dupDelApp.mainView.resultsView.filteredGroups
			if pv.groupIndex < len(groups) {
				pv.setGroup(groups[pv.groupIndex])
			}
		}
	}
}

// saveDecision saves the current decision
func (pv *PreviewView) saveDecision() {
	// Apply to all remaining groups if checked
	if pv.applyToAll.Checked && pv.currentGroup != nil && len(pv.currentGroup.KeepIndices) > 0 {
		pv.applyDecisionToAllRemainingGroups()
	}
	pv.nextGroup()
}

// applyDecisionToAllRemainingGroups applies the current group's decision to all remaining groups
func (pv *PreviewView) applyDecisionToAllRemainingGroups() {
	if pv.dupDelApp == nil || pv.dupDelApp.mainView == nil {
		return
	}

	groups := pv.dupDelApp.mainView.resultsView.filteredGroups
	keepIndices := pv.currentGroup.KeepIndices

	// Apply the same keep indices to all remaining groups
	for i := pv.groupIndex + 1; i < len(groups); i++ {
		// Only apply if the group has enough files
		if len(groups[i].Files) > 0 && len(keepIndices) > 0 {
			// Validate that keep indices are valid for this group
			validIndices := make([]int, 0)
			for _, idx := range keepIndices {
				if idx < len(groups[i].Files) {
					validIndices = append(validIndices, idx)
				}
			}
			// If we have valid indices, apply the decision
			if len(validIndices) > 0 {
				groups[i].MarkForDeletion(validIndices)
			}
		}
	}

	// Refresh the results view to show updated decisions
	pv.dupDelApp.mainView.resultsView.refreshCards()
}

// skipGroup skips the current group
func (pv *PreviewView) skipGroup() {
	pv.nextGroup()
}

// FilePreviewCard represents a card showing a file preview
type FilePreviewCard struct {
	container       *fyne.Container
	previewScroll   *container.Scroll
	previewContent  fyne.CanvasObject
	pathLabel       *widget.Label
	previewArea     fyne.CanvasObject
	sizeLabel       *widget.Label
	modifiedLabel   *widget.Label
	keepRadio       *widget.RadioGroup
	fileEntry       models.FileEntry
	keep            bool
}

// NewFilePreviewCard creates a new file preview card
func NewFilePreviewCard() *FilePreviewCard {
	fc := &FilePreviewCard{}

	// Path label
	fc.pathLabel = widget.NewLabel("")
	fc.pathLabel.TextStyle = fyne.TextStyle{Bold: true}
	fc.pathLabel.Truncation = fyne.TextTruncateEllipsis

	// Preview scroll - starts with placeholder
	placeholder := widget.NewLabel("Select a file to preview")
	fc.previewScroll = container.NewScroll(placeholder)
	fc.previewContent = fc.previewScroll

	// Size label
	fc.sizeLabel = widget.NewLabel("")

	// Modified label
	fc.modifiedLabel = widget.NewLabel("")

	// Keep radio
	fc.keepRadio = widget.NewRadioGroup([]string{"Keep", "Delete"}, func(s string) {
		fc.keep = s == "Keep"
	})
	fc.keepRadio.Horizontal = true

	// Border layout: path top, info bottom, preview center (expands)
	fc.container = container.NewBorder(
		fc.pathLabel,
		container.NewVBox(
			fc.sizeLabel,
			fc.modifiedLabel,
			fc.keepRadio,
		),
		nil,
		nil,
		fc.previewScroll,
	)

	return fc
}

// Container returns the container
func (fc *FilePreviewCard) Container() *fyne.Container {
	return fc.container
}

// SetFile sets the file to display
func (fc *FilePreviewCard) SetFile(file models.FileEntry) {
	fc.fileEntry = file

	// Set path with truncation for long paths
	fc.pathLabel.SetText(file.Path)
	fc.pathLabel.Truncation = fyne.TextTruncateEllipsis

	// Set size
	fc.sizeLabel.SetText(fmt.Sprintf("Size: %s", formatSize(file.Size)))

	// Set modified time
	fc.modifiedLabel.SetText(fmt.Sprintf("Modified: %s",
		file.ModTime.Format("2006-01-02 15:04:05")))
	fc.modifiedLabel.Truncation = fyne.TextTruncateEllipsis

	// Set preview based on file type
	fc.setPreview(file)
}

// setPreview sets the preview content based on file type
func (fc *FilePreviewCard) setPreview(file models.FileEntry) {
	var previewObj fyne.CanvasObject
	var err error

	// Use registry to get appropriate provider
	reg := preview.NewRegistry()
	provider := reg.GetProvider(file)

	previewObj, err = provider.Preview(file)
	if err != nil {
		previewObj = fc.createErrorPreview(err.Error())
	}

	// Update scroll content
	if fc.previewScroll != nil {
		fc.previewScroll.Content = previewObj
		fc.previewScroll.Refresh()
		fc.container.Refresh()
	}
}

// createErrorPreview creates an error preview widget
func (fc *FilePreviewCard) createErrorPreview(errMsg string) fyne.CanvasObject {
	return container.NewCenter(
		container.NewVBox(
			widget.NewIcon(theme.BrokenImageIcon()),
			widget.NewLabel("Error loading preview"),
			widget.NewLabel(errMsg),
		),
	)
}

// Clear clears the card
func (fc *FilePreviewCard) Clear() {
	fc.pathLabel.SetText("")
	fc.sizeLabel.SetText("")
	fc.modifiedLabel.SetText("")
	fc.keepRadio.SetSelected("")
	
	// Reset to placeholder
	if fc.previewScroll != nil {
		fc.previewScroll.Content = widget.NewLabel("Select a file to preview")
		fc.previewScroll.Refresh()
	}
}

// SetKeep sets the keep state
func (fc *FilePreviewCard) SetKeep(keep bool) {
	fc.keep = keep
	if keep {
		fc.keepRadio.SetSelected("Keep")
	} else {
		fc.keepRadio.SetSelected("Delete")
	}
}
