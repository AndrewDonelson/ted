package history

import (
	"github.com/AndrewDonelson/ted/core/buffer"
)

// InsertOperation represents an insert operation that can be undone.
type InsertOperation struct {
	Pos  buffer.Position
	Text string
}

// Undo removes the inserted text.
func (op *InsertOperation) Undo(buf *buffer.Buffer) error {
	// Calculate end position after insertion
	// This matches how Insert() calculates the cursor position
	endPos := op.Pos
	lines := splitLines(op.Text)
	if len(lines) == 1 {
		// Single line insert - cursor moves right by text length
		endPos.Col += len(lines[0])
	} else {
		// Multi-line insert - cursor moves to end of last inserted line
		endPos.Line += len(lines) - 1
		endPos.Col = len(lines[len(lines)-1])
	}

	// Delete the inserted text (from start to end)
	return buf.Delete(op.Pos, endPos)
}

// Redo reinserts the text.
func (op *InsertOperation) Redo(buf *buffer.Buffer) error {
	return buf.Insert(op.Pos, op.Text)
}

// Description returns a description of the operation.
func (op *InsertOperation) Description() string {
	if len(op.Text) == 1 {
		return "insert character"
	}
	return "insert text"
}

// DeleteOperation represents a delete operation that can be undone.
type DeleteOperation struct {
	StartPos buffer.Position
	EndPos   buffer.Position
	Deleted  string // The text that was deleted
}

// Undo restores the deleted text.
func (op *DeleteOperation) Undo(buf *buffer.Buffer) error {
	return buf.Insert(op.StartPos, op.Deleted)
}

// Redo deletes the text again.
func (op *DeleteOperation) Redo(buf *buffer.Buffer) error {
	return buf.Delete(op.StartPos, op.EndPos)
}

// Description returns a description of the operation.
func (op *DeleteOperation) Description() string {
	if op.StartPos.Line == op.EndPos.Line && op.EndPos.Col-op.StartPos.Col == 1 {
		return "delete character"
	}
	return "delete text"
}

// SetLinesOperation represents a SetLines operation (used for bulk changes).
type SetLinesOperation struct {
	OldLines []string
	NewLines []string
}

// Undo restores the old lines.
func (op *SetLinesOperation) Undo(buf *buffer.Buffer) error {
	buf.SetLines(op.OldLines)
	return nil
}

// Redo applies the new lines.
func (op *SetLinesOperation) Redo(buf *buffer.Buffer) error {
	buf.SetLines(op.NewLines)
	return nil
}

// Description returns a description of the operation.
func (op *SetLinesOperation) Description() string {
	return "set lines"
}

// splitLines splits text by newlines, similar to strings.Split but preserves empty lines.
func splitLines(text string) []string {
	if text == "" {
		return []string{""}
	}
	lines := []string{}
	current := ""
	for _, r := range text {
		if r == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(r)
		}
	}
	lines = append(lines, current)
	return lines
}
