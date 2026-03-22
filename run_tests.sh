#!/bin/bash

# DedupeDost Test Runner Script
# ==========================
# This script runs all tests for the DedupeDost project with various options.
#
# Usage:
#   bash run_tests.sh              # Run all tests with coverage
#   bash run_tests.sh -v           # Run with verbose output
#   bash run_tests.sh -race        # Run with race detector
#   bash run_tests.sh -gui         # Include GUI tests (requires X11/OpenGL)
#   bash run_tests.sh -headless    # Run only non-GUI tests (default)
#   bash run_tests.sh -html        # Generate HTML coverage report
#   bash run_tests.sh -clean       # Clean test cache before running
#   bash run_tests.sh -help        # Show this help message
#
# Note: Run with bash, not sh:
#   bash run_tests.sh    # ✓ Correct
#   sh run_tests.sh      # ✗ Wrong (sh doesn't support bash arrays)

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default options
VERBOSE=""
RACE=""
GUI_MODE=false
HTML_REPORT=false
CLEAN=false
COVERAGE_PROFILE="coverage.out"

# Print colored message
print_msg() {
    local color=$1
    local msg=$2
    echo -e "${color}${msg}${NC}"
}

# Print header
print_header() {
    echo ""
    echo "========================================"
    echo "  DedupeDost Test Suite"
    echo "========================================"
    echo ""
}

# Show help
show_help() {
    cat << EOF
DedupeDost Test Runner Script
=========================

Usage:
  ./run_tests.sh [OPTIONS]

Options:
  -v, --verbose       Run tests with verbose output
  -race               Run with race detector (slower but finds data races)
  -gui                Include GUI tests (requires X11/OpenGL dependencies)
  -headless           Run only non-GUI tests (default mode)
  -html               Generate HTML coverage report (opens in browser)
  -clean              Clean test cache before running tests
  -help, -h           Show this help message

Examples:
  ./run_tests.sh                      # Run all tests with coverage
  ./run_tests.sh -v                   # Run with verbose output
  ./run_tests.sh -race                # Run with race detector
  ./run_tests.sh -gui -html           # Run GUI tests and generate HTML report
  ./run_tests.sh -clean -v            # Clean cache and run with verbose

GUI Test Requirements (Linux):
  sudo apt-get install libgl1-mesa-dev xorg-dev

EOF
    exit 0
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -v|--verbose)
                VERBOSE="-v"
                shift
                ;;
            -race)
                RACE="-race"
                shift
                ;;
            -gui)
                GUI_MODE=true
                shift
                ;;
            -headless)
                GUI_MODE=false
                shift
                ;;
            -html)
                HTML_REPORT=true
                shift
                ;;
            -clean)
                CLEAN=true
                shift
                ;;
            -help|-h)
                show_help
                ;;
            *)
                print_msg "$RED" "Unknown option: $1"
                echo "Use -help for usage information"
                exit 1
                ;;
        esac
    done
}

# Clean test cache
clean_cache() {
    print_msg "$YELLOW" "Cleaning test cache..."
    go clean -testcache
    print_msg "$GREEN" "✓ Test cache cleaned"
}

# Check GUI dependencies
check_gui_deps() {
    if [ "$GUI_MODE" = true ]; then
        print_msg "$YELLOW" "Checking GUI dependencies..."
        
        # Check for X11/OpenGL libraries
        if ! pkg-config --exists gl 2>/dev/null; then
            print_msg "$RED" "Warning: OpenGL development libraries not found"
            print_msg "$YELLOW" "To install on Ubuntu/Debian:"
            echo "  sudo apt-get install libgl1-mesa-dev xorg-dev"
            echo ""
            print_msg "$YELLOW" "Falling back to headless mode..."
            GUI_MODE=false
        else
            print_msg "$GREEN" "✓ GUI dependencies found"
        fi
    fi
}

# Run non-GUI tests
run_headless_tests() {
    print_msg "$BLUE" "Running headless tests (non-GUI packages)..."
    echo ""
    
    # Define non-GUI packages
    PACKAGES=(
        "./pkg/..."
        "./internal/detector/..."
        "./internal/hasher/..."
        "./internal/rules/..."
        "./internal/scanner/..."
        "./internal/script/..."
        "./internal/platform/..."
        "./cmd/dedupedost/..."
        "./tests/detector/..."
        "./tests/scanner/..."
    )
    
    # Run tests
    if go test ${PACKAGES[@]} $VERBOSE $RACE -coverprofile="$COVERAGE_PROFILE"; then
        print_msg "$GREEN" "✓ Headless tests passed"
    else
        print_msg "$RED" "✗ Headless tests failed"
        exit 1
    fi
}

# Run GUI tests
run_gui_tests() {
    print_msg "$BLUE" "Running GUI tests (requires X11/OpenGL)..."
    echo ""
    
    # Define GUI packages
    GUI_PACKAGES=(
        "./internal/preview/..."
        "./internal/ui/..."
        "./cmd/dedupedost-gui/..."
        "./tests/ui/..."
        "./tests/preview/..."
    )
    
    # Run tests
    if go test ${GUI_PACKAGES[@]} $VERBOSE $RACE -coverprofile="$COVERAGE_PROFILE" -covermode=append; then
        print_msg "$GREEN" "✓ GUI tests passed"
    else
        print_msg "$RED" "✗ GUI tests failed"
        exit 1
    fi
}

# Run all tests
run_all_tests() {
    print_msg "$BLUE" "Running all tests..."
    echo ""
    
    if go test ./... $VERBOSE $RACE -coverprofile="$COVERAGE_PROFILE"; then
        print_msg "$GREEN" "✓ All tests passed"
    else
        print_msg "$RED" "✗ Some tests failed"
        exit 1
    fi
}

# Show coverage summary
show_coverage() {
    echo ""
    print_msg "$BLUE" "Coverage Summary:"
    echo "========================================"
    
    if [ -f "$COVERAGE_PROFILE" ]; then
        go tool cover -func="$COVERAGE_PROFILE" | tail -1
        echo ""
        print_msg "$YELLOW" "Detailed coverage by package:"
        go tool cover -func="$COVERAGE_PROFILE" | grep -E "^github|total:"
    else
        print_msg "$RED" "Coverage profile not found: $COVERAGE_PROFILE"
    fi
}

# Generate HTML report
generate_html() {
    if [ -f "$COVERAGE_PROFILE" ]; then
        print_msg "$BLUE" "Generating HTML coverage report..."
        go tool cover -html="$COVERAGE_PROFILE" -o coverage.html
        print_msg "$GREEN" "✓ HTML report generated: coverage.html"
        
        # Try to open in browser
        if command -v xdg-open &> /dev/null; then
            xdg-open coverage.html &
        elif command -v open &> /dev/null; then
            open coverage.html &
        fi
    else
        print_msg "$RED" "Cannot generate HTML: coverage profile not found"
    fi
}

# Main execution
main() {
    print_header
    
    parse_args "$@"
    
    if [ "$CLEAN" = true ]; then
        clean_cache
    fi
    
    check_gui_deps
    
    echo ""
    print_msg "$BLUE" "Starting tests..."
    echo ""
    
    if [ "$GUI_MODE" = true ]; then
        run_all_tests
    else
        run_headless_tests
    fi
    
    show_coverage
    
    if [ "$HTML_REPORT" = true ]; then
        generate_html
    fi
    
    echo ""
    print_msg "$GREEN" "========================================"
    print_msg "$GREEN" "  Test run completed successfully!"
    print_msg "$GREEN" "========================================"
    echo ""
}

# Run main function
main "$@"
