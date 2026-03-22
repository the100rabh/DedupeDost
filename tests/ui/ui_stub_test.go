//go:build !gui

package ui_test

import (
	"testing"
)

// TestStub tests that headless tests can run
func TestStub(t *testing.T) {
	// This is a placeholder test
	// Full UI testing requires GUI backend (OpenGL/X11)
	// 
	// To run UI tests on your system:
	// 1. Install dependencies: sudo apt-get install libgl1-mesa-dev xorg-dev
	// 2. Run: go test ./tests/ui/... -v
	t.Log("UI tests require GUI backend. Use headless tests for CI/CD.")
}
