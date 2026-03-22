package detector_test

import (
	"testing"
	"time"

	"github.com/dupdel/dup-del/internal/detector"
	"github.com/dupdel/dup-del/pkg/models"
)

// TestNewStatsCalculator tests NewStatsCalculator
func TestNewStatsCalculator(t *testing.T) {
	calc := detector.NewStatsCalculator()
	if calc == nil {
		t.Fatal("NewStatsCalculator returned nil")
	}
}

// TestStatsCalculator_Calculate tests Calculate
func TestStatsCalculator_Calculate(t *testing.T) {
	calc := detector.NewStatsCalculator()

	groups := []*models.DuplicateGroup{
		{
			ID:   "g1",
			Hash: "abc123",
			Size: 1024,
			Files: []models.FileEntry{
				{Name: "file1.txt", Size: 1024},
				{Name: "file2.txt", Size: 1024},
			},
		},
		{
			ID:   "g2",
			Hash: "def456",
			Size: 2048,
			Files: []models.FileEntry{
				{Name: "file3.txt", Size: 2048},
				{Name: "file4.txt", Size: 2048},
				{Name: "file5.txt", Size: 2048},
			},
		},
	}

	stats := calc.Calculate(groups)

	// Total files: 2 + 3 = 5
	if stats.TotalFiles != 5 {
		t.Errorf("TotalFiles = %d, want 5", stats.TotalFiles)
	}

	// Total size: 2*1024 + 3*2048 = 2048 + 6144 = 8192
	if stats.TotalSize != 8192 {
		t.Errorf("TotalSize = %d, want 8192", stats.TotalSize)
	}

	// Duplicate files: 1 + 2 = 3 (one kept per group)
	if stats.DuplicateFiles != 3 {
		t.Errorf("DuplicateFiles = %d, want 3", stats.DuplicateFiles)
	}

	// Duplicate groups: 2
	if stats.DuplicateGroups != 2 {
		t.Errorf("DuplicateGroups = %d, want 2", stats.DuplicateGroups)
	}

	// Recoverable size: 1*1024 + 2*2048 = 1024 + 4096 = 5120
	if stats.RecoverableSize != 5120 {
		t.Errorf("RecoverableSize = %d, want 5120", stats.RecoverableSize)
	}
}

// TestStatsCalculator_Calculate_Empty tests Calculate with empty groups
func TestStatsCalculator_Calculate_Empty(t *testing.T) {
	calc := detector.NewStatsCalculator()

	stats := calc.Calculate([]*models.DuplicateGroup{})

	if stats.TotalFiles != 0 {
		t.Errorf("TotalFiles = %d, want 0", stats.TotalFiles)
	}
	if stats.TotalSize != 0 {
		t.Errorf("TotalSize = %d, want 0", stats.TotalSize)
	}
	if stats.DuplicateFiles != 0 {
		t.Errorf("DuplicateFiles = %d, want 0", stats.DuplicateFiles)
	}
	if stats.DuplicateGroups != 0 {
		t.Errorf("DuplicateGroups = %d, want 0", stats.DuplicateGroups)
	}
	if stats.RecoverableSize != 0 {
		t.Errorf("RecoverableSize = %d, want 0", stats.RecoverableSize)
	}
}

// TestStatsCalculator_CalculateFromSession tests CalculateFromSession
func TestStatsCalculator_CalculateFromSession(t *testing.T) {
	calc := detector.NewStatsCalculator()

	session := models.NewScanSession("/test")
	session.DuplicateGroups = []*models.DuplicateGroup{
		{
			ID:   "g1",
			Hash: "abc123",
			Size: 1024,
			Files: []models.FileEntry{
				{Name: "file1.txt", Size: 1024},
				{Name: "file2.txt", Size: 1024},
			},
		},
	}

	stats := calc.CalculateFromSession(session)

	if stats.DuplicateGroups != 1 {
		t.Errorf("DuplicateGroups = %d, want 1", stats.DuplicateGroups)
	}
}

// TestStatsCalculator_Summary tests Summary
func TestStatsCalculator_Summary(t *testing.T) {
	calc := detector.NewStatsCalculator()

	stats := detector.DetectionStats{
		TotalFiles:      100,
		TotalSize:       1024 * 1024,
		DuplicateFiles:  10,
		DuplicateGroups: 5,
		RecoverableSize: 512 * 1024,
		ScanDuration:    time.Second,
		HashDuration:    500 * time.Millisecond,
	}

	summary := calc.Summary(stats)
	if summary == "" {
		t.Error("Summary should not be empty")
	}
}

// TestGetEfficiencyRatio tests GetEfficiencyRatio
func TestGetEfficiencyRatio(t *testing.T) {
	tests := []struct {
		name     string
		stats    detector.DetectionStats
		expected float64
	}{
		{
			name: "zero_total",
			stats: detector.DetectionStats{
				TotalFiles:     0,
				DuplicateFiles: 0,
			},
			expected: 0,
		},
		{
			name: "no_duplicates",
			stats: detector.DetectionStats{
				TotalFiles:     100,
				DuplicateFiles: 0,
			},
			expected: 0,
		},
		{
			name: "half_duplicates",
			stats: detector.DetectionStats{
				TotalFiles:     100,
				DuplicateFiles: 50,
			},
			expected: 50.0,
		},
		{
			name: "all_duplicates",
			stats: detector.DetectionStats{
				TotalFiles:     100,
				DuplicateFiles: 100,
			},
			expected: 100.0,
		},
		{
			name: "partial",
			stats: detector.DetectionStats{
				TotalFiles:     200,
				DuplicateFiles: 25,
			},
			expected: 12.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.GetEfficiencyRatio(tt.stats)
			if result != tt.expected {
				t.Errorf("GetEfficiencyRatio = %f, want %f", result, tt.expected)
			}
		})
	}
}

// TestGetSpaceEfficiency tests GetSpaceEfficiency
func TestGetSpaceEfficiency(t *testing.T) {
	tests := []struct {
		name     string
		stats    detector.DetectionStats
		expected float64
	}{
		{
			name: "zero_total",
			stats: detector.DetectionStats{
				TotalSize:       0,
				RecoverableSize: 0,
			},
			expected: 0,
		},
		{
			name: "no_recoverable",
			stats: detector.DetectionStats{
				TotalSize:       1024,
				RecoverableSize: 0,
			},
			expected: 0,
		},
		{
			name: "half_recoverable",
			stats: detector.DetectionStats{
				TotalSize:       1024,
				RecoverableSize: 512,
			},
			expected: 50.0,
		},
		{
			name: "all_recoverable",
			stats: detector.DetectionStats{
				TotalSize:       1024,
				RecoverableSize: 1024,
			},
			expected: 100.0,
		},
		{
			name: "partial",
			stats: detector.DetectionStats{
				TotalSize:       1000,
				RecoverableSize: 250,
			},
			expected: 25.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.GetSpaceEfficiency(tt.stats)
			if result != tt.expected {
				t.Errorf("GetSpaceEfficiency = %f, want %f", result, tt.expected)
			}
		})
	}
}

// TestDetectionStats tests DetectionStats struct
func TestDetectionStats(t *testing.T) {
	stats := detector.DetectionStats{
		TotalFiles:      100,
		TotalSize:       1024,
		DuplicateFiles:  10,
		DuplicateGroups: 5,
		RecoverableSize: 512,
		ScanDuration:    time.Second,
		HashDuration:    500 * time.Millisecond,
	}

	if stats.TotalFiles != 100 {
		t.Errorf("TotalFiles = %d, want 100", stats.TotalFiles)
	}
	if stats.TotalSize != 1024 {
		t.Errorf("TotalSize = %d, want 1024", stats.TotalSize)
	}
	if stats.DuplicateFiles != 10 {
		t.Errorf("DuplicateFiles = %d, want 10", stats.DuplicateFiles)
	}
	if stats.DuplicateGroups != 5 {
		t.Errorf("DuplicateGroups = %d, want 5", stats.DuplicateGroups)
	}
	if stats.RecoverableSize != 512 {
		t.Errorf("RecoverableSize = %d, want 512", stats.RecoverableSize)
	}
}
