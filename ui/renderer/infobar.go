// Package renderer implements info bar rendering.
//
// CRITICAL: The info bar MUST use inverted colors (light background, dark text)
// to visually distinguish it from the editing area.
package renderer

import (
	"fmt"
	"strings"
)

// FileInfo contains information to display in the info bar.
type FileInfo struct {
	Name       string
	Path       string
	Size       int64
	Type       string
	Encoding   string
	LineEnding string
	TabSize    int
	TotalLines int
	IsModified bool
}

// RenderInfoBar renders the info bar at the bottom of the screen.
// CRITICAL: Uses INVERTED colors (light bg #d4d4d4, dark text #1e1e1e).
func (r *Renderer) RenderInfoBar(info *FileInfo) error {
	region := r.layout.GetInfoBarRegion()
	style := GetInfoBarStyle() // INVERTED style

	// Fill entire info bar region with inverted background color first
	for x := 0; x < region.Width; x++ {
		r.screen.SetContent(region.X+x, region.Y, ' ', nil, style)
	}

	// Build info bar content
	content := r.buildInfoBarContent(info, region.Width)

	// Render content
	for i, char := range content {
		if i >= region.Width {
			break
		}
		r.screen.SetContent(region.X+i, region.Y, char, nil, style)
	}

	return nil
}

// buildInfoBarContent builds the info bar text content.
func (r *Renderer) buildInfoBarContent(info *FileInfo, width int) string {
	if info == nil {
		return "[No File]"
	}

	var parts []string

	// Filename
	filename := info.Name
	if filename == "" {
		filename = "[No Name]"
	}
	parts = append(parts, filename)

	// File size
	if info.Size > 0 {
		sizeStr := formatFileSize(info.Size)
		parts = append(parts, sizeStr)
	}

	// File type
	if info.Type != "" {
		parts = append(parts, info.Type)
	}

	// Modified status
	if info.IsModified {
		parts = append(parts, "Modified")
	} else {
		parts = append(parts, "Saved")
	}

	// Tab size
	if info.TabSize > 0 {
		parts = append(parts, fmt.Sprintf("Tab: %d", info.TabSize))
	}

	// Line ending
	if info.LineEnding != "" {
		parts = append(parts, info.LineEnding)
	}

	// Join with separators
	separator := " │ "
	content := strings.Join(parts, separator)

	// Truncate if too long
	if len(content) > width {
		content = content[:width-3] + "..."
	}

	return content
}

// formatFileSize formats file size in human-readable format.
func formatFileSize(size int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case size < KB:
		return fmt.Sprintf("%d bytes", size)
	case size < MB:
		return fmt.Sprintf("%.1f KB", float64(size)/KB)
	case size < GB:
		return fmt.Sprintf("%.1f MB", float64(size)/MB)
	default:
		return fmt.Sprintf("%.1f GB", float64(size)/GB)
	}
}

// RenderInfoBarWithContent renders the info bar with custom content.
func (r *Renderer) RenderInfoBarWithContent(content string) error {
	region := r.layout.GetInfoBarRegion()
	style := GetInfoBarStyle() // INVERTED style

	// Truncate if too long
	if len(content) > region.Width {
		content = content[:region.Width-3] + "..."
	}

	// Render content
	for i, char := range content {
		if i >= region.Width {
			break
		}
		r.screen.SetContent(region.X+i, region.Y, char, nil, style)
	}

	// Fill remaining space
	for i := len(content); i < region.Width; i++ {
		r.screen.SetContent(region.X+i, region.Y, ' ', nil, style)
	}

	return nil
}
