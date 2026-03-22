package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dupdel/dup-del/internal/ui"
)

const Version = "1.0.0"

func main() {
	// Parse flags
	guiMode := flag.Bool("gui", true, "Run in GUI mode")
	showVersion := flag.Bool("version", false, "Show version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("DupDel v%s\n", Version)
		return
	}

	if !*guiMode {
		// Run CLI mode (import from cmd/dup-del)
		runCLI()
		return
	}

	// Run GUI mode
	runApp()
}

func runApp() {
	app := ui.NewDupDelApp()
	app.Run()
}

func runCLI() {
	// Delegate to CLI package
	cliPath := os.Args[0]
	os.Args = append([]string{cliPath}, os.Args[1:]...)
	
	// Import and run CLI
	// This is a placeholder - in production, you'd call the CLI main directly
	fmt.Println("CLI mode - use the dup-del binary directly")
}
