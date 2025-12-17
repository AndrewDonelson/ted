package buffer

import (
	"reflect"
	"testing"
)

func TestNewBuffer(t *testing.T) {
	buf := NewBuffer()

	if buf.LineCount() != 1 {
		t.Errorf("NewBuffer() LineCount = %d, want 1", buf.LineCount())
	}

	line, err := buf.GetLine(0)
	if err != nil {
		t.Fatalf("NewBuffer() GetLine(0) error = %v", err)
	}
	if line != "" {
		t.Errorf("NewBuffer() GetLine(0) = %q, want \"\"", line)
	}

	cursor := buf.GetCursor()
	if cursor.Line != 0 || cursor.Col != 0 {
		t.Errorf("NewBuffer() cursor = %v, want {Line: 0, Col: 0}", cursor)
	}

	if buf.IsModified() {
		t.Error("NewBuffer() IsModified = true, want false")
	}
}

func TestBuffer_Insert(t *testing.T) {
	tests := []struct {
		name     string
		initial  []string
		pos      Position
		text     string
		want     []string
		wantErr  bool
		noModify bool // For empty inserts
	}{
		{
			name:    "insert at beginning",
			initial: []string{"hello world"},
			pos:     Position{Line: 0, Col: 0},
			text:    "foo ",
			want:    []string{"foo hello world"},
			wantErr: false,
		},
		{
			name:    "insert at end",
			initial: []string{"hello"},
			pos:     Position{Line: 0, Col: 5},
			text:    " world",
			want:    []string{"hello world"},
			wantErr: false,
		},
		{
			name:    "insert in middle",
			initial: []string{"hello world"},
			pos:     Position{Line: 0, Col: 6},
			text:    "beautiful ",
			want:    []string{"hello beautiful world"},
			wantErr: false,
		},
		{
			name:    "insert newline",
			initial: []string{"hello world"},
			pos:     Position{Line: 0, Col: 5},
			text:    "\n",
			want:    []string{"hello", " world"},
			wantErr: false,
		},
		{
			name:    "insert multiple newlines",
			initial: []string{"hello"},
			pos:     Position{Line: 0, Col: 5},
			text:    "\nworld\nfoo",
			want:    []string{"hello", "world", "foo"},
			wantErr: false,
		},
		{
			name:    "insert at invalid position",
			initial: []string{"hello"},
			pos:     Position{Line: 1, Col: 0},
			text:    "foo",
			want:    nil,
			wantErr: true,
		},
		{
			name:     "insert empty string",
			initial:  []string{"hello"},
			pos:      Position{Line: 0, Col: 0},
			text:     "",
			want:     []string{"hello"},
			wantErr:  false,
			noModify: true, // Empty insert should not modify
		},
		{
			name:    "insert at column beyond line length",
			initial: []string{"hello"},
			pos:     Position{Line: 0, Col: 10},
			text:    "foo",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "insert with newline at start",
			initial: []string{"world"},
			pos:     Position{Line: 0, Col: 0},
			text:    "hello\n",
			want:    []string{"hello", "world"},
			wantErr: false,
		},
		{
			name:    "insert across multiple existing lines",
			initial: []string{"line1", "line2", "line3"},
			pos:     Position{Line: 1, Col: 2},
			text:    "new\nmiddle",
			want:    []string{"line1", "linew", "middlene2", "line3"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := NewBuffer()
			buf.SetLines(tt.initial)

			err := buf.Insert(tt.pos, tt.text)

			if (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				got := buf.GetAllLines()
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Insert() got = %v, want %v", got, tt.want)
				}

				// Check that buffer is marked as modified (unless empty insert)
				noModify := tt.text == ""
				if noModify && buf.IsModified() {
					t.Error("Insert() IsModified = true for empty insert, want false")
				} else if !noModify && !buf.IsModified() {
					t.Error("Insert() IsModified = false, want true")
				}
			}
		})
	}
}

func TestBuffer_Delete(t *testing.T) {
	tests := []struct {
		name     string
		initial  []string
		start    Position
		end      Position
		want     []string
		wantErr  bool
		noModify bool // For no-op deletes
	}{
		{
			name:    "delete single character",
			initial: []string{"hello"},
			start:   Position{Line: 0, Col: 1},
			end:     Position{Line: 0, Col: 2},
			want:    []string{"hllo"},
			wantErr: false,
		},
		{
			name:    "delete from start of line",
			initial: []string{"hello"},
			start:   Position{Line: 0, Col: 0},
			end:     Position{Line: 0, Col: 2},
			want:    []string{"llo"},
			wantErr: false,
		},
		{
			name:    "delete to end of line",
			initial: []string{"hello"},
			start:   Position{Line: 0, Col: 2},
			end:     Position{Line: 0, Col: 5},
			want:    []string{"he"},
			wantErr: false,
		},
		{
			name:    "delete entire line content",
			initial: []string{"line1", "line2", "line3"},
			start:   Position{Line: 1, Col: 0},
			end:     Position{Line: 1, Col: 5},
			want:    []string{"line1", "line3"},
			wantErr: false,
		},
		{
			name:    "delete across two lines",
			initial: []string{"hello", "world"},
			start:   Position{Line: 0, Col: 3},
			end:     Position{Line: 1, Col: 2},
			want:    []string{"helrld"},
			wantErr: false,
		},
		{
			name:    "delete across multiple lines",
			initial: []string{"line1", "line2", "line3", "line4"},
			start:   Position{Line: 0, Col: 3},
			end:     Position{Line: 2, Col: 3},
			want:    []string{"line3", "line4"},
			wantErr: false,
		},
		{
			name:     "delete nothing (same position)",
			initial:  []string{"hello"},
			start:    Position{Line: 0, Col: 2},
			end:      Position{Line: 0, Col: 2},
			want:     []string{"hello"},
			wantErr:  false,
			noModify: true, // No-op delete should not modify
		},
		{
			name:    "delete invalid start position",
			initial: []string{"hello"},
			start:   Position{Line: 1, Col: 0},
			end:     Position{Line: 0, Col: 2},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "delete invalid end position",
			initial: []string{"hello"},
			start:   Position{Line: 0, Col: 0},
			end:     Position{Line: 1, Col: 0},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "delete start after end",
			initial: []string{"hello"},
			start:   Position{Line: 0, Col: 3},
			end:     Position{Line: 0, Col: 2},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := NewBuffer()
			buf.SetLines(tt.initial)
			buf.MarkSaved() // Reset modified flag

			err := buf.Delete(tt.start, tt.end)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				got := buf.GetAllLines()
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Delete() got = %v, want %v", got, tt.want)
				}

				// Check that buffer is marked as modified (unless no-op)
				noModify := false
				if tt.start.Line == tt.end.Line && tt.start.Col == tt.end.Col {
					noModify = true
				}
				if noModify && buf.IsModified() {
					t.Error("Delete() IsModified = true for no-op, want false")
				} else if !noModify && !buf.IsModified() {
					t.Error("Delete() IsModified = false, want true")
				}

				// Check cursor position (only if not no-op)
				if !noModify {
					cursor := buf.GetCursor()
					if cursor.Line != tt.start.Line || cursor.Col != tt.start.Col {
						t.Errorf("Delete() cursor = %v, want %v", cursor, tt.start)
					}
				}
			}
		})
	}
}

func TestBuffer_GetLine(t *testing.T) {
	buf := NewBuffer()
	buf.SetLines([]string{"line1", "line2", "line3"})

	tests := []struct {
		name    string
		lineNum int
		want    string
		wantErr bool
	}{
		{
			name:    "get first line",
			lineNum: 0,
			want:    "line1",
			wantErr: false,
		},
		{
			name:    "get middle line",
			lineNum: 1,
			want:    "line2",
			wantErr: false,
		},
		{
			name:    "get last line",
			lineNum: 2,
			want:    "line3",
			wantErr: false,
		},
		{
			name:    "get invalid negative line",
			lineNum: -1,
			want:    "",
			wantErr: true,
		},
		{
			name:    "get invalid line beyond end",
			lineNum: 10,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buf.GetLine(tt.lineNum)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("GetLine() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBuffer_LineCount(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want  int
	}{
		{
			name:  "empty buffer",
			lines: []string{""},
			want:  1,
		},
		{
			name:  "single line",
			lines: []string{"hello"},
			want:  1,
		},
		{
			name:  "multiple lines",
			lines: []string{"line1", "line2", "line3"},
			want:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := NewBuffer()
			buf.SetLines(tt.lines)

			got := buf.LineCount()
			if got != tt.want {
				t.Errorf("LineCount() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestBuffer_IsModified(t *testing.T) {
	buf := NewBuffer()

	if buf.IsModified() {
		t.Error("IsModified() = true for new buffer, want false")
	}

	buf.Insert(Position{Line: 0, Col: 0}, "test")
	if !buf.IsModified() {
		t.Error("IsModified() = false after insert, want true")
	}

	buf.MarkSaved()
	if buf.IsModified() {
		t.Error("IsModified() = true after MarkSaved, want false")
	}

	buf.Delete(Position{Line: 0, Col: 0}, Position{Line: 0, Col: 2})
	if !buf.IsModified() {
		t.Error("IsModified() = false after delete, want true")
	}
}

func TestBuffer_SetLines(t *testing.T) {
	buf := NewBuffer()
	buf.SetLines([]string{"line1", "line2"})

	if buf.LineCount() != 2 {
		t.Errorf("SetLines() LineCount = %d, want 2", buf.LineCount())
	}

	if buf.IsModified() {
		t.Error("SetLines() IsModified = true, want false")
	}

	cursor := buf.GetCursor()
	if cursor.Line != 0 || cursor.Col != 0 {
		t.Errorf("SetLines() cursor = %v, want {Line: 0, Col: 0}", cursor)
	}

	// Test empty lines
	buf.SetLines([]string{})
	if buf.LineCount() != 1 {
		t.Errorf("SetLines([]string{}) LineCount = %d, want 1", buf.LineCount())
	}
}
