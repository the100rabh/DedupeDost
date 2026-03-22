# DedupeDost - Duplicate File Scanner & Cleaner

A powerful Go application that scans directories recursively to find duplicate files, previews them side-by-side, and generates safe cleanup scripts.

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-blue.svg)
![Status](https://img.shields.io/badge/Status-Active%20Development-yellow)

## ✨ Features

### Core Features
- 🔍 **Deep Scanning**: Recursively scans all files in a directory and its subdirectories
- 📸 **Smart Preview**: Side-by-side comparison for text files, images (PNG, JPEG, GIF, BMP, WebP), and videos
- 🎯 **Rule-Based Actions**: Apply decisions to all similar files automatically
- 📝 **Safe Cleanup**: Generates a bash script instead of directly deleting files
- ⚡ **Fast Hashing**: Uses SHA-256 for accurate duplicate detection
- 📊 **Progress Tracking**: Real-time scan progress and statistics

### Preview Capabilities
- **Text Files**: Full content preview with syntax highlighting for 15+ languages (Go, Python, JS, Java, C++, etc.), word wrap, no truncation
- **Images**: Full preview for PNG, JPEG, GIF, BMP, WebP formats with zoom (25%-400%), fit-to-screen, and metadata display
- **Videos**: Thumbnail preview with metadata (duration, resolution, codec, bitrate, FPS), open in default player
- **Other Files**: Hash comparison with file metadata

### Results View
- **File Paths**: Shows full absolute path for each duplicate file
- **Keep/Delete Indicators**: Visual ✓/🗑 icons showing which files will be kept or deleted
- **Status Labels**: Clear "Keep" (green) or "Delete" (red) labels for each file
- **File Sizes**: Display size for each file in the duplicate group
- **Batch Actions**: "Select All" and "Deselect All" buttons for quick decisions

### Preview View
- **Multi-File Comparison**: Navigate through all files in a duplicate group (not just 2)
- **File Navigation**: Previous/Next buttons to cycle through duplicate copies
- **Keep/Delete Actions**: "Keep Left", "Keep Right", "Keep Both" buttons
- **Apply to All**: Apply current decision to all remaining duplicate groups
- **Progress Tracking**: Shows current position (e.g., "Comparing: 2 of 5")

### Smart Rules
- Keep Newest/Oldest (by modification time)
- Keep Shortest/Longest Path
- Keep from Specific Directory
- Keep Largest/Smallest
- Manual Selection
- Smart auto-selection based on folder patterns

## 🧪 Testing

### Quick Start

**Linux/macOS:**
```bash
./run_tests.sh              # Run all tests with coverage
./run_tests.sh -v           # Run with verbose output
./run_tests.sh -race        # Run with race detector
./run_tests.sh -html        # Generate HTML coverage report
./run_tests.sh -help        # Show all options
```

**Windows:**
```cmd
run_tests.bat              REM Run all tests with coverage
run_tests.bat -v           REM Run with verbose output
run_tests.bat -race        REM Run with race detector
run_tests.bat -html        REM Generate HTML coverage report
```

### Manual Testing

**Run all tests:**
```bash
go test ./...
```

**Run with coverage:**
```bash
go test ./... -cover
```

**Run with coverage profile:**
```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out          # View coverage by function
go tool cover -html=coverage.out          # Generate HTML report
```

**Run specific packages:**
```bash
go test ./pkg/...                         # Test pkg packages
go test ./internal/...                    # Test internal packages
go test ./cmd/...                         # Test cmd packages
```

**Run with race detector:**
```bash
go test ./... -race -cover
```

### GUI Tests

To run GUI tests (requires X11/OpenGL on Linux):

```bash
# Install dependencies (Ubuntu/Debian)
sudo apt-get install libgl1-mesa-dev xorg-dev

# Run all tests including GUI
./run_tests.sh -gui
```

### Coverage Goals

| Package | Target | Status |
|---------|--------|--------|
| pkg/types | 100% | ✅ 100.0% |
| pkg/models | 95% | ✅ 96.8% |
| internal/script | 90% | ✅ 94.7% |
| internal/detector | 90% | ✅ 92.6% |
| internal/hasher | 80% | ✅ 84.8% |
| internal/rules | 70% | ✅ 71.1% |
| internal/scanner | 70% | ✅ 71.1% |
| cmd/dedupedost | 70% | ✅ 71.7% |

Note: GUI packages (`internal/ui`, `cmd/dedupedost-gui`) require X11/OpenGL and are tested separately.

## 📦 Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/dedupedost.git
cd dedupedost

# Install dependencies
go mod download

# Build the CLI application
go build -o dedupedost ./cmd/dedupedost

# Run the application
./dedupedost --help
```

### Build GUI (Requires System Dependencies)

```bash
# Ubuntu/Debian - Install dependencies
sudo apt-get install libgl1-mesa-dev xorg-dev

# macOS - Install Xcode command line tools
xcode-select --install

# Build GUI application
go build -o dedupedost-gui ./cmd/dedupedost-gui
./dedupedost-gui
```

## 🚀 Usage

### Command Line Mode

```bash
# Basic scan of current directory
./dedupedost

# Scan a specific directory
./dedupedost --dir /path/to/scan

# With size filters
./dedupedost --dir /path/to/scan --min-size 1024 --max-size 104857600

# Scan specific file types only
./dedupedost --dir /path/to/scan --types "jpg,png,txt"

# Custom output script location
./dedupedost --dir /path/to/scan --output /tmp/cleanup.sh

# Skip hidden files
./dedupedost --dir /path/to/scan --ignore-hidden
```

### Command Line Options

| Flag | Description | Default |
|------|-------------|---------|
| `--dir`, `-d` | Directory to scan | Current directory |
| `--min-size` | Minimum file size in bytes | 0 |
| `--max-size` | Maximum file size in bytes | Unlimited |
| `--types` | Comma-separated extensions | All files |
| `--output`, `-o` | Output script path | `cleanup_duplicates.sh` |
| `--recursive`, `-r` | Scan subdirectories | `true` |
| `--ignore-hidden` | Skip hidden files | `false` |
| `--verbose`, `-v` | Verbose output | `false` |
| `--help`, `-h` | Show help | - |
| `--version` | Show version | - |

### Example Output

```
$ ./dedupedost -dir ./test_data

DedupeDost v1.0.0 - Duplicate File Scanner
======================================

Scanning: /home/user/test_data
Options: Recursive=true, Hidden=true, MinSize=0 B

Scanning files...

======================================
Scan Complete!
======================================
Total files scanned: 150
Total size: 45.2 MB
Duplicate groups: 12
Duplicate files: 28
Recoverable space: 15.3 MB
Scan duration: 234ms

Top duplicates by size:
  1. vacation.mp4 (3 copies, 5.2 MB each) - Recoverable: 10.4 MB
  2. document.pdf (2 copies, 1.5 MB each) - Recoverable: 1.5 MB
  3. IMG_001.jpg (4 copies, 256 KB each) - Recoverable: 768 KB

Generating cleanup script: cleanup_duplicates.sh

======================================
Done!
======================================
Cleanup script: cleanup_duplicates.sh

To preview deletions:
  DRY_RUN=true cleanup_duplicates.sh

To execute deletions:
  ./cleanup_duplicates.sh
```

## 🛡️ Safe by Design

DedupeDost **never deletes files directly**. Instead, it:

1. Scans and identifies duplicates using SHA-256 hashing
2. Shows results with full file paths and keep/delete indicators
3. Lets you review each duplicate group with side-by-side previews
4. Navigate through all copies in a group (not just 2)
5. Apply decisions to all remaining groups with one click
6. Generates a bash script with all deletion commands
7. You review the script before executing

### Results View Features

- **Full file paths** - See exactly where each duplicate is located
- **Visual indicators** - ✓ for files to keep, 🗑 for files to delete
- **Color-coded labels** - Green "Keep", Red "Delete"
- **File sizes** - See size of each copy
- **Batch selection** - Select All / Deselect All buttons

### Preview View Features

- **Multi-file navigation** - Cycle through all copies in a group
- **Side-by-side comparison** - Compare any two files at once
- **Smart decisions** - Keep Left, Keep Right, or Keep Both
- **Apply to all** - Apply current decision to remaining groups

### Generated Script Features

```bash
#!/bin/bash
# DedupeDost Cleanup Script
# Generated: 2024-03-16 10:30:00
# Source Directory: /home/user/documents
# Total files to delete: 28
# Total space to recover: 15.3 MB

# DRY_RUN mode: set to "true" to preview without deleting
DRY_RUN=${DRY_RUN:-false}

delete_file() {
    local file="$1"
    if [ "$DRY_RUN" = "true" ]; then
        echo "[DRY RUN] Would delete: $file"
    else
        if [ -f "$file" ]; then
            echo "Deleting: $file"
            rm -f "$file"
        else
            echo "Warning: File not found: $file"
        fi
    fi
}

# Group 1: vacation.mp4 (5.2 MB)
delete_file "/home/user/backup/vacation.mp4"
delete_file "/home/user/old/vacation.mp4"
# Keeping: /home/user/videos/vacation.mp4
```

## 📁 Supported File Types

| Type | Extensions | Preview Features |
|------|------------|------------------|
| **Text** | `.txt`, `.md`, `.json`, `.xml`, `.yaml`, `.yml`, `.csv`, `.log`, `.html`, `.css`, `.js`, `.ts`, `.go`, `.py`, `.java`, `.c`, `.cpp`, `.h`, `.rs`, `.rb`, `.php`, `.sh`, `.sql`, `.rst`, `.adoc` | Full content preview, monospace font, word wrap, no truncation |
| **Images** | `.jpg`, `.jpeg`, `.png`, `.gif`, `.bmp`, `.webp`, `.svg`, `.ico`, `.tiff`, `.tif`, `.raw`, `.heic`, `.heif` | Full image preview, zoom (25%-400%), fit-to-screen, dimensions, file size |
| **Videos** | `.mp4`, `.avi`, `.mkv`, `.mov`, `.wmv`, `.flv`, `.webm`, `.m4v`, `.mpg`, `.mpeg`, `.3gp` | Thumbnail preview, duration, resolution (4K/1080p/720p), codec, bitrate, FPS, **open in default player** |
| **Other** | All other files | Hash comparison, file metadata |

## 🏗️ How It Works

### Duplicate Detection Algorithm

```
1. Scan Phase
   └── Walk directory tree
   └── Collect file metadata (size, path, mod time)
   └── Apply filters (size, extension, hidden)

2. Group by Size (Quick Filter)
   └── Files with unique sizes are excluded
   └── Only files with matching sizes proceed

3. Hash Calculation
   └── SHA-256 hash for each candidate
   └── Parallel processing with worker pool
   └── Progress tracking

4. Create Duplicate Groups
   └── Group files by identical hash
   └── Calculate statistics
   └── Generate preview data
```

### Detection Speed

| Scenario | Files | Time | Speed |
|----------|-------|------|-------|
| Small files (1KB) | 1,000 | ~0.5s | 2,000 files/sec |
| Medium files (100KB) | 1,000 | ~2s | 500 files/sec |
| Mixed sizes | 10,000 | ~15s | 666 files/sec |

## 📂 Project Structure

```
dedupedost/
├── cmd/
│   ├── dedupedost/           # CLI application
│   │   └── main.go
│   └── dedupedost-gui/       # GUI application
│       └── main.go
├── internal/
│   ├── app/               # Application lifecycle
│   ├── detector/          # Duplicate detection
│   ├── hasher/            # SHA-256 hashing
│   ├── preview/           # File preview providers
│   │   ├── text_provider.go    # Full text preview (no truncation)
│   │   ├── image_provider.go   # PNG, JPEG, GIF, BMP, WebP preview
│   │   ├── video_provider.go   # Video thumbnail/metadata
│   │   └── registry.go         # Provider registry
│   ├── scanner/           # Directory scanning
│   ├── script/            # Bash script generation
│   └── ui/                # Fyne GUI components
│       ├── main_view.go
│       ├── scan_view.go
│       ├── results_view.go
│       ├── preview_view.go
│       └── script_view.go
├── pkg/
│   ├── models/            # Data models
│   └── types/             # Type definitions
├── tests/
│   ├── scanner/           # Scanner unit tests
│   ├── detector/          # Detector unit tests
│   ├── preview/           # Preview provider tests
│   └── ui/                # UI tests
│       └── screenshot/    # Screenshot tests
├── assets/
│   └── icons/             # Application icons
├── gen_test_files.sh      # Test file generator (includes images)
├── go.mod
├── go.sum
├── README.md
├── PRODUCT_SPEC.md
├── IMPLEMENTATION_PLAN.md
└── PROMPTS.md
```

## 🛠️ Development

### Prerequisites

- Go 1.21 or higher
- For GUI: Fyne dependencies (see [Fyne Installation](https://fyne.io/docs/getting-started/installation))

### Building

```bash
# Development build
go build -o dedupedost ./cmd/dedupedost

# Optimized release build
go build -ldflags="-s -w" -o dedupedost ./cmd/dedupedost

# Run tests
go test ./... -v

# Run specific tests
go test ./tests/scanner/... -v
```

### Running Tests

```bash
# All tests
go test ./... -v -cover

# Scanner tests
go test ./tests/scanner/... -v

# Detector tests
go test ./tests/detector/... -v

# Preview tests (image, text providers)
go test ./tests/preview/... -v

# UI screenshot tests (requires GUI tags)
go test ./tests/ui/screenshot/... -v -tags gui

# With coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Generate Test Files

```bash
# Generate test directory with duplicates (including images)
./gen_test_files.sh ./test_data

# This creates:
# - 16 duplicate sets (text, binary, images)
# - PNG, JPEG, GIF, BMP image duplicates
# - Unique files for comparison
# - Hidden files and special names
```

## 📋 Implementation Status

| Phase | Component | Status |
|-------|-----------|--------|
| Phase 1 | CLI Application | ✅ Complete |
| Phase 1 | Scanner | ✅ Complete |
| Phase 1 | Hasher | ✅ Complete |
| Phase 1 | Detector | ✅ Complete |
| Phase 1 | Script Generator | ✅ Complete |
| Phase 2 | GUI Foundation | ✅ Complete |
| Phase 2 | Main Window | ✅ Complete |
| Phase 2 | Scan View | ✅ Complete |
| Phase 2 | Results View | ✅ Complete |
| Phase 3 | Text Preview | ✅ Complete (full content, no truncation) |
| Phase 3 | Image Preview | ✅ Complete (PNG, JPEG, GIF, BMP, WebP) |
| Phase 3 | Video Preview | ✅ Complete (thumbnail, metadata, open in player) |
| Phase 4 | Rule System | ✅ Complete |
| Phase 4 | Batch Application | ✅ Complete |
| Phase 5 | Settings Persistence | 🔄 Pending |
| Phase 5 | Theme Support | 🔄 Pending |

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style

- Follow Go best practices
- Use `gofmt` for formatting
- Add tests for new features
- Document public APIs

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Fyne](https://fyne.io/) - Excellent Go GUI toolkit
- All contributors and supporters

## 📞 Support

- 📧 Email: support@dupdel.app
- 💬 Issues: [GitHub Issues](https://github.com/yourusername/dedupedost/issues)
- 📖 Documentation: See `PRODUCT_SPEC.md` and `IMPLEMENTATION_PLAN.md`

---

**Made with ❤️ using Go**
