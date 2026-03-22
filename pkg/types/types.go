package types

import (
	"strings"
)

// FileType represents the type of a file
type FileType int

const (
	FileTypeUnknown FileType = iota
	FileTypeText
	FileTypeImage
	FileTypeVideo
	FileTypeOther
)

// Supported text file extensions
var textExtensions = map[string]bool{
	".txt":  true,
	".md":   true,
	".json": true,
	".xml":  true,
	".yaml": true,
	".yml":  true,
	".csv":  true,
	".log":  true,
	".html": true,
	".css":  true,
	".js":   true,
	".ts":   true,
	".go":   true,
	".py":   true,
	".java": true,
	".c":    true,
	".cpp":  true,
	".h":    true,
	".rs":   true,
	".rb":   true,
	".php":  true,
	".sh":   true,
	".bash": true,
	".sql":  true,
	".rst":  true,
	".adoc": true,
}

// Supported image file extensions
var imageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".bmp":  true,
	".webp": true,
	".svg":  true,
	".ico":  true,
	".tiff": true,
	".tif":  true,
	".raw":  true,
	".heic": true,
	".heif": true,
}

// Supported video file extensions
var videoExtensions = map[string]bool{
	".mp4":  true,
	".avi":  true,
	".mkv":  true,
	".mov":  true,
	".wmv":  true,
	".flv":  true,
	".webm": true,
	".m4v":  true,
	".mpg":  true,
	".mpeg": true,
	".3gp":  true,
}

// GetFileTypeFromExtension determines the file type based on extension
func GetFileTypeFromExtension(ext string) FileType {
	ext = strings.ToLower(ext)
	
	if textExtensions[ext] {
		return FileTypeText
	}
	if imageExtensions[ext] {
		return FileTypeImage
	}
	if videoExtensions[ext] {
		return FileTypeVideo
	}
	return FileTypeOther
}

// SupportedTextExtensions returns a slice of supported text extensions
func SupportedTextExtensions() []string {
	return mapKeysToSlice(textExtensions)
}

// SupportedImageExtensions returns a slice of supported image extensions
func SupportedImageExtensions() []string {
	return mapKeysToSlice(imageExtensions)
}

// SupportedVideoExtensions returns a slice of supported video extensions
func SupportedVideoExtensions() []string {
	return mapKeysToSlice(videoExtensions)
}

// IsTextExtension checks if an extension is a text file
func IsTextExtension(ext string) bool {
	return textExtensions[strings.ToLower(ext)]
}

// IsImageExtension checks if an extension is an image file
func IsImageExtension(ext string) bool {
	return imageExtensions[strings.ToLower(ext)]
}

// IsVideoExtension checks if an extension is a video file
func IsVideoExtension(ext string) bool {
	return videoExtensions[strings.ToLower(ext)]
}

// Helper function to convert map keys to slice
func mapKeysToSlice(m map[string]bool) []string {
	result := make([]string, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}
