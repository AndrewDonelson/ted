// Package renderer implements text area rendering.
package renderer

import (
	"strconv"

	"github.com/AndrewDonelson/ted/core/buffer"
)

// RenderTextArea renders the buffer text in the edit area.
// It handles scrolling based on the viewport and highlights the current line.
func (r *Renderer) RenderTextArea(buf *buffer.Buffer, cursorPos buffer.Position) error {
	editRegion := r.layout.GetEditAreaRegion()
	viewport := r.layout.CalculateViewport(cursorPos.Line, buf.LineCount())

	defaultStyle := GetDefaultStyle()
	currentLineStyle := GetCurrentLineStyle()

	// Render visible lines
	for viewLine := 0; viewLine < viewport.Height; viewLine++ {
		bufferLine := viewport.StartLine + viewLine

		// Check if we've exceeded the buffer
		if bufferLine >= buf.LineCount() {
			// Fill remaining lines with empty space
			for x := 0; x < editRegion.Width; x++ {
				r.screen.SetContent(editRegion.X+x, editRegion.Y+viewLine, ' ', nil, defaultStyle)
			}
			continue
		}

		// Get line content
		lineText, err := buf.GetLine(bufferLine)
		if err != nil {
			// Skip invalid lines
			continue
		}

		// Determine style (highlight current line)
		lineStyle := defaultStyle
		if bufferLine == cursorPos.Line {
			lineStyle = currentLineStyle
		}

		// Render line content
		x := editRegion.X
		for i, char := range lineText {
			if i >= editRegion.Width {
				break // Line too long, truncate
			}
			r.screen.SetContent(x+i, editRegion.Y+viewLine, char, nil, lineStyle)
		}

		// Fill remaining space in line with background
		for x := len(lineText); x < editRegion.Width; x++ {
			r.screen.SetContent(editRegion.X+x, editRegion.Y+viewLine, ' ', nil, lineStyle)
		}
	}

	return nil
}

// RenderTextAreaWithLineNumbers renders the text area with line numbers.
func (r *Renderer) RenderTextAreaWithLineNumbers(buf *buffer.Buffer, cursorPos buffer.Position, showLineNumbers bool) error {
	editRegion := r.layout.GetEditAreaRegion()
	viewport := r.layout.CalculateViewport(cursorPos.Line, buf.LineCount())

	defaultStyle := GetDefaultStyle()
	currentLineStyle := GetCurrentLineStyle()
	lineNumberStyle := GetLineNumberStyle()

	// Calculate line number width if enabled
	lineNumberWidth := 0
	if showLineNumbers {
		lineNumberWidth = r.layout.GetLineNumberWidth(buf.LineCount())
		// Adjust edit region to account for line numbers
		editRegion.X += lineNumberWidth
		editRegion.Width -= lineNumberWidth
	}

	// Render visible lines
	for viewLine := 0; viewLine < viewport.Height; viewLine++ {
		bufferLine := viewport.StartLine + viewLine

		// Render line numbers if enabled
		if showLineNumbers {
			lineNum := bufferLine + 1                                  // 1-indexed for display
			lineNumStr := formatLineNumber(lineNum, lineNumberWidth-2) // -2 for separator
			x := 0
			for i, char := range lineNumStr {
				if i >= lineNumberWidth-1 {
					break
				}
				r.screen.SetContent(x+i, editRegion.Y+viewLine, char, nil, lineNumberStyle)
			}
			// Render separator
			r.screen.SetContent(lineNumberWidth-2, editRegion.Y+viewLine, '│', nil, lineNumberStyle)
		}

		// Check if we've exceeded the buffer
		if bufferLine >= buf.LineCount() {
			// Fill remaining space with empty
			for x := 0; x < editRegion.Width; x++ {
				r.screen.SetContent(editRegion.X+x, editRegion.Y+viewLine, ' ', nil, defaultStyle)
			}
			continue
		}

		// Get line content
		lineText, err := buf.GetLine(bufferLine)
		if err != nil {
			continue
		}

		// Determine style (highlight current line)
		lineStyle := defaultStyle
		if bufferLine == cursorPos.Line {
			lineStyle = currentLineStyle
		}

		// Render line content
		x := editRegion.X
		for i, char := range lineText {
			if i >= editRegion.Width {
				break
			}
			r.screen.SetContent(x+i, editRegion.Y+viewLine, char, nil, lineStyle)
		}

		// Fill remaining space in line
		for x := len(lineText); x < editRegion.Width; x++ {
			r.screen.SetContent(editRegion.X+x, editRegion.Y+viewLine, ' ', nil, lineStyle)
		}
	}

	return nil
}

// formatLineNumber formats a line number with right alignment.
func formatLineNumber(lineNum, width int) string {
	numStr := strconv.Itoa(lineNum)
	// Right-align by padding with spaces
	for len(numStr) < width {
		numStr = " " + numStr
	}
	return numStr
}
