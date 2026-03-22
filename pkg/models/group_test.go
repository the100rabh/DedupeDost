package models_test

import (
	"testing"
	"time"

	"github.com/the100rabh/DedupeDost/pkg/models"
	"github.com/the100rabh/DedupeDost/pkg/types"
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
		FileType:    types.FileTypeText,
		IsHidden:    false,
		IsSymlink:   false,
	}
}

// TestDuplicateGroup_GetRecoverableSize tests GetRecoverableSize method
func TestDuplicateGroup_GetRecoverableSize(t *testing.T) {
	files := []models.FileEntry{
		createTestFileEntry("/path/file1.txt", "file1.txt", 1024),
		createTestFileEntry("/path/file2.txt", "file2.txt", 1024),
		createTestFileEntry("/path/file3.txt", "file3.txt", 1024),
	}

	group := models.DuplicateGroup{
		ID:        "test-group",
		Hash:      "abc123",
		Size:      1024,
		Extension: ".txt",
		FileType:  types.FileTypeText,
		Files:     files,
	}

	// Recoverable = (3-1) * 1024 = 2048
	expected := int64(2048)
	result := group.GetRecoverableSize()
	if result != expected {
		t.Errorf("GetRecoverableSize() = %d, want %d", result, expected)
	}
}

// TestDuplicateGroup_GetRecoverableSize_SingleFile tests with single file
func TestDuplicateGroup_GetRecoverableSize_SingleFile(t *testing.T) {
	files := []models.FileEntry{
		createTestFileEntry("/path/file1.txt", "file1.txt", 1024),
	}

	group := models.DuplicateGroup{
		ID:        "test-group",
		Hash:      "abc123",
		Size:      1024,
		Extension: ".txt",
		FileType:  types.FileTypeText,
		Files:     files,
	}

	result := group.GetRecoverableSize()
	if result != 0 {
		t.Errorf("GetRecoverableSize() = %d, want 0 for single file", result)
	}
}

// TestDuplicateGroup_GetFileCount tests GetFileCount method
func TestDuplicateGroup_GetFileCount(t *testing.T) {
	files := []models.FileEntry{
		createTestFileEntry("/path/file1.txt", "file1.txt", 1024),
		createTestFileEntry("/path/file2.txt", "file2.txt", 1024),
		createTestFileEntry("/path/file3.txt", "file3.txt", 1024),
	}

	group := models.DuplicateGroup{
		Files: files,
	}

	result := group.GetFileCount()
	if result != 3 {
		t.Errorf("GetFileCount() = %d, want 3", result)
	}
}

// TestDuplicateGroup_GetDuplicateCount tests GetDuplicateCount method
func TestDuplicateGroup_GetDuplicateCount(t *testing.T) {
	files := []models.FileEntry{
		createTestFileEntry("/path/file1.txt", "file1.txt", 1024),
		createTestFileEntry("/path/file2.txt", "file2.txt", 1024),
		createTestFileEntry("/path/file3.txt", "file3.txt", 1024),
	}

	group := models.DuplicateGroup{
		Files: files,
	}

	// Duplicates = 3 - 1 = 2
	result := group.GetDuplicateCount()
	if result != 2 {
		t.Errorf("GetDuplicateCount() = %d, want 2", result)
	}
}

// TestDuplicateGroup_GetDuplicateCount_SingleFile tests with single file
func TestDuplicateGroup_GetDuplicateCount_SingleFile(t *testing.T) {
	files := []models.FileEntry{
		createTestFileEntry("/path/file1.txt", "file1.txt", 1024),
	}

	group := models.DuplicateGroup{
		Files: files,
	}

	result := group.GetDuplicateCount()
	if result != 0 {
		t.Errorf("GetDuplicateCount() = %d, want 0 for single file", result)
	}
}

// TestDuplicateGroup_GetTotalSize tests GetTotalSize method
func TestDuplicateGroup_GetTotalSize(t *testing.T) {
	files := []models.FileEntry{
		createTestFileEntry("/path/file1.txt", "file1.txt", 1024),
		createTestFileEntry("/path/file2.txt", "file2.txt", 1024),
		createTestFileEntry("/path/file3.txt", "file3.txt", 1024),
	}

	group := models.DuplicateGroup{
		Size:  1024,
		Files: files,
	}

	// Total = 3 * 1024 = 3072
	expected := int64(3072)
	result := group.GetTotalSize()
	if result != expected {
		t.Errorf("GetTotalSize() = %d, want %d", result, expected)
	}
}

// TestDuplicateGroup_GetDisplaySize tests GetDisplaySize method
func TestDuplicateGroup_GetDisplaySize(t *testing.T) {
	group := models.DuplicateGroup{
		Size: 1048576, // 1 MB
	}

	result := group.GetDisplaySize()
	expected := "1.00 MB"
	if result != expected {
		t.Errorf("GetDisplaySize() = %s, want %s", result, expected)
	}
}

// TestDuplicateGroup_GetDisplayRecoverableSize tests GetDisplayRecoverableSize method
func TestDuplicateGroup_GetDisplayRecoverableSize(t *testing.T) {
	files := []models.FileEntry{
		createTestFileEntry("/path/file1.txt", "file1.txt", 1048576),
		createTestFileEntry("/path/file2.txt", "file2.txt", 1048576),
	}

	group := models.DuplicateGroup{
		Size:  1048576, // 1 MB
		Files: files,
	}

	result := group.GetDisplayRecoverableSize()
	expected := "1.00 MB" // Recoverable = 1 * 1MB
	if result != expected {
		t.Errorf("GetDisplayRecoverableSize() = %s, want %s", result, expected)
	}
}

// TestDuplicateGroup_MarkForDeletion tests MarkForDeletion method
func TestDuplicateGroup_MarkForDeletion(t *testing.T) {
	files := []models.FileEntry{
		createTestFileEntry("/path/file1.txt", "file1.txt", 1024),
		createTestFileEntry("/path/file2.txt", "file2.txt", 1024),
		createTestFileEntry("/path/file3.txt", "file3.txt", 1024),
	}

	group := models.DuplicateGroup{
		Files: files,
	}

	// Keep only the first file (index 0)
	group.MarkForDeletion([]int{0})

	if len(group.KeepIndices) != 1 {
		t.Errorf("KeepIndices length = %d, want 1", len(group.KeepIndices))
	}

	if group.KeepIndices[0] != 0 {
		t.Errorf("KeepIndices[0] = %d, want 0", group.KeepIndices[0])
	}

	// Should have 2 files marked for deletion
	if len(group.DeletePaths) != 2 {
		t.Errorf("DeletePaths length = %d, want 2", len(group.DeletePaths))
	}

	// Verify correct files are marked for deletion
	deleteMap := make(map[string]bool)
	for _, path := range group.DeletePaths {
		deleteMap[path] = true
	}

	if !deleteMap["/path/file2.txt"] {
		t.Error("file2.txt should be marked for deletion")
	}
	if !deleteMap["/path/file3.txt"] {
		t.Error("file3.txt should be marked for deletion")
	}
	if deleteMap["/path/file1.txt"] {
		t.Error("file1.txt should NOT be marked for deletion")
	}
}

// TestDuplicateGroup_MarkForDeletion_MultipleKeep tests keeping multiple files
func TestDuplicateGroup_MarkForDeletion_MultipleKeep(t *testing.T) {
	files := []models.FileEntry{
		createTestFileEntry("/path/file1.txt", "file1.txt", 1024),
		createTestFileEntry("/path/file2.txt", "file2.txt", 1024),
		createTestFileEntry("/path/file3.txt", "file3.txt", 1024),
		createTestFileEntry("/path/file4.txt", "file4.txt", 1024),
	}

	group := models.DuplicateGroup{
		Files: files,
	}

	// Keep files at index 0 and 2
	group.MarkForDeletion([]int{0, 2})

	if len(group.KeepIndices) != 2 {
		t.Errorf("KeepIndices length = %d, want 2", len(group.KeepIndices))
	}

	// Should have 2 files marked for deletion (indices 1 and 3)
	if len(group.DeletePaths) != 2 {
		t.Errorf("DeletePaths length = %d, want 2", len(group.DeletePaths))
	}

	deleteMap := make(map[string]bool)
	for _, path := range group.DeletePaths {
		deleteMap[path] = true
	}

	if !deleteMap["/path/file2.txt"] {
		t.Error("file2.txt should be marked for deletion")
	}
	if !deleteMap["/path/file4.txt"] {
		t.Error("file4.txt should be marked for deletion")
	}
}

// TestDuplicateGroup_ClearDecisions tests ClearDecisions method
func TestDuplicateGroup_ClearDecisions(t *testing.T) {
	files := []models.FileEntry{
		createTestFileEntry("/path/file1.txt", "file1.txt", 1024),
		createTestFileEntry("/path/file2.txt", "file2.txt", 1024),
	}

	group := models.DuplicateGroup{
		Files:       files,
		KeepIndices: []int{0},
		DeletePaths: []string{"/path/file2.txt"},
		RuleApplied: &models.Rule{ID: "test-rule"},
	}

	group.ClearDecisions()

	if len(group.KeepIndices) != 0 {
		t.Errorf("KeepIndices not cleared, length = %d", len(group.KeepIndices))
	}

	if len(group.DeletePaths) != 0 {
		t.Errorf("DeletePaths not cleared, length = %d", len(group.DeletePaths))
	}

	if group.RuleApplied != nil {
		t.Error("RuleApplied not cleared")
	}
}

// TestDuplicateGroup_MarkForDeletion_EmptyKeepIndices tests with empty keep indices
func TestDuplicateGroup_MarkForDeletion_EmptyKeepIndices(t *testing.T) {
	files := []models.FileEntry{
		createTestFileEntry("/path/file1.txt", "file1.txt", 1024),
		createTestFileEntry("/path/file2.txt", "file2.txt", 1024),
	}

	group := models.DuplicateGroup{
		Files: files,
	}

	// Empty keep indices = delete all
	group.MarkForDeletion([]int{})

	if len(group.KeepIndices) != 0 {
		t.Errorf("KeepIndices length = %d, want 0", len(group.KeepIndices))
	}

	if len(group.DeletePaths) != 2 {
		t.Errorf("DeletePaths length = %d, want 2", len(group.DeletePaths))
	}
}
