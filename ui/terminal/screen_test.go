package terminal

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestNewScreen(t *testing.T) {
	// Note: This test may fail in non-terminal environments
	// We'll skip if terminal is not available
	screen, err := NewScreen()
	if err != nil {
		// Expected in non-terminal environments
		t.Skipf("Skipping test - terminal not available: %v", err)
		return
	}
	defer screen.Fini()

	if screen == nil {
		t.Fatal("NewScreen() returned nil screen")
	}

	// Test basic operations
	screen.Clear()
	width, height := screen.GetSize()
	if width <= 0 || height <= 0 {
		t.Errorf("GetSize() returned invalid dimensions: %dx%d", width, height)
	}
}

func TestTCellScreen_Clear(t *testing.T) {
	screen, err := NewScreen()
	if err != nil {
		t.Skipf("Skipping test - terminal not available: %v", err)
		return
	}
	defer screen.Fini()

	// Should not panic
	screen.Clear()
}

func TestTCellScreen_Refresh(t *testing.T) {
	screen, err := NewScreen()
	if err != nil {
		t.Skipf("Skipping test - terminal not available: %v", err)
		return
	}
	defer screen.Fini()

	if err := screen.Refresh(); err != nil {
		t.Errorf("Refresh() error = %v", err)
	}
}

func TestTCellScreen_GetSize(t *testing.T) {
	screen, err := NewScreen()
	if err != nil {
		t.Skipf("Skipping test - terminal not available: %v", err)
		return
	}
	defer screen.Fini()

	width, height := screen.GetSize()
	if width < 0 || height < 0 {
		t.Errorf("GetSize() returned negative dimensions: %dx%d", width, height)
	}
}

func TestTCellScreen_SetContent(t *testing.T) {
	screen, err := NewScreen()
	if err != nil {
		t.Skipf("Skipping test - terminal not available: %v", err)
		return
	}
	defer screen.Fini()

	style := tcell.StyleDefault
	tests := []struct {
		name    string
		x, y    int
		mainc   rune
		combc   []rune
		wantErr bool
	}{
		{
			name:    "set simple character",
			x:       0,
			y:       0,
			mainc:   'A',
			combc:   nil,
			wantErr: false,
		},
		{
			name:    "set character with combining",
			x:       5,
			y:       10,
			mainc:   'e',
			combc:   []rune{'\u0301'}, // Combining acute accent
			wantErr: false,
		},
		{
			name:    "set at negative coordinates",
			x:       -1,
			y:       -1,
			mainc:   'X',
			combc:   nil,
			wantErr: false, // tcell may handle this gracefully
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := screen.SetContent(tt.x, tt.y, tt.mainc, tt.combc, style)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTCellScreen_ShowCursor(t *testing.T) {
	screen, err := NewScreen()
	if err != nil {
		t.Skipf("Skipping test - terminal not available: %v", err)
		return
	}
	defer screen.Fini()

	// Should not panic
	screen.ShowCursor(0, 0)
	screen.ShowCursor(10, 5)
	screen.ShowCursor(-1, -1) // Edge case
}

func TestTCellScreen_HideCursor(t *testing.T) {
	screen, err := NewScreen()
	if err != nil {
		t.Skipf("Skipping test - terminal not available: %v", err)
		return
	}
	defer screen.Fini()

	// Should not panic
	screen.HideCursor()
}

func TestTCellScreen_PollEvent(t *testing.T) {
	screen, err := NewScreen()
	if err != nil {
		t.Skipf("Skipping test - terminal not available: %v", err)
		return
	}
	defer screen.Fini()

	// PollEvent blocks, so we can't test it easily in unit tests
	// This is more of an integration test concern
	// Just verify the method exists and doesn't panic immediately
	_ = screen.PollEvent
}

func TestTCellScreen_Fini(t *testing.T) {
	screen, err := NewScreen()
	if err != nil {
		t.Skipf("Skipping test - terminal not available: %v", err)
		return
	}

	// Should not panic
	screen.Fini()

	// Calling Fini multiple times should be safe
	screen.Fini()
}

func TestTCellScreen_GetRawScreen(t *testing.T) {
	screen, err := NewScreen()
	if err != nil {
		t.Skipf("Skipping test - terminal not available: %v", err)
		return
	}
	defer screen.Fini()

	raw := screen.GetRawScreen()
	if raw == nil {
		t.Error("GetRawScreen() returned nil")
	}
}
