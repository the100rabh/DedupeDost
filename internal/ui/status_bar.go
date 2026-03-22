package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// StatusBar represents the status bar at the bottom of the window
type StatusBar struct {
	container     *fyne.Container
	statusLabel   *widget.Label
	progressBar   *widget.ProgressBar
	statsLabel    *widget.Label
}

// NewStatusBar creates a new status bar
func NewStatusBar() *StatusBar {
	sb := &StatusBar{}
	
	// Status label
	sb.statusLabel = widget.NewLabel("Ready")
	
	// Progress bar (hidden initially)
	sb.progressBar = widget.NewProgressBar()
	sb.progressBar.Hide()
	
	// Stats label
	sb.statsLabel = widget.NewLabel("")
	
	// Create container
	sb.container = container.NewHBox(
		sb.statusLabel,
		container.NewHScroll(sb.progressBar),
		sb.statsLabel,
	)
	
	return sb
}

// Container returns the container
func (sb *StatusBar) Container() *fyne.Container {
	return sb.container
}

// SetStatus sets the status text
func (sb *StatusBar) SetStatus(text string) {
	sb.statusLabel.SetText(text)
}

// SetProgress sets the progress value (0.0 to 1.0)
func (sb *StatusBar) SetProgress(value float64) {
	sb.progressBar.SetValue(value)
}

// SetProgressVisible shows or hides the progress bar
func (sb *StatusBar) SetProgressVisible(visible bool) {
	if visible {
		sb.progressBar.Show()
	} else {
		sb.progressBar.Hide()
	}
}

// SetStatistics sets the statistics display
func (sb *StatusBar) SetStatistics(files, duplicates int, size int64) {
	sizeStr := formatSize(size)
	sb.statsLabel.SetText(fmt.Sprintf("Files: %d | Duplicates: %d | Size: %s", files, duplicates, sizeStr))
}

// Clear clears the status bar
func (sb *StatusBar) Clear() {
	sb.statusLabel.SetText("Ready")
	sb.progressBar.SetValue(0)
	sb.statsLabel.SetText("")
}

// formatSize formats bytes to human-readable string
func formatSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	
	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	default:
		return fmt.Sprintf("%d B", size)
	}
}
