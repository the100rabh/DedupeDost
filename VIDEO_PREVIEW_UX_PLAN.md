# Video Preview Feature - UX Design Plan

## Overview
This document outlines the user experience design for video preview functionality in DedupeDost, ensuring intuitive interaction and clear visual feedback.

---

## 1. User Stories

### Primary User Stories
- **US-01**: As a user, I want to see a thumbnail of my video files so I can visually identify them
- **US-02**: As a user, I want to see video duration so I know how long each video is
- **US-03**: As a user, I want to see video resolution/quality so I can decide which copy to keep
- **US-04**: As a user, I want to open videos in my default player so I can verify content before deleting
- **US-05**: As a user, I want to compare two videos side-by-side so I can make informed decisions

### Secondary User Stories
- **US-06**: As a user, I want to see codec information so I know which video has better compression
- **US-07**: As a user, I want to see bitrate so I can assess video quality
- **US-08**: As a user, I want clear keep/delete indicators so I know my selection status

---

## 2. Information Architecture

### Video Preview Card Layout
```
┌─────────────────────────────────────────────────────────────┐
│  [THUMBNAIL]           │  Format: MP4                      │
│  320x180px             │  Duration: 2:35                   │
│  [▶ Play Icon]         │  Resolution: 1920x1080 (1080p)   │
│                        │  Codec: H.264                     │
│                        │  Size: 45.2 MB                    │
│                        │  ─────────────────────            │
│                        │  [▶ Open in Player]               │
└─────────────────────────────────────────────────────────────┘
```

### Metadata Hierarchy
1. **Primary** (always visible):
   - Format (MP4, AVI, etc.)
   - Duration
   - Resolution with quality label
   
2. **Secondary** (when available):
   - Video codec
   - Audio codec
   - Bitrate
   - Frame rate

3. **Tertiary** (file info):
   - File size
   - File path (in results view)

---

## 3. Visual Design

### 3.1 Thumbnail Design

**States:**
| State | Visual | Description |
|-------|--------|-------------|
| Default | Extracted frame or placeholder | Shows video content or gray placeholder |
| Loading | Spinner overlay | While extracting thumbnail |
| Error | Broken image icon | When extraction fails |
| Hover | Slight brightness increase | Interactive feedback |
| Selected | Green/Red border | Based on keep/delete decision |

**Placeholder Design:**
```
┌─────────────────────────────┐
│                             │
│         ▶                   │
│    (Play Icon)              │
│                             │
│     Video Preview           │
│                             │
└─────────────────────────────┘
```

### 3.2 Quality Indicators

Visual badges for quick quality recognition:

| Quality | Resolution | Badge Color | Label |
|---------|------------|-------------|-------|
| 4K | ≥2160p | Purple | `4K` |
| 2K | ≥1440p | Blue | `2K` |
| Full HD | ≥1080p | Green | `1080p` |
| HD | ≥720p | Cyan | `720p` |
| SD | ≥480p | Gray | `480p` |
| Unknown | <480p | Gray | `Unknown` |

### 3.3 Color Scheme

**Keep/Delete Indicators:**
- Keep: Green (#4CAF50)
- Delete: Red (#F44336)
- Neutral: Gray (#9E9E9E)

**Metadata Labels:**
- Primary text: Theme default
- Secondary text: Theme muted
- Quality badges: As per quality table
- Links/Actions: Theme primary

---

## 4. Interaction Design

### 4.1 Video Card Interactions

```
┌─────────────────────────────────────────────────────────┐
│  User Flow: Video Preview                               │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  1. Click video card → Opens in Preview View           │
│                                                         │
│  2. Preview View shows:                                 │
│     - Left: First video (reference)                    │
│     - Right: Current comparison video                  │
│                                                         │
│  3. User actions:                                       │
│     - Click "Open in Player" → Opens system player     │
│     - Click "Keep Left/Right" → Marks decision         │
│     - Click "Next File" → Compare next video           │
│     - Click "Apply to All" → Apply to remaining        │
│                                                         │
│  4. Visual feedback:                                    │
│     - Green checkmark on kept files                    │
│     - Red X on deleted files                           │
│     - Results view updates immediately                 │
└─────────────────────────────────────────────────────────┘
```

### 4.2 Button States

**"Open in Player" Button:**
| State | Visual | Enabled |
|-------|--------|---------|
| Default | Play icon + text | Yes |
| Hover | Slightly brighter | Yes |
| Click | Pressed effect | Yes |
| Disabled | Grayed out | No (file missing) |

**Decision Buttons:**
| State | Keep Left | Keep Right | Keep Both |
|-------|-----------|------------|-----------|
| Default | Green outline | Green outline | Blue outline |
| Selected | Solid green | Solid green | Solid blue |
| Hover | Darker green | Darker green | Darker blue |

### 4.3 Loading States

**Thumbnail Extraction:**
```
┌─────────────────────────────┐
│                             │
│      ⟳ Loading...          │
│    (Spinning indicator)     │
│                             │
└─────────────────────────────┘
```

**Metadata Loading:**
- Show skeleton loaders for text fields
- Display "Loading..." placeholder
- Fade in actual data when ready

---

## 5. Accessibility

### 5.1 Keyboard Navigation

| Key | Action |
|-----|--------|
| Tab | Navigate between elements |
| Enter | Activate selected button |
| Space | Toggle checkbox (Apply to All) |
| Arrow Left/Right | Navigate between files |
| Escape | Close preview / Cancel action |

### 5.2 Screen Reader Support

**ARIA Labels:**
```html
<!-- Thumbnail -->
<img aria-label="Video thumbnail: vacation.mp4, 2 minutes 35 seconds, 1080p">

<!-- Metadata -->
<div aria-label="Format: MP4">
<div aria-label="Duration: 2 minutes 35 seconds">
<div aria-label="Resolution: 1920 by 1080, Full HD">

<!-- Buttons -->
<button aria-label="Open vacation.mp4 in default video player">
<button aria-label="Keep left video, delete right video">
```

### 5.3 Color Contrast

- Text on background: ≥ 4.5:1 (WCAG AA)
- Interactive elements: ≥ 3:1
- Quality badges: White text on colored background

---

## 6. Error States

### 6.1 Thumbnail Extraction Failed

```
┌─────────────────────────────┐
│        ⚠️                 │
│   Preview unavailable       │
│                             │
│   [▶ Open in Player]       │
└─────────────────────────────┘
```

**User Action:** Click "Open in Player" to verify content manually

### 6.2 Metadata Extraction Failed

```
Format: MP4
Duration: Unknown
Resolution: Unknown
Codec: Unknown
Size: 45.2 MB
─────────────────────
[▶ Open in Player]
```

**User Action:** Rely on file size and manual inspection

### 6.3 No Default Player

**Dialog:**
```
┌─────────────────────────────────────┐
│  ⚠️  No Default Video Player       │
├─────────────────────────────────────┤
│  No default video player is set    │
│  on your system.                   │
│                                    │
│  Please install a video player     │
│  (VLC, MPV, etc.) and try again.  │
│                                    │
│  [OK]  [Help]                      │
└─────────────────────────────────────┘
```

---

## 7. Responsive Design

### 7.1 Layout Breakpoints

| Screen Width | Layout |
|--------------|--------|
| ≥1200px | Side-by-side, full metadata |
| 800-1199px | Side-by-side, condensed metadata |
| <800px | Stacked (thumbnail above metadata) |

### 7.2 Mobile Considerations

**Future Enhancement:**
- Touch-optimized buttons (min 44x44px)
- Swipe gestures for navigation
- Simplified metadata display

---

## 8. Performance Guidelines

### 8.1 Thumbnail Loading

| Scenario | Target | Fallback |
|----------|--------|----------|
| FFmpeg available | <2s | Placeholder |
| No FFmpeg | <0.5s | Placeholder |
| Cached thumbnail | <0.1s | N/A |

### 8.2 Metadata Loading

| Scenario | Target |
|----------|--------|
| FFprobe available | <1s |
| No FFprobe | <0.5s (basic info) |

### 8.3 Memory Management

- Thumbnail cache: Max 50 images
- Cache eviction: LRU (Least Recently Used)
- Large video handling: Skip thumbnail for files >2GB

---

## 9. User Testing Plan

### 9.1 Test Scenarios

1. **Basic Comparison**
   - User has 3 duplicate videos
   - Task: Identify and keep highest quality

2. **Quality Decision**
   - Videos with different resolutions (4K, 1080p, 720p)
   - Task: Keep 4K version

3. **Codec Comparison**
   - Same resolution, different codecs (H.264 vs HEVC)
   - Task: Keep better compression

4. **Manual Verification**
   - Similar filenames, uncertain content
   - Task: Use "Open in Player" to verify

### 9.2 Success Metrics

| Metric | Target |
|--------|--------|
| Task completion rate | >95% |
| Average decision time | <30s per group |
| Error rate | <5% |
| User satisfaction | >4.5/5 |

---

## 10. Implementation Checklist

### Phase 1: Core Features
- [x] Video provider implementation
- [x] Thumbnail extraction
- [x] Metadata display
- [x] Open in player button

### Phase 2: UX Polish
- [x] Quality indicators (4K/1080p/etc.)
- [x] Loading states
- [x] Error handling
- [ ] Thumbnail loading spinner
- [ ] Skeleton loaders for metadata

### Phase 3: Accessibility
- [ ] ARIA labels
- [ ] Keyboard navigation
- [ ] Screen reader testing
- [ ] Color contrast verification

### Phase 4: Performance
- [x] Thumbnail caching
- [ ] Lazy loading for thumbnails
- [ ] Progress indicators for large files
- [ ] Memory usage optimization

---

## 11. Design Assets

### Icons Required
- `theme.MediaPlayIcon()` - Open in player
- `theme.VideoIcon()` - Video file type
- `theme.BrokenImageIcon()` - Error state
- Custom: Quality badges (4K, 1080p, etc.)

### Color Tokens
```go
// Quality badge colors
Color4K    = color.RGBA{128, 0, 128, 255}   // Purple
Color2K    = color.RGBA{0, 0, 255, 255}     // Blue
Color1080p = color.RGBA{0, 255, 0, 255}     // Green
Color720p  = color.RGBA{0, 255, 255, 255}   // Cyan
Color480p  = color.RGBA{128, 128, 128, 255} // Gray
```

---

## 12. Future Enhancements

### 12.1 Short-term (v1.1)
- [ ] FFmpeg bundling for consistent thumbnail extraction
- [ ] Batch thumbnail generation during scan
- [ ] Video duration in results view

### 12.2 Medium-term (v1.2)
- [ ] Timeline scrubber for manual thumbnail selection
- [ ] Multiple thumbnail preview
- [ ] Video content fingerprinting

### 12.3 Long-term (v2.0)
- [ ] Built-in video player (embedded)
- [ ] Frame-by-frame comparison
- [ ] Audio waveform visualization

---

**Document Version**: 1.0.0
**Last Updated**: 2024-03-20
**Author**: DedupeDost Team
