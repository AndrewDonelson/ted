package renderer

import (
	"strings"
	"testing"

	"github.com/AndrewDonelson/ted/ui/layout"
)

func TestRenderInfoBar(t *testing.T) {
	mockScr := newMockScreen(80, 24)
	layout := layout.NewLayout(80, 24)
	renderer := NewRenderer(mockScr, layout)

	fileInfo := &FileInfo{
		Name:       "test.txt",
		Path:       "/path/to/test.txt",
		Size:       1024,
		Type:       "Plain Text",
		Encoding:   "UTF-8",
		LineEnding: "LF",
		TabSize:    4,
		TotalLines: 10,
		IsModified: false,
	}

	if err := renderer.RenderInfoBar(fileInfo); err != nil {
		t.Errorf("RenderInfoBar() error = %v", err)
	}

	// Verify info bar was rendered
	region := layout.GetInfoBarRegion()
	if row, ok := mockScr.contents[region.Y]; ok {
		if len(row) == 0 {
			t.Error("RenderInfoBar() did not render content")
		}
	} else {
		t.Error("RenderInfoBar() did not set content at info bar region")
	}
}

func TestRenderInfoBar_NilFileInfo(t *testing.T) {
	mockScr := newMockScreen(80, 24)
	layout := layout.NewLayout(80, 24)
	renderer := NewRenderer(mockScr, layout)

	if err := renderer.RenderInfoBar(nil); err != nil {
		t.Errorf("RenderInfoBar() error = %v", err)
	}

	// Should render "[No File]"
	region := layout.GetInfoBarRegion()
	if row, ok := mockScr.contents[region.Y]; ok {
		found := false
		for _, char := range "[No File]" {
			for _, c := range row {
				if c == char {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			t.Error("RenderInfoBar(nil) did not render [No File]")
		}
	}
}

func TestRenderInfoBar_EmptyFileName(t *testing.T) {
	mockScr := newMockScreen(80, 24)
	layout := layout.NewLayout(80, 24)
	renderer := NewRenderer(mockScr, layout)

	fileInfo := &FileInfo{
		Name: "",
		Size: 0,
	}

	if err := renderer.RenderInfoBar(fileInfo); err != nil {
		t.Errorf("RenderInfoBar() error = %v", err)
	}
}

func TestRenderInfoBar_Modified(t *testing.T) {
	mockScr := newMockScreen(80, 24)
	layout := layout.NewLayout(80, 24)
	renderer := NewRenderer(mockScr, layout)

	fileInfo := &FileInfo{
		Name:       "test.txt",
		IsModified: true,
	}

	if err := renderer.RenderInfoBar(fileInfo); err != nil {
		t.Errorf("RenderInfoBar() error = %v", err)
	}

	// Verify "Modified" appears
	region := layout.GetInfoBarRegion()
	if row, ok := mockScr.contents[region.Y]; ok {
		found := false
		for _, char := range "Modified" {
			for _, c := range row {
				if c == char {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			t.Error("RenderInfoBar() did not show Modified status")
		}
	}
}

func TestRenderInfoBar_Saved(t *testing.T) {
	mockScr := newMockScreen(80, 24)
	layout := layout.NewLayout(80, 24)
	renderer := NewRenderer(mockScr, layout)

	fileInfo := &FileInfo{
		Name:       "test.txt",
		IsModified: false,
	}

	if err := renderer.RenderInfoBar(fileInfo); err != nil {
		t.Errorf("RenderInfoBar() error = %v", err)
	}
}

func TestRenderInfoBar_LongContent(t *testing.T) {
	mockScr := newMockScreen(40, 10)
	layout := layout.NewLayout(40, 10)
	renderer := NewRenderer(mockScr, layout)

	fileInfo := &FileInfo{
		Name:       "very-long-filename-that-exceeds-terminal-width.txt",
		Path:       "/very/long/path/to/file.txt",
		Size:       999999,
		Type:       "Very Long File Type Name",
		Encoding:   "UTF-8",
		LineEnding: "CRLF",
		TabSize:    8,
		TotalLines: 99999,
		IsModified: true,
	}

	if err := renderer.RenderInfoBar(fileInfo); err != nil {
		t.Errorf("RenderInfoBar() error = %v", err)
	}

	// Should truncate, not crash
	region := layout.GetInfoBarRegion()
	if row, ok := mockScr.contents[region.Y]; ok {
		if len(row) > region.Width {
			t.Errorf("Info bar content not truncated: %d > %d", len(row), region.Width)
		}
	}
}

func TestBuildInfoBarContent(t *testing.T) {
	mockScr := newMockScreen(80, 24)
	layout := layout.NewLayout(80, 24)
	renderer := NewRenderer(mockScr, layout)

	tests := []struct {
		name         string
		fileInfo     *FileInfo
		width        int
		wantContains []string
	}{
		{
			name: "complete file info",
			fileInfo: &FileInfo{
				Name:       "test.txt",
				Size:       1024,
				Type:       "Plain Text",
				LineEnding: "LF",
				TabSize:    4,
				IsModified: false,
			},
			width:        80,
			wantContains: []string{"test.txt", "1.0 KB", "Plain Text", "Saved", "Tab: 4", "LF"},
		},
		{
			name: "modified file",
			fileInfo: &FileInfo{
				Name:       "test.txt",
				IsModified: true,
			},
			width:        80,
			wantContains: []string{"test.txt", "Modified"},
		},
		{
			name:         "nil file info",
			fileInfo:     nil,
			width:        80,
			wantContains: []string{"[No File]"},
		},
		{
			name: "empty filename",
			fileInfo: &FileInfo{
				Name: "",
			},
			width:        80,
			wantContains: []string{"[No Name]"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := renderer.buildInfoBarContent(tt.fileInfo, tt.width)

			if len(content) > tt.width {
				t.Errorf("buildInfoBarContent() length = %d, want <= %d", len(content), tt.width)
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(content, want) {
					t.Errorf("buildInfoBarContent() = %q, want to contain %q", content, want)
				}
			}
		})
	}
}

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		name string
		size int64
		want string
	}{
		{
			name: "bytes",
			size: 500,
			want: "500 bytes",
		},
		{
			name: "kilobytes",
			size: 1024,
			want: "1.0 KB",
		},
		{
			name: "megabytes",
			size: 1024 * 1024,
			want: "1.0 MB",
		},
		{
			name: "gigabytes",
			size: 1024 * 1024 * 1024,
			want: "1.0 GB",
		},
		{
			name: "fractional KB",
			size: 1536,
			want: "1.5 KB",
		},
		{
			name: "zero",
			size: 0,
			want: "0 bytes",
		},
		{
			name: "large KB",
			size: 1024 * 500,
			want: "500.0 KB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatFileSize(tt.size)
			if got != tt.want {
				t.Errorf("formatFileSize(%d) = %q, want %q", tt.size, got, tt.want)
			}
		})
	}
}

func TestRenderInfoBarWithContent(t *testing.T) {
	mockScr := newMockScreen(80, 24)
	layout := layout.NewLayout(80, 24)
	renderer := NewRenderer(mockScr, layout)

	content := "Custom info bar content"
	if err := renderer.RenderInfoBarWithContent(content); err != nil {
		t.Errorf("RenderInfoBarWithContent() error = %v", err)
	}

	// Verify content was rendered
	region := layout.GetInfoBarRegion()
	if row, ok := mockScr.contents[region.Y]; ok {
		found := false
		for _, char := range content {
			for _, c := range row {
				if c == char {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			t.Error("RenderInfoBarWithContent() did not render content")
		}
	}
}

func TestRenderInfoBarWithContent_LongContent(t *testing.T) {
	mockScr := newMockScreen(40, 10)
	layout := layout.NewLayout(40, 10)
	renderer := NewRenderer(mockScr, layout)

	content := strings.Repeat("A", 100)
	if err := renderer.RenderInfoBarWithContent(content); err != nil {
		t.Errorf("RenderInfoBarWithContent() error = %v", err)
	}

	// Should truncate
	region := layout.GetInfoBarRegion()
	if row, ok := mockScr.contents[region.Y]; ok {
		if len(row) > region.Width {
			t.Errorf("Content not truncated: %d > %d", len(row), region.Width)
		}
	}
}

func TestRenderInfoBar_FillEmptySpace(t *testing.T) {
	mockScr := newMockScreen(80, 24)
	layout := layout.NewLayout(80, 24)
	renderer := NewRenderer(mockScr, layout)

	fileInfo := &FileInfo{
		Name: "short",
	}

	if err := renderer.RenderInfoBar(fileInfo); err != nil {
		t.Errorf("RenderInfoBar() error = %v", err)
	}

	// Verify empty space is filled
	region := layout.GetInfoBarRegion()
	if row, ok := mockScr.contents[region.Y]; ok {
		// Should fill entire width (content may be shorter, but spaces should fill the rest)
		// Count non-space characters to verify content was rendered
		nonSpaceCount := 0
		for _, char := range row {
			if char != ' ' {
				nonSpaceCount++
			}
		}
		if nonSpaceCount == 0 {
			t.Error("Info bar content not rendered")
		}
		// The row should have content up to region.Width (either content or spaces)
		maxX := 0
		for x := range row {
			if x > maxX {
				maxX = x
			}
		}
		// Should have content up to at least region.Width-1 (0-indexed)
		if maxX < region.Width-1 {
			t.Errorf("Info bar not fully filled: maxX=%d, want >=%d", maxX, region.Width-1)
		}
	}
}
