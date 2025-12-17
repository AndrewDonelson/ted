// Package buffer implements a text buffer for terminal text editing.
//
// The buffer stores text as a slice of lines and provides operations
// for inserting, deleting, and querying text. It maintains cursor
// position and supports text selection.
package buffer

import (
	"fmt"
	"strings"
)

// Position represents a location in the buffer.
// Line and Col are zero-indexed. Col is a byte offset, not a rune offset.
type Position struct {
	Line int // Line number (0-indexed)
	Col  int // Column number (0-indexed, byte offset)
}

// Buffer represents an in-memory text buffer.
// It stores text as a slice of lines and provides methods for
// editing operations. Buffer is not safe for concurrent use.
type Buffer struct {
	lines    []string
	cursor   Position
	modified bool
}

// NewBuffer creates a new empty buffer.
func NewBuffer() *Buffer {
	return &Buffer{
		lines:    []string{""},
		cursor:   Position{Line: 0, Col: 0},
		modified: false,
	}
}

// Insert inserts text at the specified position.
// If text contains newlines, it will be split across multiple lines.
// Returns an error if the position is invalid.
//
// Example:
//
//	err := buf.Insert(Position{Line: 0, Col: 5}, "world")
func (b *Buffer) Insert(pos Position, text string) error {
	if err := b.validatePosition(pos); err != nil {
		return err
	}

	if text == "" {
		return nil // No-op, don't mark as modified
	}

	// Split text by newlines
	lines := strings.Split(text, "\n")

	if len(lines) == 1 {
		// Single line insert
		line := b.lines[pos.Line]
		before := line[:pos.Col]
		after := line[pos.Col:]
		b.lines[pos.Line] = before + lines[0] + after
		b.cursor = Position{Line: pos.Line, Col: pos.Col + len(lines[0])}
	} else {
		// Multi-line insert
		line := b.lines[pos.Line]
		before := line[:pos.Col]
		after := line[pos.Col:]

		// Build new lines slice
		newLines := make([]string, 0, len(b.lines)+len(lines)-1)

		// Lines before the insertion point
		newLines = append(newLines, b.lines[:pos.Line]...)

		// First line: merge before + first inserted line
		newLines = append(newLines, before+lines[0])

		// Middle lines: insert as new lines
		newLines = append(newLines, lines[1:len(lines)-1]...)

		// Last line: merge last inserted line + after
		lastLine := lines[len(lines)-1] + after
		newLines = append(newLines, lastLine)

		// Remaining lines after insertion point
		if pos.Line+1 < len(b.lines) {
			newLines = append(newLines, b.lines[pos.Line+1:]...)
		}

		b.lines = newLines
		b.cursor = Position{
			Line: pos.Line + len(lines) - 1,
			Col:  len(lines[len(lines)-1]),
		}
	}

	b.modified = true
	return nil
}

// Delete deletes text between start and end positions (inclusive start, exclusive end).
// Returns an error if either position is invalid.
func (b *Buffer) Delete(start, end Position) error {
	if err := b.validatePosition(start); err != nil {
		return err
	}
	if err := b.validatePosition(end); err != nil {
		return err
	}

	if start.Line > end.Line || (start.Line == end.Line && start.Col > end.Col) {
		return fmt.Errorf("invalid delete range: start position after end position")
	}

	if start.Line == end.Line && start.Col == end.Col {
		// No-op delete - don't modify or move cursor
		return nil
	}

	if start.Line == end.Line {
		// Single line delete
		line := b.lines[start.Line]
		newLine := line[:start.Col] + line[end.Col:]
		b.lines[start.Line] = newLine

		// If line becomes empty and we deleted from start, remove the line
		// (unless it's the only line in the buffer)
		if newLine == "" && start.Col == 0 && len(b.lines) > 1 {
			newLines := make([]string, 0, len(b.lines)-1)
			newLines = append(newLines, b.lines[:start.Line]...)
			if start.Line+1 < len(b.lines) {
				newLines = append(newLines, b.lines[start.Line+1:]...)
			}
			b.lines = newLines
			// Adjust cursor if we removed a line before cursor
			if b.cursor.Line > start.Line {
				b.cursor.Line--
			} else if b.cursor.Line == start.Line {
				b.cursor.Line = start.Line
				if b.cursor.Line >= len(b.lines) {
					b.cursor.Line = len(b.lines) - 1
				}
				if b.cursor.Line < 0 {
					b.cursor.Line = 0
				}
				b.cursor.Col = 0
			}
		}

		b.cursor = start
		if b.cursor.Line >= len(b.lines) {
			b.cursor.Line = len(b.lines) - 1
		}
		if b.cursor.Line < 0 {
			b.cursor.Line = 0
		}
		if b.cursor.Line >= 0 && b.cursor.Col > len(b.lines[b.cursor.Line]) {
			b.cursor.Col = len(b.lines[b.cursor.Line])
		}
		b.modified = true
	} else {
		// Multi-line delete
		startLine := b.lines[start.Line]
		endLine := b.lines[end.Line]

		// Merge start and end lines
		newLine := startLine[:start.Col] + endLine[end.Col:]

		// Build new lines slice
		newLines := make([]string, 0, len(b.lines))
		// Lines before deletion
		if start.Line > 0 {
			newLines = append(newLines, b.lines[:start.Line]...)
		}
		// Merged line (only if it's not empty, or if it's the only line)
		if newLine != "" || len(b.lines) == 1 {
			newLines = append(newLines, newLine)
		}
		// Lines after deletion
		if end.Line+1 < len(b.lines) {
			newLines = append(newLines, b.lines[end.Line+1:]...)
		}

		// Ensure we have at least one line
		if len(newLines) == 0 {
			newLines = []string{""}
		}

		b.lines = newLines
		b.cursor = start
		// Adjust cursor if we removed lines
		if b.cursor.Line >= len(b.lines) {
			b.cursor.Line = len(b.lines) - 1
		}
		if b.cursor.Line < 0 {
			b.cursor.Line = 0
		}
		if b.cursor.Line >= 0 && b.cursor.Col > len(b.lines[b.cursor.Line]) {
			b.cursor.Col = len(b.lines[b.cursor.Line])
		}
		b.modified = true
	}

	return nil
}

// GetLine returns the text at the specified line number.
// Returns an error if the line number is invalid.
func (b *Buffer) GetLine(lineNum int) (string, error) {
	if lineNum < 0 || lineNum >= len(b.lines) {
		return "", fmt.Errorf("invalid line number: %d", lineNum)
	}
	return b.lines[lineNum], nil
}

// LineCount returns the total number of lines in the buffer.
func (b *Buffer) LineCount() int {
	return len(b.lines)
}

// GetCursor returns the current cursor position.
func (b *Buffer) GetCursor() Position {
	return b.cursor
}

// MoveCursor moves the cursor to the specified position.
// The position is validated and adjusted if necessary.
func (b *Buffer) MoveCursor(pos Position) {
	if pos.Line < 0 {
		pos.Line = 0
	}
	if pos.Line >= len(b.lines) {
		pos.Line = len(b.lines) - 1
	}
	if pos.Line < 0 {
		// Empty buffer
		pos.Line = 0
		pos.Col = 0
		b.cursor = pos
		return
	}

	maxCol := len(b.lines[pos.Line])
	if pos.Col < 0 {
		pos.Col = 0
	}
	if pos.Col > maxCol {
		pos.Col = maxCol
	}

	b.cursor = pos
}

// IsModified returns whether the buffer has been modified since the last save.
func (b *Buffer) IsModified() bool {
	return b.modified
}

// MarkSaved marks the buffer as saved (not modified).
func (b *Buffer) MarkSaved() {
	b.modified = false
}

// SetLines sets the buffer content from a slice of lines.
// This is primarily used for loading files.
func (b *Buffer) SetLines(lines []string) {
	if len(lines) == 0 {
		b.lines = []string{""}
	} else {
		b.lines = lines
	}
	b.cursor = Position{Line: 0, Col: 0}
	b.modified = false
}

// GetAllLines returns all lines in the buffer as a slice.
func (b *Buffer) GetAllLines() []string {
	lines := make([]string, len(b.lines))
	copy(lines, b.lines)
	return lines
}

// GetText returns the text between start and end positions (inclusive start, exclusive end).
// This is useful for recording what was deleted for undo operations.
func (b *Buffer) GetText(start, end Position) (string, error) {
	if err := b.validatePosition(start); err != nil {
		return "", err
	}
	if err := b.validatePosition(end); err != nil {
		return "", err
	}

	if start.Line > end.Line || (start.Line == end.Line && start.Col > end.Col) {
		return "", fmt.Errorf("invalid range: start position after end position")
	}

	if start.Line == end.Line && start.Col == end.Col {
		return "", nil // Empty range
	}

	if start.Line == end.Line {
		// Single line
		line := b.lines[start.Line]
		return line[start.Col:end.Col], nil
	}

	// Multi-line
	var result strings.Builder
	// First line: from start.Col to end of line
	result.WriteString(b.lines[start.Line][start.Col:])
	result.WriteString("\n")
	// Middle lines: full lines
	for line := start.Line + 1; line < end.Line; line++ {
		result.WriteString(b.lines[line])
		result.WriteString("\n")
	}
	// Last line: from start to end.Col
	result.WriteString(b.lines[end.Line][:end.Col])
	return result.String(), nil
}

// validatePosition checks if a position is valid for the current buffer state.
func (b *Buffer) validatePosition(pos Position) error {
	if pos.Line < 0 || pos.Line >= len(b.lines) {
		return fmt.Errorf("invalid line number: %d", pos.Line)
	}

	maxCol := len(b.lines[pos.Line])
	if pos.Col < 0 || pos.Col > maxCol {
		return fmt.Errorf("invalid column number: %d (max: %d)", pos.Col, maxCol)
	}

	return nil
}
