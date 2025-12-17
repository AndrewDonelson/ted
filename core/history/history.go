// Package history implements undo/redo functionality for the text editor.
//
// It maintains a history of operations that can be undone and redone.
// The history uses a command pattern where each operation can be reversed.
package history

import (
	"github.com/AndrewDonelson/ted/core/buffer"
)

// Operation represents a single undoable operation.
type Operation interface {
	// Undo reverses the operation, returning the buffer state after undo.
	Undo(buf *buffer.Buffer) error
	// Redo reapplies the operation, returning the buffer state after redo.
	Redo(buf *buffer.Buffer) error
	// Description returns a human-readable description of the operation.
	Description() string
}

// History manages undo/redo history for a buffer.
// It maintains separate undo and redo stacks.
type History struct {
	undoStack []Operation
	redoStack []Operation
	maxDepth  int // Maximum number of operations to keep
}

// NewHistory creates a new history manager with the specified maximum depth.
// If maxDepth is 0, a default of 100 is used.
func NewHistory(maxDepth int) *History {
	if maxDepth <= 0 {
		maxDepth = 100 // Default depth
	}
	return &History{
		undoStack: make([]Operation, 0, maxDepth),
		redoStack: make([]Operation, 0, maxDepth),
		maxDepth:  maxDepth,
	}
}

// Push adds a new operation to the undo stack.
// This clears the redo stack (new operation invalidates redo history).
func (h *History) Push(op Operation) {
	// Clear redo stack when new operation is pushed
	h.redoStack = h.redoStack[:0]

	// Add to undo stack
	h.undoStack = append(h.undoStack, op)

	// Limit stack size
	if len(h.undoStack) > h.maxDepth {
		// Remove oldest operation
		copy(h.undoStack, h.undoStack[1:])
		h.undoStack = h.undoStack[:len(h.undoStack)-1]
	}
}

// CanUndo returns whether there are operations that can be undone.
func (h *History) CanUndo() bool {
	return len(h.undoStack) > 0
}

// CanRedo returns whether there are operations that can be redone.
func (h *History) CanRedo() bool {
	return len(h.redoStack) > 0
}

// Undo undoes the last operation and moves it to the redo stack.
// Returns an error if there are no operations to undo.
func (h *History) Undo(buf *buffer.Buffer) error {
	if !h.CanUndo() {
		return ErrNoUndo
	}

	// Pop from undo stack
	op := h.undoStack[len(h.undoStack)-1]
	h.undoStack = h.undoStack[:len(h.undoStack)-1]

	// Undo the operation
	if err := op.Undo(buf); err != nil {
		// Put it back on the stack if undo failed
		h.undoStack = append(h.undoStack, op)
		return err
	}

	// Move to redo stack
	h.redoStack = append(h.redoStack, op)

	// Limit redo stack size
	if len(h.redoStack) > h.maxDepth {
		copy(h.redoStack, h.redoStack[1:])
		h.redoStack = h.redoStack[:len(h.redoStack)-1]
	}

	return nil
}

// Redo redoes the last undone operation and moves it back to the undo stack.
// Returns an error if there are no operations to redo.
func (h *History) Redo(buf *buffer.Buffer) error {
	if !h.CanRedo() {
		return ErrNoRedo
	}

	// Pop from redo stack
	op := h.redoStack[len(h.redoStack)-1]
	h.redoStack = h.redoStack[:len(h.redoStack)-1]

	// Redo the operation
	if err := op.Redo(buf); err != nil {
		// Put it back on the stack if redo failed
		h.redoStack = append(h.redoStack, op)
		return err
	}

	// Move back to undo stack
	h.undoStack = append(h.undoStack, op)

	// Limit undo stack size
	if len(h.undoStack) > h.maxDepth {
		copy(h.undoStack, h.undoStack[1:])
		h.undoStack = h.undoStack[:len(h.undoStack)-1]
	}

	return nil
}

// Clear clears all history.
func (h *History) Clear() {
	h.undoStack = h.undoStack[:0]
	h.redoStack = h.redoStack[:0]
}

// ClearRedo clears only the redo stack (used when saving).
func (h *History) ClearRedo() {
	h.redoStack = h.redoStack[:0]
}

// Depth returns the current depth of the undo stack.
func (h *History) Depth() int {
	return len(h.undoStack)
}

// Errors
var (
	ErrNoUndo = &HistoryError{msg: "no operations to undo"}
	ErrNoRedo = &HistoryError{msg: "no operations to redo"}
)

// HistoryError represents an error in history operations.
type HistoryError struct {
	msg string
}

func (e *HistoryError) Error() string {
	return e.msg
}
