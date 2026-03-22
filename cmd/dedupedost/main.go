package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/the100rabh/DedupeDost/internal/detector"
	"github.com/the100rabh/DedupeDost/internal/script"
	"github.com/the100rabh/DedupeDost/internal/scanner"
	"github.com/the100rabh/DedupeDost/pkg/models"
)

const Version = "1.0.0"

// CLI represents the command-line interface
type CLI struct {
	dir         string
	minSize     int64
	maxSize     int64
	types       string
	output      string
	recursive   bool
	ignoreHidden bool
	hashOnly    bool
	verbose     bool
	showVersion bool
	showHelp    bool
}

func main() {
	cli := parseFlags()
	
	if cli.showVersion {
		fmt.Printf("DedupeDost v%s - Duplicate File Scanner\n", Version)
		return
	}
	
	if cli.showHelp {
		printHelp()
		return
	}
	
	// Run the scanner
	if err := cli.run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// parseFlags parses command-line flags
func parseFlags() *CLI {
	return parseFlagsWithSet(flag.CommandLine, os.Args[1:])
}

// parseFlagsWithSet parses flags using a custom FlagSet (for testing)
func parseFlagsWithSet(fs *flag.FlagSet, args []string) *CLI {
	cli := &CLI{}

	fs.StringVar(&cli.dir, "dir", ".", "Directory to scan")
	fs.StringVar(&cli.dir, "d", ".", "Directory to scan (shorthand)")
	fs.Int64Var(&cli.minSize, "min-size", 0, "Minimum file size in bytes")
	fs.Int64Var(&cli.maxSize, "max-size", 0, "Maximum file size in bytes (0 = unlimited)")
	fs.StringVar(&cli.types, "types", "", "Comma-separated file extensions to scan (e.g., jpg,png,txt)")
	fs.StringVar(&cli.output, "output", "cleanup_duplicates.sh", "Output script path")
	fs.StringVar(&cli.output, "o", "cleanup_duplicates.sh", "Output script path (shorthand)")
	fs.BoolVar(&cli.recursive, "recursive", true, "Scan subdirectories")
	fs.BoolVar(&cli.recursive, "r", true, "Scan subdirectories (shorthand)")
	fs.BoolVar(&cli.ignoreHidden, "ignore-hidden", false, "Skip hidden files and directories")
	fs.BoolVar(&cli.hashOnly, "hash-only", false, "Skip preview info, hash comparison only")
	fs.BoolVar(&cli.verbose, "verbose", false, "Verbose output")
	fs.BoolVar(&cli.verbose, "v", false, "Verbose output (shorthand)")
	fs.BoolVar(&cli.showVersion, "version", false, "Show version")
	fs.BoolVar(&cli.showHelp, "help", false, "Show help")
	fs.BoolVar(&cli.showHelp, "h", false, "Show help (shorthand)")

	fs.Parse(args)

	return cli
}

// run executes the CLI command
func (c *CLI) run() error {
	// Get absolute path
	absDir, err := getAbsolutePath(c.dir)
	if err != nil {
		return fmt.Errorf("invalid directory: %w", err)
	}
	c.dir = absDir
	
	// Check directory exists
	if _, err := os.Stat(c.dir); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", c.dir)
	}
	
	fmt.Println()
	fmt.Printf("DedupeDost v%s - Duplicate File Scanner\n", Version)
	fmt.Println("======================================")
	fmt.Println()
	
	// Parse extensions
	extensions := parseExtensions(c.types)
	
	// Create scan options
	scanOptions := scanner.ScanOptions{
		Recursive:      c.recursive,
		IncludeHidden:  !c.ignoreHidden,
		FollowSymlinks: false,
		MinSize:        c.minSize,
		MaxSize:        c.maxSize,
		Extensions:     extensions,
	}
	
	// Display scan info
	fmt.Printf("Scanning: %s\n", c.dir)
	fmt.Printf("Options: Recursive=%v, Hidden=%v, MinSize=%s\n",
		scanOptions.Recursive,
		scanOptions.IncludeHidden,
		models.FormatSize(scanOptions.MinSize))
	if len(extensions) > 0 {
		fmt.Printf("        Types: %s\n", strings.Join(extensions, ", "))
	}
	fmt.Println()
	
	// Set up context with cancellation on interrupt
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// Handle Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\n\nInterrupted! Cleaning up...")
		cancel()
	}()
	
	// Create and run detector
	fmt.Println("Scanning files...")
	det := detector.NewDetector(c.dir, scanOptions)
	
	// Run detection with progress
	_ = ctx // Context available for future use
	if err := det.Detect(); err != nil {
		if err == context.Canceled {
			fmt.Println("Scan cancelled.")
			return nil
		}
		return fmt.Errorf("detection failed: %w", err)
	}
	
	// Get results
	stats := det.GetStatistics()
	groups := det.GetGroups()
	
	// Display results
	fmt.Println()
	fmt.Println("======================================")
	fmt.Println("Scan Complete!")
	fmt.Println("======================================")
	fmt.Printf("Total files scanned: %d\n", stats.TotalFiles)
	fmt.Printf("Total size: %s\n", models.FormatSize(stats.TotalSize))
	fmt.Printf("Duplicate groups: %d\n", stats.DuplicateGroups)
	fmt.Printf("Duplicate files: %d\n", stats.DuplicateFiles)
	fmt.Printf("Recoverable space: %s\n", models.FormatSize(stats.RecoverableSize))
	fmt.Printf("Scan duration: %v\n", stats.ScanDuration)
	fmt.Println()
	
	if len(groups) == 0 {
		fmt.Println("No duplicates found!")
		return nil
	}
	
	// Show top duplicates
	fmt.Println("Top duplicates by size:")
	displayTopDuplicates(groups, 5)
	fmt.Println()

	// Create scan session
	session := models.NewScanSession(c.dir)
	// Convert value slice to pointer slice
	session.DuplicateGroups = make([]*models.DuplicateGroup, len(groups))
	for i := range groups {
		session.DuplicateGroups[i] = &groups[i]
	}
	session.TotalFiles = stats.TotalFiles
	session.TotalSize = stats.TotalSize
	session.Finalize()
	
	// Auto-mark files for deletion (keep first in each group)
	for i := range session.DuplicateGroups {
		session.DuplicateGroups[i].MarkForDeletion([]int{0})
	}
	session.CalculateStatistics()
	
	// Generate script
	fmt.Printf("Generating cleanup script: %s\n", c.output)
	generator := script.NewGenerator(c.output)
	scriptPath, err := generator.Generate(session)
	if err != nil {
		return fmt.Errorf("failed to generate script: %w", err)
	}
	
	fmt.Println()
	fmt.Println("======================================")
	fmt.Println("Done!")
	fmt.Println("======================================")
	fmt.Printf("Cleanup script: %s\n", scriptPath)
	fmt.Println()
	fmt.Println("To preview deletions:")
	fmt.Printf("  DRY_RUN=true %s\n", scriptPath)
	fmt.Println()
	fmt.Println("To execute deletions:")
	fmt.Printf("  %s\n", scriptPath)
	fmt.Println()
	
	return nil
}

// displayTopDuplicates shows the top N duplicate groups by size
func displayTopDuplicates(groups []models.DuplicateGroup, n int) {
	// Sort by recoverable size (simple bubble sort for display)
	sorted := make([]models.DuplicateGroup, len(groups))
	copy(sorted, groups)
	
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i].GetRecoverableSize() < sorted[j].GetRecoverableSize() {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	
	// Display top N
	count := n
	if count > len(sorted) {
		count = len(sorted)
	}
	
	for i := 0; i < count; i++ {
		g := sorted[i]
		fmt.Printf("  %d. %s (%d copies, %s each) - Recoverable: %s\n",
			i+1,
			g.Files[0].Name,
			len(g.Files),
			g.GetDisplaySize(),
			g.GetDisplayRecoverableSize())
	}
}

// getAbsolutePath returns the absolute path of a directory
func getAbsolutePath(path string) (string, error) {
	if path == "" {
		path = "."
	}
	return os.Getwd()
}

// parseExtensions parses comma-separated extensions
func parseExtensions(types string) []string {
	if types == "" {
		return nil
	}
	
	parts := strings.Split(types, ",")
	extensions := make([]string, 0, len(parts))
	
	for _, ext := range parts {
		ext = strings.TrimSpace(ext)
		if ext != "" {
			if !strings.HasPrefix(ext, ".") {
				ext = "." + ext
			}
			extensions = append(extensions, strings.ToLower(ext))
		}
	}
	
	return extensions
}

// printHelp prints the help message
func printHelp() {
	fmt.Println()
	fmt.Printf("DedupeDost v%s - Duplicate File Scanner\n", Version)
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  dedupedost [OPTIONS]")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  -dir, -d PATH       Directory to scan (default: current directory)")
	fmt.Println("  -min-size BYTES     Minimum file size in bytes (default: 0)")
	fmt.Println("  -max-size BYTES     Maximum file size in bytes (default: unlimited)")
	fmt.Println("  -types LIST         Comma-separated extensions to scan (e.g., jpg,png,txt)")
	fmt.Println("  -output, -o PATH    Output script path (default: cleanup_duplicates.sh)")
	fmt.Println("  -recursive, -r      Scan subdirectories (default: true)")
	fmt.Println("  -ignore-hidden      Skip hidden files and directories")
	fmt.Println("  -hash-only          Skip preview info, hash comparison only")
	fmt.Println("  -verbose, -v        Verbose output")
	fmt.Println("  -version            Show version")
	fmt.Println("  -help, -h           Show this help message")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  dedupedost")
	fmt.Println("  dedupedost -dir /home/user/documents")
	fmt.Println("  dedupedost -d /home/user/photos -types jpg,png,gif")
	fmt.Println("  dedupedost -min-size 1024 -max-size 104857600")
	fmt.Println("  dedupedost -ignore-hidden -o /tmp/cleanup.sh")
	fmt.Println()
	fmt.Println("OUTPUT:")
	fmt.Println("  The application generates a bash script that can be reviewed")
	fmt.Println("  before executing. Run with DRY_RUN=true to preview deletions.")
	fmt.Println()
}
