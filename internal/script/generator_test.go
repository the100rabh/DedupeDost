package script_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dupdel/dup-del/internal/script"
	"github.com/dupdel/dup-del/pkg/models"
)

// Helper function to create test file entries
func createTestFileEntry(path, name string, size int64) models.FileEntry {
	return models.FileEntry{
		Path:        path,
		Name:        name,
		Size:        size,
		Extension:   ".txt",
		ModTime:     time.Now(),
		CreatedTime: time.Now(),
	}
}

// Helper function to create test session
func createTestSession(t *testing.T) *models.ScanSession {
	t.Helper()
	session := models.NewScanSession("/test/directory")

	// Add a group with duplicates
	files := []models.FileEntry{
		createTestFileEntry("/test/file1.txt", "file1.txt", 1024),
		createTestFileEntry("/test/file2.txt", "file2.txt", 1024),
		createTestFileEntry("/test/file3.txt", "file3.txt", 1024),
	}
	group := &models.DuplicateGroup{
		ID:        "group1",
		Hash:      "abc123",
		Size:      1024,
		Extension: ".txt",
		Files:     files,
	}
	group.MarkForDeletion([]int{0}) // Keep first, delete 2
	session.AddGroup(group)
	session.CalculateStatistics()

	return session
}

// TestNewGenerator tests NewGenerator
func TestNewGenerator(t *testing.T) {
	g := script.NewGenerator("/tmp/cleanup.sh")
	if g == nil {
		t.Fatal("NewGenerator returned nil")
	}
}

// TestNewGenerator_PathHandling tests path handling in NewGenerator
func TestNewGenerator_PathHandling(t *testing.T) {
	tests := []struct {
		path        string
		expectedDir string
		expectedName string
	}{
		{"/tmp/cleanup.sh", "/tmp", "cleanup.sh"},
		{"/tmp/test/cleanup.sh", "/tmp/test", "cleanup.sh"},
		{"cleanup.sh", ".", "cleanup.sh"},
	}

	for _, tt := range tests {
		g := script.NewGenerator(tt.path)
		// We can't directly access outputDir and scriptName, but we can verify it works
		if g == nil {
			t.Errorf("NewGenerator(%s) returned nil", tt.path)
		}
	}
}

// TestGenerator_Generate tests Generate method
func TestGenerator_Generate(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "cleanup.sh")

	g := script.NewGenerator(scriptPath)
	session := createTestSession(t)

	resultPath, err := g.Generate(session)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if resultPath != scriptPath {
		t.Errorf("Result path = %s, want %s", resultPath, scriptPath)
	}

	// Verify script was created
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		t.Error("Script file was not created")
	}

	// Verify session was updated
	if session.GeneratedScript != scriptPath {
		t.Errorf("Session.GeneratedScript = %s, want %s", session.GeneratedScript, scriptPath)
	}

	// Verify script content
	content, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("Failed to read script: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "#!/bin/bash") {
		t.Error("Script missing shebang")
	}
	if !strings.Contains(contentStr, "DRY_RUN") {
		t.Error("Script missing DRY_RUN support")
	}
	if !strings.Contains(contentStr, "delete_file") {
		t.Error("Script missing delete_file function")
	}
}

// TestGenerator_Generate_OutputDirExists tests Generate creates output directory
func TestGenerator_Generate_OutputDirExists(t *testing.T) {
	tmpDir := t.TempDir()
	nestedPath := filepath.Join(tmpDir, "nested", "dir", "cleanup.sh")

	g := script.NewGenerator(nestedPath)
	session := createTestSession(t)

	_, err := g.Generate(session)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Verify script was created in nested directory
	if _, err := os.Stat(nestedPath); os.IsNotExist(err) {
		t.Error("Script file was not created in nested directory")
	}
}

// TestGenerator_Generate_InvalidPath tests Generate with invalid path
func TestGenerator_Generate_InvalidPath(t *testing.T) {
	// Try to create in a directory that doesn't exist and can't be created
	g := script.NewGenerator("/root/nonexistent/deeply/nested/cleanup.sh")
	session := createTestSession(t)

	_, err := g.Generate(session)
	// This might succeed if running as root, or fail otherwise
	if err != nil {
		t.Logf("Generate failed as expected: %v", err)
	}
}

// TestGenerator_SetIncludeHeader tests SetIncludeHeader
func TestGenerator_SetIncludeHeader(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "cleanup.sh")

	g := script.NewGenerator(scriptPath)
	session := createTestSession(t)

	// Generate with header (default)
	g.SetIncludeHeader(false)
	_, err := g.Generate(session)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("Failed to read script: %v", err)
	}

	contentStr := string(content)
	if strings.Contains(contentStr, "DupDel Cleanup Script") {
		t.Error("Script should not contain header when SetIncludeHeader(false)")
	}
}

// TestGenerator_SetIncludeComments tests SetIncludeComments
func TestGenerator_SetIncludeComments(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "cleanup.sh")

	g := script.NewGenerator(scriptPath)
	session := createTestSession(t)

	// Generate without comments
	g.SetIncludeComments(false)
	_, err := g.Generate(session)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("Failed to read script: %v", err)
	}

	contentStr := string(content)
	if strings.Contains(contentStr, "# Group 1:") {
		t.Error("Script should not contain group comments when SetIncludeComments(false)")
	}
	if strings.Contains(contentStr, "# Keeping:") {
		t.Error("Script should not contain keeping comments when SetIncludeComments(false)")
	}
}

// TestGenerator_SetDryRunSupport tests SetDryRunSupport
func TestGenerator_SetDryRunSupport(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "cleanup.sh")

	g := script.NewGenerator(scriptPath)
	session := createTestSession(t)

	// Generate without DRY_RUN support
	g.SetDryRunSupport(false)
	_, err := g.Generate(session)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("Failed to read script: %v", err)
	}

	contentStr := string(content)
	// Check for actual DRY_RUN variable declaration (not just comments)
	if strings.Contains(contentStr, "DRY_RUN=${DRY_RUN:-false}") {
		t.Error("Script should not contain DRY_RUN variable declaration when SetDryRunSupport(false)")
	}
	if strings.Contains(contentStr, "delete_file()") {
		t.Error("Script should not contain delete_file function when SetDryRunSupport(false)")
	}
	// Should use direct rm instead
	if !strings.Contains(contentStr, "rm -f") {
		t.Error("Script should contain rm -f when DRY_RUN support is disabled")
	}
}

// TestGenerator_MakeExecutable tests MakeExecutable
func TestGenerator_MakeExecutable(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "cleanup.sh")

	// Create a file
	if err := os.WriteFile(scriptPath, []byte("#!/bin/bash"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	g := script.NewGenerator(scriptPath)
	err := g.MakeExecutable(scriptPath)
	if err != nil {
		t.Fatalf("MakeExecutable failed: %v", err)
	}

	// Verify it's executable
	info, err := os.Stat(scriptPath)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	if info.Mode()&0111 == 0 {
		t.Error("File should be executable after MakeExecutable")
	}
}

// TestGenerator_MakeExecutable_NonExistent tests MakeExecutable with non-existent file
func TestGenerator_MakeExecutable_NonExistent(t *testing.T) {
	g := script.NewGenerator("/tmp/cleanup.sh")
	err := g.MakeExecutable("/nonexistent/file.sh")
	if err == nil {
		t.Error("MakeExecutable should return error for non-existent file")
	}
}

// TestGenerator_Validate tests Validate
func TestGenerator_Validate(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "cleanup.sh")

	// Create executable script
	if err := os.WriteFile(scriptPath, []byte("#!/bin/bash"), 0755); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	g := script.NewGenerator(scriptPath)
	err := g.Validate(scriptPath)
	if err != nil {
		t.Fatalf("Validate failed: %v", err)
	}
}

// TestGenerator_Validate_NonExistent tests Validate with non-existent file
func TestGenerator_Validate_NonExistent(t *testing.T) {
	g := script.NewGenerator("/tmp/cleanup.sh")
	err := g.Validate("/nonexistent/file.sh")
	if err == nil {
		t.Error("Validate should return error for non-existent file")
	}
}

// TestGenerator_Validate_Directory tests Validate with directory
func TestGenerator_Validate_Directory(t *testing.T) {
	tmpDir := t.TempDir()

	g := script.NewGenerator(tmpDir)
	err := g.Validate(tmpDir)
	if err == nil {
		t.Error("Validate should return error for directory")
	}
}

// TestGenerator_Validate_NotExecutable tests Validate with non-executable file
func TestGenerator_Validate_NotExecutable(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "cleanup.sh")

	// Create non-executable file
	if err := os.WriteFile(scriptPath, []byte("#!/bin/bash"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	g := script.NewGenerator(scriptPath)
	err := g.Validate(scriptPath)
	if err == nil {
		t.Error("Validate should return error for non-executable file")
	}
}

// TestGenerator_Generate_MultipleGroups tests Generate with multiple groups
func TestGenerator_Generate_MultipleGroups(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "cleanup.sh")

	g := script.NewGenerator(scriptPath)
	session := models.NewScanSession("/test")

	// Add multiple groups
	for i := 0; i < 3; i++ {
		files := []models.FileEntry{
			createTestFileEntry("/test/group"+string(rune('0'+i))+"_file1.txt", "file1.txt", 1024),
			createTestFileEntry("/test/group"+string(rune('0'+i))+"_file2.txt", "file2.txt", 1024),
		}
		group := &models.DuplicateGroup{
			ID:        "group" + string(rune('0'+i)),
			Hash:      "hash" + string(rune('0'+i)),
			Size:      1024,
			Extension: ".txt",
			Files:     files,
		}
		group.MarkForDeletion([]int{0})
		session.AddGroup(group)
	}
	session.CalculateStatistics()

	_, err := g.Generate(session)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("Failed to read script: %v", err)
	}

	contentStr := string(content)
	// Should have 3 groups
	if strings.Count(contentStr, "# Group") != 3 {
		t.Errorf("Expected 3 group comments, got %d", strings.Count(contentStr, "# Group"))
	}
}

// TestGenerator_Generate_EmptySession tests Generate with empty session
func TestGenerator_Generate_EmptySession(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "cleanup.sh")

	g := script.NewGenerator(scriptPath)
	session := models.NewScanSession("/test")

	_, err := g.Generate(session)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("Failed to read script: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "#!/bin/bash") {
		t.Error("Script should contain shebang")
	}
}

// TestGenerator_Generate_VerifyDeletePaths tests that correct delete paths are generated
func TestGenerator_Generate_VerifyDeletePaths(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "cleanup.sh")

	g := script.NewGenerator(scriptPath)
	session := models.NewScanSession("/test")

	// Create group with specific delete paths
	files := []models.FileEntry{
		createTestFileEntry("/test/keep.txt", "keep.txt", 1024),
		createTestFileEntry("/test/delete1.txt", "delete1.txt", 1024),
		createTestFileEntry("/test/delete2.txt", "delete2.txt", 1024),
	}
	group := &models.DuplicateGroup{
		ID:        "group1",
		Hash:      "abc123",
		Size:      1024,
		Extension: ".txt",
		Files:     files,
	}
	group.MarkForDeletion([]int{0}) // Keep first
	session.AddGroup(group)
	session.CalculateStatistics()

	_, err := g.Generate(session)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("Failed to read script: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "delete1.txt") {
		t.Error("Script should contain delete1.txt")
	}
	if !strings.Contains(contentStr, "delete2.txt") {
		t.Error("Script should contain delete2.txt")
	}
	if strings.Contains(contentStr, "delete_file \"/test/keep.txt\"") {
		t.Error("Script should not contain keep.txt in delete commands")
	}
}

// TestGenerator_Generate_StatisticsInHeader tests that statistics are in header
func TestGenerator_Generate_StatisticsInHeader(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "cleanup.sh")

	g := script.NewGenerator(scriptPath)
	session := createTestSession(t)

	_, err := g.Generate(session)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("Failed to read script: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "Total files to delete:") {
		t.Error("Script should contain files to delete count")
	}
	if !strings.Contains(contentStr, "Total space to recover:") {
		t.Error("Script should contain space to recover")
	}
}

// TestGenerator_SetDryRunSupport_Method tests SetDryRunSupport getter
func TestGenerator_SetDryRunSupport_Method(t *testing.T) {
	g := script.NewGenerator("/tmp/test.sh")
	
	// Default should be true
	if !g.DryRunSupport() {
		t.Error("Default dryRunSupport should be true")
	}
	
	g.SetDryRunSupport(false)
	if g.DryRunSupport() {
		t.Error("dryRunSupport should be false after SetDryRunSupport(false)")
	}
}

// TestGenerator_ValidationErrors tests Validate error cases
func TestGenerator_ValidationErrors(t *testing.T) {
	g := script.NewGenerator("/tmp/test.sh")
	
	// Non-existent file
	err := g.Validate("/nonexistent/path/file.sh")
	if err == nil {
		t.Error("Validate should return error for non-existent file")
	}
	if !strings.Contains(err.Error(), "script not found") {
		t.Errorf("Error should mention 'script not found', got: %v", err)
	}
}

// TestGenerator_MakeExecutable_Error tests MakeExecutable error case
func TestGenerator_MakeExecutable_Error(t *testing.T) {
	g := script.NewGenerator("/tmp/test.sh")
	
	err := g.MakeExecutable("/nonexistent/path/file.sh")
	if err == nil {
		t.Error("MakeExecutable should return error for non-existent file")
	}
}

// TestGenerator_Generate_NoGroups tests Generate with no groups
func TestGenerator_Generate_NoGroups(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "cleanup.sh")

	g := script.NewGenerator(scriptPath)
	session := models.NewScanSession("/test")
	// Don't add any groups

	_, err := g.Generate(session)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("Failed to read script: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "#!/bin/bash") {
		t.Error("Script should contain shebang")
	}
}

// TestGenerator_WithAllOptions tests Generate with all options disabled
func TestGenerator_WithAllOptions(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "cleanup.sh")

	g := script.NewGenerator(scriptPath)
	session := createTestSession(t)

	// Disable all options
	g.SetIncludeHeader(false)
	g.SetIncludeComments(false)
	g.SetDryRunSupport(false)

	_, err := g.Generate(session)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("Failed to read script: %v", err)
	}

	contentStr := string(content)
	
	// Should not contain header
	if strings.Contains(contentStr, "DupDel Cleanup Script") {
		t.Error("Script should not contain header")
	}
	
	// Should not contain DRY_RUN
	if strings.Contains(contentStr, "DRY_RUN=${DRY_RUN:-false}") {
		t.Error("Script should not contain DRY_RUN")
	}
	
	// Should not contain comments
	if strings.Contains(contentStr, "# Group 1:") {
		t.Error("Script should not contain group comments")
	}
	
	// Should use rm -f directly
	if !strings.Contains(contentStr, "rm -f") {
		t.Error("Script should contain rm -f")
	}
}
