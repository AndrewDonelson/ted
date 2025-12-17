package renderer

import (
	"testing"

	"github.com/AndrewDonelson/ted/core/buffer"
	"github.com/AndrewDonelson/ted/ui/layout"
)

func TestRenderTextArea(t *testing.T) {
	mockScr := newMockScreen(80, 24)
	layout := layout.NewLayout(80, 24)
	renderer := NewRenderer(mockScr, layout)

	buf := buffer.NewBuffer()
	buf.SetLines([]string{"line1", "line2", "line3"})
	cursorPos := buffer.Position{Line: 1, Col: 2}

	if err := renderer.RenderTextArea(buf, cursorPos); err != nil {
		t.Errorf("RenderTextArea() error = %v", err)
	}

	// Verify text was rendered
	editRegion := layout.GetEditAreaRegion()
	found := false
	for y := editRegion.Y; y < editRegion.Y+editRegion.Height; y++ {
		if row, ok := mockScr.contents[y]; ok {
			if len(row) > 0 {
				found = true
				break
			}
		}
	}

	if !found {
		t.Error("RenderTextArea() did not render text")
	}
}

func TestRenderTextArea_EmptyBuffer(t *testing.T) {
	mockScr := newMockScreen(80, 24)
	layout := layout.NewLayout(80, 24)
	renderer := NewRenderer(mockScr, layout)

	buf := buffer.NewBuffer()
	cursorPos := buffer.Position{Line: 0, Col: 0}

	if err := renderer.RenderTextArea(buf, cursorPos); err != nil {
		t.Errorf("RenderTextArea() error = %v", err)
	}
}

func TestRenderTextArea_LongLines(t *testing.T) {
	mockScr := newMockScreen(80, 24)
	layout := layout.NewLayout(80, 24)
	renderer := NewRenderer(mockScr, layout)

	longLine := make([]rune, 200)
	for i := range longLine {
		longLine[i] = 'A'
	}

	buf := buffer.NewBuffer()
	buf.SetLines([]string{string(longLine)})
	cursorPos := buffer.Position{Line: 0, Col: 0}

	if err := renderer.RenderTextArea(buf, cursorPos); err != nil {
		t.Errorf("RenderTextArea() error = %v", err)
	}

	// Verify line was truncated
	editRegion := layout.GetEditAreaRegion()
	if row, ok := mockScr.contents[editRegion.Y]; ok {
		if len(row) > editRegion.Width {
			t.Errorf("Long line not truncated: %d > %d", len(row), editRegion.Width)
		}
	}
}

func TestRenderTextArea_CurrentLineHighlight(t *testing.T) {
	mockScr := newMockScreen(80, 24)
	layout := layout.NewLayout(80, 24)
	renderer := NewRenderer(mockScr, layout)

	buf := buffer.NewBuffer()
	buf.SetLines([]string{"line1", "line2", "line3"})
	cursorPos := buffer.Position{Line: 1, Col: 0} // Cursor on line2

	if err := renderer.RenderTextArea(buf, cursorPos); err != nil {
		t.Errorf("RenderTextArea() error = %v", err)
	}

	// Verify current line uses different style
	editRegion := layout.GetEditAreaRegion()
	viewport := layout.CalculateViewport(cursorPos.Line, buf.LineCount())
	currentLineY := editRegion.Y + (cursorPos.Line - viewport.StartLine)

	if rowStyles, ok := mockScr.styles[currentLineY]; ok {
		if len(rowStyles) == 0 {
			t.Error("Current line styles not set")
		}
	}
}

func TestRenderTextArea_Scrolling(t *testing.T) {
	mockScr := newMockScreen(80, 24)
	layout := layout.NewLayout(80, 24)
	renderer := NewRenderer(mockScr, layout)

	// Create buffer with many lines
	lines := make([]string, 100)
	for i := range lines {
		lines[i] = "line" + formatNumber(i+1)
	}

	buf := buffer.NewBuffer()
	buf.SetLines(lines)
	cursorPos := buffer.Position{Line: 50, Col: 0} // Cursor in middle

	if err := renderer.RenderTextArea(buf, cursorPos); err != nil {
		t.Errorf("RenderTextArea() error = %v", err)
	}

	// Verify viewport scrolling
	viewport := layout.CalculateViewport(cursorPos.Line, buf.LineCount())
	if viewport.StartLine > cursorPos.Line || viewport.EndLine < cursorPos.Line {
		t.Errorf("Viewport does not include cursor: start=%d, end=%d, cursor=%d",
			viewport.StartLine, viewport.EndLine, cursorPos.Line)
	}
}

func TestRenderTextArea_UTF8(t *testing.T) {
	mockScr := newMockScreen(80, 24)
	layout := layout.NewLayout(80, 24)
	renderer := NewRenderer(mockScr, layout)

	buf := buffer.NewBuffer()
	buf.SetLines([]string{"hello 世界", "测试", "café"})
	cursorPos := buffer.Position{Line: 0, Col: 0}

	if err := renderer.RenderTextArea(buf, cursorPos); err != nil {
		t.Errorf("RenderTextArea() error = %v", err)
	}
}

func TestRenderTextAreaWithLineNumbers(t *testing.T) {
	mockScr := newMockScreen(80, 24)
	layout := layout.NewLayout(80, 24)
	renderer := NewRenderer(mockScr, layout)

	buf := buffer.NewBuffer()
	buf.SetLines([]string{"line1", "line2", "line3"})
	cursorPos := buffer.Position{Line: 0, Col: 0}

	if err := renderer.RenderTextAreaWithLineNumbers(buf, cursorPos, true); err != nil {
		t.Errorf("RenderTextAreaWithLineNumbers() error = %v", err)
	}

	// Verify line numbers were rendered
	editRegion := layout.GetEditAreaRegion()
	foundLineNum := false
	for y := editRegion.Y; y < editRegion.Y+editRegion.Height; y++ {
		if row, ok := mockScr.contents[y]; ok {
			// Check for separator character
			for x, char := range row {
				if char == '│' && x < 10 {
					foundLineNum = true
					break
				}
			}
		}
		if foundLineNum {
			break
		}
	}

	if !foundLineNum {
		t.Error("RenderTextAreaWithLineNumbers() did not render line numbers")
	}
}

func TestRenderTextAreaWithLineNumbers_Disabled(t *testing.T) {
	mockScr := newMockScreen(80, 24)
	layout := layout.NewLayout(80, 24)
	renderer := NewRenderer(mockScr, layout)

	buf := buffer.NewBuffer()
	buf.SetLines([]string{"line1"})
	cursorPos := buffer.Position{Line: 0, Col: 0}

	if err := renderer.RenderTextAreaWithLineNumbers(buf, cursorPos, false); err != nil {
		t.Errorf("RenderTextAreaWithLineNumbers() error = %v", err)
	}
}

func TestFormatLineNumber(t *testing.T) {
	tests := []struct {
		name      string
		lineNum   int
		width     int
		wantLen   int
		wantRight bool
	}{
		{
			name:      "single digit",
			lineNum:   5,
			width:     5,
			wantLen:   5,
			wantRight: true,
		},
		{
			name:      "two digits",
			lineNum:   42,
			width:     5,
			wantLen:   5,
			wantRight: true,
		},
		{
			name:      "exact width",
			lineNum:   123,
			width:     3,
			wantLen:   3,
			wantRight: true,
		},
		{
			name:      "narrow width",
			lineNum:   999,
			width:     2,
			wantLen:   3, // Number is 3 digits, formatLineNumber doesn't truncate, only pads
			wantRight: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatLineNumber(tt.lineNum, tt.width)
			if len(got) != tt.wantLen {
				t.Errorf("formatLineNumber() length = %d, want %d", len(got), tt.wantLen)
			}

			if tt.wantRight && len(got) > 0 {
				// Check right alignment (should end with digit or be all spaces)
				lastChar := got[len(got)-1]
				if lastChar != ' ' && (lastChar < '0' || lastChar > '9') {
					t.Errorf("formatLineNumber() not right-aligned: %q", got)
				}
			}
		})
	}
}

func TestRenderTextArea_FillEmptySpace(t *testing.T) {
	mockScr := newMockScreen(80, 24)
	layout := layout.NewLayout(80, 24)
	renderer := NewRenderer(mockScr, layout)

	buf := buffer.NewBuffer()
	buf.SetLines([]string{"short"}) // Short line
	cursorPos := buffer.Position{Line: 0, Col: 0}

	if err := renderer.RenderTextArea(buf, cursorPos); err != nil {
		t.Errorf("RenderTextArea() error = %v", err)
	}

	// Verify empty space is filled
	editRegion := layout.GetEditAreaRegion()
	if row, ok := mockScr.contents[editRegion.Y]; ok {
		// Should have content up to editRegion.Width
		if len(row) < editRegion.Width {
			// Check if remaining space is filled with spaces
			for x := len("short"); x < editRegion.Width; x++ {
				if row[x] != ' ' {
					t.Errorf("Empty space at x=%d not filled: got %c", x, row[x])
					break
				}
			}
		}
	}
}
