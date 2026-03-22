package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// DirectoryBar represents the directory selection bar
type DirectoryBar struct {
	container     *fyne.Container
	pathEntry     *widget.Entry
	browseButton  *widget.Button
	startButton   *widget.Button
	dupDelApp     *DedupeDostApp
}

// NewDirectoryBar creates a new directory bar
func NewDirectoryBar(dupDelApp *DedupeDostApp) *DirectoryBar {
	db := &DirectoryBar{
		dupDelApp: dupDelApp,
	}

	// Path entry
	db.pathEntry = widget.NewEntry()
	db.pathEntry.SetPlaceHolder("Select a directory to scan...")
	db.pathEntry.OnChanged = func(path string) {
		if dupDelApp != nil {
			dupDelApp.SetScanDirectory(path)
		}
	}

	// Browse button
	db.browseButton = widget.NewButtonWithIcon("Browse...", theme.FolderOpenIcon(), func() {
		db.showDirectoryDialog()
	})

	// Start scan button
	db.startButton = widget.NewButtonWithIcon("Start Scan", theme.MediaPlayIcon(), func() {
		if dupDelApp != nil {
			dupDelApp.StartScan()
		}
	})
	db.startButton.Disable()

	// Create container with GridBagLayout for proper sizing
	// This allows the entry to expand and fill available space
	db.container = container.NewVBox(
		db.pathEntry,
		container.NewHBox(
			widget.NewLabel("📁 Directory:"),
			layout.NewSpacer(),
			db.browseButton,
			db.startButton,
		),
	)

	return db
}

// showDirectoryDialog shows the directory selection dialog
func (db *DirectoryBar) showDirectoryDialog() {
	dialog.ShowFolderOpen(func(lister fyne.ListableURI, err error) {
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to open directory: %w", err),
				fyne.CurrentApp().Driver().AllWindows()[0])
			return
		}
		if lister == nil {
			return
		}

		path := lister.Path()
		db.SetPath(path)
		if db.dupDelApp != nil {
			db.dupDelApp.SetScanDirectory(path)
		}
	}, fyne.CurrentApp().Driver().AllWindows()[0])
}

// Container returns the container
func (db *DirectoryBar) Container() *fyne.Container {
	return db.container
}

// SetPath sets the directory path
func (db *DirectoryBar) SetPath(path string) {
	db.pathEntry.SetText(path)
	db.pathEntry.Refresh()
}

// GetPath returns the current directory path
func (db *DirectoryBar) GetPath() string {
	return db.pathEntry.Text
}

// enableStartButton enables or disables the start button
func (db *DirectoryBar) enableStartButton(enabled bool) {
	if enabled {
		db.startButton.Enable()
	} else {
		db.startButton.Disable()
	}
}
