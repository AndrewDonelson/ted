// Package terminal implements keyboard input handling.
package terminal

import (
	"github.com/gdamore/tcell/v2"
)

// KeyAction represents an action triggered by a key press.
type KeyAction int

const (
	// KeyActionNone represents no action.
	KeyActionNone KeyAction = iota
	// KeyActionCharacter represents a character input.
	KeyActionCharacter
	// KeyActionMoveLeft represents cursor move left.
	KeyActionMoveLeft
	// KeyActionMoveRight represents cursor move right.
	KeyActionMoveRight
	// KeyActionMoveUp represents cursor move up.
	KeyActionMoveUp
	// KeyActionMoveDown represents cursor move down.
	KeyActionMoveDown
	// KeyActionBackspace represents backspace key.
	KeyActionBackspace
	// KeyActionDelete represents delete key.
	KeyActionDelete
	// KeyActionSave represents Ctrl+S (save).
	KeyActionSave
	// KeyActionQuit represents Ctrl+Q (quit).
	KeyActionQuit
	// KeyActionEnter represents Enter/Return key.
	KeyActionEnter
	// KeyActionHome represents Home key.
	KeyActionHome
	// KeyActionEnd represents End key.
	KeyActionEnd
	// KeyActionUndo represents Ctrl+Z (undo).
	KeyActionUndo
	// KeyActionRedo represents Ctrl+Y (redo).
	KeyActionRedo
	// KeyActionCut represents Ctrl+X (cut).
	KeyActionCut
	// KeyActionCopy represents Ctrl+C (copy).
	KeyActionCopy
	// KeyActionPaste represents Ctrl+V (paste).
	KeyActionPaste
	// KeyActionSelectLeft represents Shift+Left (extend selection left).
	KeyActionSelectLeft
	// KeyActionSelectRight represents Shift+Right (extend selection right).
	KeyActionSelectRight
	// KeyActionSelectUp represents Shift+Up (extend selection up).
	KeyActionSelectUp
	// KeyActionSelectDown represents Shift+Down (extend selection down).
	KeyActionSelectDown
)

// KeyEvent represents a processed keyboard event.
type KeyEvent struct {
	Action    KeyAction
	Character rune
	Key       tcell.Key
	Modifiers tcell.ModMask
}

// ProcessEvent processes a tcell event and converts it to a KeyEvent.
// Returns nil if the event is not a keyboard event or is not handled.
func ProcessEvent(ev tcell.Event) *KeyEvent {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		return processKeyEvent(ev)
	default:
		return nil
	}
}

// processKeyEvent processes a tcell key event.
func processKeyEvent(ev *tcell.EventKey) *KeyEvent {
	key := ev.Key()
	modifiers := ev.Modifiers()
	r := ev.Rune()

	// Handle special keys
	switch key {
	case tcell.KeyCtrlS:
		return &KeyEvent{Action: KeyActionSave, Key: key, Modifiers: modifiers}
	case tcell.KeyCtrlQ:
		return &KeyEvent{Action: KeyActionQuit, Key: key, Modifiers: modifiers}
	case tcell.KeyCtrlZ:
		return &KeyEvent{Action: KeyActionUndo, Key: key, Modifiers: modifiers}
	case tcell.KeyCtrlY:
		return &KeyEvent{Action: KeyActionRedo, Key: key, Modifiers: modifiers}
	case tcell.KeyCtrlX:
		return &KeyEvent{Action: KeyActionCut, Key: key, Modifiers: modifiers}
	case tcell.KeyCtrlC:
		return &KeyEvent{Action: KeyActionCopy, Key: key, Modifiers: modifiers}
	case tcell.KeyCtrlV:
		return &KeyEvent{Action: KeyActionPaste, Key: key, Modifiers: modifiers}
	case tcell.KeyLeft:
		return &KeyEvent{Action: KeyActionMoveLeft, Key: key, Modifiers: modifiers}
	case tcell.KeyRight:
		return &KeyEvent{Action: KeyActionMoveRight, Key: key, Modifiers: modifiers}
	case tcell.KeyUp:
		return &KeyEvent{Action: KeyActionMoveUp, Key: key, Modifiers: modifiers}
	case tcell.KeyDown:
		return &KeyEvent{Action: KeyActionMoveDown, Key: key, Modifiers: modifiers}
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		return &KeyEvent{Action: KeyActionBackspace, Key: key, Modifiers: modifiers}
	case tcell.KeyDelete:
		return &KeyEvent{Action: KeyActionDelete, Key: key, Modifiers: modifiers}
	case tcell.KeyEnter:
		return &KeyEvent{Action: KeyActionEnter, Key: key, Modifiers: modifiers, Character: '\n'}
	case tcell.KeyHome:
		return &KeyEvent{Action: KeyActionHome, Key: key, Modifiers: modifiers}
	case tcell.KeyEnd:
		return &KeyEvent{Action: KeyActionEnd, Key: key, Modifiers: modifiers}
	case tcell.KeyRune:
		// Regular character input
		// Ignore Alt+key combinations (they're used for menus in Phase 1+)
		if modifiers&tcell.ModAlt != 0 {
			return &KeyEvent{Action: KeyActionNone, Key: key, Modifiers: modifiers}
		}
		if r != 0 {
			return &KeyEvent{Action: KeyActionCharacter, Character: r, Key: key, Modifiers: modifiers}
		}
	}

	return &KeyEvent{Action: KeyActionNone, Key: key, Modifiers: modifiers}
}

// IsPrintable returns true if the key event represents a printable character.
func (ke *KeyEvent) IsPrintable() bool {
	return ke.Action == KeyActionCharacter && ke.Character != 0
}

// IsModifier returns true if a modifier key is pressed.
func (ke *KeyEvent) HasModifier(mod tcell.ModMask) bool {
	return ke.Modifiers&mod != 0
}
