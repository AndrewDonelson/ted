// Package editor implements the main editor controller.
//
// It coordinates all components (buffer, file, renderer, terminal) and
// manages the main event loop.
package editor

import (
	"fmt"

	"github.com/AndrewDonelson/ted/core/buffer"
	"github.com/AndrewDonelson/ted/core/clipboard"
	"github.com/AndrewDonelson/ted/core/file"
	"github.com/AndrewDonelson/ted/core/history"
	"github.com/AndrewDonelson/ted/ui/layout"
	"github.com/AndrewDonelson/ted/ui/menu"
	"github.com/AndrewDonelson/ted/ui/renderer"
	"github.com/AndrewDonelson/ted/ui/terminal"
	"github.com/gdamore/tcell/v2"
)

// EditorMode represents the editing mode.
type EditorMode int

const (
	// ModeInsert represents insert mode (default).
	ModeInsert EditorMode = iota
	// ModeOverwrite represents overwrite mode.
	ModeOverwrite
)

// Editor represents the main editor instance.
type Editor struct {
	// Core components
	buffer  *buffer.Buffer
	file    *FileState
	history *history.History

	// UI components
	layout   *layout.Layout
	renderer *renderer.Renderer
	menuBar  *menu.MenuBar
	screen   terminal.Screen

	// State
	mode       EditorMode
	isDirty    bool
	filePath   string
	fileInfo   *file.FileInfo
	lineEnding file.LineEnding

	// Selection state
	selectionStart buffer.Position // Start of selection (anchor point)
	selectionEnd   buffer.Position // End of selection (cursor position)
	hasSelection   bool            // Whether there is an active selection
}

// FileState tracks file-related state.
type FileState struct {
	Path       string
	LineEnding file.LineEnding
	Encoding   string
}

// NewEditor creates a new editor instance.
func NewEditor() (*Editor, error) {
	// Initialize terminal screen
	screen, err := terminal.NewScreen()
	if err != nil {
		return nil, fmt.Errorf("initialize screen: %w", err)
	}

	// Get screen dimensions
	width, height := screen.GetSize()

	// Initialize layout
	layout := layout.NewLayout(width, height)

	// Initialize renderer
	renderer := renderer.NewRenderer(screen, layout)

	// Initialize menu bar
	menuBar := menu.NewMenuBar()

	// Initialize buffer
	buf := buffer.NewBuffer()

	// Initialize history (undo/redo)
	hist := history.NewHistory(100) // 100 operations deep

	return &Editor{
		buffer:         buf,
		file:           &FileState{Encoding: "UTF-8"},
		history:        hist,
		layout:         layout,
		renderer:       renderer,
		menuBar:        menuBar,
		screen:         screen,
		mode:           ModeInsert,
		isDirty:        false,
		lineEnding:     file.LineEndingLF,
		hasSelection:   false,
		selectionStart: buffer.Position{Line: 0, Col: 0},
		selectionEnd:   buffer.Position{Line: 0, Col: 0},
	}, nil
}

// OpenFile opens a file and loads it into the buffer.
func (e *Editor) OpenFile(path string) error {
	lines, fileInfo, err := file.ReadFileWithInfo(path)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	e.buffer.SetLines(lines)
	e.buffer.MarkSaved() // File is loaded, not modified
	e.filePath = path
	e.fileInfo = fileInfo
	e.lineEnding = fileInfo.LineEnding
	e.isDirty = false

	// Clear history when opening a new file
	e.history.Clear()

	return nil
}

// SetFilePath sets the file path for a new file (file doesn't exist yet).
func (e *Editor) SetFilePath(path string) {
	e.filePath = path
	e.fileInfo = nil                 // No file info for new files
	e.lineEnding = file.LineEndingLF // Default to LF for new files
	e.buffer.MarkSaved()             // New file starts as "saved" (empty)
	e.isDirty = false
}

// SaveFile saves the current buffer to the file.
func (e *Editor) SaveFile() error {
	if e.filePath == "" {
		return fmt.Errorf("no file path set")
	}

	lines := e.buffer.GetAllLines()
	if err := file.WriteFile(e.filePath, lines, e.lineEnding); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	// Mark buffer as saved
	e.buffer.MarkSaved()
	e.isDirty = false

	// Clear redo stack on save (save is a checkpoint)
	// Keep undo stack so user can still undo after save
	e.history.ClearRedo()

	// Update file info after save
	if e.fileInfo != nil {
		// Update size
		var totalSize int64
		for _, line := range lines {
			totalSize += int64(len(line))
		}
		// Add line ending sizes
		lineEndingSize := int64(len(string(e.lineEnding)))
		if len(lines) > 0 {
			totalSize += lineEndingSize * int64(len(lines)-1)
		}
		e.fileInfo.Size = totalSize
	} else {
		// Create file info for new files
		var totalSize int64
		for _, line := range lines {
			totalSize += int64(len(line))
		}
		lineEndingSize := int64(len(string(e.lineEnding)))
		if len(lines) > 0 {
			totalSize += lineEndingSize * int64(len(lines)-1)
		}
		e.fileInfo = &file.FileInfo{
			Size:       totalSize,
			LineEnding: e.lineEnding,
		}
	}

	return nil
}

// Run starts the main event loop.
func (e *Editor) Run() error {
	defer e.screen.Fini()

	// Initial render
	if err := e.render(); err != nil {
		return fmt.Errorf("initial render: %w", err)
	}

	// Event loop
	for {
		ev := e.screen.PollEvent()

		// Handle nil events (shouldn't happen, but be safe)
		if ev == nil {
			continue
		}

		// Handle resize events
		if resizeEv, ok := ev.(*tcell.EventResize); ok {
			width, height := resizeEv.Size()
			e.layout.AdjustForResize(width, height)
			if err := e.render(); err != nil {
				return fmt.Errorf("render after resize: %w", err)
			}
			continue
		}

		// Process keyboard events
		keyEvent := terminal.ProcessEvent(ev)
		if keyEvent == nil {
			continue
		}

		// Handle key actions
		if err := e.handleKeyEvent(keyEvent); err != nil {
			if err == ErrQuit {
				break
			}
			return fmt.Errorf("handle key event: %w", err)
		}

		// Render after handling event (unless it was a no-op)
		if keyEvent.Action != terminal.KeyActionNone {
			if err := e.render(); err != nil {
				return fmt.Errorf("render: %w", err)
			}
		}
	}

	return nil
}

// handleKeyEvent processes a key event and updates the editor state.
func (e *Editor) handleKeyEvent(ke *terminal.KeyEvent) error {
	switch ke.Action {
	case terminal.KeyActionQuit:
		return ErrQuit
	case terminal.KeyActionSave:
		if e.filePath != "" {
			if err := e.SaveFile(); err != nil {
				// Return error so user knows save failed
				return fmt.Errorf("save file: %w", err)
			}
		}
		// If no file path, silently ignore (Save As not implemented in Phase 0)
	case terminal.KeyActionCharacter:
		if ke.IsPrintable() {
			e.clearSelection() // Clear selection when typing
			e.insertCharacter(ke.Character)
		}
	case terminal.KeyActionMoveLeft:
		e.clearSelection()
		e.buffer.MoveCursorLeft()
	case terminal.KeyActionMoveRight:
		e.clearSelection()
		e.buffer.MoveCursorRight()
	case terminal.KeyActionMoveUp:
		e.clearSelection()
		e.buffer.MoveCursorUp()
	case terminal.KeyActionMoveDown:
		e.clearSelection()
		e.buffer.MoveCursorDown()
	case terminal.KeyActionSelectLeft:
		e.startSelectionIfNeeded()
		e.buffer.MoveCursorLeft()
		e.updateSelectionEnd()
	case terminal.KeyActionSelectRight:
		e.startSelectionIfNeeded()
		e.buffer.MoveCursorRight()
		e.updateSelectionEnd()
	case terminal.KeyActionSelectUp:
		e.startSelectionIfNeeded()
		e.buffer.MoveCursorUp()
		e.updateSelectionEnd()
	case terminal.KeyActionSelectDown:
		e.startSelectionIfNeeded()
		e.buffer.MoveCursorDown()
		e.updateSelectionEnd()
	case terminal.KeyActionBackspace:
		e.handleBackspace()
	case terminal.KeyActionDelete:
		e.handleDelete()
	case terminal.KeyActionEnter:
		e.insertCharacter('\n')
	case terminal.KeyActionHome:
		e.buffer.MoveCursorToLineStart()
	case terminal.KeyActionEnd:
		e.buffer.MoveCursorToLineEnd()
	case terminal.KeyActionUndo:
		if err := e.Undo(); err != nil {
			// Silently ignore if no undo available
			return nil
		}
	case terminal.KeyActionRedo:
		if err := e.Redo(); err != nil {
			// Silently ignore if no redo available
			return nil
		}
	case terminal.KeyActionCut:
		if err := e.Cut(); err != nil {
			return fmt.Errorf("cut: %w", err)
		}
	case terminal.KeyActionCopy:
		if err := e.Copy(); err != nil {
			return fmt.Errorf("copy: %w", err)
		}
	case terminal.KeyActionPaste:
		if err := e.Paste(); err != nil {
			return fmt.Errorf("paste: %w", err)
		}
	}

	return nil
}

// insertCharacter inserts a character at the current cursor position.
func (e *Editor) insertCharacter(r rune) {
	pos := e.buffer.GetCursor()
	text := string(r)

	// Record operation for undo
	op := &history.InsertOperation{
		Pos:  pos,
		Text: text,
	}

	// Perform insertion
	if err := e.buffer.Insert(pos, text); err != nil {
		// Ignore insertion errors for now
		return
	}

	// Push to history
	e.history.Push(op)
}

// handleBackspace handles the backspace key.
func (e *Editor) handleBackspace() {
	pos := e.buffer.GetCursor()
	var start, end buffer.Position
	var deletedText string

	if pos.Col > 0 {
		// Delete character before cursor
		start = buffer.Position{Line: pos.Line, Col: pos.Col - 1}
		end = pos
	} else if pos.Line > 0 {
		// Join with previous line
		prevLineLen := 0
		if line, err := e.buffer.GetLine(pos.Line - 1); err == nil {
			prevLineLen = len(line)
		}
		start = buffer.Position{Line: pos.Line - 1, Col: prevLineLen}
		end = buffer.Position{Line: pos.Line, Col: 0}
	} else {
		// At start of document, nothing to delete
		return
	}

	// Get text that will be deleted for undo
	var err error
	deletedText, err = e.buffer.GetText(start, end)
	if err != nil {
		return
	}

	// Record operation for undo
	op := &history.DeleteOperation{
		StartPos: start,
		EndPos:   end,
		Deleted:  deletedText,
	}

	// Perform deletion
	if err := e.buffer.Delete(start, end); err != nil {
		return
	}

	// Update cursor position
	if pos.Col > 0 {
		e.buffer.MoveCursorLeft()
	} else {
		e.buffer.MoveCursor(start)
	}

	// Push to history
	e.history.Push(op)
}

// handleDelete handles the delete key.
func (e *Editor) handleDelete() {
	pos := e.buffer.GetCursor()
	line, err := e.buffer.GetLine(pos.Line)
	if err != nil {
		return
	}

	var start, end buffer.Position
	var deletedText string

	if pos.Col < len(line) {
		// Delete character at cursor
		start = pos
		end = buffer.Position{Line: pos.Line, Col: pos.Col + 1}
	} else if pos.Line < e.buffer.LineCount()-1 {
		// Join with next line
		start = pos
		end = buffer.Position{Line: pos.Line + 1, Col: 0}
	} else {
		// At end of document, nothing to delete
		return
	}

	// Get text that will be deleted for undo
	deletedText, err = e.buffer.GetText(start, end)
	if err != nil {
		return
	}

	// Record operation for undo
	op := &history.DeleteOperation{
		StartPos: start,
		EndPos:   end,
		Deleted:  deletedText,
	}

	// Perform deletion
	if err := e.buffer.Delete(start, end); err != nil {
		return
	}

	// Push to history
	e.history.Push(op)
}

// Undo undoes the last operation.
func (e *Editor) Undo() error {
	return e.history.Undo(e.buffer)
}

// Redo redoes the last undone operation.
func (e *Editor) Redo() error {
	return e.history.Redo(e.buffer)
}

// Copy copies the selected text (or current line if no selection) to clipboard.
func (e *Editor) Copy() error {
	var text string
	var err error

	if e.hasSelection {
		// Copy selected text
		start, end := e.getSelectionRange()
		text, err = e.buffer.GetText(start, end)
		if err != nil {
			return fmt.Errorf("get selected text: %w", err)
		}
	} else {
		// Copy current line if no selection
		pos := e.buffer.GetCursor()
		text, err = e.buffer.GetLine(pos.Line)
		if err != nil {
			return fmt.Errorf("get line: %w", err)
		}
	}

	// Copy to clipboard
	if err := clipboard.Write(text); err != nil {
		return fmt.Errorf("write clipboard: %w", err)
	}

	return nil
}

// Cut cuts the selected text (or current line if no selection) to clipboard.
func (e *Editor) Cut() error {
	var start, end buffer.Position
	var deletedText string
	var err error

	if e.hasSelection {
		// Cut selected text
		start, end = e.getSelectionRange()
		deletedText, err = e.buffer.GetText(start, end)
		if err != nil {
			return fmt.Errorf("get selected text: %w", err)
		}
	} else {
		// Cut current line if no selection
		pos := e.buffer.GetCursor()
		line, err := e.buffer.GetLine(pos.Line)
		if err != nil {
			return fmt.Errorf("get line: %w", err)
		}

		start = buffer.Position{Line: pos.Line, Col: 0}
		end = buffer.Position{Line: pos.Line + 1, Col: 0}
		if pos.Line == e.buffer.LineCount()-1 {
			// Last line - delete to end
			end = buffer.Position{Line: pos.Line, Col: len(line)}
		}

		deletedText, err = e.buffer.GetText(start, end)
		if err != nil {
			return fmt.Errorf("get text: %w", err)
		}
	}

	// Copy to clipboard
	if err := clipboard.Write(deletedText); err != nil {
		return fmt.Errorf("write clipboard: %w", err)
	}

	// Record operation for undo
	op := &history.DeleteOperation{
		StartPos: start,
		EndPos:   end,
		Deleted:  deletedText,
	}

	// Perform deletion
	if err := e.buffer.Delete(start, end); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	// Clear selection
	e.clearSelection()

	// Adjust cursor
	e.buffer.MoveCursor(start)

	// Push to history
	e.history.Push(op)

	return nil
}

// Paste pastes text from clipboard at the current cursor position.
func (e *Editor) Paste() error {
	// Read from clipboard
	text, err := clipboard.Read()
	if err != nil {
		return fmt.Errorf("read clipboard: %w", err)
	}

	if text == "" {
		return nil // Nothing to paste
	}

	// Record operation for undo
	pos := e.buffer.GetCursor()
	op := &history.InsertOperation{
		Pos:  pos,
		Text: text,
	}

	// Insert text
	if err := e.buffer.Insert(pos, text); err != nil {
		return fmt.Errorf("insert: %w", err)
	}

	// Push to history
	e.history.Push(op)

	return nil
}

// clearSelection clears the current selection.
func (e *Editor) clearSelection() {
	e.hasSelection = false
}

// startSelectionIfNeeded starts a selection if one doesn't exist.
func (e *Editor) startSelectionIfNeeded() {
	if !e.hasSelection {
		e.hasSelection = true
		e.selectionStart = e.buffer.GetCursor()
	}
}

// updateSelectionEnd updates the end of the selection to the current cursor position.
func (e *Editor) updateSelectionEnd() {
	if e.hasSelection {
		e.selectionEnd = e.buffer.GetCursor()
	}
}

// getSelectionRange returns the normalized selection range (start <= end).
func (e *Editor) getSelectionRange() (start, end buffer.Position) {
	if !e.hasSelection {
		return buffer.Position{}, buffer.Position{}
	}

	start = e.selectionStart
	end = e.selectionEnd

	// Normalize: ensure start <= end
	if start.Line > end.Line || (start.Line == end.Line && start.Col > end.Col) {
		start, end = end, start
	}

	return start, end
}

// render renders all UI components.
func (e *Editor) render() error {
	cursorPos := e.buffer.GetCursor()

	// Build file info for info bar
	fileInfo := e.buildFileInfo()

	// Render everything
	return e.renderer.RenderAll(e.buffer, cursorPos, fileInfo)
}

// buildFileInfo builds the file info for the info bar.
func (e *Editor) buildFileInfo() *renderer.FileInfo {
	// Use buffer's modified flag as source of truth
	isModified := e.buffer.IsModified()

	info := &renderer.FileInfo{
		Name:       e.getFileName(),
		Path:       e.filePath,
		Encoding:   e.file.Encoding,
		LineEnding: string(e.lineEnding),
		TabSize:    4, // Default for Phase 0
		TotalLines: e.buffer.LineCount(),
		IsModified: isModified,
	}

	if e.fileInfo != nil {
		info.Size = e.fileInfo.Size
		info.Type = e.detectFileType()
	} else if e.filePath != "" {
		// For new files, still detect type from extension
		info.Type = e.detectFileType()
		info.Size = 0 // New file has no size yet
	}

	return info
}

// getFileName returns the filename for display.
func (e *Editor) getFileName() string {
	if e.filePath == "" {
		return ""
	}
	// Extract just the filename from path
	lastSlash := -1
	for i := len(e.filePath) - 1; i >= 0; i-- {
		if e.filePath[i] == '/' || e.filePath[i] == '\\' {
			lastSlash = i
			break
		}
	}
	if lastSlash >= 0 {
		return e.filePath[lastSlash+1:]
	}
	return e.filePath
}

// detectFileType detects the file type from the extension.
func (e *Editor) detectFileType() string {
	if e.filePath == "" {
		return ""
	}
	// Simple detection based on extension
	ext := ""
	for i := len(e.filePath) - 1; i >= 0; i-- {
		if e.filePath[i] == '.' {
			ext = e.filePath[i:]
			break
		}
		if e.filePath[i] == '/' || e.filePath[i] == '\\' {
			break
		}
	}

	switch ext {
	case ".go":
		return "Go"
	case ".js", ".jsx":
		return "JavaScript"
	case ".ts", ".tsx":
		return "TypeScript"
	case ".py":
		return "Python"
	case ".md":
		return "Markdown"
	case ".txt":
		return "Plain Text"
	default:
		return "Plain Text"
	}
}

// ErrQuit is returned when the user quits the editor.
var ErrQuit = fmt.Errorf("quit")
