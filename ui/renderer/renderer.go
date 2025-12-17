// Package renderer implements rendering components for the editor UI.
//
// It handles rendering of the menu bar, text area, info bar, and other
// visual elements to the terminal screen.
package renderer

import (
	"github.com/AndrewDonelson/ted/core/buffer"
	"github.com/AndrewDonelson/ted/ui/layout"
	"github.com/AndrewDonelson/ted/ui/terminal"
	"github.com/gdamore/tcell/v2"
)

// Renderer handles all rendering operations for the editor.
type Renderer struct {
	screen terminal.Screen
	layout *layout.Layout
}

// NewRenderer creates a new renderer with the given screen and layout.
func NewRenderer(screen terminal.Screen, layout *layout.Layout) *Renderer {
	return &Renderer{
		screen: screen,
		layout: layout,
	}
}

// Clear clears the entire screen.
func (r *Renderer) Clear() {
	r.screen.Clear()
}

// Refresh updates the display with any pending changes.
func (r *Renderer) Refresh() error {
	return r.screen.Refresh()
}

// RenderAll renders all UI components.
func (r *Renderer) RenderAll(buf *buffer.Buffer, cursorPos buffer.Position, fileInfo *FileInfo) error {
	r.Clear()

	// Fill entire screen with background color first
	if err := r.fillScreen(); err != nil {
		return err
	}

	// Render menu bar
	if err := r.RenderMenuBar(); err != nil {
		return err
	}

	// Render text area
	if err := r.RenderTextArea(buf, cursorPos); err != nil {
		return err
	}

	// Render info bar (CRITICAL: inverted colors)
	if err := r.RenderInfoBar(fileInfo); err != nil {
		return err
	}

	// Show cursor
	viewport := r.layout.CalculateViewport(cursorPos.Line, buf.LineCount())
	screenX, screenY := r.layout.BufferToScreen(cursorPos.Line, cursorPos.Col, viewport)
	if screenX >= 0 && screenY >= 0 {
		r.screen.ShowCursor(screenX, screenY)
	}

	return r.Refresh()
}

// fillScreen fills the entire screen with the default background color.
func (r *Renderer) fillScreen() error {
	screenWidth, screenHeight := r.screen.GetSize()
	defaultStyle := GetDefaultStyle()
	menuBarStyle := GetMenuBarStyle()
	infoBarStyle := GetInfoBarStyle()

	menuBarRegion := r.layout.GetMenuBarRegion()
	infoBarRegion := r.layout.GetInfoBarRegion()

	for y := 0; y < screenHeight; y++ {
		var style tcell.Style
		if y < menuBarRegion.Y+menuBarRegion.Height {
			// Menu bar area - use menu bar style
			style = menuBarStyle
		} else if y >= infoBarRegion.Y {
			// Info bar area - use inverted info bar style
			style = infoBarStyle
		} else {
			// Edit area - use default style
			style = defaultStyle
		}

		// Fill entire width with appropriate background
		for x := 0; x < screenWidth; x++ {
			if err := r.screen.SetContent(x, y, ' ', nil, style); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetDefaultStyle returns the default text style.
func GetDefaultStyle() tcell.Style {
	return tcell.StyleDefault.
		Foreground(tcell.Color252). // Light gray (#d4d4d4)
		Background(tcell.Color235)  // Dark gray (#1e1e1e)
}

// GetMenuBarStyle returns the style for the menu bar.
func GetMenuBarStyle() tcell.Style {
	return tcell.StyleDefault.
		Foreground(tcell.Color252). // Light gray
		Background(tcell.Color240)  // Slightly lighter dark gray (#252525)
}

// GetInfoBarStyle returns the INVERTED style for the info bar.
// CRITICAL: This must use inverted colors (light bg, dark text).
func GetInfoBarStyle() tcell.Style {
	return tcell.StyleDefault.
		Background(tcell.Color252). // Light gray background (#d4d4d4)
		Foreground(tcell.Color235)  // Dark gray text (#1e1e1e)
}

// GetLineNumberStyle returns the style for line numbers.
func GetLineNumberStyle() tcell.Style {
	return tcell.StyleDefault.
		Foreground(tcell.Color245). // Muted gray (#858585)
		Background(tcell.Color235)  // Dark gray
}

// GetCurrentLineStyle returns the style for the current line highlight.
func GetCurrentLineStyle() tcell.Style {
	return tcell.StyleDefault.
		Foreground(tcell.Color252).
		Background(tcell.Color240) // Subtle highlight (#2a2a2a)
}

// GetCursorStyle returns the style for the cursor.
func GetCursorStyle() tcell.Style {
	return tcell.StyleDefault.
		Foreground(tcell.Color235). // Dark background
		Background(tcell.Color255)  // White cursor (#ffffff)
}
