package ui

import (
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/dupdel/dup-del/pkg/models"
)

// ResultsView represents the results list view
type ResultsView struct {
	container      *fyne.Container
	content        *fyne.Container
	searchEntry    *widget.Entry
	sortSelect     *widget.Select
	selectAllBtn   *widget.Button
	deselectAllBtn *widget.Button
	generateBtn    *widget.Button
	groups         []models.DuplicateGroup
	filteredGroups []models.DuplicateGroup
	selectedGroup  int
	dupDelApp      *DupDelApp
	cards          []*resultsCard
	scroll         *container.Scroll
}

// resultsCard is a simple card for displaying a duplicate group
type resultsCard struct {
	container   *fyne.Container
	icon        *widget.Icon
	titleLabel  *widget.Label
	detailLabel *widget.Label
	filesList   *fyne.Container
	group       models.DuplicateGroup
}

// HoverableCard wraps content with hover and tap support
type HoverableCard struct {
	widget.BaseWidget
	content   fyne.CanvasObject
	onTap     func()
	hovering  bool
	hoverRect *canvas.Rectangle
}

// NewHoverableCard creates a new hoverable card
func NewHoverableCard(content fyne.CanvasObject, onTap func()) *HoverableCard {
	hc := &HoverableCard{
		content:   content,
		onTap:     onTap,
		hoverRect: canvas.NewRectangle(color.Transparent),
	}
	hc.ExtendBaseWidget(hc)
	return hc
}

// CreateRenderer implements fyne.Widget
func (hc *HoverableCard) CreateRenderer() fyne.WidgetRenderer {
	stack := container.NewStack(hc.hoverRect, hc.content)
	return &hoverableCardRenderer{widget: hc, stack: stack}
}

// Tapped implements fyne.Tappable
func (hc *HoverableCard) Tapped(_ *fyne.PointEvent) {
	if hc.onTap != nil {
		hc.onTap()
	}
}

// MouseIn implements desktop.Hoverable
func (hc *HoverableCard) MouseIn(_ *desktop.MouseEvent) {
	hc.hovering = true
	hc.hoverRect.FillColor = &color.NRGBA{R: 128, G: 128, B: 128, A: 25}
	hc.hoverRect.Refresh()
}

// MouseOut implements desktop.Hoverable  
func (hc *HoverableCard) MouseOut() {
	hc.hovering = false
	hc.hoverRect.FillColor = color.Transparent
	hc.hoverRect.Refresh()
}

// MouseMoved implements desktop.Hoverable
func (hc *HoverableCard) MouseMoved(_ *desktop.MouseEvent) {
}

type hoverableCardRenderer struct {
	widget *HoverableCard
	stack  *fyne.Container
}

func (r *hoverableCardRenderer) Layout(size fyne.Size) {
	r.stack.Resize(size)
}

func (r *hoverableCardRenderer) MinSize() fyne.Size {
	return r.stack.MinSize()
}

func (r *hoverableCardRenderer) Refresh() {
	r.stack.Refresh()
}

func (r *hoverableCardRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.stack}
}

func (r *hoverableCardRenderer) Destroy() {
}

// NewResultsView creates a new results view
func NewResultsView(dupDelApp *DupDelApp) *ResultsView {
	rv := &ResultsView{
		dupDelApp:      dupDelApp,
		groups:         make([]models.DuplicateGroup, 0),
		filteredGroups: make([]models.DuplicateGroup, 0),
		cards:          make([]*resultsCard, 0),
		selectedGroup:  -1,
	}

	// Search entry - full width on top row
	rv.searchEntry = widget.NewEntry()
	rv.searchEntry.SetPlaceHolder("Search files...")
	rv.searchEntry.OnChanged = func(text string) {
		rv.filterGroups(text)
	}

	// Sort select
	rv.sortSelect = widget.NewSelect([]string{
		"Size (largest first)",
		"Count (most copies first)",
		"Name (alphabetical)",
	}, func(value string) {
		rv.sortGroups(value)
	})

	// Select all button
	rv.selectAllBtn = widget.NewButtonWithIcon("Select All", theme.ContentAddIcon(), func() {
		rv.selectAll()
	})

	// Deselect all button
	rv.deselectAllBtn = widget.NewButtonWithIcon("Deselect All", theme.ContentRemoveIcon(), func() {
		rv.deselectAll()
	})

	// Generate button
	rv.generateBtn = widget.NewButtonWithIcon("Generate Script", theme.DocumentCreateIcon(), func() {
		rv.dupDelApp.GenerateScript()
	})
	rv.generateBtn.Importance = widget.HighImportance
	rv.generateBtn.Disable()

	// Content container (will hold cards)
	rv.content = container.NewVBox()

	// Scroll container
	rv.scroll = container.NewScroll(rv.content)

	// Toolbar - search on top row, controls on bottom row
	toolbar := container.NewVBox(
		rv.searchEntry,
		container.NewHBox(
			widget.NewLabel("Sort:"),
			rv.sortSelect,
			rv.selectAllBtn,
			rv.deselectAllBtn,
		),
	)

	// Bottom toolbar
	bottomToolbar := container.NewHBox(
		widget.NewLabel(""),
		rv.generateBtn,
	)

	// Create container
	rv.container = container.NewBorder(
		toolbar,
		bottomToolbar,
		nil,
		nil,
		rv.scroll,
	)

	return rv
}

// Container returns the container
func (rv *ResultsView) Container() *fyne.Container {
	return rv.container
}

// FilterPanel returns the filter panel
func (rv *ResultsView) FilterPanel() *FilterPanel {
	return nil // Would need to pass reference in constructor
}

// SetGroups sets the duplicate groups to display
func (rv *ResultsView) SetGroups(groups []models.DuplicateGroup) {
	rv.groups = groups
	rv.filterGroups("")
	rv.generateBtn.Enable()
}

// filterGroups filters groups based on search text
func (rv *ResultsView) filterGroups(text string) {
	text = strings.ToLower(text)

	rv.filteredGroups = make([]models.DuplicateGroup, 0)
	for _, group := range rv.groups {
		if text == "" {
			rv.filteredGroups = append(rv.filteredGroups, group)
			continue
		}

		// Search in file names
		for _, file := range group.Files {
			if strings.Contains(strings.ToLower(file.Name), text) {
				rv.filteredGroups = append(rv.filteredGroups, group)
				break
			}
			if strings.Contains(strings.ToLower(file.Path), text) {
				rv.filteredGroups = append(rv.filteredGroups, group)
				break
			}
		}
	}

	rv.refreshCards()
}

// sortGroups sorts groups based on criteria
func (rv *ResultsView) sortGroups(criteria string) {
	switch criteria {
	case "Size (largest first)":
		// Sort by size descending
		for i := 0; i < len(rv.filteredGroups)-1; i++ {
			for j := i + 1; j < len(rv.filteredGroups); j++ {
				if rv.filteredGroups[i].Size < rv.filteredGroups[j].Size {
					rv.filteredGroups[i], rv.filteredGroups[j] = rv.filteredGroups[j], rv.filteredGroups[i]
				}
			}
		}
	case "Count (most copies first)":
		// Sort by file count descending
		for i := 0; i < len(rv.filteredGroups)-1; i++ {
			for j := i + 1; j < len(rv.filteredGroups); j++ {
				if len(rv.filteredGroups[i].Files) < len(rv.filteredGroups[j].Files) {
					rv.filteredGroups[i], rv.filteredGroups[j] = rv.filteredGroups[j], rv.filteredGroups[i]
				}
			}
		}
	case "Name (alphabetical)":
		// Sort by name ascending
		for i := 0; i < len(rv.filteredGroups)-1; i++ {
			for j := i + 1; j < len(rv.filteredGroups); j++ {
				if rv.filteredGroups[i].Files[0].Name > rv.filteredGroups[j].Files[0].Name {
					rv.filteredGroups[i], rv.filteredGroups[j] = rv.filteredGroups[j], rv.filteredGroups[i]
				}
			}
		}
	}

	rv.refreshCards()
}

// refreshCards refreshes the card display
func (rv *ResultsView) refreshCards() {
	// Update existing cards or create new ones
	if len(rv.cards) == len(rv.filteredGroups) {
		// Same number of groups, update existing cards
		for i, group := range rv.filteredGroups {
			rv.updateCard(rv.cards[i], group, i)
		}
	} else {
		// Different number of groups, recreate all cards
		rv.content.Objects = make([]fyne.CanvasObject, 0)
		rv.cards = make([]*resultsCard, 0)

		for i, group := range rv.filteredGroups {
			card := rv.createCard(group, i)
			rv.cards = append(rv.cards, card)

			// Wrap in hoverable card for click support with subtle hover
			hoverCard := NewHoverableCard(card.container, func() {
				rv.onGroupSelected(i)
			})

			rv.content.Add(hoverCard)
		}
	}

	if rv.scroll != nil {
		rv.scroll.Refresh()
	}
}

// updateCard updates an existing card with new group data
func (rv *ResultsView) updateCard(card *resultsCard, group models.DuplicateGroup, index int) {
	card.group = group

	// Update icon based on file type
	switch group.FileType {
	case 1: // Text
		card.icon.SetResource(theme.DocumentIcon())
	case 2: // Image
		card.icon.SetResource(theme.FileIcon())
	case 3: // Video
		card.icon.SetResource(theme.MediaVideoIcon())
	default:
		card.icon.SetResource(theme.FileIcon())
	}

	// Update labels
	card.titleLabel.SetText(fmt.Sprintf("%s (%d copies)",
		group.Files[0].Name, len(group.Files)))
	card.detailLabel.SetText(fmt.Sprintf("%s each • Recoverable: %s",
		formatSize(group.Size),
		formatSize(group.GetRecoverableSize())))

	// Update file list with new decisions
	if card.filesList != nil {
		newFileList := rv.createFileList(group)
		card.filesList.Objects = newFileList.Objects
		card.filesList.Refresh()
	}

	card.container.Refresh()
}

// createCard creates a card for a duplicate group
func (rv *ResultsView) createCard(group models.DuplicateGroup, index int) *resultsCard {
	card := &resultsCard{}

	// Icon based on file type
	switch group.FileType {
	case 1: // Text
		card.icon = widget.NewIcon(theme.DocumentIcon())
	case 2: // Image
		card.icon = widget.NewIcon(theme.FileIcon())
	case 3: // Video
		card.icon = widget.NewIcon(theme.MediaVideoIcon())
	default:
		card.icon = widget.NewIcon(theme.FileIcon())
	}

	// Title label
	card.titleLabel = widget.NewLabel(fmt.Sprintf("%s (%d copies)",
		group.Files[0].Name, len(group.Files)))
	card.titleLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Detail label
	card.detailLabel = widget.NewLabel(fmt.Sprintf("%s each • Recoverable: %s",
		formatSize(group.Size),
		formatSize(group.GetRecoverableSize())))
	card.detailLabel.TextStyle = fyne.TextStyle{Italic: true}

	card.group = group

	// Create file list with paths and keep/delete indicators
	card.filesList = rv.createFileList(group)

	// Create card content
	card.container = container.NewVBox(
		container.NewHBox(
			card.icon,
			container.NewVBox(
				card.titleLabel,
				card.detailLabel,
			),
			layout.NewSpacer(),
			widget.NewIcon(theme.NavigateNextIcon()),
		),
		widget.NewSeparator(),
		card.filesList,
	)

	return card
}

// createFileList creates a list of files with paths and keep/delete indicators
func (rv *ResultsView) createFileList(group models.DuplicateGroup) *fyne.Container {
	fileList := container.NewVBox()

	for _, file := range group.Files {
		// Check if this file is marked for deletion
		isDeleted := false
		for _, delPath := range group.DeletePaths {
			if delPath == file.Path {
				isDeleted = true
				break
			}
		}

		// Create indicator icon with color
		var indicator *widget.Icon
		if isDeleted {
			indicator = widget.NewIcon(theme.DeleteIcon())
		} else {
			indicator = widget.NewIcon(theme.ConfirmIcon())
		}

		// File path label - truncate long paths
		pathLabel := widget.NewLabel(file.Path)
		pathLabel.Truncation = fyne.TextTruncateEllipsis

		// Keep/Delete status label
		statusText := "Keep"
		if isDeleted {
			statusText = "Delete"
		}
		statusLabel := widget.NewLabel(statusText)
		statusLabel.TextStyle = fyne.TextStyle{Bold: true}

		// File size label
		sizeLabel := widget.NewLabel(fmt.Sprintf("(%s)", file.GetDisplaySize()))

		// File entry
		fileEntry := container.NewHBox(
			indicator,
			container.NewVBox(
				pathLabel,
				container.NewHBox(
					statusLabel,
					layout.NewSpacer(),
					sizeLabel,
				),
			),
		)

		fileList.Add(fileEntry)
	}

	return fileList
}

// onGroupSelected handles group selection
func (rv *ResultsView) onGroupSelected(index int) {
	if index < 0 || index >= len(rv.filteredGroups) {
		return
	}

	rv.selectedGroup = index
	group := rv.filteredGroups[index]

	// Show preview
	if rv.dupDelApp != nil && rv.dupDelApp.mainView != nil {
		rv.dupDelApp.mainView.previewView.setGroup(group)
		rv.dupDelApp.mainView.showPreview()
	}
}

// selectAll selects all groups
func (rv *ResultsView) selectAll() {
	// Mark all files for deletion except first in each group
	for i := range rv.groups {
		rv.groups[i].MarkForDeletion([]int{0})
	}
	rv.refreshCards()
}

// deselectAll deselects all groups
func (rv *ResultsView) deselectAll() {
	for i := range rv.groups {
		rv.groups[i].ClearDecisions()
	}
	rv.refreshCards()
}

// getSelectedGroups returns groups with deletion decisions
func (rv *ResultsView) getSelectedGroups() []models.DuplicateGroup {
	var selected []models.DuplicateGroup
	for _, group := range rv.groups {
		if len(group.DeletePaths) > 0 {
			selected = append(selected, group)
		}
	}
	return selected
}
