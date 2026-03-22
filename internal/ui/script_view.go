package ui

import (
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/dupdel/dup-del/internal/script"
	"github.com/dupdel/dup-del/pkg/models"
)

// ScriptView represents the script preview view
type ScriptView struct {
	container       *fyne.Container
	scriptScroll    *container.Scroll
	scriptText      *widget.RichText
	pathEntry       *widget.Entry
	browseButton    *widget.Button
	executableCheck *widget.Check
	backupCheck     *widget.Check
	dryRunCheck     *widget.Check
	saveButton      *widget.Button
	copyButton      *widget.Button
	generator       *script.Generator
	session         *models.ScanSession
	dupDelApp       *DupDelApp
}

// NewScriptView creates a new script view
func NewScriptView(dupDelApp *DupDelApp) *ScriptView {
	sv := &ScriptView{
		dupDelApp: dupDelApp,
	}
	
	// Script text (read-only)
	sv.scriptText = widget.NewRichTextFromMarkdown("No script generated yet.\n\nRun a scan first, then generate the cleanup script.")
	
	// Path entry
	cwd, _ := os.Getwd()
	defaultPath := filepath.Join(cwd, "cleanup_duplicates.sh")
	sv.pathEntry = widget.NewEntry()
	sv.pathEntry.SetText(defaultPath)
	sv.pathEntry.Disable() // Read-only display

	// Browse button
	sv.browseButton = widget.NewButtonWithIcon("Browse...", theme.FolderOpenIcon(), func() {
		sv.showSaveDialog()
	})

	// Executable check
	sv.executableCheck = widget.NewCheck("Make script executable", nil)
	sv.executableCheck.SetChecked(true)

	// Backup check
	sv.backupCheck = widget.NewCheck("Create backup of existing script", nil)

	// Dry run check
	sv.dryRunCheck = widget.NewCheck("Include dry-run support (DRY_RUN=true)", nil)
	sv.dryRunCheck.SetChecked(true)

	// Save button
	sv.saveButton = widget.NewButtonWithIcon("Save Script", theme.DocumentSaveIcon(), func() {
		sv.save()
	})
	sv.saveButton.Importance = widget.HighImportance
	sv.saveButton.Disable()

	// Copy button
	sv.copyButton = widget.NewButtonWithIcon("Copy to Clipboard", theme.ContentCopyIcon(), func() {
		sv.copyToClipboard()
	})
	sv.copyButton.Disable()

	// Path box - entry on top, label+button below (similar to DirectoryBar)
	pathBox := container.NewVBox(
		sv.pathEntry,
		container.NewHBox(
			widget.NewLabel("Output:"),
			layout.NewSpacer(),
			sv.browseButton,
		),
		container.NewHBox(
			sv.executableCheck,
			sv.backupCheck,
			sv.dryRunCheck,
		),
	)
	
	// Buttons box
	buttonsBox := container.NewHBox(
		sv.copyButton,
		layout.NewSpacer(),
		sv.saveButton,
	)

	// Script scroll - expands to fill available space
	sv.scriptScroll = container.NewScroll(sv.scriptText)

	// Create container with Border layout for proper expansion
	sv.container = container.NewBorder(
		pathBox,        // Top
		buttonsBox,     // Bottom
		nil, nil,       // No left/right
		sv.scriptScroll, // Center (expands)
	)

	return sv
}

// Container returns the container
func (sv *ScriptView) Container() *fyne.Container {
	return sv.container
}

// SetScriptContent sets the script content for display
func (sv *ScriptView) SetScriptContent(content string) {
	sv.scriptText = widget.NewRichTextFromMarkdown("```bash\n" + content + "\n```")
	if sv.scriptScroll != nil {
		sv.scriptScroll.Content = sv.scriptText
		sv.scriptScroll.Refresh()
	}
}

// generate generates the cleanup script
func (sv *ScriptView) generate() {
	if sv.dupDelApp == nil || sv.dupDelApp.mainView == nil {
		return
	}
	
	// Get groups from results
	groups := sv.dupDelApp.mainView.resultsView.groups
	if len(groups) == 0 {
		dialog.ShowError(fmt.Errorf("no duplicate groups found"), 
			fyne.CurrentApp().Driver().AllWindows()[0])
		return
	}
	
	// Create session
	sv.session = models.NewScanSession(sv.dupDelApp.GetScanDirectory())
	sv.session.DuplicateGroups = groups
	
	// Mark files for deletion (keep first in each group by default)
	for i := range sv.session.DuplicateGroups {
		sv.session.DuplicateGroups[i].MarkForDeletion([]int{0})
	}
	sv.session.CalculateStatistics()
	
	// Create generator
	outputPath := sv.pathEntry.Text
	sv.generator = script.NewGenerator(outputPath)
	sv.generator.SetDryRunSupport(sv.dryRunCheck.Checked)
	
	// Generate script content (don't save yet)
	content := sv.generatorForPreview(sv.session)

	// Display script
	sv.scriptText = widget.NewRichTextFromMarkdown("```bash\n" + content + "\n```")
	if sv.scriptScroll != nil {
		sv.scriptScroll.Content = sv.scriptText
		sv.scriptScroll.Refresh()
	}

	// Enable buttons
	sv.saveButton.Enable()
	sv.copyButton.Enable()
}

// generatorForPreview generates script content for preview
func (sv *ScriptView) generatorForPreview(session *models.ScanSession) string {
	var content string
	
	// Header
	content += "#!/bin/bash\n\n"
	content += "# ============================================\n"
	content += "# DupDel Cleanup Script\n"
	content += "# ============================================\n"
	content += fmt.Sprintf("# Generated: %s\n", session.StartTime.Format("2006-01-02 15:04:05"))
	content += fmt.Sprintf("# Source Directory: %s\n", session.SourceDir)
	content += fmt.Sprintf("# Total files to delete: %d\n", session.FilesToDelete)
	content += fmt.Sprintf("# Total space to recover: %s\n", session.GetRecoverableSizeHuman())
	content += "# ============================================\n\n"
	
	// DRY_RUN support
	if sv.dryRunCheck.Checked {
		content += `DRY_RUN=${DRY_RUN:-false}

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
	
	// Groups
	for i, group := range session.DuplicateGroups {
		content += fmt.Sprintf("# Group %d: %s (%s)\n", 
			i+1, group.Files[0].Name, formatSize(group.Size))
		
		for _, path := range group.DeletePaths {
			if sv.dryRunCheck.Checked {
				content += fmt.Sprintf("delete_file \"%s\"\n", path)
			} else {
				content += fmt.Sprintf("rm -f \"%s\"\n", path)
			}
		}
		content += fmt.Sprintf("# Keeping: %s\n\n", group.Files[0].Path)
	}
	
	// Summary
	content += "# ============================================\n"
	content += "# Cleanup Complete!\n"
	content += "# ============================================\n"
	
	return content
}

// save saves the script to file
func (sv *ScriptView) save() {
	outputPath := sv.pathEntry.Text
	
	// Create generator
	sv.generator = script.NewGenerator(outputPath)
	sv.generator.SetDryRunSupport(sv.dryRunCheck.Checked)
	
	// Generate script
	scriptPath, err := sv.generator.Generate(sv.session)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to save script: %w", err),
			fyne.CurrentApp().Driver().AllWindows()[0])
		return
	}
	
	// Make executable if requested
	if sv.executableCheck.Checked {
		sv.generator.MakeExecutable(scriptPath)
	}
	
	// Show success
	dialog.ShowInformation("Success", 
		fmt.Sprintf("Script saved to:\n%s\n\nTo preview: DRY_RUN=true %s\nTo execute: %s",
			scriptPath, scriptPath, scriptPath),
		fyne.CurrentApp().Driver().AllWindows()[0])
}

// copyToClipboard copies the script to clipboard
func (sv *ScriptView) copyToClipboard() {
	// Get clipboard
	clipboard := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()
	
	// Get script content
	content := sv.generatorForPreview(sv.session)
	
	// Copy
	clipboard.SetContent(content)
	
	// Show confirmation
	dialog.ShowInformation("Copied", "Script copied to clipboard", 
		fyne.CurrentApp().Driver().AllWindows()[0])
}

// showSaveDialog shows the save file dialog
func (sv *ScriptView) showSaveDialog() {
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			return
		}
		if writer == nil {
			return
		}
		
		path := writer.URI().Path()
		sv.pathEntry.SetText(path)
		writer.Close()
	}, fyne.CurrentApp().Driver().AllWindows()[0])
}
