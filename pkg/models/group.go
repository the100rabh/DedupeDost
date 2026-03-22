package models

import "github.com/dupdel/dup-del/pkg/types"

// DuplicateGroup represents a group of duplicate files
type DuplicateGroup struct {
	ID          string         // Unique group identifier
	Hash        string         // Common hash value
	Size        int64          // File size (all files have same size)
	Extension   string         // File extension
	FileType    types.FileType // Type of files in group
	Files       []FileEntry    // All duplicate files
	KeepIndices []int          // Indices of files to keep
	DeletePaths []string       // Paths marked for deletion
	RuleApplied *Rule          // Rule applied to this group (if any)
}

// GetRecoverableSize returns the total size that can be recovered by deleting duplicates
func (g *DuplicateGroup) GetRecoverableSize() int64 {
	if len(g.Files) <= 1 {
		return 0
	}
	// Recoverable = (total files - 1) * file size
	return g.Size * int64(len(g.Files)-1)
}

// GetFileCount returns the number of files in this group
func (g *DuplicateGroup) GetFileCount() int {
	return len(g.Files)
}

// GetDuplicateCount returns the number of duplicate files (excluding the original)
func (g *DuplicateGroup) GetDuplicateCount() int {
	if len(g.Files) <= 1 {
		return 0
	}
	return len(g.Files) - 1
}

// GetTotalSize returns the total size of all files in the group
func (g *DuplicateGroup) GetTotalSize() int64 {
	return g.Size * int64(len(g.Files))
}

// GetDisplaySize returns human-readable size of one file in the group
func (g *DuplicateGroup) GetDisplaySize() string {
	return FormatSize(g.Size)
}

// GetDisplayRecoverableSize returns human-readable recoverable size
func (g *DuplicateGroup) GetDisplayRecoverableSize() string {
	return FormatSize(g.GetRecoverableSize())
}

// MarkForDeletion marks specific files for deletion based on indices to keep
func (g *DuplicateGroup) MarkForDeletion(keepIndices []int) {
	g.KeepIndices = keepIndices
	g.DeletePaths = make([]string, 0)

	keepMap := make(map[int]bool)
	for _, idx := range keepIndices {
		keepMap[idx] = true
	}

	for i, file := range g.Files {
		if !keepMap[i] {
			g.DeletePaths = append(g.DeletePaths, file.Path)
		}
	}
}

// ClearDecisions clears all keep/delete decisions
func (g *DuplicateGroup) ClearDecisions() {
	g.KeepIndices = nil
	g.DeletePaths = nil
	g.RuleApplied = nil
}
