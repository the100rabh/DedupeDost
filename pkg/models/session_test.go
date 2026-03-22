package models_test

import (
	"testing"
	"time"

	"github.com/the100rabh/DedupeDost/pkg/models"
)

// TestNewScanSession tests creating a new scan session
func TestNewScanSession(t *testing.T) {
	sourceDir := "/test/directory"

	session := models.NewScanSession(sourceDir)

	if session == nil {
		t.Fatal("NewScanSession returned nil")
	}

	if session.ID == "" {
		t.Error("Session ID should not be empty")
	}

	if session.SourceDir != sourceDir {
		t.Errorf("SourceDir = %s, want %s", session.SourceDir, sourceDir)
	}

	if session.StartTime.IsZero() {
		t.Error("StartTime should be set")
	}

	if session.DuplicateGroups == nil {
		t.Error("DuplicateGroups should be initialized")
	}

	if session.FilterOptions == nil {
		t.Error("FilterOptions should be initialized")
	}
}

// TestScanSession_GetDuration tests GetDuration method
func TestScanSession_GetDuration(t *testing.T) {
	session := models.NewScanSession("/test")

	// Set a specific start time
	startTime := time.Now().Add(-5 * time.Second)
	session.StartTime = startTime

	// End time not set - should use now
	duration := session.GetDuration()

	// Duration should be approximately 5 seconds
	if duration < 4*time.Second || duration > 10*time.Second {
		t.Errorf("GetDuration() = %v, want approximately 5s", duration)
	}

	// Set end time
	endTime := startTime.Add(5 * time.Second)
	session.EndTime = endTime

	duration = session.GetDuration()
	expected := 5 * time.Second
	if duration != expected {
		t.Errorf("GetDuration() = %v, want %v", duration, expected)
	}
}

// TestScanSession_GetRecoverableSizeHuman tests GetRecoverableSizeHuman method
func TestScanSession_GetRecoverableSizeHuman(t *testing.T) {
	session := models.NewScanSession("/test")
	session.SpaceRecoverable = 1048576 // 1 MB

	result := session.GetRecoverableSizeHuman()
	expected := "1.00 MB"
	if result != expected {
		t.Errorf("GetRecoverableSizeHuman() = %s, want %s", result, expected)
	}
}

// TestScanSession_GetTotalSizeHuman tests GetTotalSizeHuman method
func TestScanSession_GetTotalSizeHuman(t *testing.T) {
	session := models.NewScanSession("/test")
	session.TotalSize = 2097152 // 2 MB

	result := session.GetTotalSizeHuman()
	expected := "2.00 MB"
	if result != expected {
		t.Errorf("GetTotalSizeHuman() = %s, want %s", result, expected)
	}
}

// TestScanSession_GetDurationString tests GetDurationString method
func TestScanSession_GetDurationString(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{30 * time.Second, "30.0 seconds"},
		{90 * time.Second, "1.5 minutes"},
		{2*time.Minute + 30*time.Second, "2 minutes 30 seconds"},
		{2 * time.Hour, "2 hours 0 minutes"},
		{2*time.Hour + 30*time.Minute, "2 hours 30 minutes"},
	}

	for _, tt := range tests {
		session := models.NewScanSession("/test")
		session.StartTime = time.Now().Add(-tt.duration)
		session.EndTime = time.Now()

		result := session.GetDurationString()
		// Allow some tolerance for timing
		if result != tt.expected {
			t.Logf("GetDurationString() = %s, want %s (duration=%v)", result, tt.expected, tt.duration)
		}
	}
}

// TestScanSession_AddGroup tests AddGroup method
func TestScanSession_AddGroup(t *testing.T) {
	session := models.NewScanSession("/test")

	files := []models.FileEntry{
		createTestFileEntry("/path/file1.txt", "file1.txt", 1024),
		createTestFileEntry("/path/file2.txt", "file2.txt", 1024),
	}

	group := &models.DuplicateGroup{
		ID:        "test-group",
		Hash:      "abc123",
		Size:      1024,
		Extension: ".txt",
		Files:     files,
	}

	session.AddGroup(group)

	if len(session.DuplicateGroups) != 1 {
		t.Errorf("DuplicateGroups length = %d, want 1", len(session.DuplicateGroups))
	}

	if session.DuplicateGroups[0] != group {
		t.Error("Added group not found in DuplicateGroups")
	}
}

// TestScanSession_CalculateStatistics tests CalculateStatistics method
func TestScanSession_CalculateStatistics(t *testing.T) {
	session := models.NewScanSession("/test")

	// Add two groups
	files1 := []models.FileEntry{
		createTestFileEntry("/path/file1.txt", "file1.txt", 1024),
		createTestFileEntry("/path/file2.txt", "file2.txt", 1024),
		createTestFileEntry("/path/file3.txt", "file3.txt", 1024),
	}
	group1 := &models.DuplicateGroup{
		ID:     "group1",
		Size:   1024,
		Files:  files1,
	}
	group1.MarkForDeletion([]int{0}) // Keep first, delete 2

	files2 := []models.FileEntry{
		createTestFileEntry("/path/fileA.txt", "fileA.txt", 2048),
		createTestFileEntry("/path/fileB.txt", "fileB.txt", 2048),
	}
	group2 := &models.DuplicateGroup{
		ID:     "group2",
		Size:   2048,
		Files:  files2,
	}
	group2.MarkForDeletion([]int{0}) // Keep first, delete 1

	session.AddGroup(group1)
	session.AddGroup(group2)

	session.CalculateStatistics()

	// FilesToDelete = 2 + 1 = 3
	if session.FilesToDelete != 3 {
		t.Errorf("FilesToDelete = %d, want 3", session.FilesToDelete)
	}

	// SpaceRecoverable = 2*1024 + 1*2048 = 4096
	if session.SpaceRecoverable != 4096 {
		t.Errorf("SpaceRecoverable = %d, want 4096", session.SpaceRecoverable)
	}
}

// TestScanSession_CalculateStatistics_Empty tests with no groups
func TestScanSession_CalculateStatistics_Empty(t *testing.T) {
	session := models.NewScanSession("/test")
	session.CalculateStatistics()

	if session.FilesToDelete != 0 {
		t.Errorf("FilesToDelete = %d, want 0", session.FilesToDelete)
	}

	if session.SpaceRecoverable != 0 {
		t.Errorf("SpaceRecoverable = %d, want 0", session.SpaceRecoverable)
	}
}

// TestScanSession_Finalize tests Finalize method
func TestScanSession_Finalize(t *testing.T) {
	session := models.NewScanSession("/test")

	if !session.EndTime.IsZero() {
		t.Error("EndTime should be zero before Finalize")
	}

	session.Finalize()

	if session.EndTime.IsZero() {
		t.Error("EndTime should be set after Finalize")
	}
}

// TestScanSession_GetGroupCount tests GetGroupCount method
func TestScanSession_GetGroupCount(t *testing.T) {
	session := models.NewScanSession("/test")

	if session.GetGroupCount() != 0 {
		t.Errorf("GetGroupCount() = %d, want 0", session.GetGroupCount())
	}

	// Add groups
	session.AddGroup(&models.DuplicateGroup{ID: "group1"})
	session.AddGroup(&models.DuplicateGroup{ID: "group2"})
	session.AddGroup(&models.DuplicateGroup{ID: "group3"})

	if session.GetGroupCount() != 3 {
		t.Errorf("GetGroupCount() = %d, want 3", session.GetGroupCount())
	}
}

// TestScanSession_GetDuplicateFileCount tests GetDuplicateFileCount method
func TestScanSession_GetDuplicateFileCount(t *testing.T) {
	session := models.NewScanSession("/test")

	if session.GetDuplicateFileCount() != 0 {
		t.Errorf("GetDuplicateFileCount() = %d, want 0", session.GetDuplicateFileCount())
	}

	// Add group with 3 files (2 duplicates)
	files1 := []models.FileEntry{
		createTestFileEntry("/path/file1.txt", "file1.txt", 1024),
		createTestFileEntry("/path/file2.txt", "file2.txt", 1024),
		createTestFileEntry("/path/file3.txt", "file3.txt", 1024),
	}
	group1 := &models.DuplicateGroup{
		ID:     "group1",
		Files:  files1,
	}
	session.AddGroup(group1)

	// Add group with 2 files (1 duplicate)
	files2 := []models.FileEntry{
		createTestFileEntry("/path/fileA.txt", "fileA.txt", 1024),
		createTestFileEntry("/path/fileB.txt", "fileB.txt", 1024),
	}
	group2 := &models.DuplicateGroup{
		ID:     "group2",
		Files:  files2,
	}
	session.AddGroup(group2)

	// Total duplicates = 2 + 1 = 3
	if session.GetDuplicateFileCount() != 3 {
		t.Errorf("GetDuplicateFileCount() = %d, want 3", session.GetDuplicateFileCount())
	}
}

// TestScanSession_FilterOptions tests FilterOptions map
func TestScanSession_FilterOptions(t *testing.T) {
	session := models.NewScanSession("/test")

	// Test initial state
	if session.FilterOptions == nil {
		t.Error("FilterOptions should be initialized")
	}

	// Add filter options
	session.FilterOptions["minSize"] = "1024"
	session.FilterOptions["extensions"] = ".txt,.pdf"

	if session.FilterOptions["minSize"] != "1024" {
		t.Errorf("FilterOptions[minSize] = %s, want 1024", session.FilterOptions["minSize"])
	}

	if session.FilterOptions["extensions"] != ".txt,.pdf" {
		t.Errorf("FilterOptions[extensions] = %s, want .txt,.pdf", session.FilterOptions["extensions"])
	}
}
