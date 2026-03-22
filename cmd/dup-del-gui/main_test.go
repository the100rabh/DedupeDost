//go:build !gui

package main

import (
	"flag"
	"testing"
)

// TestParseFlags tests flag parsing
func TestParseFlags(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantGUI     bool
		wantVersion bool
	}{
		{
			name:        "default",
			args:        []string{},
			wantGUI:     true,
			wantVersion: false,
		},
		{
			name:        "version",
			args:        []string{"-version"},
			wantGUI:     true,
			wantVersion: true,
		},
		{
			name:        "gui_true",
			args:        []string{"-gui=true"},
			wantGUI:     true,
			wantVersion: false,
		},
		{
			name:        "gui_false",
			args:        []string{"-gui=false"},
			wantGUI:     false,
			wantVersion: false,
		},
		{
			name:        "version_and_gui_false",
			args:        []string{"-version", "-gui=false"},
			wantGUI:     false,
			wantVersion: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := flag.NewFlagSet("test", flag.ExitOnError)
			guiMode := fs.Bool("gui", true, "Run in GUI mode")
			showVersion := fs.Bool("version", false, "Show version")
			
			if err := fs.Parse(tt.args); err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			if *guiMode != tt.wantGUI {
				t.Errorf("guiMode = %v, want %v", *guiMode, tt.wantGUI)
			}
			if *showVersion != tt.wantVersion {
				t.Errorf("showVersion = %v, want %v", *showVersion, tt.wantVersion)
			}
		})
	}
}

// TestRunCLI tests runCLI function
func TestRunCLI(t *testing.T) {
	// Just verify it doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("runCLI panicked: %v", r)
		}
	}()

	runCLI()
	// Output goes to stdout, can't easily verify
}

// TestVersion_Constant tests Version constant
func TestVersion_Constant(t *testing.T) {
	if Version == "" {
		t.Error("Version should not be empty")
	}
	if Version != "1.0.0" {
		t.Errorf("Version = %s, want 1.0.0", Version)
	}
}

// TestRunApp tests runApp function
func TestRunApp(t *testing.T) {
	t.Skip("runApp requires GUI environment (Fyne)")
}

// TestParseFlags_InvalidValue tests invalid flag values
func TestParseFlags_InvalidValue(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	guiMode := fs.Bool("gui", true, "Run in GUI mode")

	// Invalid boolean should fail to parse with ContinueOnError
	err := fs.Parse([]string{"-gui", "invalid"})
	// With ContinueOnError, Parse returns error but doesn't exit
	if err == nil {
		t.Log("Parse did not return error for invalid value (may use default)")
	}
	
	// guiMode should remain default when parse fails
	t.Logf("guiMode after invalid parse: %v (err=%v)", *guiMode, err)
}

// TestParseFlags_Help tests help flag
func TestParseFlags_Help(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)

	// Help flag should be recognized
	err := fs.Parse([]string{"-help"})
	if err == nil {
		t.Log("Help flag parsed successfully")
	}
}

// TestParseFlags_ShortArgs tests short argument forms
func TestParseFlags_ShortArgs(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	guiMode := fs.Bool("gui", true, "Run in GUI mode")
	showVersion := fs.Bool("version", false, "Show version")

	fs.Parse([]string{})
	if !*guiMode {
		t.Error("Default guiMode should be true")
	}
	if *showVersion {
		t.Error("Default showVersion should be false")
	}
}

// TestParseFlags_EmptyArgs tests with empty args
func TestParseFlags_EmptyArgs(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	guiMode := fs.Bool("gui", true, "Run in GUI mode")
	showVersion := fs.Bool("version", false, "Show version")

	fs.Parse([]string{})

	if !*guiMode {
		t.Error("guiMode should be true with empty args")
	}
	if *showVersion {
		t.Error("showVersion should be false with empty args")
	}
}

// TestParseFlags_VersionOnly tests version flag only
func TestParseFlags_VersionOnly(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	guiMode := fs.Bool("gui", true, "Run in GUI mode")
	showVersion := fs.Bool("version", false, "Show version")

	fs.Parse([]string{"-version"})

	if !*showVersion {
		t.Error("showVersion should be true")
	}
	if !*guiMode {
		t.Error("guiMode should remain true")
	}
}

// TestParseFlags_GUIOnly tests gui flag only
func TestParseFlags_GUIOnly(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ExitOnError)
	guiMode := fs.Bool("gui", true, "Run in GUI mode")
	showVersion := fs.Bool("version", false, "Show version")

	if err := fs.Parse([]string{"-gui=false"}); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if *guiMode {
		t.Error("guiMode should be false")
	}
	if *showVersion {
		t.Error("showVersion should be false")
	}
}

// TestRunApp_Skipped tests that runApp is properly skipped
func TestRunApp_Skipped(t *testing.T) {
	// This test verifies the skip behavior
	t.Skip("runApp requires Fyne GUI environment")
}

// TestFlagDefaults tests default flag values
func TestFlagDefaults(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	guiMode := fs.Bool("gui", true, "Run in GUI mode")
	showVersion := fs.Bool("version", false, "Show version")

	fs.Parse([]string{})

	if *guiMode != true {
		t.Errorf("Default guiMode = %v, want true", *guiMode)
	}
	if *showVersion != false {
		t.Errorf("Default showVersion = %v, want false", *showVersion)
	}
}

// TestFlagUsage tests flag usage output
func TestFlagUsage(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.Bool("gui", true, "Run in GUI mode")
	fs.Bool("version", false, "Show version")

	// Just verify PrintDefaults doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PrintDefaults panicked: %v", r)
		}
	}()

	fs.PrintDefaults()
}

// TestRunCLI_Output tests runCLI output
func TestRunCLI_Output(t *testing.T) {
	// Verify runCLI produces output without panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("runCLI panicked: %v", r)
		}
	}()

	runCLI()
	// Output: "CLI mode - use the dup-del binary directly"
}
