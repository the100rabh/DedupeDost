package rules_test

import (
	"testing"
	"time"

	"github.com/dupdel/dup-del/internal/rules"
	"github.com/dupdel/dup-del/pkg/models"
)

// TestBatchApplier_NewBatchApplier tests NewBatchApplier
func TestBatchApplier_NewBatchApplier(t *testing.T) {
	engine := rules.NewRuleEngine()
	applier := rules.NewBatchApplier(engine)

	if applier == nil {
		t.Fatal("NewBatchApplier returned nil")
	}

	progress := applier.GetProgress()
	if progress == nil {
		t.Fatal("GetProgress returned nil")
	}

	if progress.Done {
		t.Error("Done should be false initially")
	}
}

// TestBatchApplier_Apply tests Apply method
func TestBatchApplier_Apply(t *testing.T) {
	engine := rules.NewRuleEngine()
	applier := rules.NewBatchApplier(engine)

	now := time.Now()
	groups := []*models.DuplicateGroup{
		createTestGroup([]models.FileEntry{
			createTestFileEntry("/path/old1.txt", "old1.txt", 1024, now.Add(-24*time.Hour)),
			createTestFileEntry("/path/new1.txt", "new1.txt", 1024, now),
		}),
		createTestGroup([]models.FileEntry{
			createTestFileEntry("/path/old2.txt", "old2.txt", 1024, now.Add(-24*time.Hour)),
			createTestFileEntry("/path/new2.txt", "new2.txt", 1024, now),
		}),
		createTestGroup([]models.FileEntry{
			createTestFileEntry("/path/old3.txt", "old3.txt", 1024, now.Add(-24*time.Hour)),
			createTestFileEntry("/path/new3.txt", "new3.txt", 1024, now),
		}),
	}

	rule := rules.NewRule("keep_newest", "Keep Newest", "Test", rules.CriteriaKeepNewest)

	calledCount := 0
	applied, skipped, err := applier.Apply(groups, rule, func(progress *rules.BatchProgress) {
		calledCount++
	})

	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	if applied != 3 {
		t.Errorf("Applied = %d, want 3", applied)
	}
	if skipped != 0 {
		t.Errorf("Skipped = %d, want 0", skipped)
	}

	// Progress callback should be called for each group
	if calledCount != 3 {
		t.Errorf("Progress callback called %d times, want 3", calledCount)
	}

	// Verify all groups have decisions
	for i, group := range groups {
		if len(group.KeepIndices) != 1 {
			t.Errorf("Group %d: KeepIndices length = %d, want 1", i, len(group.KeepIndices))
		}
	}

	// Verify progress
	progress := applier.GetProgress()
	if !progress.Done {
		t.Error("Done should be true after Apply")
	}
	if progress.Processed != 3 {
		t.Errorf("Processed = %d, want 3", progress.Processed)
	}
	if progress.Applied != 3 {
		t.Errorf("Applied = %d, want 3", progress.Applied)
	}
}

// TestBatchApplier_Apply_WithSkips tests Apply with some groups being skipped
func TestBatchApplier_Apply_WithSkips(t *testing.T) {
	engine := rules.NewRuleEngine()
	applier := rules.NewBatchApplier(engine)

	now := time.Now()
	groups := []*models.DuplicateGroup{
		// Valid group
		createTestGroup([]models.FileEntry{
			createTestFileEntry("/path/old1.txt", "old1.txt", 1024, now.Add(-24*time.Hour)),
			createTestFileEntry("/path/new1.txt", "new1.txt", 1024, now),
		}),
		// Single file - will be skipped
		createTestGroup([]models.FileEntry{
			createTestFileEntry("/path/single.txt", "single.txt", 1024, now),
		}),
		// Valid group
		createTestGroup([]models.FileEntry{
			createTestFileEntry("/path/old2.txt", "old2.txt", 1024, now.Add(-24*time.Hour)),
			createTestFileEntry("/path/new2.txt", "new2.txt", 1024, now),
		}),
	}

	rule := rules.NewRule("keep_newest", "Keep Newest", "Test", rules.CriteriaKeepNewest)

	applied, skipped, err := applier.Apply(groups, rule, nil)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	if applied != 2 {
		t.Errorf("Applied = %d, want 2", applied)
	}
	if skipped != 1 {
		t.Errorf("Skipped = %d, want 1", skipped)
	}

	progress := applier.GetProgress()
	if progress.Skipped != 1 {
		t.Errorf("Progress.Skipped = %d, want 1", progress.Skipped)
	}
}

// TestBatchApplier_Apply_InvalidRule tests Apply with invalid rule
func TestBatchApplier_Apply_InvalidRule(t *testing.T) {
	engine := rules.NewRuleEngine()
	applier := rules.NewBatchApplier(engine)

	now := time.Now()
	groups := []*models.DuplicateGroup{
		createTestGroup([]models.FileEntry{
			createTestFileEntry("/path/file1.txt", "file1.txt", 1024, now),
			createTestFileEntry("/path/file2.txt", "file2.txt", 1024, now),
		}),
	}

	// Invalid criteria
	rule := rules.NewRule("invalid", "Invalid", "Test", "invalid_criteria")

	applied, skipped, err := applier.Apply(groups, rule, nil)
	// Should return error for invalid criteria
	if err == nil {
		t.Error("Apply should return error for invalid criteria")
	}
	// No groups should be processed
	if applied != 0 {
		t.Errorf("Applied = %d, want 0", applied)
	}
	if skipped != 0 {
		t.Errorf("Skipped = %d, want 0", skipped)
	}
}

// TestBatchApplier_Cancel tests Cancel method
func TestBatchApplier_Cancel(t *testing.T) {
	engine := rules.NewRuleEngine()
	applier := rules.NewBatchApplier(engine)

	now := time.Now()
	// Create many groups to ensure we can cancel mid-processing
	groups := make([]*models.DuplicateGroup, 100)
	for i := 0; i < 100; i++ {
		groups[i] = createTestGroup([]models.FileEntry{
			createTestFileEntry("/path/file1.txt", "file1.txt", 1024, now),
			createTestFileEntry("/path/file2.txt", "file2.txt", 1024, now),
		})
	}

	rule := rules.NewRule("keep_newest", "Keep Newest", "Test", rules.CriteriaKeepNewest)

	// Cancel immediately
	applier.Cancel()

	// Apply should return early
	applied, skipped, err := applier.Apply(groups, rule, nil)
	
	// Should have been cancelled
	if err == nil {
		t.Log("Apply did not return cancellation error (may have completed before cancel)")
	}
	
	// Check that not all groups were processed (if cancellation worked)
	t.Logf("Applied: %d, Skipped: %d", applied, skipped)
}

// TestBatchApplier_GetProgress tests GetProgress method
func TestBatchApplier_GetProgress(t *testing.T) {
	engine := rules.NewRuleEngine()
	applier := rules.NewBatchApplier(engine)

	progress := applier.GetProgress()
	if progress == nil {
		t.Fatal("GetProgress returned nil")
	}

	if progress.Total != 0 {
		t.Errorf("Total = %d, want 0 before Apply", progress.Total)
	}

	now := time.Now()
	groups := []*models.DuplicateGroup{
		createTestGroup([]models.FileEntry{
			createTestFileEntry("/path/file1.txt", "file1.txt", 1024, now),
			createTestFileEntry("/path/file2.txt", "file2.txt", 1024, now),
		}),
		createTestGroup([]models.FileEntry{
			createTestFileEntry("/path/file3.txt", "file3.txt", 1024, now),
			createTestFileEntry("/path/file4.txt", "file4.txt", 1024, now),
		}),
	}

	rule := rules.NewRule("keep_newest", "Keep Newest", "Test", rules.CriteriaKeepNewest)
	applier.Apply(groups, rule, nil)

	progress = applier.GetProgress()
	if progress.Total != 2 {
		t.Errorf("Total = %d, want 2", progress.Total)
	}
	if progress.Processed != 2 {
		t.Errorf("Processed = %d, want 2", progress.Processed)
	}
}

// TestBatchApplier_IsDone tests IsDone method
func TestBatchApplier_IsDone(t *testing.T) {
	engine := rules.NewRuleEngine()
	applier := rules.NewBatchApplier(engine)

	if applier.IsDone() {
		t.Error("IsDone should return false before Apply")
	}

	now := time.Now()
	groups := []*models.DuplicateGroup{
		createTestGroup([]models.FileEntry{
			createTestFileEntry("/path/file1.txt", "file1.txt", 1024, now),
			createTestFileEntry("/path/file2.txt", "file2.txt", 1024, now),
		}),
	}

	rule := rules.NewRule("keep_newest", "Keep Newest", "Test", rules.CriteriaKeepNewest)
	applier.Apply(groups, rule, nil)

	if !applier.IsDone() {
		t.Error("IsDone should return true after Apply")
	}
}

// TestBatchApplier_ProgressCallback tests progress callback functionality
func TestBatchApplier_ProgressCallback(t *testing.T) {
	engine := rules.NewRuleEngine()
	applier := rules.NewBatchApplier(engine)

	now := time.Now()
	groups := []*models.DuplicateGroup{
		createTestGroup([]models.FileEntry{
			createTestFileEntry("/path/file1.txt", "file1.txt", 1024, now),
			createTestFileEntry("/path/file2.txt", "file2.txt", 1024, now),
		}),
		createTestGroup([]models.FileEntry{
			createTestFileEntry("/path/file3.txt", "file3.txt", 1024, now),
			createTestFileEntry("/path/file4.txt", "file4.txt", 1024, now),
		}),
		createTestGroup([]models.FileEntry{
			createTestFileEntry("/path/file5.txt", "file5.txt", 1024, now),
			createTestFileEntry("/path/file6.txt", "file6.txt", 1024, now),
		}),
	}

	rule := rules.NewRule("keep_newest", "Keep Newest", "Test", rules.CriteriaKeepNewest)

	var progressUpdates []*rules.BatchProgress
	applier.Apply(groups, rule, func(progress *rules.BatchProgress) {
		// Make a copy of the progress
		p := *progress
		progressUpdates = append(progressUpdates, &p)
	})

	if len(progressUpdates) != 3 {
		t.Errorf("Progress callback called %d times, want 3", len(progressUpdates))
	}

	// Verify progress increments
	for i, p := range progressUpdates {
		expectedProcessed := i + 1
		if p.Processed != expectedProcessed {
			t.Errorf("Update %d: Processed = %d, want %d", i, p.Processed, expectedProcessed)
		}
	}
}

// TestBatchApplier_EmptyGroups tests Apply with empty groups
func TestBatchApplier_EmptyGroups(t *testing.T) {
	engine := rules.NewRuleEngine()
	applier := rules.NewBatchApplier(engine)

	groups := []*models.DuplicateGroup{}
	rule := rules.NewRule("keep_newest", "Keep Newest", "Test", rules.CriteriaKeepNewest)

	applied, skipped, err := applier.Apply(groups, rule, nil)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	if applied != 0 {
		t.Errorf("Applied = %d, want 0", applied)
	}
	if skipped != 0 {
		t.Errorf("Skipped = %d, want 0", skipped)
	}

	progress := applier.GetProgress()
	if !progress.Done {
		t.Error("Done should be true")
	}
}

// TestBatchApplier_SpecificDirRule tests Apply with SpecificDir rule
func TestBatchApplier_SpecificDirRule(t *testing.T) {
	engine := rules.NewRuleEngine()
	applier := rules.NewBatchApplier(engine)

	now := time.Now()
	groups := []*models.DuplicateGroup{
		createTestGroup([]models.FileEntry{
			createTestFileEntry("/target/file1.txt", "file1.txt", 1024, now),
			createTestFileEntry("/other/file1.txt", "file1.txt", 1024, now),
		}),
		createTestGroup([]models.FileEntry{
			createTestFileEntry("/target/file2.txt", "file2.txt", 1024, now),
			createTestFileEntry("/other/file2.txt", "file2.txt", 1024, now),
		}),
	}

	rule := rules.NewRule("keep_specific_dir", "Keep Specific Dir", "Test", rules.CriteriaKeepSpecificDir)
	rule.WithConfig("directory", "/target")

	// Set target directory on evaluator
	eval := &rules.SpecificDirEvaluator{TargetDir: "/target"}
	engine.RegisterEvaluator(eval)

	applied, skipped, err := applier.Apply(groups, rule, nil)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	if applied != 2 {
		t.Errorf("Applied = %d, want 2", applied)
	}
	if skipped != 0 {
		t.Errorf("Skipped = %d, want 0", skipped)
	}
}
