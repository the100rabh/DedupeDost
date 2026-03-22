package models

// RuleCriteria defines the criteria for selecting files to keep/delete
type RuleCriteria string

const (
	KeepNewest       RuleCriteria = "newest"
	KeepOldest       RuleCriteria = "oldest"
	KeepShortestPath RuleCriteria = "shortest_path"
	KeepLongestPath  RuleCriteria = "longest_path"
	KeepSpecificDir  RuleCriteria = "specific_dir"
	ManualSelection  RuleCriteria = "manual"
)

// Rule represents a rule for automatically selecting files to keep/delete
type Rule struct {
	ID          string                 // Unique rule identifier
	Name        string                 // Human-readable name
	Description string                 // Rule description
	Criteria    RuleCriteria           // Selection criteria
	ApplyToAll  bool                   // Whether to apply to all remaining groups
	Config      map[string]interface{} // Additional configuration
}

// FileDecision represents a decision for a single file
type FileDecision struct {
	FilePath string // File path
	Decision string // "keep" or "delete"
}

// GroupDecisions represents decisions for a duplicate group
type GroupDecisions struct {
	GroupID     string       // Group identifier
	KeepIndices []int        // Indices of files to keep
	DeletePaths []string     // Paths to delete
	ApplyToAll  bool         // Whether to apply to all remaining groups
	Rule        *Rule        // Rule applied (if any)
}

// NewRule creates a new rule with the given parameters
func NewRule(id, name, description string, criteria RuleCriteria) *Rule {
	return &Rule{
		ID:          id,
		Name:        name,
		Description: description,
		Criteria:    criteria,
		ApplyToAll:  false,
		Config:      make(map[string]interface{}),
	}
}

// WithConfig adds configuration to the rule
func (r *Rule) WithConfig(key string, value interface{}) *Rule {
	r.Config[key] = value
	return r
}

// WithApplyToAll sets whether the rule applies to all remaining groups
func (r *Rule) WithApplyToAll(apply bool) *Rule {
	r.ApplyToAll = apply
	return r
}

// GetConfig retrieves a configuration value
func (r *Rule) GetConfig(key string) (interface{}, bool) {
	value, ok := r.Config[key]
	return value, ok
}
