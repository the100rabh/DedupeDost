package ui

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/dupdel/dup-del/internal/detector"
	"github.com/dupdel/dup-del/internal/scanner"
)

// ScanView represents the scan progress view
type ScanView struct {
	container     *fyne.Container
	scroll        *container.Scroll
	pathLabel     *widget.Label
	progressBar   *widget.ProgressBar
	progressText  *widget.Label
	filesLabel    *widget.Label
	currentFile   *widget.Label
	statsBox      *fyne.Container
	stopButton    *widget.Button
	startButton   *widget.Button
	dupDelApp     *DupDelApp
	scanner       *scanner.Scanner
	detector      *detector.Detector
	isScanning    int32
}

// NewScanView creates a new scan view
func NewScanView(dupDelApp *DupDelApp) *ScanView {
	sv := &ScanView{
		dupDelApp: dupDelApp,
	}
	
	// Path label
	sv.pathLabel = widget.NewLabel("No directory selected")
	sv.pathLabel.TextStyle = fyne.TextStyle{Bold: true}
	
	// Progress bar
	sv.progressBar = widget.NewProgressBar()
	sv.progressBar.Hide()
	
	// Progress text
	sv.progressText = widget.NewLabel("")
	
	// Files label
	sv.filesLabel = widget.NewLabel("Files: 0")
	
	// Current file label
	sv.currentFile = widget.NewLabel("")
	sv.currentFile.Truncation = fyne.TextTruncateEllipsis
	
	// Stats box
	sv.statsBox = container.NewVBox(
		widget.NewLabel("Statistics will appear here..."),
	)

	// Stop button
	sv.stopButton = widget.NewButtonWithIcon("Stop", theme.MediaStopIcon(), func() {
		sv.stopScan()
	})
	sv.stopButton.Hide()
	sv.stopButton.Importance = widget.DangerImportance

	// Button container
	buttonBox := container.NewHBox(
		layout.NewSpacer(),
		sv.stopButton,
		layout.NewSpacer(),
	)
	
	// Create container
	content := container.NewVBox(
		widget.NewSeparator(),
		sv.pathLabel,
		widget.NewSeparator(),
		sv.progressBar,
		sv.progressText,
		sv.filesLabel,
		sv.currentFile,
		widget.NewSeparator(),
		sv.statsBox,
		widget.NewSeparator(),
		buttonBox,
	)

	sv.scroll = container.NewScroll(content)
	sv.container = container.NewVBox(sv.scroll)

	return sv
}

// Container returns the container
func (sv *ScanView) Container() *fyne.Container {
	return sv.container
}

// startScan starts the scanning process
func (sv *ScanView) startScan(dir string, options scanner.ScanOptions) {
	if dir == "" {
		dialog.ShowError(fmt.Errorf("please select a directory first"), 
			fyne.CurrentApp().Driver().AllWindows()[0])
		return
	}
	
	if atomic.LoadInt32(&sv.isScanning) == 1 {
		return
	}
	
	atomic.StoreInt32(&sv.isScanning, 1)

	// Update UI
	sv.pathLabel.SetText("📁 " + dir)
	sv.progressBar.Show()
	sv.progressBar.SetValue(0)
	sv.progressText.SetText("Initializing scan...")
	sv.filesLabel.SetText("Files: 0")
	sv.currentFile.SetText("")
	sv.stopButton.Show()

	// Disable directory bar start button
	if sv.dupDelApp != nil && sv.dupDelApp.mainView != nil {
		sv.dupDelApp.mainView.directoryBar.enableStartButton(false)
	}
	
	// Update status bar
	if sv.dupDelApp != nil && sv.dupDelApp.mainView != nil {
		sv.dupDelApp.mainView.statusBar.SetStatus("Scanning...")
		sv.dupDelApp.mainView.statusBar.SetProgressVisible(true)
	}
	
	// Start scanning in goroutine
	go func() {
		sv.runScan(dir, options)
	}()
}

// runScan runs the scan process
func (sv *ScanView) runScan(dir string, options scanner.ScanOptions) {
	startTime := time.Now()
	
	// Create detector
	sv.detector = detector.NewDetector(dir, options)
	
	// Run detection
	err := sv.detector.Detect()
	
	// Update UI on main thread
	fyne.Do(func() {
		atomic.StoreInt32(&sv.isScanning, 0)
		sv.stopButton.Hide()

		if sv.dupDelApp != nil && sv.dupDelApp.mainView != nil {
			sv.dupDelApp.mainView.statusBar.SetProgressVisible(false)
			// Re-enable directory bar start button
			sv.dupDelApp.mainView.directoryBar.enableStartButton(true)
		}
		
		if err != nil {
			if err == context.Canceled {
				sv.progressText.SetText("Scan cancelled")
				sv.dupDelApp.mainView.statusBar.SetStatus("Cancelled")
			} else {
				sv.progressText.SetText("Error: " + err.Error())
				sv.dupDelApp.mainView.statusBar.SetStatus("Error")
				dialog.ShowError(fmt.Errorf("scan failed: %w", err),
					fyne.CurrentApp().Driver().AllWindows()[0])
			}
			return
		}
		
		// Scan completed successfully
		stats := sv.detector.GetStatistics()
		duration := time.Since(startTime)
		
		sv.progressBar.SetValue(1.0)
		sv.progressText.SetText(fmt.Sprintf("Scan completed in %v", duration.Round(time.Millisecond)))
		sv.filesLabel.SetText(fmt.Sprintf("Files: %d", stats.TotalFiles))
		
		// Update stats box
		sv.statsBox = container.NewVBox(
			widget.NewLabel(fmt.Sprintf("✅ Total files: %d", stats.TotalFiles)),
			widget.NewLabel(fmt.Sprintf("📦 Total size: %s", formatSize(stats.TotalSize))),
			widget.NewLabel(fmt.Sprintf("🔄 Duplicate groups: %d", stats.DuplicateGroups)),
			widget.NewLabel(fmt.Sprintf("📄 Duplicate files: %d", stats.DuplicateFiles)),
			widget.NewLabel(fmt.Sprintf("💾 Recoverable: %s", formatSize(stats.RecoverableSize))),
		)
		sv.statsBox.Refresh()
		
		// Update status bar
		sv.dupDelApp.mainView.statusBar.SetStatistics(
			stats.TotalFiles,
			stats.DuplicateFiles,
			stats.RecoverableSize,
		)
		sv.dupDelApp.mainView.statusBar.SetStatus("Scan complete")

		// Update results view
		sv.dupDelApp.mainView.resultsView.SetGroups(sv.detector.GetGroups())

		// Show results tab
		sv.dupDelApp.mainView.showResults()
	})
}

// stopScan stops the scanning process
func (sv *ScanView) stopScan() {
	if sv.detector != nil {
		sv.detector.Cancel()
	}
}

// IsScanning returns whether a scan is in progress
func (sv *ScanView) IsScanning() bool {
	return atomic.LoadInt32(&sv.isScanning) == 1
}
