package rules

// RuleCriteria defines the criteria for selecting files
type RuleCriteria string

const (
	CriteriaKeepNewest       RuleCriteria = "keep_newest"
	CriteriaKeepOldest       RuleCriteria = "keep_oldest"
	CriteriaKeepShortestPath RuleCriteria = "keep_shortest_path"
	CriteriaKeepLongestPath  RuleCriteria = "keep_longest_path"
	CriteriaKeepSpecificDir  RuleCriteria = "keep_specific_dir"
	CriteriaKeepLargest      RuleCriteria = "keep_largest"
	CriteriaKeepSmallest     RuleCriteria = "keep_smallest"
	CriteriaManualSelection  RuleCriteria = "manual"
)

// BuiltinRules contains all built-in rule definitions
var BuiltinRules = []Rule{
	{
		ID:          "keep_newest",
		Name:        "Keep Newest",
		Description: "Keep the most recently modified file and delete older copies",
		Criteria:    CriteriaKeepNewest,
		Icon:        "clock_recent",
	},
	{
		ID:          "keep_oldest",
		Name:        "Keep Oldest",
		Description: "Keep the oldest file and delete newer copies",
		Criteria:    CriteriaKeepOldest,
		Icon:        "clock_old",
	},
	{
		ID:          "keep_shortest_path",
		Name:        "Keep Shortest Path",
		Description: "Keep the file with the shortest path (root level preferred)",
		Criteria:    CriteriaKeepShortestPath,
		Icon:        "path_short",
	},
	{
		ID:          "keep_longest_path",
		Name:        "Keep Longest Path",
		Description: "Keep the file with the longest path (most specific location)",
		Criteria:    CriteriaKeepLongestPath,
		Icon:        "path_long",
	},
	{
		ID:          "keep_specific_dir",
		Name:        "Keep from Specific Directory",
		Description: "Keep files from a specific directory, delete others",
		Criteria:    CriteriaKeepSpecificDir,
		Icon:        "folder_select",
		ConfigSchema: map[string]interface{}{
			"directory": "",
		},
	},
	{
		ID:          "keep_largest",
		Name:        "Keep Largest",
		Description: "Keep the largest file (useful for different quality versions)",
		Criteria:    CriteriaKeepLargest,
		Icon:        "size_large",
	},
	{
		ID:          "keep_smallest",
		Name:        "Keep Smallest",
		Description: "Keep the smallest file (for space efficiency)",
		Criteria:    CriteriaKeepSmallest,
		Icon:        "size_small",
	},
	{
		ID:          "manual",
		Name:        "Manual Selection",
		Description: "Manually select which files to keep for each group",
		Criteria:    CriteriaManualSelection,
		Icon:        "hand_select",
	},
}

// Rule represents a rule for automatic file selection
type Rule struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Criteria     RuleCriteria           `json:"criteria"`
	Icon         string                 `json:"icon"`
	ConfigSchema map[string]interface{} `json:"config_schema,omitempty"`
	Config       map[string]interface{} `json:"config,omitempty"`
	ApplyToAll   bool                   `json:"apply_to_all"`
}

// NewRule creates a new rule
func NewRule(id, name, description string, criteria RuleCriteria) *Rule {
	return &Rule{
		ID:          id,
		Name:        name,
		Description: description,
		Criteria:    criteria,
		Config:      make(map[string]interface{}),
	}
}

// WithConfig adds configuration to the rule
func (r *Rule) WithConfig(key string, value interface{}) *Rule {
	if r.Config == nil {
		r.Config = make(map[string]interface{})
	}
	r.Config[key] = value
	return r
}

// WithApplyToAll sets whether the rule applies to all remaining groups
func (r *Rule) WithApplyToAll(apply bool) *Rule {
	r.ApplyToAll = apply
	return r
}

// GetRuleByID returns a rule by its ID
func GetRuleByID(id string) *Rule {
	for _, rule := range BuiltinRules {
		if rule.ID == id {
			return &rule
		}
	}
	return nil
}

// GetRulesByCriteria returns rules matching the given criteria
func GetRulesByCriteria(criteria RuleCriteria) []Rule {
	var rules []Rule
	for _, rule := range BuiltinRules {
		if rule.Criteria == criteria {
			rules = append(rules, rule)
		}
	}
	return rules
}

// GetAllRules returns all available rules
func GetAllRules() []Rule {
	return BuiltinRules
}
