# Screenshot-Based UI Testing

This package provides screenshot-based visual regression testing for DupDel UI components.

## How It Works

1. **Capture** - Takes a screenshot of a UI component
2. **Compare** - Compares against a "golden" (expected) screenshot
3. **Report** - Flags visual differences > 2%

## Directory Structure

```
tests/ui/screenshot/
├── screenshot.go          # Screenshot utilities
├── screenshot_test.go     # Screenshot tests
└── testdata/
    ├── golden/            # Expected screenshots
    │   ├── directory_bar.png
    │   ├── filter_panel.png
    │   └── ...
    └── diff/              # Generated diff images (on failure)
        └── *.diff.png
```

## Running Tests

### First Run (Generate Golden Images)

```bash
go test ./tests/ui/screenshot/... -v -tags gui -update
```

This creates the golden images in `testdata/golden/`.

### Subsequent Runs (Compare Against Golden)

```bash
go test ./tests/ui/screenshot/... -v -tags gui
```

### Update Golden Images

If you've made intentional UI changes:

```bash
go test ./tests/ui/screenshot/... -v -tags gui -update
```

## Test Output

### Passing Test
```
=== RUN   TestScreenshot_DirectoryBar
--- PASS: TestScreenshot_DirectoryBar (0.05s)
```

### Failing Test
```
=== RUN   TestScreenshot_DirectoryBar
    screenshot_test.go:25: Screenshot differs from golden image by 15.32%
        Golden: testdata/golden/directory_bar.png
        Diff saved to: testdata/diff/directory_bar.diff.png
--- FAIL: TestScreenshot_DirectoryBar (0.05s)
```

## Understanding Diff Images

Diff images highlight differences in **red**:
- **Original colors** = Same pixels (no change)
- **Red** = Different pixels (visual change)

## Best Practices

### 1. Stable Test Data
Ensure UI components have consistent data:
```go
// Set fixed values before screenshot
statusBar.SetStatus("Ready")
statusBar.SetProgress(0.5)
```

### 2. Fixed Window Size
Resize windows to consistent size:
```go
window.Resize(fyne.NewSize(800, 600))
```

### 3. Wait for Rendering
Allow time for async rendering:
```go
test.WaitForRendering(window.Canvas())
```

### 4. Update Golden Images Responsibly
Only update golden images when:
- ✅ Intentional UI changes
- ✅ Bug fixes in rendering
- ❌ NOT for test failures due to bugs

## CI/CD Integration

### GitHub Actions

```yaml
name: Visual Tests
on: [push, pull_request]

jobs:
  visual-tests:
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
    
    - name: Run visual tests
      run: go test ./tests/ui/screenshot/... -v -tags gui
    
    - name: Upload diff images
      if: failure()
      uses: actions/upload-artifact@v2
      with:
        name: visual-test-diffs
        path: tests/ui/screenshot/testdata/diff/*.png
```

### GitLab CI

```yaml
visual-tests:
  image: golang:1.21
  before_script:
    - apt-get update && apt-get install -y libgl1-mesa-dev xorg-dev
  script:
    - go test ./tests/ui/screenshot/... -v -tags gui
  artifacts:
    when: on_failure
    paths:
      - tests/ui/screenshot/testdata/diff/*.png
```

## API Reference

### Capture Screenshot

```go
img, err := screenshot.Capture(canvasObject)
```

### Save Image

```go
err := screenshot.SaveImage(img, "path/to/file.png")
```

### Load Image

```go
img, err := screenshot.LoadImage("path/to/file.png")
```

### Compare Images

```go
diffPercent, err := screenshot.CompareImages(img1, img2)
// Returns percentage of different pixels (0.0 - 100.0)
```

### Assert Against Golden

```go
screenshot.AssertMatchesGolden(t, img, "path/to/golden.png", updateFlag)
```

### Create Diff Image

```go
err := screenshot.CreateDiffImage(img1, img2, "path/to/diff.png")
```

## Troubleshooting

### Issue: "Screenshot capture not supported"
**Cause:** Running without GUI backend
**Solution:** Use `-tags gui` flag

### Issue: "Golden image not found"
**Cause:** First run or missing golden images
**Solution:** Run with `-update` flag

### Issue: High diff percentage on minor changes
**Cause:** Anti-aliasing or font rendering differences
**Solution:** Increase tolerance in `AssertMatchesGolden` (default 2%)

### Issue: Tests fail on different platforms
**Cause:** Platform-specific rendering differences
**Solution:** Maintain separate golden images per platform:
```go
goldenPath := screenshot.GoldenPath("component_" + runtime.GOOS)
```

## Advanced Usage

### Custom Tolerance

```go
func TestScreenshot_CustomTolerance(t *testing.T) {
    img, _ := screenshot.Capture(component)
    golden, _ := screenshot.LoadImage("golden.png")
    
    diff, _ := screenshot.CompareImages(img, golden)
    
    // Allow 5% difference
    if diff > 5.0 {
        t.Errorf("Diff too high: %.2f%%", diff)
    }
}
```

### Region-Specific Comparison

```go
func TestScreenshot_Region(t *testing.T) {
    img, _ := screenshot.Capture(component)
    
    // Compare only specific region
    region := image.Rect(100, 100, 200, 200)
    cropped := img.(interface {
        SubImage(image.Rectangle) image.Image
    }).SubImage(region)
    
    // Compare cropped region
    // ...
}
```

### Animated GIF Diffs

```go
func CreateAnimatedDiff(img1, img2, path string) error {
    // Create animated GIF showing before/after
    // Implementation left as exercise
}
```

## Performance Tips

1. **Cache golden images** - Load once, compare multiple times
2. **Parallel tests** - Use `t.Parallel()` for faster execution
3. **Selective testing** - Test only changed components in CI
4. **Compress images** - Use optimized PNG encoding

## Security Considerations

- Golden images are committed to repository
- Ensure no sensitive data in screenshots
- Review diff images before committing updates
