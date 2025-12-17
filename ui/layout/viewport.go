// Package layout implements layout calculations for the editor UI.
//
// It manages screen regions (menu bar, edit area, info bar) and
// viewport calculations for scrolling.
package layout

// Region represents a rectangular region on the screen.
type Region struct {
	X      int // Top-left X coordinate (0-indexed)
	Y      int // Top-left Y coordinate (0-indexed)
	Width  int // Width in characters
	Height int // Height in lines
}

// Viewport represents the visible portion of the buffer.
type Viewport struct {
	StartLine int // First visible line (0-indexed)
	EndLine   int // Last visible line (0-indexed, inclusive)
	OffsetX   int // Horizontal scroll offset in characters
	Width     int // Viewport width in characters
	Height    int // Viewport height in lines
}

// Layout manages the screen layout and viewport calculations.
type Layout struct {
	width      int
	height     int
	menuHeight int // Height of menu bar (typically 1)
	infoHeight int // Height of info bar (typically 1)
}

// NewLayout creates a new layout with the given screen dimensions.
func NewLayout(width, height int) *Layout {
	return &Layout{
		width:      width,
		height:     height,
		menuHeight: 1, // Menu bar takes 1 line
		infoHeight: 1, // Info bar takes 1 line
	}
}

// AdjustForResize updates the layout dimensions for a terminal resize.
func (l *Layout) AdjustForResize(newWidth, newHeight int) {
	l.width = newWidth
	l.height = newHeight
}

// GetMenuBarRegion returns the region for the menu bar.
func (l *Layout) GetMenuBarRegion() Region {
	return Region{
		X:      0,
		Y:      0,
		Width:  l.width,
		Height: l.menuHeight,
	}
}

// GetEditAreaRegion returns the region for the editable text area.
func (l *Layout) GetEditAreaRegion() Region {
	editY := l.menuHeight
	editHeight := l.height - l.menuHeight - l.infoHeight

	// Ensure minimum height
	if editHeight < 1 {
		editHeight = 1
	}

	return Region{
		X:      0,
		Y:      editY,
		Width:  l.width,
		Height: editHeight,
	}
}

// GetInfoBarRegion returns the region for the info bar at the bottom.
func (l *Layout) GetInfoBarRegion() Region {
	infoY := l.height - l.infoHeight

	// Ensure info bar is visible
	if infoY < 0 {
		infoY = 0
	}

	return Region{
		X:      0,
		Y:      infoY,
		Width:  l.width,
		Height: l.infoHeight,
	}
}

// GetLineNumberWidth calculates the width needed for line numbers.
// It considers the total number of lines in the buffer.
func (l *Layout) GetLineNumberWidth(totalLines int) int {
	if totalLines < 1 {
		return 0
	}

	// Calculate digits needed
	digits := 1
	n := totalLines
	for n >= 10 {
		n /= 10
		digits++
	}

	// Add padding: digits + 1 space + 1 separator = digits + 2
	return digits + 2
}

// CalculateViewport calculates the viewport based on cursor position and total lines.
// It ensures the cursor is visible and centers it if possible.
func (l *Layout) CalculateViewport(cursorLine, totalLines int) Viewport {
	editRegion := l.GetEditAreaRegion()
	viewportHeight := editRegion.Height

	// Handle empty buffer
	if totalLines == 0 {
		return Viewport{
			StartLine: 0,
			EndLine:   0,
			OffsetX:   0,
			Width:     editRegion.Width,
			Height:    viewportHeight,
		}
	}

	// Clamp cursor line to valid range
	if cursorLine < 0 {
		cursorLine = 0
	}
	if cursorLine >= totalLines {
		cursorLine = totalLines - 1
	}

	// Calculate start line to keep cursor visible
	startLine := cursorLine - viewportHeight/2
	if startLine < 0 {
		startLine = 0
	}

	// Calculate end line
	endLine := startLine + viewportHeight - 1
	if endLine >= totalLines {
		endLine = totalLines - 1
		// Adjust start line if we're at the end
		startLine = endLine - viewportHeight + 1
		if startLine < 0 {
			startLine = 0
		}
	}

	return Viewport{
		StartLine: startLine,
		EndLine:   endLine,
		OffsetX:   0, // Horizontal scrolling not implemented in Phase 0
		Width:     editRegion.Width,
		Height:    viewportHeight,
	}
}

// ScreenToBuffer converts screen coordinates to buffer position.
// Returns the line number and column in the buffer.
func (l *Layout) ScreenToBuffer(screenX, screenY int) (line, col int) {
	editRegion := l.GetEditAreaRegion()

	// Check if coordinates are in edit area
	if screenY < editRegion.Y || screenY >= editRegion.Y+editRegion.Height {
		return -1, -1
	}

	// Convert screen Y to buffer line (relative to viewport)
	bufferLine := screenY - editRegion.Y

	// Column is screen X (no horizontal scrolling in Phase 0)
	col = screenX

	return bufferLine, col
}

// BufferToScreen converts buffer position to screen coordinates.
// Returns screen X and Y coordinates, or -1, -1 if not visible.
func (l *Layout) BufferToScreen(bufferLine, bufferCol int, viewport Viewport) (screenX, screenY int) {
	editRegion := l.GetEditAreaRegion()

	// Check if line is in viewport
	if bufferLine < viewport.StartLine || bufferLine > viewport.EndLine {
		return -1, -1
	}

	// Calculate screen coordinates
	screenY = editRegion.Y + (bufferLine - viewport.StartLine)
	screenX = bufferCol + viewport.OffsetX

	// Check bounds
	if screenX < 0 || screenX >= editRegion.Width {
		return -1, -1
	}

	return screenX, screenY
}

// GetWidth returns the current layout width.
func (l *Layout) GetWidth() int {
	return l.width
}

// GetHeight returns the current layout height.
func (l *Layout) GetHeight() int {
	return l.height
}

// GetMinimumSize returns the minimum terminal size required.
func GetMinimumSize() (width, height int) {
	return 40, 10 // Minimum 40 columns, 10 rows
}

// IsSizeValid checks if the current size meets minimum requirements.
func (l *Layout) IsSizeValid() bool {
	minWidth, minHeight := GetMinimumSize()
	return l.width >= minWidth && l.height >= minHeight
}
