# DupDel UI Testing Guide

## Overview

This document describes how to test the DupDel GUI application.

## Test Files

### `screenshot/screenshot_test.go`
Screenshot-based visual regression tests. Compares UI against golden images.

### `ui_headless_test.go`
Unit tests for UI components that require GUI backend. Tagged with `//go:build gui`.

### `ui_stub_test.go`
Placeholder test for CI/CD environments without GUI support. Always runs.

## Running Tests

### With GUI Backend (Full Testing)

**Prerequisites:**
```bash
# Ubuntu/Debian
sudo apt-get install libgl1-mesa-dev xorg-dev

# macOS
xcode-select --install
```

**Run all UI tests:**
```bash
go test ./tests/ui/... -v -tags gui
```

**Run screenshot tests with golden image update:**
```bash
go test ./tests/ui/screenshot/... -v -tags gui -update
```

**Run specific test:**
```bash
go test ./tests/ui/... -v -tags gui -run TestFilterPanel
```

### Without GUI Backend (CI/CD)

**Run stub test only:**
```bash
go test ./tests/ui/ui_stub_test.go -v
```

**Default behavior (no gui tag):**
```bash
go test ./tests/ui/... -v  # Only runs stub test
```

## Test Categories

### 1. Component Creation Tests
Test that UI components can be instantiated:

```go
func TestComponent_Creation(t *testing.T) {
    component := ui.NewComponent()
    if component == nil {
        t.Fatal("Component should not be nil")
    }
}
```

### 2. Component Interaction Tests
Test user interactions with components:

```go
func TestComponent_Interaction(t *testing.T) {
    app := test.NewApp()
    window := test.NewWindow(component)
    
    // Simulate click
    test.Tap(button)
    
    // Verify result
    if !clicked {
        t.Error("Button should have been clicked")
    }
}
```

### 3. Component State Tests
Test component state changes:

```go
func TestComponent_State(t *testing.T) {
    component.SetState(value)
    if component.GetState() != value {
        t.Error("State should match")
    }
}
```

### 4. Integration Tests
Test complete workflows:

```go
func TestIntegration_Workflow(t *testing.T) {
    // 1. Select directory
    // 2. Start scan
    // 3. View results
    // 4. Generate script
}
```

## Test Utilities

### Fyne Test Helpers

```go
import "fyne.io/fyne/v2/test"

// Tap a button
test.Tap(button)

// Type in an entry
test.Type(entry, "text")

// Drag from one point to another
test.Drag(canvasObject, startX, startY, endX, endY)

// Scroll an object
test.Scroll(canvasObject, dx, dy)
```

### Test App and Window

```go
import (
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/test"
)

// Create test app
myApp := app.New()

// Create test window
window := test.NewWindow(content)
defer window.Close()
```

## Common Test Patterns

### Testing Button Clicks

```go
func TestButton_Click(t *testing.T) {
    clicked := false
    button := widget.NewButton("Click", func() {
        clicked = true
    })
    
    test.Tap(button)
    
    if !clicked {
        t.Error("Button should have been clicked")
    }
}
```

### Testing Text Entry

```go
func TestEntry_Input(t *testing.T) {
    entry := widget.NewEntry()
    
    test.Type(entry, "Hello World")
    
    if entry.Text != "Hello World" {
        t.Errorf("Expected 'Hello World', got '%s'", entry.Text)
    }
}
```

### Testing Selection

```go
func TestSelect_Selection(t *testing.T) {
    select := widget.NewSelect([]string{"A", "B", "C"}, nil)
    
    select.SetSelected("B")
    
    if select.Selected != "B" {
        t.Error("Selection should be B")
    }
}
```

### Testing Lists

```go
func TestList_Selection(t *testing.T) {
    list := widget.NewList(
        func() int { return 10 },
        func() fyne.CanvasObject { return widget.NewLabel("") },
        func(id, item) { },
    )
    
    // Select item
    list.Select(5)
    
    // Verify selection
    if list.Length() != 10 {
        t.Error("List should have 10 items")
    }
}
```

## Testing DupDel Specific Features

### Testing Directory Selection

```go
func TestDirectoryBar_Selection(t *testing.T) {
    app := ui.NewDupDelApp()
    dirBar := ui.NewDirectoryBar(app)
    
    dirBar.SetPath("/test/path")
    
    if dirBar.GetPath() != "/test/path" {
        t.Error("Path should be set")
    }
}
```

### Testing Filter Panel

```go
func TestFilterPanel_Options(t *testing.T) {
    panel := ui.NewFilterPanel()
    
    opts := panel.GetOptions()
    
    if !opts.Recursive {
        t.Error("Recursive should be true")
    }
}
```

### Testing Results View

```go
func TestResultsView_Groups(t *testing.T) {
    app := ui.NewDupDelApp()
    view := ui.NewResultsView(app)
    
    groups := []models.DuplicateGroup{
        // Add test groups
    }
    
    view.SetGroups(groups)
    
    // Verify groups are displayed
}
```

### Testing Preview View

```go
func TestPreviewView_Preview(t *testing.T) {
    app := ui.NewDupDelApp()
    view := ui.NewPreviewView(app)
    
    group := models.DuplicateGroup{
        // Add test group
    }
    
    view.SetGroup(group)
    
    // Verify preview is shown
}
```

### Testing Script Generation

```go
func TestScriptView_Generation(t *testing.T) {
    app := ui.NewDupDelApp()
    view := ui.NewScriptView(app)
    
    // Set up session with groups
    session := models.NewScanSession("/test")
    
    view.SetSession(session)
    view.Generate()
    
    // Verify script is generated
}
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v2
    
    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.21
    
    - name: Install dependencies
      run: |
        sudo apt-get update
        sudo apt-get install -y libgl1-mesa-dev xorg-dev
    
    - name: Run tests
      run: go test ./... -v -cover
```

### GitLab CI Example

```yaml
test:
  image: golang:1.21
  before_script:
    - apt-get update && apt-get install -y libgl1-mesa-dev xorg-dev
  script:
    - go test ./... -v -cover
```

## Debugging Test Failures

### Enable Test Logging

```bash
go test ./tests/ui/... -v -args -test.v
```

### Capture Screenshots

```go
func TestWithScreenshot(t *testing.T) {
    window := test.NewWindow(content)
    
    // ... test code ...
    
    // Capture screenshot on failure
    if t.Failed() {
        screenshot := window.Capture()
        // Save or analyze screenshot
    }
}
```

### Check Canvas State

```go
func TestCanvasState(t *testing.T) {
    window := test.NewWindow(content)
    
    // Check object count
    if len(window.Canvas().Content().(*fyne.Container).Objects) != 3 {
        t.Error("Should have 3 objects")
    }
}
```

## Best Practices

1. **Use test.App.New()** for test apps, not app.New()
2. **Always close windows** with `defer window.Close()`
3. **Use test.Tap()** instead of calling OnTapped directly
4. **Wait for async operations** with time.Sleep or channels
5. **Clean up resources** after each test
6. **Test edge cases** (empty lists, invalid input, etc.)
7. **Use meaningful test names** that describe what's being tested

## Common Issues

### Issue: "Package gl was not found"
**Solution:** Install OpenGL development libraries
```bash
sudo apt-get install libgl1-mesa-dev
```

### Issue: "X11/Xlib.h: No such file or directory"
**Solution:** Install X11 development libraries
```bash
sudo apt-get install xorg-dev
```

### Issue: Tests hang indefinitely
**Solution:** Add timeout to test
```go
func TestWithTimeout(t *testing.T) {
    done := make(chan bool)
    
    go func() {
        // Test code
        done <- true
    }()
    
    select {
    case <-done:
        // Success
    case <-time.After(5 * time.Second):
        t.Fatal("Test timed out")
    }
}
```

### Issue: Canvas is nil in tests
**Solution:** Use test.NewWindow() which creates a proper test canvas
```go
window := test.NewWindow(content)
```

## Additional Resources

- [Fyne Testing Documentation](https://fyne.io/docs/testing/)
- [Fyne Test Package](https://pkg.go.dev/fyne.io/fyne/v2/test)
- [Go Testing Package](https://pkg.go.dev/testing)
