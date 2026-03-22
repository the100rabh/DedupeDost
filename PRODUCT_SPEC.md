# DupDel - Product Specification

## 1. Executive Summary

**Product Name**: DupDel (Duplicate File Deleter)

**Version**: 1.0.0

**Type**: Desktop Application with GUI

**Platform**: Cross-platform (Windows, macOS, Linux)

**Technology Stack**: Go (Golang) + Fyne GUI Framework

**Purpose**: A user-friendly application that scans directories for duplicate files, provides visual comparison tools, and generates safe cleanup scripts without directly modifying user data.

---

## 2. Problem Statement

Users accumulate duplicate files over time through:
- Multiple downloads of the same content
- Backup copies forgotten in different locations
- Sync conflicts creating multiple versions
- Photo/video duplicates from device transfers
- Copy-paste operations without cleanup

These duplicates waste storage space, create confusion, and make file management difficult. Existing solutions either:
- Delete files directly (risky)
- Lack visual preview capabilities
- Have poor user interfaces
- Don't support batch rule application

---

## 3. Solution Overview

DupDel provides a safe, visual, and user-controlled approach to duplicate file management:

1. **Scan**: Recursively analyze directory contents
2. **Detect**: Identify duplicates using hash-based comparison
3. **Preview**: Show files side-by-side for informed decisions
4. **Decide**: User controls which files to mark for deletion
5. **Generate**: Create a reviewable bash script for cleanup

---

## 4. Target Users

| User Type | Use Case |
|-----------|----------|
| Home Users | Clean up personal photo collections, documents |
| Developers | Remove duplicate code files, dependencies |
| Designers | Manage duplicate assets, images, videos |
| IT Administrators | Clean shared drives, user directories |
| Power Users | Batch process large file collections |

---

## 5. Functional Requirements

### 5.1 File Scanning

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-001 | Scan all files in selected directory | Must Have |
| FR-002 | Recursively scan subdirectories | Must Have |
| FR-003 | Filter by file extensions | Should Have |
| FR-004 | Filter by file size range | Should Have |
| FR-005 | Option to include/exclude hidden files | Should Have |
| FR-006 | Display real-time scan progress | Must Have |
| FR-007 | Show scan statistics (files scanned, size) | Must Have |
| FR-008 | Cancel ongoing scan | Should Have |
| FR-009 | Resume interrupted scan | Could Have |

### 5.2 Duplicate Detection

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-010 | Group files by size (initial filter) | Must Have |
| FR-011 | Calculate SHA-256 hash for comparison | Must Have |
| FR-012 | Handle large files efficiently (chunked hashing) | Must Have |
| FR-013 | Detect duplicates across different directories | Must Have |
| FR-014 | Handle symbolic links appropriately | Should Have |
| FR-015 | Detect near-duplicates (future enhancement) | Won't Have (v1) |

### 5.3 File Preview

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-020 | Side-by-side file comparison view | Must Have |
| FR-021 | Text file preview with full content (no truncation) | Must Have |
| FR-022 | Image preview (PNG, JPEG, GIF, BMP, WebP) with zoom | Must Have |
| FR-023 | Video thumbnail and metadata display | Must Have |
| FR-024 | Show file metadata (size, dates, path) | Must Have |
| FR-025 | Navigate between duplicate groups | Must Have |
| FR-026 | Preview unsupported file types (hash only) | Should Have |
| FR-027 | Text word wrap and monospace font | Must Have |
| FR-028 | Image fit-to-screen and zoom controls | Must Have |
| FR-029 | Navigate through all files in a group (>2 files) | Must Have |
| FR-030 | Show file path in results view | Must Have |
| FR-031 | Visual keep/delete indicators | Must Have |
| FR-032 | Video open in default player | Must Have |
| FR-033 | Video quality indicator (4K/1080p/720p) | Should Have |

### 5.4 User Actions

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-035 | Select files to mark for deletion | Must Have |
| FR-036 | Keep all files in a group option | Must Have |
| FR-037 | Auto-select oldest file for deletion | Should Have |
| FR-038 | Auto-select deepest path for deletion | Should Have |
| FR-039 | "Apply to All" rule for similar files | Must Have |
| FR-040 | Undo last action | Should Have |
| FR-041 | Review all marked files before script generation | Must Have |
| FR-042 | Keep Left / Keep Right / Keep Both actions | Must Have |
| FR-043 | Batch Select All / Deselect All | Should Have |

### 5.5 Rule System

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-040 | Apply rule to all remaining groups | Must Have |
| FR-041 | Rule: Keep newest, delete older | Must Have |
| FR-042 | Rule: Keep oldest, delete newer | Should Have |
| FR-043 | Rule: Keep shortest path, delete others | Should Have |
| FR-044 | Rule: Keep specific directory, delete others | Could Have |
| FR-045 | Rule: Always ask (no auto) | Must Have |
| FR-046 | Save custom rules for future sessions | Could Have |

### 5.6 Script Generation

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-050 | Generate bash script for deletions | Must Have |
| FR-051 | Script uses absolute paths | Must Have |
| FR-052 | Script includes comments for each group | Must Have |
| FR-053 | Script includes header with metadata | Must Have |
| FR-054 | Script is executable (chmod +x ready) | Must Have |
| FR-055 | Option to generate dry-run script | Should Have |
| FR-056 | Export deletion report (CSV/JSON) | Should Have |
| FR-057 | Create backup of script before generation | Should Have |

### 5.7 User Interface

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-060 | Modern, clean GUI design | Must Have |
| FR-061 | Dark/Light theme support | Should Have |
| FR-062 | Responsive layout | Must Have |
| FR-063 | Keyboard shortcuts | Should Have |
| FR-064 | Multi-language support (i18n) | Could Have |
| FR-065 | Accessibility features | Could Have |
| FR-066 | System tray integration | Won't Have (v1) |

---

## 6. Non-Functional Requirements

### 6.1 Performance

| ID | Requirement | Target |
|----|-------------|--------|
| NFR-001 | Scan speed | >1000 files/second |
| NFR-002 | Memory usage | <500MB for 100K files |
| NFR-003 | Hash calculation | Async, non-blocking UI |
| NFR-004 | Large file handling | Stream hashing for files >100MB |
| NFR-005 | Startup time | <3 seconds |

### 6.2 Reliability

| ID | Requirement | Target |
|----|-------------|--------|
| NFR-010 | Crash-free sessions | >99.9% |
| NFR-011 | Handle permission errors gracefully | 100% |
| NFR-012 | Handle corrupted files | 100% |
| NFR-013 | Recover from interrupted scans | Yes |

### 6.3 Security

| ID | Requirement | Description |
|----|-------------|-------------|
| NFR-020 | No direct file deletion | App never deletes files directly |
| NFR-021 | User review required | Script must be reviewed before execution |
| NFR-022 | No network calls | Application works offline |
| NFR-023 | No data collection | Privacy-first design |

### 6.4 Compatibility

| Platform | Minimum Version |
|----------|-----------------|
| Windows | Windows 10 |
| macOS | macOS 10.15 (Catalina) |
| Linux | Kernel 4.0+, GTK3 |

---

## 7. User Interface Specification

### 7.1 Main Window Layout

```
┌─────────────────────────────────────────────────────────────────┐
│  [App Icon] DupDel                              ─ □ ×          │
├─────────────────────────────────────────────────────────────────┤
│  MENU BAR                                                       │
│  File    Edit    View    Tools    Help                          │
├─────────────────────────────────────────────────────────────────┤
│  TOOLBAR                                                        │
│  [📁 Open] [▶️ Scan] [⏹ Stop] [💾 Save] [⚙ Settings] [❓ Help] │
├─────────────────────────────────────────────────────────────────┤
│  DIRECTORY SELECTION                                            │
│  📁 Path: [/home/user/documents____________] [Browse...]        │
├─────────────────────────────────────────────────────────────────┤
│  FILTER OPTIONS                                                 │
│  ☑ Recursive  ☐ Hidden  Min: [1KB▼]  Max: [None▼]  Types: [...]│
├─────────────────────────────────────────────────────────────────┤
│  MAIN CONTENT AREA (Tabbed)                                     │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │ [Scan] [Results] [Preview] [Script Preview]                │ │
│  │                                                           │ │
│  │ Content changes based on selected tab                     │ │
│  │                                                           │ │
│  └───────────────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────────┤
│  STATUS BAR                                                     │
│  Ready | Files: 0 | Duplicates: 0 | Size: 0 B | [Progress Bar] │
└─────────────────────────────────────────────────────────────────┘
```

### 7.2 Scan View

```
┌─────────────────────────────────────────────────────────────────┐
│  SCANNING DIRECTORY                                             │
│                                                                 │
│  📁 /home/user/documents                                        │
│                                                                 │
│  Progress: ████████████████░░░░░░░░ 65%                         │
│                                                                 │
│  Files scanned: 12,453 / 19,158                                 │
│  Current: /home/user/documents/photos/2024/IMG_1234.jpg        │
│                                                                 │
│  Statistics:                                                    │
│  ├── Total size: 4.5 GB                                         │
│  ├── Files processed: 12,453                                    │
│  ├── Potential duplicates: 234                                  │
│  └── Elapsed time: 00:02:34                                     │
│                                                                 │
│              [⏹ Stop Scan]                                      │
└─────────────────────────────────────────────────────────────────┘
```

### 7.3 Results View

```
┌─────────────────────────────────────────────────────────────────┐
│  DUPLICATE GROUPS (45 groups, 127 files, 523.4 MB)              │
├─────────────────────────────────────────────────────────────────┤
│  Search: [________________]  Sort: [Size ▼]  View: [List ◡]    │
├─────────────────────────────────────────────────────────────────┤
│  ┌───────────────────────────────────────────────────────────┐ │
│  │ 📄 document.pdf                                            │ │
│  │    3 copies | 2.4 MB each | Total: 4.8 MB recoverable    │ │
│  │    /home/user/docs/document.pdf                          │ │
│  │    /home/user/backup/document.pdf                        │ │
│  │    /home/user/old/document.pdf                           │ │
│  └───────────────────────────────────────────────────────────┘ │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │ 🖼️ IMG_001.jpg                                             │ │
│  │    5 copies | 1.2 MB each | Total: 4.8 MB recoverable    │ │
│  └───────────────────────────────────────────────────────────┘ │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │ 🎥 vacation.mp4                                            │ │
│  │    2 copies | 156 MB each | Total: 156 MB recoverable    │ │
│  └───────────────────────────────────────────────────────────┘ │
│                                                                 │
│  [Select All] [Deselect All] [Invert Selection]                │
│                                                                 │
│  Total recoverable: 523.4 MB across 127 files                  │
│                    [Review & Generate Script]                   │
└─────────────────────────────────────────────────────────────────┘
```

### 7.4 Preview/Comparison View

```
┌─────────────────────────────────────────────────────────────────┐
│  Group 3 of 45  |  🖼️ IMG_001.jpg  |  ◀ Previous  Next ▶       │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌────────────────────────────┬────────────────────────────────┐│
│  │ FILE 1                     │ FILE 2                         ││
│  │ /home/user/photos/IMG_001.jpg                               ││
│  │                            │ /home/user/backup/IMG_001.jpg  ││
│  │ ┌──────────────────────┐  │ ┌──────────────────────────────┐││
│  │ │                      │  │ │                              │││
│  │ │   [IMAGE PREVIEW]    │  │ │     [IMAGE PREVIEW]          │││
│  │ │                      │  │ │                              │││
│  │ │   Zoom: 100%         │  │ │     Zoom: 100%               │││
│  │ └──────────────────────┘  │ └──────────────────────────────┘││
│  │                            │                                ││
│  │ Size: 1.2 MB               │ Size: 1.2 MB                   ││
│  │ Created: 2024-01-15 10:30 │ Created: 2024-01-15 10:30     ││
│  │ Modified: 2024-01-15 10:30│ Modified: 2024-02-01 14:20    ││
│  │ Hash: abc123...           │ Hash: abc123...                ││
│  │                            │                                ││
│  │ ○ Keep  ● Delete          │ ● Keep  ○ Delete               ││
│  └────────────────────────────┴────────────────────────────────┘│
│                                                                 │
│  Actions:                                                       │
│  [Keep This] [Delete This] [Keep Both] [Auto-Select Oldest]    │
│                                                                 │
│  Apply rule to remaining: ☑ Yes, apply to all similar files    │
│                                                                 │
│              [Save & Continue]     [Skip This Group]            │
└─────────────────────────────────────────────────────────────────┘
```

### 7.5 Script Preview View

```
┌─────────────────────────────────────────────────────────────────┐
│  GENERATED CLEANUP SCRIPT                                       │
├─────────────────────────────────────────────────────────────────┤
│  Output: [./cleanup_duplicates.sh____________] [Change...]      │
├─────────────────────────────────────────────────────────────────┤
│  ┌───────────────────────────────────────────────────────────┐ │
│  │ #!/bin/bash                                                │ │
│  │ # DupDel Cleanup Script                                    │ │
│  │ # Generated: 2024-03-15 14:30:00                          │ │
│  │ # Source Directory: /home/user/documents                   │ │
│  │ # Total files to delete: 127                              │ │
│  │ # Total space to recover: 523.4 MB                        │ │
│  │                                                            │ │
│  │ # Group 1: document.pdf (2.4 MB)                          │ │
│  │ rm -f "/home/user/backup/document.pdf"                    │ │
│  │ rm -f "/home/user/old/document.pdf"                       │ │
│  │                                                            │ │
│  │ # Group 2: IMG_001.jpg (1.2 MB)                           │ │
│  │ rm -f "/home/user/backup/IMG_001.jpg"                     │ │
│  │ rm -f "/home/user/downloads/IMG_001.jpg"                  │ │
│  │ ...                                                        │ │
│  └───────────────────────────────────────────────────────────┘ │
│                                                                 │
│  Options:                                                       │
│  ☑ Make script executable                                       │
│  ☑ Create backup of script                                      │
│  ☐ Include dry-run option                                       │
│                                                                 │
│              [Save Script]  [Copy to Clipboard]  [Cancel]       │
└─────────────────────────────────────────────────────────────────┘
```

### 7.6 Settings Dialog

```
┌─────────────────────────────────────────────────────────────────┐
│  Settings                                              [✕]      │
├─────────────────────────────────────────────────────────────────┤
│  [General] [Scan] [Preview] [Output] [Advanced]                 │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  SCAN SETTINGS                                                  │
│                                                                 │
│  Default behavior:                                              │
│  ☑ Scan subdirectories recursively                              │
│  ☐ Include hidden files and folders                             │
│  ☐ Follow symbolic links                                        │
│                                                                 │
│  Size filters:                                                  │
│  Minimum file size: [1      KB ▼]                               │
│  Maximum file size: [Unlimited ▼]                               │
│                                                                 │
│  File types to scan:                                            │
│  ☑ Images (jpg, png, gif, bmp, webp)                           │
│  ☑ Videos (mp4, avi, mkv, mov)                                 │
│  ☑ Documents (pdf, doc, docx)                                  │
│  ☑ Text files (txt, md, json, xml)                             │
│  ☑ Archives (zip, tar, gz, rar)                                │
│  ☐ All other files                                             │
│                                                                 │
│  HASH SETTINGS                                                  │
│                                                                 │
│  Hash algorithm: [SHA-256 ▼]                                    │
│  ☑ Use size-based prefiltering                                  │
│  Chunk size for large files: [4 MB ▼]                          │
│                                                                 │
│              [Save]  [Reset to Defaults]  [Cancel]              │
└─────────────────────────────────────────────────────────────────┘
```

---

## 8. Data Models

### 8.1 File Entry

```go
type FileEntry struct {
    Path        string    // Absolute file path
    Name        string    // File name with extension
    Size        int64     // File size in bytes
    Extension   string    // File extension (lowercase)
    ModTime     time.Time // Last modified time
    CreatedTime time.Time // Creation time (if available)
    Hash        string    // SHA-256 hash
    FileType    FileType  // Enum: Text, Image, Video, Other
    IsHidden    bool      // Whether file is hidden
    IsSymlink   bool      // Whether file is a symbolic link
}
```

### 8.2 Duplicate Group

```go
type DuplicateGroup struct {
    ID          string      // Unique group identifier
    Hash        string      // Common hash value
    Size        int64       // File size (all files have same size)
    Extension   string      // File extension
    FileType    FileType    // Type of files in group
    Files       []FileEntry // All duplicate files
    KeepIndices []int       // Indices of files to keep
    DeletePaths []string    // Paths marked for deletion
    RuleApplied *Rule       // Rule applied to this group (if any)
}
```

### 8.3 Rule

```go
type Rule struct {
    ID          string      // Rule identifier
    Name        string      // Human-readable name
    Description string      // Rule description
    Criteria    RuleCriteria // How to select files to keep
    ApplyToAll  bool        // Whether to apply to remaining groups
}

type RuleCriteria string
const (
    KeepNewest     RuleCriteria = "newest"
    KeepOldest     RuleCriteria = "oldest"
    KeepShortestPath RuleCriteria = "shortest_path"
    KeepLongestPath  RuleCriteria = "longest_path"
    KeepSpecificDir  RuleCriteria = "specific_dir"
    ManualSelection  RuleCriteria = "manual"
)
```

### 8.4 Scan Session

```go
type ScanSession struct {
    ID            string        // Session identifier
    SourceDir     string        // Scanned directory
    StartTime     time.Time     // When scan started
    EndTime       time.Time     // When scan ended
    TotalFiles    int           // Total files scanned
    TotalSize     int64         // Total size in bytes
    DuplicateGroups []DuplicateGroup // Found duplicate groups
    FilesToDelete   int           // Count of files marked for deletion
    SpaceRecoverable int64       // Recoverable space in bytes
    GeneratedScript string        // Path to generated script
}
```

---

## 9. Technical Architecture

### 9.1 Component Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                         PRESENTATION LAYER                       │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │   Main      │  │   Scan      │  │   Preview   │             │
│  │   Window    │  │   View      │  │   View      │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │   Results   │  │   Script    │  │   Settings  │             │
│  │   View      │  │   Preview   │  │   Dialog    │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                        APPLICATION LAYER                         │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │   Scan      │  │   Rule      │  │   Script    │             │
│  │   Manager   │  │   Engine    │  │   Generator │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
│  ┌─────────────┐  ┌─────────────┐                              │
│  │   Session   │  │   Event     │                              │
│  │   Manager   │  │   Bus       │                              │
│  └─────────────┘  └─────────────┘                              │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                         DOMAIN LAYER                             │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │   File      │  │   Duplicate │  │   Rule      │             │
│  │   Entry     │  │   Group     │  │   Criteria  │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
│  ┌─────────────┐  ┌─────────────┐                              │
│  │   Scan      │  │   Script    │                              │
│  │   Session   │  │   Template  │                              │
│  └─────────────┘  └─────────────┘                              │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                        INFRASTRUCTURE LAYER                      │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │   File      │  │   Hash      │  │   Image     │             │
│  │   System    │  │   Calculator│  │   Processor │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │   Video     │  │   Text      │  │   Config    │             │
│  │   Processor │  │   Renderer  │  │   Store     │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
└─────────────────────────────────────────────────────────────────┘
```

### 9.2 Package Structure

```
dup-del/
├── cmd/
│   └── dup-del/
│       └── main.go              # Application entry point
├── internal/
│   ├── app/
│   │   ├── app.go               # Application lifecycle
│   │   └── config.go            # Configuration management
│   ├── scanner/
│   │   ├── scanner.go           # Directory traversal
│   │   ├── filter.go            # File filtering
│   │   └── progress.go          # Progress tracking
│   ├── hasher/
│   │   ├── hasher.go            # Hash calculation
│   │   └── pool.go              # Worker pool for hashing
│   ├── detector/
│   │   ├── detector.go          # Duplicate detection logic
│   │   └── grouper.go           # Group duplicates
│   ├── preview/
│   │   ├── preview.go           # Preview interface
│   │   ├── text.go              # Text file preview
│   │   ├── image.go             # Image preview
│   │   └── video.go             # Video thumbnail
│   ├── rules/
│   │   ├── rule.go              # Rule definitions
│   │   ├── engine.go            # Rule application engine
│   │   └── builtin.go           # Built-in rules
│   ├── script/
│   │   ├── generator.go         # Script generation
│   │   └── template.go          # Script templates
│   └── ui/
│       ├── ui.go                # UI initialization
│       ├── main_window.go       # Main window
│       ├── scan_view.go         # Scan view
│       ├── results_view.go      # Results list
│       ├── preview_view.go      # Comparison view
│       ├── script_view.go       # Script preview
│       ├── settings_view.go     # Settings dialog
│       └── widgets/             # Custom widgets
│           ├── file_card.go
│           ├── progress_bar.go
│           └── theme.go
├── pkg/
│   ├── models/
│   │   ├── file.go              # FileEntry model
│   │   ├── group.go             # DuplicateGroup model
│   │   ├── rule.go              # Rule model
│   │   └── session.go           # ScanSession model
│   └── types/
│       ├── types.go             # Type definitions
│       └── constants.go         # Application constants
├── assets/
│   ├── icons/                   # Application icons
│   ├── themes/                  # Theme definitions
│   └── templates/               # Script templates
├── tests/
│   ├── scanner/
│   ├── hasher/
│   ├── detector/
│   └── script/
├── go.mod
├── go.sum
├── README.md
└── PRODUCT_SPEC.md
```

---

## 10. Implementation Phases

### Phase 1: Core Functionality (Weeks 1-3)

| Task | Description | Priority |
|------|-------------|----------|
| P1-T1 | Project setup and dependencies | Critical |
| P1-T2 | File system scanner | Critical |
| P1-T3 | Hash calculation module | Critical |
| P1-T4 | Duplicate detection algorithm | Critical |
| P1-T5 | Basic CLI interface | High |
| P1-T6 | Unit tests for core modules | High |

### Phase 2: GUI Foundation (Weeks 4-5)

| Task | Description | Priority |
|------|-------------|----------|
| P2-T1 | Main window layout | Critical |
| P2-T2 | Directory selection dialog | Critical |
| P2-T3 | Scan progress view | Critical |
| P2-T4 | Results list view | High |
| P2-T5 | Basic settings dialog | Medium |

### Phase 3: Preview System (Weeks 6-7)

| Task | Description | Priority |
|------|-------------|----------|
| P3-T1 | Text file preview | Critical |
| P3-T2 | Image preview with zoom | Critical |
| P3-T3 | Video thumbnail generation | High |
| P3-T4 | Side-by-side comparison view | Critical |
| P3-T5 | File metadata display | High |

### Phase 4: Rule System (Week 8)

| Task | Description | Priority |
|------|-------------|----------|
| P4-T1 | Rule definition interface | Critical |
| P4-T2 | Built-in rules implementation | Critical |
| P4-T3 | "Apply to All" functionality | Critical |
| P4-T4 | Manual selection interface | High |

### Phase 5: Script Generation (Week 9)

| Task | Description | Priority |
|------|-------------|----------|
| P5-T1 | Bash script generator | Critical |
| P5-T2 | Script preview view | Critical |
| P5-T3 | Export options (CSV, JSON) | Medium |
| P5-T4 | Script validation | High |

### Phase 6: Polish & Release (Weeks 10-11)

| Task | Description | Priority |
|------|-------------|----------|
| P6-T1 | Theme support (dark/light) | Medium |
| P6-T2 | Keyboard shortcuts | Medium |
| P6-T3 | Performance optimization | High |
| P6-T4 | Documentation | High |
| P6-T5 | Beta testing | High |
| P6-T6 | Release preparation | Critical |

---

## 11. Testing Strategy

### 11.1 Unit Tests

- File scanning with various directory structures
- Hash calculation accuracy
- Duplicate detection correctness
- Rule application logic
- Script generation format

### 11.2 Integration Tests

- End-to-end scan workflow
- UI component interactions
- File preview rendering
- Script execution (dry-run)

### 11.3 Performance Tests

- Large directory scans (100K+ files)
- Memory usage under load
- Hash calculation throughput
- UI responsiveness during scan

### 11.4 User Acceptance Tests

- Common use case scenarios
- Edge cases (empty dirs, permission errors)
- Cross-platform compatibility
- Accessibility compliance

---

## 12. Success Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| Scan accuracy | 100% | Hash collision rate |
| False positive rate | 0% | User reports |
| User satisfaction | >4.5/5 | User surveys |
| Crash rate | <0.1% | Error reports |
| Average scan speed | >1000 files/sec | Performance tests |
| Memory efficiency | <500MB/100K files | Profiling |

---

## 13. Risks and Mitigations

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Hash collisions | High | Very Low | Use SHA-256, add size verification |
| Large file handling | Medium | Medium | Chunked hashing, progress updates |
| UI freezing during scan | Medium | Medium | Async operations, worker pools |
| Permission errors | Low | High | Graceful handling, user notifications |
| Cross-platform issues | Medium | Medium | CI/CD with all target platforms |
| Video thumbnail generation | Low | Medium | Use external tools, fallback to icons |

---

## 14. Future Enhancements

| Feature | Description | Priority | Timeline |
|---------|-------------|----------|----------|
| Near-duplicate detection | Perceptual hashing for similar images | Low | v1.2 |
| Cloud storage support | Scan Google Drive, Dropbox, OneDrive | Low | v2.0 |
| Plugin system | Custom file type handlers | Low | v2.0 |
| Content-aware preview | PDF, Office documents | Medium | v1.1 |
| Scheduled scans | Automatic duplicate monitoring | Low | v1.2 |
| Network drive support | Scan SMB, NFS shares | Low | v1.1 |
| Enhanced image formats | RAW, HEIC, HEIF full support | Medium | v1.1 |
| FFmpeg bundled | Include FFmpeg for better video support | Medium | v1.1 |

---

## 15. Appendix

### 15.1 Glossary

| Term | Definition |
|------|------------|
| Duplicate Group | A set of files with identical content (hash) |
| Hash | SHA-256 cryptographic hash of file contents |
| Rule | A predefined criteria for selecting files to keep/delete |
| Session | A complete scan operation from start to script generation |

### 15.2 References

- [Fyne Documentation](https://fyne.io/docs/)
- [Go Cryptographic Hashing](https://pkg.go.dev/crypto/sha256)
- [Bash Best Practices](https://google.github.io/styleguide/shellguide.html)

---

**Document Version**: 1.0.0  
**Last Updated**: 2024-03-15  
**Author**: DupDel Team
