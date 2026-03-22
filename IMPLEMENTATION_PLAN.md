# DupDel - Detailed Implementation Plan

## Overview

This document provides a comprehensive, task-level implementation plan based on the product specification. The project is divided into 6 phases spanning 11 weeks.

---

## Phase 1: Core Functionality (Weeks 1-3)

### Week 1: Project Setup & File System Scanner

#### Task 1.1: Project Initialization
**Duration**: 0.5 days  
**Dependencies**: None  
**Deliverables**:
- [ ] Initialize Go module (`go mod init github.com/yourusername/dup-del`)
- [ ] Create directory structure per spec
- [ ] Set up `.gitignore` for Go projects
- [ ] Create initial `go.mod` with Go 1.21+ requirement
- [ ] Add LICENSE file (MIT)
- [ ] Create initial README.md placeholder

**Files to Create**:
```
dup-del/
├── cmd/dup-del/main.go
├── internal/
├── pkg/
├── assets/
├── tests/
├── go.mod
├── .gitignore
└── LICENSE
```

---

#### Task 1.2: Define Core Data Models
**Duration**: 1 day  
**Dependencies**: Task 1.1  
**Deliverables**:
- [ ] `pkg/models/file.go` - FileEntry struct
- [ ] `pkg/models/group.go` - DuplicateGroup struct
- [ ] `pkg/models/rule.go` - Rule and RuleCriteria
- [ ] `pkg/models/session.go` - ScanSession struct
- [ ] `pkg/types/types.go` - FileType enum, constants
- [ ] Unit tests for model constructors

**Implementation Details**:
```go
// pkg/types/types.go
type FileType int

const (
    FileTypeUnknown FileType = iota
    FileTypeText
    FileTypeImage
    FileTypeVideo
    FileTypeOther
)

func GetFileTypeFromExtension(ext string) FileType
func SupportedTextExtensions() []string
func SupportedImageExtensions() []string
func SupportedVideoExtensions() []string
```

---

#### Task 1.3: File System Scanner
**Duration**: 2 days  
**Dependencies**: Task 1.2  
**Deliverables**:
- [ ] `internal/scanner/scanner.go` - Directory traversal
- [ ] `internal/scanner/filter.go` - File filtering logic
- [ ] `internal/scanner/progress.go` - Progress tracking
- [ ] Unit tests for scanner with mock filesystem

**Implementation Details**:
```go
// internal/scanner/scanner.go
type Scanner struct {
    rootDir     string
    options     ScanOptions
    progress    *Progress
    fileChan    chan FileEntry
    errorChan   chan error
    doneChan    chan struct{}
}

type ScanOptions struct {
    Recursive      bool
    IncludeHidden  bool
    FollowSymlinks bool
    MinSize        int64
    MaxSize        int64
    Extensions     []string
}

func (s *Scanner) Scan(ctx context.Context) (<-chan FileEntry, error)
func (s *Scanner) CountFiles() (int64, error)
```

**Key Features**:
- Context-based cancellation
- Concurrent file walking
- Progress callbacks
- Error handling for permission issues
- Size and extension filtering

---

#### Task 1.4: Scanner Unit Tests
**Duration**: 0.5 days
**Dependencies**: Task 1.3
**Deliverables**:
- [ ] `tests/scanner/scanner_test.go`
- [ ] Test with empty directory
- [ ] Test with nested directories
- [ ] Test with hidden files
- [ ] Test with symlinks
- [ ] Test with permission errors
- [ ] Test size filtering
- [ ] Test extension filtering
- [ ] Test context cancellation

---

### Week 2: Hash Calculation & Duplicate Detection

#### Task 1.5: Hash Calculator
**Duration**: 1.5 days  
**Dependencies**: Task 1.2  
**Deliverables**:
- [ ] `internal/hasher/hasher.go` - SHA-256 calculation
- [ ] `internal/hasher/pool.go` - Worker pool for parallel hashing
- [ ] `internal/hasher/cache.go` - Hash caching for resume capability
- [ ] Unit tests for hash accuracy

**Implementation Details**:
```go
// internal/hasher/hasher.go
type Hasher struct {
    algorithm   string
    chunkSize   int64
    workerCount int
}

type HashResult struct {
    Path      string
    Hash      string
    Size      int64
    Error     error
    Cancelled bool
}

func (h *Hasher) CalculateHash(filePath string) (string, error)
func (h *Hasher) CalculateHashWithProgress(filePath string, callback func(percent float64)) (string, error)
func (h *Hasher) HashMultiple(files []string, maxWorkers int) <-chan HashResult
```

**Key Features**:
- SHA-256 as default algorithm
- Chunked reading for large files (>100MB)
- Worker pool for parallel processing
- Context-based cancellation
- Hash verification option

---

#### Task 1.6: Duplicate Detector
**Duration**: 2 days  
**Dependencies**: Task 1.3, Task 1.5  
**Deliverables**:
- [ ] `internal/detector/detector.go` - Main detection logic
- [ ] `internal/detector/grouper.go` - Group duplicates together
- [ ] `internal/detector/stats.go` - Statistics calculation
- [ ] Integration tests

**Implementation Details**:
```go
// internal/detector/detector.go
type Detector struct {
    scanner     *scanner.Scanner
    hasher      *hasher.Hasher
    sizeMap     map[int64][]FileEntry    // Files grouped by size
    hashMap     map[string][]FileEntry   // Files grouped by hash
    groups      []DuplicateGroup
}

func (d *Detector) Detect(ctx context.Context) error
func (d *Detector) GetGroups() []DuplicateGroup
func (d *Detector) GetStatistics() DetectionStats
func (d *Detector) Cancel()
```

**Algorithm**:
1. Scan all files, collect FileEntry objects
2. Group files by size (quick filter - files with unique sizes can't be duplicates)
3. For size groups with 2+ files, calculate hashes
4. Group files by hash
5. Create DuplicateGroup for hash groups with 2+ files
6. Calculate statistics (total duplicates, recoverable space)

---

#### Task 1.7: Detection Tests & Benchmarking
**Duration**: 0.5 days
**Dependencies**: Task 1.6
**Deliverables**:
- [ ] `tests/detector/detector_test.go`
- [ ] Test with known duplicate files
- [ ] Test with no duplicates
- [ ] Test with all duplicates
- [ ] Test with mixed file types
- [ ] Test all files included in groups (no truncation)
- [ ] Benchmark detection speed
- [ ] Memory usage profiling

---

### Week 3: CLI Interface & Core Completion

#### Task 1.8: CLI Application
**Duration**: 2 days  
**Dependencies**: Task 1.6  
**Deliverables**:
- [ ] `cmd/dup-del/main.go` - CLI entry point
- [ ] `cmd/dup-del/flags.go` - Command line flag parsing
- [ ] `cmd/dup-del/output.go` - Console output formatting
- [ ] `cmd/dup-del/script.go` - CLI script generation
- [ ] Help documentation

**Implementation Details**:
```go
// Command line flags
--dir, -d string       Directory to scan (default: current)
--min-size int         Minimum file size in bytes (default: 0)
--max-size int         Maximum file size in bytes (default: unlimited)
--types string         Comma-separated extensions (default: all)
--output, -o string    Output script path (default: cleanup_duplicates.sh)
--recursive, -r bool   Scan subdirectories (default: true)
--ignore-hidden bool   Skip hidden files (default: false)
--hash-only bool       Skip preview info (default: false)
--verbose, -v bool     Verbose output (default: false)
--version              Show version
--help, -h             Show help
```

**CLI Output Format**:
```
$ dup-del --dir /home/user/docs

DupDel v1.0.0 - Duplicate File Scanner
Scanning: /home/user/docs
Options: Recursive=true, Hidden=false, MinSize=0

Scanning files... ████████████████████ 100% (15,234 files, 4.5 GB)
Calculating hashes... ████████████████████ 100%

Found 45 duplicate groups (127 files, 523.4 MB recoverable)

Top duplicates by size:
  1. video.mp4 (2 copies, 156 MB each) - Recoverable: 156 MB
  2. document.pdf (3 copies, 2.4 MB each) - Recoverable: 4.8 MB
  3. IMG_001.jpg (5 copies, 1.2 MB each) - Recoverable: 4.8 MB

Generating cleanup script: cleanup_duplicates.sh
Done! Review the script before running.
```

---

#### Task 1.9: Script Generator (Core)
**Duration**: 1.5 days  
**Dependencies**: Task 1.6  
**Deliverables**:
- [ ] `internal/script/generator.go` - Bash script generation
- [ ] `internal/script/template.go` - Script templates
- [ ] `internal/script/validator.go` - Script validation
- [ ] Tests for script output

**Implementation Details**:
```go
// internal/script/generator.go
type Generator struct {
    outputDir     string
    scriptName    string
    includeHeader bool
    includeComments bool
}

type ScriptTemplate struct {
    GeneratedTime    time.Time
    SourceDirectory  string
    TotalFiles       int
    TotalSize        int64
    Groups           []DuplicateGroup
    DeletePaths      []string
}

func (g *Generator) Generate(session *ScanSession) (string, error)
func (g *Generator) Validate(scriptPath string) error
func (g *Generator) MakeExecutable(scriptPath string) error
```

**Script Template**:
```bash
#!/bin/bash
# ============================================
# DupDel Cleanup Script
# ============================================
# Generated: {{.GeneratedTime}}
# Source Directory: {{.SourceDirectory}}
# Total files to delete: {{.TotalFiles}}
# Total space to recover: {{.TotalSize}}
# ============================================

set -e  # Exit on error

DRY_RUN=${DRY_RUN:-false}

delete_file() {
    local file="$1"
    if [ "$DRY_RUN" = "true" ]; then
        echo "[DRY RUN] Would delete: $file"
    else
        echo "Deleting: $file"
        rm -f "$file"
    fi
}

# Group {{.Group.ID}}: {{.Group.Files.0.Name}} ({{.Group.Size}})
{{range .Group.DeletePaths}}
delete_file "{{.}}"
{{end}}

# ... more groups

echo "Cleanup complete!"
```

---

#### Task 1.10: Phase 1 Integration & Testing
**Duration**: 0.5 days  
**Dependencies**: Tasks 1.7, 1.8, 1.9  
**Deliverables**:
- [ ] End-to-end CLI test
- [ ] Performance benchmark (100K files target)
- [ ] Memory profiling
- [ ] Bug fixes from integration testing
- [ ] Phase 1 documentation update

**Milestone 1 Complete**: ✅ Core CLI application functional

---

## Phase 2: GUI Foundation (Weeks 4-5)

### Week 4: Fyne Setup & Main Window

#### Task 2.1: Fyne Integration
**Duration**: 0.5 days  
**Dependencies**: Phase 1 complete  
**Deliverables**:
- [ ] Add Fyne dependency (`go get fyne.io/fyne/v2`)
- [ ] Add Fyne dependencies (image, theme)
- [ ] Create basic Fyne app initialization
- [ ] Test on target platforms

**Dependencies to Add**:
```go
// go.mod additions
require (
    fyne.io/fyne/v2 v2.4.0
    golang.org/x/image v0.15.0
    github.com/disintegration/imaging v1.6.2
)
```

---

#### Task 2.2: Main Window Layout
**Duration**: 2 days  
**Dependencies**: Task 2.1  
**Deliverables**:
- [ ] `internal/ui/main_window.go` - Main window structure
- [ ] `internal/ui/menu.go` - Menu bar implementation
- [ ] `internal/ui/toolbar.go` - Toolbar with actions
- [ ] `internal/ui/status_bar.go` - Status bar
- [ ] `internal/ui/theme.go` - Custom theme

**Implementation Details**:
```go
// internal/ui/main_window.go
type MainWindow struct {
    window      fyne.Window
    app         fyne.App
    content     *MainContent
    statusBar   *StatusBar
    menu        *MainMenu
}

type MainContent struct {
    container      *fyne.Container
    tabs           *container.AppTabs
    directoryBar   *DirectoryBar
    filterPanel    *FilterPanel
    scanView       *ScanView
    resultsView    *ResultsView
    previewView    *PreviewView
    scriptView     *ScriptView
}

func NewMainWindow(app fyne.App) *MainWindow
func (m *MainWindow) Show()
func (m *MainWindow) SetupContent()
```

**Layout Structure**:
```
MainWindow
├── MenuBar (File, Edit, View, Tools, Help)
├── Toolbar (Open, Scan, Stop, Save, Settings, Help)
├── MainContent (VBoxLayout)
│   ├── DirectoryBar (HBoxLayout)
│   │   ├── Label "Directory:"
│   │   ├── Entry (path)
│   │   └── Button "Browse..."
│   ├── FilterPanel (HBoxLayout)
│   │   ├── Check "Recursive"
│   │   ├── Check "Hidden"
│   │   ├── Select "Min Size"
│   │   ├── Select "Max Size"
│   │   └── Button "File Types..."
│   └── AppTabs
│       ├── Tab "Scan" → ScanView
│       ├── Tab "Results" → ResultsView
│       ├── Tab "Preview" → PreviewView
│       └── Tab "Script" → ScriptView
└── StatusBar
    ├── Label (status text)
    ├── ProgressIndicator
    └── Label (statistics)
```

---

#### Task 2.3: Directory Selection Dialog
**Duration**: 1 day  
**Dependencies**: Task 2.2  
**Deliverables**:
- [ ] `internal/ui/widgets/dir_dialog.go` - Directory picker
- [ ] `internal/ui/widgets/recent_dirs.go` - Recent directories list
- [ ] Integration with main window

**Implementation Details**:
```go
// internal/ui/widgets/dir_dialog.go
type DirectoryDialog struct {
    dialog      *widget.PopUp
    onSelect    func(path string)
    onCancel    func()
    recentList  *widget.List
    driveList   *widget.List
}

func NewDirectoryDialog(callback func(string)) *DirectoryDialog
func (d *DirectoryDialog) Show()
func (d *DirectoryDialog) Hide()
```

---

#### Task 2.4: Filter Panel
**Duration**: 1 day  
**Dependencies**: Task 2.2  
**Deliverables**:
- [ ] `internal/ui/filter_panel.go` - Filter controls
- [ ] `internal/ui/file_type_dialog.go` - File type selector
- [ ] Validation for filter inputs

**Implementation Details**:
```go
// internal/ui/filter_panel.go
type FilterPanel struct {
    widget.BaseWidget
    recursiveCheck    *widget.Check
    hiddenCheck       *widget.Check
    minSizeSelect     *widget.Select
    maxSizeSelect     *widget.Select
    fileTypesButton   *widget.Button
    onFilterChange    func(FilterOptions)
}

type FilterOptions struct {
    Recursive     bool
    IncludeHidden bool
    MinSize       int64
    MaxSize       int64
    Extensions    []string
}

func NewFilterPanel(callback func(FilterOptions)) *FilterPanel
func (f *FilterPanel) GetOptions() FilterOptions
func (f *FilterPanel) SetEnabled(enabled bool)
```

---

#### Task 2.5: Status Bar
**Duration**: 0.5 days  
**Dependencies**: Task 2.2  
**Deliverables**:
- [ ] `internal/ui/status_bar.go` - Status bar component
- [ ] Progress indicator integration
- [ ] Statistics display

**Implementation Details**:
```go
// internal/ui/status_bar.go
type StatusBar struct {
    widget.BaseWidget
    statusLabel   *widget.Label
    progressBar   *widget.ProgressBar
    statsLabel    *widget.Label
}

func NewStatusBar() *StatusBar
func (s *StatusBar) SetStatus(text string)
func (s *StatusBar) SetProgress(value float32)
func (s *StatusBar) SetStatistics(files, duplicates int, size int64)
```

---

### Week 5: Scan & Results Views

#### Task 2.6: Scan View
**Duration**: 2 days  
**Dependencies**: Task 2.2, Task 1.3  
**Deliverables**:
- [ ] `internal/ui/scan_view.go` - Scan progress display
- [ ] `internal/ui/scan_progress.go` - Detailed progress widget
- [ ] Scanner integration with UI
- [ ] Cancel functionality

**Implementation Details**:
```go
// internal/ui/scan_view.go
type ScanView struct {
    widget.BaseWidget
    pathLabel       *widget.Label
    progressBar     *widget.ProgressBar
    progressText    *widget.Label
    filesLabel      *widget.Label
    currentFile     *widget.Label
    statsBox        *fyne.Container
    stopButton      *widget.Button
    scanner         *scanner.Scanner
    cancelChan      chan struct{}
}

func NewScanView() *ScanView
func (s *ScanView) StartScan(path string, options ScanOptions)
func (s *ScanView) StopScan()
func (s *ScanView) OnProgress(update ScanProgress)
```

**UI Components**:
- Large progress bar (centered)
- Files scanned counter (X / Y)
- Current file being processed
- Statistics box (total size, potential duplicates found)
- Elapsed time display
- Stop button (prominent)

---

#### Task 2.7: Results View (List)
**Duration**: 2 days  
**Dependencies**: Task 2.6, Task 1.6  
**Deliverables**:
- [ ] `internal/ui/results_view.go` - Results container
- [ ] `internal/ui/widgets/group_card.go` - Duplicate group card
- [ ] `internal/ui/widgets/group_list.go` - Scrollable list
- [ ] Selection management
- [ ] Sort functionality

**Implementation Details**:
```go
// internal/ui/results_view.go
type ResultsView struct {
    widget.BaseWidget
    list            *widget.List
    groups          []DuplicateGroup
    selectedGroup   int
    searchEntry     *widget.Entry
    sortSelect      *widget.Select
    viewSelect      *widget.Select
    selectAllBtn    *widget.Button
    deselectAllBtn  *widget.Button
    onGroupSelect   func(int)
    onGenerate      func()
}

func NewResultsView() *ResultsView
func (r *ResultsView) SetGroups(groups []DuplicateGroup)
func (r *ResultsView) GetSelectedGroup() *DuplicateGroup
func (r *ResultsView) Filter(query string)
func (r *ResultsView) Sort(by SortOption)
```

**Group Card Layout**:
```
┌─────────────────────────────────────────────────────────┐
│  📄 document.pdf                              [Expand ▼]│
│     3 copies | 2.4 MB each | Recoverable: 4.8 MB       │
│  ─────────────────────────────────────────────────────  │
│  ☐ /home/user/docs/document.pdf                         │
│  ☐ /home/user/backup/document.pdf                       │
│  ☐ /home/user/old/document.pdf                          │
└─────────────────────────────────────────────────────────┘
```

---

#### Task 2.8: Results View Integration
**Duration**: 1 day  
**Dependencies**: Task 2.7  
**Deliverables**:
- [ ] Connect results view to detector
- [ ] Group selection callback to preview
- [ ] Bulk selection actions
- [ ] Search/filter functionality
- [ ] Sort options (size, count, name)

---

#### Task 2.9: Phase 2 Integration & Testing
**Duration**: 0.5 days  
**Dependencies**: Tasks 2.8  
**Deliverables**:
- [ ] Full GUI flow test (select → scan → view results)
- [ ] UI responsiveness testing
- [ ] Cross-platform GUI testing
- [ ] Bug fixes

**Milestone 2 Complete**: ✅ GUI foundation functional (scan + results)

---

## Phase 3: Preview System (Weeks 6-7)

### Week 6: Text & Image Preview

#### Task 3.1: Preview View Container
**Duration**: 1 day  
**Dependencies**: Phase 2 complete  
**Deliverables**:
- [ ] `internal/ui/preview_view.go` - Main preview container
- [ ] `internal/ui/preview_base.go` - Base preview interface
- [ ] Navigation between groups
- [ ] File metadata display

**Implementation Details**:
```go
// internal/ui/preview_view.go
type PreviewView struct {
    widget.BaseWidget
    container       *fyne.Container
    groupLabel      *widget.Label
    groupIndex      int
    totalGroups     int
    prevButton      *widget.Button
    nextButton      *widget.Button
    fileCard1       *FileCard
    fileCard2       *FileCard
    actionsBox      *fyne.Container
    ruleCheckbox    *widget.Check
    currentGroup    *DuplicateGroup
    previewer       PreviewProvider
}

func NewPreviewView() *PreviewView
func (p *PreviewView) SetGroup(group *DuplicateGroup)
func (p *PreviewView) NextGroup()
func (p *PreviewView) PreviousGroup()
func (p *PreviewView) GetDecisions() GroupDecisions
```

---

#### Task 3.2: File Card Widget
**Duration**: 1 day  
**Dependencies**: Task 3.1  
**Deliverables**:
- [ ] `internal/ui/widgets/file_card.go` - File display card
- [ ] Preview area (dynamic based on file type)
- [ ] Metadata display (path, size, dates)
- [ ] Keep/Delete radio selection

**Implementation Details**:
```go
// internal/ui/widgets/file_card.go
type FileCard struct {
    widget.BaseWidget
    pathLabel       *widget.Label
    previewArea     *fyne.Container
    sizeLabel       *widget.Label
    modifiedLabel   *widget.Label
    createdLabel    *widget.Label
    hashLabel       *widget.Label
    keepRadio       *widget.RadioGroup
    fileEntry       FileEntry
    previewContent  fyne.CanvasObject
}

func NewFileCard() *FileCard
func (f *FileCard) SetFile(entry FileEntry)
func (f *FileCard) SetPreview(content fyne.CanvasObject)
func (f *FileCard) GetDecision() FileDecision
func (f *FileCard) SetDecision(decision FileDecision)
```

---

#### Task 3.3: Text Preview Provider
**Duration**: 1.5 days
**Dependencies**: Task 3.2
**Deliverables**:
- [ ] `internal/preview/text_provider.go` - Text file preview
- [ ] Full content preview (no truncation)
- [ ] Monospace font for code files
- [ ] Word wrap support
- [ ] Scrollable text area

**Implementation Details**:
```go
// internal/preview/text_provider.go
type TextPreviewProvider struct {
    wordWrap    bool
    monospace   bool
}

func (t *TextPreviewProvider) Preview(file FileEntry) (fyne.CanvasObject, error)
func (t *TextPreviewProvider) ReadFile(path string) (string, error)
```

**Features**:
- Full file content preview (no line limit)
- Monospace font for code files
- Word wrap toggle
- Scrollable view for large files
- Copy to clipboard button
- UTF-8 encoding support

---

#### Task 3.4: Image Preview Provider
**Duration**: 1.5 days
**Dependencies**: Task 3.2
**Deliverables**:
- [ ] `internal/preview/image_provider.go` - Image file preview
- [ ] `internal/preview/image_zoom.go` - Zoom/pan controls
- [ ] Support for PNG, JPEG, GIF, BMP, WebP
- [ ] Image metadata (dimensions, file size)

**Implementation Details**:
```go
// internal/preview/image_provider.go
type ImagePreviewProvider struct {
    zoomLevel       float32
    enableZoom      bool
    cache           map[string]image.Image
    maxCacheSize    int
}

func (i *ImagePreviewProvider) Preview(file FileEntry) (fyne.CanvasObject, error)
func (i *ImagePreviewProvider) LoadImage(path string) (image.Image, error)
func (i *ImagePreviewProvider) GetMetadata(path string) (ImageMetadata, error)
```

**Features**:
- Full image preview for PNG, JPEG, GIF, BMP, WebP
- Zoom controls (25% - 400%)
- Fit-to-screen button
- Image dimensions display
- File size display
- Image caching for performance

---

#### Task 3.5: Preview Provider Interface
**Duration**: 0.5 days  
**Dependencies**: Tasks 3.3, 3.4  
**Deliverables**:
- [ ] `internal/preview/provider.go` - Provider interface
- [ ] Provider registry/factory
- [ ] Fallback for unsupported types

**Implementation Details**:
```go
// internal/preview/provider.go
type PreviewProvider interface {
    CanPreview(fileType FileType, extension string) bool
    Preview(file FileEntry) (fyne.CanvasObject, error)
    GetMetadata(file FileEntry) (Metadata, error)
}

type PreviewRegistry struct {
    providers []PreviewProvider
}

func (r *PreviewRegistry) Register(provider PreviewProvider)
func (r *PreviewRegistry) GetProvider(file FileEntry) PreviewProvider
func (r *PreviewRegistry) Preview(file FileEntry) (fyne.CanvasObject, error)
```

---

#### Task 3.5.1: Preview Provider Tests
**Duration**: 0.5 days
**Dependencies**: Task 3.5
**Deliverables**:
- [ ] `tests/preview/text_provider_test.go`
- [ ] `tests/preview/image_provider_test.go`
- [ ] Test text preview (full content, no truncation)
- [ ] Test image preview (PNG, JPEG, GIF, BMP, WebP)
- [ ] Test zoom functionality
- [ ] Test metadata extraction
- [ ] Test error handling for corrupt files
- [ ] Test cache behavior

---

### Week 7: Video Preview & Actions

#### Task 3.6: Video Preview Provider
**Duration**: 2 days
**Dependencies**: Task 3.5
**Deliverables**:
- [x] `internal/preview/video_provider.go` - Video file preview
- [x] `internal/preview/video_metadata.go` - Video metadata structures
- [x] Thumbnail extraction (FFmpeg or placeholder)
- [x] Video metadata (duration, codec, resolution, bitrate, FPS)
- [x] `internal/platform/open_video.go` - Open in default player
- [x] Quality indicator (4K/1080p/720p/480p)

**Implementation Details**:
```go
// internal/preview/video_provider.go
type VideoProvider struct {
    ffmpegPath   string
    ffprobePath  string
    useFFmpeg    bool
    thumbnailCache map[string]interface{}
}

type VideoMetadata struct {
    Duration    time.Duration
    Width       int
    Height      int
    VideoCodec  string
    AudioCodec  string
    Bitrate     int64
    FrameRate   float64
    HasAudio    bool
    HasVideo    bool
    Format      string
}

func (v *VideoProvider) Preview(file FileEntry) (fyne.CanvasObject, error)
func (v *VideoProvider) ExtractThumbnail(path string) (image.Image, error)
func (v *VideoProvider) GetVideoMetadata(path string) (*VideoMetadata, error)
func (v *VideoProvider) OpenInPlayer(path string) error
```

**Thumbnail Extraction Strategy**:
1. Check thumbnail cache first
2. Try FFmpeg if available (`ffmpeg -i video.mp4 -ss 00:00:05 -vframes 1 thumb.jpg`)
3. Try embedded thumbnail extraction (MP4/MKV)
4. Fallback to placeholder with play icon

**Metadata Display**:
- Duration (HH:MM:SS or MM:SS format)
- Resolution with quality label (1920x1080 (1080p))
- Video codec (H.264, HEVC, VP9, etc.)
- Audio codec (AAC, MP3, etc.)
- Bitrate (Mbps/Kbps)
- Frame rate (fps)
- File size
- **Open in Player** button (platform-specific)

**Platform Support**:
- Windows: `start` command
- macOS: `open` command
- Linux: `xdg-open` or fallback to VLC/MPV

---

#### Task 3.7: Side-by-Side Comparison
**Duration**: 1.5 days
**Dependencies**: Task 3.1, Task 3.6
**Deliverables**:
- [x] Split view layout (HSplit container)
- [x] File navigation for groups with >2 files
- [x] Previous/Next file buttons
- [x] Dynamic action buttons based on file count

**Implementation Details**:
```go
// internal/ui/preview_view.go
type PreviewView struct {
    container     *fyne.Container
    fileCard1     *FilePreviewCard  // Always shows first file
    fileCard2     *FilePreviewCard  // Cycles through remaining files
    fileIndex     int               // Current comparison index
    actionsBox    *fyne.Container   // Dynamic action buttons
}

func (pv *PreviewView) updateFileCards()
func (pv *PreviewView) updateActionButtons()
func (pv *PreviewView) previousFile()
func (pv *PreviewView) nextFile()
```

**Multi-File Support**:
- Left card always shows first file (considered "original")
- Right card cycles through files 2, 3, 4, etc.
- "Previous File" / "Next File" navigation
- Progress indicator: "Comparing: 2 of 5"

---

#### Task 3.8: Action Buttons & Decisions
**Duration**: 1 day
**Dependencies**: Task 3.7
**Deliverables**:
- [x] Action button panel (dynamically built)
- [x] Decision tracking (MarkForDeletion)
- [x] "Apply to All" checkbox
- [x] Keep Left / Keep Right / Keep Both buttons
- [x] Results view refresh on decision change

**Implementation Details**:
```go
// internal/ui/preview_view.go
func (pv *PreviewView) keepFile(index int)
func (pv *PreviewView) keepBothFiles(index1, index2 int)
func (pv *PreviewView) refreshResultsView()

// internal/ui/results_view.go
func (rv *ResultsView) createFileList(group) *fyne.Container
func (rv *ResultsView) refreshCards()
```

**Results View Features**:
- Full file path display with ellipsis truncation
- ✓ icon for files to keep (green)
- 🗑 icon for files to delete (red)
- "Keep" / "Delete" status labels
- File size display
- Real-time updates when decisions change
    keepBothBtn     *widget.Button
    autoSelectBtn   *widget.Button
    applyToAllCheck *widget.Check
    saveBtn         *widget.Button
    skipBtn         *widget.Button
    onDecision      func(GroupDecisions)
}

type GroupDecisions struct {
    GroupID       string
    KeepIndices   []int
    DeletePaths   []string
    ApplyToAll    bool
    Rule          *Rule
}

func NewPreviewActions() *PreviewActions
func (p *PreviewActions) SetGroup(group *DuplicateGroup)
func (p *PreviewActions) GetDecisions() GroupDecisions
```

**Action Buttons**:
- **Keep This**: Mark specific file to keep
- **Delete This**: Mark specific file for deletion
- **Keep Both/All**: Don't delete any in this group
- **Auto-Select**: Apply rule (oldest/newest/shortest path)
- **Apply to All**: Checkbox for rule application
- **Save & Continue**: Save decision, go to next
- **Skip**: Skip this group for now

---

#### Task 3.9: Phase 3 Integration & Testing
**Duration**: 0.5 days  
**Dependencies**: Task 3.8  
**Deliverables**:
- [ ] Full preview flow test
- [ ] All file type tests
- [ ] Performance testing (large images/videos)
- [ ] Memory leak checks
- [ ] Bug fixes

**Milestone 3 Complete**: ✅ Full preview system with side-by-side comparison

---

## Phase 4: Rule System (Week 8)

### Week 8: Rules Engine & Application

#### Task 4.1: Rule Definitions
**Duration**: 1 day  
**Dependencies**: Phase 3 complete  
**Deliverables**:
- [ ] `internal/rules/builtin.go` - Built-in rules
- [ ] `internal/rules/rule.go` - Rule interface
- [ ] Rule registration system

**Implementation Details**:
```go
// internal/rules/builtin.go
var BuiltinRules = []Rule{
    {
        ID:          "keep_newest",
        Name:        "Keep Newest",
        Description: "Keep the most recently modified file, delete older copies",
        Criteria:    KeepNewest,
        Icon:        "newest_icon",
    },
    {
        ID:          "keep_oldest",
        Name:        "Keep Oldest",
        Description: "Keep the oldest file, delete newer copies",
        Criteria:    KeepOldest,
        Icon:        "oldest_icon",
    },
    {
        ID:          "keep_shortest_path",
        Name:        "Keep Shortest Path",
        Description: "Keep file with shortest path, delete others",
        Criteria:    KeepShortestPath,
        Icon:        "path_icon",
    },
    {
        ID:          "keep_specific_dir",
        Name:        "Keep from Specific Directory",
        Description: "Keep files from selected directory, delete others",
        Criteria:    KeepSpecificDir,
        Icon:        "folder_icon",
        Config:      SpecificDirConfig{},
    },
    {
        ID:          "manual",
        Name:        "Manual Selection",
        Description: "Manually select which files to keep",
        Criteria:    ManualSelection,
        Icon:        "manual_icon",
    },
}
```

---

#### Task 4.2: Rule Engine
**Duration**: 1.5 days  
**Dependencies**: Task 4.1  
**Deliverables**:
- [ ] `internal/rules/engine.go` - Rule application engine
- [ ] `internal/rules/evaluator.go` - Rule evaluation logic
- [ ] Unit tests for each rule

**Implementation Details**:
```go
// internal/rules/engine.go
type RuleEngine struct {
    rules       map[string]Rule
    evaluators  map[RuleCriteria]RuleEvaluator
}

type RuleEvaluator interface {
    Evaluate(group *DuplicateGroup, rule Rule) (*RuleResult, error)
}

func (e *RuleEngine) ApplyRule(group *DuplicateGroup, ruleID string) (*RuleResult, error)
func (e *RuleEngine) ApplyToAll(groups []DuplicateGroup, ruleID string) ([]RuleResult, error)
func (e *RuleEngine) RegisterEvaluator(criteria RuleCriteria, evaluator RuleEvaluator)
```

**Evaluator Implementations**:
```go
// KeepNewestEvaluator
func (e *KeepNewestEvaluator) Evaluate(group *DuplicateGroup, rule Rule) (*RuleResult, error) {
    // Sort by modification time (newest first)
    // Keep index 0, mark rest for deletion
}

// KeepShortestPathEvaluator
func (e *KeepShortestPathEvaluator) Evaluate(group *DuplicateGroup, rule Rule) (*RuleResult, error) {
    // Sort by path length (shortest first)
    // Keep index 0, mark rest for deletion
}
```

---

#### Task 4.3: Rule Selection Dialog
**Duration**: 1 day  
**Dependencies**: Task 4.2  
**Deliverables**:
- [ ] `internal/ui/rule_dialog.go` - Rule selection UI
- [ ] Rule preview/description
- [ ] Configuration for specific rules
- [ ] Apply to all toggle

**Implementation Details**:
```go
// internal/ui/rule_dialog.go
type RuleDialog struct {
    widget.BaseWidget
    ruleList      *widget.List
    rulePreview   *fyne.Container
    applyToAll    *widget.Check
    configBox     *fyne.Container
    applyButton   *widget.Button
    cancelButton  *widget.Button
    onSelect      func(ruleID string, applyToAll bool)
}

func NewRuleDialog(callback func(string, bool)) *RuleDialog
func (r *RuleDialog) Show()
func (r *RuleDialog) SetSelectedRule(ruleID string)
```

**UI Layout**:
```
┌─────────────────────────────────────────────────────────┐
│  Select Rule                                   [✕]      │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌─────────────────────┬───────────────────────────────┐│
│  │ Available Rules:    │ Preview:                      ││
│  │                     │                               ││
│  │ 🕐 Keep Newest      │ Will keep the most recently   ││
│  │ 📜 Keep Oldest      │ modified file and delete      ││
│  │ 📁 Shortest Path    │ older copies.                 ││
│  │ 📂 Specific Dir     │                               ││
│  │ ✋ Manual           │ Files to keep: 1              ││
│  │                     │ Files to delete: 2            ││
│  │                     │                               ││
│  └─────────────────────┴───────────────────────────────┘│
│                                                         │
│  ☑ Apply this rule to all remaining duplicate groups   │
│                                                         │
│              [Apply Rule]           [Cancel]            │
└─────────────────────────────────────────────────────────┘
```

---

#### Task 4.4: "Apply to All" Functionality
**Duration**: 1 day  
**Dependencies**: Task 4.3  
**Deliverables**:
- [ ] Batch rule application
- [ ] Progress tracking for batch operations
- [ ] Undo batch operation
- [ ] Skip groups where rule doesn't apply

**Implementation Details**:
```go
// internal/rules/batch.go
type BatchRuleApplier struct {
    engine      *RuleEngine
    progress    *Progress
    results     []RuleResult
    skipped     []SkipReason
}

type SkipReason struct {
    GroupID   string
    Reason    string  // "Different file type", "Rule not applicable"
}

func (b *BatchRuleApplier) Apply(ruleID string, groups []DuplicateGroup) error
func (b *BatchRuleApplier) GetResults() []RuleResult
func (b *BatchRuleApplier) GetSkipped() []SkipReason
func (b *BatchRuleApplier) Undo() error
```

---

#### Task 4.5: Decision Review Panel
**Duration**: 0.5 days  
**Dependencies**: Task 4.4  
**Deliverables**:
- [ ] Summary of all decisions
- [ ] Edit individual decisions
- [ ] Filter by decision type
- [ ] Export decision list

**Implementation Details**:
```go
// internal/ui/review_panel.go
type ReviewPanel struct {
    widget.BaseWidget
    list          *widget.List
    decisions     []GroupDecisions
    filterSelect  *widget.Select
    editButton    *widget.Button
    clearButton   *widget.Button
}

func NewReviewPanel() *ReviewPanel
func (r *ReviewPanel) SetDecisions(decisions []GroupDecisions)
func (r *ReviewPanel) Filter(by DecisionType)
func (r *ReviewPanel) EditDecision(index int)
```

---

#### Task 4.6: Phase 4 Integration & Testing
**Duration**: 0.5 days  
**Dependencies**: Task 4.5  
**Deliverables**:
- [ ] Full rule system test
- [ ] All built-in rules tested
- [ ] Batch application tested
- [ ] Edge case handling
- [ ] Bug fixes

**Milestone 4 Complete**: ✅ Complete rule system with batch application

---

## Phase 5: Script Generation (Week 9)

### Week 9: Script Generation & Export

#### Task 5.1: Enhanced Script Generator
**Duration**: 1.5 days  
**Dependencies**: Phase 4 complete  
**Deliverables**:
- [ ] Enhanced script templates
- [ ] Group comments in script
- [ ] Statistics in header
- [ ] Safety checks

**Implementation Details**:
```go
// internal/script/generator.go (enhanced)
type Generator struct {
    outputDir         string
    scriptName        string
    includeHeader     bool
    includeComments   bool
    includeDryRun     bool
    makeExecutable    bool
    createBackup      bool
}

type ScriptContent struct {
    Header          ScriptHeader
    Groups          []ScriptGroup
    Summary         ScriptSummary
    DryRunSupport   bool
}

type ScriptHeader struct {
    GeneratedTime   time.Time
    Version         string
    SourceDirectory string
    TotalFiles      int
    TotalSize       int64
    SizeHuman       string
}

type ScriptGroup struct {
    GroupID       string
    FileName      string
    FileSize      int64
    DeletePaths   []string
    KeepPath      string
}

func (g *Generator) Generate(session *ScanSession) (string, error)
func (g *Generator) GenerateWithDryRun(session *ScanSession) (string, error)
```

---

#### Task 5.2: Script Preview View
**Duration**: 1.5 days  
**Dependencies**: Task 5.1  
**Deliverables**:
- [ ] `internal/ui/script_view.go` - Script preview tab
- [ ] Syntax highlighting for bash
- [ ] Line numbers
- [ ] Copy to clipboard
- [ ] Save to file

**Implementation Details**:
```go
// internal/ui/script_view.go
type ScriptView struct {
    widget.BaseWidget
    scriptText    *widget.RichText
    pathEntry     *widget.Entry
    browseButton  *widget.Button
    executableCheck *widget.Check
    backupCheck   *widget.Check
    dryRunCheck   *widget.Check
    saveButton    *widget.Button
    copyButton    *widget.Button
    generator     *script.Generator
}

func NewScriptView() *ScriptView
func (s *ScriptView) SetSession(session *ScanSession)
func (s *ScriptView) Preview() (string, error)
func (s *ScriptView) Save() error
func (s *ScriptView) CopyToClipboard() error
```

**UI Layout**:
```
┌─────────────────────────────────────────────────────────┐
│  Generated Cleanup Script                               │
├─────────────────────────────────────────────────────────┤
│  Output: [/home/user/cleanup_duplicates.sh___] [Browse] │
├─────────────────────────────────────────────────────────┤
│  ┌───────────────────────────────────────────────────┐  │
│  │ 1  #!/bin/bash                                    │  │
│  │ 2  # ============================================  │  │
│  │ 3  # DupDel Cleanup Script                        │  │
│  │ 4  # Generated: 2024-03-15 14:30:00              │  │
│  │ 5  # ============================================  │  │
│  │ 6                                                 │  │
│  │ 7  # Group 1: document.pdf (2.4 MB)              │  │
│  │ 8  rm -f "/home/user/backup/document.pdf"        │  │
│  │ 9  rm -f "/home/user/old/document.pdf"           │  │
│  │ 10                                                │  │
│  │ 11 # Group 2: IMG_001.jpg (1.2 MB)               │  │
│  │ 12 rm -f "/home/user/backup/IMG_001.jpg"         │  │
│  │ ...                                               │  │
│  └───────────────────────────────────────────────────┘  │
├─────────────────────────────────────────────────────────┤
│  Options:                                               │
│  ☑ Make script executable                               │
│  ☑ Create backup of existing script                     │
│  ☑ Include dry-run support                              │
│                                                         │
│          [Save Script]  [Copy]  [Cancel]                │
└─────────────────────────────────────────────────────────┘
```

---

#### Task 5.3: Export Options
**Duration**: 1 day  
**Dependencies**: Task 5.2  
**Deliverables**:
- [ ] CSV export for report
- [ ] JSON export for programmatic use
- [ ] HTML report (optional)
- [ ] Print dialog

**Implementation Details**:
```go
// internal/export/csv.go
type CSVExporter struct {
    delimiter     rune
    includeHeader bool
}

func (e *CSVExporter) Export(session *ScanSession, path string) error

// internal/export/json.go
type JSONExporter struct {
    prettyPrint   bool
}

func (e *JSONExporter) Export(session *ScanSession, path string) error

// Export formats:
// - CSV: GroupID, FileName, FileSize, KeepPath, DeletePath
// - JSON: Full session data with all metadata
// - HTML: Formatted report with tables and styling
```

---

#### Task 5.4: Script Validation
**Duration**: 0.5 days  
**Dependencies**: Task 5.1  
**Deliverables**:
- [ ] Path existence check
- [ ] Permission check
- [ ] Syntax validation (shellcheck integration optional)
- [ ] Warning for files that no longer exist

**Implementation Details**:
```go
// internal/script/validator.go
type Validator struct {
    warnings []ValidationWarning
    errors   []ValidationError
}

type ValidationWarning struct {
    Path    string
    Message string
}

type ValidationError struct {
    Path    string
    Message string
    Fatal   bool
}

func (v *Validator) Validate(script ScriptContent) error
func (v *Validator) CheckPathsExist(script ScriptContent)
func (v *Validator) CheckPermissions(script ScriptContent)
func (v *Validator) GetWarnings() []ValidationWarning
```

---

#### Task 5.5: Phase 5 Integration & Testing
**Duration**: 0.5 days  
**Dependencies**: Task 5.4  
**Deliverables**:
- [ ] Full script generation flow
- [ ] All export formats tested
- [ ] Script execution test (dry-run)
- [ ] Validation tests
- [ ] Bug fixes

**Milestone 5 Complete**: ✅ Complete script generation and export system

---

## Phase 6: Polish & Release (Weeks 10-11)

### Week 10: Polish & Optimization

#### Task 6.1: Theme Support
**Duration**: 1.5 days  
**Dependencies**: Phase 5 complete  
**Deliverables**:
- [ ] Dark theme
- [ ] Light theme
- [ ] Theme switcher in settings
- [ ] System theme detection (optional)

**Implementation Details**:
```go
// internal/ui/themes.go
type ThemeManager struct {
    currentTheme  fyne.Theme
    themes        map[string]fyne.Theme
}

func NewThemeManager() *ThemeManager
func (t *ThemeManager) LoadTheme(name string) fyne.Theme
func (t *ThemeManager) SetTheme(name string)
func (t *ThemeManager) GetSystemTheme() string
func (t *ThemeManager) AutoDetect() string

// Built-in themes:
// - dark (default)
// - light
// - system (auto-detect)
```

---

#### Task 6.2: Keyboard Shortcuts
**Duration**: 1 day  
**Dependencies**: Task 6.1  
**Deliverables**:
- [ ] Global shortcuts
- [ ] Context-specific shortcuts
- [ ] Shortcut customization (optional)
- [ ] Help dialog showing shortcuts

**Shortcut Map**:
```
Global:
  Ctrl+O      Open directory
  Ctrl+S      Start scan
  Ctrl+E      Stop scan
  Ctrl+G      Generate script
  Ctrl+,      Settings
  F1          Help
  Ctrl+Q      Quit

Results View:
  Ctrl+A      Select all
  Ctrl+D      Deselect all
  Delete      Mark selected for deletion
  Enter       Open preview

Preview View:
  ←/→         Navigate groups
  K           Keep selected
  D           Delete selected
  B           Keep both/all
  A           Toggle apply to all
  S           Save & continue
  Esc         Skip
```

---

#### Task 6.3: Performance Optimization
**Duration**: 1.5 days  
**Dependencies**: Task 6.2  
**Deliverables**:
- [ ] Hash calculation profiling
- [ ] Memory optimization
- [ ] UI responsiveness improvements
- [ ] Lazy loading for results
- [ ] Image caching

**Optimization Targets**:
- Scan speed: >1000 files/second
- Memory: <500MB for 100K files
- UI: 60 FPS during operations
- Hash cache for resume capability

**Profiling Commands**:
```bash
# CPU profiling
go test -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof ./...
go tool pprof mem.prof

# Benchmark
go test -bench=. ./internal/scanner
go test -bench=. ./internal/hasher
go test -bench=. ./internal/detector
```

---

#### Task 6.4: Settings Persistence
**Duration**: 0.5 days  
**Dependencies**: Task 6.3  
**Deliverables**:
- [ ] Configuration file (YAML/JSON)
- [ ] Load settings on startup
- [ ] Save settings on change
- [ ] Reset to defaults option

**Config File Location**:
- Linux: `~/.config/dup-del/config.yaml`
- macOS: `~/Library/Application Support/dup-del/config.yaml`
- Windows: `%APPDATA%\dup-del\config.yaml`

---

#### Task 6.5.1: Test File Generator
**Duration**: 0.5 days
**Dependencies**: None
**Deliverables**:
- [ ] `gen_test_files.sh` - Shell script to generate test data
- [ ] Create 16+ duplicate sets (text, binary, images)
- [ ] Generate PNG, JPEG, GIF, BMP image duplicates
- [ ] Create nested directory structure
- [ ] Include hidden files and special names
- [ ] Document usage in README

**Test Data Structure**:
```
test_data/
├── documents/      # Text duplicates
├── photos/         # Image duplicates (PNG, JPG, GIF, BMP)
├── videos/         # Video duplicates
├── backup/         # Cross-directory duplicates
├── archive/        # Old duplicates
└── work/           # Work-related duplicates
```

---

#### Task 6.6: Help System
**Duration**: 0.5 days  
**Dependencies**: Task 6.5.1  
**Deliverables**:
- [ ] Help dialog
- [ ] Keyboard shortcuts reference
- [ ] FAQ section
- [ ] About dialog with version info

---

### Week 11: Testing & Release

#### Task 6.7: Comprehensive Testing
**Duration**: 2 days
**Dependencies**: Task 6.6
**Deliverables**:
- [ ] Unit test coverage >80%
- [ ] Integration test suite
- [ ] UI test scenarios
- [ ] Cross-platform testing
- [ ] Performance benchmark report

**Test Checklist**:
- [ ] Empty directory scan
- [ ] Single file directory
- [ ] Deep nested directories (10+ levels)
- [ ] Hidden files/directories
- [ ] Symbolic links
- [ ] Permission denied scenarios
- [ ] Large files (>1GB)
- [ ] Many small files (100K+)
- [ ] All supported file types (PNG, JPEG, GIF, BMP, WebP, text)
- [ ] Interrupted scan recovery
- [ ] Script execution (dry-run)
- [ ] Theme switching
- [ ] Settings persistence

---

#### Task 6.8: Documentation
**Duration**: 1 day
**Dependencies**: Task 6.7
**Deliverables**:
- [ ] README.md final version
- [ ] PRODUCT_SPEC.md final version
- [ ] IMPLEMENTATION_PLAN.md (this document)
- [ ] CONTRIBUTING.md
- [ ] CHANGELOG.md
- [ ] User guide (wiki or docs folder)
- [ ] Developer setup guide

---

#### Task 6.9: Release Preparation
**Duration**: 1.5 days
**Dependencies**: Task 6.8
**Deliverables**:
- [ ] Version tagging (git tag v1.0.0)
- [ ] Release notes
- [ ] Binary builds for all platforms
- [ ] Package distribution (optional: .deb, .rpm, .pkg)
- [ ] Website/landing page (optional)
- [ ] Announcement preparation

**Build Commands**:
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o dup-del-linux-amd64 ./cmd/dup-del
GOOS=linux GOARCH=arm64 go build -o dup-del-linux-arm64 ./cmd/dup-del

# macOS
GOOS=darwin GOARCH=amd64 go build -o dup-del-macos-amd64 ./cmd/dup-del
GOOS=darwin GOARCH=arm64 go build -o dup-del-macos-arm64 ./cmd/dup-del

# Windows
GOOS=windows GOARCH=amd64 go build -o dup-del-windows-amd64.exe ./cmd/dup-del
```

---

#### Task 6.10: Beta Testing
**Duration**: Ongoing (parallel with 6.9)
**Dependencies**: Task 6.7
**Deliverables**:
- [ ] Beta tester recruitment
- [ ] Feedback collection system
- [ ] Bug triage process
- [ ] Critical bug fixes

---

#### Task 6.11: Phase 6 & Project Completion
**Duration**: 0.5 days
**Dependencies**: Tasks 6.9, 6.10
**Deliverables**:
- [ ] Final QA pass
- [ ] Release candidate build
- [ ] v1.0.0 release
- [ ] Post-release monitoring plan

**Milestone 6 Complete**: ✅ Production release ready

---

## Summary Timeline

| Phase | Weeks | Key Deliverables |
|-------|-------|------------------|
| 1. Core Functionality | 1-3 | CLI scanner, hasher, detector, script generator |
| 2. GUI Foundation | 4-5 | Main window, scan view, results view |
| 3. Preview System | 6-7 | Text/image/video preview, side-by-side comparison |
| 4. Rule System | 8 | Rule engine, batch application, decision review |
| 5. Script Generation | 9 | Enhanced generator, script preview, export options |
| 6. Polish & Release | 10-11 | Themes, shortcuts, optimization, testing, release |

---

## Critical Path

```
Week 1-2: Data Models → Scanner → Hasher → Detector
                    ↓
Week 3: CLI + Script Generator (Milestone 1)
                    ↓
Week 4-5: Fyne Setup → Main Window → Scan View → Results View (Milestone 2)
                    ↓
Week 6-7: Preview System → Text/Image/Video → Side-by-Side (Milestone 3)
                    ↓
Week 8: Rule Engine → Batch Application (Milestone 4)
                    ↓
Week 9: Script Generator → Script Preview → Export (Milestone 5)
                    ↓
Week 10-11: Polish → Testing → Release (Milestone 6)
```

---

## Risk Mitigation

| Risk | Mitigation Strategy |
|------|---------------------|
| Fyne learning curve | Start Fyne integration early (Week 4), build prototypes |
| Video thumbnail complexity | Implement fallback chain, make FFmpeg optional |
| Performance issues | Profile early, implement worker pools, lazy loading |
| Cross-platform bugs | CI/CD with all platforms, early testing |
| Scope creep | Strict adherence to v1.0 requirements, defer nice-to-haves |

---

## Success Criteria

- [ ] All Milestones 1-6 achieved
- [ ] Unit test coverage >80%
- [ ] Performance targets met
- [ ] Zero critical bugs
- [ ] Cross-platform builds successful
- [ ] Documentation complete
- [ ] v1.0.0 released

---

**Document Version**: 1.0.0  
**Last Updated**: 2024-03-15  
**Total Estimated Effort**: 11 weeks (55 person-days)
