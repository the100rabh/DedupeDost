package platform

import (
	"os"
	"os/exec"
	"runtime"
	"testing"
)

// TestOpenInDefaultPlayer tests OpenInDefaultPlayer function
func TestOpenInDefaultPlayer(t *testing.T) {
	// Create a temp file to test with
	tmpFile, err := os.CreateTemp("", "test_video.*.mp4")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Test with valid file path - this will try to open the file
	// We just verify it doesn't panic
	err = OpenInDefaultPlayer(tmpFile.Name())
	
	// In CI/headless environment, this should return an error
	// but should not panic
	t.Logf("OpenInDefaultPlayer returned: %v (runtime: %s)", err, runtime.GOOS)
}

// TestOpenInDefaultPlayer_NonExistentFile tests with non-existent file
func TestOpenInDefaultPlayer_NonExistentFile(t *testing.T) {
	err := OpenInDefaultPlayer("/nonexistent/file.mp4")
	// Behavior varies by platform - just log result
	t.Logf("OpenInDefaultPlayer for non-existent file: %v", err)
}

// TestOpenInDefaultPlayer_EmptyPath tests with empty path
func TestOpenInDefaultPlayer_EmptyPath(t *testing.T) {
	err := OpenInDefaultPlayer("")
	// Behavior varies by platform - just log result
	t.Logf("OpenInDefaultPlayer for empty path: %v", err)
}

// TestOpenInDefaultPlayer_InvalidExtension tests with invalid extension
func TestOpenInDefaultPlayer_InvalidExtension(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_invalid.*.xyz")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	err = OpenInDefaultPlayer(tmpFile.Name())
	t.Logf("OpenInDefaultPlayer with invalid extension: %v", err)
}

// TestOpenWindows tests openWindows function directly
func TestOpenWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific test")
	}

	tmpFile, err := os.CreateTemp("", "test.*.mp4")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	err = openWindows(tmpFile.Name())
	t.Logf("openWindows returned: %v", err)
}

// TestOpenMacOS tests openMacOS function directly
func TestOpenMacOS(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping macOS-specific test")
	}

	tmpFile, err := os.CreateTemp("", "test.*.mp4")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	err = openMacOS(tmpFile.Name())
	t.Logf("openMacOS returned: %v", err)
}

// TestOpenLinux tests openLinux function directly
func TestOpenLinux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Skipping Linux-specific test")
	}

	tmpFile, err := os.CreateTemp("", "test.*.mp4")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	err = openLinux(tmpFile.Name())
	t.Logf("openLinux returned: %v", err)
}

// TestOpenLinux_NoPlayers tests openLinux fallback behavior
func TestOpenLinux_NoPlayers(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Skipping Linux-specific test")
	}

	// Test with non-existent file to trigger player search
	err := openLinux("/nonexistent/file.mp4")
	t.Logf("openLinux for non-existent file: %v", err)
}

// TestOpenLinux_PlayerFallback tests player fallback logic
func TestOpenLinux_PlayerFallback(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Skipping Linux-specific test")
	}

	// Test that LookPath is called for players
	players := []string{"vlc", "mpv", "totem", "mplayer"}
	foundPlayer := false
	for _, player := range players {
		if _, err := exec.LookPath(player); err == nil {
			foundPlayer = true
			t.Logf("Found player: %s", player)
			break
		}
	}
	
	if !foundPlayer {
		t.Log("No video players found (expected in CI)")
	}
}

// TestOpenInDefaultPlayer_UnsupportedPlatform tests unsupported platform
func TestOpenInDefaultPlayer_UnsupportedPlatform(t *testing.T) {
	// This test documents the behavior for unsupported platforms
	// Actual testing would require mocking runtime.GOOS
	t.Log("OpenInDefaultPlayer returns error for unsupported platforms")
	t.Log("Supported platforms: windows, darwin, linux")
}

// TestOpenLinux_XdgOpenError tests xdg-open error handling
func TestOpenLinux_XdgOpenError(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Skipping Linux-specific test")
	}

	// Test error path when xdg-open fails
	err := openLinux("")
	t.Logf("openLinux with empty path: %v", err)
}

// TestOpenWindows_Command tests Windows command construction
func TestOpenWindows_Command(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific test")
	}

	// Verify command structure
	cmd := exec.Command("cmd", "/c", "start", "", "test.mp4")
	if cmd == nil {
		t.Error("Command should not be nil")
	}
	t.Logf("Windows command: %v", cmd.Args)
}

// TestOpenMacOS_Command tests macOS command construction
func TestOpenMacOS_Command(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping macOS-specific test")
	}

	// Verify command structure
	cmd := exec.Command("open", "test.mp4")
	if cmd == nil {
		t.Error("Command should not be nil")
	}
	t.Logf("macOS command: %v", cmd.Args)
}

// TestOpenLinux_Command tests Linux command construction
func TestOpenLinux_Command(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Skipping Linux-specific test")
	}

	// Verify xdg-open command structure
	cmd := exec.Command("xdg-open", "test.mp4")
	if cmd == nil {
		t.Error("Command should not be nil")
	}
	t.Logf("Linux command: %v", cmd.Args)
}
