package rules

import (
	"context"
	"strings"

	"github.com/dupdel/dup-del/pkg/models"
)

// BatchApplier applies rules to multiple groups
type BatchApplier struct {
	engine      *RuleEngine
	progress    *BatchProgress
	ctx         context.Context
	cancel      context.CancelFunc
}

// BatchProgress tracks batch operation progress
type BatchProgress struct {
	Total       int
	Processed   int
	Applied     int
	Skipped     int
	CurrentGroup string
	Done        bool
}

// NewBatchApplier creates a new batch applier
func NewBatchApplier(engine *RuleEngine) *BatchApplier {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &BatchApplier{
		engine:   engine,
		progress: &BatchProgress{},
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Apply applies a rule to multiple groups
func (b *BatchApplier) Apply(groups []models.DuplicateGroup, rule *Rule, progressCallback func(*BatchProgress)) (int, int, error) {
	b.progress = &BatchProgress{
		Total:     len(groups),
		Processed: 0,
		Applied:   0,
		Skipped:   0,
		Done:      false,
	}
	
	evaluator := b.getEvaluator(rule.Criteria)
	if evaluator == nil {
		return 0, 0, ErrNoEvaluator
	}
	
	for i := range groups {
		// Check for cancellation
		select {
		case <-b.ctx.Done():
			return b.progress.Applied, b.progress.Skipped, b.ctx.Err()
		default:
		}
		
		group := &groups[i]
		b.progress.CurrentGroup = group.ID
		
		if !evaluator.IsApplicable(group) {
			b.progress.Skipped++
			b.progress.Processed++
			
			if progressCallback != nil {
				progressCallback(b.progress)
			}
			continue
		}
		
		keepIndex, err := evaluator.Evaluate(group)
		if err != nil {
			b.progress.Skipped++
			b.progress.Processed++
			
			if progressCallback != nil {
				progressCallback(b.progress)
			}
			continue
		}
		
		group.MarkForDeletion([]int{keepIndex})
		group.RuleApplied = &models.Rule{
			ID:       rule.ID,
			Name:     rule.Name,
			Criteria: models.RuleCriteria(rule.Criteria),
		}
		
		b.progress.Applied++
		b.progress.Processed++
		
		if progressCallback != nil {
			progressCallback(b.progress)
		}
	}
	
	b.progress.Done = true
	return b.progress.Applied, b.progress.Skipped, nil
}

// getEvaluator gets the evaluator for a criteria
func (b *BatchApplier) getEvaluator(criteria RuleCriteria) RuleEvaluator {
	evaluators := map[RuleCriteria]RuleEvaluator{
		CriteriaKeepNewest:       &NewestEvaluator{},
		CriteriaKeepOldest:       &OldestEvaluator{},
		CriteriaKeepShortestPath: &ShortestPathEvaluator{},
		CriteriaKeepLongestPath:  &LongestPathEvaluator{},
		CriteriaKeepLargest:      &LargestEvaluator{},
		CriteriaKeepSmallest:     &SmallestEvaluator{},
	}
	return evaluators[criteria]
}

// Cancel cancels the batch operation
func (b *BatchApplier) Cancel() {
	b.cancel()
}

// GetProgress returns the current progress
func (b *BatchApplier) GetProgress() *BatchProgress {
	return b.progress
}

// IsDone returns whether the batch operation is complete
func (b *BatchApplier) IsDone() bool {
	return b.progress.Done
}

// Error definitions
var (
	ErrNoEvaluator = &RuleError{"no evaluator available for criteria"}
)

// RuleError represents a rule-related error
type RuleError struct {
	Message string
}

func (e *RuleError) Error() string {
	return e.Message
}

// SmartSelector provides smart auto-selection based on patterns
type SmartSelector struct {
	patterns []SelectionPattern
}

// SelectionPattern defines a pattern for smart selection
type SelectionPattern struct {
	Name        string
	Description string
	MatchFunc   func(models.FileEntry) bool
	Priority    int
	Score       int
}

// NewSmartSelector creates a new smart selector
func NewSmartSelector() *SmartSelector {
	ss := &SmartSelector{
		patterns: []SelectionPattern{
			{
				Name:        "Keep Documents Folder",
				Description: "Prefer files in Documents folders",
				MatchFunc:   matchFolderPattern("Documents", "docs", "documents"),
				Priority:    1,
				Score:       10,
			},
			{
				Name:        "Keep Photos Folder",
				Description: "Prefer files in Photos folders",
				MatchFunc:   matchFolderPattern("Photos", "Pictures", "images"),
				Priority:    1,
				Score:       10,
			},
			{
				Name:        "Avoid Backup Folders",
				Description: "Avoid files in Backup/Old folders",
				MatchFunc:   matchFolderPattern("Backup", "Backups", "Old", "Archive"),
				Priority:    2,
				Score:       -10,
			},
			{
				Name:        "Avoid Temp Folders",
				Description: "Avoid files in Temp/Cache folders",
				MatchFunc:   matchFolderPattern("Temp", "tmp", "cache", ".cache"),
				Priority:    2,
				Score:       -15,
			},
			{
				Name:        "Clean Filename",
				Description: "Prefer files without 'copy' or 'duplicate' in name",
				MatchFunc:   matchCleanFilename,
				Priority:    3,
				Score:       5,
			},
		},
	}
	
	return ss
}

// SelectBest selects the best file from a group based on patterns
func (ss *SmartSelector) SelectBest(group *models.DuplicateGroup) (keepIndex int, score int, reason string) {
	if len(group.Files) == 0 {
		return 0, 0, "no files"
	}
	
	if len(group.Files) == 1 {
		return 0, 0, "only one file"
	}
	
	bestIdx := 0
	bestScore := 0
	bestReason := "default selection"
	
	for i, file := range group.Files {
		fileScore := ss.scoreFile(file)
		
		if fileScore > bestScore || i == 0 {
			bestScore = fileScore
			bestIdx = i
			bestReason = ss.getReason(file)
		}
	}
	
	return bestIdx, bestScore, bestReason
}

// scoreFile scores a file based on patterns
func (ss *SmartSelector) scoreFile(file models.FileEntry) int {
	score := 0
	
	for _, pattern := range ss.patterns {
		if pattern.MatchFunc(file) {
			score += pattern.Score
		}
	}
	
	return score
}

// getReason gets the reason for selection
func (ss *SmartSelector) getReason(file models.FileEntry) string {
	for _, pattern := range ss.patterns {
		if pattern.MatchFunc(file) {
			return pattern.Description
		}
	}
	return "default selection"
}

// matchFolderPattern creates a match function for folder patterns
func matchFolderPattern(patterns ...string) func(models.FileEntry) bool {
	return func(file models.FileEntry) bool {
		path := strings.ToLower(file.Path)
		for _, p := range patterns {
			if strings.Contains(path, strings.ToLower(p)) {
				return true
			}
		}
		return false
	}
}

// matchCleanFilename checks if filename is clean
func matchCleanFilename(file models.FileEntry) bool {
	name := strings.ToLower(file.Name)
	
	// Check for copy/duplicate patterns
	badPatterns := []string{"copy", "duplicate", "dup", "backup", "old", "new"}
	for _, p := range badPatterns {
		if strings.Contains(name, p) {
			return false
		}
	}
	
	// Check for numbered copies
	if strings.Contains(name, "(1)") || strings.Contains(name, "(2)") {
		return false
	}
	
	return true
}
