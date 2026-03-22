package models_test

import (
	"testing"

	"github.com/the100rabh/DedupeDost/pkg/models"
)

// TestRuleCriteria constants
func TestRuleCriteria_Constants(t *testing.T) {
	if models.KeepNewest != "newest" {
		t.Errorf("KeepNewest = %s, want newest", models.KeepNewest)
	}
	if models.KeepOldest != "oldest" {
		t.Errorf("KeepOldest = %s, want oldest", models.KeepOldest)
	}
	if models.KeepShortestPath != "shortest_path" {
		t.Errorf("KeepShortestPath = %s, want shortest_path", models.KeepShortestPath)
	}
	if models.KeepLongestPath != "longest_path" {
		t.Errorf("KeepLongestPath = %s, want longest_path", models.KeepLongestPath)
	}
	if models.KeepSpecificDir != "specific_dir" {
		t.Errorf("KeepSpecificDir = %s, want specific_dir", models.KeepSpecificDir)
	}
	if models.ManualSelection != "manual" {
		t.Errorf("ManualSelection = %s, want manual", models.ManualSelection)
	}
}

// TestNewRule tests creating a new rule
func TestNewRule(t *testing.T) {
	rule := models.NewRule("rule1", "Keep Newest", "Keep the most recently modified file", models.KeepNewest)

	if rule == nil {
		t.Fatal("NewRule returned nil")
	}

	if rule.ID != "rule1" {
		t.Errorf("ID = %s, want rule1", rule.ID)
	}

	if rule.Name != "Keep Newest" {
		t.Errorf("Name = %s, want Keep Newest", rule.Name)
	}

	if rule.Description != "Keep the most recently modified file" {
		t.Errorf("Description = %s, want Keep the most recently modified file", rule.Description)
	}

	if rule.Criteria != models.KeepNewest {
		t.Errorf("Criteria = %s, want newest", rule.Criteria)
	}

	if rule.ApplyToAll {
		t.Error("ApplyToAll should be false by default")
	}

	if rule.Config == nil {
		t.Error("Config should be initialized")
	}

	if len(rule.Config) != 0 {
		t.Errorf("Config should be empty initially, got %d items", len(rule.Config))
	}
}

// TestRule_WithConfig tests WithConfig method
func TestRule_WithConfig(t *testing.T) {
	rule := models.NewRule("rule1", "Test Rule", "Test", models.KeepNewest)

	// Add config
	result := rule.WithConfig("key1", "value1")

	// Should return same rule for chaining
	if result != rule {
		t.Error("WithConfig should return same rule for chaining")
	}

	// Verify config was added
	value, ok := rule.GetConfig("key1")
	if !ok {
		t.Error("Config key1 not found")
	}
	if value != "value1" {
		t.Errorf("Config key1 = %v, want value1", value)
	}
}

// TestRule_WithConfig_MultipleValues tests adding multiple config values
func TestRule_WithConfig_MultipleValues(t *testing.T) {
	rule := models.NewRule("rule1", "Test Rule", "Test", models.KeepNewest)

	rule.WithConfig("string", "value").
		WithConfig("int", 42).
		WithConfig("bool", true).
		WithConfig("float", 3.14)

	tests := []struct {
		key      string
		expected interface{}
	}{
		{"string", "value"},
		{"int", 42},
		{"bool", true},
		{"float", 3.14},
	}

	for _, tt := range tests {
		value, ok := rule.GetConfig(tt.key)
		if !ok {
			t.Errorf("Config key %s not found", tt.key)
			continue
		}
		if value != tt.expected {
			t.Errorf("Config %s = %v, want %v", tt.key, value, tt.expected)
		}
	}
}

// TestRule_WithApplyToAll tests WithApplyToAll method
func TestRule_WithApplyToAll(t *testing.T) {
	rule := models.NewRule("rule1", "Test Rule", "Test", models.KeepNewest)

	// Default should be false
	if rule.ApplyToAll {
		t.Error("ApplyToAll should be false by default")
	}

	// Set to true
	result := rule.WithApplyToAll(true)
	if result != rule {
		t.Error("WithApplyToAll should return same rule for chaining")
	}
	if !rule.ApplyToAll {
		t.Error("ApplyToAll should be true after WithApplyToAll(true)")
	}

	// Set to false
	rule.WithApplyToAll(false)
	if rule.ApplyToAll {
		t.Error("ApplyToAll should be false after WithApplyToAll(false)")
	}
}

// TestRule_WithApplyToAll_Chaining tests method chaining
func TestRule_WithApplyToAll_Chaining(t *testing.T) {
	rule := models.NewRule("rule1", "Test Rule", "Test", models.KeepNewest)

	// Chain multiple calls
	result := rule.WithApplyToAll(true).WithConfig("key", "value")

	if result != rule {
		t.Error("Chaining should return same rule")
	}

	if !rule.ApplyToAll {
		t.Error("ApplyToAll should be true")
	}

	_, ok := rule.GetConfig("key")
	if !ok {
		t.Error("Config should be set")
	}
}

// TestRule_GetConfig tests GetConfig method
func TestRule_GetConfig(t *testing.T) {
	rule := models.NewRule("rule1", "Test Rule", "Test", models.KeepNewest)

	// Non-existent key should return false
	_, ok := rule.GetConfig("nonexistent")
	if ok {
		t.Error("GetConfig should return false for non-existent key")
	}

	// Add a value and retrieve it
	rule.WithConfig("existing", "value")
	value, ok := rule.GetConfig("existing")
	if !ok {
		t.Error("GetConfig should return true for existing key")
	}
	if value != "value" {
		t.Errorf("GetConfig = %v, want value", value)
	}
}

// TestRule_GetConfig_DifferentTypes tests config with different types
func TestRule_GetConfig_DifferentTypes(t *testing.T) {
	rule := models.NewRule("rule1", "Test Rule", "Test", models.KeepNewest)

	rule.WithConfig("string", "hello").
		WithConfig("number", 123).
		WithConfig("slice", []string{"a", "b", "c"}).
		WithConfig("map", map[string]int{"x": 1})

	tests := []struct {
		key      string
		expected interface{}
	}{
		{"string", "hello"},
		{"number", 123},
		{"slice", []string{"a", "b", "c"}},
		{"map", map[string]int{"x": 1}},
	}

	for _, tt := range tests {
		value, ok := rule.GetConfig(tt.key)
		if !ok {
			t.Errorf("Config key %s not found", tt.key)
			continue
		}
		if value == nil {
			t.Errorf("Config %s is nil", tt.key)
		}
	}
}

// TestFileDecision tests FileDecision struct
func TestFileDecision(t *testing.T) {
	decision := models.FileDecision{
		FilePath: "/path/to/file.txt",
		Decision: "keep",
	}

	if decision.FilePath != "/path/to/file.txt" {
		t.Errorf("FilePath = %s, want /path/to/file.txt", decision.FilePath)
	}

	if decision.Decision != "keep" {
		t.Errorf("Decision = %s, want keep", decision.Decision)
	}
}

// TestGroupDecisions tests GroupDecisions struct
func TestGroupDecisions(t *testing.T) {
	rule := models.NewRule("rule1", "Test", "Test", models.KeepNewest)

	decisions := models.GroupDecisions{
		GroupID:     "group1",
		KeepIndices: []int{0, 2},
		DeletePaths: []string{"/path/file2.txt", "/path/file4.txt"},
		ApplyToAll:  true,
		Rule:        rule,
	}

	if decisions.GroupID != "group1" {
		t.Errorf("GroupID = %s, want group1", decisions.GroupID)
	}

	if len(decisions.KeepIndices) != 2 {
		t.Errorf("KeepIndices length = %d, want 2", len(decisions.KeepIndices))
	}

	if len(decisions.DeletePaths) != 2 {
		t.Errorf("DeletePaths length = %d, want 2", len(decisions.DeletePaths))
	}

	if !decisions.ApplyToAll {
		t.Error("ApplyToAll should be true")
	}

	if decisions.Rule != rule {
		t.Error("Rule should match")
	}
}

// TestGroupDecisions_Empty tests empty GroupDecisions
func TestGroupDecisions_Empty(t *testing.T) {
	decisions := models.GroupDecisions{
		GroupID:     "group1",
		KeepIndices: []int{},
		DeletePaths: []string{},
		ApplyToAll:  false,
		Rule:        nil,
	}

	if decisions.GroupID != "group1" {
		t.Errorf("GroupID = %s, want group1", decisions.GroupID)
	}

	if decisions.ApplyToAll {
		t.Error("ApplyToAll should be false")
	}

	if decisions.Rule != nil {
		t.Error("Rule should be nil")
	}
}
