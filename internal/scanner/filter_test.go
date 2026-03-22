package scanner_test

import (
	"testing"

	"github.com/dupdel/dup-del/internal/scanner"
)

// TestShouldIncludeHidden tests ShouldIncludeHidden function
func TestShouldIncludeHidden(t *testing.T) {
	tests := []struct {
		name          string
		includeHidden bool
		expected      bool
	}{
		{"file.txt", true, true},
		{"file.txt", false, true},
		{".hidden", true, true},
		{".hidden", false, false},
		{".gitignore", false, false},
		{"normal", false, true},
	}

	for _, tt := range tests {
		result := scanner.ShouldIncludeHidden(tt.name, tt.includeHidden)
		if result != tt.expected {
			t.Errorf("ShouldIncludeHidden(%s, %v) = %v, want %v",
				tt.name, tt.includeHidden, result, tt.expected)
		}
	}
}

// TestShouldIncludeBySize tests ShouldIncludeBySize function
func TestShouldIncludeBySize(t *testing.T) {
	tests := []struct {
		size    int64
		minSize int64
		maxSize int64
		expected bool
	}{
		{100, 0, 0, true},       // No limits
		{100, 50, 200, true},    // Within range
		{50, 50, 200, true},     // At min boundary
		{200, 50, 200, true},    // At max boundary
		{49, 50, 200, false},    // Below min
		{201, 50, 200, false},   // Above max
		{100, 0, 50, false},     // Above max, no min
		{100, 150, 0, false},    // Below min, no max
		{0, 0, 0, true},         // Zero size, no limits
		{-1, 0, 0, false},       // Negative size is below minSize 0
	}

	for _, tt := range tests {
		result := scanner.ShouldIncludeBySize(tt.size, tt.minSize, tt.maxSize)
		if result != tt.expected {
			t.Errorf("ShouldIncludeBySize(%d, %d, %d) = %v, want %v",
				tt.size, tt.minSize, tt.maxSize, result, tt.expected)
		}
	}
}

// TestShouldIncludeByExtension tests ShouldIncludeByExtension function
func TestShouldIncludeByExtension(t *testing.T) {
	tests := []struct {
		ext      string
		allowed  []string
		expected bool
	}{
		{".txt", []string{".txt", ".md"}, true},
		{".TXT", []string{".txt", ".md"}, true}, // Case insensitive
		{".md", []string{".txt", ".md"}, true},
		{".pdf", []string{".txt", ".md"}, false},
		{".txt", []string{}, true}, // Empty allowed list = include all
		{".txt", nil, true},        // Nil allowed list = include all
		{"", []string{".txt"}, false},
		{".txt", []string{".TXT"}, true}, // Case insensitive pattern
	}

	for _, tt := range tests {
		result := scanner.ShouldIncludeByExtension(tt.ext, tt.allowed)
		if result != tt.expected {
			t.Errorf("ShouldIncludeByExtension(%s, %v) = %v, want %v",
				tt.ext, tt.allowed, result, tt.expected)
		}
	}
}

// TestIsHidden tests IsHidden function
func TestIsHidden(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{".hidden", true},
		{".gitignore", true},
		{".config", true},
		{"file.txt", false},
		{"normal", false},
		{"", false},
		{".", true},
		{"..", true},
	}

	for _, tt := range tests {
		result := scanner.IsHidden(tt.name)
		if result != tt.expected {
			t.Errorf("IsHidden(%s) = %v, want %v", tt.name, result, tt.expected)
		}
	}
}

// TestNormalizeExtension tests NormalizeExtension function
func TestNormalizeExtension(t *testing.T) {
	tests := []struct {
		ext      string
		expected string
	}{
		{".txt", ".txt"},
		{"txt", ".txt"},
		{".TXT", ".txt"},
		{"TXT", ".txt"},
		{".Md", ".md"},
		{"", ""},
		{".", "."},
		{".long", ".long"},
		{".tar.gz", ".tar.gz"},
	}

	for _, tt := range tests {
		result := scanner.NormalizeExtension(tt.ext)
		if result != tt.expected {
			t.Errorf("NormalizeExtension(%s) = %s, want %s", tt.ext, result, tt.expected)
		}
	}
}

// TestMatchExtension tests MatchExtension function
func TestMatchExtension(t *testing.T) {
	tests := []struct {
		fileExt  string
		pattern  string
		expected bool
	}{
		{".txt", ".txt", true},
		{".TXT", ".txt", true},
		{"txt", ".txt", true},
		{".txt", "txt", true},
		{".txt", ".md", false},
		{".txt", ".TXT", true},
		{"", ".txt", false},
		{".txt", "", false},
		{"", "", true},
	}

	for _, tt := range tests {
		result := scanner.MatchExtension(tt.fileExt, tt.pattern)
		if result != tt.expected {
			t.Errorf("MatchExtension(%s, %s) = %v, want %v",
				tt.fileExt, tt.pattern, result, tt.expected)
		}
	}
}

// TestMatchExtensions tests MatchExtensions function
func TestMatchExtensions(t *testing.T) {
	tests := []struct {
		fileExt  string
		patterns []string
		expected bool
	}{
		{".txt", []string{".txt", ".md"}, true},
		{".TXT", []string{".txt", ".md"}, true},
		{".md", []string{".txt", ".md"}, true},
		{".pdf", []string{".txt", ".md"}, false},
		{".txt", []string{}, false},
		{".txt", nil, false},
		{".txt", []string{".TXT"}, true},
		{".go", []string{".go", ".rs", ".py"}, true},
	}

	for _, tt := range tests {
		result := scanner.MatchExtensions(tt.fileExt, tt.patterns)
		if result != tt.expected {
			t.Errorf("MatchExtensions(%s, %v) = %v, want %v",
				tt.fileExt, tt.patterns, result, tt.expected)
		}
	}
}

// TestFilterBySize tests FilterBySize function (via ShouldIncludeBySize)
func TestFilterBySize_EdgeCases(t *testing.T) {
	// Test with maxSize = 0 (unlimited)
	if !scanner.ShouldIncludeBySize(1000000, 0, 0) {
		t.Error("Should include file when no size limits")
	}

	// Test with very large file
	if !scanner.ShouldIncludeBySize(1024*1024*1024, 0, 0) {
		t.Error("Should include very large file when no limits")
	}

	// Test with minSize only
	if scanner.ShouldIncludeBySize(50, 100, 0) {
		t.Error("Should not include file below minSize")
	}
	if !scanner.ShouldIncludeBySize(150, 100, 0) {
		t.Error("Should include file above minSize")
	}

	// Test with maxSize only
	if !scanner.ShouldIncludeBySize(50, 0, 100) {
		t.Error("Should include file below maxSize")
	}
	if scanner.ShouldIncludeBySize(150, 0, 100) {
		t.Error("Should not include file above maxSize")
	}
}

// TestFilterByExtension_CaseSensitivity tests case insensitivity
func TestFilterByExtension_CaseSensitivity(t *testing.T) {
	// All variations should match
	extensions := []string{".TXT", ".txt", ".Txt", ".tXt"}
	patterns := []string{".txt", ".TXT"}

	for _, ext := range extensions {
		for _, pattern := range patterns {
			if !scanner.ShouldIncludeByExtension(ext, []string{pattern}) {
				t.Errorf("ShouldIncludeByExtension(%s, [%s]) should be true", ext, pattern)
			}
		}
	}
}

// TestNormalizeExtension_PreservesDot tests that dot is preserved
func TestNormalizeExtension_PreservesDot(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{".tar.gz", ".tar.gz"},
		{"tar.gz", ".tar.gz"},
		{".TAR.GZ", ".tar.gz"},
	}

	for _, tt := range tests {
		result := scanner.NormalizeExtension(tt.input)
		if result != tt.expected {
			t.Errorf("NormalizeExtension(%s) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}

// TestMatchExtensions_LargeList tests with large pattern list
func TestMatchExtensions_LargeList(t *testing.T) {
	patterns := []string{
		".txt", ".md", ".json", ".xml", ".yaml", ".yml",
		".csv", ".log", ".html", ".css", ".js", ".ts",
		".go", ".py", ".java", ".c", ".cpp", ".h", ".rs",
	}

	if !scanner.MatchExtensions(".go", patterns) {
		t.Error("Should match .go in large list")
	}
	if scanner.MatchExtensions(".exe", patterns) {
		t.Error("Should not match .exe in large list")
	}
}
