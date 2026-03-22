package detector

import (
	"fmt"

	"github.com/dupdel/dup-del/pkg/models"
)

// StatsCalculator calculates detection statistics
type StatsCalculator struct {
	stats *DetectionStats
}

// NewStatsCalculator creates a new StatsCalculator
func NewStatsCalculator() *StatsCalculator {
	return &StatsCalculator{
		stats: &DetectionStats{},
	}
}

// Calculate calculates all statistics from the given groups
func (c *StatsCalculator) Calculate(groups []*models.DuplicateGroup) DetectionStats {
	totalFiles := 0
	totalSize := int64(0)
	duplicateFiles := 0
	duplicateGroups := len(groups)
	recoverableSize := int64(0)

	for _, group := range groups {
		fileCount := group.GetFileCount()
		totalFiles += fileCount
		totalSize += group.GetTotalSize()
		duplicateFiles += group.GetDuplicateCount()
		recoverableSize += group.GetRecoverableSize()
	}

	return DetectionStats{
		TotalFiles:      totalFiles,
		TotalSize:       totalSize,
		DuplicateFiles:  duplicateFiles,
		DuplicateGroups: duplicateGroups,
		RecoverableSize: recoverableSize,
	}
}

// CalculateFromSession calculates statistics from a scan session
func (c *StatsCalculator) CalculateFromSession(session *models.ScanSession) DetectionStats {
	return c.Calculate(session.DuplicateGroups)
}

// Summary generates a text summary of the statistics
func (c *StatsCalculator) Summary(stats DetectionStats) string {
	return formatStats(stats)
}

// formatStats formats statistics for display
func formatStats(stats DetectionStats) string {
	lines := []string{
		"=== Detection Statistics ===",
		fmt.Sprintf("Total files scanned: %d", stats.TotalFiles),
		fmt.Sprintf("Total size: %s", models.FormatSize(stats.TotalSize)),
		fmt.Sprintf("Duplicate groups: %d", stats.DuplicateGroups),
		fmt.Sprintf("Duplicate files: %d", stats.DuplicateFiles),
		fmt.Sprintf("Recoverable space: %s", models.FormatSize(stats.RecoverableSize)),
		fmt.Sprintf("Scan duration: %v", stats.ScanDuration),
		fmt.Sprintf("Hash duration: %v", stats.HashDuration),
	}
	
	var result string
	for _, line := range lines {
		result += line + "\n"
	}
	
	return result
}

// GetEfficiencyRatio returns the ratio of duplicates to total files
func GetEfficiencyRatio(stats DetectionStats) float64 {
	if stats.TotalFiles == 0 {
		return 0
	}
	return float64(stats.DuplicateFiles) / float64(stats.TotalFiles) * 100
}

// GetSpaceEfficiency returns the percentage of space that can be recovered
func GetSpaceEfficiency(stats DetectionStats) float64 {
	if stats.TotalSize == 0 {
		return 0
	}
	return float64(stats.RecoverableSize) / float64(stats.TotalSize) * 100
}
