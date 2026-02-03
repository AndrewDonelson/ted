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
		wantFound bool
		wantMenu  int // expected active menu index
	}{
		{
			name:      "File menu",
			key:       'F',
			wantFound: true,
			wantMenu:  0,
		},
		{
			name:      "Edit menu",
			key:       'E',
			wantFound: true,
			wantMenu:  1,
		},
		{
			name:      "Search menu",
			key:       'S',
			wantFound: true,
			wantMenu:  2,
		},
		{
			name:      "View menu",
			key:       'V',
			wantFound: true,
			wantMenu:  3,
		},
		{
			name:      "Help menu",
			key:       'H',
			wantFound: true,
			wantMenu:  4,
		},
		{
			name:      "non-existent key",
			key:       'X',
			wantFound: false,
		},
		{
			name:      "lowercase key also works",
			key:       'f',
			wantFound: true, // Now case-insensitive
			wantMenu:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mb.CloseMenu() // Reset state before each test
			found := mb.FindMenuByKey(tt.key)

			if found != tt.wantFound {
				t.Errorf("FindMenuByKey(%c) = %v, want %v", tt.key, found, tt.wantFound)
			}

			if tt.wantFound {
				if !mb.IsOpen() {
					t.Errorf("FindMenuByKey(%c) should open menu", tt.key)
				}
				if mb.GetActiveMenu() != tt.wantMenu {
					t.Errorf("FindMenuByKey(%c) active menu = %d, want %d", tt.key, mb.GetActiveMenu(), tt.wantMenu)
				}
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

	// Count non-separator items
	nonSepItems := 0
	for _, item := range fileMenu.Items {
		if !item.IsSeparator {
			nonSepItems++
		}
	}

	// Verify File menu has expected non-separator items
	expectedItems := []string{"New", "Open...", "Save", "Save As...", "Close", "Quit"}
	if nonSepItems != len(expectedItems) {
		t.Errorf("File menu has %d non-separator items, want %d", nonSepItems, len(expectedItems))
	}

	// Verify expected items exist (in order, ignoring separators)
	itemIdx := 0
	for _, item := range fileMenu.Items {
		if !item.IsSeparator {
			if itemIdx < len(expectedItems) && item.Label != expectedItems[itemIdx] {
				t.Errorf("File menu item %d = %q, want %q", itemIdx, item.Label, expectedItems[itemIdx])
			}
			itemIdx++
		}
	}
}

func TestMenu_Shortcuts(t *testing.T) {
	mb := NewMenuBar()
	fileMenu := mb.GetMenu(0)

	if fileMenu == nil {
		t.Fatal("GetMenu(0) returned nil")
	}

	// Verify non-separator items have shortcuts (except some may not have them)
	hasAtLeastOneShortcut := false
	for _, item := range fileMenu.Items {
		if !item.IsSeparator && item.Shortcut != "" {
			hasAtLeastOneShortcut = true
			break
		}
	}
	if !hasAtLeastOneShortcut {
		t.Error("File menu has no items with shortcuts")
	}
}

func TestMenu_Actions(t *testing.T) {
	mb := NewMenuBar()
	fileMenu := mb.GetMenu(0)

	if fileMenu == nil {
		t.Fatal("GetMenu(0) returned nil")
	}

	// Verify non-separator items have actions
	for _, item := range fileMenu.Items {
		if !item.IsSeparator && item.Action == "" {
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

func TestMenuBar_OpenCloseMenu(t *testing.T) {
	mb := NewMenuBar()

	// Initially closed
	if mb.IsOpen() {
		t.Error("Menu should be closed initially")
	}
	if mb.GetActiveMenu() != -1 {
		t.Error("Active menu should be -1 when closed")
	}

	// Open menu
	mb.OpenMenu(0)
	if !mb.IsOpen() {
		t.Error("Menu should be open after OpenMenu(0)")
	}
	if mb.GetActiveMenu() != 0 {
		t.Errorf("Active menu = %d, want 0", mb.GetActiveMenu())
	}

	// Close menu
	mb.CloseMenu()
	if mb.IsOpen() {
		t.Error("Menu should be closed after CloseMenu()")
	}
}

func TestMenuBar_Toggle(t *testing.T) {
	mb := NewMenuBar()

	// Toggle to open
	mb.Toggle()
	if !mb.IsOpen() {
		t.Error("Menu should be open after first Toggle()")
	}

	// Toggle to close
	mb.Toggle()
	if mb.IsOpen() {
		t.Error("Menu should be closed after second Toggle()")
	}
}

func TestMenuBar_MoveLeftRight(t *testing.T) {
	mb := NewMenuBar()
	mb.OpenMenu(0)

	// Move right
	mb.MoveRight()
	if mb.GetActiveMenu() != 1 {
		t.Errorf("After MoveRight, active menu = %d, want 1", mb.GetActiveMenu())
	}

	// Move right to end and wrap
	mb.MoveRight() // 2
	mb.MoveRight() // 3
	mb.MoveRight() // 4
	mb.MoveRight() // 0 (wrap)
	if mb.GetActiveMenu() != 0 {
		t.Errorf("After wrapping, active menu = %d, want 0", mb.GetActiveMenu())
	}

	// Move left to wrap
	mb.MoveLeft() // 4
	if mb.GetActiveMenu() != 4 {
		t.Errorf("After MoveLeft wrap, active menu = %d, want 4", mb.GetActiveMenu())
	}
}

func TestMenuBar_MoveUpDown(t *testing.T) {
	mb := NewMenuBar()
	mb.OpenMenu(0) // File menu

	// Initially at first non-separator item
	initialItem := mb.GetActiveItem()
	if initialItem < 0 {
		t.Error("Active item should be >= 0 after opening menu")
	}

	// Move down
	mb.MoveDown()
	if mb.GetActiveItem() <= initialItem {
		t.Error("Active item should increase after MoveDown")
	}

	// Move up
	mb.MoveUp()
	if mb.GetActiveItem() != initialItem {
		t.Errorf("Active item = %d, want %d after MoveUp", mb.GetActiveItem(), initialItem)
	}
}

func TestMenuBar_SelectItem(t *testing.T) {
	mb := NewMenuBar()
	mb.OpenMenu(0) // File menu

	// Select first item
	action := mb.SelectItem()
	if action == ActionNone {
		t.Error("SelectItem should return an action")
	}

	// Menu should be closed after selection
	if mb.IsOpen() {
		t.Error("Menu should be closed after SelectItem")
	}
}

func TestMenuBar_GetMenuPosition(t *testing.T) {
	mb := NewMenuBar()

	// First menu should be at position 1 (with padding)
	pos0 := mb.GetMenuPosition(0)
	if pos0 < 0 {
		t.Errorf("First menu position = %d, should be >= 0", pos0)
	}

	// Second menu should be after first
	pos1 := mb.GetMenuPosition(1)
	if pos1 <= pos0 {
		t.Errorf("Second menu position (%d) should be > first (%d)", pos1, pos0)
	}
}

func TestMenuBar_GetDropdownWidth(t *testing.T) {
	mb := NewMenuBar()
	mb.OpenMenu(0) // File menu

	width := mb.GetDropdownWidth()
	if width <= 0 {
		t.Error("Dropdown width should be > 0")
	}
}

func TestMenuItem_Separator(t *testing.T) {
	mb := NewMenuBar()
	fileMenu := mb.GetMenu(0)

	if fileMenu == nil {
		t.Fatal("GetMenu(0) returned nil")
	}

	// File menu should have some separators
	hasSeparator := false
	for _, item := range fileMenu.Items {
		if item.IsSeparator {
			hasSeparator = true
			break
		}
	}
	if !hasSeparator {
		t.Error("File menu should have at least one separator")
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
