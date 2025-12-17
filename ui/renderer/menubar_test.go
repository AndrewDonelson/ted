package renderer

import (
	"strings"
	"testing"

	"github.com/AndrewDonelson/ted/ui/layout"
)

func TestRenderMenuBar(t *testing.T) {
	mockScr := newMockScreen(80, 24)
	layout := layout.NewLayout(80, 24)
	renderer := NewRenderer(mockScr, layout)

	if err := renderer.RenderMenuBar(); err != nil {
		t.Errorf("RenderMenuBar() error = %v", err)
	}

	// Verify menu items were rendered
	region := layout.GetMenuBarRegion()
	foundMenu := false
	for y := region.Y; y < region.Y+region.Height; y++ {
		if row, ok := mockScr.contents[y]; ok {
			for x := range row {
				if x < len("File  Edit  Search  View  Help") {
					foundMenu = true
					break
				}
			}
		}
		if foundMenu {
			break
		}
	}

	if !foundMenu {
		t.Error("RenderMenuBar() did not render menu items")
	}
}

func TestRenderMenuBarWithStatus(t *testing.T) {
	mockScr := newMockScreen(80, 24)
	layout := layout.NewLayout(80, 24)
	renderer := NewRenderer(mockScr, layout)

	if err := renderer.RenderMenuBarWithStatus("INS", "UTF-8", 5, 10); err != nil {
		t.Errorf("RenderMenuBarWithStatus() error = %v", err)
	}
}

func TestRenderMenuBar_NarrowTerminal(t *testing.T) {
	mockScr := newMockScreen(40, 10)
	layout := layout.NewLayout(40, 10)
	renderer := NewRenderer(mockScr, layout)

	if err := renderer.RenderMenuBar(); err != nil {
		t.Errorf("RenderMenuBar() error = %v", err)
	}
}

func TestRenderMenuBar_VeryNarrowTerminal(t *testing.T) {
	mockScr := newMockScreen(20, 10)
	layout := layout.NewLayout(20, 10)
	renderer := NewRenderer(mockScr, layout)

	if err := renderer.RenderMenuBar(); err != nil {
		t.Errorf("RenderMenuBar() error = %v", err)
	}
}

func TestFormatStatus(t *testing.T) {
	tests := []struct {
		name     string
		mode     string
		encoding string
		line     int
		col      int
		want     string
	}{
		{
			name:     "basic status",
			mode:     "INS",
			encoding: "UTF-8",
			line:     5,
			col:      10,
			want:     "INS │ UTF-8 │ LN 6, COL 11",
		},
		{
			name:     "overwrite mode",
			mode:     "OVR",
			encoding: "UTF-8",
			line:     0,
			col:      0,
			want:     "OVR │ UTF-8 │ LN 1, COL 1",
		},
		{
			name:     "large line numbers",
			mode:     "INS",
			encoding: "UTF-8",
			line:     999,
			col:      999,
			want:     "INS │ UTF-8 │ LN 1000, COL 1000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatStatus(tt.mode, tt.encoding, tt.line, tt.col)
			if got != tt.want {
				t.Errorf("formatStatus() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want string
	}{
		{
			name: "zero",
			n:    0,
			want: "0",
		},
		{
			name: "single digit",
			n:    5,
			want: "5",
		},
		{
			name: "two digits",
			n:    42,
			want: "42",
		},
		{
			name: "three digits",
			n:    123,
			want: "123",
		},
		{
			name: "large number",
			n:    99999,
			want: "99999",
		},
		{
			name: "negative",
			n:    -5,
			want: "0", // Should clamp to 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatNumber(tt.n)
			if got != tt.want {
				t.Errorf("formatNumber(%d) = %q, want %q", tt.n, got, tt.want)
			}
		})
	}
}

func TestSetMenuBarContent(t *testing.T) {
	mockScr := newMockScreen(80, 24)
	layout := layout.NewLayout(80, 24)
	renderer := NewRenderer(mockScr, layout)

	content := "Test Menu Content"
	style := GetMenuBarStyle()

	if err := renderer.SetMenuBarContent(content, style); err != nil {
		t.Errorf("SetMenuBarContent() error = %v", err)
	}

	// Verify content was set
	region := layout.GetMenuBarRegion()
	found := false
	for i, char := range content {
		if i >= region.Width {
			break
		}
		if row, ok := mockScr.contents[region.Y]; ok {
			if row[i] == char {
				found = true
				break
			}
		}
	}

	if !found {
		t.Error("SetMenuBarContent() did not set content")
	}
}

func TestSetMenuBarContent_LongContent(t *testing.T) {
	mockScr := newMockScreen(20, 10)
	layout := layout.NewLayout(20, 10)
	renderer := NewRenderer(mockScr, layout)

	content := strings.Repeat("A", 100) // Longer than screen width
	style := GetMenuBarStyle()

	if err := renderer.SetMenuBarContent(content, style); err != nil {
		t.Errorf("SetMenuBarContent() error = %v", err)
	}

	// Should truncate, not crash
	region := layout.GetMenuBarRegion()
	if len(mockScr.contents[region.Y]) > region.Width {
		t.Errorf("Content not truncated: %d > %d", len(mockScr.contents[region.Y]), region.Width)
	}
}
