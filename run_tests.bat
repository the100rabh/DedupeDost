@echo off
REM DedupeDost Test Runner Script (Windows)
REM ====================================
REM This script runs all tests for the DedupeDost project.
REM
REM Usage:
REM   run_tests.bat              - Run all tests with coverage
REM   run_tests.bat -v           - Run with verbose output
REM   run_tests.bat -race        - Run with race detector
REM   run_tests.bat -html        - Generate HTML coverage report
REM   run_tests.bat -clean       - Clean test cache before running
REM   run_tests.bat -help        - Show help message

setlocal enabledelayedexpansion

REM Default options
set VERBOSE=
set RACE=
set HTML_REPORT=0
set CLEAN=0
set COVERAGE_PROFILE=coverage.out

REM Parse arguments
:parse_args
if "%~1"=="" goto :run_tests
if /i "%~1"=="-v" set VERBOSE=-v& shift & goto :parse_args
if /i "%~1"=="--verbose" set VERBOSE=-v& shift & goto :parse_args
if /i "%~1"=="-race" set RACE=-race& shift & goto :parse_args
if /i "%~1"=="-html" set HTML_REPORT=1& shift & goto :parse_args
if /i "%~1"=="-clean" set CLEAN=1& shift & goto :parse_args
if /i "%~1"=="-help" goto :show_help
if /i "%~1"=="-h" goto :show_help
echo Unknown option: %~1
echo Use -help for usage information
goto :eof

:show_help
echo DedupeDost Test Runner Script (Windows)
echo ===================================
echo.
echo Usage:
echo   run_tests.bat [OPTIONS]
echo.
echo Options:
echo   -v, --verbose    Run tests with verbose output
echo   -race            Run with race detector
echo   -html            Generate HTML coverage report
echo   -clean           Clean test cache before running
echo   -help, -h        Show this help message
echo.
echo Examples:
echo   run_tests.bat                  Run all tests with coverage
echo   run_tests.bat -v               Run with verbose output
echo   run_tests.bat -race            Run with race detector
echo   run_tests.bat -clean -v        Clean cache and run verbose
goto :eof

:run_tests
echo.
echo ========================================
echo   DedupeDost Test Suite
echo ========================================
echo.

if %CLEAN%==1 (
    echo Cleaning test cache...
    go clean -testcache
    echo [OK] Test cache cleaned
    echo.
)

echo Running tests...
echo.

go test ./pkg/... ./internal/detector/... ./internal/hasher/... ./internal/rules/... ./internal/scanner/... ./internal/script/... ./internal/platform/... ./cmd/dedupedost/... %VERBOSE% %RACE% -coverprofile=%COVERAGE_PROFILE%

if errorlevel 1 (
    echo.
    echo [FAIL] Some tests failed
    goto :eof
)

echo.
echo [OK] All tests passed
echo.

echo ========================================
echo Coverage Summary:
echo ========================================
go tool cover -func=%COVERAGE_PROFILE% | findstr /C:"total:"

if %HTML_REPORT%==1 (
    echo.
    echo Generating HTML coverage report...
    go tool cover -html=%COVERAGE_PROFILE% -o coverage.html
    echo [OK] HTML report generated: coverage.html
    start coverage.html
)

echo.
echo ========================================
echo   Test run completed successfully!
echo ========================================
echo.

endlocal
