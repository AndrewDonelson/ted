// Package dialog implements modal dialogs for the editor.
//
// It provides a framework for creating input dialogs, confirmation dialogs,
// and other modal interactions that overlay the editor content.
package dialog

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
)

// Dialog represents a modal dialog that can be shown, hidden, and rendered.
type Dialog interface {
	// Show opens the dialog with the given screen dimensions
	Show(screenWidth, screenHeight int)

	// Hide closes the dialog
	Hide()

	// IsOpen returns whether the dialog is currently open
	IsOpen() bool

	// HandleInput processes a key event and returns true if the dialog handled it
	HandleInput(key tcell.Key, mod tcell.ModMask, ch rune) bool

	// Render renders the dialog to the given screen
	Render(screen Screen, style tcell.Style)

	// GetResult returns the dialog result (implementation-specific)
	GetResult() interface{}

	// IsConfirmed returns whether the dialog was confirmed (OK/Yes)
	IsConfirmed() bool
}

// Screen is the interface for screen operations used by dialogs.
// This is compatible with terminal.Screen interface.
type Screen interface {
	SetContent(x, y int, mainc rune, combc []rune, style tcell.Style) error
}

// BaseDialog provides common dialog functionality.
type BaseDialog struct {
	isOpen     bool
	confirmed  bool
	cancelled  bool
	screenW    int
	screenH    int
	width      int
	height     int
	x          int
	y          int
	title      string
	focusIndex int // Which element has focus (0=first button/input, 1=second, etc.)
}

// Show opens the dialog and calculates position.
func (d *BaseDialog) Show(screenWidth, screenHeight int) {
	d.isOpen = true
	d.confirmed = false
	d.cancelled = false
	d.screenW = screenWidth
	d.screenH = screenHeight
	d.focusIndex = 0

	// Center the dialog
	d.x = (screenWidth - d.width) / 2
	d.y = (screenHeight - d.height) / 2

	// Ensure dialog stays on screen
	if d.x < 0 {
		d.x = 0
	}
	if d.y < 0 {
		d.y = 0
	}
	if d.x+d.width > screenWidth {
		d.x = screenWidth - d.width
	}
	if d.y+d.height > screenHeight {
		d.y = screenHeight - d.height
	}
}

// Hide closes the dialog.
func (d *BaseDialog) Hide() {
	d.isOpen = false
}

// IsOpen returns whether the dialog is currently open.
func (d *BaseDialog) IsOpen() bool {
	return d.isOpen
}

// IsConfirmed returns whether the dialog was confirmed.
func (d *BaseDialog) IsConfirmed() bool {
	return d.confirmed
}

// IsCancelled returns whether the dialog was cancelled.
func (d *BaseDialog) IsCancelled() bool {
	return d.cancelled
}

// SetCancelled marks the dialog as cancelled.
func (d *BaseDialog) SetCancelled() {
	d.cancelled = true
	d.confirmed = false
	d.isOpen = false
}

// SetConfirmed marks the dialog as confirmed.
func (d *BaseDialog) SetConfirmed() {
	d.confirmed = true
	d.cancelled = false
	d.isOpen = false
}

// DrawBorder draws the dialog border with title.
func (d *BaseDialog) DrawBorder(screen Screen, style tcell.Style) {
	// Draw top border with title
	for x := 0; x < d.width; x++ {
		char := '─'
		if x == 0 {
			char = '┌'
		} else if x == d.width-1 {
			char = '┐'
		}
		screen.SetContent(d.x+x, d.y, char, []rune{}, style)
	}

	// Draw title if present
	if d.title != "" {
		titleStyle := style.Bold(true)
		titleX := d.x + 2
		for i, ch := range d.title {
			if titleX+i < d.x+d.width-2 {
				screen.SetContent(titleX+i, d.y, ch, []rune{}, titleStyle)
			}
		}
	}

	// Draw side borders
	for y := 1; y < d.height-1; y++ {
		screen.SetContent(d.x, d.y+y, '│', []rune{}, style)
		screen.SetContent(d.x+d.width-1, d.y+y, '│', []rune{}, style)
	}

	// Draw bottom border
	for x := 0; x < d.width; x++ {
		char := '─'
		if x == 0 {
			char = '└'
		} else if x == d.width-1 {
			char = '┘'
		}
		screen.SetContent(d.x+x, d.y+d.height-1, char, []rune{}, style)
	}
}

// DrawButton draws a button at the specified position.
func (d *BaseDialog) DrawButton(screen Screen, x, y, index int, label string, style tcell.Style, isFocused bool) {
	buttonStyle := style
	if isFocused {
		buttonStyle = style.Reverse(true).Bold(true)
	}

	// Draw button with brackets
	buttonText := fmt.Sprintf("[ %s ]", label)
	for i, ch := range buttonText {
		if x+i < d.x+d.width-1 {
			screen.SetContent(x+i, y, ch, []rune{}, buttonStyle)
		}
	}
}

// DrawText draws text within the dialog bounds.
func (d *BaseDialog) DrawText(screen Screen, x, y int, text string, style tcell.Style) {
	for i, ch := range text {
		if x+i < d.x+d.width-1 && x+i >= d.x+1 {
			screen.SetContent(x+i, y, ch, []rune{}, style)
		}
	}
}

// Clear clears the dialog area.
func (d *BaseDialog) Clear(screen Screen, bgStyle tcell.Style) {
	for y := 0; y < d.height; y++ {
		for x := 0; x < d.width; x++ {
			screen.SetContent(d.x+x, d.y+y, ' ', []rune{}, bgStyle)
		}
	}
}

// InputDialog is a dialog for text input with OK/Cancel buttons.
type InputDialog struct {
	BaseDialog
	prompt    string
	input     string
	onConfirm func(string)
	onCancel  func()
	cursorPos int
	maxLength int
}

// NewInputDialog creates a new input dialog.
// width should be at least 40 for good usability.
func NewInputDialog(title, prompt string, defaultValue string, onConfirm func(string), onCancel func()) *InputDialog {
	// Calculate dimensions
	promptWidth := len(prompt)
	inputWidth := 40
	if inputWidth < promptWidth {
		inputWidth = promptWidth + 5
	}
	if inputWidth > 60 {
		inputWidth = 60
	}

	d := &InputDialog{
		BaseDialog: BaseDialog{
			title:  title,
			width:  inputWidth + 4, // Padding on sides
			height: 6,              // Title + prompt + input + buttons + spacing
		},
		prompt:    prompt,
		input:     defaultValue,
		onConfirm: onConfirm,
		onCancel:  onCancel,
		cursorPos: len(defaultValue),
		maxLength: inputWidth,
	}

	return d
}

// HandleInput processes keyboard input for the dialog.
func (d *InputDialog) HandleInput(key tcell.Key, mod tcell.ModMask, ch rune) bool {
	switch key {
	case tcell.KeyEscape:
		d.SetCancelled()
		if d.onCancel != nil {
			d.onCancel()
		}
		return true

	case tcell.KeyEnter:
		d.SetConfirmed()
		if d.onConfirm != nil {
			d.onConfirm(d.input)
		}
		return true

	case tcell.KeyTab:
		// Cycle focus between input field (0) and buttons (1, 2)
		d.focusIndex = (d.focusIndex + 1) % 3
		return true

	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if d.focusIndex == 0 && len(d.input) > 0 && d.cursorPos > 0 {
			// Remove character before cursor
			if d.cursorPos >= len(d.input) {
				d.input = d.input[:len(d.input)-1]
			} else {
				d.input = d.input[:d.cursorPos-1] + d.input[d.cursorPos:]
			}
			d.cursorPos--
		}
		return true

	case tcell.KeyDelete:
		if d.focusIndex == 0 && d.cursorPos < len(d.input) {
			d.input = d.input[:d.cursorPos] + d.input[d.cursorPos+1:]
		}
		return true

	case tcell.KeyLeft:
		if d.focusIndex == 0 && d.cursorPos > 0 {
			d.cursorPos--
		}
		return true

	case tcell.KeyRight:
		if d.focusIndex == 0 && d.cursorPos < len(d.input) {
			d.cursorPos++
		}
		return true

	case tcell.KeyHome:
		if d.focusIndex == 0 {
			d.cursorPos = 0
		}
		return true

	case tcell.KeyEnd:
		if d.focusIndex == 0 {
			d.cursorPos = len(d.input)
		}
		return true

	case tcell.KeyRune:
		if d.focusIndex == 0 && ch != 0 && len(d.input) < d.maxLength {
			// Insert character at cursor position
			if d.cursorPos >= len(d.input) {
				d.input += string(ch)
			} else {
				d.input = d.input[:d.cursorPos] + string(ch) + d.input[d.cursorPos:]
			}
			d.cursorPos++
		}
		return true
	}

	return false
}

// Render draws the input dialog.
func (d *InputDialog) Render(screen Screen, style tcell.Style) {
	if !d.isOpen {
		return
	}

	// Clear dialog area
	d.Clear(screen, style)

	// Draw border
	d.DrawBorder(screen, style)

	// Draw prompt
	promptX := d.x + 2
	promptY := d.y + 2
	d.DrawText(screen, promptX, promptY, d.prompt, style)

	// Draw input field with border
	inputY := promptY + 1
	inputStartX := d.x + 2
	inputEndX := d.x + d.width - 2

	// Draw input field background
	inputStyle := style.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
	for x := inputStartX; x < inputEndX; x++ {
		screen.SetContent(x, inputY, ' ', []rune{}, inputStyle)
	}

	// Draw input text
	displayStart := 0
	if d.cursorPos > inputEndX-inputStartX-2 {
		displayStart = d.cursorPos - (inputEndX - inputStartX - 2)
	}
	displayText := d.input
	if displayStart > 0 {
		displayText = d.input[displayStart:]
	}
	if len(displayText) > inputEndX-inputStartX-2 {
		displayText = displayText[:inputEndX-inputStartX-2]
	}

	for i, ch := range displayText {
		screen.SetContent(inputStartX+i, inputY, ch, []rune{}, inputStyle)
	}

	// Draw cursor if input field is focused
	if d.focusIndex == 0 {
		cursorX := inputStartX + d.cursorPos - displayStart
		if cursorX < inputEndX {
			cursorStyle := style.Reverse(true)
			if d.cursorPos < len(d.input) {
				screen.SetContent(cursorX, inputY, rune(d.input[d.cursorPos]), []rune{}, cursorStyle)
			} else {
				screen.SetContent(cursorX, inputY, ' ', []rune{}, cursorStyle)
			}
		}
	}

	// Draw buttons
	buttonY := inputY + 2
	buttonWidth := 10 // "[ OK ]" or "[ Cancel ]"
	spacing := (d.width - 2*buttonWidth) / 3

	okX := d.x + spacing
	cancelX := okX + buttonWidth + spacing

	d.DrawButton(screen, okX, buttonY, 1, "OK", style, d.focusIndex == 1)
	d.DrawButton(screen, cancelX, buttonY, 2, "Cancel", style, d.focusIndex == 2)
}

// GetResult returns the input value.
func (d *InputDialog) GetResult() interface{} {
	return d.input
}

// GetInput returns the current input string.
func (d *InputDialog) GetInput() string {
	return d.input
}

// SetInput sets the input value programmatically.
func (d *InputDialog) SetInput(value string) {
	d.input = value
	d.cursorPos = len(value)
}

// ConfirmDialog is a dialog for yes/no confirmation.
type ConfirmDialog struct {
	BaseDialog
	message   string
	onConfirm func()
	onCancel  func()
}

// NewConfirmDialog creates a new confirmation dialog.
func NewConfirmDialog(title, message string, onConfirm func(), onCancel func()) *ConfirmDialog {
	// Calculate dimensions based on message
	lines := strings.Split(message, "\n")
	maxWidth := len(title)
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	// Add padding and button space
	width := maxWidth + 8
	if width < 40 {
		width = 40
	}

	height := len(lines) + 5 // Message lines + title + buttons + spacing

	return &ConfirmDialog{
		BaseDialog: BaseDialog{
			title:  title,
			width:  width,
			height: height,
		},
		message:   message,
		onConfirm: onConfirm,
		onCancel:  onCancel,
	}
}

// HandleInput processes keyboard input for the dialog.
func (d *ConfirmDialog) HandleInput(key tcell.Key, mod tcell.ModMask, ch rune) bool {
	switch key {
	case tcell.KeyEscape:
		d.SetCancelled()
		if d.onCancel != nil {
			d.onCancel()
		}
		return true

	case tcell.KeyEnter:
		if d.focusIndex == 0 || d.focusIndex == 1 {
			d.SetConfirmed()
			if d.onConfirm != nil {
				d.onConfirm()
			}
		} else {
			d.SetCancelled()
			if d.onCancel != nil {
				d.onCancel()
			}
		}
		return true

	case tcell.KeyTab:
		// Cycle between Yes (0) and No (1)
		d.focusIndex = (d.focusIndex + 1) % 2
		return true

	case tcell.KeyLeft, tcell.KeyRight:
		// Move between Yes and No
		d.focusIndex = (d.focusIndex + 1) % 2
		return true
	}

	return false
}

// Render draws the confirmation dialog.
func (d *ConfirmDialog) Render(screen Screen, style tcell.Style) {
	if !d.isOpen {
		return
	}

	// Clear dialog area
	d.Clear(screen, style)

	// Draw border
	d.DrawBorder(screen, style)

	// Draw message (may be multi-line)
	lines := strings.Split(d.message, "\n")
	messageStartY := d.y + 2
	for i, line := range lines {
		lineX := d.x + (d.width-len(line))/2 // Center text
		if lineX < d.x+1 {
			lineX = d.x + 1
		}
		d.DrawText(screen, lineX, messageStartY+i, line, style)
	}

	// Draw buttons
	buttonY := messageStartY + len(lines) + 1
	buttonWidth := 10
	spacing := (d.width - 2*buttonWidth) / 3

	yesX := d.x + spacing
	noX := yesX + buttonWidth + spacing

	d.DrawButton(screen, yesX, buttonY, 0, "Yes", style, d.focusIndex == 0)
	d.DrawButton(screen, noX, buttonY, 1, "No", style, d.focusIndex == 1)
}

// GetResult returns nil for confirmation dialog (use IsConfirmed).
func (d *ConfirmDialog) GetResult() interface{} {
	return nil
}

// GoToLineDialog is a specialized dialog for jumping to a line number.
type GoToLineDialog struct {
	*InputDialog
	maxLine int
}

// NewGoToLineDialog creates a new "Go to Line" dialog.
func NewGoToLineDialog(maxLine int, onConfirm func(int), onCancel func()) *GoToLineDialog {
	prompt := fmt.Sprintf("Line number (1-%d):", maxLine)

	inputDlg := NewInputDialog(
		"Go to Line",
		prompt,
		"",
		func(lineStr string) {
			var lineNum int
			fmt.Sscanf(lineStr, "%d", &lineNum)
			if lineNum < 1 {
				lineNum = 1
			}
			if lineNum > maxLine {
				lineNum = maxLine
			}
			onConfirm(lineNum)
		},
		onCancel,
	)

	return &GoToLineDialog{
		InputDialog: inputDlg,
		maxLine:     maxLine,
	}
}

// OpenFileDialog is a specialized dialog for opening files.
type OpenFileDialog struct {
	*InputDialog
}

// NewOpenFileDialog creates a new "Open File" dialog.
func NewOpenFileDialog(defaultPath string, onConfirm func(string), onCancel func()) *OpenFileDialog {
	inputDlg := NewInputDialog(
		"Open File",
		"File path:",
		defaultPath,
		onConfirm,
		onCancel,
	)

	return &OpenFileDialog{
		InputDialog: inputDlg,
	}
}

// SaveAsDialog is a specialized dialog for saving files.
type SaveAsDialog struct {
	*InputDialog
}

// NewSaveAsDialog creates a new "Save As" dialog.
func NewSaveAsDialog(defaultPath string, onConfirm func(string), onCancel func()) *SaveAsDialog {
	inputDlg := NewInputDialog(
		"Save As",
		"File path:",
		defaultPath,
		onConfirm,
		onCancel,
	)

	return &SaveAsDialog{
		InputDialog: inputDlg,
	}
}

// UnsavedChangesDialog is a specialized confirmation dialog for unsaved changes.
type UnsavedChangesDialog struct {
	*ConfirmDialog
}

// NewUnsavedChangesDialog creates a dialog for unsaved changes.
func NewUnsavedChangesDialog(filename string, onSave func(), onDiscard func(), onCancel func()) *UnsavedChangesDialog {
	var message string
	if filename != "" {
		message = fmt.Sprintf("'%s' has unsaved changes.\nDo you want to save before closing?", filename)
	} else {
		message = "This file has unsaved changes.\nDo you want to save before closing?"
	}

	// Override button labels
	confirmDlg := &ConfirmDialog{
		BaseDialog: BaseDialog{
			title:  "Unsaved Changes",
			width:  50,
			height: 6,
		},
		message:   message,
		onConfirm: onSave,
		onCancel:  onCancel,
	}

	return &UnsavedChangesDialog{
		ConfirmDialog: confirmDlg,
	}
}

// DialogManager manages a stack of dialogs.
type DialogManager struct {
	dialogs []Dialog
}

// NewDialogManager creates a new dialog manager.
func NewDialogManager() *DialogManager {
	return &DialogManager{
		dialogs: make([]Dialog, 0),
	}
}

// Push adds a dialog to the stack and opens it.
func (dm *DialogManager) Push(d Dialog, screenW, screenH int) {
	dm.dialogs = append(dm.dialogs, d)
	d.Show(screenW, screenH)
}

// Pop removes and returns the top dialog.
func (dm *DialogManager) Pop() Dialog {
	if len(dm.dialogs) == 0 {
		return nil
	}

	d := dm.dialogs[len(dm.dialogs)-1]
	dm.dialogs = dm.dialogs[:len(dm.dialogs)-1]
	d.Hide()
	return d
}

// Peek returns the top dialog without removing it.
func (dm *DialogManager) Peek() Dialog {
	if len(dm.dialogs) == 0 {
		return nil
	}
	return dm.dialogs[len(dm.dialogs)-1]
}

// IsEmpty returns whether there are no open dialogs.
func (dm *DialogManager) IsEmpty() bool {
	return len(dm.dialogs) == 0
}

// HandleInput routes input to the top dialog.
func (dm *DialogManager) HandleInput(key tcell.Key, mod tcell.ModMask, ch rune) bool {
	if dm.IsEmpty() {
		return false
	}

	d := dm.Peek()
	handled := d.HandleInput(key, mod, ch)

	// If dialog closed, pop it
	if !d.IsOpen() {
		dm.Pop()
	}

	return handled
}

// Render renders all open dialogs (top one last = on top).
func (dm *DialogManager) Render(screen Screen, style tcell.Style) {
	for _, d := range dm.dialogs {
		d.Render(screen, style)
	}
}

// HasOpenDialog returns whether any dialog is currently open.
func (dm *DialogManager) HasOpenDialog() bool {
	for _, d := range dm.dialogs {
		if d.IsOpen() {
			return true
		}
	}
	return false
}
