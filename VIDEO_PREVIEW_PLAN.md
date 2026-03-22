# Video Preview Implementation Plan

**Status**: ✅ **COMPLETE**

## Overview
Add side-by-side video preview with thumbnail extraction, metadata display, and open-in-player functionality.

## Implementation Summary

All features have been implemented:
- ✅ Video thumbnail extraction (FFmpeg or placeholder fallback)
- ✅ Video metadata display (duration, resolution, codec, bitrate, FPS)
- ✅ Open in default player (cross-platform)
- ✅ Side-by-side video comparison
- ✅ Quality indicators (4K/1080p/720p/480p)
- ✅ Thumbnail caching for performance
- ✅ Unit tests for video provider

## Features

### 1. Video Thumbnail Extraction
- Extract frame from video file as preview thumbnail
- Support multiple formats: MP4, AVI, MKV, MOV, WebM, WMV, FLV
- Cache thumbnails for performance
- Fallback to file icon if thumbnail extraction fails

### 2. Video Metadata Display
- Duration (length)
- Resolution (width x height)
- Codec (H.264, HEVC, VP9, etc.)
- Bitrate
- Frame rate (FPS)
- File size
- Creation/modification date

### 3. Open in Default Player
- Button to open video in system's default video player
- Cross-platform support:
  - Windows: Default app for .mp4/.avi/etc.
  - macOS: QuickTime Player or default
  - Linux: VLC, Totem, or xdg-utils

### 4. Side-by-Side Comparison
- Same layout as image/text preview
- Compare two videos at once
- Keep/Delete decisions
- Navigate through multiple video copies

## Implementation Tasks

### Task 1: Video Preview Provider
**File**: `internal/preview/video_provider.go`

```go
type VideoProvider struct {
    thumbnailCache map[string]image.Image
    useFFmpeg      bool
    ffmpegPath     string
}

func (v *VideoProvider) Preview(file models.FileEntry) (fyne.CanvasObject, error)
func (v *VideoProvider) ExtractThumbnail(path string) (image.Image, error)
func (v *VideoProvider) GetMetadata(path string) (VideoMetadata, error)
```

### Task 2: Video Metadata Structure
**File**: `internal/preview/video_metadata.go`

```go
type VideoMetadata struct {
    Duration   time.Duration
    Width      int
    Height     int
    Codec      string
    Bitrate    int64
    FrameRate  float64
    HasAudio   bool
}
```

### Task 3: Thumbnail Extraction Methods
1. **FFmpeg** (preferred, if available)
   ```bash
   ffmpeg -i video.mp4 -ss 00:00:01 -vframes 1 thumb.jpg
   ```
2. **go-ffmpeg bindings** (optional dependency)
3. **GStreamer** (Linux fallback)
4. **Built-in decoders** (limited format support)
5. **File icon fallback**

### Task 4: Metadata Extraction
1. **FFprobe** (FFmpeg companion tool)
   ```bash
   ffprobe -v quiet -print_format json -show_format -show_streams video.mp4
   ```
2. **go-ffprobe** Go wrapper
3. **Pure Go parsers** for MP4/MKV containers

### Task 5: UI Components
**File**: `internal/ui/video_preview_card.go`

```go
type VideoPreviewCard struct {
    container       *fyne.Container
    thumbnail       *canvas.Image
    metadataLabel   *widget.Label
    durationLabel   *widget.Label
    resolutionLabel *widget.Label
    openButton      *widget.Button
    fileEntry       models.FileEntry
}

func NewVideoPreviewCard() *VideoPreviewCard
func (v *VideoPreviewCard) SetFile(file models.FileEntry)
func (v *VideoPreviewCard) Clear()
```

### Task 6: Open in Default Player
**File**: `internal/platform/open_video.go`

```go
func OpenInDefaultPlayer(filePath string) error

// Platform-specific implementations:
// - Windows: exec.Command("cmd", "/c", "start", filePath)
// - macOS:   exec.Command("open", filePath)
// - Linux:   exec.Command("xdg-open", filePath)
```

### Task 7: Integration with PreviewView
- Detect video file type
- Use VideoPreviewCard instead of generic FilePreviewCard
- Add "Open in Player" button to action panel
- Show thumbnail in side-by-side view

## Dependencies

### Required
- None (pure Go implementation with fallbacks)

### Optional (enhanced features)
- `github.com/u2takey/ffmpeg-go` - FFmpeg Go bindings
- `github.com/vansante/go-ffprobe` - FFprobe Go wrapper
- `golang.org/x/image` - Image processing

## Implementation Phases

### Phase 1: Basic Implementation (2 days)
- [ ] Video provider skeleton
- [ ] File icon fallback
- [ ] Metadata display (basic: duration, resolution)
- [ ] Open in player button
- [ ] Integration with PreviewView

### Phase 2: Thumbnail Extraction (2 days)
- [ ] FFmpeg detection
- [ ] Thumbnail extraction with FFmpeg
- [ ] Thumbnail caching
- [ ] Error handling and fallbacks

### Phase 3: Enhanced Metadata (1 day)
- [ ] FFprobe integration
- [ ] Codec detection
- [ ] Bitrate/FPS display
- [ ] Audio track info

### Phase 4: Testing & Polish (1 day)
- [ ] Unit tests for video provider
- [ ] Test with various video formats
- [ ] Cross-platform testing
- [ ] Documentation

## File Structure

```
internal/
├── preview/
│   ├── video_provider.go      # New: Video preview logic
│   ├── video_metadata.go      # New: Metadata structures
│   ├── video_thumbnail.go     # New: Thumbnail extraction
│   └── registry.go            # Modified: Register video provider
├── platform/
│   └── open_video.go          # New: Open in default player
ui/
├── video_preview_card.go      # New: Video card widget
└── preview_view.go            # Modified: Support video cards
tests/
└── preview/
    └── video_provider_test.go # New: Video tests
```

## API Design

### VideoProvider Interface
```go
type VideoProvider interface {
    PreviewProvider
    ExtractThumbnail(path string) (image.Image, error)
    GetVideoMetadata(path string) (VideoMetadata, error)
    OpenInPlayer(path string) error
}
```

### VideoPreviewCard API
```go
func NewVideoPreviewCard() *VideoPreviewCard
func (v *VideoPreviewCard) SetFile(file models.FileEntry)
func (v *VideoPreviewCard) Clear()
func (v *VideoPreviewCard) SetKeep(keep bool)
```

## Error Handling

1. **Thumbnail extraction fails**: Show file icon with overlay
2. **Metadata extraction fails**: Show basic file info (size, name)
3. **No default player**: Show error dialog with message
4. **Corrupt video file**: Show error in preview area

## Performance Considerations

1. **Thumbnail caching**: Store extracted thumbnails in memory
2. **Lazy loading**: Extract thumbnails only when card is visible
3. **Async extraction**: Don't block UI during thumbnail generation
4. **Size limits**: Skip thumbnail for videos > 2GB (too slow)

## Testing Strategy

1. **Unit Tests**:
   - VideoProvider.CanPreview() for various extensions
   - VideoMetadata parsing
   - OpenInDefaultPlayer() path handling

2. **Integration Tests**:
   - Preview creation with sample video files
   - Thumbnail extraction (if FFmpeg available)
   - Metadata extraction accuracy

3. **Manual Testing**:
   - Test with real video files (MP4, AVI, MKV)
   - Cross-platform player opening
   - Various resolutions and codecs

## Success Criteria

- [ ] Video files show thumbnail preview
- [ ] Metadata displays correctly (duration, resolution, codec)
- [ ] "Open in Player" button works on all platforms
- [ ] Side-by-side comparison works for video duplicates
- [ ] Performance acceptable (<2s thumbnail extraction)
- [ ] Graceful fallback when FFmpeg not available
