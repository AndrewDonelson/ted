package buffer

import (
	"testing"
)

func TestBuffer_MoveCursorLeft(t *testing.T) {
	tests := []struct {
		name     string
		initial  []string
		startPos Position
		wantPos  Position
	}{
		{
			name:     "move left within line",
			initial:  []string{"hello"},
			startPos: Position{Line: 0, Col: 3},
			wantPos:  Position{Line: 0, Col: 2},
		},
		{
			name:     "move left from start of line",
			initial:  []string{"line1", "line2"},
			startPos: Position{Line: 1, Col: 0},
			wantPos:  Position{Line: 0, Col: 5},
		},
		{
			name:     "move left from start of document",
			initial:  []string{"hello"},
			startPos: Position{Line: 0, Col: 0},
			wantPos:  Position{Line: 0, Col: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := NewBuffer()
			buf.SetLines(tt.initial)
			buf.MoveCursor(tt.startPos)

			buf.MoveCursorLeft()

			got := buf.GetCursor()
			if got.Line != tt.wantPos.Line || got.Col != tt.wantPos.Col {
				t.Errorf("MoveCursorLeft() cursor = %v, want %v", got, tt.wantPos)
			}
		})
	}
}

func TestBuffer_MoveCursorRight(t *testing.T) {
	tests := []struct {
		name     string
		initial  []string
		startPos Position
		wantPos  Position
	}{
		{
			name:     "move right within line",
			initial:  []string{"hello"},
			startPos: Position{Line: 0, Col: 2},
			wantPos:  Position{Line: 0, Col: 3},
		},
		{
			name:     "move right from end of line",
			initial:  []string{"line1", "line2"},
			startPos: Position{Line: 0, Col: 5},
			wantPos:  Position{Line: 1, Col: 0},
		},
		{
			name:     "move right from end of document",
			initial:  []string{"hello"},
			startPos: Position{Line: 0, Col: 5},
			wantPos:  Position{Line: 0, Col: 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := NewBuffer()
			buf.SetLines(tt.initial)
			buf.MoveCursor(tt.startPos)

			buf.MoveCursorRight()

			got := buf.GetCursor()
			if got.Line != tt.wantPos.Line || got.Col != tt.wantPos.Col {
				t.Errorf("MoveCursorRight() cursor = %v, want %v", got, tt.wantPos)
			}
		})
	}
}

func TestBuffer_MoveCursorUp(t *testing.T) {
	tests := []struct {
		name     string
		initial  []string
		startPos Position
		wantPos  Position
	}{
		{
			name:     "move up preserving column",
			initial:  []string{"line1", "line2"},
			startPos: Position{Line: 1, Col: 3},
			wantPos:  Position{Line: 0, Col: 3},
		},
		{
			name:     "move up adjusting column",
			initial:  []string{"short", "very long line"},
			startPos: Position{Line: 1, Col: 10},
			wantPos:  Position{Line: 0, Col: 5},
		},
		{
			name:     "move up from first line",
			initial:  []string{"line1", "line2"},
			startPos: Position{Line: 0, Col: 3},
			wantPos:  Position{Line: 0, Col: 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := NewBuffer()
			buf.SetLines(tt.initial)
			buf.MoveCursor(tt.startPos)

			buf.MoveCursorUp()

			got := buf.GetCursor()
			if got.Line != tt.wantPos.Line || got.Col != tt.wantPos.Col {
				t.Errorf("MoveCursorUp() cursor = %v, want %v", got, tt.wantPos)
			}
		})
	}
}

func TestBuffer_MoveCursorDown(t *testing.T) {
	tests := []struct {
		name     string
		initial  []string
		startPos Position
		wantPos  Position
	}{
		{
			name:     "move down preserving column",
			initial:  []string{"line1", "line2"},
			startPos: Position{Line: 0, Col: 3},
			wantPos:  Position{Line: 1, Col: 3},
		},
		{
			name:     "move down adjusting column",
			initial:  []string{"very long line", "short"},
			startPos: Position{Line: 0, Col: 10},
			wantPos:  Position{Line: 1, Col: 5},
		},
		{
			name:     "move down from last line",
			initial:  []string{"line1", "line2"},
			startPos: Position{Line: 1, Col: 3},
			wantPos:  Position{Line: 1, Col: 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := NewBuffer()
			buf.SetLines(tt.initial)
			buf.MoveCursor(tt.startPos)

			buf.MoveCursorDown()

			got := buf.GetCursor()
			if got.Line != tt.wantPos.Line || got.Col != tt.wantPos.Col {
				t.Errorf("MoveCursorDown() cursor = %v, want %v", got, tt.wantPos)
			}
		})
	}
}

func TestBuffer_MoveCursorToLineStart(t *testing.T) {
	buf := NewBuffer()
	buf.SetLines([]string{"hello", "world"})
	buf.MoveCursor(Position{Line: 0, Col: 3})

	buf.MoveCursorToLineStart()

	got := buf.GetCursor()
	if got.Line != 0 || got.Col != 0 {
		t.Errorf("MoveCursorToLineStart() cursor = %v, want {Line: 0, Col: 0}", got)
	}
}

func TestBuffer_MoveCursorToLineEnd(t *testing.T) {
	buf := NewBuffer()
	buf.SetLines([]string{"hello", "world"})
	buf.MoveCursor(Position{Line: 0, Col: 2})

	buf.MoveCursorToLineEnd()

	got := buf.GetCursor()
	if got.Line != 0 || got.Col != 5 {
		t.Errorf("MoveCursorToLineEnd() cursor = %v, want {Line: 0, Col: 5}", got)
	}
}

func TestBuffer_MoveCursorToDocumentStart(t *testing.T) {
	buf := NewBuffer()
	buf.SetLines([]string{"line1", "line2", "line3"})
	buf.MoveCursor(Position{Line: 2, Col: 3})

	buf.MoveCursorToDocumentStart()

	got := buf.GetCursor()
	if got.Line != 0 || got.Col != 0 {
		t.Errorf("MoveCursorToDocumentStart() cursor = %v, want {Line: 0, Col: 0}", got)
	}
}

func TestBuffer_MoveCursorToDocumentEnd(t *testing.T) {
	buf := NewBuffer()
	buf.SetLines([]string{"line1", "line2", "line3"})
	buf.MoveCursor(Position{Line: 0, Col: 0})

	buf.MoveCursorToDocumentEnd()

	got := buf.GetCursor()
	wantLine := 2
	wantCol := len("line3")
	if got.Line != wantLine || got.Col != wantCol {
		t.Errorf("MoveCursorToDocumentEnd() cursor = %v, want {Line: %d, Col: %d}", got, wantLine, wantCol)
	}
}

func TestBuffer_MoveCursor(t *testing.T) {
	tests := []struct {
		name    string
		initial []string
		pos     Position
		wantPos Position
	}{
		{
			name:    "valid position",
			initial: []string{"hello", "world"},
			pos:     Position{Line: 1, Col: 2},
			wantPos: Position{Line: 1, Col: 2},
		},
		{
			name:    "adjust negative line",
			initial: []string{"hello"},
			pos:     Position{Line: -1, Col: 2},
			wantPos: Position{Line: 0, Col: 2},
		},
		{
			name:    "adjust line beyond end",
			initial: []string{"hello"},
			pos:     Position{Line: 10, Col: 2},
			wantPos: Position{Line: 0, Col: 2},
		},
		{
			name:    "adjust negative column",
			initial: []string{"hello"},
			pos:     Position{Line: 0, Col: -1},
			wantPos: Position{Line: 0, Col: 0},
		},
		{
			name:    "adjust column beyond line length",
			initial: []string{"hello"},
			pos:     Position{Line: 0, Col: 10},
			wantPos: Position{Line: 0, Col: 5},
		},
		{
			name:    "empty buffer",
			initial: []string{""},
			pos:     Position{Line: 0, Col: 5},
			wantPos: Position{Line: 0, Col: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := NewBuffer()
			buf.SetLines(tt.initial)

			buf.MoveCursor(tt.pos)

			got := buf.GetCursor()
			if got.Line != tt.wantPos.Line || got.Col != tt.wantPos.Col {
				t.Errorf("MoveCursor() cursor = %v, want %v", got, tt.wantPos)
			}
		})
	}
}
