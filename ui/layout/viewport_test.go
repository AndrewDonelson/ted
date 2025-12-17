package layout

import (
	"reflect"
	"testing"
)

func TestNewLayout(t *testing.T) {
	l := NewLayout(80, 24)

	if l.width != 80 {
		t.Errorf("NewLayout() width = %d, want 80", l.width)
	}
	if l.height != 24 {
		t.Errorf("NewLayout() height = %d, want 24", l.height)
	}
}

func TestLayout_GetMenuBarRegion(t *testing.T) {
	l := NewLayout(80, 24)
	region := l.GetMenuBarRegion()

	want := Region{X: 0, Y: 0, Width: 80, Height: 1}
	if !reflect.DeepEqual(region, want) {
		t.Errorf("GetMenuBarRegion() = %v, want %v", region, want)
	}
}

func TestLayout_GetEditAreaRegion(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
		want   Region
	}{
		{
			name:   "normal size",
			width:  80,
			height: 24,
			want:   Region{X: 0, Y: 1, Width: 80, Height: 22},
		},
		{
			name:   "minimum size",
			width:  40,
			height: 10,
			want:   Region{X: 0, Y: 1, Width: 40, Height: 8},
		},
		{
			name:   "very small",
			width:  40,
			height: 3,
			want:   Region{X: 0, Y: 1, Width: 40, Height: 1}, // Minimum 1 line
		},
		{
			name:   "wide terminal",
			width:  200,
			height: 50,
			want:   Region{X: 0, Y: 1, Width: 200, Height: 48},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLayout(tt.width, tt.height)
			got := l.GetEditAreaRegion()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEditAreaRegion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLayout_GetInfoBarRegion(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
		want   Region
	}{
		{
			name:   "normal size",
			width:  80,
			height: 24,
			want:   Region{X: 0, Y: 23, Width: 80, Height: 1},
		},
		{
			name:   "minimum size",
			width:  40,
			height: 10,
			want:   Region{X: 0, Y: 9, Width: 40, Height: 1},
		},
		{
			name:   "very small",
			width:  40,
			height: 2,
			want:   Region{X: 0, Y: 1, Width: 40, Height: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLayout(tt.width, tt.height)
			got := l.GetInfoBarRegion()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetInfoBarRegion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLayout_GetLineNumberWidth(t *testing.T) {
	tests := []struct {
		name       string
		totalLines int
		want       int
	}{
		{
			name:       "single digit",
			totalLines: 5,
			want:       3, // 1 digit + 2 (space + separator)
		},
		{
			name:       "two digits",
			totalLines: 50,
			want:       4, // 2 digits + 2
		},
		{
			name:       "three digits",
			totalLines: 500,
			want:       5, // 3 digits + 2
		},
		{
			name:       "four digits",
			totalLines: 5000,
			want:       6, // 4 digits + 2
		},
		{
			name:       "zero lines",
			totalLines: 0,
			want:       0,
		},
		{
			name:       "one line",
			totalLines: 1,
			want:       3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLayout(80, 24)
			got := l.GetLineNumberWidth(tt.totalLines)

			if got != tt.want {
				t.Errorf("GetLineNumberWidth(%d) = %d, want %d", tt.totalLines, got, tt.want)
			}
		})
	}
}

func TestLayout_CalculateViewport(t *testing.T) {
	tests := []struct {
		name       string
		width      int
		height     int
		cursorLine int
		totalLines int
		want       Viewport
	}{
		{
			name:       "cursor at start",
			width:      80,
			height:     24,
			cursorLine: 0,
			totalLines: 100,
			want: Viewport{
				StartLine: 0,
				EndLine:   21, // 22 lines visible (24 - 1 menu - 1 info)
				OffsetX:   0,
				Width:     80,
				Height:    22,
			},
		},
		{
			name:       "cursor in middle",
			width:      80,
			height:     24,
			cursorLine: 50,
			totalLines: 100,
			want: Viewport{
				StartLine: 39, // 50 - 11 (half of 22)
				EndLine:   60, // 39 + 22 - 1
				OffsetX:   0,
				Width:     80,
				Height:    22,
			},
		},
		{
			name:       "cursor at end",
			width:      80,
			height:     24,
			cursorLine: 99,
			totalLines: 100,
			want: Viewport{
				StartLine: 78, // 99 - 21 (to show last line)
				EndLine:   99,
				OffsetX:   0,
				Width:     80,
				Height:    22,
			},
		},
		{
			name:       "fewer lines than viewport",
			width:      80,
			height:     24,
			cursorLine: 5,
			totalLines: 10,
			want: Viewport{
				StartLine: 0,
				EndLine:   9,
				OffsetX:   0,
				Width:     80,
				Height:    22,
			},
		},
		{
			name:       "empty buffer",
			width:      80,
			height:     24,
			cursorLine: 0,
			totalLines: 0,
			want: Viewport{
				StartLine: 0,
				EndLine:   0,
				OffsetX:   0,
				Width:     80,
				Height:    22,
			},
		},
		{
			name:       "cursor beyond end",
			width:      80,
			height:     24,
			cursorLine: 150,
			totalLines: 100,
			want: Viewport{
				StartLine: 78, // Clamped to last viewport
				EndLine:   99,
				OffsetX:   0,
				Width:     80,
				Height:    22,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLayout(tt.width, tt.height)
			got := l.CalculateViewport(tt.cursorLine, tt.totalLines)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculateViewport() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestLayout_AdjustForResize(t *testing.T) {
	l := NewLayout(80, 24)

	l.AdjustForResize(120, 50)

	if l.width != 120 {
		t.Errorf("AdjustForResize() width = %d, want 120", l.width)
	}
	if l.height != 50 {
		t.Errorf("AdjustForResize() height = %d, want 50", l.height)
	}

	// Verify regions are updated
	editRegion := l.GetEditAreaRegion()
	if editRegion.Width != 120 {
		t.Errorf("GetEditAreaRegion() width = %d, want 120", editRegion.Width)
	}
	if editRegion.Height != 48 {
		t.Errorf("GetEditAreaRegion() height = %d, want 48", editRegion.Height)
	}
}

func TestLayout_ScreenToBuffer(t *testing.T) {
	l := NewLayout(80, 24)

	tests := []struct {
		name     string
		screenX  int
		screenY  int
		wantLine int
		wantCol  int
	}{
		{
			name:     "top of edit area",
			screenX:  10,
			screenY:  1, // First line of edit area
			wantLine: 0,
			wantCol:  10,
		},
		{
			name:     "middle of edit area",
			screenX:  20,
			screenY:  12, // Middle of edit area
			wantLine: 11,
			wantCol:  20,
		},
		{
			name:     "in menu bar",
			screenX:  10,
			screenY:  0,
			wantLine: -1,
			wantCol:  -1,
		},
		{
			name:     "in info bar",
			screenX:  10,
			screenY:  23,
			wantLine: -1,
			wantCol:  -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			line, col := l.ScreenToBuffer(tt.screenX, tt.screenY)

			if line != tt.wantLine || col != tt.wantCol {
				t.Errorf("ScreenToBuffer(%d, %d) = (%d, %d), want (%d, %d)",
					tt.screenX, tt.screenY, line, col, tt.wantLine, tt.wantCol)
			}
		})
	}
}

func TestLayout_BufferToScreen(t *testing.T) {
	l := NewLayout(80, 24)
	viewport := l.CalculateViewport(10, 100) // Cursor at line 10, 100 total lines

	tests := []struct {
		name       string
		bufferLine int
		bufferCol  int
		wantX      int
		wantY      int
	}{
		{
			name:       "first visible line",
			bufferLine: viewport.StartLine,
			bufferCol:  5,
			wantX:      5,
			wantY:      1, // Edit area starts at Y=1
		},
		{
			name:       "last visible line",
			bufferLine: viewport.EndLine,
			bufferCol:  10,
			wantX:      10,
			wantY:      1 + (viewport.EndLine - viewport.StartLine),
		},
		{
			name:       "line before viewport",
			bufferLine: viewport.StartLine - 1,
			bufferCol:  5,
			wantX:      -1,
			wantY:      -1,
		},
		{
			name:       "line after viewport",
			bufferLine: viewport.EndLine + 1,
			bufferCol:  5,
			wantX:      -1,
			wantY:      -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, y := l.BufferToScreen(tt.bufferLine, tt.bufferCol, viewport)

			if x != tt.wantX || y != tt.wantY {
				t.Errorf("BufferToScreen(%d, %d) = (%d, %d), want (%d, %d)",
					tt.bufferLine, tt.bufferCol, x, y, tt.wantX, tt.wantY)
			}
		})
	}
}

func TestLayout_IsSizeValid(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
		want   bool
	}{
		{
			name:   "valid size",
			width:  80,
			height: 24,
			want:   true,
		},
		{
			name:   "minimum size",
			width:  40,
			height: 10,
			want:   true,
		},
		{
			name:   "too narrow",
			width:  30,
			height: 10,
			want:   false,
		},
		{
			name:   "too short",
			width:  40,
			height: 5,
			want:   false,
		},
		{
			name:   "both too small",
			width:  30,
			height: 5,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLayout(tt.width, tt.height)
			got := l.IsSizeValid()

			if got != tt.want {
				t.Errorf("IsSizeValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMinimumSize(t *testing.T) {
	width, height := GetMinimumSize()

	if width != 40 {
		t.Errorf("GetMinimumSize() width = %d, want 40", width)
	}
	if height != 10 {
		t.Errorf("GetMinimumSize() height = %d, want 10", height)
	}
}
