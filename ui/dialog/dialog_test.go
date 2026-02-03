package dialog

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

// mockScreen is a mock implementation of the Screen interface for testing.
type mockScreen struct {
	contents map[int]map[int]rune
	styles   map[int]map[int]tcell.Style
}

func newMockScreen() *mockScreen {
	return &mockScreen{
		contents: make(map[int]map[int]rune),
		styles:   make(map[int]map[int]tcell.Style),
	}
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

func TestNewInputDialog(t *testing.T) {
	dlg := NewInputDialog("Test Title", "Enter value:", "default", nil, nil)

	if dlg.title != "Test Title" {
		t.Errorf("title = %q, want %q", dlg.title, "Test Title")
	}

	if dlg.prompt != "Enter value:" {
		t.Errorf("prompt = %q, want %q", dlg.prompt, "Enter value:")
	}

	if dlg.input != "default" {
		t.Errorf("input = %q, want %q", dlg.input, "default")
	}

	if dlg.cursorPos != 7 { // len("default")
		t.Errorf("cursorPos = %d, want %d", dlg.cursorPos, 7)
	}
}

func TestInputDialog_ShowHide(t *testing.T) {
	dlg := NewInputDialog("Test", "Prompt:", "", nil, nil)

	if dlg.IsOpen() {
		t.Error("dialog should not be open initially")
	}

	dlg.Show(80, 24)

	if !dlg.IsOpen() {
		t.Error("dialog should be open after Show()")
	}

	if dlg.x <= 0 || dlg.y <= 0 {
		t.Errorf("dialog position (%d, %d) should be positive", dlg.x, dlg.y)
	}

	dlg.Hide()

	if dlg.IsOpen() {
		t.Error("dialog should not be open after Hide()")
	}
}

func TestInputDialog_HandleInput_Character(t *testing.T) {
	dlg := NewInputDialog("Test", "Prompt:", "", nil, nil)
	dlg.Show(80, 24)

	// Type "abc"
	dlg.HandleInput(tcell.KeyRune, 0, 'a')
	dlg.HandleInput(tcell.KeyRune, 0, 'b')
	dlg.HandleInput(tcell.KeyRune, 0, 'c')

	if dlg.input != "abc" {
		t.Errorf("input = %q, want %q", dlg.input, "abc")
	}

	if dlg.cursorPos != 3 {
		t.Errorf("cursorPos = %d, want %d", dlg.cursorPos, 3)
	}
}

func TestInputDialog_HandleInput_Backspace(t *testing.T) {
	dlg := NewInputDialog("Test", "Prompt:", "abc", nil, nil)
	dlg.Show(80, 24)

	// Backspace once
	dlg.HandleInput(tcell.KeyBackspace, 0, 0)

	if dlg.input != "ab" {
		t.Errorf("input = %q, want %q", dlg.input, "ab")
	}

	if dlg.cursorPos != 2 {
		t.Errorf("cursorPos = %d, want %d", dlg.cursorPos, 2)
	}
}

func TestInputDialog_HandleInput_Enter(t *testing.T) {
	confirmed := false
	var result string

	dlg := NewInputDialog("Test", "Prompt:", "hello", func(s string) {
		confirmed = true
		result = s
	}, nil)

	dlg.Show(80, 24)
	dlg.HandleInput(tcell.KeyEnter, 0, 0)

	if !dlg.IsConfirmed() {
		t.Error("dialog should be confirmed after Enter")
	}

	if !confirmed {
		t.Error("onConfirm callback should have been called")
	}

	if result != "hello" {
		t.Errorf("result = %q, want %q", result, "hello")
	}
}

func TestInputDialog_HandleInput_Escape(t *testing.T) {
	cancelled := false

	dlg := NewInputDialog("Test", "Prompt:", "hello", nil, func() {
		cancelled = true
	})

	dlg.Show(80, 24)
	dlg.HandleInput(tcell.KeyEscape, 0, 0)

	if !dlg.IsCancelled() {
		t.Error("dialog should be cancelled after Escape")
	}

	if !cancelled {
		t.Error("onCancel callback should have been called")
	}
}

func TestInputDialog_GetResult(t *testing.T) {
	dlg := NewInputDialog("Test", "Prompt:", "test value", nil, nil)

	result := dlg.GetResult()
	if result != "test value" {
		t.Errorf("GetResult() = %v, want %v", result, "test value")
	}
}

func TestNewConfirmDialog(t *testing.T) {
	dlg := NewConfirmDialog("Confirm", "Are you sure?", nil, nil)

	if dlg.title != "Confirm" {
		t.Errorf("title = %q, want %q", dlg.title, "Confirm")
	}

	if dlg.message != "Are you sure?" {
		t.Errorf("message = %q, want %q", dlg.message, "Are you sure?")
	}
}

func TestConfirmDialog_HandleInput_Tab(t *testing.T) {
	dlg := NewConfirmDialog("Test", "Message", nil, nil)
	dlg.Show(80, 24)

	if dlg.focusIndex != 0 {
		t.Errorf("focusIndex = %d, want %d", dlg.focusIndex, 0)
	}

	dlg.HandleInput(tcell.KeyTab, 0, 0)

	if dlg.focusIndex != 1 {
		t.Errorf("focusIndex = %d, want %d after Tab", dlg.focusIndex, 1)
	}
}

func TestConfirmDialog_HandleInput_Enter(t *testing.T) {
	confirmed := false

	dlg := NewConfirmDialog("Test", "Message", func() {
		confirmed = true
	}, nil)

	dlg.Show(80, 24)
	dlg.HandleInput(tcell.KeyEnter, 0, 0)

	if !dlg.IsConfirmed() {
		t.Error("dialog should be confirmed")
	}

	if !confirmed {
		t.Error("onConfirm callback should have been called")
	}
}

func TestNewGoToLineDialog(t *testing.T) {
	dlg := NewGoToLineDialog(100, func(line int) {
		// Test callback
	}, nil)

	if dlg.maxLine != 100 {
		t.Errorf("maxLine = %d, want %d", dlg.maxLine, 100)
	}

	if dlg.title != "Go to Line" {
		t.Errorf("title = %q, want %q", dlg.title, "Go to Line")
	}
}

func TestGoToLineDialog_Confirm(t *testing.T) {
	var result int

	dlg := NewGoToLineDialog(100, func(line int) {
		result = line
	}, nil)

	dlg.Show(80, 24)

	// Type "42"
	dlg.HandleInput(tcell.KeyRune, 0, '4')
	dlg.HandleInput(tcell.KeyRune, 0, '2')
	dlg.HandleInput(tcell.KeyEnter, 0, 0)

	if result != 42 {
		t.Errorf("result = %d, want %d", result, 42)
	}
}

func TestNewOpenFileDialog(t *testing.T) {
	dlg := NewOpenFileDialog("/home/user/", nil, nil)

	if dlg.title != "Open File" {
		t.Errorf("title = %q, want %q", dlg.title, "Open File")
	}

	if dlg.input != "/home/user/" {
		t.Errorf("input = %q, want %q", dlg.input, "/home/user/")
	}
}

func TestNewSaveAsDialog(t *testing.T) {
	dlg := NewSaveAsDialog("/home/user/file.txt", nil, nil)

	if dlg.title != "Save As" {
		t.Errorf("title = %q, want %q", dlg.title, "Save As")
	}

	if dlg.input != "/home/user/file.txt" {
		t.Errorf("input = %q, want %q", dlg.input, "/home/user/file.txt")
	}
}

func TestNewUnsavedChangesDialog(t *testing.T) {
	dlg := NewUnsavedChangesDialog("test.txt", nil, nil, nil)

	if dlg.title != "Unsaved Changes" {
		t.Errorf("title = %q, want %q", dlg.title, "Unsaved Changes")
	}
}

func TestDialogManager(t *testing.T) {
	dm := NewDialogManager()

	if !dm.IsEmpty() {
		t.Error("new dialog manager should be empty")
	}

	if dm.HasOpenDialog() {
		t.Error("new dialog manager should have no open dialogs")
	}

	if dm.Peek() != nil {
		t.Error("Peek() should return nil for empty manager")
	}

	if dm.Pop() != nil {
		t.Error("Pop() should return nil for empty manager")
	}
}

func TestDialogManager_PushPop(t *testing.T) {
	dm := NewDialogManager()
	dlg := NewInputDialog("Test", "Prompt:", "", nil, nil)

	dm.Push(dlg, 80, 24)

	if dm.IsEmpty() {
		t.Error("dialog manager should not be empty after push")
	}

	if !dm.HasOpenDialog() {
		t.Error("dialog manager should have open dialog after push")
	}

	if dm.Peek() != dlg {
		t.Error("Peek() should return the pushed dialog")
	}

	popped := dm.Pop()
	if popped != dlg {
		t.Error("Pop() should return the pushed dialog")
	}

	if !dm.IsEmpty() {
		t.Error("dialog manager should be empty after pop")
	}

	// Dialog should be closed after pop (not open)
	if popped.IsOpen() {
		t.Error("dialog should be closed after pop")
	}
}

func TestDialogManager_HandleInput(t *testing.T) {
	dm := NewDialogManager()
	confirmed := false

	dlg := NewConfirmDialog("Test", "Message", func() {
		confirmed = true
	}, nil)

	dm.Push(dlg, 80, 24)

	// Send Enter key
	handled := dm.HandleInput(tcell.KeyEnter, 0, 0)

	if !handled {
		t.Error("dialog should handle Enter key")
	}

	if !confirmed {
		t.Error("dialog should have been confirmed")
	}

	// Dialog should be closed now
	if dm.HasOpenDialog() {
		t.Error("dialog manager should have no open dialogs after confirmation")
	}
}

func TestDialogManager_MultipleDialogs(t *testing.T) {
	dm := NewDialogManager()

	dlg1 := NewConfirmDialog("First", "First message", nil, nil)
	dlg2 := NewInputDialog("Second", "Second prompt:", "", nil, nil)

	dm.Push(dlg1, 80, 24)
	dm.Push(dlg2, 80, 24)

	// Top dialog should be dlg2
	if dm.Peek() != dlg2 {
		t.Error("top dialog should be dlg2")
	}

	// Handle input should go to dlg2
	dm.HandleInput(tcell.KeyEscape, 0, 0)

	// Now dlg1 should be on top
	if dm.Peek() != dlg1 {
		t.Error("top dialog should be dlg1 after dlg2 closed")
	}
}

func TestBaseDialog_DrawBorder(t *testing.T) {
	d := &BaseDialog{
		title:  "Test",
		width:  20,
		height: 10,
		x:      5,
		y:      5,
	}

	screen := newMockScreen()
	style := tcell.StyleDefault

	d.DrawBorder(screen, style)

	// Check corners
	if screen.contents[5][5] != '┌' {
		t.Errorf("top-left corner = %q, want %q", screen.contents[5][5], '┌')
	}

	if screen.contents[5][24] != '┐' {
		t.Errorf("top-right corner = %q, want %q", screen.contents[5][24], '┐')
	}

	if screen.contents[14][5] != '└' {
		t.Errorf("bottom-left corner = %q, want %q", screen.contents[14][5], '└')
	}

	if screen.contents[14][24] != '┘' {
		t.Errorf("bottom-right corner = %q, want %q", screen.contents[14][24], '┘')
	}
}

func TestBaseDialog_DrawButton(t *testing.T) {
	d := &BaseDialog{
		width:  40,
		height: 10,
		x:      0,
		y:      0,
	}

	screen := newMockScreen()
	style := tcell.StyleDefault

	d.DrawButton(screen, 5, 5, 0, "OK", style, false)

	// Check button text "[ OK ]"
	expected := "[ OK ]"
	for i, ch := range expected {
		if screen.contents[5][5+i] != ch {
			t.Errorf("button char %d = %q, want %q", i, screen.contents[5][5+i], ch)
		}
	}
}

func TestBaseDialog_DrawText(t *testing.T) {
	d := &BaseDialog{
		width:  40,
		height: 10,
		x:      0,
		y:      0,
	}

	screen := newMockScreen()
	style := tcell.StyleDefault

	d.DrawText(screen, 5, 5, "Hello", style)

	for i, ch := range "Hello" {
		if screen.contents[5][5+i] != ch {
			t.Errorf("text char %d = %q, want %q", i, screen.contents[5][5+i], ch)
		}
	}
}

func TestBaseDialog_Clear(t *testing.T) {
	d := &BaseDialog{
		width:  10,
		height: 5,
		x:      0,
		y:      0,
	}

	screen := newMockScreen()

	// Set some content first
	screen.SetContent(2, 2, 'X', []rune{}, tcell.StyleDefault)

	style := tcell.StyleDefault
	d.Clear(screen, style)

	// Check that content was cleared
	if screen.contents[2][2] != ' ' {
		t.Errorf("cleared cell = %q, want %q", screen.contents[2][2], ' ')
	}
}
