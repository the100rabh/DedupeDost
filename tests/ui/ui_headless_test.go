//go:build gui

package ui_test

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"

	"github.com/dupdel/dup-del/internal/ui"
)

// TestFilterPanel_Creation tests filter panel creation
func TestFilterPanel_Creation(t *testing.T) {
	test.NewTempApp(t)

	var filterPanel *ui.FilterPanel
	fyne.DoAndWait(func() {
		filterPanel = ui.NewFilterPanel()
	})

	if filterPanel == nil {
		t.Fatal("FilterPanel should not be nil")
	}

	if filterPanel.Container() == nil {
		t.Fatal("FilterPanel container should not be nil")
	}
}

// TestFilterPanel_GetOptions tests getting filter options
func TestFilterPanel_GetOptions(t *testing.T) {
	test.NewTempApp(t)

	var filterPanel *ui.FilterPanel
	fyne.DoAndWait(func() {
		filterPanel = ui.NewFilterPanel()
	})

	opts := filterPanel.GetOptions()

	// Default options
	if opts.Recursive != true {
		t.Error("Recursive should be true by default")
	}

	if opts.IncludeHidden != false {
		t.Error("IncludeHidden should be false by default")
	}
}

// TestStatusBar_Creation tests status bar creation
func TestStatusBar_Creation(t *testing.T) {
	test.NewTempApp(t)

	var statusBar *ui.StatusBar
	fyne.DoAndWait(func() {
		statusBar = ui.NewStatusBar()
	})

	if statusBar == nil {
		t.Fatal("StatusBar should not be nil")
	}

	if statusBar.Container() == nil {
		t.Fatal("StatusBar container should not be nil")
	}
}

// TestStatusBar_SetStatus tests setting status text
func TestStatusBar_SetStatus(t *testing.T) {
	test.NewTempApp(t)

	var statusBar *ui.StatusBar
	fyne.DoAndWait(func() {
		statusBar = ui.NewStatusBar()
		statusBar.SetStatus("Testing...")
		statusBar.SetStatus("Ready")
	})
}

// TestStatusBar_SetProgress tests setting progress
func TestStatusBar_SetProgress(t *testing.T) {
	test.NewTempApp(t)

	var statusBar *ui.StatusBar
	fyne.DoAndWait(func() {
		statusBar = ui.NewStatusBar()

		// Should not panic with various values
		statusBar.SetProgress(0.0)
		statusBar.SetProgress(0.5)
		statusBar.SetProgress(1.0)
		statusBar.SetProgress(0.25)
	})
}

// TestStatusBar_SetStatistics tests setting statistics
func TestStatusBar_SetStatistics(t *testing.T) {
	test.NewTempApp(t)

	var statusBar *ui.StatusBar
	fyne.DoAndWait(func() {
		statusBar = ui.NewStatusBar()

		// Should not panic
		statusBar.SetStatistics(100, 10, 1024*1024)
		statusBar.SetStatistics(0, 0, 0)
	})
}

// TestHoverableCard_Creation tests hoverable card creation
func TestHoverableCard_Creation(t *testing.T) {
	test.NewTempApp(t)

	var card *ui.HoverableCard
	clicked := false

	fyne.DoAndWait(func() {
		card = ui.NewHoverableCard(nil, func() {
			clicked = true
		})
	})

	if card == nil {
		t.Fatal("HoverableCard should not be nil")
	}

	// Test tap
	fyne.DoAndWait(func() {
		card.Tapped(nil)
	})
	if !clicked {
		t.Error("Card should have been clicked")
	}
}

// TestHoverableCard_Hover tests hover behavior
func TestHoverableCard_Hover(t *testing.T) {
	test.NewTempApp(t)

	var card *ui.HoverableCard
	fyne.DoAndWait(func() {
		card = ui.NewHoverableCard(nil, func() {})
	})

	// Should not panic
	fyne.DoAndWait(func() {
		card.MouseIn(nil)
		card.MouseOut()
		card.MouseMoved(nil)
	})
}

// TestFilePreviewCard_Creation tests file preview card creation
func TestFilePreviewCard_Creation(t *testing.T) {
	test.NewTempApp(t)

	var card *ui.FilePreviewCard
	fyne.DoAndWait(func() {
		card = ui.NewFilePreviewCard()
	})

	if card == nil {
		t.Fatal("FilePreviewCard should not be nil")
	}

	if card.Container() == nil {
		t.Fatal("FilePreviewCard container should not be nil")
	}

	// Should not panic
	fyne.DoAndWait(func() {
		card.Clear()
	})
}

// TestFilePreviewCard_SetKeep tests setting keep state
func TestFilePreviewCard_SetKeep(t *testing.T) {
	test.NewTempApp(t)

	var card *ui.FilePreviewCard
	fyne.DoAndWait(func() {
		card = ui.NewFilePreviewCard()
	})

	// Should not panic
	fyne.DoAndWait(func() {
		card.SetKeep(true)
		card.SetKeep(false)
	})
}

// TestFormatSize tests the formatSize helper function
func TestFormatSize(t *testing.T) {
	// Placeholder - formatSize is not exported
	t.Skip("formatSize is not exported")
}
