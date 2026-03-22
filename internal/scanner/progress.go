package scanner

import (
	"sync"
	"time"
)

// ProgressCallback is a function type for progress updates
type ProgressCallback func(progress ScanProgress)

// ProgressTracker tracks and reports scan progress
type ProgressTracker struct {
	progress  ScanProgress
	callbacks []ProgressCallback
	mu        sync.RWMutex
	startTime time.Time
}

// NewProgressTracker creates a new progress tracker
func NewProgressTracker() *ProgressTracker {
	return &ProgressTracker{
		progress:  ScanProgress{},
		callbacks: make([]ProgressCallback, 0),
		startTime: time.Now(),
	}
}

// AddCallback adds a callback to be called on progress updates
func (p *ProgressTracker) AddCallback(callback ProgressCallback) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.callbacks = append(p.callbacks, callback)
}

// Update updates the progress and notifies callbacks
func (p *ProgressTracker) Update(filesScanned, totalFiles int64, currentFile string, totalSize int64) {
	p.mu.Lock()
	
	p.progress.FilesScanned = filesScanned
	p.progress.TotalFiles = totalFiles
	p.progress.CurrentFile = currentFile
	p.progress.TotalSize = totalSize
	
	if totalFiles > 0 {
		p.progress.PercentComplete = float64(filesScanned) / float64(totalFiles) * 100
	}
	
	// Copy callbacks to call outside lock
	callbacks := make([]ProgressCallback, len(p.callbacks))
	copy(callbacks, p.callbacks)
	
	p.mu.Unlock()
	
	// Notify callbacks
	for _, callback := range callbacks {
		callback(p.progress)
	}
}

// GetProgress returns the current progress
func (p *ProgressTracker) GetProgress() ScanProgress {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.progress
}

// GetElapsed returns the elapsed time since tracking started
func (p *ProgressTracker) GetElapsed() time.Duration {
	return time.Since(p.startTime)
}

// GetFilesPerSecond calculates the average files scanned per second
func (p *ProgressTracker) GetFilesPerSecond() float64 {
	elapsed := p.GetElapsed().Seconds()
	if elapsed == 0 {
		return 0
	}
	
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	return float64(p.progress.FilesScanned) / elapsed
}

// GetETA estimates the remaining time based on current progress
func (p *ProgressTracker) GetETA() time.Duration {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	if p.progress.TotalFiles == 0 || p.progress.FilesScanned == 0 {
		return 0
	}
	
	elapsed := time.Since(p.startTime)
	rate := float64(p.progress.FilesScanned) / elapsed.Seconds()
	
	if rate == 0 {
		return 0
	}
	
	remaining := p.progress.TotalFiles - p.progress.FilesScanned
	etaSeconds := float64(remaining) / rate
	
	return time.Duration(etaSeconds) * time.Second
}

// Reset resets the progress tracker
func (p *ProgressTracker) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	p.progress = ScanProgress{}
	p.startTime = time.Now()
}
