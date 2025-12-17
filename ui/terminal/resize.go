// Package terminal implements resize event handling.
package terminal

import (
	"github.com/gdamore/tcell/v2"
)

// ResizeHandler handles terminal resize events.
type ResizeHandler struct {
	onResize func(width, height int)
}

// NewResizeHandler creates a new resize handler.
func NewResizeHandler(onResize func(width, height int)) *ResizeHandler {
	return &ResizeHandler{
		onResize: onResize,
	}
}

// HandleEvent processes an event and calls the resize callback if it's a resize event.
// Returns true if the event was a resize event and was handled.
func (h *ResizeHandler) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventResize:
		width, height := ev.Size()
		if h.onResize != nil {
			h.onResize(width, height)
		}
		return true
	}
	return false
}

// Dimensions represents screen dimensions.
type Dimensions struct {
	Width  int
	Height int
}

// GetDimensions extracts dimensions from a resize event.
func GetDimensions(ev tcell.Event) (Dimensions, bool) {
	resizeEv, ok := ev.(*tcell.EventResize)
	if !ok {
		return Dimensions{}, false
	}

	width, height := resizeEv.Size()
	return Dimensions{Width: width, Height: height}, true
}

