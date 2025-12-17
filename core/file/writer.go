// Package file implements file writing with atomic operations.
package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WriteFile writes lines to a file atomically (using temp file + rename).
// It preserves the specified line ending style.
// Returns an error if the file cannot be written.
//
// Example:
//
//	lines := []string{"line1", "line2", "line3"}
//	err := WriteFile("example.txt", lines, LineEndingLF)
func WriteFile(path string, lines []string, lineEnding LineEnding) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	// Validate path
	cleanPath, err := validatePath(path)
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(cleanPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create directory %q: %w", dir, err)
	}

	// Convert line ending to string
	ending := lineEndingToString(lineEnding)

	// Build file content
	var content strings.Builder
	for i, line := range lines {
		content.WriteString(line)
		// Add line ending after each line except the last
		if i < len(lines)-1 {
			content.WriteString(ending)
		}
		// If last line is empty, it represents a trailing newline that was already
		// added after the previous line, so we don't add another one
	}

	// Atomic write: write to temp file, then rename
	return atomicWrite(cleanPath, []byte(content.String()))
}

// WriteFilePreserveEnding writes lines to a file, preserving the original line ending.
// If the file doesn't exist, it defaults to LF.
func WriteFilePreserveEnding(path string, lines []string) error {
	var lineEnding LineEnding = LineEndingLF

	// Try to detect existing line ending
	if _, err := os.Stat(path); err == nil {
		lineEnding = detectLineEnding(path)
	}

	return WriteFile(path, lines, lineEnding)
}

// atomicWrite writes data to a file atomically using a temporary file and rename.
// This ensures the file is either completely written or not written at all.
func atomicWrite(path string, data []byte) error {
	// Create temp file in same directory
	dir := filepath.Dir(path)
	tmpFile, err := os.CreateTemp(dir, filepath.Base(path)+".tmp.*")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	// Write data to temp file
	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("write temp file: %w", err)
	}

	// Sync to ensure data is written to disk
	if err := tmpFile.Sync(); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("sync temp file: %w", err)
	}

	// Close temp file
	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("close temp file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("rename temp file to %q: %w", path, err)
	}

	return nil
}

// lineEndingToString converts a LineEnding to its string representation.
func lineEndingToString(ending LineEnding) string {
	switch ending {
	case LineEndingCRLF:
		return "\r\n"
	case LineEndingCR:
		return "\r"
	case LineEndingLF:
		return "\n"
	default:
		return "\n" // Default to LF
	}
}
