package detector_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/the100rabh/DedupeDost/internal/detector"
	"github.com/the100rabh/DedupeDost/internal/scanner"
)

// Helper function to create test directory structure
func setupTestDir(t *testing.T) (string, func()) {
	t.Helper()

	dir, err := os.MkdirTemp("", "detector_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(dir)
	}

	return dir, cleanup
}

// Helper function to create test files with specific content
func createTestFile(t *testing.T, dir, name string, content []byte) string {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	return path
}

// TestDetector_NoDuplicates tests detection with no duplicates
func TestDetector_NoDuplicates(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	// Create files with unique content
	createTestFile(t, dir, "file1.txt", []byte("unique content 1"))
	createTestFile(t, dir, "file2.txt", []byte("unique content 2"))
	createTestFile(t, dir, "file3.txt", []byte("unique content 3"))

	d := detector.NewDetector(dir, scanner.ScanOptions{
		Recursive: true,
	})

	err := d.Detect()
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	groups := d.GetGroups()
	if len(groups) != 0 {
		t.Errorf("Expected 0 duplicate groups, got %d", len(groups))
	}

	stats := d.GetStatistics()
	if stats.DuplicateGroups != 0 {
		t.Errorf("Expected 0 duplicate groups in stats, got %d", stats.DuplicateGroups)
	}
	if stats.DuplicateFiles != 0 {
		t.Errorf("Expected 0 duplicate files, got %d", stats.DuplicateFiles)
	}
}

// TestDetector_SimpleDuplicates tests detection of simple duplicates
func TestDetector_SimpleDuplicates(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	// Create duplicate content
	content := []byte("This is duplicate content for testing")

	// Create 3 duplicate files
	file1 := createTestFile(t, dir, "original.txt", content)
	file2 := createTestFile(t, dir, "copy1.txt", content)
	file3 := createTestFile(t, dir, "copy2.txt", content)

	// Create a unique file
	createTestFile(t, dir, "unique.txt", []byte("unique content"))

	d := detector.NewDetector(dir, scanner.ScanOptions{
		Recursive: true,
	})

	err := d.Detect()
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	groups := d.GetGroups()
	if len(groups) != 1 {
		t.Fatalf("Expected 1 duplicate group, got %d", len(groups))
	}

	group := groups[0]
	if len(group.Files) != 3 {
		t.Errorf("Expected 3 files in duplicate group, got %d", len(group.Files))
		t.Logf("Files in group: %v", group.Files)
	}

	// Verify all duplicate files are in the group
	filePaths := make(map[string]bool)
	for _, f := range group.Files {
		filePaths[f.Path] = true
	}

	if !filePaths[file1] {
		t.Errorf("Missing file1 in duplicate group: %s", file1)
	}
	if !filePaths[file2] {
		t.Errorf("Missing file2 in duplicate group: %s", file2)
	}
	if !filePaths[file3] {
		t.Errorf("Missing file3 in duplicate group: %s", file3)
	}

	stats := d.GetStatistics()
	if stats.DuplicateGroups != 1 {
		t.Errorf("Expected 1 duplicate group in stats, got %d", stats.DuplicateGroups)
	}
	if stats.DuplicateFiles != 2 {
		t.Errorf("Expected 2 duplicate files (excluding original), got %d", stats.DuplicateFiles)
	}
}

// TestDetector_MultipleGroups tests detection of multiple duplicate groups
func TestDetector_MultipleGroups(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	// Group 1: Text files
	content1 := []byte("Content for group 1")
	createTestFile(t, dir, "text1.txt", content1)
	createTestFile(t, dir, "text2.txt", content1)
	createTestFile(t, dir, "text3.txt", content1)

	// Group 2: Config files
	content2 := []byte(`{"config": "value"}`)
	createTestFile(t, dir, "config1.json", content2)
	createTestFile(t, dir, "config2.json", content2)

	// Group 3: Log files
	content3 := []byte("Log entry 123")
	createTestFile(t, dir, "app1.log", content3)
	createTestFile(t, dir, "app2.log", content3)
	createTestFile(t, dir, "app3.log", content3)
	createTestFile(t, dir, "app4.log", content3)

	d := detector.NewDetector(dir, scanner.ScanOptions{
		Recursive: true,
	})

	err := d.Detect()
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	groups := d.GetGroups()
	if len(groups) != 3 {
		t.Fatalf("Expected 3 duplicate groups, got %d", len(groups))
	}

	// Debug: print all groups
	for i, g := range groups {
		t.Logf("Group %d: %d files, size=%d, ext=%s", i, len(g.Files), g.Size, g.Extension)
		for _, f := range g.Files {
			t.Logf("  - %s", f.Name)
		}
	}

	// Verify each group has correct file count
	groupSizes := make(map[int]int)
	for _, group := range groups {
		groupSizes[len(group.Files)]++
	}

	// We have: 1 group with 3 files (.txt), 1 group with 2 files (.json), 1 group with 4 files (.log)
	if groupSizes[3] != 1 {
		t.Errorf("Expected 1 group with 3 files, got %d", groupSizes[3])
	}
	if groupSizes[2] != 1 {
		t.Errorf("Expected 1 group with 2 files, got %d", groupSizes[2])
	}
	if groupSizes[4] != 1 {
		t.Errorf("Expected 1 group with 4 files, got %d", groupSizes[4])
	}

	stats := d.GetStatistics()
	if stats.DuplicateGroups != 3 {
		t.Errorf("Expected 3 duplicate groups in stats, got %d", stats.DuplicateGroups)
	}
	// Duplicate files = (3-1) + (2-1) + (4-1) = 2 + 1 + 3 = 6
	if stats.DuplicateFiles != 6 {
		t.Errorf("Expected 6 duplicate files, got %d", stats.DuplicateFiles)
	}
}

// TestDetector_NestedDirectories tests detection in nested directories
func TestDetector_NestedDirectories(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	// Create nested structure
	subdir1 := filepath.Join(dir, "sub1")
	subdir2 := filepath.Join(dir, "sub1", "sub2")
	subdir3 := filepath.Join(dir, "sub3")
	os.MkdirAll(subdir2, 0755)
	os.MkdirAll(subdir3, 0755)

	// Same content in different directories
	content := []byte("Duplicate across directories")
	createTestFile(t, dir, "root.txt", content)
	createTestFile(t, subdir1, "level1.txt", content)
	createTestFile(t, subdir2, "level2.txt", content)
	createTestFile(t, subdir3, "other.txt", content)

	d := detector.NewDetector(dir, scanner.ScanOptions{
		Recursive: true,
	})

	err := d.Detect()
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	groups := d.GetGroups()
	if len(groups) != 1 {
		t.Fatalf("Expected 1 duplicate group, got %d", len(groups))
	}

	group := groups[0]
	if len(group.Files) != 4 {
		t.Errorf("Expected 4 files in duplicate group, got %d", len(group.Files))
		for i, f := range group.Files {
			t.Logf("File %d: %s", i, f.Path)
		}
	}
}

// TestDetector_EmptyDirectory tests detection in empty directory
func TestDetector_EmptyDirectory(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	d := detector.NewDetector(dir, scanner.ScanOptions{
		Recursive: true,
	})

	err := d.Detect()
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	groups := d.GetGroups()
	if len(groups) != 0 {
		t.Errorf("Expected 0 duplicate groups, got %d", len(groups))
	}
}

// TestDetector_SingleFile tests detection with single file (no duplicates possible)
func TestDetector_SingleFile(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	createTestFile(t, dir, "only.txt", []byte("single file"))

	d := detector.NewDetector(dir, scanner.ScanOptions{
		Recursive: true,
	})

	err := d.Detect()
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	groups := d.GetGroups()
	if len(groups) != 0 {
		t.Errorf("Expected 0 duplicate groups, got %d", len(groups))
	}
}

// TestDetector_SameSizeDifferentContent tests that same-size files with different content are not marked as duplicates
func TestDetector_SameSizeDifferentContent(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	// Create files with same size but different content
	createTestFile(t, dir, "file1.txt", []byte("AAAA"))
	createTestFile(t, dir, "file2.txt", []byte("BBBB"))
	createTestFile(t, dir, "file3.txt", []byte("CCCC"))

	d := detector.NewDetector(dir, scanner.ScanOptions{
		Recursive: true,
	})

	err := d.Detect()
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	groups := d.GetGroups()
	if len(groups) != 0 {
		t.Errorf("Expected 0 duplicate groups (same size, different content), got %d", len(groups))
		for _, g := range groups {
			t.Logf("Unexpected group with %d files", len(g.Files))
		}
	}
}

// TestDetector_ExactDuplicatesAllFilesIncluded tests that ALL files in a duplicate set are included
func TestDetector_ExactDuplicatesAllFilesIncluded(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	// Create 10 exact duplicates
	content := []byte("Exact duplicate content for comprehensive testing")
	expectedFiles := make([]string, 10)

	for i := 0; i < 10; i++ {
		expectedFiles[i] = createTestFile(t, dir, filepath.Join("duplicate_"+string(rune('a'+i))+".txt"), content)
	}

	d := detector.NewDetector(dir, scanner.ScanOptions{
		Recursive: true,
	})

	err := d.Detect()
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	groups := d.GetGroups()
	if len(groups) != 1 {
		t.Fatalf("Expected 1 duplicate group, got %d", len(groups))
	}

	group := groups[0]
	if len(group.Files) != 10 {
		t.Errorf("Expected ALL 10 duplicate files in group, got %d", len(group.Files))
		t.Log("Files found in group:")
		for i, f := range group.Files {
			t.Logf("  %d: %s", i, f.Path)
		}
		t.Log("Expected files:")
		for i, p := range expectedFiles {
			t.Logf("  %d: %s", i, p)
		}
	}

	// Verify each expected file is in the group
	foundFiles := make(map[string]bool)
	for _, f := range group.Files {
		foundFiles[f.Path] = true
	}

	for _, expected := range expectedFiles {
		if !foundFiles[expected] {
			t.Errorf("Missing expected file: %s", expected)
		}
	}

	stats := d.GetStatistics()
	if stats.DuplicateFiles != 9 {
		t.Errorf("Expected 9 duplicate files (10-1), got %d", stats.DuplicateFiles)
	}
}

// TestDetector_MixedSizesAndDuplicates tests complex scenario with mixed sizes
func TestDetector_MixedSizesAndDuplicates(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	// Small duplicates (same size)
	small := []byte("small")
	createTestFile(t, dir, "small1.txt", small)
	createTestFile(t, dir, "small2.txt", small)

	// Medium duplicates (same size)
	medium := []byte("medium content here")
	createTestFile(t, dir, "medium1.txt", medium)
	createTestFile(t, dir, "medium2.txt", medium)
	createTestFile(t, dir, "medium3.txt", medium)

	// Large duplicates (same size)
	large := []byte("This is a larger file with more content for testing purposes")
	createTestFile(t, dir, "large1.txt", large)
	createTestFile(t, dir, "large2.txt", large)

	// Unique files of various sizes
	createTestFile(t, dir, "unique_small.txt", []byte("tiny"))
	createTestFile(t, dir, "unique_medium.txt", []byte("somewhat larger unique"))
	createTestFile(t, dir, "unique_large.txt", []byte("This is a unique large file that has no duplicates in this test"))

	d := detector.NewDetector(dir, scanner.ScanOptions{
		Recursive: true,
	})

	err := d.Detect()
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	groups := d.GetGroups()
	if len(groups) != 3 {
		t.Fatalf("Expected 3 duplicate groups, got %d", len(groups))
	}

	// Count total files in all groups
	totalInGroups := 0
	for _, g := range groups {
		totalInGroups += len(g.Files)
	}

	// Should be 2 + 3 + 2 = 7 files in duplicate groups
	if totalInGroups != 7 {
		t.Errorf("Expected 7 total files in duplicate groups, got %d", totalInGroups)
	}

	stats := d.GetStatistics()
	// Duplicate files = (2-1) + (3-1) + (2-1) = 1 + 2 + 1 = 4
	if stats.DuplicateFiles != 4 {
		t.Errorf("Expected 4 duplicate files, got %d", stats.DuplicateFiles)
	}
}
