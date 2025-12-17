package history

import (
	"testing"

	"github.com/AndrewDonelson/ted/core/buffer"
)

func TestNewHistory(t *testing.T) {
	h := NewHistory(50)
	if h == nil {
		t.Fatal("NewHistory() returned nil")
	}
	if h.maxDepth != 50 {
		t.Errorf("NewHistory() maxDepth = %d, want 50", h.maxDepth)
	}
	if h.CanUndo() {
		t.Error("NewHistory() CanUndo() = true, want false")
	}
	if h.CanRedo() {
		t.Error("NewHistory() CanRedo() = true, want false")
	}
}

func TestNewHistory_DefaultDepth(t *testing.T) {
	h := NewHistory(0)
	if h.maxDepth != 100 {
		t.Errorf("NewHistory(0) maxDepth = %d, want 100", h.maxDepth)
	}
}

func TestHistory_Push(t *testing.T) {
	h := NewHistory(10)

	op := &InsertOperation{
		Pos:  buffer.Position{Line: 0, Col: 0},
		Text: "test",
	}

	h.Push(op)

	if !h.CanUndo() {
		t.Error("After Push(), CanUndo() = false, want true")
	}
	if h.CanRedo() {
		t.Error("After Push(), CanRedo() = true, want false")
	}
	if h.Depth() != 1 {
		t.Errorf("After Push(), Depth() = %d, want 1", h.Depth())
	}
}

func TestHistory_Undo(t *testing.T) {
	h := NewHistory(10)
	buf := buffer.NewBuffer()
	buf.SetLines([]string{"hello", "world"})

	// Insert text
	op := &InsertOperation{
		Pos:  buffer.Position{Line: 0, Col: 5},
		Text: " there",
	}
	buf.Insert(op.Pos, op.Text)
	h.Push(op)

	// Undo
	if err := h.Undo(buf); err != nil {
		t.Fatalf("Undo() error = %v", err)
	}

	// Verify text was removed
	line, _ := buf.GetLine(0)
	if line != "hello" {
		t.Errorf("After Undo(), line = %q, want %q", line, "hello")
	}

	if h.CanUndo() {
		t.Error("After Undo(), CanUndo() = true, want false")
	}
	if !h.CanRedo() {
		t.Error("After Undo(), CanRedo() = false, want true")
	}
}

func TestHistory_Redo(t *testing.T) {
	h := NewHistory(10)
	buf := buffer.NewBuffer()
	buf.SetLines([]string{"hello", "world"})

	// Insert text
	op := &InsertOperation{
		Pos:  buffer.Position{Line: 0, Col: 5},
		Text: " there",
	}
	buf.Insert(op.Pos, op.Text)
	h.Push(op)

	// Undo then redo
	h.Undo(buf)
	if err := h.Redo(buf); err != nil {
		t.Fatalf("Redo() error = %v", err)
	}

	// Verify text was restored
	line, _ := buf.GetLine(0)
	if line != "hello there" {
		t.Errorf("After Redo(), line = %q, want %q", line, "hello there")
	}

	if !h.CanUndo() {
		t.Error("After Redo(), CanUndo() = false, want true")
	}
	if h.CanRedo() {
		t.Error("After Redo(), CanRedo() = true, want false")
	}
}

func TestHistory_Undo_NoOperations(t *testing.T) {
	h := NewHistory(10)
	buf := buffer.NewBuffer()

	if err := h.Undo(buf); err != ErrNoUndo {
		t.Errorf("Undo() with no operations error = %v, want ErrNoUndo", err)
	}
}

func TestHistory_Redo_NoOperations(t *testing.T) {
	h := NewHistory(10)
	buf := buffer.NewBuffer()

	if err := h.Redo(buf); err != ErrNoRedo {
		t.Errorf("Redo() with no operations error = %v, want ErrNoRedo", err)
	}
}

func TestHistory_MaxDepth(t *testing.T) {
	h := NewHistory(3)
	buf := buffer.NewBuffer()

	// Push more than maxDepth operations
	for i := 0; i < 5; i++ {
		op := &InsertOperation{
			Pos:  buffer.Position{Line: 0, Col: i},
			Text: "x",
		}
		buf.Insert(op.Pos, op.Text)
		h.Push(op)
	}

	// Should only keep maxDepth operations
	if h.Depth() > 3 {
		t.Errorf("History depth = %d, want <= 3", h.Depth())
	}
}

func TestHistory_ClearRedoOnPush(t *testing.T) {
	h := NewHistory(10)
	buf := buffer.NewBuffer()

	// Insert and undo
	op1 := &InsertOperation{Pos: buffer.Position{Line: 0, Col: 0}, Text: "a"}
	buf.Insert(op1.Pos, op1.Text)
	h.Push(op1)
	h.Undo(buf)

	if !h.CanRedo() {
		t.Fatal("Expected CanRedo() = true after undo")
	}

	// Push new operation - should clear redo
	op2 := &InsertOperation{Pos: buffer.Position{Line: 0, Col: 0}, Text: "b"}
	buf.Insert(op2.Pos, op2.Text)
	h.Push(op2)

	if h.CanRedo() {
		t.Error("After Push(), CanRedo() = true, want false (redo stack cleared)")
	}
}

func TestDeleteOperation_UndoRedo(t *testing.T) {
	h := NewHistory(10)
	buf := buffer.NewBuffer()
	buf.SetLines([]string{"hello", "world"})

	// Delete text
	start := buffer.Position{Line: 0, Col: 2}
	end := buffer.Position{Line: 0, Col: 4}
	deleted, _ := buf.GetText(start, end)
	op := &DeleteOperation{
		StartPos: start,
		EndPos:   end,
		Deleted:  deleted,
	}
	buf.Delete(start, end)
	h.Push(op)

	// Undo
	if err := h.Undo(buf); err != nil {
		t.Fatalf("Undo() error = %v", err)
	}

	line, _ := buf.GetLine(0)
	if line != "hello" {
		t.Errorf("After Undo(), line = %q, want %q", line, "hello")
	}

	// Redo
	if err := h.Redo(buf); err != nil {
		t.Fatalf("Redo() error = %v", err)
	}

	line, _ = buf.GetLine(0)
	if line != "heo" {
		t.Errorf("After Redo(), line = %q, want %q", line, "heo")
	}
}
