package platform

import (
	"fmt"
	"os/exec"
	"runtime"
)

// OpenInDefaultPlayer opens a file in the system's default application
func OpenInDefaultPlayer(filePath string) error {
	switch runtime.GOOS {
	case "windows":
		return openWindows(filePath)
	case "darwin":
		return openMacOS(filePath)
	case "linux":
		return openLinux(filePath)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// openWindows opens file using Windows default application
func openWindows(filePath string) error {
	// Use "start" command which opens file with default application
	cmd := exec.Command("cmd", "/c", "start", "", filePath)
	return cmd.Start()
}

// openMacOS opens file using macOS default application (QuickTime, etc.)
func openMacOS(filePath string) error {
	// Use "open" command which opens file with default application
	cmd := exec.Command("open", filePath)
	return cmd.Start()
}

// openLinux opens file using xdg-open (freedesktop.org standard)
func openLinux(filePath string) error {
	// Try xdg-open first (works on most Linux desktop environments)
	cmd := exec.Command("xdg-open", filePath)
	err := cmd.Start()
	if err != nil {
		// Fallback to common media players
		players := []string{"vlc", "mpv", "totem", "mplayer"}
		for _, player := range players {
			if path, err := exec.LookPath(player); err == nil {
				cmd := exec.Command(path, filePath)
				return cmd.Start()
			}
		}
		return fmt.Errorf("no suitable video player found: %w", err)
	}
	return nil
}
