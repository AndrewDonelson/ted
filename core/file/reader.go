// Package file implements file I/O operations for the editor.
//
// It provides functions for reading and writing files with proper
// UTF-8 handling, line ending detection, and error handling.
// This package has no UI dependencies and is purely focused on I/O.
package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LineEnding represents the type of line ending in a file.
type LineEnding string

const (
	// LineEndingLF represents Unix-style line endings (\n)
	LineEndingLF LineEnding = "LF"
	// LineEndingCRLF represents Windows-style line endings (\r\n)
	LineEndingCRLF LineEnding = "CRLF"
	// LineEndingCR represents old Mac-style line endings (\r)
	LineEndingCR LineEnding = "CR"
	// LineEndingUnknown represents unknown or mixed line endings
	LineEndingUnknown LineEnding = "Unknown"
)

// FileInfo contains metadata about a file.
type FileInfo struct {
	Path       string
	Size       int64
	LineEnding LineEnding
	Encoding   string // Always "UTF-8" for now
}

// ReadFile reads a file and returns its contents as a slice of lines.
// It handles UTF-8 encoding, detects line endings, and validates the file path.
// Returns an error if the file cannot be read or the path is invalid.
//
// Example:
//
//	lines, err := ReadFile("example.txt")
//	if err != nil {
//	    log.Fatal(err)
//	}
func ReadFile(path string) ([]string, error) {
	// Validate and clean path
	cleanPath, err := validatePath(path)
	if err != nil {
		return nil, err
	}

	// Check if file exists and is readable
	info, err := os.Stat(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("stat file %q: %w", cleanPath, err)
	}

	if info.IsDir() {
		return nil, fmt.Errorf("path %q is a directory", cleanPath)
	}

	// Read file
	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("read file %q: %w", cleanPath, err)
	}

	// Split into lines (handle different line endings)
	lines := splitLines(data)
	return lines, nil
}

// ReadFileWithInfo reads a file and returns both the contents and file metadata.
func ReadFileWithInfo(path string) ([]string, *FileInfo, error) {
	lines, err := ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, nil, fmt.Errorf("stat file %q: %w", path, err)
	}

	// Detect line ending
	lineEnding := detectLineEnding(path)

	fileInfo := &FileInfo{
		Path:       path,
		Size:       info.Size(),
		LineEnding: lineEnding,
		Encoding:   "UTF-8",
	}

	return lines, fileInfo, nil
}

// validatePath validates and cleans a file path.
// Returns an absolute path or an error if the path is invalid.
func validatePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	cleanPath := filepath.Clean(path)
	if !filepath.IsAbs(cleanPath) {
		var err error
		cleanPath, err = filepath.Abs(cleanPath)
		if err != nil {
			return "", fmt.Errorf("resolve path %q: %w", path, err)
		}
	}

	return cleanPath, nil
}

// splitLines splits file data into lines, handling different line endings.
// It preserves the original line ending style by detecting it from the data.
func splitLines(data []byte) []string {
	if len(data) == 0 {
		return []string{""}
	}

	// Convert to string (assuming UTF-8)
	content := string(data)

	// Normalize line endings to \n for splitting
	// We'll detect the original style separately
	normalized := strings.ReplaceAll(content, "\r\n", "\n")
	normalized = strings.ReplaceAll(normalized, "\r", "\n")

	// Split by \n
	lines := strings.Split(normalized, "\n")

	// Handle case where file ends with newline (remove empty last line)
	// But preserve it if it's the only content
	if len(lines) > 1 && lines[len(lines)-1] == "" && strings.HasSuffix(content, "\n") {
		// File ends with newline, keep the empty line
		return lines
	}

	// If last line is empty and file doesn't end with newline, it's likely
	// just the split artifact, but we'll keep it for consistency
	return lines
}

// detectLineEnding detects the line ending style of a file by reading a sample.
func detectLineEnding(path string) LineEnding {
	data, err := os.ReadFile(path)
	if err != nil {
		return LineEndingUnknown
	}

	content := string(data)
	if strings.Contains(content, "\r\n") {
		return LineEndingCRLF
	}
	if strings.Contains(content, "\r") && !strings.Contains(content, "\n") {
		return LineEndingCR
	}
	if strings.Contains(content, "\n") {
		return LineEndingLF
	}

	// No line endings found (single line file)
	return LineEndingLF // Default to LF
}
