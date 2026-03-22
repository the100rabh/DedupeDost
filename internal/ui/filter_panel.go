package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// FilterOptions contains scan filter options
type FilterOptions struct {
	Recursive     bool
	IncludeHidden bool
	MinSize       int64
	MaxSize       int64
	Extensions    []string
}

// FilterPanel represents the filter options panel
type FilterPanel struct {
	container       *fyne.Container
	recursiveCheck  *widget.Check
	hiddenCheck     *widget.Check
	minSizeSelect   *widget.Select
	maxSizeSelect   *widget.Select
	typesButton     *widget.Button
	fileTypes       map[string]bool
}

// NewFilterPanel creates a new filter panel
func NewFilterPanel() *FilterPanel {
	fp := &FilterPanel{
		fileTypes: map[string]bool{
			"Images":  true,
			"Videos":  true,
			"Text":    true,
			"Other":   true,
		},
	}
	
	// Recursive check
	fp.recursiveCheck = widget.NewCheck("Recursive", nil)
	fp.recursiveCheck.SetChecked(true)
	
	// Include hidden check
	fp.hiddenCheck = widget.NewCheck("Include Hidden", nil)
	
	// Min size select
	fp.minSizeSelect = widget.NewSelect([]string{
		"No limit",
		"1 KB",
		"10 KB",
		"100 KB",
		"1 MB",
		"10 MB",
	}, nil)
	fp.minSizeSelect.SetSelected("No limit")
	
	// Max size select
	fp.maxSizeSelect = widget.NewSelect([]string{
		"No limit",
		"1 MB",
		"10 MB",
		"100 MB",
		"1 GB",
		"10 GB",
	}, nil)
	fp.maxSizeSelect.SetSelected("No limit")
	
	// File types button
	fp.typesButton = widget.NewButton("File Types...", func() {
		fp.showFileTypesDialog()
	})
	
	// Create container
	fp.container = container.NewHBox(
		fp.recursiveCheck,
		fp.hiddenCheck,
		widget.NewLabel("Min:"),
		fp.minSizeSelect,
		widget.NewLabel("Max:"),
		fp.maxSizeSelect,
		fp.typesButton,
	)
	
	return fp
}

// showFileTypesDialog shows the file types selection dialog
func (fp *FilterPanel) showFileTypesDialog() {
	imageCheck := widget.NewCheck("Images (jpg, png, gif, etc.)", nil)
	imageCheck.SetChecked(fp.fileTypes["Images"])
	
	videoCheck := widget.NewCheck("Videos (mp4, avi, mkv, etc.)", nil)
	videoCheck.SetChecked(fp.fileTypes["Videos"])
	
	textCheck := widget.NewCheck("Text (txt, md, json, code, etc.)", nil)
	textCheck.SetChecked(fp.fileTypes["Text"])
	
	otherCheck := widget.NewCheck("Other files", nil)
	otherCheck.SetChecked(fp.fileTypes["Other"])
	
	content := container.NewVBox(
		widget.NewLabel("Select file types to scan:"),
		imageCheck,
		videoCheck,
		textCheck,
		otherCheck,
	)
	
	d := dialog.NewCustom("File Types", "OK", content, fyne.CurrentApp().Driver().AllWindows()[0])
	d.SetOnClosed(func() {
		fp.fileTypes["Images"] = imageCheck.Checked
		fp.fileTypes["Videos"] = videoCheck.Checked
		fp.fileTypes["Text"] = textCheck.Checked
		fp.fileTypes["Other"] = otherCheck.Checked
	})
	d.Show()
}

// Container returns the container
func (fp *FilterPanel) Container() *fyne.Container {
	return fp.container
}

// GetOptions returns the current filter options
func (fp *FilterPanel) GetOptions() FilterOptions {
	opts := FilterOptions{
		Recursive:     fp.recursiveCheck.Checked,
		IncludeHidden: fp.hiddenCheck.Checked,
		Extensions:    fp.getFileExtensions(),
	}
	
	// Parse min size
	switch fp.minSizeSelect.Selected {
	case "1 KB":
		opts.MinSize = 1024
	case "10 KB":
		opts.MinSize = 10 * 1024
	case "100 KB":
		opts.MinSize = 100 * 1024
	case "1 MB":
		opts.MinSize = 1024 * 1024
	case "10 MB":
		opts.MinSize = 10 * 1024 * 1024
	}
	
	// Parse max size
	switch fp.maxSizeSelect.Selected {
	case "1 MB":
		opts.MaxSize = 1024 * 1024
	case "10 MB":
		opts.MaxSize = 10 * 1024 * 1024
	case "100 MB":
		opts.MaxSize = 100 * 1024 * 1024
	case "1 GB":
		opts.MaxSize = 1024 * 1024 * 1024
	case "10 GB":
		opts.MaxSize = 10 * 1024 * 1024 * 1024
	}
	
	return opts
}

// getFileExtensions returns file extensions based on selected types
func (fp *FilterPanel) getFileExtensions() []string {
	var extensions []string
	
	if fp.fileTypes["Images"] {
		extensions = append(extensions, []string{
			".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg", ".ico", ".tiff",
		}...)
	}
	
	if fp.fileTypes["Videos"] {
		extensions = append(extensions, []string{
			".mp4", ".avi", ".mkv", ".mov", ".wmv", ".flv", ".webm",
		}...)
	}
	
	if fp.fileTypes["Text"] {
		extensions = append(extensions, []string{
			".txt", ".md", ".json", ".xml", ".yaml", ".yml", ".csv", ".log",
			".html", ".css", ".js", ".ts", ".go", ".py", ".java", ".c", ".cpp",
		}...)
	}
	
	if fp.fileTypes["Other"] {
		// Return empty to include all
		return []string{}
	}
	
	return extensions
}

// SetEnabled enables or disables the filter panel
func (fp *FilterPanel) SetEnabled(enabled bool) {
	// Note: Fyne widgets don't have SetEnabled in older versions
	// This is a no-op for compatibility
	_ = enabled
}
