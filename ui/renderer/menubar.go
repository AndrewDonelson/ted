// Package renderer implements menu bar rendering.
package renderer

import (
	"fmt"

	"github.com/AndrewDonelson/ted/ui/menu"
	"github.com/gdamore/tcell/v2"
)

// RenderMenuBar renders the static menu bar at the top of the screen.
// It displays "File Edit Search View Help" on the left and status indicators on the right.
func (r *Renderer) RenderMenuBar() error {
	region := r.layout.GetMenuBarRegion()
	style := GetMenuBarStyle()

	// Fill entire menu bar region with background color first
	for x := 0; x < region.Width; x++ {
		r.screen.SetContent(region.X+x, region.Y, ' ', nil, style)
	}

	// Menu items
	menuItems := "File  Edit  Search  View  Help"
	menuX := 0

	// Render menu items
	for i, char := range menuItems {
		if menuX+i < region.Width {
			r.screen.SetContent(menuX+i, region.Y, char, nil, style)
		}
	}

	// Status indicators on the right
	statusText := "INS │ UTF-8 │ LN 1, COL 1"
	statusX := region.Width - len(statusText)

	// Render status (only if there's space)
	if statusX > len(menuItems) && statusX >= 0 {
		for i, char := range statusText {
			if statusX+i < region.Width {
				r.screen.SetContent(statusX+i, region.Y, char, nil, style)
			}
		}
	}

	return nil
}

// RenderMenuBarWithStatus renders the menu bar with custom status information.
func (r *Renderer) RenderMenuBarWithStatus(mode string, encoding string, line, col int) error {
	region := r.layout.GetMenuBarRegion()
	style := GetMenuBarStyle()

	// Fill entire menu bar region with background color first
	for x := 0; x < region.Width; x++ {
		r.screen.SetContent(region.X+x, region.Y, ' ', nil, style)
	}

	// Menu items
	menuItems := "File  Edit  Search  View  Help"
	menuX := 0

	// Render menu items
	for i, char := range menuItems {
		if menuX+i < region.Width {
			r.screen.SetContent(menuX+i, region.Y, char, nil, style)
		}
	}

	// Build status text
	statusText := formatStatus(mode, encoding, line, col)
	statusX := region.Width - len(statusText)

	// Render status (only if there's space)
	if statusX > len(menuItems) && statusX >= 0 {
		for i, char := range statusText {
			if statusX+i < region.Width {
				r.screen.SetContent(statusX+i, region.Y, char, nil, style)
			}
		}
	}

	return nil
}

// RenderInteractiveMenuBar renders the menu bar with highlighting for active menu.
func (r *Renderer) RenderInteractiveMenuBar(menuBar *menu.MenuBar) error {
	region := r.layout.GetMenuBarRegion()
	style := GetMenuBarStyle()
	activeStyle := GetMenuActiveStyle()

	// Fill entire menu bar region with background color first
	for x := 0; x < region.Width; x++ {
		r.screen.SetContent(region.X+x, region.Y, ' ', nil, style)
	}

	// Render each menu label
	menus := menuBar.GetMenus()
	for i, m := range menus {
		x := menuBar.GetMenuPosition(i)
		menuStyle := style

		// Highlight active menu
		if menuBar.IsOpen() && menuBar.GetActiveMenu() == i {
			menuStyle = activeStyle
		}

		// Render menu label with underlined first character
		for j, char := range m.Label {
			if x+j < region.Width {
				charStyle := menuStyle
				if j == 0 {
					// Underline the hotkey character
					charStyle = charStyle.Underline(true)
				}
				r.screen.SetContent(x+j, region.Y, char, nil, charStyle)
			}
		}
	}

	return nil
}

// RenderDropdownMenu renders the dropdown menu separately (called after text area)
func (r *Renderer) RenderDropdownMenu(menuBar *menu.MenuBar) error {
	if !menuBar.IsOpen() || menuBar.GetActiveMenu() < 0 {
		return nil
	}

	return r.renderDropdownMenu(menuBar)
}

// renderDropdownMenu is the internal implementation that renders the dropdown
func (r *Renderer) renderDropdownMenu(menuBar *menu.MenuBar) error {
	normalStyle := GetDropdownStyle()
	selectedStyle := GetDropdownSelectedStyle()
	separatorStyle := GetDropdownSeparatorStyle()
	shortcutStyle := GetDropdownShortcutStyle()

	menus := menuBar.GetMenus()
	activeMenuIndex := menuBar.GetActiveMenu()
	if activeMenuIndex < 0 || activeMenuIndex >= len(menus) {
		return nil
	}

	activeMenu := menus[activeMenuIndex]
	activeItemIndex := menuBar.GetActiveItem()

	// Calculate dropdown position and size
	x := menuBar.GetMenuPosition(activeMenuIndex)
	y := 1 // Below menu bar

	// Calculate width based on longest item
	width := len(activeMenu.Label)
	for _, item := range activeMenu.Items {
		itemWidth := len(item.Label) + len(item.Shortcut) + 4 // padding + spacing
		if itemWidth > width {
			width = itemWidth
		}
	}
	if width < 20 {
		width = 20
	}

	// Render each menu item
	for i, item := range activeMenu.Items {
		itemY := y + i
		itemStyle := normalStyle

		if i == activeItemIndex && !item.IsSeparator {
			itemStyle = selectedStyle
		}

		// Fill item background
		for px := 0; px < width; px++ {
			r.screen.SetContent(x+px, itemY, ' ', nil, itemStyle)
		}

		if item.IsSeparator {
			// Render separator line
			for px := 1; px < width-1; px++ {
				r.screen.SetContent(x+px, itemY, '─', nil, separatorStyle)
			}
		} else {
			// Render item label
			labelX := x + 2 // Left padding
			for j, char := range item.Label {
				r.screen.SetContent(labelX+j, itemY, char, nil, itemStyle)
			}

			// Render shortcut on the right
			if item.Shortcut != "" {
				shortcutX := x + width - len(item.Shortcut) - 2
				scStyle := shortcutStyle
				if i == activeItemIndex {
					scStyle = itemStyle // Use selected style for shortcut too
				}
				for j, char := range item.Shortcut {
					r.screen.SetContent(shortcutX+j, itemY, char, nil, scStyle)
				}
			}
		}
	}

	// Draw border around dropdown
	borderStyle := GetDropdownBorderStyle()
	height := len(activeMenu.Items)

	// Top border
	r.screen.SetContent(x-1, y-1, '┌', nil, borderStyle)
	for px := 0; px < width; px++ {
		r.screen.SetContent(x+px, y-1, '─', nil, borderStyle)
	}
	r.screen.SetContent(x+width, y-1, '┐', nil, borderStyle)

	// Side borders
	for py := 0; py < height; py++ {
		r.screen.SetContent(x-1, y+py, '│', nil, borderStyle)
		r.screen.SetContent(x+width, y+py, '│', nil, borderStyle)
	}

	// Bottom border
	r.screen.SetContent(x-1, y+height, '└', nil, borderStyle)
	for px := 0; px < width; px++ {
		r.screen.SetContent(x+px, y+height, '─', nil, borderStyle)
	}
	r.screen.SetContent(x+width, y+height, '┘', nil, borderStyle)

	return nil
}

// GetMenuActiveStyle returns the style for the active/highlighted menu.
func GetMenuActiveStyle() tcell.Style {
	return tcell.StyleDefault.
		Foreground(tcell.ColorBlack).
		Background(tcell.ColorWhite)
}

// GetDropdownStyle returns the style for dropdown menu items.
func GetDropdownStyle() tcell.Style {
	return tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorNavy)
}

// GetDropdownSelectedStyle returns the style for selected dropdown items.
func GetDropdownSelectedStyle() tcell.Style {
	return tcell.StyleDefault.
		Foreground(tcell.ColorBlack).
		Background(tcell.ColorAqua)
}

// GetDropdownSeparatorStyle returns the style for dropdown separators.
func GetDropdownSeparatorStyle() tcell.Style {
	return tcell.StyleDefault.
		Foreground(tcell.ColorGray).
		Background(tcell.ColorNavy)
}

// GetDropdownShortcutStyle returns the style for shortcut text.
func GetDropdownShortcutStyle() tcell.Style {
	return tcell.StyleDefault.
		Foreground(tcell.ColorSilver).
		Background(tcell.ColorNavy)
}

// GetDropdownBorderStyle returns the style for dropdown borders.
func GetDropdownBorderStyle() tcell.Style {
	return tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorNavy)
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
