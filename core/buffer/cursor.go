// Package buffer implements cursor movement operations.
package buffer

// MoveCursorLeft moves the cursor one character to the left.
// If at the start of a line, moves to the end of the previous line.
func (b *Buffer) MoveCursorLeft() {
	pos := b.cursor

	if pos.Col > 0 {
		// Move left within the same line
		pos.Col--
	} else if pos.Line > 0 {
		// Move to end of previous line
		pos.Line--
		pos.Col = len(b.lines[pos.Line])
	}

	b.MoveCursor(pos)
}

// MoveCursorRight moves the cursor one character to the right.
// If at the end of a line, moves to the start of the next line.
func (b *Buffer) MoveCursorRight() {
	pos := b.cursor
	currentLineLen := len(b.lines[pos.Line])

	if pos.Col < currentLineLen {
		// Move right within the same line
		pos.Col++
	} else if pos.Line < len(b.lines)-1 {
		// Move to start of next line
		pos.Line++
		pos.Col = 0
	}

	b.MoveCursor(pos)
}

// MoveCursorUp moves the cursor one line up.
// The column position is preserved if possible, otherwise adjusted.
func (b *Buffer) MoveCursorUp() {
	pos := b.cursor

	if pos.Line > 0 {
		pos.Line--
		// Preserve column position if possible
		maxCol := len(b.lines[pos.Line])
		if pos.Col > maxCol {
			pos.Col = maxCol
		}
	}

	b.MoveCursor(pos)
}

// MoveCursorDown moves the cursor one line down.
// The column position is preserved if possible, otherwise adjusted.
func (b *Buffer) MoveCursorDown() {
	pos := b.cursor

	if pos.Line < len(b.lines)-1 {
		pos.Line++
		// Preserve column position if possible
		maxCol := len(b.lines[pos.Line])
		if pos.Col > maxCol {
			pos.Col = maxCol
		}
	}

	b.MoveCursor(pos)
}

// MoveCursorToLineStart moves the cursor to the start of the current line.
func (b *Buffer) MoveCursorToLineStart() {
	pos := b.cursor
	pos.Col = 0
	b.MoveCursor(pos)
}

// MoveCursorToLineEnd moves the cursor to the end of the current line.
func (b *Buffer) MoveCursorToLineEnd() {
	pos := b.cursor
	pos.Col = len(b.lines[pos.Line])
	b.MoveCursor(pos)
}

// MoveCursorToDocumentStart moves the cursor to the start of the document.
func (b *Buffer) MoveCursorToDocumentStart() {
	b.MoveCursor(Position{Line: 0, Col: 0})
}

// MoveCursorToDocumentEnd moves the cursor to the end of the document.
func (b *Buffer) MoveCursorToDocumentEnd() {
	if len(b.lines) == 0 {
		b.MoveCursor(Position{Line: 0, Col: 0})
		return
	}

	lastLine := len(b.lines) - 1
	lastCol := len(b.lines[lastLine])
	b.MoveCursor(Position{Line: lastLine, Col: lastCol})
}
