package menu

import (
	"testing"
)

func TestNewMenuBar(t *testing.T) {
	mb := NewMenuBar()
	if mb == nil {
		t.Fatal("NewMenuBar() returned nil")
	}

	if len(mb.menus) == 0 {
		t.Error("NewMenuBar() returned empty menus")
	}
}

func TestMenuBar_GetMenus(t *testing.T) {
	mb := NewMenuBar()
	menus := mb.GetMenus()

	if len(menus) == 0 {
		t.Error("GetMenus() returned empty slice")
	}

	// Verify expected menus exist
	expectedMenus := []string{"File", "Edit", "Search", "View", "Help"}
	if len(menus) != len(expectedMenus) {
		t.Errorf("GetMenus() returned %d menus, want %d", len(menus), len(expectedMenus))
	}

	for i, expected := range expectedMenus {
		if i < len(menus) && menus[i].Label != expected {
			t.Errorf("GetMenus()[%d].Label = %q, want %q", i, menus[i].Label, expected)
		}
	}
}

func TestMenuBar_GetMenuCount(t *testing.T) {
	mb := NewMenuBar()
	count := mb.GetMenuCount()

	if count == 0 {
		t.Error("GetMenuCount() returned 0")
	}

	if count != len(mb.menus) {
		t.Errorf("GetMenuCount() = %d, want %d", count, len(mb.menus))
	}
}

func TestMenuBar_GetMenu(t *testing.T) {
	mb := NewMenuBar()

	tests := []struct {
		name      string
		index     int
		wantLabel string
		wantNil   bool
	}{
		{
			name:      "first menu",
			index:     0,
			wantLabel: "File",
			wantNil:   false,
		},
		{
			name:      "last menu",
			index:     4,
			wantLabel: "Help",
			wantNil:   false,
		},
		{
			name:    "negative index",
			index:   -1,
			wantNil: true,
		},
		{
			name:    "out of bounds",
			index:   100,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menu := mb.GetMenu(tt.index)

			if tt.wantNil {
				if menu != nil {
					t.Errorf("GetMenu(%d) = %v, want nil", tt.index, menu)
				}
				return
			}

			if menu == nil {
				t.Fatalf("GetMenu(%d) returned nil", tt.index)
			}

			if menu.Label != tt.wantLabel {
				t.Errorf("GetMenu(%d).Label = %q, want %q", tt.index, menu.Label, tt.wantLabel)
			}
		})
	}
}

func TestMenuBar_FindMenuByKey(t *testing.T) {
	mb := NewMenuBar()

	tests := []struct {
		name      string
		key       rune
		wantLabel string
		wantNil   bool
	}{
		{
			name:      "File menu",
			key:       'F',
			wantLabel: "File",
			wantNil:   false,
		},
		{
			name:      "Edit menu",
			key:       'E',
			wantLabel: "Edit",
			wantNil:   false,
		},
		{
			name:      "Search menu",
			key:       'S',
			wantLabel: "Search",
			wantNil:   false,
		},
		{
			name:      "View menu",
			key:       'V',
			wantLabel: "View",
			wantNil:   false,
		},
		{
			name:      "Help menu",
			key:       'H',
			wantLabel: "Help",
			wantNil:   false,
		},
		{
			name:    "non-existent key",
			key:     'X',
			wantNil: true,
		},
		{
			name:    "lowercase key",
			key:     'f',
			wantNil: true, // Keys are case-sensitive
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menu := mb.FindMenuByKey(tt.key)

			if tt.wantNil {
				if menu != nil {
					t.Errorf("FindMenuByKey(%c) = %v, want nil", tt.key, menu)
				}
				return
			}

			if menu == nil {
				t.Fatalf("FindMenuByKey(%c) returned nil", tt.key)
			}

			if menu.Label != tt.wantLabel {
				t.Errorf("FindMenuByKey(%c).Label = %q, want %q", tt.key, menu.Label, tt.wantLabel)
			}
		})
	}
}

func TestMenuBar_GetMenuLabels(t *testing.T) {
	mb := NewMenuBar()
	labels := mb.GetMenuLabels()

	if labels == "" {
		t.Error("GetMenuLabels() returned empty string")
	}

	// Should contain all menu labels
	expectedLabels := []string{"File", "Edit", "Search", "View", "Help"}
	for _, expected := range expectedLabels {
		if !contains(labels, expected) {
			t.Errorf("GetMenuLabels() = %q, want to contain %q", labels, expected)
		}
	}
}

func TestMenu_Items(t *testing.T) {
	mb := NewMenuBar()
	fileMenu := mb.GetMenu(0)

	if fileMenu == nil {
		t.Fatal("GetMenu(0) returned nil")
	}

	if len(fileMenu.Items) == 0 {
		t.Error("File menu has no items")
	}

	// Verify File menu has expected items
	expectedItems := []string{"New", "Open...", "Save", "Save As...", "Close", "Quit"}
	if len(fileMenu.Items) != len(expectedItems) {
		t.Errorf("File menu has %d items, want %d", len(fileMenu.Items), len(expectedItems))
	}

	for i, expected := range expectedItems {
		if i < len(fileMenu.Items) && fileMenu.Items[i].Label != expected {
			t.Errorf("File menu.Items[%d].Label = %q, want %q", i, fileMenu.Items[i].Label, expected)
		}
	}
}

func TestMenu_Shortcuts(t *testing.T) {
	mb := NewMenuBar()
	fileMenu := mb.GetMenu(0)

	if fileMenu == nil {
		t.Fatal("GetMenu(0) returned nil")
	}

	// Verify shortcuts exist
	for _, item := range fileMenu.Items {
		if item.Shortcut == "" {
			t.Errorf("Menu item %q has no shortcut", item.Label)
		}
	}
}

func TestMenu_Actions(t *testing.T) {
	mb := NewMenuBar()
	fileMenu := mb.GetMenu(0)

	if fileMenu == nil {
		t.Fatal("GetMenu(0) returned nil")
	}

	// Verify actions exist
	for _, item := range fileMenu.Items {
		if item.Action == "" {
			t.Errorf("Menu item %q has no action", item.Label)
		}
	}
}

func TestMenuBar_AllMenusHaveItems(t *testing.T) {
	mb := NewMenuBar()
	menus := mb.GetMenus()

	for _, menu := range menus {
		if len(menu.Items) == 0 {
			t.Errorf("Menu %q has no items", menu.Label)
		}
	}
}

func TestMenuBar_AllMenusHaveKeys(t *testing.T) {
	mb := NewMenuBar()
	menus := mb.GetMenus()

	for _, menu := range menus {
		if menu.Key == 0 {
			t.Errorf("Menu %q has no key", menu.Label)
		}
	}
}

func TestMenuBar_MenuKeysAreUnique(t *testing.T) {
	mb := NewMenuBar()
	menus := mb.GetMenus()

	keys := make(map[rune]bool)
	for _, menu := range menus {
		if keys[menu.Key] {
			t.Errorf("Duplicate menu key: %c", menu.Key)
		}
		keys[menu.Key] = true
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
