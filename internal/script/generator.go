package script

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dupdel/dup-del/pkg/models"
)

// Generator generates bash cleanup scripts
type Generator struct {
	outputDir       string
	scriptName      string
	includeHeader   bool
	includeComments bool
	dryRunSupport   bool
}

// NewGenerator creates a new script Generator
func NewGenerator(outputPath string) *Generator {
	dir := filepath.Dir(outputPath)
	name := filepath.Base(outputPath)
	
	return &Generator{
		outputDir:       dir,
		scriptName:      name,
		includeHeader:   true,
		includeComments: true,
		dryRunSupport:   true,
	}
}

// Generate generates a cleanup script from a scan session
func (g *Generator) Generate(session *models.ScanSession) (string, error) {
	scriptPath := filepath.Join(g.outputDir, g.scriptName)
	
	// Ensure output directory exists
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}
	
	// Generate script content
	content := g.generateContent(session)
	
	// Write to file
	if err := os.WriteFile(scriptPath, []byte(content), 0755); err != nil {
		return "", fmt.Errorf("failed to write script: %w", err)
	}
	
	// Update session with script path
	session.GeneratedScript = scriptPath
	
	return scriptPath, nil
}

// generateContent generates the script content
func (g *Generator) generateContent(session *models.ScanSession) string {
	var content string
	
	// Shebang
	content += "#!/bin/bash\n\n"
	
	// Header
	if g.includeHeader {
		content += g.generateHeader(session)
	}
	
	// DRY_RUN support
	if g.dryRunSupport {
		content += `
# DRY_RUN mode: set to "true" to preview without deleting
DRY_RUN=${DRY_RUN:-false}

# delete_file function - safely deletes a file
delete_file() {
    local file="$1"
    if [ "$DRY_RUN" = "true" ]; then
        echo "[DRY RUN] Would delete: $file"
    else
        if [ -f "$file" ]; then
            echo "Deleting: $file"
            rm -f "$file"
        else
            echo "Warning: File not found: $file"
        fi
    fi
}

`
	}
	
	// Group deletions
	for i, group := range session.DuplicateGroups {
		if g.includeComments {
			content += fmt.Sprintf("# Group %d: %s (%s)\n", 
				i+1, group.Files[0].Name, models.FormatSize(group.Size))
		}
		
		for _, path := range group.DeletePaths {
			if g.dryRunSupport {
				content += fmt.Sprintf("delete_file \"%s\"\n", path)
			} else {
				content += fmt.Sprintf("rm -f \"%s\"\n", path)
			}
		}
		
		if g.includeComments {
			content += fmt.Sprintf("# Keeping: %s\n", group.Files[0].Path)
		}
		content += "\n"
	}
	
	// Summary
	content += g.generateSummary(session)
	
	return content
}

// generateHeader generates the script header
func (g *Generator) generateHeader(session *models.ScanSession) string {
	var header string
	
	header += "# ============================================\n"
	header += "# DupDel Cleanup Script\n"
	header += "# ============================================\n"
	header += fmt.Sprintf("# Generated: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	header += fmt.Sprintf("# Source Directory: %s\n", session.SourceDir)
	header += fmt.Sprintf("# Total files to delete: %d\n", session.FilesToDelete)
	header += fmt.Sprintf("# Total space to recover: %s\n", session.GetRecoverableSizeHuman())
	header += "# ============================================\n\n"
	
	return header
}

// generateSummary generates the script summary
func (g *Generator) generateSummary(session *models.ScanSession) string {
	var summary string
	
	summary += "# ============================================\n"
	summary += "# Cleanup Summary\n"
	summary += "# ============================================\n"
	summary += fmt.Sprintf("# Duplicate groups processed: %d\n", session.GetGroupCount())
	summary += fmt.Sprintf("# Files marked for deletion: %d\n", session.FilesToDelete)
	summary += fmt.Sprintf("# Space to be recovered: %s\n", session.GetRecoverableSizeHuman())
	summary += "#\n"
	summary += "# To execute this script:\n"
	summary += "#   1. Review the files to be deleted\n"
	summary += "#   2. Run: DRY_RUN=true ./cleanup.sh (preview)\n"
	summary += "#   3. Run: ./cleanup.sh (actual deletion)\n"
	summary += "# ============================================\n"
	
	return summary
}

// MakeExecutable ensures the script is executable
func (g *Generator) MakeExecutable(path string) error {
	return os.Chmod(path, 0755)
}

// Validate validates the generated script
func (g *Generator) Validate(path string) error {
	// Check file exists
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("script not found: %w", err)
	}
	
	// Check it's a file
	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a file")
	}
	
	// Check it's executable
	if info.Mode()&0111 == 0 {
		return fmt.Errorf("script is not executable")
	}
	
	return nil
}

// SetIncludeHeader sets whether to include header comments
func (g *Generator) SetIncludeHeader(include bool) {
	g.includeHeader = include
}

// SetIncludeComments sets whether to include group comments
func (g *Generator) SetIncludeComments(include bool) {
	g.includeComments = include
}

// SetDryRunSupport sets whether to include DRY_RUN support
func (g *Generator) SetDryRunSupport(include bool) {
	g.dryRunSupport = include
}
