# DupDel Implementation Prompts

This document contains detailed prompts for implementing each task in the implementation plan.

---

## Phase 1: Core Functionality (Weeks 1-3)

### Task 1.1: Project Initialization

```
Create a Go project structure for a duplicate file scanner application called "DupDel".

Requirements:
1. Initialize Go module with name "github.com/dupdel/dup-del" and Go 1.21+ requirement
2. Create the following directory structure:
   - cmd/dup-del/ (main entry point)
   - internal/app/ (application logic)
   - internal/scanner/ (file scanning)
   - internal/hasher/ (hash calculation)
   - internal/detector/ (duplicate detection)
   - internal/preview/ (file preview)
   - internal/rules/ (rule engine)
   - internal/script/ (script generation)
   - internal/ui/ (GUI components)
   - pkg/models/ (data models)
   - pkg/types/ (type definitions)
   - assets/icons/ (application icons)
   - tests/ (test files)

3. Create a .gitignore file for Go projects (include: binaries, dependencies, OS files)
4. Create a LICENSE file with MIT license
5. Create an initial main.go in cmd/dup-del/ that prints "DupDel v1.0.0" and exits

Output: Complete project structure with all directories and initial files.
```

---

### Task 1.2: Define Core Data Models

```
Create the core data models for the DupDel application in the pkg/models/ directory.

Files to create:

1. pkg/types/types.go:
   - Define FileType enum (FileTypeUnknown, FileTypeText, FileTypeImage, FileTypeVideo, FileTypeOther)
   - Define file extension lists for each type
   - Create GetFileTypeFromExtension() function
   - Create helper functions: SupportedTextExtensions(), SupportedImageExtensions(), SupportedVideoExtensions()

2. pkg/models/file.go:
   - Create FileEntry struct with fields:
     * Path (string) - absolute file path
     * Name (string) - file name with extension
     * Size (int64) - file size in bytes
     * Extension (string) - lowercase extension
     * ModTime (time.Time) - last modified time
     * CreatedTime (time.Time) - creation time
     * Hash (string) - SHA-256 hash
     * FileType (types.FileType) - type of file
     * IsHidden (bool) - whether file is hidden
     * IsSymlink (bool) - whether file is a symbolic link
   - Add constructor function NewFileEntry(path string) (*FileEntry, error)
   - Add method GetDisplaySize() string that returns human-readable size

3. pkg/models/group.go:
   - Create DuplicateGroup struct with fields:
     * ID (string) - unique group identifier
     * Hash (string) - common hash value
     * Size (int64) - file size
     * Extension (string) - file extension
     * FileType (types.FileType) - type of files
     * Files ([]FileEntry) - all duplicate files
     * KeepIndices ([]int) - indices of files to keep
     * DeletePaths ([]string) - paths marked for deletion
     * RuleApplied (*Rule) - rule applied to this group
   - Add method GetRecoverableSize() int64
   - Add method GetFileCount() int

4. pkg/models/rule.go:
   - Define RuleCriteria type as string enum
   - Define constants: KeepNewest, KeepOldest, KeepShortestPath, KeepLongestPath, KeepSpecificDir, ManualSelection
   - Create Rule struct with fields:
     * ID (string)
     * Name (string)
     * Description (string)
     * Criteria (RuleCriteria)
     * ApplyToAll (bool)
     * Config (map[string]interface{})
   - Create FileDecision struct with fields:
     * FilePath (string)
     * Decision (string) - "keep" or "delete"

5. pkg/models/session.go:
   - Create ScanSession struct with fields:
     * ID (string)
     * SourceDir (string)
     * StartTime (time.Time)
     * EndTime (time.Time)
     * TotalFiles (int)
     * TotalSize (int64)
     * DuplicateGroups ([]DuplicateGroup)
     * FilesToDelete (int)
     * SpaceRecoverable (int64)
     * GeneratedScript (string)
   - Add method GetDuration() time.Duration
   - Add method GetRecoverableSizeHuman() string

Output: All model files with proper Go code, constructors, and methods.
```

---

### Task 1.3: File System Scanner

```
Create the file system scanner component for DupDel.

Files to create:

1. internal/scanner/scanner.go:
   - Create ScanOptions struct with fields:
     * Recursive (bool)
     * IncludeHidden (bool)
     * FollowSymlinks (bool)
     * MinSize (int64)
     * MaxSize (int64)
     * Extensions ([]string)
   - Create ScanProgress struct with fields:
     * FilesScanned (int64)
     * TotalFiles (int64)
     * CurrentFile (string)
     * TotalSize (int64)
     * PercentComplete (float64)
   - Create Scanner struct with fields:
     * rootDir (string)
     * options (ScanOptions)
     * progress (*ScanProgress)
     * ctx (context.Context)
     * cancel (context.CancelFunc)
   - Implement methods:
     * NewScanner(rootDir string, options ScanOptions) *Scanner
     * Scan() (<-chan FileEntry, <-chan error, error) - returns channels for files and errors
     * CountFiles() (int64, error) - pre-scan to count files
     * GetProgress() ScanProgress
     * Cancel()
   - Use filepath.WalkDir for directory traversal
   - Implement context-based cancellation
   - Handle permission errors gracefully (log and continue)

2. internal/scanner/filter.go:
   - Create Filter struct with static methods:
     * ShouldIncludeHidden(name string) bool
     * ShouldIncludeBySize(size, minSize, maxSize int64) bool
     * ShouldIncludeByExtension(ext string, allowed []string) bool
     * IsHidden(name string) bool
     * NormalizeExtension(ext string) string

3. internal/scanner/progress.go:
   - Create ProgressTracker struct
   - Implement methods for tracking and reporting progress
   - Add callback support for progress updates

Output: Complete scanner implementation with tests.
```

---

### Task 1.4: Scanner Unit Tests

```
Create comprehensive unit tests for the scanner component.

File to create: tests/scanner/scanner_test.go

Test cases to implement:
1. TestScanner_EmptyDirectory - scan empty directory, expect no files
2. TestScanner_SingleFile - scan directory with single file
3. TestScanner_NestedDirectories - scan directory with nested subdirectories
4. TestScanner_HiddenFiles - verify hidden files are included/excluded based on option
5. TestScanner_Symlinks - verify symlinks are handled correctly
6. TestScanner_PermissionErrors - verify scanner continues on permission denied
7. TestScanner_SizeFilter_MinSize - verify minimum size filter works
8. TestScanner_SizeFilter_MaxSize - verify maximum size filter works
9. TestScanner_ExtensionFilter - verify extension filtering works
10. TestScanner_Recursive - verify recursive scanning works
11. TestScanner_NonRecursive - verify non-recursive scanning works
12. TestScanner_Cancel - verify cancellation works mid-scan
13. TestScanner_Progress - verify progress tracking is accurate
14. TestFilter_ShouldIncludeHidden - test hidden file detection
15. TestFilter_ShouldIncludeBySize - test size filtering
16. TestFilter_ShouldIncludeByExtension - test extension filtering

Use os.MkdirTemp for creating test directories.
Create helper functions for setting up test fixtures.

Output: Complete test file with all test cases passing.
```

---

### Task 1.5: Hash Calculator

```
Create the hash calculation component for DupDel.

Files to create:

1. internal/hasher/hasher.go:
   - Create Hasher struct with fields:
     * algorithm (string) - default "sha256"
     * chunkSize (int64) - default 4MB for large files
     * largeFileThreshold (int64) - default 100MB
   - Create HashResult struct with fields:
     * Path (string)
     * Hash (string)
     * Size (int64)
     * Error (error)
     * Cancelled (bool)
   - Implement methods:
     * NewHasher() *Hasher
     * CalculateHash(filePath string) (string, error)
     * CalculateHashWithProgress(filePath string, progressChan chan<- float64) (string, error)
     * HashMultiple(files []string, maxWorkers int) <-chan HashResult
   - Use crypto/sha256 for hashing
   - Implement chunked reading for large files
   - Handle file read errors gracefully

2. internal/hasher/pool.go:
   - Create WorkerPool struct for parallel hashing
   - Implement worker goroutines
   - Manage work distribution and result collection
   - Support cancellation

3. internal/hasher/cache.go:
   - Create HashCache struct for caching computed hashes
   - Implement in-memory cache with LRU eviction
   - Add persistence to disk (optional, JSON file)
   - Methods: Get(path string) string, Set(path, hash string), Clear()

Output: Complete hasher implementation with worker pool support.
```

---

### Task 1.6: Duplicate Detector

```
Create the duplicate detection component for DupDel.

Files to create:

1. internal/detector/detector.go:
   - Create Detector struct with fields:
     * scanner (*scanner.Scanner)
     * hasher (*hasher.Hasher)
     * sizeMap (map[int64][]models.FileEntry)
     * hashMap (map[string][]models.FileEntry)
     * groups ([]models.DuplicateGroup)
     * ctx (context.Context)
     * cancel (context.CancelFunc)
     * stats (*DetectionStats)
   - Create DetectionStats struct with fields:
     * TotalFiles (int)
     * TotalSize (int64)
     * DuplicateFiles (int)
     * DuplicateGroups (int)
     * RecoverableSize (int64)
     * ScanDuration (time.Duration)
   - Implement methods:
     * NewDetector(scanDir string, options scanner.ScanOptions) *Detector
     * Detect() error - main detection algorithm
     * GetGroups() []models.DuplicateGroup
     * GetStatistics() DetectionStats
     * Cancel()
   - Algorithm:
     1. Scan all files
     2. Group by size (files with unique sizes can't be duplicates)
     3. For size groups with 2+ files, calculate hashes
     4. Group by hash
     5. Create DuplicateGroup for hash groups with 2+ files
     6. Calculate statistics

2. internal/detector/grouper.go:
   - Create Grouper helper struct
   - Implement grouping logic for size and hash
   - Generate unique group IDs

3. internal/detector/stats.go:
   - Implement statistics calculation
   - Calculate recoverable space (sum of duplicate file sizes minus one original)

Output: Complete detector implementation with efficient algorithms.
```

---

### Task 1.7: Detection Tests & Benchmarking

```
Create tests and benchmarks for the detector component.

Files to create:

1. tests/detector/detector_test.go:
   - TestDetector_NoDuplicates - directory with unique files
   - TestDetector_AllDuplicates - directory where all files are duplicates
   - TestDetector_MixedDuplicates - directory with some duplicates
   - TestDetector_DifferentSizes - files of different sizes
   - TestDetector_SameSize_DifferentContent - same size but different content
   - TestDetector_EmptyDirectory - empty directory
   - TestDetector_LargeFiles - test with large files (>10MB)
   - TestDetector_MixedFileTypes - test with text, images, etc.
   - TestDetector_Cancel - test cancellation during detection
   - TestDetector_Statistics - verify statistics are accurate

2. tests/detector/benchmark_test.go:
   - BenchmarkDetector_SmallFiles - 100 files, 1KB each
   - BenchmarkDetector_MediumFiles - 1000 files, 100KB each
   - BenchmarkDetector_LargeFiles - 100 files, 10MB each
   - BenchmarkDetector_ManyFiles - 10000 files, mixed sizes
   - BenchmarkDetector_DeepNesting - files in deeply nested directories

Output: Complete test suite with passing tests and benchmark results.
```

---

### Task 1.8: CLI Application

```
Create the CLI entry point for DupDel.

Files to create:

1. cmd/dup-del/main.go:
   - Parse command line flags:
     * --dir, -d string: Directory to scan (default: current)
     * --min-size int: Minimum file size in bytes (default: 0)
     * --max-size int: Maximum file size (default: unlimited)
     * --types string: Comma-separated extensions (default: all)
     * --output, -o string: Output script path
     * --recursive, -r bool: Scan subdirectories (default: true)
     * --ignore-hidden bool: Skip hidden files (default: false)
     * --hash-only bool: Skip preview info (default: false)
     * --verbose, -v bool: Verbose output (default: false)
     * --version: Show version
     * --help, -h: Show help
   - Implement main() function with flag parsing
   - Create and run detector
   - Display progress with text-based progress bar
   - Generate script on completion

2. cmd/dup-del/output.go:
   - Create progress display functions
   - Create results summary display
   - Implement text-based progress bar
   - Format file sizes for display

3. cmd/dup-del/script.go:
   - Handle script generation from CLI
   - Display script location on completion

Output: Working CLI application that can scan directories and generate cleanup scripts.
```

---

### Task 1.9: Script Generator (Core)

```
Create the bash script generator component.

Files to create:

1. internal/script/generator.go:
   - Create Generator struct with fields:
     * outputDir (string)
     * scriptName (string)
     * includeHeader (bool)
     * includeComments (bool)
   - Create ScriptData struct for template data
   - Implement methods:
     * NewGenerator(outputPath string) *Generator
     * Generate(session *models.ScanSession) (string, error)
     * MakeExecutable(path string) error
     * Validate(path string) error

2. internal/script/template.go:
   - Define bash script template as constant
   - Include:
     * Shebang line
     * Header comments with metadata
     * DRY_RUN support
     * delete_file function
     * Grouped rm commands with comments
     * Summary at end

3. internal/script/validator.go:
   - Validate generated script paths exist
   - Check for potential issues
   - Return warnings and errors

Output: Script generator that produces valid, safe bash scripts.
```

---

### Task 1.10: Phase 1 Integration & Testing

```
Perform integration testing for Phase 1 components.

Tasks:
1. Create end-to-end test that:
   - Creates test directory with known duplicates
   - Runs scanner
   - Runs hasher
   - Runs detector
   - Generates script
   - Verifies script content

2. Performance test:
   - Create 1000 test files
   - Measure scan time
   - Measure hash time
   - Measure total detection time
   - Verify targets (>1000 files/second)

3. Memory profiling:
   - Run detection on large directory
   - Profile memory usage
   - Identify leaks

4. Bug fixes based on test results

Output: Phase 1 complete with all tests passing.
```

---

## Phase 2: GUI Foundation (Weeks 4-5)

### Task 2.1: Fyne Integration

```
Set up Fyne GUI framework for DupDel.

Tasks:
1. Add dependencies to go.mod:
   - fyne.io/fyne/v2 v2.4.0
   - golang.org/x/image
   - github.com/disintegration/imaging

2. Run: go mod tidy

3. Create basic Fyne app in internal/ui/app.go:
   - Initialize fyne app
   - Create window
   - Set up basic event loop

4. Test on current platform

Output: Fyne dependencies installed and basic app running.
```

---

### Task 2.2: Main Window Layout

```
Create the main window layout for DupDel GUI.

Files to create:

1. internal/ui/main_window.go:
   - Create MainWindow struct with:
     * window (fyne.Window)
     * app (fyne.App)
     * content (*MainContent)
     * statusBar (*StatusBar)
   - Implement:
     * NewMainWindow(app fyne.App) *MainWindow
     * Show()
     * SetupContent()
     * SetupMenu()
     * SetupToolbar()

2. internal/ui/menu.go:
   - Create MainMenu struct
   - Implement File menu (Open, Save, Quit)
   - Implement Edit menu
   - Implement View menu
   - Implement Help menu

3. internal/ui/toolbar.go:
   - Create Toolbar struct
   - Add buttons: Open, Scan, Stop, Save, Settings, Help
   - Implement button callbacks

4. internal/ui/status_bar.go:
   - Create StatusBar struct
   - Add status label
   - Add progress indicator
   - Add statistics label
   - Implement update methods

Output: Main window with menu, toolbar, and status bar.
```

---

### Task 2.3: Directory Selection Dialog

```
Create directory selection dialog for DupDel.

Files to create:

1. internal/ui/widgets/dir_dialog.go:
   - Create DirectoryDialog struct
   - Implement file folder dialog using fyne dialog
   - Add recent directories list
   - Implement onSelect and onCancel callbacks

2. internal/ui/widgets/recent_dirs.go:
   - Track recently selected directories
   - Persist to config file
   - Display in directory dialog

Output: Working directory selection dialog.
```

---

### Task 2.4: Filter Panel

```
Create filter panel for scan options.

Files to create:

1. internal/ui/filter_panel.go:
   - Create FilterPanel struct with:
     * recursiveCheck (*widget.Check)
     * hiddenCheck (*widget.Check)
     * minSizeSelect (*widget.Select)
     * maxSizeSelect (*widget.Select)
     * fileTypesButton (*widget.Button)
   - Implement FilterOptions struct
   - Add onFilterChange callback
   - Implement GetOptions() method

2. internal/ui/file_type_dialog.go:
   - Create dialog for selecting file types
   - Checkboxes for: Images, Videos, Documents, Text, Archives, Other
   - Custom extension input

Output: Filter panel with all scan options.
```

---

### Task 2.5: Status Bar

```
Create status bar component.

File to enhance: internal/ui/status_bar.go

Implement:
- SetStatus(text string) - update status text
- SetProgress(value float32) - update progress bar (0.0-1.0)
- SetStatistics(files, duplicates int, size int64) - update stats
- SetMode(mode string) - change status bar mode

Output: Fully functional status bar.
```

---

### Task 2.6: Scan View

```
Create scan progress view.

Files to create:

1. internal/ui/scan_view.go:
   - Create ScanView struct with:
     * pathLabel (*widget.Label)
     * progressBar (*widget.ProgressBar)
     * progressText (*widget.Label)
     * filesLabel (*widget.Label)
     * currentFile (*widget.Label)
     * statsBox (*fyne.Container)
     * stopButton (*widget.Button)
   - Implement:
     * NewScanView() *ScanView
     * StartScan(path string, options ScanOptions)
     * StopScan()
     * OnProgress(update ScanProgress)

2. internal/ui/scan_progress.go:
   - Create detailed progress widget
   - Show files per second
   - Show ETA

Output: Scan view with real-time progress updates.
```

---

### Task 2.7: Results View (List)

```
Create results view showing duplicate groups.

Files to create:

1. internal/ui/results_view.go:
   - Create ResultsView struct with:
     * list (*widget.List)
     * groups ([]DuplicateGroup)
     * selectedGroup (int)
     * searchEntry (*widget.Entry)
     * sortSelect (*widget.Select)
   - Implement:
     * NewResultsView() *ResultsView
     * SetGroups(groups []DuplicateGroup)
     * GetSelectedGroup() *DuplicateGroup
     * Filter(query string)
     * Sort(by SortOption)

2. internal/ui/widgets/group_card.go:
   - Create GroupCard widget
   - Display group summary
   - Show file list
   - Add selection checkboxes

Output: Results view with sortable, filterable list.
```

---

### Task 2.8: Results View Integration

```
Integrate results view with detector and other components.

Tasks:
1. Connect results view to detector output
2. Implement group selection callback
3. Add bulk selection actions (Select All, Deselect All)
4. Implement search/filter functionality
5. Add sort options (by size, count, name)
6. Add "Generate Script" button

Output: Fully integrated results view.
```

---

### Task 2.9: Phase 2 Integration & Testing

```
Perform Phase 2 integration testing.

Tasks:
1. Test full GUI flow: select directory → scan → view results
2. Test UI responsiveness during scan
3. Test on all target platforms
4. Fix any bugs found

Output: Phase 2 complete with working GUI foundation.
```

---

## Phase 3: Preview System (Weeks 6-7)

### Task 3.1: Preview View Container

```
Create main preview view container.

Files to create:

1. internal/ui/preview_view.go:
   - Create PreviewView struct with:
     * container (*fyne.Container)
     * groupLabel (*widget.Label)
     * groupIndex (int)
     * totalGroups (int)
     * prevButton (*widget.Button)
     * nextButton (*widget.Button)
     * fileCard1 (*FileCard)
     * fileCard2 (*FileCard)
     * actionsBox (*fyne.Container)
   - Implement:
     * NewPreviewView() *PreviewView
     * SetGroup(group *DuplicateGroup)
     * NextGroup()
     * PreviousGroup()
     * GetDecisions() GroupDecisions

Output: Preview view container with navigation.
```

---

### Task 3.2: File Card Widget

```
Create file card widget for preview.

Files to create:

1. internal/ui/widgets/file_card.go:
   - Create FileCard struct with:
     * pathLabel (*widget.Label)
     * previewArea (*fyne.Container)
     * sizeLabel (*widget.Label)
     * modifiedLabel (*widget.Label)
     * createdLabel (*widget.Label)
     * hashLabel (*widget.Label)
     * keepRadio (*widget.RadioGroup)
   - Implement:
     * NewFileCard() *FileCard
     * SetFile(entry FileEntry)
     * SetPreview(content fyne.CanvasObject)
     * GetDecision() FileDecision
     * SetDecision(decision FileDecision)

Output: Reusable file card widget.
```

---

### Task 3.3: Text Preview Provider

```
Create text file preview provider.

Files to create:

1. internal/preview/text.go:
   - Create TextPreviewProvider struct
   - Implement Preview(file FileEntry) method
   - Read first N lines of text file
   - Detect character encoding
   - Detect language for syntax highlighting

2. internal/preview/text_highlight.go:
   - Add syntax highlighting for common languages
   - Support: Go, Python, JavaScript, Java, C, etc.

Output: Text preview with syntax highlighting.
```

---

### Task 3.4: Image Preview Provider

```
Create image file preview provider.

Files to create:

1. internal/preview/image.go:
   - Create ImagePreviewProvider struct
   - Implement Preview(file FileEntry) method
   - Load and display images
   - Add zoom controls (25% - 400%)
   - Add pan support
   - Generate thumbnails
   - Extract EXIF metadata

Output: Image preview with zoom/pan.
```

---

### Task 3.5: Preview Provider Interface

```
Create preview provider interface and registry.

Files to create:

1. internal/preview/provider.go:
   - Define PreviewProvider interface
   - Create PreviewRegistry struct
   - Implement provider registration
   - Implement GetProvider(file FileEntry)
   - Add fallback provider for unsupported types

Output: Provider registry with multiple providers.
```

---

### Task 3.6: Video Preview Provider

```
Create video file preview provider.

Files to create:

1. internal/preview/video.go:
   - Create VideoPreviewProvider struct
   - Implement Preview(file FileEntry) method
   - Extract thumbnail using:
     1. FFmpeg if available
     2. Embedded thumbnail
     3. Fallback to icon
   - Extract video metadata (duration, codec, resolution)

Output: Video preview with thumbnail and metadata.
```

---

### Task 3.7: Side-by-Side Comparison

```
Create side-by-side comparison view.

Files to create:

1. internal/ui/widgets/split_preview.go:
   - Create SplitPreview struct
   - Implement split container
   - Add synchronized zoom for images
   - Add diff view for text files
   - Add zoom slider

Output: Side-by-side comparison with sync features.
```

---

### Task 3.8: Action Buttons & Decisions

```
Create action buttons for preview view.

Files to create:

1. internal/ui/preview_actions.go:
   - Create PreviewActions struct
   - Add buttons: Keep This, Delete This, Keep Both, Auto-Select
   - Add "Apply to All" checkbox
   - Add Save & Continue, Skip buttons
   - Implement decision tracking

Output: Action buttons with decision tracking.
```

---

### Task 3.9: Phase 3 Integration & Testing

```
Perform Phase 3 integration testing.

Tasks:
1. Test preview for all file types
2. Test side-by-side comparison
3. Test decision tracking
4. Performance test with large files
5. Fix bugs

Output: Phase 3 complete with working preview system.
```

---

## Phase 4: Rule System (Week 8)

### Task 4.1: Rule Definitions

```
Create built-in rule definitions.

Files to create:

1. internal/rules/builtin.go:
   - Define BuiltinRules slice with:
     * Keep Newest
     * Keep Oldest
     * Keep Shortest Path
     * Keep Longest Path
     * Keep from Specific Directory
     * Manual Selection
   - Each rule has ID, Name, Description, Criteria, Icon

Output: Built-in rules defined.
```

---

### Task 4.2: Rule Engine

```
Create rule application engine.

Files to create:

1. internal/rules/engine.go:
   - Create RuleEngine struct
   - Implement ApplyRule(group, ruleID)
   - Implement ApplyToAll(groups, ruleID)
   - Register evaluators for each criteria type

2. internal/rules/evaluators.go:
   - Create evaluator for each RuleCriteria
   - Implement Evaluate method for each

Output: Rule engine with all evaluators.
```

---

### Task 4.3: Rule Selection Dialog

```
Create rule selection dialog.

Files to create:

1. internal/ui/rule_dialog.go:
   - Create RuleDialog struct
   - Display list of available rules
   - Show rule preview
   - Add "Apply to All" checkbox
   - Implement onSelect callback

Output: Rule selection dialog.
```

---

### Task 4.4: Apply to All Functionality

```
Implement batch rule application.

Files to create:

1. internal/rules/batch.go:
   - Create BatchRuleApplier struct
   - Implement Apply(ruleID, groups)
   - Track results and skipped groups
   - Implement Undo()

Output: Batch rule application.
```

---

### Task 4.5: Decision Review Panel

```
Create decision review panel.

Files to create:

1. internal/ui/review_panel.go:
   - Create ReviewPanel struct
   - Display all decisions in list
   - Add filter by decision type
   - Add edit decision functionality

Output: Decision review panel.
```

---

### Task 4.6: Phase 4 Integration & Testing

```
Perform Phase 4 integration testing.

Tasks:
1. Test all built-in rules
2. Test batch application
3. Test decision review
4. Fix bugs

Output: Phase 4 complete with working rule system.
```

---

## Phase 5: Script Generation (Week 9)

### Task 5.1: Enhanced Script Generator

```
Enhance script generator with more features.

Enhance: internal/script/generator.go

Add:
- Dry-run support in script
- Better formatting
- Group comments
- Statistics in header
- Safety checks

Output: Enhanced script generator.
```

---

### Task 5.2: Script Preview View

```
Create script preview view.

Files to create:

1. internal/ui/script_view.go:
   - Create ScriptView struct
   - Display generated script with syntax highlighting
   - Add line numbers
   - Add save button
   - Add copy to clipboard
   - Add options (executable, backup, dry-run)

Output: Script preview view.
```

---

### Task 5.3: Export Options

```
Create export functionality.

Files to create:

1. internal/export/csv.go:
   - Create CSVExporter struct
   - Implement Export(session, path)

2. internal/export/json.go:
   - Create JSONExporter struct
   - Implement Export(session, path)

3. internal/export/html.go (optional):
   - Create HTMLExporter struct
   - Implement Export(session, path)

Output: Export functionality for multiple formats.
```

---

### Task 5.4: Script Validation

```
Create script validation.

Enhance: internal/script/validator.go

Add:
- Path existence check
- Permission check
- Warning generation

Output: Script validation.
```

---

### Task 5.5: Phase 5 Integration & Testing

```
Perform Phase 5 integration testing.

Tasks:
1. Test script generation
2. Test all export formats
3. Test script validation
4. Fix bugs

Output: Phase 5 complete.
```

---

## Phase 6: Polish & Release (Weeks 10-11)

### Task 6.1: Theme Support

```
Add theme support.

Files to create:

1. internal/ui/themes.go:
   - Create ThemeManager struct
   - Implement dark theme
   - Implement light theme
   - Add theme switcher

Output: Theme support.
```

---

### Task 6.2: Keyboard Shortcuts

```
Add keyboard shortcuts.

Files to create:

1. internal/ui/shortcuts.go:
   - Define all keyboard shortcuts
   - Implement shortcut handler
   - Add help dialog showing shortcuts

Output: Keyboard shortcuts.
```

---

### Task 6.3: Performance Optimization

```
Optimize performance.

Tasks:
1. Profile hash calculation
2. Optimize memory usage
3. Implement lazy loading
4. Add image caching
5. Meet performance targets

Output: Optimized application.
```

---

### Task 6.4: Settings Persistence

```
Add settings persistence.

Files to create:

1. internal/app/config.go:
   - Create Config struct
   - Implement load from YAML/JSON
   - Implement save to file
   - Define config file locations per platform

Output: Settings persistence.
```

---

### Task 6.5: Help System

```
Create help system.

Files to create:

1. internal/ui/help.go:
   - Create help dialog
   - Add FAQ section
   - Add about dialog

Output: Help system.
```

---

### Task 6.6: Comprehensive Testing

```
Perform comprehensive testing.

Tasks:
1. Run all unit tests
2. Run integration tests
3. Test on all platforms
4. Performance benchmark
5. Fix all bugs

Output: Fully tested application.
```

---

### Task 6.7: Documentation

```
Create documentation.

Files to update/create:
1. README.md - final version
2. CONTRIBUTING.md
3. CHANGELOG.md
4. User guide

Output: Complete documentation.
```

---

### Task 6.8: Release Preparation

```
Prepare for release.

Tasks:
1. Create git tag v1.0.0
2. Write release notes
3. Build binaries for all platforms
4. Create distribution packages

Output: Release ready.
```

---

## Usage Instructions

To use these prompts:

1. **Sequential Implementation**: Follow the prompts in order, completing each phase before moving to the next.

2. **Copy-Paste**: Copy each prompt and use it to generate the code for that specific task.

3. **Verification**: After completing each task, run tests to verify the implementation.

4. **Integration**: At the end of each phase, perform integration testing.

5. **Documentation**: Update documentation as features are completed.
