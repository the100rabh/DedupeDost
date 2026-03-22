package rules_test

import (
	"testing"
	"time"

	"github.com/dupdel/dup-del/internal/rules"
	"github.com/dupdel/dup-del/pkg/models"
	"github.com/dupdel/dup-del/pkg/types"
)

// Helper function to create test file entries
func createTestFileEntry(path, name string, size int64, modTime time.Time) models.FileEntry {
	return models.FileEntry{
		Path:        path,
		Name:        name,
		Size:        size,
		Extension:   ".txt",
		ModTime:     modTime,
		CreatedTime: modTime,
		FileType:    types.FileTypeText,
		IsHidden:    false,
		IsSymlink:   false,
	}
}

func createTestGroup(files []models.FileEntry) *models.DuplicateGroup {
	return &models.DuplicateGroup{
		ID:        "test-group",
		Hash:      "abc123",
		Size:      files[0].Size,
		Extension: ".txt",
		FileType:  types.FileTypeText,
		Files:     files,
	}
}

// TestBuiltinRules tests that all builtin rules are defined
func TestBuiltinRules(t *testing.T) {
	rules := rules.GetAllRules()

	if len(rules) == 0 {
		t.Fatal("GetAllRules returned empty list")
	}

	// Check for expected rules
	expectedIDs := []string{
		"keep_newest",
		"keep_oldest",
		"keep_shortest_path",
		"keep_longest_path",
		"keep_specific_dir",
		"keep_largest",
		"keep_smallest",
		"manual",
	}

	foundIDs := make(map[string]bool)
	for _, rule := range rules {
		foundIDs[rule.ID] = true
	}

	for _, expectedID := range expectedIDs {
		if !foundIDs[expectedID] {
			t.Errorf("Expected rule ID %s not found", expectedID)
		}
	}
}

// TestGetRuleByID tests GetRuleByID function
func TestGetRuleByID(t *testing.T) {
	rule := rules.GetRuleByID("keep_newest")
	if rule == nil {
		t.Fatal("GetRuleByID returned nil for keep_newest")
	}

	if rule.ID != "keep_newest" {
		t.Errorf("ID = %s, want keep_newest", rule.ID)
	}

	if rule.Name != "Keep Newest" {
		t.Errorf("Name = %s, want Keep Newest", rule.Name)
	}

	// Non-existent rule should return nil
	rule = rules.GetRuleByID("nonexistent")
	if rule != nil {
		t.Error("GetRuleByID should return nil for non-existent ID")
	}
}

// TestGetRulesByCriteria tests GetRulesByCriteria function
func TestGetRulesByCriteria(t *testing.T) {
	ruleList := rules.GetRulesByCriteria(rules.CriteriaKeepNewest)
	if len(ruleList) == 0 {
		t.Error("GetRulesByCriteria returned empty list for CriteriaKeepNewest")
	}

	found := false
	for _, r := range ruleList {
		if r.Criteria == rules.CriteriaKeepNewest {
			found = true
			break
		}
	}
	if !found {
		t.Error("No rule with CriteriaKeepNewest found")
	}
}

// TestNewRule tests creating a new rule
func TestNewRule(t *testing.T) {
	rule := rules.NewRule("test_rule", "Test Rule", "Test description", rules.CriteriaKeepNewest)

	if rule == nil {
		t.Fatal("NewRule returned nil")
	}

	if rule.ID != "test_rule" {
		t.Errorf("ID = %s, want test_rule", rule.ID)
	}

	if rule.Name != "Test Rule" {
		t.Errorf("Name = %s, want Test Rule", rule.Name)
	}

	if rule.Criteria != rules.CriteriaKeepNewest {
		t.Errorf("Criteria = %s, want CriteriaKeepNewest", rule.Criteria)
	}
}

// TestRule_WithConfig tests WithConfig method
func TestRule_WithConfig(t *testing.T) {
	rule := rules.NewRule("test", "Test", "Test", rules.CriteriaKeepNewest)

	result := rule.WithConfig("key", "value")
	if result != rule {
		t.Error("WithConfig should return same rule for chaining")
	}

	val, ok := rule.Config["key"]
	if !ok {
		t.Error("Config key not found")
	}
	if val != "value" {
		t.Errorf("Config value = %v, want value", val)
	}
}

// TestRule_WithApplyToAll tests WithApplyToAll method
func TestRule_WithApplyToAll(t *testing.T) {
	rule := rules.NewRule("test", "Test", "Test", rules.CriteriaKeepNewest)

	if rule.ApplyToAll {
		t.Error("ApplyToAll should be false by default")
	}

	rule.WithApplyToAll(true)
	if !rule.ApplyToAll {
		t.Error("ApplyToAll should be true after WithApplyToAll(true)")
	}
}

// TestRuleEngine_NewRuleEngine tests NewRuleEngine
func TestRuleEngine_NewRuleEngine(t *testing.T) {
	engine := rules.NewRuleEngine()
	if engine == nil {
		t.Fatal("NewRuleEngine returned nil")
	}
}

// TestRuleEngine_ApplyRule_KeepNewest tests ApplyRule with KeepNewest
func TestRuleEngine_ApplyRule_KeepNewest(t *testing.T) {
	now := time.Now()
	files := []models.FileEntry{
		createTestFileEntry("/path/old.txt", "old.txt", 1024, now.Add(-24*time.Hour)),
		createTestFileEntry("/path/new.txt", "new.txt", 1024, now),
		createTestFileEntry("/path/older.txt", "older.txt", 1024, now.Add(-48*time.Hour)),
	}
	group := createTestGroup(files)

	engine := rules.NewRuleEngine()
	rule := rules.NewRule("keep_newest", "Keep Newest", "Test", rules.CriteriaKeepNewest)

	err := engine.ApplyRule(group, rule)
	if err != nil {
		t.Fatalf("ApplyRule failed: %v", err)
	}

	// Should keep index 1 (new.txt)
	if len(group.KeepIndices) != 1 {
		t.Fatalf("KeepIndices length = %d, want 1", len(group.KeepIndices))
	}
	if group.KeepIndices[0] != 1 {
		t.Errorf("KeepIndices[0] = %d, want 1", group.KeepIndices[0])
	}

	// Should have 2 files marked for deletion
	if len(group.DeletePaths) != 2 {
		t.Errorf("DeletePaths length = %d, want 2", len(group.DeletePaths))
	}
}

// TestRuleEngine_ApplyRule_KeepOldest tests ApplyRule with KeepOldest
func TestRuleEngine_ApplyRule_KeepOldest(t *testing.T) {
	now := time.Now()
	files := []models.FileEntry{
		createTestFileEntry("/path/new.txt", "new.txt", 1024, now),
		createTestFileEntry("/path/old.txt", "old.txt", 1024, now.Add(-24*time.Hour)),
		createTestFileEntry("/path/newest.txt", "newest.txt", 1024, now.Add(1*time.Hour)),
	}
	group := createTestGroup(files)

	engine := rules.NewRuleEngine()
	rule := rules.NewRule("keep_oldest", "Keep Oldest", "Test", rules.CriteriaKeepOldest)

	err := engine.ApplyRule(group, rule)
	if err != nil {
		t.Fatalf("ApplyRule failed: %v", err)
	}

	// Should keep index 1 (old.txt)
	if len(group.KeepIndices) != 1 {
		t.Fatalf("KeepIndices length = %d, want 1", len(group.KeepIndices))
	}
	if group.KeepIndices[0] != 1 {
		t.Errorf("KeepIndices[0] = %d, want 1", group.KeepIndices[0])
	}
}

// TestRuleEngine_ApplyRule_ShortestPath tests ApplyRule with ShortestPath
func TestRuleEngine_ApplyRule_ShortestPath(t *testing.T) {
	now := time.Now()
	files := []models.FileEntry{
		createTestFileEntry("/a/b/c/deep/file.txt", "file.txt", 1024, now),
		createTestFileEntry("/file.txt", "file.txt", 1024, now),
		createTestFileEntry("/a/b/shallow/file.txt", "file.txt", 1024, now),
	}
	group := createTestGroup(files)

	engine := rules.NewRuleEngine()
	rule := rules.NewRule("keep_shortest_path", "Keep Shortest Path", "Test", rules.CriteriaKeepShortestPath)

	err := engine.ApplyRule(group, rule)
	if err != nil {
		t.Fatalf("ApplyRule failed: %v", err)
	}

	// Should keep index 1 (/file.txt - shortest path)
	if len(group.KeepIndices) != 1 {
		t.Fatalf("KeepIndices length = %d, want 1", len(group.KeepIndices))
	}
	if group.KeepIndices[0] != 1 {
		t.Errorf("KeepIndices[0] = %d, want 1", group.KeepIndices[0])
	}
}

// TestRuleEngine_ApplyRule_LongestPath tests ApplyRule with LongestPath
func TestRuleEngine_ApplyRule_LongestPath(t *testing.T) {
	now := time.Now()
	files := []models.FileEntry{
		createTestFileEntry("/file.txt", "file.txt", 1024, now),
		createTestFileEntry("/a/b/c/deep/file.txt", "file.txt", 1024, now),
		createTestFileEntry("/a/b/shallow/file.txt", "file.txt", 1024, now),
	}
	group := createTestGroup(files)

	engine := rules.NewRuleEngine()
	rule := rules.NewRule("keep_longest_path", "Keep Longest Path", "Test", rules.CriteriaKeepLongestPath)

	err := engine.ApplyRule(group, rule)
	if err != nil {
		t.Fatalf("ApplyRule failed: %v", err)
	}

	// Should keep index 1 (/a/b/c/deep/file.txt - longest path)
	if len(group.KeepIndices) != 1 {
		t.Fatalf("KeepIndices length = %d, want 1", len(group.KeepIndices))
	}
	if group.KeepIndices[0] != 1 {
		t.Errorf("KeepIndices[0] = %d, want 1", group.KeepIndices[0])
	}
}

// TestRuleEngine_ApplyRule_SpecificDir tests ApplyRule with SpecificDir
func TestRuleEngine_ApplyRule_SpecificDir(t *testing.T) {
	now := time.Now()
	files := []models.FileEntry{
		createTestFileEntry("/other/file.txt", "file.txt", 1024, now),
		createTestFileEntry("/target/file.txt", "file.txt", 1024, now),
		createTestFileEntry("/another/file.txt", "file.txt", 1024, now),
	}
	group := createTestGroup(files)

	engine := rules.NewRuleEngine()
	rule := rules.NewRule("keep_specific_dir", "Keep Specific Dir", "Test", rules.CriteriaKeepSpecificDir)
	rule.WithConfig("directory", "/target")

	// Set target directory on evaluator
	eval := &rules.SpecificDirEvaluator{TargetDir: "/target"}
	engine.RegisterEvaluator(eval)

	err := engine.ApplyRule(group, rule)
	if err != nil {
		t.Fatalf("ApplyRule failed: %v", err)
	}

	// Should keep index 1 (/target/file.txt)
	if len(group.KeepIndices) != 1 {
		t.Fatalf("KeepIndices length = %d, want 1", len(group.KeepIndices))
	}
	if group.KeepIndices[0] != 1 {
		t.Errorf("KeepIndices[0] = %d, want 1", group.KeepIndices[0])
	}
}

// TestRuleEngine_ApplyRule_Largest tests ApplyRule with Largest
func TestRuleEngine_ApplyRule_Largest(t *testing.T) {
	now := time.Now()
	files := []models.FileEntry{
		createTestFileEntry("/path/small.txt", "small.txt", 512, now),
		createTestFileEntry("/path/large.txt", "large.txt", 2048, now),
		createTestFileEntry("/path/medium.txt", "medium.txt", 1024, now),
	}
	group := createTestGroup(files)

	engine := rules.NewRuleEngine()
	rule := rules.NewRule("keep_largest", "Keep Largest", "Test", rules.CriteriaKeepLargest)

	err := engine.ApplyRule(group, rule)
	if err != nil {
		t.Fatalf("ApplyRule failed: %v", err)
	}

	// Should keep index 1 (large.txt - 2048 bytes)
	if len(group.KeepIndices) != 1 {
		t.Fatalf("KeepIndices length = %d, want 1", len(group.KeepIndices))
	}
	if group.KeepIndices[0] != 1 {
		t.Errorf("KeepIndices[0] = %d, want 1", group.KeepIndices[0])
	}
}

// TestRuleEngine_ApplyRule_Smallest tests ApplyRule with Smallest
func TestRuleEngine_ApplyRule_Smallest(t *testing.T) {
	now := time.Now()
	files := []models.FileEntry{
		createTestFileEntry("/path/large.txt", "large.txt", 2048, now),
		createTestFileEntry("/path/small.txt", "small.txt", 512, now),
		createTestFileEntry("/path/medium.txt", "medium.txt", 1024, now),
	}
	group := createTestGroup(files)

	engine := rules.NewRuleEngine()
	rule := rules.NewRule("keep_smallest", "Keep Smallest", "Test", rules.CriteriaKeepSmallest)

	err := engine.ApplyRule(group, rule)
	if err != nil {
		t.Fatalf("ApplyRule failed: %v", err)
	}

	// Should keep index 1 (small.txt - 512 bytes)
	if len(group.KeepIndices) != 1 {
		t.Fatalf("KeepIndices length = %d, want 1", len(group.KeepIndices))
	}
	if group.KeepIndices[0] != 1 {
		t.Errorf("KeepIndices[0] = %d, want 1", group.KeepIndices[0])
	}
}

// TestRuleEngine_ApplyRule_NoEvaluator tests error when no evaluator
func TestRuleEngine_ApplyRule_NoEvaluator(t *testing.T) {
	files := []models.FileEntry{
		createTestFileEntry("/path/file1.txt", "file1.txt", 1024, time.Now()),
		createTestFileEntry("/path/file2.txt", "file2.txt", 1024, time.Now()),
	}
	group := createTestGroup(files)

	engine := rules.NewRuleEngine()
	rule := rules.NewRule("invalid", "Invalid", "Test", "invalid_criteria")

	err := engine.ApplyRule(group, rule)
	if err == nil {
		t.Error("ApplyRule should return error for invalid criteria")
	}
}

// TestRuleEngine_ApplyRule_NotApplicable tests when rule is not applicable
func TestRuleEngine_ApplyRule_NotApplicable(t *testing.T) {
	// Single file - rule not applicable
	files := []models.FileEntry{
		createTestFileEntry("/path/file1.txt", "file1.txt", 1024, time.Now()),
	}
	group := createTestGroup(files)

	engine := rules.NewRuleEngine()
	rule := rules.NewRule("keep_newest", "Keep Newest", "Test", rules.CriteriaKeepNewest)

	err := engine.ApplyRule(group, rule)
	if err == nil {
		t.Error("ApplyRule should return error when rule not applicable")
	}
}

// TestRuleEngine_ApplyToAll tests ApplyToAll method
func TestRuleEngine_ApplyToAll(t *testing.T) {
	now := time.Now()
	
	// Create multiple groups
	groups := []*models.DuplicateGroup{
		createTestGroup([]models.FileEntry{
			createTestFileEntry("/path/old1.txt", "old1.txt", 1024, now.Add(-24*time.Hour)),
			createTestFileEntry("/path/new1.txt", "new1.txt", 1024, now),
		}),
		createTestGroup([]models.FileEntry{
			createTestFileEntry("/path/old2.txt", "old2.txt", 1024, now.Add(-24*time.Hour)),
			createTestFileEntry("/path/new2.txt", "new2.txt", 1024, now),
		}),
	}

	engine := rules.NewRuleEngine()
	rule := rules.NewRule("keep_newest", "Keep Newest", "Test", rules.CriteriaKeepNewest)

	applied, skipped, err := engine.ApplyToAll(groups, rule)
	if err != nil {
		t.Fatalf("ApplyToAll failed: %v", err)
	}

	if applied != 2 {
		t.Errorf("Applied = %d, want 2", applied)
	}
	if skipped != 0 {
		t.Errorf("Skipped = %d, want 0", skipped)
	}

	// Verify both groups have decisions
	for i, group := range groups {
		if len(group.KeepIndices) != 1 {
			t.Errorf("Group %d: KeepIndices length = %d, want 1", i, len(group.KeepIndices))
		}
	}
}

// TestRuleEngine_GetHistory tests GetHistory method
func TestRuleEngine_GetHistory(t *testing.T) {
	engine := rules.NewRuleEngine()
	
	// Initial history should be empty
	history := engine.GetHistory()
	if len(history) != 0 {
		t.Errorf("Initial history length = %d, want 0", len(history))
	}

	// Apply a rule
	now := time.Now()
	files := []models.FileEntry{
		createTestFileEntry("/path/old.txt", "old.txt", 1024, now.Add(-24*time.Hour)),
		createTestFileEntry("/path/new.txt", "new.txt", 1024, now),
	}
	group := createTestGroup(files)
	rule := rules.NewRule("keep_newest", "Keep Newest", "Test", rules.CriteriaKeepNewest)

	engine.ApplyRule(group, rule)

	// History should have one entry
	history = engine.GetHistory()
	if len(history) != 1 {
		t.Errorf("History length = %d, want 1", len(history))
	}
}

// TestRuleEngine_UndoLast tests UndoLast method
func TestRuleEngine_UndoLast(t *testing.T) {
	now := time.Now()
	files := []models.FileEntry{
		createTestFileEntry("/path/old.txt", "old.txt", 1024, now.Add(-24*time.Hour)),
		createTestFileEntry("/path/new.txt", "new.txt", 1024, now),
	}
	group := createTestGroup(files)

	engine := rules.NewRuleEngine()
	rule := rules.NewRule("keep_newest", "Keep Newest", "Test", rules.CriteriaKeepNewest)

	// Apply rule
	engine.ApplyRule(group, rule)
	if len(group.KeepIndices) == 0 {
		t.Fatal("Rule was not applied")
	}

	// Undo
	success := engine.UndoLast([]*models.DuplicateGroup{group})
	if !success {
		t.Error("UndoLast should return true")
	}

	// Decisions should be cleared
	if len(group.KeepIndices) != 0 {
		t.Errorf("KeepIndices not cleared, length = %d", len(group.KeepIndices))
	}
	if len(group.DeletePaths) != 0 {
		t.Errorf("DeletePaths not cleared, length = %d", len(group.DeletePaths))
	}
}

// TestRuleEngine_UndoLast_EmptyHistory tests UndoLast with empty history
func TestRuleEngine_UndoLast_EmptyHistory(t *testing.T) {
	engine := rules.NewRuleEngine()
	
	success := engine.UndoLast([]*models.DuplicateGroup{})
	if success {
		t.Error("UndoLast should return false with empty history")
	}
}

// TestRuleEngine_ClearHistory tests ClearHistory method
func TestRuleEngine_ClearHistory(t *testing.T) {
	engine := rules.NewRuleEngine()
	
	// Apply a rule to add history
	now := time.Now()
	files := []models.FileEntry{
		createTestFileEntry("/path/old.txt", "old.txt", 1024, now.Add(-24*time.Hour)),
		createTestFileEntry("/path/new.txt", "new.txt", 1024, now),
	}
	group := createTestGroup(files)
	rule := rules.NewRule("keep_newest", "Keep Newest", "Test", rules.CriteriaKeepNewest)
	engine.ApplyRule(group, rule)

	// Clear history
	engine.ClearHistory()
	
	history := engine.GetHistory()
	if len(history) != 0 {
		t.Errorf("History length after clear = %d, want 0", len(history))
	}
}

// TestNewestEvaluator tests NewestEvaluator
func TestNewestEvaluator(t *testing.T) {
	eval := &rules.NewestEvaluator{}
	
	now := time.Now()
	files := []models.FileEntry{
		createTestFileEntry("/path/old.txt", "old.txt", 1024, now.Add(-24*time.Hour)),
		createTestFileEntry("/path/new.txt", "new.txt", 1024, now),
		createTestFileEntry("/path/older.txt", "older.txt", 1024, now.Add(-48*time.Hour)),
	}
	group := createTestGroup(files)

	idx, err := eval.Evaluate(group)
	if err != nil {
		t.Fatalf("Evaluate failed: %v", err)
	}
	if idx != 1 {
		t.Errorf("Evaluate = %d, want 1", idx)
	}

	if !eval.IsApplicable(group) {
		t.Error("IsApplicable should return true")
	}

	desc := eval.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}
}

// TestOldestEvaluator tests OldestEvaluator
func TestOldestEvaluator(t *testing.T) {
	eval := &rules.OldestEvaluator{}
	
	now := time.Now()
	files := []models.FileEntry{
		createTestFileEntry("/path/new.txt", "new.txt", 1024, now),
		createTestFileEntry("/path/old.txt", "old.txt", 1024, now.Add(-24*time.Hour)),
	}
	group := createTestGroup(files)

	idx, err := eval.Evaluate(group)
	if err != nil {
		t.Fatalf("Evaluate failed: %v", err)
	}
	if idx != 1 {
		t.Errorf("Evaluate = %d, want 1", idx)
	}
}

// TestShortestPathEvaluator tests ShortestPathEvaluator
func TestShortestPathEvaluator(t *testing.T) {
	eval := &rules.ShortestPathEvaluator{}
	
	now := time.Now()
	files := []models.FileEntry{
		createTestFileEntry("/a/b/c/deep/file.txt", "file.txt", 1024, now),
		createTestFileEntry("/file.txt", "file.txt", 1024, now),
	}
	group := createTestGroup(files)

	idx, err := eval.Evaluate(group)
	if err != nil {
		t.Fatalf("Evaluate failed: %v", err)
	}
	if idx != 1 {
		t.Errorf("Evaluate = %d, want 1", idx)
	}
}

// TestLongestPathEvaluator tests LongestPathEvaluator
func TestLongestPathEvaluator(t *testing.T) {
	eval := &rules.LongestPathEvaluator{}
	
	now := time.Now()
	files := []models.FileEntry{
		createTestFileEntry("/file.txt", "file.txt", 1024, now),
		createTestFileEntry("/a/b/c/deep/file.txt", "file.txt", 1024, now),
	}
	group := createTestGroup(files)

	idx, err := eval.Evaluate(group)
	if err != nil {
		t.Fatalf("Evaluate failed: %v", err)
	}
	if idx != 1 {
		t.Errorf("Evaluate = %d, want 1", idx)
	}
}

// TestSpecificDirEvaluator tests SpecificDirEvaluator
func TestSpecificDirEvaluator(t *testing.T) {
	eval := &rules.SpecificDirEvaluator{TargetDir: "/target"}
	
	now := time.Now()
	files := []models.FileEntry{
		createTestFileEntry("/other/file.txt", "file.txt", 1024, now),
		createTestFileEntry("/target/file.txt", "file.txt", 1024, now),
	}
	group := createTestGroup(files)

	idx, err := eval.Evaluate(group)
	if err != nil {
		t.Fatalf("Evaluate failed: %v", err)
	}
	if idx != 1 {
		t.Errorf("Evaluate = %d, want 1", idx)
	}

	eval.SetDirectory("/other")
	if eval.TargetDir != "/other" {
		t.Errorf("TargetDir = %s, want /other", eval.TargetDir)
	}
}

// TestLargestEvaluator tests LargestEvaluator
func TestLargestEvaluator(t *testing.T) {
	eval := &rules.LargestEvaluator{}
	
	now := time.Now()
	files := []models.FileEntry{
		createTestFileEntry("/path/small.txt", "small.txt", 512, now),
		createTestFileEntry("/path/large.txt", "large.txt", 2048, now),
	}
	group := createTestGroup(files)

	idx, err := eval.Evaluate(group)
	if err != nil {
		t.Fatalf("Evaluate failed: %v", err)
	}
	if idx != 1 {
		t.Errorf("Evaluate = %d, want 1", idx)
	}
}

// TestSmallestEvaluator tests SmallestEvaluator
func TestSmallestEvaluator(t *testing.T) {
	eval := &rules.SmallestEvaluator{}
	
	now := time.Now()
	files := []models.FileEntry{
		createTestFileEntry("/path/large.txt", "large.txt", 2048, now),
		createTestFileEntry("/path/small.txt", "small.txt", 512, now),
	}
	group := createTestGroup(files)

	idx, err := eval.Evaluate(group)
	if err != nil {
		t.Fatalf("Evaluate failed: %v", err)
	}
	if idx != 1 {
		t.Errorf("Evaluate = %d, want 1", idx)
	}
}

// TestEvaluators_EmptyGroup tests evaluators with empty group
func TestEvaluators_EmptyGroup(t *testing.T) {
	group := &models.DuplicateGroup{
		ID:        "empty",
		Files:     []models.FileEntry{},
	}

	evaluators := []rules.RuleEvaluator{
		&rules.NewestEvaluator{},
		&rules.OldestEvaluator{},
		&rules.ShortestPathEvaluator{},
		&rules.LongestPathEvaluator{},
		&rules.SpecificDirEvaluator{TargetDir: "/test"},
		&rules.LargestEvaluator{},
		&rules.SmallestEvaluator{},
	}

	for _, eval := range evaluators {
		_, err := eval.Evaluate(group)
		if err == nil {
			t.Errorf("%T should return error for empty group", eval)
		}
	}
}

// TestEvaluators_IsApplicable tests IsApplicable for all evaluators
func TestEvaluators_IsApplicable(t *testing.T) {
	now := time.Now()
	
	// Single file - not applicable
	singleFileGroup := createTestGroup([]models.FileEntry{
		createTestFileEntry("/path/file1.txt", "file1.txt", 1024, now),
	})
	
	// Multiple files - applicable
	multiFileGroup := createTestGroup([]models.FileEntry{
		createTestFileEntry("/path/file1.txt", "file1.txt", 1024, now),
		createTestFileEntry("/path/file2.txt", "file2.txt", 1024, now),
	})

	evaluators := []struct {
		eval rules.RuleEvaluator
		name string
	}{
		{&rules.NewestEvaluator{}, "NewestEvaluator"},
		{&rules.OldestEvaluator{}, "OldestEvaluator"},
		{&rules.ShortestPathEvaluator{}, "ShortestPathEvaluator"},
		{&rules.LongestPathEvaluator{}, "LongestPathEvaluator"},
		{&rules.LargestEvaluator{}, "LargestEvaluator"},
		{&rules.SmallestEvaluator{}, "SmallestEvaluator"},
	}

	for _, e := range evaluators {
		if e.eval.IsApplicable(singleFileGroup) {
			t.Errorf("%s: IsApplicable should return false for single file", e.name)
		}
		if !e.eval.IsApplicable(multiFileGroup) {
			t.Errorf("%s: IsApplicable should return true for multiple files", e.name)
		}
	}
}
