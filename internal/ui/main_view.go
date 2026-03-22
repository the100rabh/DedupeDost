package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"

	"github.com/the100rabh/DedupeDost/internal/scanner"
)

// NewMainView creates the main view container
func NewMainView(dupDelApp *DedupeDostApp) *MainView {
	mv := &MainView{}

	// Create components
	mv.directoryBar = NewDirectoryBar(dupDelApp)
	mv.filterPanel = NewFilterPanel()
	mv.scanView = NewScanView(dupDelApp)
	mv.resultsView = NewResultsView(dupDelApp)
	mv.previewView = NewPreviewView(dupDelApp)
	mv.scriptView = NewScriptView(dupDelApp)
	mv.statusBar = NewStatusBar()

	// Create tabs
	mv.tabs = container.NewAppTabs(
		container.NewTabItemWithIcon("Scan", theme.FolderIcon(), mv.scanView.container),
		container.NewTabItemWithIcon("Results", theme.ListIcon(), mv.resultsView.container),
		container.NewTabItemWithIcon("Preview", theme.FileIcon(), mv.previewView.container),
		container.NewTabItemWithIcon("Script", theme.FileIcon(), mv.scriptView.container),
	)
	mv.tabs.SetTabLocation(container.TabLocationTop)

	// Create main container with proper sizing
	topArea := container.NewVBox(
		mv.directoryBar.container,
		container.NewHScroll(mv.filterPanel.container),
	)
	
	mv.container = container.NewBorder(
		topArea,
		mv.statusBar.container,
		nil,
		nil,
		mv.tabs,
	)

	return mv
}

// Container returns the main container
func (mv *MainView) Container() fyne.CanvasObject {
	return mv.container
}

// showDirectoryDialog shows the directory selection dialog
func (mv *MainView) showDirectoryDialog() {
	dialog.ShowFolderOpen(func(lister fyne.ListableURI, err error) {
		if err != nil {
			return
		}
		if lister == nil {
			return
		}

		path := lister.Path()
		mv.directoryBar.SetPath(path)
	}, fyne.CurrentApp().Driver().AllWindows()[0])
}

// startScan starts the scanning process
func (mv *MainView) startScan(dir string, options FilterOptions) {
	mv.tabs.SelectIndex(0) // Switch to scan tab
	
	// Convert FilterOptions to ScanOptions
	scanOpts := scanner.ScanOptions{
		Recursive:      options.Recursive,
		IncludeHidden:  options.IncludeHidden,
		FollowSymlinks: false,
		MinSize:        options.MinSize,
		MaxSize:        options.MaxSize,
		Extensions:     options.Extensions,
	}
	
	mv.scanView.startScan(dir, scanOpts)
}

// stopScan stops the scanning process
func (mv *MainView) stopScan() {
	mv.scanView.stopScan()
}

// generateScript generates the cleanup script
func (mv *MainView) generateScript() {
	mv.scriptView.generate()
	mv.tabs.SelectIndex(3) // Switch to script tab
}

// showResults switches to results tab
func (mv *MainView) showResults() {
	mv.tabs.SelectIndex(1)
}

// showPreview switches to preview tab
func (mv *MainView) showPreview() {
	mv.tabs.SelectIndex(2)
}
