package detector_test

import (
	"testing"
	"time"

	"github.com/dupdel/dup-del/internal/detector"
	"github.com/dupdel/dup-del/pkg/models"
	"github.com/dupdel/dup-del/pkg/types"
)

// Helper function to create test file entries
func createTestFileEntry(path, name string, size int64, hash string) models.FileEntry {
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
		Hash:        hash,
	}
}

// TestNewGrouper tests NewGrouper
func TestNewGrouper(t *testing.T) {
	g := detector.NewGrouper()
	if g == nil {
		t.Fatal("NewGrouper returned nil")
	}
}

// TestGrouper_GroupBySize tests GroupBySize
func TestGrouper_GroupBySize(t *testing.T) {
	g := detector.NewGrouper()

	files := []models.FileEntry{
		createTestFileEntry("/path/file1.txt", "file1.txt", 1024, "hash1"),
		createTestFileEntry("/path/file2.txt", "file2.txt", 1024, "hash2"),
		createTestFileEntry("/path/file3.txt", "file3.txt", 2048, "hash3"),
		createTestFileEntry("/path/file4.txt", "file4.txt", 2048, "hash4"),
		createTestFileEntry("/path/file5.txt", "file5.txt", 512, "hash5"),
	}

	sizeMap := g.GroupBySize(files)

	if len(sizeMap) != 3 {
		t.Errorf("Expected 3 size groups, got %d", len(sizeMap))
	}

	if len(sizeMap[1024]) != 2 {
		t.Errorf("Expected 2 files with size 1024, got %d", len(sizeMap[1024]))
	}
	if len(sizeMap[2048]) != 2 {
		t.Errorf("Expected 2 files with size 2048, got %d", len(sizeMap[2048]))
	}
	if len(sizeMap[512]) != 1 {
		t.Errorf("Expected 1 file with size 512, got %d", len(sizeMap[512]))
	}
}

// TestGrouper_GroupBySize_Empty tests GroupBySize with empty list
func TestGrouper_GroupBySize_Empty(t *testing.T) {
	g := detector.NewGrouper()

	sizeMap := g.GroupBySize([]models.FileEntry{})

	if len(sizeMap) != 0 {
		t.Errorf("Expected empty size map, got %d entries", len(sizeMap))
	}
}

// TestGrouper_GroupByHash tests GroupByHash
func TestGrouper_GroupByHash(t *testing.T) {
	g := detector.NewGrouper()

	files := []models.FileEntry{
		createTestFileEntry("/path/file1.txt", "file1.txt", 1024, "abc123"),
		createTestFileEntry("/path/file2.txt", "file2.txt", 1024, "abc123"),
		createTestFileEntry("/path/file3.txt", "file3.txt", 1024, "def456"),
		createTestFileEntry("/path/file4.txt", "file4.txt", 1024, "def456"),
		createTestFileEntry("/path/file5.txt", "file5.txt", 1024, "ghi789"),
	}

	hashMap := g.GroupByHash(files)

	if len(hashMap) != 3 {
		t.Errorf("Expected 3 hash groups, got %d", len(hashMap))
	}

	if len(hashMap["abc123"]) != 2 {
		t.Errorf("Expected 2 files with hash abc123, got %d", len(hashMap["abc123"]))
	}
	if len(hashMap["def456"]) != 2 {
		t.Errorf("Expected 2 files with hash def456, got %d", len(hashMap["def456"]))
	}
	if len(hashMap["ghi789"]) != 1 {
		t.Errorf("Expected 1 file with hash ghi789, got %d", len(hashMap["ghi789"]))
	}
}

// TestGrouper_GroupByHash_EmptyHash tests GroupByHash with empty hash
func TestGrouper_GroupByHash_EmptyHash(t *testing.T) {
	g := detector.NewGrouper()

	files := []models.FileEntry{
		createTestFileEntry("/path/file1.txt", "file1.txt", 1024, ""),
		createTestFileEntry("/path/file2.txt", "file2.txt", 1024, ""),
		createTestFileEntry("/path/file3.txt", "file3.txt", 1024, "abc123"),
	}

	hashMap := g.GroupByHash(files)

	// Files with empty hash should be excluded
	if len(hashMap) != 1 {
		t.Errorf("Expected 1 hash group, got %d", len(hashMap))
	}
	if _, ok := hashMap[""]; ok {
		t.Error("Empty hash should not be included")
	}
}

// TestGrouper_FilterDuplicates tests FilterDuplicates
func TestGrouper_FilterDuplicates(t *testing.T) {
	g := detector.NewGrouper()

	groupMap := map[int64][]models.FileEntry{
		1024: {
			createTestFileEntry("/path/file1.txt", "file1.txt", 1024, "hash1"),
			createTestFileEntry("/path/file2.txt", "file2.txt", 1024, "hash2"),
		},
		2048: {
			createTestFileEntry("/path/file3.txt", "file3.txt", 2048, "hash3"),
		},
		512: {
			createTestFileEntry("/path/file4.txt", "file4.txt", 512, "hash4"),
			createTestFileEntry("/path/file5.txt", "file5.txt", 512, "hash5"),
			createTestFileEntry("/path/file6.txt", "file6.txt", 512, "hash6"),
		},
	}

	filtered := g.FilterDuplicates(groupMap)

	if len(filtered) != 2 {
		t.Errorf("Expected 2 groups after filtering, got %d", len(filtered))
	}

	if _, ok := filtered[1024]; !ok {
		t.Error("Size 1024 group should be included")
	}
	if _, ok := filtered[2048]; ok {
		t.Error("Size 2048 group should be excluded (only 1 file)")
	}
	if _, ok := filtered[512]; !ok {
		t.Error("Size 512 group should be included")
	}
}

// TestGrouper_FilterDuplicates_Empty tests FilterDuplicates with empty map
func TestGrouper_FilterDuplicates_Empty(t *testing.T) {
	g := detector.NewGrouper()

	filtered := g.FilterDuplicates(map[int64][]models.FileEntry{})

	if len(filtered) != 0 {
		t.Errorf("Expected empty filtered map, got %d entries", len(filtered))
	}
}

// TestGrouper_CreateDuplicateGroups tests CreateDuplicateGroups
func TestGrouper_CreateDuplicateGroups(t *testing.T) {
	g := detector.NewGrouper()

	hashMap := map[string][]models.FileEntry{
		"abc123": {
			createTestFileEntry("/path/file1.txt", "file1.txt", 1024, "abc123"),
			createTestFileEntry("/path/file2.txt", "file2.txt", 1024, "abc123"),
		},
		"def456": {
			createTestFileEntry("/path/file3.txt", "file3.txt", 2048, "def456"),
			createTestFileEntry("/path/file4.txt", "file4.txt", 2048, "def456"),
			createTestFileEntry("/path/file5.txt", "file5.txt", 2048, "def456"),
		},
		"single": {
			createTestFileEntry("/path/file6.txt", "file6.txt", 512, "single"),
		},
	}

	groups := g.CreateDuplicateGroups(hashMap)

	if len(groups) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(groups))
	}

	// Verify group properties
	for _, group := range groups {
		if group.ID == "" {
			t.Error("Group ID should not be empty")
		}
		if group.Hash == "" {
			t.Error("Group Hash should not be empty")
		}
		if len(group.Files) < 2 {
			t.Errorf("Group should have at least 2 files, got %d", len(group.Files))
		}
	}
}

// TestGrouper_CreateDuplicateGroups_Empty tests CreateDuplicateGroups with empty map
func TestGrouper_CreateDuplicateGroups_Empty(t *testing.T) {
	g := detector.NewGrouper()

	groups := g.CreateDuplicateGroups(map[string][]models.FileEntry{})

	if len(groups) != 0 {
		t.Errorf("Expected 0 groups, got %d", len(groups))
	}
}

// TestGrouper_generateGroupID tests unique ID generation
func TestGrouper_generateGroupID(t *testing.T) {
	g := detector.NewGrouper()

	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		// Use reflection or create groups to trigger ID generation
		hashMap := map[string][]models.FileEntry{
			"hash": {
				createTestFileEntry("/path/file1.txt", "file1.txt", 1024, "hash"),
				createTestFileEntry("/path/file2.txt", "file2.txt", 1024, "hash"),
			},
		}
		groups := g.CreateDuplicateGroups(hashMap)
		if len(groups) > 0 {
			if ids[groups[0].ID] {
				t.Errorf("Duplicate group ID generated: %s", groups[0].ID)
			}
			ids[groups[0].ID] = true
		}
	}
}

// TestSortGroups tests SortGroups
func TestSortGroups(t *testing.T) {
	files1 := []models.FileEntry{
		createTestFileEntry("/path/file1.txt", "file1.txt", 1024, "hash1"),
		createTestFileEntry("/path/file2.txt", "file2.txt", 1024, "hash1"),
	}
	files2 := []models.FileEntry{
		createTestFileEntry("/path/file3.txt", "file3.txt", 2048, "hash2"),
		createTestFileEntry("/path/file4.txt", "file4.txt", 2048, "hash2"),
		createTestFileEntry("/path/file5.txt", "file5.txt", 2048, "hash2"),
	}
	files3 := []models.FileEntry{
		createTestFileEntry("/path/file6.txt", "file6.txt", 512, "hash3"),
		createTestFileEntry("/path/file7.txt", "file7.txt", 512, "hash3"),
	}

	groups := []models.DuplicateGroup{
		{ID: "g1", Hash: "hash1", Size: 1024, Extension: ".txt", Files: files1},
		{ID: "g2", Hash: "hash2", Size: 2048, Extension: ".md", Files: files2},
		{ID: "g3", Hash: "hash3", Size: 512, Extension: ".go", Files: files3},
	}

	// Test SortBySize
	groupsBySize := make([]models.DuplicateGroup, len(groups))
	copy(groupsBySize, groups)
	detector.SortGroups(groupsBySize, detector.SortBySize)
	if groupsBySize[0].Size != 2048 {
		t.Errorf("SortBySize: first group size = %d, want 2048", groupsBySize[0].Size)
	}

	// Test SortByCount
	groupsByCount := make([]models.DuplicateGroup, len(groups))
	copy(groupsByCount, groups)
	detector.SortGroups(groupsByCount, detector.SortByCount)
	if len(groupsByCount[0].Files) != 3 {
		t.Errorf("SortByCount: first group count = %d, want 3", len(groupsByCount[0].Files))
	}

	// Test SortByRecoverable
	groupsByRecoverable := make([]models.DuplicateGroup, len(groups))
	copy(groupsByRecoverable, groups)
	detector.SortGroups(groupsByRecoverable, detector.SortByRecoverable)
	// Recoverable = (count-1) * size
	// g1: 1 * 1024 = 1024
	// g2: 2 * 2048 = 4096
	// g3: 1 * 512 = 512
	if groupsByRecoverable[0].Size != 2048 {
		t.Errorf("SortByRecoverable: first group size = %d, want 2048", groupsByRecoverable[0].Size)
	}

	// Test SortByName
	groupsByName := make([]models.DuplicateGroup, len(groups))
	copy(groupsByName, groups)
	detector.SortGroups(groupsByName, detector.SortByName)
	if groupsByName[0].Extension != ".go" {
		t.Errorf("SortByName: first extension = %s, want .go", groupsByName[0].Extension)
	}
}

// TestCalculateHash tests CalculateHash
func TestCalculateHash(t *testing.T) {
	data := []byte("Hello, World!")
	hash := detector.CalculateHash(data)

	if hash == "" {
		t.Error("CalculateHash returned empty string")
	}

	// Same data should produce same hash
	hash2 := detector.CalculateHash(data)
	if hash != hash2 {
		t.Error("Same data should produce same hash")
	}

	// Different data should produce different hash
	hash3 := detector.CalculateHash([]byte("Different data"))
	if hash == hash3 {
		t.Error("Different data should produce different hash")
	}
}

// TestFilesHaveSameContent tests FilesHaveSameContent
func TestFilesHaveSameContent(t *testing.T) {
	f1 := createTestFileEntry("/path/file1.txt", "file1.txt", 1024, "abc123")
	f2 := createTestFileEntry("/path/file2.txt", "file2.txt", 1024, "abc123")
	f3 := createTestFileEntry("/path/file3.txt", "file3.txt", 2048, "abc123")
	f4 := createTestFileEntry("/path/file4.txt", "file4.txt", 1024, "def456")

	if !detector.FilesHaveSameContent(f1, f2) {
		t.Error("Files with same size and hash should have same content")
	}
	if detector.FilesHaveSameContent(f1, f3) {
		t.Error("Files with different size should not have same content")
	}
	if detector.FilesHaveSameContent(f1, f4) {
		t.Error("Files with different hash should not have same content")
	}
}

// TestFindPotentialDuplicates tests FindPotentialDuplicates
func TestFindPotentialDuplicates(t *testing.T) {
	files := []models.FileEntry{
		createTestFileEntry("/path/file1.txt", "file1.txt", 1024, ""),
		createTestFileEntry("/path/file2.txt", "file2.txt", 1024, ""),
		createTestFileEntry("/path/file3.txt", "file3.txt", 2048, ""),
		createTestFileEntry("/path/file4.txt", "file4.txt", 512, ""),
		createTestFileEntry("/path/file5.txt", "file5.txt", 512, ""),
		createTestFileEntry("/path/file6.txt", "file6.txt", 512, ""),
	}

	potential := detector.FindPotentialDuplicates(files)

	if len(potential) != 5 {
		t.Errorf("Expected 5 potential duplicates, got %d", len(potential))
	}
}

// TestFindPotentialDuplicates_NoDuplicates tests with no duplicates
func TestFindPotentialDuplicates_NoDuplicates(t *testing.T) {
	files := []models.FileEntry{
		createTestFileEntry("/path/file1.txt", "file1.txt", 1024, ""),
		createTestFileEntry("/path/file2.txt", "file2.txt", 2048, ""),
		createTestFileEntry("/path/file3.txt", "file3.txt", 512, ""),
	}

	potential := detector.FindPotentialDuplicates(files)

	if len(potential) != 0 {
		t.Errorf("Expected 0 potential duplicates, got %d", len(potential))
	}
}

// TestEstimateDuplicates tests EstimateDuplicates
func TestEstimateDuplicates(t *testing.T) {
	files := []models.FileEntry{
		createTestFileEntry("/path/file1.txt", "file1.txt", 1024, ""),
		createTestFileEntry("/path/file2.txt", "file2.txt", 1024, ""),
		createTestFileEntry("/path/file3.txt", "file3.txt", 1024, ""),
		createTestFileEntry("/path/file4.txt", "file4.txt", 2048, ""),
		createTestFileEntry("/path/file5.txt", "file5.txt", 512, ""),
		createTestFileEntry("/path/file6.txt", "file6.txt", 512, ""),
	}

	estimate := detector.EstimateDuplicates(files)

	// 1024: 3 files = 2 duplicates
	// 2048: 1 file = 0 duplicates
	// 512: 2 files = 1 duplicate
	// Total: 3 duplicates
	if estimate != 3 {
		t.Errorf("Expected estimate 3, got %d", estimate)
	}
}

// TestEstimateDuplicates_NoDuplicates tests with no duplicates
func TestEstimateDuplicates_NoDuplicates(t *testing.T) {
	files := []models.FileEntry{
		createTestFileEntry("/path/file1.txt", "file1.txt", 1024, ""),
		createTestFileEntry("/path/file2.txt", "file2.txt", 2048, ""),
		createTestFileEntry("/path/file3.txt", "file3.txt", 512, ""),
	}

	estimate := detector.EstimateDuplicates(files)

	if estimate != 0 {
		t.Errorf("Expected estimate 0, got %d", estimate)
	}
}

// TestGetFileAge tests GetFileAge
func TestGetFileAge(t *testing.T) {
	now := time.Now()
	past := now.Add(-24 * time.Hour)

	entry := models.FileEntry{
		ModTime: past,
	}

	age := detector.GetFileAge(entry)

	// Age should be approximately 24 hours
	expectedMin := 23 * time.Hour
	expectedMax := 25 * time.Hour
	if age < expectedMin || age > expectedMax {
		t.Errorf("GetFileAge = %v, want between %v and %v", age, expectedMin, expectedMax)
	}
}

// TestGetPathDepth tests GetPathDepth
func TestGetPathDepth(t *testing.T) {
	tests := []struct {
		path     string
		expected int
	}{
		{"/file.txt", 1},
		{"/dir/file.txt", 2},
		{"/dir/subdir/file.txt", 3},
		{"/a/b/c/d/file.txt", 5},
		{"file.txt", 0},
		{"/", 1},
		{"", 0},
		{"C:\\Windows\\file.txt", 2}, // Windows path
	}

	for _, tt := range tests {
		depth := detector.GetPathDepth(tt.path)
		if depth != tt.expected {
			t.Errorf("GetPathDepth(%s) = %d, want %d", tt.path, depth, tt.expected)
		}
	}
}

// TestSortBy_Constants tests SortBy constants
func TestSortBy_Constants(t *testing.T) {
	if detector.SortBySize != 0 {
		t.Errorf("SortBySize = %d, want 0", detector.SortBySize)
	}
	if detector.SortByCount != 1 {
		t.Errorf("SortByCount = %d, want 1", detector.SortByCount)
	}
	if detector.SortByRecoverable != 2 {
		t.Errorf("SortByRecoverable = %d, want 2", detector.SortByRecoverable)
	}
	if detector.SortByName != 3 {
		t.Errorf("SortByName = %d, want 3", detector.SortByName)
	}
}
