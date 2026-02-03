package terminal

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestProcessEvent(t *testing.T) {
	tests := []struct {
		name       string
		ev         tcell.Event
		wantAction KeyAction
		wantChar   rune
		wantNil    bool
	}{
		{
			name:       "character input",
			ev:         tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone),
			wantAction: KeyActionCharacter,
			wantChar:   'a',
			wantNil:    false,
		},
		{
			name:       "Ctrl+S",
			ev:         tcell.NewEventKey(tcell.KeyCtrlS, 0, tcell.ModNone),
			wantAction: KeyActionSave,
			wantChar:   0,
			wantNil:    false,
		},
		{
			name:       "Ctrl+Q",
			ev:         tcell.NewEventKey(tcell.KeyCtrlQ, 0, tcell.ModNone),
			wantAction: KeyActionQuit,
			wantChar:   0,
			wantNil:    false,
		},
		{
			name:       "arrow left",
			ev:         tcell.NewEventKey(tcell.KeyLeft, 0, tcell.ModNone),
			wantAction: KeyActionMoveLeft,
			wantChar:   0,
			wantNil:    false,
		},
		{
			name:       "arrow right",
			ev:         tcell.NewEventKey(tcell.KeyRight, 0, tcell.ModNone),
			wantAction: KeyActionMoveRight,
			wantChar:   0,
			wantNil:    false,
		},
		{
			name:       "arrow up",
			ev:         tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone),
			wantAction: KeyActionMoveUp,
			wantChar:   0,
			wantNil:    false,
		},
		{
			name:       "arrow down",
			ev:         tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone),
			wantAction: KeyActionMoveDown,
			wantChar:   0,
			wantNil:    false,
		},
		{
			name:       "backspace",
			ev:         tcell.NewEventKey(tcell.KeyBackspace, 0, tcell.ModNone),
			wantAction: KeyActionBackspace,
			wantChar:   0,
			wantNil:    false,
		},
		{
			name:       "backspace2",
			ev:         tcell.NewEventKey(tcell.KeyBackspace2, 0, tcell.ModNone),
			wantAction: KeyActionBackspace,
			wantChar:   0,
			wantNil:    false,
		},
		{
			name:       "delete",
			ev:         tcell.NewEventKey(tcell.KeyDelete, 0, tcell.ModNone),
			wantAction: KeyActionDelete,
			wantChar:   0,
			wantNil:    false,
		},
		{
			name:       "enter",
			ev:         tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone),
			wantAction: KeyActionEnter,
			wantChar:   '\n',
			wantNil:    false,
		},
		{
			name:       "home",
			ev:         tcell.NewEventKey(tcell.KeyHome, 0, tcell.ModNone),
			wantAction: KeyActionHome,
			wantChar:   0,
			wantNil:    false,
		},
		{
			name:       "end",
			ev:         tcell.NewEventKey(tcell.KeyEnd, 0, tcell.ModNone),
			wantAction: KeyActionEnd,
			wantChar:   0,
			wantNil:    false,
		},
		{
			name:       "non-keyboard event",
			ev:         tcell.NewEventResize(80, 24),
			wantAction: KeyActionNone,
			wantChar:   0,
			wantNil:    true,
		},
		{
			name:       "F1 key (help)",
			ev:         tcell.NewEventKey(tcell.KeyF1, 0, tcell.ModNone),
			wantAction: KeyActionHelp,
			wantChar:   0,
			wantNil:    false,
		},
		{
			name:       "unicode character",
			ev:         tcell.NewEventKey(tcell.KeyRune, '世', tcell.ModNone),
			wantAction: KeyActionCharacter,
			wantChar:   '世',
			wantNil:    false,
		},
		{
			name:       "zero rune",
			ev:         tcell.NewEventKey(tcell.KeyRune, 0, tcell.ModNone),
			wantAction: KeyActionNone,
			wantChar:   0,
			wantNil:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProcessEvent(tt.ev)

			if tt.wantNil {
				if got != nil {
					t.Errorf("ProcessEvent() = %v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Fatal("ProcessEvent() returned nil")
			}

			if got.Action != tt.wantAction {
				t.Errorf("ProcessEvent() Action = %v, want %v", got.Action, tt.wantAction)
			}

			if got.Character != tt.wantChar {
				t.Errorf("ProcessEvent() Character = %c, want %c", got.Character, tt.wantChar)
			}
		})
	}
}

func TestKeyEvent_IsPrintable(t *testing.T) {
	tests := []struct {
		name string
		ke   *KeyEvent
		want bool
	}{
		{
			name: "printable character",
			ke:   &KeyEvent{Action: KeyActionCharacter, Character: 'a'},
			want: true,
		},
		{
			name: "printable unicode",
			ke:   &KeyEvent{Action: KeyActionCharacter, Character: '世'},
			want: true,
		},
		{
			name: "zero character",
			ke:   &KeyEvent{Action: KeyActionCharacter, Character: 0},
			want: false,
		},
		{
			name: "non-character action",
			ke:   &KeyEvent{Action: KeyActionMoveLeft, Character: 'a'},
			want: false,
		},
		{
			name: "enter action",
			ke:   &KeyEvent{Action: KeyActionEnter, Character: '\n'},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ke.IsPrintable()
			if got != tt.want {
				t.Errorf("IsPrintable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyEvent_HasModifier(t *testing.T) {
	tests := []struct {
		name string
		ke   *KeyEvent
		mod  tcell.ModMask
		want bool
	}{
		{
			name: "has ctrl modifier",
			ke:   &KeyEvent{Modifiers: tcell.ModCtrl},
			mod:  tcell.ModCtrl,
			want: true,
		},
		{
			name: "has alt modifier",
			ke:   &KeyEvent{Modifiers: tcell.ModAlt},
			mod:  tcell.ModAlt,
			want: true,
		},
		{
			name: "has shift modifier",
			ke:   &KeyEvent{Modifiers: tcell.ModShift},
			mod:  tcell.ModShift,
			want: true,
		},
		{
			name: "no modifier",
			ke:   &KeyEvent{Modifiers: tcell.ModNone},
			mod:  tcell.ModCtrl,
			want: false,
		},
		{
			name: "multiple modifiers",
			ke:   &KeyEvent{Modifiers: tcell.ModCtrl | tcell.ModShift},
			mod:  tcell.ModCtrl,
			want: true,
		},
		{
			name: "multiple modifiers - check shift",
			ke:   &KeyEvent{Modifiers: tcell.ModCtrl | tcell.ModShift},
			mod:  tcell.ModShift,
			want: true,
		},
		{
			name: "multiple modifiers - check alt",
			ke:   &KeyEvent{Modifiers: tcell.ModCtrl | tcell.ModShift},
			mod:  tcell.ModAlt,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ke.HasModifier(tt.mod)
			if got != tt.want {
				t.Errorf("HasModifier(%v) = %v, want %v", tt.mod, got, tt.want)
			}
		})
	}
}

func TestProcessEvent_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		ev       tcell.Event
		wantNil  bool
		validate func(*testing.T, *KeyEvent)
	}{
		{
			name:    "nil event",
			ev:      nil,
			wantNil: true,
		},
		{
			name: "key with modifiers",
			ev:   tcell.NewEventKey(tcell.KeyRune, 'A', tcell.ModShift),
			validate: func(t *testing.T, ke *KeyEvent) {
				if ke == nil {
					t.Fatal("ProcessEvent() returned nil")
				}
				// Note: tcell may not preserve modifiers for KeyRune events
				// This is expected behavior - modifiers are typically handled separately
				if ke.Action != KeyActionCharacter {
					t.Errorf("Action = %v, want KeyActionCharacter", ke.Action)
				}
				if ke.Character != 'A' {
					t.Errorf("Character = %c, want 'A'", ke.Character)
				}
			},
		},
		{
			name: "special keys with modifiers",
			ev:   tcell.NewEventKey(tcell.KeyCtrlS, 0, tcell.ModCtrl),
			validate: func(t *testing.T, ke *KeyEvent) {
				if ke == nil {
					t.Fatal("ProcessEvent() returned nil")
				}
				if ke.Action != KeyActionSave {
					t.Errorf("Action = %v, want KeyActionSave", ke.Action)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProcessEvent(tt.ev)

			if tt.wantNil {
				if got != nil {
					t.Errorf("ProcessEvent() = %v, want nil", got)
				}
				return
			}

			if tt.validate != nil {
				tt.validate(t, got)
			}
		})
	}
}
