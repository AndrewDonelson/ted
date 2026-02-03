package renderer

import (
	"testing"

	"github.com/AndrewDonelson/ted/core/buffer"
	"github.com/AndrewDonelson/ted/ui/layout"
	"github.com/gdamore/tcell/v2"
)

// mockScreen is a mock implementation of terminal.Screen for testing.
type mockScreen struct {
	width      int
	height     int
	contents   map[int]map[int]rune
	styles     map[int]map[int]tcell.Style
	cursorX    int
	cursorY    int
	cursorShow bool
	cleared    bool
	refreshed  bool
}

func newMockScreen(width, height int) *mockScreen {
	return &mockScreen{
		width:    width,
		height:   height,
		contents: make(map[int]map[int]rune),
		styles:   make(map[int]map[int]tcell.Style),
	}
}

func (m *mockScreen) Clear() {
	m.cleared = true
	m.contents = make(map[int]map[int]rune)
	m.styles = make(map[int]map[int]tcell.Style)
}

func (m *mockScreen) Refresh() error {
	m.refreshed = true
	return nil
}

func (m *mockScreen) GetSize() (width, height int) {
	return m.width, m.height
}

func (m *mockScreen) SetContent(x, y int, mainc rune, combc []rune, style tcell.Style) error {
	if m.contents[y] == nil {
		m.contents[y] = make(map[int]rune)
		m.styles[y] = make(map[int]tcell.Style)
	}
	m.contents[y][x] = mainc
	m.styles[y][x] = style
	return nil
}

func (m *mockScreen) ShowCursor(x, y int) {
	m.cursorX = x
	m.cursorY = y
	m.cursorShow = true
}

func (m *mockScreen) HideCursor() {
	m.cursorShow = false
}

func (m *mockScreen) PollEvent() tcell.Event {
	return nil // Not used in tests
}

func (m *mockScreen) Fini() {
	// No-op for mock
}

func (m *mockScreen) GetRawScreen() tcell.Screen {
	// Return nil for mock - dialogs won't be tested here
	return nil
}

func TestNewRenderer(t *testing.T) {
	mockScr := newMockScreen(80, 24)
	layout := layout.NewLayout(80, 24)

	renderer := NewRenderer(mockScr, layout)
	if renderer == nil {
		t.Fatal("NewRenderer() returned nil")
	}
}

func TestRenderer_Clear(t *testing.T) {
	mockScr := newMockScreen(80, 24)
	layout := layout.NewLayout(80, 24)
	renderer := NewRenderer(mockScr, layout)

	renderer.Clear()

	if !mockScr.cleared {
		t.Error("Clear() did not call screen.Clear()")
	}
}

func TestRenderer_Refresh(t *testing.T) {
	mockScr := newMockScreen(80, 24)
	layout := layout.NewLayout(80, 24)
	renderer := NewRenderer(mockScr, layout)

	if err := renderer.Refresh(); err != nil {
		t.Errorf("Refresh() error = %v", err)
	}

	if !mockScr.refreshed {
		t.Error("Refresh() did not call screen.Refresh()")
	}
}

func TestGetDefaultStyle(t *testing.T) {
	style := GetDefaultStyle()
	if style == tcell.StyleDefault {
		t.Error("GetDefaultStyle() returned default style")
	}
}

func TestGetMenuBarStyle(t *testing.T) {
	style := GetMenuBarStyle()
	if style == tcell.StyleDefault {
		t.Error("GetMenuBarStyle() returned default style")
	}
}

func TestGetInfoBarStyle(t *testing.T) {
	style := GetInfoBarStyle()
	if style == tcell.StyleDefault {
		t.Error("GetInfoBarStyle() returned default style")
	}

	// Verify it's actually inverted (different from default)
	defaultStyle := GetDefaultStyle()
	if style == defaultStyle {
		t.Error("GetInfoBarStyle() should use inverted colors")
	}
}

func TestGetLineNumberStyle(t *testing.T) {
	style := GetLineNumberStyle()
	if style == tcell.StyleDefault {
		t.Error("GetLineNumberStyle() returned default style")
	}
}

func TestGetCurrentLineStyle(t *testing.T) {
	style := GetCurrentLineStyle()
	if style == tcell.StyleDefault {
		t.Error("GetCurrentLineStyle() returned default style")
	}
}

func TestGetCursorStyle(t *testing.T) {
	style := GetCursorStyle()
	if style == tcell.StyleDefault {
		t.Error("GetCursorStyle() returned default style")
	}
}

func TestRenderer_RenderAll(t *testing.T) {
	mockScr := newMockScreen(80, 24)
	layout := layout.NewLayout(80, 24)
	renderer := NewRenderer(mockScr, layout)

	buf := buffer.NewBuffer()
	buf.SetLines([]string{"line1", "line2", "line3"})
	cursorPos := buffer.Position{Line: 0, Col: 0}
	fileInfo := &FileInfo{
		Name:       "test.txt",
		Size:       100,
		Type:       "Plain Text",
		Encoding:   "UTF-8",
		LineEnding: "LF",
		TabSize:    4,
		TotalLines: 3,
		IsModified: false,
	}

	if err := renderer.RenderAll(buf, cursorPos, fileInfo); err != nil {
		t.Errorf("RenderAll() error = %v", err)
	}

	if !mockScr.cleared {
		t.Error("RenderAll() did not clear screen")
	}

	if !mockScr.refreshed {
		t.Error("RenderAll() did not refresh screen")
	}

	if !mockScr.cursorShow {
		t.Error("RenderAll() did not show cursor")
	}
}

func TestRenderer_RenderAll_EmptyBuffer(t *testing.T) {
	mockScr := newMockScreen(80, 24)
	layout := layout.NewLayout(80, 24)
	renderer := NewRenderer(mockScr, layout)

	buf := buffer.NewBuffer()
	cursorPos := buffer.Position{Line: 0, Col: 0}
	fileInfo := &FileInfo{Name: "[No Name]"}

	if err := renderer.RenderAll(buf, cursorPos, fileInfo); err != nil {
		t.Errorf("RenderAll() error = %v", err)
	}
}

func TestRenderer_RenderAll_NilFileInfo(t *testing.T) {
	mockScr := newMockScreen(80, 24)
	layout := layout.NewLayout(80, 24)
	renderer := NewRenderer(mockScr, layout)

	buf := buffer.NewBuffer()
	buf.SetLines([]string{"test"})
	cursorPos := buffer.Position{Line: 0, Col: 0}

	if err := renderer.RenderAll(buf, cursorPos, nil); err != nil {
		t.Errorf("RenderAll() error = %v", err)
	}
}

func TestRenderer_RenderAll_CursorOutOfBounds(t *testing.T) {
	mockScr := newMockScreen(80, 24)
	layout := layout.NewLayout(80, 24)
	renderer := NewRenderer(mockScr, layout)

	buf := buffer.NewBuffer()
	buf.SetLines([]string{"line1"})
	cursorPos := buffer.Position{Line: 100, Col: 100} // Out of bounds
	fileInfo := &FileInfo{Name: "test.txt"}

	if err := renderer.RenderAll(buf, cursorPos, fileInfo); err != nil {
		t.Errorf("RenderAll() error = %v", err)
	}
}
