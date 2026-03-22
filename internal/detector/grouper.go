package detector

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"time"

	"github.com/dupdel/dup-del/pkg/models"
)

// Grouper provides utilities for grouping files
type Grouper struct {
	groupCounter int
}

// NewGrouper creates a new Grouper
func NewGrouper() *Grouper {
	return &Grouper{}
}

// GroupBySize groups files by their size
func (g *Grouper) GroupBySize(files []models.FileEntry) map[int64][]models.FileEntry {
	sizeMap := make(map[int64][]models.FileEntry)
	
	for _, file := range files {
		sizeMap[file.Size] = append(sizeMap[file.Size], file)
	}
	
	return sizeMap
}

// GroupByHash groups files by their hash
func (g *Grouper) GroupByHash(files []models.FileEntry) map[string][]models.FileEntry {
	hashMap := make(map[string][]models.FileEntry)
	
	for _, file := range files {
		if file.Hash != "" {
			hashMap[file.Hash] = append(hashMap[file.Hash], file)
		}
	}
	
	return hashMap
}

// FilterDuplicates filters to only groups with 2+ files
func (g *Grouper) FilterDuplicates(groupMap map[int64][]models.FileEntry) map[int64][]models.FileEntry {
	result := make(map[int64][]models.FileEntry)
	
	for size, files := range groupMap {
		if len(files) >= 2 {
			result[size] = files
		}
	}
	
	return result
}

// CreateDuplicateGroups creates DuplicateGroup objects from hash groups
func (g *Grouper) CreateDuplicateGroups(hashMap map[string][]models.FileEntry) []models.DuplicateGroup {
	var groups []models.DuplicateGroup
	
	for hash, files := range hashMap {
		if len(files) < 2 {
			continue
		}
		
		group := models.DuplicateGroup{
			ID:        g.generateGroupID(),
			Hash:      hash,
			Size:      files[0].Size,
			Extension: files[0].Extension,
			FileType:  files[0].FileType,
			Files:     files,
		}
		
		groups = append(groups, group)
	}
	
	return groups
}

// generateGroupID generates a unique group ID
func (g *Grouper) generateGroupID() string {
	g.groupCounter++
	return fmt.Sprintf("grp_%d", g.groupCounter)
}

// SortGroups sorts duplicate groups by various criteria
func SortGroups(groups []models.DuplicateGroup, by SortBy) {
	switch by {
	case SortBySize:
		sort.Slice(groups, func(i, j int) bool {
			return groups[i].Size > groups[j].Size
		})
	case SortByCount:
		sort.Slice(groups, func(i, j int) bool {
			return len(groups[i].Files) > len(groups[j].Files)
		})
	case SortByRecoverable:
		sort.Slice(groups, func(i, j int) bool {
			return groups[i].GetRecoverableSize() > groups[j].GetRecoverableSize()
		})
	case SortByName:
		sort.Slice(groups, func(i, j int) bool {
			return groups[i].Extension < groups[j].Extension
		})
	}
}

// SortBy defines sorting criteria for groups
type SortBy int

const (
	SortBySize SortBy = iota
	SortByCount
	SortByRecoverable
	SortByName
)

// CalculateHash calculates SHA-256 hash of a string (for testing)
func CalculateHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// FilesHaveSameContent compares two files by their metadata
func FilesHaveSameContent(f1, f2 models.FileEntry) bool {
	return f1.Size == f2.Size && f1.Hash == f2.Hash
}

// FindPotentialDuplicates finds files that might be duplicates based on size
func FindPotentialDuplicates(files []models.FileEntry) []models.FileEntry {
	sizeCount := make(map[int64]int)
	
	// Count files by size
	for _, file := range files {
		sizeCount[file.Size]++
	}
	
	// Filter to files with matching sizes
	var potential []models.FileEntry
	for _, file := range files {
		if sizeCount[file.Size] >= 2 {
			potential = append(potential, file)
		}
	}
	
	return potential
}

// EstimateDuplicates estimates the number of duplicates without full hashing
func EstimateDuplicates(files []models.FileEntry) int {
	sizeMap := make(map[int64]int)
	
	for _, file := range files {
		sizeMap[file.Size]++
	}
	
	estimate := 0
	for _, count := range sizeMap {
		if count >= 2 {
			estimate += count - 1
		}
	}
	
	return estimate
}

// GetFileAge returns the age of a file
func GetFileAge(entry models.FileEntry) time.Duration {
	return time.Since(entry.ModTime)
}

// GetPathDepth returns the depth of a file path
func GetPathDepth(path string) int {
	depth := 0
	for _, c := range path {
		if c == '/' || c == '\\' {
			depth++
		}
	}
	return depth
}
