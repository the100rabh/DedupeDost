package script

// Template provides script templates
type Template struct{}

// BashScriptTemplate is the default bash script template
const BashScriptTemplate = `#!/bin/bash
# ============================================
# DupDel Cleanup Script
# ============================================
# Generated: {{.GeneratedTime}}
# Source Directory: {{.SourceDirectory}}
# Total files to delete: {{.TotalFiles}}
# Total space to recover: {{.TotalSize}}
# ============================================

set -e  # Exit on error

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

echo "Starting cleanup..."
echo "DRY_RUN mode: $DRY_RUN"
echo ""

{{range .Groups}}
# Group {{.ID}}: {{.Name}} ({{.Size}})
{{range .DeletePaths}}delete_file "{{.}}"
{{end}}
# Keeping: {{.KeepPath}}

{{end}}
echo ""
echo "Cleanup complete!"
echo "Files deleted: {{.TotalFiles}}"
echo "Space recovered: {{.TotalSize}}"
`

// DryRunTemplate adds dry-run capability
const DryRunTemplate = `
# To run in dry-run mode (preview only):
#   DRY_RUN=true ./{{.ScriptName}}
#
# To run actual deletion:
#   ./{{.ScriptName}}
`

// SafeDeleteTemplate includes additional safety checks
const SafeDeleteTemplate = `
# Safety check - confirm before deleting
if [ "$DRY_RUN" != "true" ]; then
    read -p "This will permanently delete {{.TotalFiles}} files. Continue? (y/N) " confirm
    if [ "$confirm" != "y" ] && [ "$confirm" != "Y" ]; then
        echo "Aborted."
        exit 1
    fi
fi
`

// BackupTemplate creates backups before deletion
const BackupTemplate = `
# Create backup directory
BACKUP_DIR="{{.BackupDir}}"
mkdir -p "$BACKUP_DIR"

# Backup function
backup_and_delete() {
    local file="$1"
    local basename=$(basename "$file")
    if [ -f "$file" ]; then
        cp "$file" "$BACKUP_DIR/$basename"
        rm -f "$file"
        echo "Backed up and deleted: $file"
    fi
}
`
