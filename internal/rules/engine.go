package rules

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dupdel/dup-del/pkg/models"
)
type RuleEvaluator interface {
	// Evaluate determines which file to keep
	Evaluate(group *models.DuplicateGroup) (keepIndex int, err error)
	
	// Description returns a description of what this evaluator does
	Description() string
	
	// IsApplicable checks if this evaluator can be applied to the group
	IsApplicable(group *models.DuplicateGroup) bool
}

// RuleEngine applies rules to duplicate groups
type RuleEngine struct {
	evaluators map[RuleCriteria]RuleEvaluator
	history    []RuleApplication
}

// RuleApplication records a rule application
type RuleApplication struct {
	GroupID   string
	RuleID    string
	KeepIndex int
}

// NewRuleEngine creates a new rule engine
func NewRuleEngine() *RuleEngine {
	engine := &RuleEngine{
		evaluators: make(map[RuleCriteria]RuleEvaluator),
		history:    make([]RuleApplication, 0),
	}
	
	// Register default evaluators
	engine.RegisterEvaluator(&NewestEvaluator{})
	engine.RegisterEvaluator(&OldestEvaluator{})
	engine.RegisterEvaluator(&ShortestPathEvaluator{})
	engine.RegisterEvaluator(&LongestPathEvaluator{})
	engine.RegisterEvaluator(&SpecificDirEvaluator{})
	engine.RegisterEvaluator(&LargestEvaluator{})
	engine.RegisterEvaluator(&SmallestEvaluator{})
	
	return engine
}

// RegisterEvaluator registers an evaluator for a criteria
func (e *RuleEngine) RegisterEvaluator(eval RuleEvaluator) {
	// Get criteria from evaluator type
	var criteria RuleCriteria
	switch eval.(type) {
	case *NewestEvaluator:
		criteria = CriteriaKeepNewest
	case *OldestEvaluator:
		criteria = CriteriaKeepOldest
	case *ShortestPathEvaluator:
		criteria = CriteriaKeepShortestPath
	case *LongestPathEvaluator:
		criteria = CriteriaKeepLongestPath
	case *SpecificDirEvaluator:
		criteria = CriteriaKeepSpecificDir
	case *LargestEvaluator:
		criteria = CriteriaKeepLargest
	case *SmallestEvaluator:
		criteria = CriteriaKeepSmallest
	}
	
	e.evaluators[criteria] = eval
}

// ApplyRule applies a rule to a single group
func (e *RuleEngine) ApplyRule(group *models.DuplicateGroup, rule *Rule) error {
	evaluator, ok := e.evaluators[rule.Criteria]
	if !ok {
		return fmt.Errorf("no evaluator for criteria: %s", rule.Criteria)
	}
	
	if !evaluator.IsApplicable(group) {
		return fmt.Errorf("rule not applicable to this group")
	}
	
	keepIndex, err := evaluator.Evaluate(group)
	if err != nil {
		return err
	}
	
	// Mark files for deletion
	group.MarkForDeletion([]int{keepIndex})
	group.RuleApplied = &models.Rule{
		ID:       rule.ID,
		Name:     rule.Name,
		Criteria: models.RuleCriteria(rule.Criteria),
	}
	
	// Record in history
	e.history = append(e.history, RuleApplication{
		GroupID:   group.ID,
		RuleID:    rule.ID,
		KeepIndex: keepIndex,
	})
	
	return nil
}

// ApplyToAll applies a rule to all groups
func (e *RuleEngine) ApplyToAll(groups []*models.DuplicateGroup, rule *Rule) (applied int, skipped int, err error) {
	evaluator, ok := e.evaluators[rule.Criteria]
	if !ok {
		return 0, 0, fmt.Errorf("no evaluator for criteria: %s", rule.Criteria)
	}

	for _, group := range groups {
		if !evaluator.IsApplicable(group) {
			skipped++
			continue
		}

		keepIndex, evalErr := evaluator.Evaluate(group)
		if evalErr != nil {
			skipped++
			continue
		}

		group.MarkForDeletion([]int{keepIndex})
		group.RuleApplied = &models.Rule{
			ID:       rule.ID,
			Name:     rule.Name,
			Criteria: models.RuleCriteria(rule.Criteria),
		}

		e.history = append(e.history, RuleApplication{
			GroupID:   group.ID,
			RuleID:    rule.ID,
			KeepIndex: keepIndex,
		})

		applied++
	}

	return applied, skipped, nil
}

// GetHistory returns the history of rule applications
func (e *RuleEngine) GetHistory() []RuleApplication {
	return e.history
}

// UndoLast undoes the last rule application
func (e *RuleEngine) UndoLast(groups []models.DuplicateGroup) bool {
	if len(e.history) == 0 {
		return false
	}
	
	last := e.history[len(e.history)-1]
	e.history = e.history[:len(e.history)-1]
	
	// Find and clear the group
	for i := range groups {
		if groups[i].ID == last.GroupID {
			groups[i].ClearDecisions()
			return true
		}
	}
	
	return false
}

// ClearHistory clears the application history
func (e *RuleEngine) ClearHistory() {
	e.history = make([]RuleApplication, 0)
}

// NewestEvaluator keeps the most recently modified file
type NewestEvaluator struct{}

func (e *NewestEvaluator) Evaluate(group *models.DuplicateGroup) (int, error) {
	if len(group.Files) == 0 {
		return 0, fmt.Errorf("no files in group")
	}
	
	bestIdx := 0
	bestTime := group.Files[0].ModTime
	
	for i, file := range group.Files {
		if file.ModTime.After(bestTime) {
			bestTime = file.ModTime
			bestIdx = i
		}
	}
	
	return bestIdx, nil
}

func (e *NewestEvaluator) Description() string {
	return "Keeps the most recently modified file"
}

func (e *NewestEvaluator) IsApplicable(group *models.DuplicateGroup) bool {
	return len(group.Files) >= 2
}

// OldestEvaluator keeps the oldest file
type OldestEvaluator struct{}

func (e *OldestEvaluator) Evaluate(group *models.DuplicateGroup) (int, error) {
	if len(group.Files) == 0 {
		return 0, fmt.Errorf("no files in group")
	}
	
	bestIdx := 0
	bestTime := group.Files[0].ModTime
	
	for i, file := range group.Files {
		if file.ModTime.Before(bestTime) {
			bestTime = file.ModTime
			bestIdx = i
		}
	}
	
	return bestIdx, nil
}

func (e *OldestEvaluator) Description() string {
	return "Keeps the oldest file"
}

func (e *OldestEvaluator) IsApplicable(group *models.DuplicateGroup) bool {
	return len(group.Files) >= 2
}

// ShortestPathEvaluator keeps the file with the shortest path
type ShortestPathEvaluator struct{}

func (e *ShortestPathEvaluator) Evaluate(group *models.DuplicateGroup) (int, error) {
	if len(group.Files) == 0 {
		return 0, fmt.Errorf("no files in group")
	}
	
	bestIdx := 0
	bestDepth := getPathDepth(group.Files[0].Path)
	
	for i, file := range group.Files {
		depth := getPathDepth(file.Path)
		if depth < bestDepth {
			bestDepth = depth
			bestIdx = i
		}
	}
	
	return bestIdx, nil
}

func (e *ShortestPathEvaluator) Description() string {
	return "Keeps the file with the shortest path"
}

func (e *ShortestPathEvaluator) IsApplicable(group *models.DuplicateGroup) bool {
	return len(group.Files) >= 2
}

// LongestPathEvaluator keeps the file with the longest path
type LongestPathEvaluator struct{}

func (e *LongestPathEvaluator) Evaluate(group *models.DuplicateGroup) (int, error) {
	if len(group.Files) == 0 {
		return 0, fmt.Errorf("no files in group")
	}
	
	bestIdx := 0
	bestDepth := getPathDepth(group.Files[0].Path)
	
	for i, file := range group.Files {
		depth := getPathDepth(file.Path)
		if depth > bestDepth {
			bestDepth = depth
			bestIdx = i
		}
	}
	
	return bestIdx, nil
}

func (e *LongestPathEvaluator) Description() string {
	return "Keeps the file with the longest path"
}

func (e *LongestPathEvaluator) IsApplicable(group *models.DuplicateGroup) bool {
	return len(group.Files) >= 2
}

// SpecificDirEvaluator keeps files from a specific directory
type SpecificDirEvaluator struct {
	TargetDir string
}

func (e *SpecificDirEvaluator) Evaluate(group *models.DuplicateGroup) (int, error) {
	if len(group.Files) == 0 {
		return 0, fmt.Errorf("no files in group")
	}
	
	if e.TargetDir == "" {
		return 0, fmt.Errorf("no target directory specified")
	}
	
	// First, try exact match
	for i, file := range group.Files {
		if strings.HasPrefix(file.Path, e.TargetDir) {
			return i, nil
		}
	}
	
	// Fallback to first file
	return 0, nil
}

func (e *SpecificDirEvaluator) Description() string {
	return fmt.Sprintf("Keeps files from: %s", e.TargetDir)
}

func (e *SpecificDirEvaluator) IsApplicable(group *models.DuplicateGroup) bool {
	return len(group.Files) >= 2 && e.TargetDir != ""
}

// SetDirectory sets the target directory
func (e *SpecificDirEvaluator) SetDirectory(dir string) {
	e.TargetDir = dir
}

// LargestEvaluator keeps the largest file
type LargestEvaluator struct{}

func (e *LargestEvaluator) Evaluate(group *models.DuplicateGroup) (int, error) {
	if len(group.Files) == 0 {
		return 0, fmt.Errorf("no files in group")
	}
	
	bestIdx := 0
	bestSize := group.Files[0].Size
	
	for i, file := range group.Files {
		if file.Size > bestSize {
			bestSize = file.Size
			bestIdx = i
		}
	}
	
	return bestIdx, nil
}

func (e *LargestEvaluator) Description() string {
	return "Keeps the largest file"
}

func (e *LargestEvaluator) IsApplicable(group *models.DuplicateGroup) bool {
	return len(group.Files) >= 2
}

// SmallestEvaluator keeps the smallest file
type SmallestEvaluator struct{}

func (e *SmallestEvaluator) Evaluate(group *models.DuplicateGroup) (int, error) {
	if len(group.Files) == 0 {
		return 0, fmt.Errorf("no files in group")
	}
	
	bestIdx := 0
	bestSize := group.Files[0].Size
	
	for i, file := range group.Files {
		if file.Size < bestSize {
			bestSize = file.Size
			bestIdx = i
		}
	}
	
	return bestIdx, nil
}

func (e *SmallestEvaluator) Description() string {
	return "Keeps the smallest file"
}

func (e *SmallestEvaluator) IsApplicable(group *models.DuplicateGroup) bool {
	return len(group.Files) >= 2
}

// getPathDepth returns the depth of a path
func getPathDepth(path string) int {
	return strings.Count(path, "/") + strings.Count(path, "\\")
}

// sortFilesByModTime sorts files by modification time
func sortFilesByModTime(files []models.FileEntry, ascending bool) {
	sort.Slice(files, func(i, j int) bool {
		if ascending {
			return files[i].ModTime.Before(files[j].ModTime)
		}
		return files[i].ModTime.After(files[j].ModTime)
	})
}
