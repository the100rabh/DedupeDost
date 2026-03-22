# DupDel - Phase 3 & 4 Implementation Prompts

## Phase 3: Enhanced Preview System

### Task 3.1: Text Preview Provider with Syntax Highlighting

```
Create an enhanced text preview provider with syntax highlighting.

File: internal/preview/text_provider.go

Requirements:
1. Create TextProvider struct that implements PreviewProvider interface
2. Support reading text files with encoding detection (UTF-8, ASCII, Latin-1)
3. Implement syntax highlighting for common languages:
   - Go, Python, JavaScript, TypeScript, Java, C, C++, Rust
   - HTML, CSS, JSON, YAML, XML, SQL, Shell
4. Show line numbers
5. Configurable max lines (default 100)
6. Word wrap option
7. Search within text

Implementation:
- Use regexp for basic syntax highlighting
- Map file extensions to languages
- Return fyne.CanvasObject with colored text
- Handle large files efficiently (read only visible portion)

Output: Complete text preview provider with syntax highlighting.
```

### Task 3.2: Image Preview Provider

```
Create an image preview provider with zoom and pan.

File: internal/preview/image_provider.go

Requirements:
1. Create ImageProvider struct
2. Support formats: JPG, PNG, GIF, BMP, WebP, SVG, ICO, TIFF
3. Implement zoom (25% - 400%)
4. Implement pan (drag to move)
5. Show image dimensions and file size
6. Extract and display EXIF metadata when available:
   - Date taken
   - Camera model
   - Resolution
   - Orientation
7. Generate thumbnails for quick preview
8. Cache loaded images

Implementation:
- Use golang.org/x/image for decoding
- Use exif package for metadata
- Implement mouse wheel zoom
- Implement drag-to-pan
- Show zoom level indicator

Output: Complete image preview provider with zoom/pan.
```

### Task 3.3: Video Preview Provider

```
Create a video preview provider with thumbnail extraction.

File: internal/preview/video_provider.go

Requirements:
1. Create VideoProvider struct
2. Support formats: MP4, AVI, MKV, MOV, WMV, FLV, WebM
3. Extract thumbnail from video:
   - Try FFmpeg first (if available)
   - Try embedded thumbnail
   - Fallback to placeholder with duration
4. Display video metadata:
   - Duration
   - Resolution (width x height)
   - Codec (H.264, HEVC, VP9, etc.)
   - Bitrate
   - Frame rate
   - Audio codec
5. Show file path and size

Implementation:
- Use exec.Command for FFmpeg when available
- Parse video metadata using ffprobe or go-mp4
- Display thumbnail with overlay duration
- Handle missing codecs gracefully

Output: Complete video preview provider.
```

### Task 3.4: Preview Provider Registry

```
Create a preview provider registry system.

File: internal/preview/registry.go

Requirements:
1. Define PreviewProvider interface:
   - CanPreview(fileType, extension) bool
   - Preview(file FileEntry) (CanvasObject, error)
   - GetMetadata(file FileEntry) (Metadata, error)
2. Create Registry struct to manage providers
3. Register providers for file types
4. Auto-detect best provider for file
5. Fallback provider for unsupported types

Implementation:
- Singleton pattern for registry
- Thread-safe registration
- Priority-based provider selection
- Default fallback with file icon

Output: Provider registry with auto-detection.
```

### Task 3.5: Side-by-Side Comparison View

```
Create enhanced side-by-side comparison.

File: internal/ui/compare_view.go

Requirements:
1. Create CompareView struct
2. Synchronized scrolling for text files
3. Synchronized zoom for images
4. Visual diff highlighting for text:
   - Green for identical lines
   - Red for different content
   - Yellow for whitespace differences
5. Show file metadata comparison table
6. Quick navigation between differences

Implementation:
- Use container.Split for side-by-side
- Link scroll offsets for text
- Link zoom levels for images
- Implement line-by-line comparison
- Highlight differences

Output: Side-by-side comparison with sync features.
```

---

## Phase 4: Rule System

### Task 4.1: Built-in Rules Definition

```
Create built-in rule definitions.

File: internal/rules/builtin_rules.go

Requirements:
1. Define BuiltinRules slice with these rules:
   - Keep Newest (by modification time)
   - Keep Oldest (by modification time)
   - Keep Shortest Path (fewest directory levels)
   - Keep Longest Path (most specific location)
   - Keep from Specific Directory (user-configured)
   - Keep Largest (for different quality versions)
   - Keep Smallest (for space efficiency)
   - Manual Selection (user chooses)

2. Each rule has:
   - ID (unique identifier)
   - Name (display name)
   - Description (what it does)
   - Icon (fyne resource)
   - Criteria (enum)
   - Config schema (for configurable rules)

Output: Complete built-in rules definitions.
```

### Task 4.2: Rule Engine

```
Create the rule application engine.

File: internal/rules/engine.go

Requirements:
1. Create RuleEngine struct
2. Implement ApplyRule(group, ruleID) method
3. Implement ApplyToAll(groups, ruleID) for batch application
4. Register evaluators for each criteria type
5. Support undo/redo of rule applications
6. Track which rules were applied to which groups

Implementation:
- Strategy pattern for rule evaluators
- Command pattern for undo/redo
- Event system for rule application notifications
- Thread-safe for concurrent application

Output: Rule engine with all evaluators.
```

### Task 4.3: Rule Evaluators

```
Create rule evaluators for each criteria.

File: internal/rules/evaluators.go

Requirements:
1. Create evaluator for each RuleCriteria:
   - NewestEvaluator: Sort by ModTime descending
   - OldestEvaluator: Sort by ModTime ascending
   - ShortestPathEvaluator: Sort by path depth
   - LongestPathEvaluator: Sort by path depth descending
   - SpecificDirEvaluator: Sort by directory match
   - LargestEvaluator: Sort by size descending
   - SmallestEvaluator: Sort by size ascending

2. Each evaluator implements:
   - Evaluate(group) (keepIndex, error)
   - Description() string
   - IsApplicable(group) bool

Output: All rule evaluators implemented.
```

### Task 4.4: Rule Selection Dialog

```
Create rule selection dialog UI.

File: internal/ui/rule_dialog.go

Requirements:
1. Create RuleDialog struct
2. Display list of available rules with icons
3. Show rule preview/description
4. Show which files would be kept/deleted
5. "Apply to All" checkbox
6. Configuration panel for configurable rules
7. Apply and Cancel buttons

Implementation:
- Use widget.List for rule selection
- Preview panel shows example application
- Dynamic config panel based on selected rule
- Confirmation before batch application

Output: Rule selection dialog.
```

### Task 4.5: Batch Rule Application

```
Implement batch rule application.

File: internal/rules/batch_applier.go

Requirements:
1. Create BatchApplier struct
2. Apply rule to multiple groups
3. Track progress during batch operation
4. Support cancellation
5. Generate report of applied rules
6. Support undo of batch operation

Implementation:
- Progress callbacks for UI updates
- Context-based cancellation
- Result aggregation
- Undo stack for reversibility

Output: Batch rule application with progress.
```

### Task 4.6: Decision Review Panel

```
Create decision review panel.

File: internal/ui/review_panel.go

Requirements:
1. Create ReviewPanel struct
2. Display all decisions in a list
3. Filter by decision type (keep/delete/skip)
4. Edit individual decisions
5. Show summary statistics
6. Export decisions to CSV/JSON

Implementation:
- Use widget.List for decisions
- Inline editing for decisions
- Summary with counts and sizes
- Export functionality

Output: Decision review panel.
```

### Task 4.7: Smart Auto-Selection

```
Create smart auto-selection features.

File: internal/rules/smart_select.go

Requirements:
1. Auto-select based on common patterns:
   - Keep files in "Documents", "Photos" folders
   - Delete files in "Backup", "Old", "Temp" folders
   - Keep files with cleaner names (no "copy", "duplicate")
   - Keep files with better quality (higher resolution for images)
2. Confidence scoring for suggestions
3. User can accept/reject suggestions
4. Learn from user decisions (optional)

Implementation:
- Pattern matching for folder names
- Filename analysis
- Quality metrics for media files
- Scoring system

Output: Smart auto-selection system.
```

---

## Usage

Copy each prompt and use it to generate the implementation code.
Follow the order for proper dependencies.
