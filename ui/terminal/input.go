// Package terminal implements keyboard input handling.
package terminal

import (
	"unicode"

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
	// KeyActionSelectAll represents Ctrl+A (select all).
	KeyActionSelectAll
	// KeyActionNew represents Ctrl+N (new file).
	KeyActionNew
	// KeyActionOpen represents Ctrl+O (open file).
	KeyActionOpen
	// KeyActionFind represents Ctrl+F (find).
	KeyActionFind
	// KeyActionReplace represents Ctrl+H (replace).
	KeyActionReplace
	// KeyActionGoToLine represents Ctrl+G (go to line).
	KeyActionGoToLine
	// KeyActionToggleLineNumbers represents Ctrl+L (toggle line numbers).
	KeyActionToggleLineNumbers
	// KeyActionHelp represents F1 (help).
	KeyActionHelp
	// KeyActionMenuToggle represents F10 (toggle menu).
	KeyActionMenuToggle
	// KeyActionMenuAlt represents Alt+key for menu shortcuts.
	KeyActionMenuAlt
	// KeyActionEscape represents Escape key.
	KeyActionEscape
	// Line Operations
	// KeyActionDeleteLine represents Ctrl+Shift+K (delete line).
	KeyActionDeleteLine
	// KeyActionDuplicateLine represents Ctrl+D (duplicate line).
	KeyActionDuplicateLine
	// KeyActionMoveLineUp represents Alt+Up (move line up).
	KeyActionMoveLineUp
	// KeyActionMoveLineDown represents Alt+Down (move line down).
	KeyActionMoveLineDown
	// KeyActionInsertLineAbove represents Ctrl+Shift+Enter (insert line above).
	KeyActionInsertLineAbove
	// KeyActionInsertLineBelow represents Ctrl+Enter (insert line below).
	KeyActionInsertLineBelow
	// Word Navigation
	// KeyActionWordLeft represents Ctrl+Left (move to previous word).
	KeyActionWordLeft
	// KeyActionWordRight represents Ctrl+Right (move to next word).
	KeyActionWordRight
	// Page Navigation
	// KeyActionPageUp represents Page Up key.
	KeyActionPageUp
	// KeyActionPageDown represents Page Down key.
	KeyActionPageDown
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

	// Handle special keys first
	switch key {
	case tcell.KeyEscape:
		return &KeyEvent{Action: KeyActionEscape, Key: key, Modifiers: modifiers}
	case tcell.KeyF1:
		return &KeyEvent{Action: KeyActionHelp, Key: key, Modifiers: modifiers}
	case tcell.KeyF10:
		return &KeyEvent{Action: KeyActionMenuToggle, Key: key, Modifiers: modifiers}
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
	case tcell.KeyCtrlA:
		return &KeyEvent{Action: KeyActionSelectAll, Key: key, Modifiers: modifiers}
	case tcell.KeyCtrlN:
		return &KeyEvent{Action: KeyActionNew, Key: key, Modifiers: modifiers}
	case tcell.KeyCtrlO:
		return &KeyEvent{Action: KeyActionOpen, Key: key, Modifiers: modifiers}
	case tcell.KeyCtrlF:
		return &KeyEvent{Action: KeyActionFind, Key: key, Modifiers: modifiers}
	case tcell.KeyCtrlH:
		return &KeyEvent{Action: KeyActionReplace, Key: key, Modifiers: modifiers}
	case tcell.KeyCtrlG:
		return &KeyEvent{Action: KeyActionGoToLine, Key: key, Modifiers: modifiers}
	case tcell.KeyCtrlL:
		return &KeyEvent{Action: KeyActionToggleLineNumbers, Key: key, Modifiers: modifiers}
	case tcell.KeyLeft:
		if modifiers&tcell.ModCtrl != 0 {
			return &KeyEvent{Action: KeyActionWordLeft, Key: key, Modifiers: modifiers}
		}
		if modifiers&tcell.ModShift != 0 {
			return &KeyEvent{Action: KeyActionSelectLeft, Key: key, Modifiers: modifiers}
		}
		return &KeyEvent{Action: KeyActionMoveLeft, Key: key, Modifiers: modifiers}
	case tcell.KeyRight:
		if modifiers&tcell.ModCtrl != 0 {
			return &KeyEvent{Action: KeyActionWordRight, Key: key, Modifiers: modifiers}
		}
		if modifiers&tcell.ModShift != 0 {
			return &KeyEvent{Action: KeyActionSelectRight, Key: key, Modifiers: modifiers}
		}
		return &KeyEvent{Action: KeyActionMoveRight, Key: key, Modifiers: modifiers}
	case tcell.KeyUp:
		if modifiers&tcell.ModAlt != 0 {
			return &KeyEvent{Action: KeyActionMoveLineUp, Key: key, Modifiers: modifiers}
		}
		if modifiers&tcell.ModShift != 0 {
			return &KeyEvent{Action: KeyActionSelectUp, Key: key, Modifiers: modifiers}
		}
		return &KeyEvent{Action: KeyActionMoveUp, Key: key, Modifiers: modifiers}
	case tcell.KeyDown:
		if modifiers&tcell.ModAlt != 0 {
			return &KeyEvent{Action: KeyActionMoveLineDown, Key: key, Modifiers: modifiers}
		}
		if modifiers&tcell.ModShift != 0 {
			return &KeyEvent{Action: KeyActionSelectDown, Key: key, Modifiers: modifiers}
		}
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
	case tcell.KeyPgUp:
		return &KeyEvent{Action: KeyActionPageUp, Key: key, Modifiers: modifiers}
	case tcell.KeyPgDn:
		return &KeyEvent{Action: KeyActionPageDown, Key: key, Modifiers: modifiers}
	case tcell.KeyCtrlK:
		// Ctrl+Shift+K for delete line
		if modifiers&tcell.ModShift != 0 {
			return &KeyEvent{Action: KeyActionDeleteLine, Key: key, Modifiers: modifiers}
		}
		return &KeyEvent{Action: KeyActionNone, Key: key, Modifiers: modifiers}
	case tcell.KeyCtrlD:
		return &KeyEvent{Action: KeyActionDuplicateLine, Key: key, Modifiers: modifiers}
	case tcell.KeyCtrlJ:
		// Ctrl+J for insert line below, Ctrl+Shift+J for insert line above
		// Note: We'll also handle Ctrl+Enter in the main loop since tcell
		// may treat it differently
		if modifiers&tcell.ModShift != 0 {
			return &KeyEvent{Action: KeyActionInsertLineAbove, Key: key, Modifiers: modifiers}
		}
		return &KeyEvent{Action: KeyActionInsertLineBelow, Key: key, Modifiers: modifiers}
	case tcell.KeyRune:
		// Check for Alt+key combinations (for menu shortcuts)
		if modifiers&tcell.ModAlt != 0 && r != 0 {
			upperR := unicode.ToUpper(r)
			return &KeyEvent{Action: KeyActionMenuAlt, Character: upperR, Key: key, Modifiers: modifiers}
		}
		// Regular character input
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
