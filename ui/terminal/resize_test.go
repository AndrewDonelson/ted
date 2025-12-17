package terminal

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestNewResizeHandler(t *testing.T) {
	called := false
	var capturedWidth, capturedHeight int

	onResize := func(width, height int) {
		called = true
		capturedWidth = width
		capturedHeight = height
	}

	handler := NewResizeHandler(onResize)
	if handler == nil {
		t.Fatal("NewResizeHandler() returned nil")
	}

	// Create a mock resize event
	ev := tcell.NewEventResize(80, 24)

	// Handle the event
	result := handler.HandleEvent(ev)

	if !result {
		t.Error("HandleEvent() returned false for resize event")
	}

	if !called {
		t.Error("onResize callback was not called")
	}

	if capturedWidth != 80 || capturedHeight != 24 {
		t.Errorf("onResize called with (%d, %d), want (80, 24)", capturedWidth, capturedHeight)
	}
}

func TestResizeHandler_HandleEvent_NonResize(t *testing.T) {
	called := false
	onResize := func(width, height int) {
		called = true
	}

	handler := NewResizeHandler(onResize)

	// Create a non-resize event (key event)
	ev := tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone)

	result := handler.HandleEvent(ev)

	if result {
		t.Error("HandleEvent() returned true for non-resize event")
	}

	if called {
		t.Error("onResize callback was called for non-resize event")
	}
}

func TestResizeHandler_HandleEvent_NilCallback(t *testing.T) {
	handler := &ResizeHandler{
		onResize: nil,
	}

	ev := tcell.NewEventResize(80, 24)

	// Should not panic
	result := handler.HandleEvent(ev)

	if !result {
		t.Error("HandleEvent() should return true for resize event even with nil callback")
	}
}

func TestGetDimensions(t *testing.T) {
	tests := []struct {
		name   string
		ev     tcell.Event
		want   Dimensions
		wantOk bool
	}{
		{
			name:   "resize event",
			ev:     tcell.NewEventResize(80, 24),
			want:   Dimensions{Width: 80, Height: 24},
			wantOk: true,
		},
		{
			name:   "key event",
			ev:     tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone),
			want:   Dimensions{},
			wantOk: false,
		},
		{
			name:   "large dimensions",
			ev:     tcell.NewEventResize(200, 100),
			want:   Dimensions{Width: 200, Height: 100},
			wantOk: true,
		},
		{
			name:   "small dimensions",
			ev:     tcell.NewEventResize(40, 10),
			want:   Dimensions{Width: 40, Height: 10},
			wantOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := GetDimensions(tt.ev)

			if ok != tt.wantOk {
				t.Errorf("GetDimensions() ok = %v, want %v", ok, tt.wantOk)
			}

			if ok && got != tt.want {
				t.Errorf("GetDimensions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResizeHandler_MultipleCalls(t *testing.T) {
	callCount := 0
	var dimensions []Dimensions

	onResize := func(width, height int) {
		callCount++
		dimensions = append(dimensions, Dimensions{Width: width, Height: height})
	}

	handler := NewResizeHandler(onResize)

	// Handle multiple resize events
	ev1 := tcell.NewEventResize(80, 24)
	ev2 := tcell.NewEventResize(100, 30)
	ev3 := tcell.NewEventResize(120, 40)

	handler.HandleEvent(ev1)
	handler.HandleEvent(ev2)
	handler.HandleEvent(ev3)

	if callCount != 3 {
		t.Errorf("onResize called %d times, want 3", callCount)
	}

	if len(dimensions) != 3 {
		t.Errorf("captured %d dimensions, want 3", len(dimensions))
	}

	expected := []Dimensions{
		{Width: 80, Height: 24},
		{Width: 100, Height: 30},
		{Width: 120, Height: 40},
	}

	for i, dim := range dimensions {
		if dim != expected[i] {
			t.Errorf("dimensions[%d] = %v, want %v", i, dim, expected[i])
		}
	}
}
