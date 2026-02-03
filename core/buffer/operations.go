// Package buffer implements line operations and advanced navigation.
package buffer

import (
	"strings"
	"unicode"
)

// DeleteLine deletes the current line and returns its content.
// The cursor moves to the start of the next line, or the previous line if deleting the last line.
// Returns the deleted line content and any error.
func (b *Buffer) DeleteLine() (string, error) {
	if len(b.lines) == 0 {
		return "", nil
	}

	lineNum := b.cursor.Line
	if lineNum < 0 || lineNum >= len(b.lines) {
		return "", nil
	}

	deletedLine := b.lines[lineNum]

	// Remove the line
	newLines := make([]string, 0, len(b.lines)-1)
	newLines = append(newLines, b.lines[:lineNum]...)
	if lineNum+1 < len(b.lines) {
		newLines = append(newLines, b.lines[lineNum+1:]...)
	}

	// Ensure we have at least one line
	if len(newLines) == 0 {
		newLines = []string{""}
	}

	b.lines = newLines

	// Adjust cursor position
	if lineNum >= len(b.lines) {
		// Deleted last line, move to new last line
		b.cursor.Line = len(b.lines) - 1
		b.cursor.Col = 0
	} else {
		// Stay on same line number (which is now the next line)
		b.cursor.Line = lineNum
		b.cursor.Col = 0
	}

	b.modified = true
	return deletedLine, nil
}

// DuplicateLine creates a copy of the current line below it.
// The cursor moves to the duplicated line at the same column position.
func (b *Buffer) DuplicateLine() error {
	if len(b.lines) == 0 {
		return nil
	}

	lineNum := b.cursor.Line
	if lineNum < 0 || lineNum >= len(b.lines) {
		return nil
	}

	line := b.lines[lineNum]

	// Insert copy of line after current line
	newLines := make([]string, 0, len(b.lines)+1)
	newLines = append(newLines, b.lines[:lineNum+1]...)
	newLines = append(newLines, line)
	if lineNum+1 < len(b.lines) {
		newLines = append(newLines, b.lines[lineNum+1:]...)
	}

	b.lines = newLines

	// Move cursor to the duplicated line
	b.cursor.Line = lineNum + 1
	// Keep column position
	if b.cursor.Col > len(line) {
		b.cursor.Col = len(line)
	}

	b.modified = true
	return nil
}

// MoveLineUp swaps the current line with the one above it.
// The cursor moves with the line.
func (b *Buffer) MoveLineUp() error {
	if len(b.lines) < 2 {
		return nil
	}

	lineNum := b.cursor.Line
	if lineNum <= 0 {
		// Already at top, can't move up
		return nil
	}

	// Swap current line with line above
	b.lines[lineNum], b.lines[lineNum-1] = b.lines[lineNum-1], b.lines[lineNum]

	// Move cursor up with the line
	b.cursor.Line = lineNum - 1

	b.modified = true
	return nil
}

// MoveLineDown swaps the current line with the one below it.
// The cursor moves with the line.
func (b *Buffer) MoveLineDown() error {
	if len(b.lines) < 2 {
		return nil
	}

	lineNum := b.cursor.Line
	if lineNum >= len(b.lines)-1 {
		// Already at bottom, can't move down
		return nil
	}

	// Swap current line with line below
	b.lines[lineNum], b.lines[lineNum+1] = b.lines[lineNum+1], b.lines[lineNum]

	// Move cursor down with the line
	b.cursor.Line = lineNum + 1

	b.modified = true
	return nil
}

// InsertLineAbove inserts a new empty line above the current line.
// The cursor moves to the start of the new line.
func (b *Buffer) InsertLineAbove() error {
	lineNum := b.cursor.Line

	// Insert empty line above
	newLines := make([]string, 0, len(b.lines)+1)
	newLines = append(newLines, b.lines[:lineNum]...)
	newLines = append(newLines, "")
	newLines = append(newLines, b.lines[lineNum:]...)

	b.lines = newLines

	// Move cursor to the new line
	b.cursor.Line = lineNum
	b.cursor.Col = 0

	b.modified = true
	return nil
}

// InsertLineBelow inserts a new empty line below the current line.
// The cursor moves to the start of the new line.
func (b *Buffer) InsertLineBelow() error {
	lineNum := b.cursor.Line

	// Insert empty line below
	newLines := make([]string, 0, len(b.lines)+1)
	newLines = append(newLines, b.lines[:lineNum+1]...)
	newLines = append(newLines, "")
	if lineNum+1 < len(b.lines) {
		newLines = append(newLines, b.lines[lineNum+1:]...)
	}

	b.lines = newLines

	// Move cursor to the new line
	b.cursor.Line = lineNum + 1
	b.cursor.Col = 0

	b.modified = true
	return nil
}

// isWordChar returns true if the rune is a word character (alphanumeric or underscore).
func isWordChar(r byte) bool {
	return unicode.IsLetter(rune(r)) || unicode.IsDigit(rune(r)) || r == '_'
}

// MoveCursorWordLeft moves the cursor to the start of the previous word.
// A word is a sequence of word characters (alphanumeric + underscore).
func (b *Buffer) MoveCursorWordLeft() {
	pos := b.cursor
	line := b.lines[pos.Line]

	// If at the start of a line, move to end of previous line
	if pos.Col == 0 {
		if pos.Line > 0 {
			pos.Line--
			pos.Col = len(b.lines[pos.Line])
			b.MoveCursor(pos)
		}
		return
	}

	// Check if we're currently on a word character
	onWord := isWordChar(line[pos.Col-1])

	if onWord {
		// We're in the middle of a word, skip to start of current word
		for pos.Col > 0 && isWordChar(line[pos.Col-1]) {
			pos.Col--
		}
	} else {
		// We're on non-word chars (spaces/punctuation), skip them
		for pos.Col > 0 && !isWordChar(line[pos.Col-1]) {
			pos.Col--
		}
		// Then skip the word we land on
		for pos.Col > 0 && isWordChar(line[pos.Col-1]) {
			pos.Col--
		}
	}

	b.MoveCursor(pos)
}

// MoveCursorWordRight moves the cursor to the start of the next word.
// A word is a sequence of word characters (alphanumeric + underscore).
func (b *Buffer) MoveCursorWordRight() {
	pos := b.cursor
	line := b.lines[pos.Line]

	// If at the end of a line, move to start of next line
	if pos.Col >= len(line) {
		if pos.Line < len(b.lines)-1 {
			pos.Line++
			pos.Col = 0
			b.MoveCursor(pos)
		}
		return
	}

	// Check if we're currently on a word character
	onWord := isWordChar(line[pos.Col])

	if onWord {
		// We're in the middle of a word, skip to end of current word
		for pos.Col < len(line) && isWordChar(line[pos.Col]) {
			pos.Col++
		}
	}

	// Skip non-word characters (whitespace, punctuation)
	for pos.Col < len(line) && !isWordChar(line[pos.Col]) {
		pos.Col++
	}

	b.MoveCursor(pos)
}

// MoveCursorPageUp moves the cursor up by the specified number of lines.
// Typically used with viewport height to scroll by page.
func (b *Buffer) MoveCursorPageUp(pageSize int) {
	if pageSize <= 0 {
		pageSize = 10 // Default page size
	}

	pos := b.cursor
	pos.Line -= pageSize

	if pos.Line < 0 {
		pos.Line = 0
	}

	// Adjust column to fit new line
	if pos.Line < len(b.lines) {
		maxCol := len(b.lines[pos.Line])
		if pos.Col > maxCol {
			pos.Col = maxCol
		}
	}

	b.MoveCursor(pos)
}

// MoveCursorPageDown moves the cursor down by the specified number of lines.
// Typically used with viewport height to scroll by page.
func (b *Buffer) MoveCursorPageDown(pageSize int) {
	if pageSize <= 0 {
		pageSize = 10 // Default page size
	}

	pos := b.cursor
	pos.Line += pageSize

	if pos.Line >= len(b.lines) {
		pos.Line = len(b.lines) - 1
	}

	// Adjust column to fit new line
	if pos.Line < len(b.lines) {
		maxCol := len(b.lines[pos.Line])
		if pos.Col > maxCol {
			pos.Col = maxCol
		}
	}

	b.MoveCursor(pos)
}

// GetCurrentLineIndentation returns the leading whitespace of the current line.
// This is useful for auto-indentation when inserting new lines.
func (b *Buffer) GetCurrentLineIndentation() string {
	if b.cursor.Line < 0 || b.cursor.Line >= len(b.lines) {
		return ""
	}

	line := b.lines[b.cursor.Line]
	var indent strings.Builder
	for i := 0; i < len(line); i++ {
		if line[i] == ' ' || line[i] == '\t' {
			indent.WriteByte(line[i])
		} else {
			break
		}
	}

	return indent.String()
}
