package buffer

import (
	"testing"
)

func TestDeleteLine(t *testing.T) {
	tests := []struct {
		name        string
		initial     []string
		cursorLine  int
		cursorCol   int
		wantLines   []string
		wantDeleted string
		wantCursor  Position
	}{
		{
			name:        "delete single line",
			initial:     []string{"hello", "world"},
			cursorLine:  0,
			cursorCol:   2,
			wantLines:   []string{"world"},
			wantDeleted: "hello",
			wantCursor:  Position{Line: 0, Col: 0},
		},
		{
			name:        "delete last line",
			initial:     []string{"hello", "world"},
			cursorLine:  1,
			cursorCol:   2,
			wantLines:   []string{"hello"},
			wantDeleted: "world",
			wantCursor:  Position{Line: 0, Col: 0},
		},
		{
			name:        "delete from middle",
			initial:     []string{"line1", "line2", "line3"},
			cursorLine:  1,
			cursorCol:   2,
			wantLines:   []string{"line1", "line3"},
			wantDeleted: "line2",
			wantCursor:  Position{Line: 1, Col: 0},
		},
		{
			name:        "delete only line leaves empty line",
			initial:     []string{"only"},
			cursorLine:  0,
			cursorCol:   0,
			wantLines:   []string{""},
			wantDeleted: "only",
			wantCursor:  Position{Line: 0, Col: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBuffer()
			b.SetLines(tt.initial)
			b.MoveCursor(Position{Line: tt.cursorLine, Col: tt.cursorCol})

			deleted, _ := b.DeleteLine()

			if deleted != tt.wantDeleted {
				t.Errorf("DeleteLine() deleted = %q, want %q", deleted, tt.wantDeleted)
			}

			gotLines := b.GetAllLines()
			if !slicesEqual(gotLines, tt.wantLines) {
				t.Errorf("DeleteLine() lines = %v, want %v", gotLines, tt.wantLines)
			}

			gotCursor := b.GetCursor()
			if gotCursor != tt.wantCursor {
				t.Errorf("DeleteLine() cursor = %v, want %v", gotCursor, tt.wantCursor)
			}
		})
	}
}

func TestDuplicateLine(t *testing.T) {
	tests := []struct {
		name       string
		initial    []string
		cursorLine int
		cursorCol  int
		wantLines  []string
		wantCursor Position
	}{
		{
			name:       "duplicate single line",
			initial:    []string{"hello"},
			cursorLine: 0,
			cursorCol:  2,
			wantLines:  []string{"hello", "hello"},
			wantCursor: Position{Line: 1, Col: 2},
		},
		{
			name:       "duplicate first of multiple",
			initial:    []string{"line1", "line2"},
			cursorLine: 0,
			cursorCol:  2,
			wantLines:  []string{"line1", "line1", "line2"},
			wantCursor: Position{Line: 1, Col: 2},
		},
		{
			name:       "duplicate last line",
			initial:    []string{"line1", "line2"},
			cursorLine: 1,
			cursorCol:  2,
			wantLines:  []string{"line1", "line2", "line2"},
			wantCursor: Position{Line: 2, Col: 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBuffer()
			b.SetLines(tt.initial)
			b.MoveCursor(Position{Line: tt.cursorLine, Col: tt.cursorCol})

			b.DuplicateLine()

			gotLines := b.GetAllLines()
			if !slicesEqual(gotLines, tt.wantLines) {
				t.Errorf("DuplicateLine() lines = %v, want %v", gotLines, tt.wantLines)
			}

			gotCursor := b.GetCursor()
			if gotCursor != tt.wantCursor {
				t.Errorf("DuplicateLine() cursor = %v, want %v", gotCursor, tt.wantCursor)
			}
		})
	}
}

func TestMoveLineUp(t *testing.T) {
	tests := []struct {
		name       string
		initial    []string
		cursorLine int
		wantLines  []string
		wantCursor Position
	}{
		{
			name:       "move line up from second line",
			initial:    []string{"line1", "line2"},
			cursorLine: 1,
			wantLines:  []string{"line2", "line1"},
			wantCursor: Position{Line: 0, Col: 0},
		},
		{
			name:       "move line up from middle",
			initial:    []string{"line1", "line2", "line3"},
			cursorLine: 1,
			wantLines:  []string{"line2", "line1", "line3"},
			wantCursor: Position{Line: 0, Col: 0},
		},
		{
			name:       "move up from first line does nothing",
			initial:    []string{"line1", "line2"},
			cursorLine: 0,
			wantLines:  []string{"line1", "line2"},
			wantCursor: Position{Line: 0, Col: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBuffer()
			b.SetLines(tt.initial)
			b.MoveCursor(Position{Line: tt.cursorLine, Col: 0})

			b.MoveLineUp()

			gotLines := b.GetAllLines()
			if !slicesEqual(gotLines, tt.wantLines) {
				t.Errorf("MoveLineUp() lines = %v, want %v", gotLines, tt.wantLines)
			}

			gotCursor := b.GetCursor()
			if gotCursor != tt.wantCursor {
				t.Errorf("MoveLineUp() cursor = %v, want %v", gotCursor, tt.wantCursor)
			}
		})
	}
}

func TestMoveLineDown(t *testing.T) {
	tests := []struct {
		name       string
		initial    []string
		cursorLine int
		wantLines  []string
		wantCursor Position
	}{
		{
			name:       "move line down from first line",
			initial:    []string{"line1", "line2"},
			cursorLine: 0,
			wantLines:  []string{"line2", "line1"},
			wantCursor: Position{Line: 1, Col: 0},
		},
		{
			name:       "move line down from second of three",
			initial:    []string{"line1", "line2", "line3"},
			cursorLine: 1,
			wantLines:  []string{"line1", "line3", "line2"},
			wantCursor: Position{Line: 2, Col: 0},
		},
		{
			name:       "move down from last line does nothing",
			initial:    []string{"line1", "line2"},
			cursorLine: 1,
			wantLines:  []string{"line1", "line2"},
			wantCursor: Position{Line: 1, Col: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBuffer()
			b.SetLines(tt.initial)
			b.MoveCursor(Position{Line: tt.cursorLine, Col: 0})

			b.MoveLineDown()

			gotLines := b.GetAllLines()
			if !slicesEqual(gotLines, tt.wantLines) {
				t.Errorf("MoveLineDown() lines = %v, want %v", gotLines, tt.wantLines)
			}

			gotCursor := b.GetCursor()
			if gotCursor != tt.wantCursor {
				t.Errorf("MoveLineDown() cursor = %v, want %v", gotCursor, tt.wantCursor)
			}
		})
	}
}

func TestInsertLineAbove(t *testing.T) {
	tests := []struct {
		name       string
		initial    []string
		cursorLine int
		cursorCol  int
		wantLines  []string
		wantCursor Position
	}{
		{
			name:       "insert above first line",
			initial:    []string{"hello"},
			cursorLine: 0,
			cursorCol:  3,
			wantLines:  []string{"", "hello"},
			wantCursor: Position{Line: 0, Col: 0},
		},
		{
			name:       "insert above second line",
			initial:    []string{"line1", "line2"},
			cursorLine: 1,
			cursorCol:  3,
			wantLines:  []string{"line1", "", "line2"},
			wantCursor: Position{Line: 1, Col: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBuffer()
			b.SetLines(tt.initial)
			b.MoveCursor(Position{Line: tt.cursorLine, Col: tt.cursorCol})

			b.InsertLineAbove()

			gotLines := b.GetAllLines()
			if !slicesEqual(gotLines, tt.wantLines) {
				t.Errorf("InsertLineAbove() lines = %v, want %v", gotLines, tt.wantLines)
			}

			gotCursor := b.GetCursor()
			if gotCursor != tt.wantCursor {
				t.Errorf("InsertLineAbove() cursor = %v, want %v", gotCursor, tt.wantCursor)
			}
		})
	}
}

func TestInsertLineBelow(t *testing.T) {
	tests := []struct {
		name       string
		initial    []string
		cursorLine int
		cursorCol  int
		wantLines  []string
		wantCursor Position
	}{
		{
			name:       "insert below single line",
			initial:    []string{"hello"},
			cursorLine: 0,
			cursorCol:  3,
			wantLines:  []string{"hello", ""},
			wantCursor: Position{Line: 1, Col: 0},
		},
		{
			name:       "insert below first of two",
			initial:    []string{"line1", "line2"},
			cursorLine: 0,
			cursorCol:  3,
			wantLines:  []string{"line1", "", "line2"},
			wantCursor: Position{Line: 1, Col: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBuffer()
			b.SetLines(tt.initial)
			b.MoveCursor(Position{Line: tt.cursorLine, Col: tt.cursorCol})

			b.InsertLineBelow()

			gotLines := b.GetAllLines()
			if !slicesEqual(gotLines, tt.wantLines) {
				t.Errorf("InsertLineBelow() lines = %v, want %v", gotLines, tt.wantLines)
			}

			gotCursor := b.GetCursor()
			if gotCursor != tt.wantCursor {
				t.Errorf("InsertLineBelow() cursor = %v, want %v", gotCursor, tt.wantCursor)
			}
		})
	}
}

func TestMoveCursorWordLeft(t *testing.T) {
	tests := []struct {
		name      string
		lines     []string
		startLine int
		startCol  int
		wantLine  int
		wantCol   int
	}{
		{
			name:      "move to start of current word",
			lines:     []string{"hello world"},
			startLine: 0,
			startCol:  8,
			wantLine:  0,
			wantCol:   6,
		},
		{
			name:      "move to previous word from word start",
			lines:     []string{"hello world"},
			startLine: 0,
			startCol:  6,
			wantLine:  0,
			wantCol:   0,
		},
		{
			name:      "at line start moves to end of previous line",
			lines:     []string{"hello", "world"},
			startLine: 1,
			startCol:  0,
			wantLine:  0,
			wantCol:   5,
		},
		{
			name:      "skip non-word chars to previous word",
			lines:     []string{"hello world"},
			startLine: 0,
			startCol:  6,
			wantLine:  0,
			wantCol:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBuffer()
			b.SetLines(tt.lines)
			b.MoveCursor(Position{Line: tt.startLine, Col: tt.startCol})

			b.MoveCursorWordLeft()

			gotCursor := b.GetCursor()
			if gotCursor.Line != tt.wantLine || gotCursor.Col != tt.wantCol {
				t.Errorf("MoveCursorWordLeft() cursor = {Line: %d, Col: %d}, want {Line: %d, Col: %d}",
					gotCursor.Line, gotCursor.Col, tt.wantLine, tt.wantCol)
			}
		})
	}
}

func TestMoveCursorWordRight(t *testing.T) {
	tests := []struct {
		name      string
		lines     []string
		startLine int
		startCol  int
		wantLine  int
		wantCol   int
	}{
		{
			name:      "move to start of next word from word start",
			lines:     []string{"hello world"},
			startLine: 0,
			startCol:  0,
			wantLine:  0,
			wantCol:   6,
		},
		{
			name:      "move from middle of word to next word",
			lines:     []string{"hello world"},
			startLine: 0,
			startCol:  2,
			wantLine:  0,
			wantCol:   6,
		},
		{
			name:      "at end of line moves to next line",
			lines:     []string{"hello", "world"},
			startLine: 0,
			startCol:  5,
			wantLine:  1,
			wantCol:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBuffer()
			b.SetLines(tt.lines)
			b.MoveCursor(Position{Line: tt.startLine, Col: tt.startCol})

			b.MoveCursorWordRight()

			gotCursor := b.GetCursor()
			if gotCursor.Line != tt.wantLine || gotCursor.Col != tt.wantCol {
				t.Errorf("MoveCursorWordRight() cursor = {Line: %d, Col: %d}, want {Line: %d, Col: %d}",
					gotCursor.Line, gotCursor.Col, tt.wantLine, tt.wantCol)
			}
		})
	}
}

func TestMoveCursorPageUp(t *testing.T) {
	tests := []struct {
		name      string
		lines     int
		startLine int
		pageSize  int
		wantLine  int
	}{
		{
			name:      "move up by page",
			lines:     20,
			startLine: 15,
			pageSize:  10,
			wantLine:  5,
		},
		{
			name:      "stops at top",
			lines:     20,
			startLine: 5,
			pageSize:  10,
			wantLine:  0,
		},
		{
			name:      "default page size",
			lines:     20,
			startLine: 15,
			pageSize:  0,
			wantLine:  5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBuffer()
			lines := make([]string, tt.lines)
			for i := range lines {
				lines[i] = "line"
			}
			b.SetLines(lines)
			b.MoveCursor(Position{Line: tt.startLine, Col: 0})

			b.MoveCursorPageUp(tt.pageSize)

			gotCursor := b.GetCursor()
			if gotCursor.Line != tt.wantLine {
				t.Errorf("MoveCursorPageUp() cursor line = %d, want %d", gotCursor.Line, tt.wantLine)
			}
		})
	}
}

func TestMoveCursorPageDown(t *testing.T) {
	tests := []struct {
		name      string
		lines     int
		startLine int
		pageSize  int
		wantLine  int
	}{
		{
			name:      "move down by page",
			lines:     20,
			startLine: 5,
			pageSize:  10,
			wantLine:  15,
		},
		{
			name:      "stops at bottom",
			lines:     20,
			startLine: 15,
			pageSize:  10,
			wantLine:  19,
		},
		{
			name:      "default page size",
			lines:     20,
			startLine: 5,
			pageSize:  0,
			wantLine:  15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBuffer()
			lines := make([]string, tt.lines)
			for i := range lines {
				lines[i] = "line"
			}
			b.SetLines(lines)
			b.MoveCursor(Position{Line: tt.startLine, Col: 0})

			b.MoveCursorPageDown(tt.pageSize)

			gotCursor := b.GetCursor()
			if gotCursor.Line != tt.wantLine {
				t.Errorf("MoveCursorPageDown() cursor line = %d, want %d", gotCursor.Line, tt.wantLine)
			}
		})
	}
}

func TestGetCurrentLineIndentation(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "no indentation",
			content: "hello",
			want:    "",
		},
		{
			name:    "spaces indentation",
			content: "    hello",
			want:    "    ",
		},
		{
			name:    "tabs indentation",
			content: "\thello",
			want:    "\t",
		},
		{
			name:    "mixed indentation",
			content: "  \t  hello",
			want:    "  \t  ",
		},
		{
			name:    "empty line",
			content: "",
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBuffer()
			b.SetLines([]string{tt.content})
			b.MoveCursor(Position{Line: 0, Col: 5})

			got := b.GetCurrentLineIndentation()
			if got != tt.want {
				t.Errorf("GetCurrentLineIndentation() = %q, want %q", got, tt.want)
			}
		})
	}
}

// Helper function
func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
