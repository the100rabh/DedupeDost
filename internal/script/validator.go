package script

import (
	"fmt"
	"os"
	"path/filepath"
)

// Validator validates generated scripts
type Validator struct {
	warnings []ValidationWarning
	errors   []ValidationError
}

// ValidationWarning represents a non-critical issue
type ValidationWarning struct {
	Path    string
	Message string
}

// ValidationError represents a critical issue
type ValidationError struct {
	Path    string
	Message string
	Fatal   bool
}

// Error implements the error interface for ValidationError
func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error at %s: %s", e.Path, e.Message)
}

// NewValidator creates a new Validator
func NewValidator() *Validator {
	return &Validator{
		warnings: make([]ValidationWarning, 0),
		errors:   make([]ValidationError, 0),
	}
}

// Validate validates a script file
func (v *Validator) Validate(scriptPath string) error {
	v.warnings = make([]ValidationWarning, 0)
	v.errors = make([]ValidationError, 0)
	
	// Check file exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		v.errors = append(v.errors, ValidationError{
			Path:    scriptPath,
			Message: "Script file does not exist",
			Fatal:   true,
		})
		return v.errors[0]
	}
	
	// Check file is readable
	content, err := os.ReadFile(scriptPath)
	if err != nil {
		v.errors = append(v.errors, ValidationError{
			Path:    scriptPath,
			Message: "Cannot read script file",
			Fatal:   true,
		})
		return v.errors[0]
	}
	
	// Check for shebang
	if len(content) < 2 || string(content[:2]) != "#!" {
		v.warnings = append(v.warnings, ValidationWarning{
			Path:    scriptPath,
			Message: "Script missing shebang line",
		})
	}
	
	// Check for rm commands
	if !containsString(string(content), "rm -f") {
		v.warnings = append(v.warnings, ValidationWarning{
			Path:    scriptPath,
			Message: "Script contains no rm commands",
		})
	}
	
	// Check file permissions
	info, err := os.Stat(scriptPath)
	if err == nil {
		if info.Mode()&0111 == 0 {
			v.warnings = append(v.warnings, ValidationWarning{
				Path:    scriptPath,
				Message: "Script is not executable",
			})
		}
	}
	
	if len(v.errors) > 0 {
		return v.errors[0]
	}
	
	return nil
}

// GetWarnings returns validation warnings
func (v *Validator) GetWarnings() []ValidationWarning {
	return v.warnings
}

// GetErrors returns validation errors
func (v *Validator) GetErrors() []ValidationError {
	return v.errors
}

// ValidatePaths checks if paths in the script exist
func (v *Validator) ValidatePaths(scriptPath string) {
	content, err := os.ReadFile(scriptPath)
	if err != nil {
		return
	}
	
	// Simple check - look for paths in rm commands
	lines := splitLines(string(content))
	for _, line := range lines {
		if containsString(line, "rm -f") || containsString(line, "delete_file") {
			// Extract path from command
			path := extractPath(line)
			if path != "" {
				if _, err := os.Stat(path); os.IsNotExist(err) {
					v.warnings = append(v.warnings, ValidationWarning{
						Path:    path,
						Message: "File no longer exists",
					})
				}
			}
		}
	}
}

// Helper functions
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func extractPath(line string) string {
	// Find path between quotes
	start := -1
	end := -1
	for i, c := range line {
		if c == '"' {
			if start == -1 {
				start = i + 1
			} else {
				end = i
				break
			}
		}
	}
	
	if start != -1 && end != -1 && end > start {
		return line[start:end]
	}
	
	return ""
}

// GetOutputPath returns the recommended output path for a script
func GetOutputPath(sourceDir string) string {
	baseName := "cleanup_duplicates.sh"
	
	// Try current directory first
	cwd, err := os.Getwd()
	if err == nil {
		return filepath.Join(cwd, baseName)
	}
	
	// Fallback to source directory
	return filepath.Join(sourceDir, baseName)
}
