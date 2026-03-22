package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ScanSession represents a complete scan operation
type ScanSession struct {
	ID               string             // Unique session identifier
	SourceDir        string             // Scanned directory
	StartTime        time.Time          // When scan started
	EndTime          time.Time          // When scan ended
	TotalFiles       int                // Total files scanned
	TotalSize        int64              // Total size in bytes
	DuplicateGroups  []*DuplicateGroup  // Found duplicate groups
	FilesToDelete    int                // Count of files marked for deletion
	SpaceRecoverable int64              // Recoverable space in bytes
	GeneratedScript  string             // Path to generated script
	FilterOptions    map[string]string  // Applied filter options
}

// NewScanSession creates a new scan session for the given directory
func NewScanSession(sourceDir string) *ScanSession {
	return &ScanSession{
		ID:              uuid.New().String(),
		SourceDir:       sourceDir,
		StartTime:       time.Now(),
		DuplicateGroups: make([]*DuplicateGroup, 0),
		FilterOptions:   make(map[string]string),
	}
}

// GetDuration returns the duration of the scan session
func (s *ScanSession) GetDuration() time.Duration {
	endTime := s.EndTime
	if endTime.IsZero() {
		endTime = time.Now()
	}
	return endTime.Sub(s.StartTime)
}

// GetRecoverableSizeHuman returns human-readable recoverable size
func (s *ScanSession) GetRecoverableSizeHuman() string {
	return FormatSize(s.SpaceRecoverable)
}

// GetTotalSizeHuman returns human-readable total size
func (s *ScanSession) GetTotalSizeHuman() string {
	return FormatSize(s.TotalSize)
}

// GetDurationString returns formatted duration string
func (s *ScanSession) GetDurationString() string {
	duration := s.GetDuration()
	
	if duration < time.Minute {
		return fmt.Sprintf("%.1f seconds", duration.Seconds())
	}
	if duration < time.Hour {
		minutes := int(duration.Minutes())
		seconds := int(duration.Seconds()) % 60
		return fmt.Sprintf("%d minutes %d seconds", minutes, seconds)
	}
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	return fmt.Sprintf("%d hours %d minutes", hours, minutes)
}

// AddGroup adds a duplicate group to the session
func (s *ScanSession) AddGroup(group *DuplicateGroup) {
	s.DuplicateGroups = append(s.DuplicateGroups, group)
}

// CalculateStatistics calculates session statistics
func (s *ScanSession) CalculateStatistics() {
	totalFiles := 0
	totalToDelete := 0
	spaceRecoverable := int64(0)

	for _, group := range s.DuplicateGroups {
		totalFiles += group.GetFileCount()
		totalToDelete += len(group.DeletePaths)
		spaceRecoverable += group.GetRecoverableSize()
	}

	s.FilesToDelete = totalToDelete
	s.SpaceRecoverable = spaceRecoverable
}

// Finalize marks the session as complete
func (s *ScanSession) Finalize() {
	s.EndTime = time.Now()
	s.CalculateStatistics()
}

// GetGroupCount returns the number of duplicate groups
func (s *ScanSession) GetGroupCount() int {
	return len(s.DuplicateGroups)
}

// GetDuplicateFileCount returns the total number of duplicate files
func (s *ScanSession) GetDuplicateFileCount() int {
	count := 0
	for _, group := range s.DuplicateGroups {
		count += group.GetDuplicateCount()
	}
	return count
}
