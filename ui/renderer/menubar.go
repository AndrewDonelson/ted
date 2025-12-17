// Package renderer implements menu bar rendering.
package renderer

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

// RenderMenuBar renders the static menu bar at the top of the screen.
// It displays "File Edit Search View Help" on the left and status indicators on the right.
func (r *Renderer) RenderMenuBar() error {
	region := r.layout.GetMenuBarRegion()
	style := GetMenuBarStyle()

	// Menu items
	menuItems := "File  Edit  Search  View  Help"
	menuX := 0

	// Render menu items
	for i, char := range menuItems {
		r.screen.SetContent(menuX+i, region.Y, char, nil, style)
	}

	// Status indicators on the right
	statusText := "INS │ UTF-8 │ LN 1, COL 1"
	statusX := region.Width - len(statusText)

	// Render status (only if there's space)
	if statusX > len(menuItems) {
		for i, char := range statusText {
			r.screen.SetContent(statusX+i, region.Y, char, nil, style)
		}
	}

	return nil
}

// RenderMenuBarWithStatus renders the menu bar with custom status information.
func (r *Renderer) RenderMenuBarWithStatus(mode string, encoding string, line, col int) error {
	region := r.layout.GetMenuBarRegion()
	style := GetMenuBarStyle()

	// Menu items
	menuItems := "File  Edit  Search  View  Help"
	menuX := 0

	// Render menu items
	for i, char := range menuItems {
		r.screen.SetContent(menuX+i, region.Y, char, nil, style)
	}

	// Build status text
	statusText := formatStatus(mode, encoding, line, col)
	statusX := region.Width - len(statusText)

	// Render status (only if there's space)
	if statusX > len(menuItems) {
		for i, char := range statusText {
			r.screen.SetContent(statusX+i, region.Y, char, nil, style)
		}
	}

	return nil
}

// formatStatus formats the status indicators for the menu bar.
func formatStatus(mode, encoding string, line, col int) string {
	lineStr := formatNumber(line + 1) // 1-indexed for display
	colStr := formatNumber(col + 1)   // 1-indexed for display
	return fmt.Sprintf("%s │ %s │ LN %s, COL %s", mode, encoding, lineStr, colStr)
}

// formatNumber formats a number as a string.
func formatNumber(n int) string {
	if n < 0 {
		n = 0
	}
	return fmt.Sprintf("%d", n)
}

// SetMenuBarContent allows setting custom menu bar content.
func (r *Renderer) SetMenuBarContent(content string, style tcell.Style) error {
	region := r.layout.GetMenuBarRegion()
	for i, char := range content {
		if i >= region.Width {
			break
		}
		r.screen.SetContent(i, region.Y, char, nil, style)
	}
	return nil
}
