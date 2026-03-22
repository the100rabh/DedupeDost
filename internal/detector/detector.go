package detector

import (
	"context"
	"fmt"
	"time"

	"github.com/dupdel/dup-del/internal/hasher"
	"github.com/dupdel/dup-del/internal/scanner"
	"github.com/dupdel/dup-del/pkg/models"
)

// DetectionStats contains statistics about the detection process
type DetectionStats struct {
	TotalFiles       int           `json:"total_files"`
	TotalSize        int64         `json:"total_size"`
	FilesWithHash    int           `json:"files_with_hash"`
	DuplicateFiles   int           `json:"duplicate_files"`
	DuplicateGroups  int           `json:"duplicate_groups"`
	RecoverableSize  int64         `json:"recoverable_size"`
	ScanDuration     time.Duration `json:"scan_duration"`
	HashDuration     time.Duration `json:"hash_duration"`
}

// Detector finds duplicate files
type Detector struct {
	scanDir   string
	options   scanner.ScanOptions
	sizeMap   map[int64][]models.FileEntry
	hashMap   map[string][]models.FileEntry
	groups    []models.DuplicateGroup
	stats     *DetectionStats
	ctx       context.Context
	cancel    context.CancelFunc
	hasher    *hasher.Hasher
	startTime time.Time
}

// NewDetector creates a new Detector
func NewDetector(scanDir string, options scanner.ScanOptions) *Detector {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &Detector{
		scanDir:   scanDir,
		options:   options,
		sizeMap:   make(map[int64][]models.FileEntry),
		hashMap:   make(map[string][]models.FileEntry),
		groups:    make([]models.DuplicateGroup, 0),
		stats:     &DetectionStats{},
		ctx:       ctx,
		cancel:    cancel,
		hasher:    hasher.NewHasher(),
	}
}

// Detect runs the duplicate detection algorithm
func (d *Detector) Detect() error {
	d.startTime = time.Now()
	
	// Step 1: Scan all files and group by size
	if err := d.scanAndGroupBySize(); err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}
	
	// Check for cancellation
	select {
	case <-d.ctx.Done():
		return d.ctx.Err()
	default:
	}
	
	// Step 2: Calculate hashes for files with matching sizes
	hashStartTime := time.Now()
	if err := d.calculateHashes(); err != nil {
		return fmt.Errorf("hash calculation failed: %w", err)
	}
	d.stats.HashDuration = time.Since(hashStartTime)
	
	// Check for cancellation
	select {
	case <-d.ctx.Done():
		return d.ctx.Err()
	default:
	}
	
	// Step 3: Group by hash and create duplicate groups
	d.groupByHash()
	
	// Step 4: Calculate final statistics
	d.calculateStats()
	
	return nil
}

// scanAndGroupBySize scans files and groups them by size
func (d *Detector) scanAndGroupBySize() error {
	scanner := scanner.NewScanner(d.scanDir, d.options)
	
	// Pre-count files for progress
	totalFiles, err := scanner.CountFiles()
	if err != nil {
		return err
	}
	d.stats.TotalFiles = int(totalFiles)
	
	fileChan, errorChan, err := scanner.Scan()
	if err != nil {
		return err
	}
	
	for {
		select {
		case <-d.ctx.Done():
			return d.ctx.Err()
		case entry, ok := <-fileChan:
			if !ok {
				return nil
			}
			
			// Group by size
			files := d.sizeMap[entry.Size]
			files = append(files, entry)
			d.sizeMap[entry.Size] = files
			
			d.stats.TotalSize += entry.Size
		case err, ok := <-errorChan:
			if ok && err != nil {
				// Log but continue
				fmt.Printf("Scan error: %v\n", err)
			}
		}
	}
}

// calculateHashes calculates hashes for files that have size matches
func (d *Detector) calculateHashes() error {
	// Collect files that need hashing (those with matching sizes)
	var filesToHash []string
	fileMap := make(map[string]models.FileEntry)
	
	for _, files := range d.sizeMap {
		if len(files) < 2 {
			continue // Skip unique sizes
		}
		
		for _, file := range files {
			filesToHash = append(filesToHash, file.Path)
			fileMap[file.Path] = file
		}
	}
	
	if len(filesToHash) == 0 {
		return nil
	}
	
	// Hash files in parallel
	resultChan := d.hasher.HashMultipleWithContext(d.ctx, filesToHash, 4)
	
	for result := range resultChan {
		if result.Cancelled {
			return d.ctx.Err()
		}
		
		if result.Error != nil {
			continue // Skip files that failed to hash
		}
		
		entry := fileMap[result.Path]
		entry.Hash = result.Hash
		
		// Group by hash
		files := d.hashMap[result.Hash]
		files = append(files, entry)
		d.hashMap[result.Hash] = files
		
		d.stats.FilesWithHash++
	}
	
	return nil
}

// groupByHash creates duplicate groups from hash groups
func (d *Detector) groupByHash() {
	groupID := 0
	
	for hash, files := range d.hashMap {
		if len(files) < 2 {
			continue // Not a duplicate
		}
		
		group := models.DuplicateGroup{
			ID:        fmt.Sprintf("grp_%d", groupID),
			Hash:      hash,
			Size:      files[0].Size,
			Extension: files[0].Extension,
			FileType:  files[0].FileType,
			Files:     files,
		}
		
		d.groups = append(d.groups, group)
		groupID++
	}
}

// calculateStats calculates final detection statistics
func (d *Detector) calculateStats() {
	d.stats.DuplicateGroups = len(d.groups)
	
	duplicateFiles := 0
	recoverableSize := int64(0)
	
	for _, group := range d.groups {
		duplicateFiles += group.GetDuplicateCount()
		recoverableSize += group.GetRecoverableSize()
	}
	
	d.stats.DuplicateFiles = duplicateFiles
	d.stats.RecoverableSize = recoverableSize
	d.stats.ScanDuration = time.Since(d.startTime)
}

// GetGroups returns the detected duplicate groups
func (d *Detector) GetGroups() []models.DuplicateGroup {
	return d.groups
}

// GetStatistics returns the detection statistics
func (d *Detector) GetStatistics() DetectionStats {
	return *d.stats
}

// Cancel stops the detection process
func (d *Detector) Cancel() {
	d.cancel()
}

// GetProgress returns the current detection progress
func (d *Detector) GetProgress() DetectionProgress {
	return DetectionProgress{
		FilesScanned:    d.stats.TotalFiles,
		FilesWithHash:   d.stats.FilesWithHash,
		DuplicateGroups: len(d.groups),
	}
}

// DetectionProgress represents the current detection progress
type DetectionProgress struct {
	FilesScanned    int
	FilesWithHash   int
	DuplicateGroups int
}
