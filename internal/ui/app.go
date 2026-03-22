package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
)

// DupDelApp represents the main application
type DupDelApp struct {
	app         fyne.App
	mainWindow  fyne.Window
	mainView    *MainView
	scanDir     string
}

// MainView represents the main content view
type MainView struct {
	container    *fyne.Container
	directoryBar *DirectoryBar
	filterPanel  *FilterPanel
	tabs         *container.AppTabs
	scanView     *ScanView
	resultsView  *ResultsView
	previewView  *PreviewView
	scriptView   *ScriptView
	statusBar    *StatusBar
}

// NewDupDelApp creates a new DupDel application
func NewDupDelApp() *DupDelApp {
	a := app.NewWithID("io.github.dupdel.app")

	d := &DupDelApp{
		app: a,
	}

	return d
}

// NewDupDelAppWithFyneApp creates a DupDelApp with a specific Fyne app (for testing)
func NewDupDelAppWithFyneApp(fyneApp fyne.App) *DupDelApp {
	d := &DupDelApp{
		app: fyneApp,
	}

	return d
}

// Run starts the application
func (d *DupDelApp) Run() {
	d.setupMainWindow()
	d.mainWindow.ShowAndRun()
}

// setupMainWindow creates and configures the main window
func (d *DupDelApp) setupMainWindow() {
	d.mainWindow = d.app.NewWindow("DupDel - Duplicate File Scanner")
	d.mainWindow.Resize(fyne.NewSize(1400, 900))
	d.mainWindow.CenterOnScreen()

	// Create main view
	d.mainView = NewMainView(d)

	// Set content
	d.mainWindow.SetContent(d.mainView.container)

	// Set menu
	d.mainWindow.SetMainMenu(d.createMainMenu())

	// Set close interceptor
	d.mainWindow.SetCloseIntercept(func() {
		d.onClose()
	})
}

// createMainMenu creates the application menu
func (d *DupDelApp) createMainMenu() *fyne.MainMenu {
	// File menu
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("Open Directory...", func() {
			d.mainView.showDirectoryDialog()
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Quit", func() {
			d.app.Quit()
		}),
	)
	
	// Edit menu
	editMenu := fyne.NewMenu("Edit",
		fyne.NewMenuItem("Select All", func() {
			d.mainView.resultsView.selectAll()
		}),
		fyne.NewMenuItem("Deselect All", func() {
			d.mainView.resultsView.deselectAll()
		}),
	)
	
	// View menu
	viewMenu := fyne.NewMenu("View",
		fyne.NewMenuItem("Scan View", func() {
			d.mainView.tabs.SelectIndex(0)
		}),
		fyne.NewMenuItem("Results View", func() {
			d.mainView.tabs.SelectIndex(1)
		}),
		fyne.NewMenuItem("Preview View", func() {
			d.mainView.tabs.SelectIndex(2)
		}),
		fyne.NewMenuItem("Script View", func() {
			d.mainView.tabs.SelectIndex(3)
		}),
	)
	
	// Help menu
	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("About", func() {
			d.showAboutDialog()
		}),
	)
	
	return fyne.NewMainMenu(fileMenu, editMenu, viewMenu, helpMenu)
}

// showAboutDialog shows the about dialog
func (d *DupDelApp) showAboutDialog() {
	dialog.ShowInformation("About DupDel", 
		"DupDel v1.0.0\n\nDuplicate File Scanner\n\nA powerful tool to find and remove duplicate files.\n\nFeatures:\n- Recursive directory scanning\n- SHA-256 hash-based detection\n- Side-by-side file preview\n- Safe bash script generation\n\nMade with ❤️ using Go and Fyne", 
		d.mainWindow)
}

// onClose handles application close
func (d *DupDelApp) onClose() {
	// Cleanup if needed
	d.app.Quit()
}

// SetScanDirectory sets the scan directory
func (d *DupDelApp) SetScanDirectory(dir string) {
	d.scanDir = dir
	if d.mainView != nil && d.mainView.directoryBar != nil {
		d.mainView.directoryBar.SetPath(dir)
		d.mainView.directoryBar.enableStartButton(true)
	}
}

// GetScanDirectory returns the current scan directory
func (d *DupDelApp) GetScanDirectory() string {
	return d.scanDir
}

// StartScan starts the scanning process
func (d *DupDelApp) StartScan() {
	if d.scanDir == "" {
		dialog.ShowError(fmt.Errorf("please select a directory first"), d.mainWindow)
		return
	}
	
	d.mainView.startScan(d.scanDir, d.mainView.filterPanel.GetOptions())
}

// StopScan stops the scanning process
func (d *DupDelApp) StopScan() {
	d.mainView.stopScan()
}

// GenerateScript generates the cleanup script
func (d *DupDelApp) GenerateScript() {
	d.mainView.generateScript()
}
