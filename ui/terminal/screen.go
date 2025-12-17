// Package terminal implements terminal interface using tcell.
//
// It provides screen management, event handling, and terminal state
// management for the editor.
package terminal

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

// Screen represents a terminal screen interface.
type Screen interface {
	// Clear clears the entire screen.
	Clear()

	// Refresh updates the display with any pending changes.
	Refresh() error

	// GetSize returns the current screen dimensions (width, height).
	GetSize() (width, height int)

	// SetContent sets the content at a specific position.
	SetContent(x, y int, mainc rune, combc []rune, style tcell.Style) error

	// ShowCursor sets the cursor position.
	ShowCursor(x, y int)

	// HideCursor hides the cursor.
	HideCursor()

	// PollEvent waits for and returns the next event.
	PollEvent() tcell.Event

	// Fini finalizes the screen and restores the terminal state.
	Fini()
}

// TCellScreen wraps tcell.Screen to implement our Screen interface.
type TCellScreen struct {
	screen tcell.Screen
}

// NewScreen creates and initializes a new terminal screen.
// Returns an error if the screen cannot be initialized.
func NewScreen() (*TCellScreen, error) {
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, fmt.Errorf("create screen: %w", err)
	}

	if err := s.Init(); err != nil {
		return nil, fmt.Errorf("init screen: %w", err)
	}

	// Set default style
	s.SetStyle(tcell.StyleDefault)

	// Clear the screen
	s.Clear()

	return &TCellScreen{screen: s}, nil
}

// Clear clears the entire screen.
func (s *TCellScreen) Clear() {
	s.screen.Clear()
}

// Refresh updates the display with any pending changes.
func (s *TCellScreen) Refresh() error {
	s.screen.Show()
	return nil
}

// GetSize returns the current screen dimensions (width, height).
func (s *TCellScreen) GetSize() (width, height int) {
	return s.screen.Size()
}

// SetContent sets the content at a specific position.
func (s *TCellScreen) SetContent(x, y int, mainc rune, combc []rune, style tcell.Style) error {
	s.screen.SetContent(x, y, mainc, combc, style)
	return nil
}

// ShowCursor sets the cursor position.
func (s *TCellScreen) ShowCursor(x, y int) {
	s.screen.ShowCursor(x, y)
}

// HideCursor hides the cursor.
func (s *TCellScreen) HideCursor() {
	s.screen.HideCursor()
}

// PollEvent waits for and returns the next event.
func (s *TCellScreen) PollEvent() tcell.Event {
	return s.screen.PollEvent()
}

// Fini finalizes the screen and restores the terminal state.
func (s *TCellScreen) Fini() {
	s.screen.Fini()
}

// GetRawScreen returns the underlying tcell.Screen for advanced operations.
// This should be used sparingly and only when necessary.
func (s *TCellScreen) GetRawScreen() tcell.Screen {
	return s.screen
}

