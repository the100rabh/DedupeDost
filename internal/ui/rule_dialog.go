package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/dupdel/dup-del/internal/rules"
	"github.com/dupdel/dup-del/pkg/models"
)

// RuleDialog displays the rule selection dialog
type RuleDialog struct {
	dialog        *dialog.CustomDialog
	ruleList      *widget.List
	rulePreview   *widget.RichText
	applyToAll    *widget.Check
	configBox     *fyne.Container
	applyButton   *widget.Button
	cancelButton  *widget.Button
	selectedRule  *rules.Rule
	rules         []rules.Rule
	onSelect      func(rule *rules.Rule, applyToAll bool)
}

// NewRuleDialog creates a new rule dialog
func NewRuleDialog(callback func(*rules.Rule, bool)) *RuleDialog {
	rd := &RuleDialog{
		rules:    rules.GetAllRules(),
		onSelect: callback,
	}
	
	// Rule list
	rd.ruleList = widget.NewList(
		func() int { return len(rd.rules) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.FolderIcon()),
				container.NewVBox(
					widget.NewLabelWithStyle("Rule Name", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
					widget.NewLabel("Description"),
				),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id < 0 || id >= len(rd.rules) {
				return
			}
			
			rule := rd.rules[id]
			container := item.(*fyne.Container)
			icon := container.Objects[0].(*widget.Icon)
			textContainer := container.Objects[1].(*fyne.Container)
			nameLabel := textContainer.Objects[0].(*widget.Label)
			descLabel := textContainer.Objects[1].(*widget.Label)
			
			// Set icon based on criteria
			icon.SetResource(getRuleIcon(rule.Criteria))
			
			nameLabel.SetText(rule.Name)
			descLabel.SetText(rule.Description)
		},
	)
	
	// Rule preview
	rd.rulePreview = widget.NewRichTextFromMarkdown("Select a rule to see preview")
	
	// Apply to all check
	rd.applyToAll = widget.NewCheck("Apply to all remaining groups", nil)
	rd.applyToAll.SetChecked(true)
	
	// Config box (hidden by default)
	rd.configBox = container.NewVBox()
	rd.configBox.Hide()
	
	// Apply button
	rd.applyButton = widget.NewButtonWithIcon("Apply Rule", theme.ConfirmIcon(), func() {
		if rd.selectedRule != nil && rd.onSelect != nil {
			rd.onSelect(rd.selectedRule, rd.applyToAll.Checked)
			rd.dialog.Hide()
		}
	})
	rd.applyButton.Disable()
	
	// Cancel button
	rd.cancelButton = widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
		rd.dialog.Hide()
	})
	
	// Selection handler
	rd.ruleList.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(rd.rules) {
			rd.selectedRule = &rd.rules[id]
			rd.applyButton.Enable()
			rd.updatePreview()
		}
	}
	
	// Create content
	content := container.NewBorder(
		widget.NewLabel("Select a rule:"),
		container.NewHBox(
			rd.applyToAll,
			layout.NewSpacer(),
			rd.cancelButton,
			rd.applyButton,
		),
		nil,
		nil,
		container.NewHSplit(
			rd.ruleList,
			container.NewVBox(
				widget.NewLabel("Preview:"),
				rd.rulePreview,
				rd.configBox,
			),
		),
	)
	
	// Create dialog
	rd.dialog = dialog.NewCustom("Select Rule", "Cancel", content, nil)
	rd.dialog.SetOnClosed(func() {
		rd.selectedRule = nil
	})
	
	return rd
}

// Show shows the rule dialog
func (rd *RuleDialog) Show() {
	rd.dialog.Show()
}

// updatePreview updates the rule preview
func (rd *RuleDialog) updatePreview() {
	if rd.selectedRule == nil {
		return
	}
	
	markdown := fmt.Sprintf(`## %s

%s

**Criteria:** %s

**How it works:**
This rule will automatically select which file to keep based on the criteria above.
`, rd.selectedRule.Name, rd.selectedRule.Description, rd.selectedRule.Criteria)
	
	rd.rulePreview = widget.NewRichTextFromMarkdown(markdown)
}

// getRuleIcon returns an icon for the rule criteria
func getRuleIcon(criteria rules.RuleCriteria) fyne.Resource {
	switch criteria {
	case rules.CriteriaKeepNewest:
		return theme.HistoryIcon()
	case rules.CriteriaKeepOldest:
		return theme.HistoryIcon()
	case rules.CriteriaKeepShortestPath:
		return theme.FolderIcon()
	case rules.CriteriaKeepSpecificDir:
		return theme.FolderOpenIcon()
	default:
		return theme.FolderIcon()
	}
}

// RuleSelectionPanel displays rule selection inline
type RuleSelectionPanel struct {
	container   *fyne.Container
	ruleSelect  *widget.Select
	applyToAll  *widget.Check
	applyButton *widget.Button
	onApply     func(rule *rules.Rule, applyToAll bool)
}

// NewRuleSelectionPanel creates a new rule selection panel
func NewRuleSelectionPanel(onApply func(*rules.Rule, bool)) *RuleSelectionPanel {
	rsp := &RuleSelectionPanel{
		onApply: onApply,
	}
	
	// Build rule options
	ruleOptions := make([]string, len(rules.GetAllRules()))
	for i, rule := range rules.GetAllRules() {
		ruleOptions[i] = rule.Name
	}
	
	// Rule select
	rsp.ruleSelect = widget.NewSelect(ruleOptions, func(value string) {
		// Find selected rule
		for _, rule := range rules.GetAllRules() {
			if rule.Name == value {
				break
			}
		}
	})
	rsp.ruleSelect.SetSelected("")

	// Apply to all check
	rsp.applyToAll = widget.NewCheck("Apply to all remaining groups", nil)
	rsp.applyToAll.SetChecked(true)
	
	// Apply button
	rsp.applyButton = widget.NewButtonWithIcon("Apply", theme.ConfirmIcon(), func() {
		if rsp.ruleSelect.Selected != "" && rsp.onApply != nil {
			// Find selected rule
			for _, rule := range rules.GetAllRules() {
				if rule.Name == rsp.ruleSelect.Selected {
					rsp.onApply(&rule, rsp.applyToAll.Checked)
					break
				}
			}
		}
	})
	rsp.applyButton.Disable()
	
	rsp.ruleSelect.OnChanged = func(s string) {
		if s != "" {
			rsp.applyButton.Enable()
		} else {
			rsp.applyButton.Disable()
		}
	}
	
	// Create container
	rsp.container = container.NewHBox(
		widget.NewLabel("Apply Rule:"),
		rsp.ruleSelect,
		rsp.applyToAll,
		layout.NewSpacer(),
		rsp.applyButton,
	)
	
	return rsp
}

// Container returns the panel container
func (rsp *RuleSelectionPanel) Container() *fyne.Container {
	return rsp.container
}

// Reset resets the panel state
func (rsp *RuleSelectionPanel) Reset() {
	rsp.ruleSelect.SetSelected("")
	rsp.applyToAll.SetChecked(true)
	rsp.applyButton.Disable()
}

// ReviewPanel displays decisions for review
type ReviewPanel struct {
	container   *fyne.Container
	list        *widget.List
	filterSelect *widget.Select
	decisions   []models.GroupDecisions
	statsLabel  *widget.Label
}

// NewReviewPanel creates a new review panel
func NewReviewPanel() *ReviewPanel {
	rp := &ReviewPanel{
		decisions: make([]models.GroupDecisions, 0),
	}
	
	// Filter select
	rp.filterSelect = widget.NewSelect([]string{
		"All",
		"Marked for Deletion",
		"Keeping All",
		"Skipped",
	}, func(value string) {
		rp.list.Refresh()
	})
	rp.filterSelect.SetSelected("All")
	
	// Stats label
	rp.statsLabel = widget.NewLabel("")
	
	// List
	rp.list = widget.NewList(
		func() int { return len(rp.decisions) },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id < 0 || id >= len(rp.decisions) {
				return
			}
			
			label := item.(*widget.Label)
			dec := rp.decisions[id]
			
			status := "⏭️ Skipped"
			if len(dec.DeletePaths) > 0 {
				status = fmt.Sprintf("🗑️ %d files to delete", len(dec.DeletePaths))
			} else if len(dec.KeepIndices) > 0 {
				status = "✅ Keeping all"
			}
			
			label.SetText(fmt.Sprintf("%s - %s", dec.GroupID, status))
		},
	)
	
	// Create container
	rp.container = container.NewBorder(
		container.NewHBox(
			widget.NewLabel("Filter:"),
			rp.filterSelect,
			layout.NewSpacer(),
			rp.statsLabel,
		),
		nil,
		nil,
		nil,
		rp.list,
	)
	
	return rp
}

// SetDecisions sets the decisions to display
func (rp *ReviewPanel) SetDecisions(decisions []models.GroupDecisions) {
	rp.decisions = decisions
	rp.updateStats()
	rp.list.Refresh()
}

// updateStats updates the statistics display
func (rp *ReviewPanel) updateStats() {
	total := len(rp.decisions)
	toDelete := 0
	keeping := 0
	
	for _, dec := range rp.decisions {
		if len(dec.DeletePaths) > 0 {
			toDelete += len(dec.DeletePaths)
		} else {
			keeping++
		}
	}
	
	rp.statsLabel.SetText(fmt.Sprintf("Total: %d | To delete: %d | Keeping: %d", 
		total, toDelete, keeping))
}

// Container returns the panel container
func (rp *ReviewPanel) Container() *fyne.Container {
	return rp.container
}
