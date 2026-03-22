//go:build gui

package screenshot_test

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"

	"github.com/dupdel/dup-del/internal/script"
	"github.com/dupdel/dup-del/internal/ui"
	"github.com/dupdel/dup-del/pkg/models"
	"github.com/dupdel/dup-del/pkg/types"
	"github.com/dupdel/dup-del/tests/ui/screenshot"
)

var update = flag.Bool("update", false, "Update golden images")

// TestScreenshot_DirectoryBar tests DirectoryBar screenshot
func TestScreenshot_DirectoryBar(t *testing.T) {
	app := test.NewTempApp(t)

	var dirBar *ui.DirectoryBar
	var window fyne.Window

	fyne.DoAndWait(func() {
		dupDelApp := ui.NewDupDelAppWithFyneApp(app)
		dirBar = ui.NewDirectoryBar(dupDelApp)
		window = test.NewWindow(dirBar.Container())
	})
	defer window.Close()

	img, err := screenshot.Capture(window.Canvas())
	if err != nil {
		t.Skipf("Screenshot capture not supported: %v", err)
	}

	goldenPath := screenshot.GoldenPath("directory_bar")
	screenshot.AssertMatchesGolden(t, img, goldenPath, *update)
}

// TestScreenshot_FilterPanel tests FilterPanel screenshot
func TestScreenshot_FilterPanel(t *testing.T) {
	test.NewTempApp(t)

	var filterPanel *ui.FilterPanel
	var window fyne.Window
	
	fyne.DoAndWait(func() {
		filterPanel = ui.NewFilterPanel()
		window = test.NewWindow(filterPanel.Container())
	})
	defer window.Close()

	img, err := screenshot.Capture(window.Canvas())
	if err != nil {
		t.Skipf("Screenshot capture not supported: %v", err)
	}

	goldenPath := screenshot.GoldenPath("filter_panel")
	screenshot.AssertMatchesGolden(t, img, goldenPath, *update)
}

// TestScreenshot_StatusBar tests StatusBar screenshot
func TestScreenshot_StatusBar(t *testing.T) {
	test.NewTempApp(t)

	var statusBar *ui.StatusBar
	var window fyne.Window

	fyne.DoAndWait(func() {
		statusBar = ui.NewStatusBar()
		window = test.NewWindow(statusBar.Container())

		// Set some test values
		statusBar.SetStatus("Testing...")
		statusBar.SetProgress(0.5)
		statusBar.SetStatistics(100, 10, 1024*1024)
	})
	defer window.Close()

	img, err := screenshot.Capture(window.Canvas())
	if err != nil {
		t.Skipf("Screenshot capture not supported: %v", err)
	}

	goldenPath := screenshot.GoldenPath("status_bar")
	screenshot.AssertMatchesGolden(t, img, goldenPath, *update)
}

// TestScreenshot_ScanView tests ScanView screenshot
func TestScreenshot_ScanView(t *testing.T) {
	app := test.NewTempApp(t)

	var scanView *ui.ScanView
	var window fyne.Window

	fyne.DoAndWait(func() {
		dupDelApp := ui.NewDupDelAppWithFyneApp(app)
		scanView = ui.NewScanView(dupDelApp)
		window = test.NewWindow(scanView.Container())
	})
	defer window.Close()

	img, err := screenshot.Capture(window.Canvas())
	if err != nil {
		t.Skipf("Screenshot capture not supported: %v", err)
	}

	goldenPath := screenshot.GoldenPath("scan_view")
	screenshot.AssertMatchesGolden(t, img, goldenPath, *update)
}

// TestScreenshot_ResultsView tests ResultsView screenshot
func TestScreenshot_ResultsView(t *testing.T) {
	app := test.NewTempApp(t)

	var resultsView *ui.ResultsView
	var window fyne.Window

	fyne.DoAndWait(func() {
		dupDelApp := ui.NewDupDelAppWithFyneApp(app)
		resultsView = ui.NewResultsView(dupDelApp)
		window = test.NewWindow(resultsView.Container())
		window.Resize(fyne.NewSize(1200, 800))

		// Create test duplicate groups with actual files
		tmpDir := t.TempDir()

		// Group 1: Duplicate text files (3 copies)
		file1Path := filepath.Join(tmpDir, "readme.txt")
		file2Path := filepath.Join(tmpDir, "readme_copy.txt")
		file3Path := filepath.Join(tmpDir, "readme_backup.txt")
		os.WriteFile(file1Path, []byte("This is a sample readme file\nWith multiple lines\nFor testing duplicates"), 0644)
		os.WriteFile(file2Path, []byte("This is a sample readme file\nWith multiple lines\nFor testing duplicates"), 0644)
		os.WriteFile(file3Path, []byte("This is a sample readme file\nWith multiple lines\nFor testing duplicates"), 0644)

		// Group 2: Duplicate video files (2 copies)
		mp4Data := []byte{
			0x00, 0x00, 0x00, 0x14, 'f', 't', 'y', 'p',
			'i', 's', 'o', 'm', 0x00, 0x00, 0x00, 0x01,
			'i', 's', 'o', 'm',
			0x00, 0x00, 0x00, 0x08, 'm', 'o', 'o', 'v',
		}
		config1Path := filepath.Join(tmpDir, "vacation.mp4")
		config2Path := filepath.Join(tmpDir, "vacation_copy.mp4")
		os.WriteFile(config1Path, mp4Data, 0644)
		os.WriteFile(config2Path, mp4Data, 0644)

		// Create file entries
		f1, _ := models.NewFileEntry(file1Path)
		f2, _ := models.NewFileEntry(file2Path)
		f3, _ := models.NewFileEntry(file3Path)
		v1, _ := models.NewFileEntry(config1Path)
		v2, _ := models.NewFileEntry(config2Path)

		// Create duplicate groups
		group1 := models.DuplicateGroup{
			ID:        "group1",
			Hash:      "abc123",
			Size:      f1.Size,
			Extension: ".txt",
			Files:     []models.FileEntry{*f1, *f2, *f3},
		}
		group2 := models.DuplicateGroup{
			ID:        "group2",
			Hash:      "def456",
			Size:      v1.Size,
			Extension: ".mp4",
			FileType:  types.FileTypeVideo,
			Files:     []models.FileEntry{*v1, *v2},
		}

		// Mark files for deletion (keep first in each group)
		group1.MarkForDeletion([]int{0})
		group2.MarkForDeletion([]int{0})

		// Set groups on results view
		resultsView.SetGroups([]models.DuplicateGroup{group1, group2})
	})
	defer window.Close()

	img, err := screenshot.Capture(window.Canvas())
	if err != nil {
		t.Skipf("Screenshot capture not supported: %v", err)
	}

	goldenPath := screenshot.GoldenPath("results_view")
	screenshot.AssertMatchesGolden(t, img, goldenPath, *update)
}

// TestScreenshot_PreviewView tests PreviewView screenshot
func TestScreenshot_PreviewView(t *testing.T) {
	app := test.NewTempApp(t)

	var previewView *ui.PreviewView
	var window fyne.Window

	fyne.DoAndWait(func() {
		dupDelApp := ui.NewDupDelAppWithFyneApp(app)
		previewView = ui.NewPreviewView(dupDelApp)
		window = test.NewWindow(previewView.Container())
		window.Resize(fyne.NewSize(1000, 700))

		// Create test duplicate group with 3 text files (testing >2 files support)
		tmpDir := t.TempDir()
		file1Path := filepath.Join(tmpDir, "readme_original.txt")
		file2Path := filepath.Join(tmpDir, "readme_copy1.txt")
		file3Path := filepath.Join(tmpDir, "readme_copy2.txt")
		os.WriteFile(file1Path, []byte("This is the original readme file\nWith multiple lines\nFor testing duplicates"), 0644)
		os.WriteFile(file2Path, []byte("This is the original readme file\nWith multiple lines\nFor testing duplicates"), 0644)
		os.WriteFile(file3Path, []byte("This is the original readme file\nWith multiple lines\nFor testing duplicates"), 0644)

		file1, _ := models.NewFileEntry(file1Path)
		file2, _ := models.NewFileEntry(file2Path)
		file3, _ := models.NewFileEntry(file3Path)

		group := models.DuplicateGroup{
			Files: []models.FileEntry{*file1, *file2, *file3},
			Size:  file1.Size,
		}
		// Mark file1 to keep, others for deletion
		group.MarkForDeletion([]int{0})

		// Set the group on preview view
		previewView.SetGroups([]models.DuplicateGroup{group})
	})
	defer window.Close()

	img, err := screenshot.Capture(window.Canvas())
	if err != nil {
		t.Skipf("Screenshot capture not supported: %v", err)
	}

	goldenPath := screenshot.GoldenPath("preview_view")
	screenshot.AssertMatchesGolden(t, img, goldenPath, *update)
}

// TestScreenshot_VideoPreview tests video preview screenshot
func TestScreenshot_VideoPreview(t *testing.T) {
	app := test.NewTempApp(t)

	var previewView *ui.PreviewView
	var window fyne.Window

	fyne.DoAndWait(func() {
		dupDelApp := ui.NewDupDelAppWithFyneApp(app)
		previewView = ui.NewPreviewView(dupDelApp)
		window = test.NewWindow(previewView.Container())
		window.Resize(fyne.NewSize(1000, 700))

		// Create test duplicate video files (minimal MP4 structure)
		tmpDir := t.TempDir()
		
		// Create minimal MP4 files for testing
		mp4Data := []byte{
			// ftyp box
			0x00, 0x00, 0x00, 0x14, 'f', 't', 'y', 'p',
			'i', 's', 'o', 'm', 0x00, 0x00, 0x00, 0x01,
			'i', 's', 'o', 'm',
			// moov box
			0x00, 0x00, 0x00, 0x08, 'm', 'o', 'o', 'v',
		}
		
		file1Path := filepath.Join(tmpDir, "vacation_original.mp4")
		file2Path := filepath.Join(tmpDir, "vacation_copy.mp4")
		os.WriteFile(file1Path, mp4Data, 0644)
		os.WriteFile(file2Path, mp4Data, 0644)

		file1, _ := models.NewFileEntry(file1Path)
		file2, _ := models.NewFileEntry(file2Path)

		group := models.DuplicateGroup{
			Files:    []models.FileEntry{*file1, *file2},
			Size:     file1.Size,
			FileType: types.FileTypeVideo,
		}
		// Mark file1 to keep
		group.MarkForDeletion([]int{0})

		// Set the group on preview view
		previewView.SetGroups([]models.DuplicateGroup{group})
	})
	defer window.Close()

	img, err := screenshot.Capture(window.Canvas())
	if err != nil {
		t.Skipf("Screenshot capture not supported: %v", err)
	}

	goldenPath := screenshot.GoldenPath("video_preview")
	screenshot.AssertMatchesGolden(t, img, goldenPath, *update)
}

// TestScreenshot_ScriptView tests ScriptView screenshot
func TestScreenshot_ScriptView(t *testing.T) {
	app := test.NewTempApp(t)

	var scriptView *ui.ScriptView
	var window fyne.Window

	fyne.DoAndWait(func() {
		dupDelApp := ui.NewDupDelAppWithFyneApp(app)
		scriptView = ui.NewScriptView(dupDelApp)
		window = test.NewWindow(scriptView.Container())
		window.Resize(fyne.NewSize(900, 600))

		// Create test duplicate groups with actual files
		tmpDir := t.TempDir()
		file1Path := filepath.Join(tmpDir, "readme.txt")
		file2Path := filepath.Join(tmpDir, "readme_copy.txt")
		file3Path := filepath.Join(tmpDir, "readme_backup.txt")
		os.WriteFile(file1Path, []byte("This is a sample readme file\nWith multiple lines\nFor testing duplicates"), 0644)
		os.WriteFile(file2Path, []byte("This is a sample readme file\nWith multiple lines\nFor testing duplicates"), 0644)
		os.WriteFile(file3Path, []byte("This is a sample readme file\nWith multiple lines\nFor testing duplicates"), 0644)

		// Create file entries
		f1, _ := models.NewFileEntry(file1Path)
		f2, _ := models.NewFileEntry(file2Path)
		f3, _ := models.NewFileEntry(file3Path)

		// Create duplicate group
		group := models.DuplicateGroup{
			ID:        "group1",
			Hash:      "abc123",
			Size:      f1.Size,
			Extension: ".txt",
			Files:     []models.FileEntry{*f1, *f2, *f3},
		}

		// Create session and mark files for deletion
		session := models.NewScanSession(tmpDir)
		session.DuplicateGroups = []models.DuplicateGroup{group}
		session.DuplicateGroups[0].MarkForDeletion([]int{0})
		session.CalculateStatistics()

		// Generate script content for preview
		outputPath := filepath.Join(tmpDir, "cleanup.sh")
		generator := script.NewGenerator(outputPath)
		generator.SetDryRunSupport(true)
		content := generatorForPreview(generator, session)

		// Display script in the view
		scriptView.SetScriptContent(content)
	})
	defer window.Close()

	img, err := screenshot.Capture(window.Canvas())
	if err != nil {
		t.Skipf("Screenshot capture not supported: %v", err)
	}

	goldenPath := screenshot.GoldenPath("script_view")
	screenshot.AssertMatchesGolden(t, img, goldenPath, *update)
}

// TestScreenshot_FilePreviewCard tests FilePreviewCard screenshot
func TestScreenshot_FilePreviewCard(t *testing.T) {
	test.NewTempApp(t)

	var card *ui.FilePreviewCard
	var window fyne.Window
	
	fyne.DoAndWait(func() {
		card = ui.NewFilePreviewCard()
		window = test.NewWindow(card.Container())
	})
	defer window.Close()

	img, err := screenshot.Capture(window.Canvas())
	if err != nil {
		t.Skipf("Screenshot capture not supported: %v", err)
	}

	goldenPath := screenshot.GoldenPath("file_preview_card")
	screenshot.AssertMatchesGolden(t, img, goldenPath, *update)
}

// TestScreenshot_HoverableCard tests HoverableCard screenshot
func TestScreenshot_HoverableCard(t *testing.T) {
	test.NewTempApp(t)

	var card fyne.CanvasObject
	var window fyne.Window
	
	fyne.DoAndWait(func() {
		content := widget.NewLabel("Test Content")
		card = ui.NewHoverableCard(content, func() {})
		window = test.NewWindow(card)
	})
	defer window.Close()

	img, err := screenshot.Capture(window.Canvas())
	if err != nil {
		t.Skipf("Screenshot capture not supported: %v", err)
	}

	goldenPath := screenshot.GoldenPath("hoverable_card")
	screenshot.AssertMatchesGolden(t, img, goldenPath, *update)
}

// TestScreenshot_MainWindow tests full main window screenshot
func TestScreenshot_MainWindow(t *testing.T) {
	app := test.NewTempApp(t)

	var mainView *ui.MainView
	var window fyne.Window

	fyne.DoAndWait(func() {
		dupDelApp := ui.NewDupDelAppWithFyneApp(app)
		mainView = ui.NewMainView(dupDelApp)
		window = test.NewWindow(mainView.Container())
		window.Resize(fyne.NewSize(1200, 800))
	})
	defer window.Close()

	img, err := screenshot.Capture(window.Canvas())
	if err != nil {
		t.Skipf("Screenshot capture not supported: %v", err)
	}

	goldenPath := screenshot.GoldenPath("main_window")
	screenshot.AssertMatchesGolden(t, img, goldenPath, *update)
}

// generatorForPreview generates script content for preview (helper for tests)
func generatorForPreview(gen *script.Generator, session *models.ScanSession) string {
	var content string

	content += "#!/bin/bash\n\n"
	content += "# ============================================\n"
	content += "# DupDel Cleanup Script\n"
	content += "# ============================================\n"
	content += fmt.Sprintf("# Generated: %s\n", session.StartTime.Format("2006-01-02 15:04:05"))
	content += fmt.Sprintf("# Source Directory: %s\n", session.SourceDir)
	content += fmt.Sprintf("# Total files to delete: %d\n", session.FilesToDelete)
	content += fmt.Sprintf("# Total space to recover: %s\n", session.GetRecoverableSizeHuman())
	content += "# ============================================\n\n"

	// DRY_RUN support
	content += `DRY_RUN=${DRY_RUN:-false}

delete_file() {
    local file="$1"
    if [ "$DRY_RUN" = "true" ]; then
        echo "[DRY RUN] Would delete: $file"
    else
        rm -f "$file"
    fi
}

`

	// Groups
	for i, group := range session.DuplicateGroups {
		content += fmt.Sprintf("# Group %d\n", i+1)
		for _, path := range group.DeletePaths {
			content += fmt.Sprintf("delete_file \"%s\"\n", path)
		}
		content += fmt.Sprintf("# Keeping: %s\n\n", group.Files[0].Path)
	}

	return content
}
