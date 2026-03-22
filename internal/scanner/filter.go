package scanner

import "strings"

// Filter provides utility functions for filtering files
type Filter struct{}

// ShouldIncludeHidden checks if a hidden file should be included
func ShouldIncludeHidden(name string, includeHidden bool) bool {
	if !includeHidden && strings.HasPrefix(name, ".") {
		return false
	}
	return true
}

// ShouldIncludeBySize checks if a file size is within the allowed range
func ShouldIncludeBySize(size, minSize, maxSize int64) bool {
	return FilterBySize(size, minSize, maxSize)
}

// ShouldIncludeByExtension checks if a file extension is in the allowed list
func ShouldIncludeByExtension(ext string, allowed []string) bool {
	return FilterByExtension(ext, allowed)
}

// IsHidden checks if a file name indicates it's hidden
func IsHidden(name string) bool {
	return strings.HasPrefix(name, ".")
}

// NormalizeExtension normalizes a file extension to lowercase with dot
func NormalizeExtension(ext string) string {
	if ext == "" {
		return ""
	}
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	return strings.ToLower(ext)
}

// MatchExtension checks if a file extension matches a pattern
func MatchExtension(fileExt, pattern string) bool {
	return NormalizeExtension(fileExt) == NormalizeExtension(pattern)
}

// MatchExtensions checks if a file extension matches any pattern in the list
func MatchExtensions(fileExt string, patterns []string) bool {
	for _, pattern := range patterns {
		if MatchExtension(fileExt, pattern) {
			return true
		}
	}
	return false
}
