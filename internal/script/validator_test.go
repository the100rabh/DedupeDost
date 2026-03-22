package script_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/the100rabh/DedupeDost/internal/script"
)

// TestNewValidator tests NewValidator
func TestNewValidator(t *testing.T) {
	v := script.NewValidator()
	if v == nil {
		t.Fatal("NewValidator returned nil")
	}

	warnings := v.GetWarnings()
	if warnings == nil {
		t.Error("GetWarnings should return empty slice, not nil")
	}
	if len(warnings) != 0 {
		t.Errorf("Initial warnings length = %d, want 0", len(warnings))
	}

	errors := v.GetErrors()
	if errors == nil {
		t.Error("GetErrors should return empty slice, not nil")
	}
	if len(errors) != 0 {
		t.Errorf("Initial errors length = %d, want 0", len(errors))
	}
}

// TestValidator_Validate tests Validate
func TestValidator_Validate(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.sh")

	// Create a valid script
	validScript := `#!/bin/bash
rm -f "/tmp/test.txt"
`
	if err := os.WriteFile(scriptPath, []byte(validScript), 0755); err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	v := script.NewValidator()
	err := v.Validate(scriptPath)
	if err != nil {
		t.Errorf("Validate should not return error for valid script: %v", err)
	}

	warnings := v.GetWarnings()
	if len(warnings) != 0 {
		t.Errorf("Expected 0 warnings for valid script, got %d", len(warnings))
	}
}

// TestValidator_Validate_NonExistent tests Validate with non-existent file
func TestValidator_Validate_NonExistent(t *testing.T) {
	v := script.NewValidator()
	err := v.Validate("/nonexistent/script.sh")
	if err == nil {
		t.Error("Validate should return error for non-existent file")
	}

	errors := v.GetErrors()
	if len(errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(errors))
	}
	if !errors[0].Fatal {
		t.Error("Error should be fatal")
	}
}

// TestValidator_Validate_NoShebang tests Validate with missing shebang
func TestValidator_Validate_NoShebang(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.sh")

	// Create script without shebang
	noShebangScript := `echo "Hello"
rm -f "/tmp/test.txt"
`
	if err := os.WriteFile(scriptPath, []byte(noShebangScript), 0644); err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	v := script.NewValidator()
	err := v.Validate(scriptPath)
	if err != nil {
		t.Errorf("Validate should not return error: %v", err)
	}

	warnings := v.GetWarnings()
	foundShebangWarning := false
	for _, w := range warnings {
		if w.Message == "Script missing shebang line" {
			foundShebangWarning = true
			break
		}
	}
	if !foundShebangWarning {
		t.Error("Should warn about missing shebang")
	}
}

// TestValidator_Validate_NoRmCommands tests Validate with no rm commands
func TestValidator_Validate_NoRmCommands(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.sh")

	// Create script without rm commands
	noRmScript := `#!/bin/bash
echo "Hello"
`
	if err := os.WriteFile(scriptPath, []byte(noRmScript), 0755); err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	v := script.NewValidator()
	err := v.Validate(scriptPath)
	if err != nil {
		t.Errorf("Validate should not return error: %v", err)
	}

	warnings := v.GetWarnings()
	foundRmWarning := false
	for _, w := range warnings {
		if w.Message == "Script contains no rm commands" {
			foundRmWarning = true
			break
		}
	}
	if !foundRmWarning {
		t.Error("Should warn about no rm commands")
	}
}

// TestValidator_Validate_NotExecutable tests Validate with non-executable file
func TestValidator_Validate_NotExecutable(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.sh")

	validScript := `#!/bin/bash
rm -f "/tmp/test.txt"
`
	if err := os.WriteFile(scriptPath, []byte(validScript), 0644); err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	v := script.NewValidator()
	err := v.Validate(scriptPath)
	if err != nil {
		t.Errorf("Validate should not return error: %v", err)
	}

	warnings := v.GetWarnings()
	foundExecWarning := false
	for _, w := range warnings {
		if w.Message == "Script is not executable" {
			foundExecWarning = true
			break
		}
	}
	if !foundExecWarning {
		t.Error("Should warn about not executable")
	}
}

// TestValidationError_Error tests ValidationError.Error
func TestValidationError_Error(t *testing.T) {
	err := script.ValidationError{
		Path:    "/test/path.sh",
		Message: "test error",
		Fatal:   true,
	}

	expected := "validation error at /test/path.sh: test error"
	if err.Error() != expected {
		t.Errorf("Error() = %s, want %s", err.Error(), expected)
	}
}

// TestValidationWarning tests ValidationWarning struct
func TestValidationWarning(t *testing.T) {
	warning := script.ValidationWarning{
		Path:    "/test/path.sh",
		Message: "test warning",
	}

	if warning.Path != "/test/path.sh" {
		t.Errorf("Path = %s, want /test/path.sh", warning.Path)
	}
	if warning.Message != "test warning" {
		t.Errorf("Message = %s, want test warning", warning.Message)
	}
}

// TestValidationError tests ValidationError struct
func TestValidationError(t *testing.T) {
	err := script.ValidationError{
		Path:    "/test/path.sh",
		Message: "test error",
		Fatal:   true,
	}

	if err.Path != "/test/path.sh" {
		t.Errorf("Path = %s, want /test/path.sh", err.Path)
	}
	if err.Message != "test error" {
		t.Errorf("Message = %s, want test error", err.Message)
	}
	if !err.Fatal {
		t.Error("Fatal should be true")
	}
}

// TestValidator_GetWarnings tests GetWarnings
func TestValidator_GetWarnings(t *testing.T) {
	v := script.NewValidator()

	// Initially empty
	warnings := v.GetWarnings()
	if len(warnings) != 0 {
		t.Errorf("Initial warnings = %d, want 0", len(warnings))
	}

	// Create a script that generates warnings
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.sh")
	if err := os.WriteFile(scriptPath, []byte("echo test"), 0644); err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	v.Validate(scriptPath)
	warnings = v.GetWarnings()
	if len(warnings) == 0 {
		t.Error("Expected warnings after validation")
	}
}

// TestValidator_GetErrors tests GetErrors
func TestValidator_GetErrors(t *testing.T) {
	v := script.NewValidator()

	// Initially empty
	errors := v.GetErrors()
	if len(errors) != 0 {
		t.Errorf("Initial errors = %d, want 0", len(errors))
	}

	// Create a non-existent path
	v.Validate("/nonexistent/file.sh")
	errors = v.GetErrors()
	if len(errors) == 0 {
		t.Error("Expected errors after validation")
	}
}

// TestValidator_ValidatePaths tests ValidatePaths
func TestValidator_ValidatePaths(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.sh")

	// Create script with valid path
	tmpFile := filepath.Join(tmpDir, "to_delete.txt")
	if err := os.WriteFile(tmpFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	validScript := `#!/bin/bash
rm -f "` + tmpFile + `"
`
	if err := os.WriteFile(scriptPath, []byte(validScript), 0755); err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	v := script.NewValidator()
	v.ValidatePaths(scriptPath)

	warnings := v.GetWarnings()
	// File exists, so no warnings about missing files
	for _, w := range warnings {
		if w.Message == "File no longer exists" {
			t.Error("Should not warn about existing file")
		}
	}
}

// TestValidator_ValidatePaths_NonExistentFile tests ValidatePaths with non-existent files
func TestValidator_ValidatePaths_NonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.sh")

	// Create script with non-existent path
	validScript := `#!/bin/bash
rm -f "/nonexistent/file.txt"
`
	if err := os.WriteFile(scriptPath, []byte(validScript), 0755); err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	v := script.NewValidator()
	v.ValidatePaths(scriptPath)

	warnings := v.GetWarnings()
	foundMissingWarning := false
	for _, w := range warnings {
		if w.Message == "File no longer exists" {
			foundMissingWarning = true
			break
		}
	}
	if !foundMissingWarning {
		t.Error("Should warn about non-existent file")
	}
}

// TestValidator_ValidatePaths_ReadError tests ValidatePaths with read error
func TestValidator_ValidatePaths_ReadError(t *testing.T) {
	v := script.NewValidator()
	// Non-existent script file
	v.ValidatePaths("/nonexistent/script.sh")

	// Should not panic, just return silently
	warnings := v.GetWarnings()
	if len(warnings) != 0 {
		t.Errorf("Expected 0 warnings for unreadable script, got %d", len(warnings))
	}
}

// TestHelperFunctions tests helper functions indirectly
func TestHelperFunctions(t *testing.T) {
	// Test containsString via validation
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.sh")

	scriptContent := `#!/bin/bash
rm -f "/tmp/file.txt"
delete_file "/tmp/file2.txt"
`
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	v := script.NewValidator()
	v.ValidatePaths(scriptPath)

	// Should find both paths
	warnings := v.GetWarnings()
	// Both files don't exist, so should have 2 warnings
	missingCount := 0
	for _, w := range warnings {
		if w.Message == "File no longer exists" {
			missingCount++
		}
	}
	if missingCount != 2 {
		t.Errorf("Expected 2 missing file warnings, got %d", missingCount)
	}
}

// TestGetOutputPath tests GetOutputPath
func TestGetOutputPath(t *testing.T) {
	path := script.GetOutputPath("/test/dir")
	if path == "" {
		t.Error("GetOutputPath should return non-empty string")
	}

	// Should end with the base name
	expectedEnd := "cleanup_duplicates.sh"
	if len(path) < len(expectedEnd) || path[len(path)-len(expectedEnd):] != expectedEnd {
		t.Errorf("GetOutputPath = %s, should end with %s", path, expectedEnd)
	}
}

// TestValidator_Validate_UnreadableFile tests Validate with unreadable file
func TestValidator_Validate_UnreadableFile(t *testing.T) {
	// This test may not work in all environments due to permissions
	// Skip if running as root
	if os.Geteuid() == 0 {
		t.Skip("Skipping test when running as root")
	}

	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.sh")

	validScript := `#!/bin/bash
rm -f "/tmp/test.txt"
`
	if err := os.WriteFile(scriptPath, []byte(validScript), 0644); err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	// Make file unreadable
	if err := os.Chmod(scriptPath, 0000); err != nil {
		t.Fatalf("Failed to change permissions: %v", err)
	}
	defer os.Chmod(scriptPath, 0644)

	v := script.NewValidator()
	err := v.Validate(scriptPath)
	if err == nil {
		t.Error("Validate should return error for unreadable file")
	}

	errors := v.GetErrors()
	if len(errors) == 0 {
		t.Error("Expected errors for unreadable file")
	}
}

// TestValidator_Validate_MultipleWarnings tests multiple warnings
func TestValidator_Validate_MultipleWarnings(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.sh")

	// Create script with multiple issues
	badScript := `echo "no shebang, no rm"
`
	if err := os.WriteFile(scriptPath, []byte(badScript), 0644); err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	v := script.NewValidator()
	err := v.Validate(scriptPath)
	if err != nil {
		t.Errorf("Validate should not return error: %v", err)
	}

	warnings := v.GetWarnings()
	if len(warnings) < 2 {
		t.Errorf("Expected at least 2 warnings, got %d", len(warnings))
	}
}

// TestValidator_Reset tests that Validate resets warnings/errors
func TestValidator_Reset(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.sh")

	// First validation with errors
	v := script.NewValidator()
	v.Validate("/nonexistent/file.sh")

	errors := v.GetErrors()
	if len(errors) == 0 {
		t.Fatal("Expected errors from first validation")
	}

	// Second validation with valid script
	validScript := `#!/bin/bash
rm -f "/tmp/test.txt"
`
	if err := os.WriteFile(scriptPath, []byte(validScript), 0755); err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	v.Validate(scriptPath)

	// Errors should be cleared
	errors = v.GetErrors()
	if len(errors) != 0 {
		t.Errorf("Errors should be cleared after new validation, got %d", len(errors))
	}
}
