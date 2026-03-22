package preview

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/the100rabh/DedupeDost/internal/platform"
	"github.com/the100rabh/DedupeDost/pkg/models"
	"github.com/the100rabh/DedupeDost/pkg/types"
)

// VideoMetadata contains metadata about a video file
type VideoMetadata struct {
	Duration     time.Duration `json:"duration"`
	Width        int           `json:"width"`
	Height       int           `json:"height"`
	VideoCodec   string        `json:"video_codec"`
	AudioCodec   string        `json:"audio_codec"`
	Bitrate      int64         `json:"bitrate"`
	FrameRate    float64       `json:"frame_rate"`
	HasAudio     bool          `json:"has_audio"`
	HasVideo     bool          `json:"has_video"`
	Format       string        `json:"format"`
	FileSize     int64         `json:"file_size"`
	CreationTime time.Time     `json:"creation_time"`
}

// GetDisplayDuration returns human-readable duration
func (m *VideoMetadata) GetDisplayDuration() string {
	if m.Duration == 0 {
		return "Unknown"
	}

	hours := int(m.Duration.Hours())
	minutes := int(m.Duration.Minutes()) % 60
	seconds := int(m.Duration.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

// GetDisplayResolution returns human-readable resolution
func (m *VideoMetadata) GetDisplayResolution() string {
	if m.Width == 0 || m.Height == 0 {
		return "Unknown"
	}

	// Determine quality label
	quality := ""
	if m.Height >= 2160 {
		quality = " (4K)"
	} else if m.Height >= 1440 {
		quality = " (2K)"
	} else if m.Height >= 1080 {
		quality = " (1080p)"
	} else if m.Height >= 720 {
		quality = " (720p)"
	} else if m.Height >= 480 {
		quality = " (480p)"
	}

	return fmt.Sprintf("%dx%d%s", m.Width, m.Height, quality)
}

// GetDisplayBitrate returns human-readable bitrate
func (m *VideoMetadata) GetDisplayBitrate() string {
	if m.Bitrate == 0 {
		return "Unknown"
	}

	if m.Bitrate >= 1000000 {
		return fmt.Sprintf("%.1f Mbps", float64(m.Bitrate)/1000000)
	}
	if m.Bitrate >= 1000 {
		return fmt.Sprintf("%.1f Kbps", float64(m.Bitrate)/1000)
	}
	return fmt.Sprintf("%d bps", m.Bitrate)
}

// GetDisplayFrameRate returns human-readable frame rate
func (m *VideoMetadata) GetDisplayFrameRate() string {
	if m.FrameRate == 0 {
		return "Unknown"
	}
	return fmt.Sprintf("%.1f fps", m.FrameRate)
}

// VideoProvider provides preview for video files
type VideoProvider struct {
	ffmpegPath     string
	ffprobePath    string
	useFFmpeg      bool
	thumbnailCache map[string]interface{} // Cache thumbnail images
}

// NewVideoProvider creates a new video preview provider
func NewVideoProvider() *VideoProvider {
	vp := &VideoProvider{
		useFFmpeg:      false,
		thumbnailCache: make(map[string]interface{}),
	}

	// Check for FFmpeg/FFprobe availability
	vp.ffmpegPath = findFFmpeg()
	vp.ffprobePath = findFFprobe()
	vp.useFFmpeg = vp.ffmpegPath != "" && vp.ffprobePath != ""

	return vp
}

// findFFmpeg searches for ffmpeg executable
func findFFmpeg() string {
	// Check common locations
	paths := []string{
		"ffmpeg",
		"/usr/bin/ffmpeg",
		"/usr/local/bin/ffmpeg",
		"C:\\ffmpeg\\bin\\ffmpeg.exe",
	}

	for _, path := range paths {
		if _, err := exec.LookPath(path); err == nil {
			return path
		}
	}
	return ""
}

// findFFprobe searches for ffprobe executable
func findFFprobe() string {
	// Check common locations
	paths := []string{
		"ffprobe",
		"/usr/bin/ffprobe",
		"/usr/local/bin/ffprobe",
		"C:\\ffmpeg\\bin\\ffprobe.exe",
	}

	for _, path := range paths {
		if _, err := exec.LookPath(path); err == nil {
			return path
		}
	}
	return ""
}

// CanPreview checks if this provider can preview the file
func (vp *VideoProvider) CanPreview(fileType types.FileType, extension string) bool {
	return fileType == types.FileTypeVideo
}

// Preview creates a preview widget for the video file
func (vp *VideoProvider) Preview(file models.FileEntry) (fyne.CanvasObject, error) {
	// Get metadata
	metadata, err := vp.GetMetadata(file)
	if err != nil {
		return vp.createErrorWidget(err), nil
	}

	// Extract or get cached thumbnail
	thumbnail, err := vp.ExtractThumbnail(file.Path)
	if err != nil {
		// Use placeholder if thumbnail extraction fails
		thumbnail = vp.createPlaceholderThumbnail()
	}

	// Create thumbnail image widget
	thumbnailImg := canvas.NewImageFromImage(thumbnail)
	thumbnailImg.FillMode = canvas.ImageFillContain
	thumbnailImg.SetMinSize(fyne.NewSize(280, 160))

	// Get video metadata
	videoMeta, ok := metadata.VideoMetadata.(*VideoMetadata)

	// Create metadata labels with truncation
	var metaLabel, durationLabel, resolutionLabel, codecLabel, sizeLabel *widget.Label

	if ok && videoMeta != nil {
		// Truncate format to avoid overflow
		format := strings.ToUpper(videoMeta.Format)
		if len(format) > 10 {
			format = format[:7] + "..."
		}
		metaLabel = widget.NewLabel(fmt.Sprintf("Format: %s", format))
		metaLabel.Truncation = fyne.TextTruncateEllipsis
		
		durationLabel = widget.NewLabel(fmt.Sprintf("Duration: %s", videoMeta.GetDisplayDuration()))
		durationLabel.Truncation = fyne.TextTruncateEllipsis
		
		resolutionLabel = widget.NewLabel(fmt.Sprintf("Resolution: %s", videoMeta.GetDisplayResolution()))
		resolutionLabel.Truncation = fyne.TextTruncateEllipsis
		
		codecLabel = widget.NewLabel("")
		if videoMeta.VideoCodec != "" {
			codecLabel.SetText(fmt.Sprintf("Codec: %s", videoMeta.VideoCodec))
		}
		codecLabel.Truncation = fyne.TextTruncateEllipsis
	} else {
		metaLabel = widget.NewLabel(fmt.Sprintf("Format: %s", strings.ToUpper(file.Extension)))
		metaLabel.Truncation = fyne.TextTruncateEllipsis
		durationLabel = widget.NewLabel("Duration: Unknown")
		durationLabel.Truncation = fyne.TextTruncateEllipsis
		resolutionLabel = widget.NewLabel("Resolution: Unknown")
		resolutionLabel.Truncation = fyne.TextTruncateEllipsis
		codecLabel = widget.NewLabel("")
		codecLabel.Truncation = fyne.TextTruncateEllipsis
	}

	sizeLabel = widget.NewLabel(fmt.Sprintf("Size: %s", file.GetDisplaySize()))
	sizeLabel.Truncation = fyne.TextTruncateEllipsis

	// Open in player button - shorter text
	openBtn := widget.NewButtonWithIcon("Open", theme.MediaPlayIcon(), func() {
		vp.OpenInPlayer(file.Path)
	})

	// Create info box - compact layout
	infoBox := container.NewVBox(
		metaLabel,
		durationLabel,
		resolutionLabel,
		codecLabel,
		sizeLabel,
		widget.NewSeparator(),
		openBtn,
	)

	// Create main layout with scroll for overflow
	content := container.NewHBox(
		container.NewVBox(thumbnailImg),
		widget.NewSeparator(),
		infoBox,
		layout.NewSpacer(),
	)

	return container.NewScroll(content), nil
}

// ExtractThumbnail extracts a thumbnail from the video file
func (vp *VideoProvider) ExtractThumbnail(path string) (image.Image, error) {
	// Check cache
	if thumb, ok := vp.thumbnailCache[path]; ok {
		if img, ok := thumb.(image.Image); ok {
			return img, nil
		}
	}

	// Check if file exists
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("video file not found: %s", path)
	}

	// Try FFmpeg if available
	if vp.useFFmpeg {
		img, err := vp.extractThumbnailFFmpeg(path)
		if err == nil {
			vp.thumbnailCache[path] = img
			return img, nil
		}
	}

	// Fallback: try to read embedded thumbnail (for MP4/MKV)
	img, err := vp.extractEmbeddedThumbnail(path)
	if err == nil {
		vp.thumbnailCache[path] = img
		return img, nil
	}

	return nil, fmt.Errorf("thumbnail extraction failed")
}

// extractThumbnailFFmpeg extracts thumbnail using FFmpeg
func (vp *VideoProvider) extractThumbnailFFmpeg(path string) (image.Image, error) {
	// Create temp file for thumbnail
	tmpFile, err := os.CreateTemp("", "thumb_*.jpg")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Run FFmpeg to extract frame at 10% of duration or 5 seconds
	cmd := exec.Command(vp.ffmpegPath,
		"-i", path,
		"-ss", "00:00:05",
		"-vframes", "1",
		"-vf", "scale=320:-1",
		"-f", "image2",
		"-y",
		tmpFile.Name(),
	)

	// Suppress FFmpeg output
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// Read the thumbnail
	file, err := os.Open(tmpFile.Name())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// extractEmbeddedThumbnail tries to read embedded thumbnail from container
func (vp *VideoProvider) extractEmbeddedThumbnail(path string) (image.Image, error) {
	// For now, return a placeholder
	// Full implementation would parse MP4/MKV containers
	return vp.createPlaceholderThumbnail(), nil
}

// createPlaceholderThumbnail creates a placeholder thumbnail
func (vp *VideoProvider) createPlaceholderThumbnail() image.Image {
	// Create a simple colored rectangle as placeholder
	img := image.NewRGBA(image.Rect(0, 0, 320, 180))
	for y := 0; y < 180; y++ {
		for x := 0; x < 320; x++ {
			img.Set(x, y, color.RGBA{50, 50, 50, 255})
		}
	}

	// Draw play icon (simplified)
	centerX, centerY := 160, 90
	for y := centerY - 30; y <= centerY+30; y++ {
		for x := centerX - 20; x <= centerX+20; x++ {
			if (x-centerX)*(x-centerX)+(y-centerY)*(y-centerY) <= 30*30 {
				img.Set(x, y, color.RGBA{255, 255, 255, 200})
			}
		}
	}

	return img
}

// GetMetadata extracts metadata from the video file
func (vp *VideoProvider) GetMetadata(file models.FileEntry) (Metadata, error) {
	videoMeta, err := vp.GetVideoMetadata(file.Path)
	if err != nil {
		// Return basic metadata from file
		return Metadata{
			Properties: map[string]string{
				"Name": file.Name,
				"Size": file.GetDisplaySize(),
				"Type": "Video",
			},
		}, nil
	}

	properties := map[string]string{
		"Duration":   videoMeta.GetDisplayDuration(),
		"Resolution": videoMeta.GetDisplayResolution(),
		"Size":       file.GetDisplaySize(),
		"Format":     strings.ToUpper(videoMeta.Format),
	}

	if videoMeta.VideoCodec != "" {
		properties["Video Codec"] = videoMeta.VideoCodec
	}
	if videoMeta.AudioCodec != "" {
		properties["Audio Codec"] = videoMeta.AudioCodec
	}
	if videoMeta.Bitrate > 0 {
		properties["Bitrate"] = videoMeta.GetDisplayBitrate()
	}
	if videoMeta.FrameRate > 0 {
		properties["Frame Rate"] = videoMeta.GetDisplayFrameRate()
	}

	return Metadata{
		Properties:    properties,
		VideoMetadata: videoMeta,
	}, nil
}

// GetVideoMetadata extracts detailed video metadata using FFprobe
func (vp *VideoProvider) GetVideoMetadata(path string) (*VideoMetadata, error) {
	if !vp.useFFmpeg {
		return vp.getMetadataFallback(path)
	}

	return vp.getMetadataFFprobe(path)
}

// getMetadataFFprobe uses ffprobe to get metadata
func (vp *VideoProvider) getMetadataFFprobe(path string) (*VideoMetadata, error) {
	// Run ffprobe
	cmd := exec.Command(vp.ffprobePath,
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		path,
	)

	output, err := cmd.Output()
	if err != nil {
		return vp.getMetadataFallback(path)
	}

	// Parse JSON output
	var ffprobeOutput struct {
		Streams []struct {
			CodecType     string `json:"codec_type"`
			CodecName     string `json:"codec_name"`
			Width         int    `json:"width"`
			Height        int    `json:"height"`
			RFrameRate    string `json:"r_frame_rate"`
			BitRate       string `json:"bit_rate"`
			HasAudio      bool   `json:"has_audio"`
			Tags          struct {
				CreationTime string `json:"creation_time"`
			} `json:"tags"`
		} `json:"streams"`
		Format struct {
			Duration       string `json:"duration"`
			BitRate        string `json:"bit_rate"`
			FormatName     string `json:"format_name"`
			FormatLongName string `json:"format_long_name"`
			Size           string `json:"size"`
		} `json:"format"`
	}

	if err := json.Unmarshal(output, &ffprobeOutput); err != nil {
		return vp.getMetadataFallback(path)
	}

	// Build metadata
	meta := &VideoMetadata{
		Format: ffprobeOutput.Format.FormatName,
	}

	// Parse duration
	if ffprobeOutput.Format.Duration != "" {
		if duration, err := strconv.ParseFloat(ffprobeOutput.Format.Duration, 64); err == nil {
			meta.Duration = time.Duration(duration * float64(time.Second))
		}
	}

	// Parse bitrate
	if ffprobeOutput.Format.BitRate != "" {
		if bitrate, err := strconv.ParseInt(ffprobeOutput.Format.BitRate, 10, 64); err == nil {
			meta.Bitrate = bitrate
		}
	}

	// Parse streams
	for _, stream := range ffprobeOutput.Streams {
		switch stream.CodecType {
		case "video":
			meta.HasVideo = true
			meta.Width = stream.Width
			meta.Height = stream.Height
			meta.VideoCodec = stream.CodecName

			// Parse frame rate
			if stream.RFrameRate != "" {
				parts := strings.Split(stream.RFrameRate, "/")
				if len(parts) == 2 {
					if num, err := strconv.ParseFloat(parts[0], 64); err == nil {
						if den, err := strconv.ParseFloat(parts[1], 64); err == nil && den > 0 {
							meta.FrameRate = num / den
						}
					}
				}
			}
		case "audio":
			meta.HasAudio = true
			meta.AudioCodec = stream.CodecName
		}
	}

	return meta, nil
}

// getMetadataFallback provides basic metadata when ffprobe is unavailable
func (vp *VideoProvider) getMetadataFallback(path string) (*VideoMetadata, error) {
	// Get file info
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	meta := &VideoMetadata{
		Format:   strings.TrimPrefix(filepath.Ext(path), "."),
		FileSize: info.Size(),
	}

	// Try to get duration using ffprobe alternative or return unknown
	meta.Duration = 0 // Unknown without ffprobe

	return meta, nil
}

// OpenInPlayer opens the video file in the system's default player
func (vp *VideoProvider) OpenInPlayer(path string) error {
	return platform.OpenInDefaultPlayer(path)
}

// createErrorWidget creates an error display widget
func (vp *VideoProvider) createErrorWidget(err error) fyne.CanvasObject {
	return container.NewCenter(
		container.NewVBox(
			widget.NewIcon(theme.BrokenImageIcon()),
			widget.NewLabel("Error loading video preview"),
			widget.NewLabel(err.Error()),
		),
	)
}
