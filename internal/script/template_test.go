package script_test

import (
	"strings"
	"testing"

	"github.com/the100rabh/DedupeDost/internal/script"
)

// TestTemplate_Constants tests template constants
func TestTemplate_Constants(t *testing.T) {
	if script.BashScriptTemplate == "" {
		t.Error("BashScriptTemplate should not be empty")
	}

	if script.DryRunTemplate == "" {
		t.Error("DryRunTemplate should not be empty")
	}

	if script.SafeDeleteTemplate == "" {
		t.Error("SafeDeleteTemplate should not be empty")
	}

	if script.BackupTemplate == "" {
		t.Error("BackupTemplate should not be empty")
	}
}

// TestBashScriptTemplate_Content tests BashScriptTemplate content
func TestBashScriptTemplate_Content(t *testing.T) {
	template := script.BashScriptTemplate

	// Check for required elements
	required := []string{
		"#!/bin/bash",
		"DedupeDost Cleanup Script",
		"{{.GeneratedTime}}",
		"{{.SourceDirectory}}",
		"{{.TotalFiles}}",
		"{{.TotalSize}}",
		"DRY_RUN",
		"delete_file()",
		"{{range .Groups}}",
		"{{range .DeletePaths}}",
	}

	for _, req := range required {
		if !strings.Contains(template, req) {
			t.Errorf("BashScriptTemplate missing required element: %s", req)
		}
	}
}

// TestDryRunTemplate_Content tests DryRunTemplate content
func TestDryRunTemplate_Content(t *testing.T) {
	template := script.DryRunTemplate

	required := []string{
		"DRY_RUN=true",
		"{{.ScriptName}}",
	}

	for _, req := range required {
		if !strings.Contains(template, req) {
			t.Errorf("DryRunTemplate missing required element: %s", req)
		}
	}
}

// TestSafeDeleteTemplate_Content tests SafeDeleteTemplate content
func TestSafeDeleteTemplate_Content(t *testing.T) {
	template := script.SafeDeleteTemplate

	required := []string{
		"confirm",
		"{{.TotalFiles}}",
		"Continue",
	}

	for _, req := range required {
		if !strings.Contains(template, req) {
			t.Errorf("SafeDeleteTemplate missing required element: %s", req)
		}
	}
}

// TestBackupTemplate_Content tests BackupTemplate content
func TestBackupTemplate_Content(t *testing.T) {
	template := script.BackupTemplate

	required := []string{
		"BACKUP_DIR",
		"{{.BackupDir}}",
		"backup_and_delete",
		"cp",
		"rm",
	}

	for _, req := range required {
		if !strings.Contains(template, req) {
			t.Errorf("BackupTemplate missing required element: %s", req)
		}
	}
}

// TestBashScriptTemplate_Structure tests template structure
func TestBashScriptTemplate_Structure(t *testing.T) {
	template := script.BashScriptTemplate

	// Check structure
	if !strings.HasPrefix(template, "#!/bin/bash") {
		t.Error("BashScriptTemplate should start with shebang")
	}

	// Check for error handling
	if !strings.Contains(template, "set -e") {
		t.Error("BashScriptTemplate should include error handling")
	}

	// Check for file existence check
	if !strings.Contains(template, "[ -f \"$file\" ]") {
		t.Error("BashScriptTemplate should check file existence")
	}

	// Check for warning on missing file
	if !strings.Contains(template, "Warning: File not found") {
		t.Error("BashScriptTemplate should warn on missing file")
	}
}

// TestTemplate_NoEmptyLines tests that templates don't have excessive empty lines
func TestTemplate_NoEmptyLines(t *testing.T) {
	templates := []string{
		script.BashScriptTemplate,
		script.DryRunTemplate,
		script.SafeDeleteTemplate,
		script.BackupTemplate,
	}

	for i, template := range templates {
		lines := strings.Split(template, "\n")
		consecutiveEmpty := 0
		maxConsecutive := 0

		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				consecutiveEmpty++
				if consecutiveEmpty > maxConsecutive {
					maxConsecutive = consecutiveEmpty
				}
			} else {
				consecutiveEmpty = 0
			}
		}

		if maxConsecutive > 3 {
			t.Errorf("Template %d has %d consecutive empty lines", i, maxConsecutive)
		}
	}
}

// TestBashScriptTemplate_DryRunSupport tests DRY_RUN support in template
func TestBashScriptTemplate_DryRunSupport(t *testing.T) {
	template := script.BashScriptTemplate

	// Check for DRY_RUN variable
	if !strings.Contains(template, "DRY_RUN=${DRY_RUN:-false}") {
		t.Error("BashScriptTemplate should define DRY_RUN variable")
	}

	// Check for dry run mode check
	if !strings.Contains(template, `"$DRY_RUN" = "true"`) {
		t.Error("BashScriptTemplate should check DRY_RUN mode")
	}

	// Check for dry run message
	if !strings.Contains(template, "[DRY RUN]") {
		t.Error("BashScriptTemplate should include dry run message")
	}
}

// TestBashScriptTemplate_GroupIteration tests group iteration in template
func TestBashScriptTemplate_GroupIteration(t *testing.T) {
	template := script.BashScriptTemplate

	// Check for group iteration
	if !strings.Contains(template, "{{range .Groups}}") {
		t.Error("BashScriptTemplate should iterate over groups")
	}

	// Check for group properties
	groupProps := []string{
		"{{.ID}}",
		"{{.Name}}",
		"{{.Size}}",
		"{{.KeepPath}}",
	}

	for _, prop := range groupProps {
		if !strings.Contains(template, prop) {
			t.Errorf("BashScriptTemplate missing group property: %s", prop)
		}
	}
}

// TestBashScriptTemplate_DeletePaths tests delete paths iteration
func TestBashScriptTemplate_DeletePaths(t *testing.T) {
	template := script.BashScriptTemplate

	// Check for delete paths iteration
	if !strings.Contains(template, "{{range .DeletePaths}}") {
		t.Error("BashScriptTemplate should iterate over delete paths")
	}

	// Check for delete_file call
	if !strings.Contains(template, "delete_file") {
		t.Error("BashScriptTemplate should call delete_file")
	}
}

// TestBashScriptTemplate_Summary tests summary section
func TestBashScriptTemplate_Summary(t *testing.T) {
	template := script.BashScriptTemplate

	// Check for summary
	if !strings.Contains(template, "Cleanup complete") {
		t.Error("BashScriptTemplate should include completion message")
	}

	if !strings.Contains(template, "Files deleted") {
		t.Error("BashScriptTemplate should report files deleted")
	}

	if !strings.Contains(template, "Space recovered") {
		t.Error("BashScriptTemplate should report space recovered")
	}
}

// TestBackupTemplate_BackupDirectory tests backup directory creation
func TestBackupTemplate_BackupDirectory(t *testing.T) {
	template := script.BackupTemplate

	if !strings.Contains(template, "mkdir -p") {
		t.Error("BackupTemplate should create backup directory")
	}

	if !strings.Contains(template, "$BACKUP_DIR") {
		t.Error("BackupTemplate should use BACKUP_DIR variable")
	}
}

// TestSafeDeleteTemplate_Confirmation tests confirmation prompt
func TestSafeDeleteTemplate_Confirmation(t *testing.T) {
	template := script.SafeDeleteTemplate

	if !strings.Contains(template, "read -p") {
		t.Error("SafeDeleteTemplate should prompt for confirmation")
	}

	if !strings.Contains(template, "Continue?") {
		t.Error("SafeDeleteTemplate should ask for confirmation")
	}

	if !strings.Contains(template, "Aborted") {
		t.Error("SafeDeleteTemplate should handle abort")
	}
}

// TestTemplate_Variables tests template variables
func TestTemplate_Variables(t *testing.T) {
	// Collect all template variables
	allTemplates := script.BashScriptTemplate + script.DryRunTemplate +
		script.SafeDeleteTemplate + script.BackupTemplate

	expectedVars := []string{
		"{{.GeneratedTime}}",
		"{{.SourceDirectory}}",
		"{{.TotalFiles}}",
		"{{.TotalSize}}",
		"{{.ScriptName}}",
		"{{.BackupDir}}",
	}

	for _, v := range expectedVars {
		if !strings.Contains(allTemplates, v) {
			t.Errorf("Missing template variable: %s", v)
		}
	}
}

// TestTemplate_GoTemplateSyntax tests Go template syntax
func TestTemplate_GoTemplateSyntax(t *testing.T) {
	templates := map[string]string{
		"BashScriptTemplate": script.BashScriptTemplate,
		"DryRunTemplate":     script.DryRunTemplate,
		"SafeDeleteTemplate": script.SafeDeleteTemplate,
		"BackupTemplate":     script.BackupTemplate,
	}

	for name, template := range templates {
		// Check for balanced {{ }}
		openCount := strings.Count(template, "{{")
		closeCount := strings.Count(template, "}}")
		if openCount != closeCount {
			t.Errorf("%s: unbalanced {{ }} (%d open, %d close)",
				name, openCount, closeCount)
		}

		// Check for balanced range/end
		rangeCount := strings.Count(template, "{{range")
		endCount := strings.Count(template, "{{end}}")
		if rangeCount != endCount {
			t.Errorf("%s: unbalanced range/end (%d range, %d end)",
				name, rangeCount, endCount)
		}
	}
}
