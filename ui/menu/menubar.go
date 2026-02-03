// Package menu implements the menu system for the editor.
package menu

// MenuAction represents the action triggered by a menu item.
type MenuAction string

const (
	// File menu actions
	ActionFileNew    MenuAction = "file.new"
	ActionFileOpen   MenuAction = "file.open"
	ActionFileSave   MenuAction = "file.save"
	ActionFileSaveAs MenuAction = "file.saveas"
	ActionFileClose  MenuAction = "file.close"
	ActionFileQuit   MenuAction = "file.quit"

	// Edit menu actions
	ActionEditUndo          MenuAction = "edit.undo"
	ActionEditRedo          MenuAction = "edit.redo"
	ActionEditCut           MenuAction = "edit.cut"
	ActionEditCopy          MenuAction = "edit.copy"
	ActionEditPaste         MenuAction = "edit.paste"
	ActionEditSelectAll     MenuAction = "edit.selectall"
	ActionEditDeleteLine    MenuAction = "edit.deleteline"
	ActionEditDuplicateLine MenuAction = "edit.duplicateline"
	ActionEditMoveLineUp    MenuAction = "edit.movelineup"
	ActionEditMoveLineDown  MenuAction = "edit.movelinedown"

	// Search menu actions
	ActionSearchFind     MenuAction = "search.find"
	ActionSearchReplace  MenuAction = "search.replace"
	ActionSearchGoToLine MenuAction = "search.gotoline"

	// View menu actions
	ActionViewLineNumbers MenuAction = "view.linenumbers"
	ActionViewWordWrap    MenuAction = "view.wordwrap"

	// Help menu actions
	ActionHelpShortcuts MenuAction = "help.shortcuts"
	ActionHelpAbout     MenuAction = "help.about"

	// No action
	ActionNone MenuAction = ""
)

// MenuBar represents the top menu bar.
type MenuBar struct {
	menus         []Menu
	isOpen        bool  // Whether any menu is open
	activeMenu    int   // Index of the active/hovered menu (-1 for none)
	activeItem    int   // Index of the active/hovered item (-1 for none)
	menuPositions []int // X positions of each menu label for click detection
}

// Menu represents a single menu (File, Edit, etc.).
type Menu struct {
	Label string
	Key   rune // Alt+Key activation
	Items []MenuItem
}

// MenuItem represents a single menu item.
type MenuItem struct {
	Label       string
	Shortcut    string     // e.g., "Ctrl+S"
	Action      MenuAction // Action identifier
	IsSeparator bool       // Whether this is a separator line
}

// NewMenuBar creates a new menu bar with the default menus.
func NewMenuBar() *MenuBar {
	mb := &MenuBar{
		menus: []Menu{
			{
				Label: "File",
				Key:   'F',
				Items: []MenuItem{
					{Label: "New", Shortcut: "Ctrl+N", Action: ActionFileNew},
					{Label: "Open...", Shortcut: "Ctrl+O", Action: ActionFileOpen},
					{IsSeparator: true},
					{Label: "Save", Shortcut: "Ctrl+S", Action: ActionFileSave},
					{Label: "Save As...", Shortcut: "Ctrl+Shift+S", Action: ActionFileSaveAs},
					{IsSeparator: true},
					{Label: "Close", Shortcut: "Ctrl+W", Action: ActionFileClose},
					{Label: "Quit", Shortcut: "Ctrl+Q", Action: ActionFileQuit},
				},
			},
			{
				Label: "Edit",
				Key:   'E',
				Items: []MenuItem{
					{Label: "Undo", Shortcut: "Ctrl+Z", Action: ActionEditUndo},
					{Label: "Redo", Shortcut: "Ctrl+Y", Action: ActionEditRedo},
					{IsSeparator: true},
					{Label: "Cut", Shortcut: "Ctrl+X", Action: ActionEditCut},
					{Label: "Copy", Shortcut: "Ctrl+C", Action: ActionEditCopy},
					{Label: "Paste", Shortcut: "Ctrl+V", Action: ActionEditPaste},
					{IsSeparator: true},
					{Label: "Select All", Shortcut: "Ctrl+A", Action: ActionEditSelectAll},
					{IsSeparator: true},
					{Label: "Delete Line", Shortcut: "Ctrl+Shift+K", Action: ActionEditDeleteLine},
					{Label: "Duplicate Line", Shortcut: "Ctrl+D", Action: ActionEditDuplicateLine},
					{Label: "Move Line Up", Shortcut: "Alt+Up", Action: ActionEditMoveLineUp},
					{Label: "Move Line Down", Shortcut: "Alt+Down", Action: ActionEditMoveLineDown},
				},
			},
			{
				Label: "Search",
				Key:   'S',
				Items: []MenuItem{
					{Label: "Find...", Shortcut: "Ctrl+F", Action: ActionSearchFind},
					{Label: "Replace...", Shortcut: "Ctrl+H", Action: ActionSearchReplace},
					{IsSeparator: true},
					{Label: "Go to Line...", Shortcut: "Ctrl+G", Action: ActionSearchGoToLine},
				},
			},
			{
				Label: "View",
				Key:   'V',
				Items: []MenuItem{
					{Label: "Toggle Line Numbers", Shortcut: "Ctrl+L", Action: ActionViewLineNumbers},
					{Label: "Toggle Word Wrap", Shortcut: "", Action: ActionViewWordWrap},
				},
			},
			{
				Label: "Help",
				Key:   'H',
				Items: []MenuItem{
					{Label: "Keyboard Shortcuts", Shortcut: "F1", Action: ActionHelpShortcuts},
					{IsSeparator: true},
					{Label: "About Ted", Action: ActionHelpAbout},
				},
			},
		},
		isOpen:     false,
		activeMenu: -1,
		activeItem: -1,
	}
	mb.calculateMenuPositions()
	return mb
}

// calculateMenuPositions calculates the X positions for each menu.
func (mb *MenuBar) calculateMenuPositions() {
	mb.menuPositions = make([]int, len(mb.menus))
	x := 1 // Start with 1 character padding
	for i, menu := range mb.menus {
		mb.menuPositions[i] = x
		x += len(menu.Label) + 2 // Label + 2 spaces padding
	}
}

// GetMenus returns all menus.
func (mb *MenuBar) GetMenus() []Menu {
	return mb.menus
}

// GetMenuCount returns the number of menus.
func (mb *MenuBar) GetMenuCount() int {
	return len(mb.menus)
}

// GetMenu returns the menu at the given index.
func (mb *MenuBar) GetMenu(index int) *Menu {
	if index < 0 || index >= len(mb.menus) {
		return nil
	}
	return &mb.menus[index]
}

// IsOpen returns whether a menu is currently open.
func (mb *MenuBar) IsOpen() bool {
	return mb.isOpen
}

// GetActiveMenu returns the index of the active menu (-1 if none).
func (mb *MenuBar) GetActiveMenu() int {
	return mb.activeMenu
}

// GetActiveItem returns the index of the active item (-1 if none).
func (mb *MenuBar) GetActiveItem() int {
	return mb.activeItem
}

// GetMenuPosition returns the X position of a menu by index.
func (mb *MenuBar) GetMenuPosition(index int) int {
	if index < 0 || index >= len(mb.menuPositions) {
		return 0
	}
	return mb.menuPositions[index]
}

// OpenMenu opens the menu at the given index.
func (mb *MenuBar) OpenMenu(index int) {
	if index >= 0 && index < len(mb.menus) {
		mb.isOpen = true
		mb.activeMenu = index
		mb.activeItem = 0 // Select first item
		// Skip separators
		mb.skipSeparators(1)
	}
}

// CloseMenu closes any open menu.
func (mb *MenuBar) CloseMenu() {
	mb.isOpen = false
	mb.activeMenu = -1
	mb.activeItem = -1
}

// Toggle toggles the menu bar open/closed.
func (mb *MenuBar) Toggle() {
	if mb.isOpen {
		mb.CloseMenu()
	} else {
		mb.OpenMenu(0) // Open first menu
	}
}

// MoveLeft moves to the previous menu.
func (mb *MenuBar) MoveLeft() {
	if !mb.isOpen || mb.activeMenu < 0 {
		return
	}
	mb.activeMenu--
	if mb.activeMenu < 0 {
		mb.activeMenu = len(mb.menus) - 1
	}
	mb.activeItem = 0
	mb.skipSeparators(1)
}

// MoveRight moves to the next menu.
func (mb *MenuBar) MoveRight() {
	if !mb.isOpen || mb.activeMenu < 0 {
		return
	}
	mb.activeMenu++
	if mb.activeMenu >= len(mb.menus) {
		mb.activeMenu = 0
	}
	mb.activeItem = 0
	mb.skipSeparators(1)
}

// MoveUp moves to the previous item in the current menu.
func (mb *MenuBar) MoveUp() {
	if !mb.isOpen || mb.activeMenu < 0 {
		return
	}
	menu := &mb.menus[mb.activeMenu]
	mb.activeItem--
	if mb.activeItem < 0 {
		mb.activeItem = len(menu.Items) - 1
	}
	mb.skipSeparators(-1)
}

// MoveDown moves to the next item in the current menu.
func (mb *MenuBar) MoveDown() {
	if !mb.isOpen || mb.activeMenu < 0 {
		return
	}
	menu := &mb.menus[mb.activeMenu]
	mb.activeItem++
	if mb.activeItem >= len(menu.Items) {
		mb.activeItem = 0
	}
	mb.skipSeparators(1)
}

// skipSeparators skips separator items in the given direction.
func (mb *MenuBar) skipSeparators(direction int) {
	if mb.activeMenu < 0 || mb.activeMenu >= len(mb.menus) {
		return
	}
	menu := &mb.menus[mb.activeMenu]
	if len(menu.Items) == 0 {
		return
	}

	// Prevent infinite loop if all items are separators
	count := 0
	for count < len(menu.Items) && menu.Items[mb.activeItem].IsSeparator {
		mb.activeItem += direction
		if mb.activeItem < 0 {
			mb.activeItem = len(menu.Items) - 1
		} else if mb.activeItem >= len(menu.Items) {
			mb.activeItem = 0
		}
		count++
	}
}

// SelectItem returns the action of the currently selected item and closes the menu.
func (mb *MenuBar) SelectItem() MenuAction {
	if !mb.isOpen || mb.activeMenu < 0 || mb.activeItem < 0 {
		return ActionNone
	}
	menu := &mb.menus[mb.activeMenu]
	if mb.activeItem >= len(menu.Items) {
		return ActionNone
	}
	item := &menu.Items[mb.activeItem]
	if item.IsSeparator {
		return ActionNone
	}
	action := item.Action
	mb.CloseMenu()
	return action
}

// HandleClick handles a click at the given screen coordinates.
// Returns true if the click was handled, and the action if an item was selected.
func (mb *MenuBar) HandleClick(x, y int) (handled bool, action MenuAction) {
	// Check if click is in menu bar (y == 0)
	if y == 0 {
		// Find which menu was clicked
		for i := range mb.menus {
			startX := mb.menuPositions[i]
			endX := startX + len(mb.menus[i].Label)
			if x >= startX && x < endX {
				if mb.isOpen && mb.activeMenu == i {
					// Clicking on open menu closes it
					mb.CloseMenu()
				} else {
					mb.OpenMenu(i)
				}
				return true, ActionNone
			}
		}
		// Click on menu bar but not on any menu - close any open menu
		if mb.isOpen {
			mb.CloseMenu()
			return true, ActionNone
		}
		return false, ActionNone
	}

	// Check if click is in dropdown menu
	if mb.isOpen && mb.activeMenu >= 0 {
		menu := &mb.menus[mb.activeMenu]
		menuX := mb.menuPositions[mb.activeMenu]
		menuWidth := mb.getDropdownWidth(mb.activeMenu)

		// Check if within dropdown bounds
		if x >= menuX && x < menuX+menuWidth && y > 0 && y <= len(menu.Items) {
			itemIndex := y - 1
			if itemIndex >= 0 && itemIndex < len(menu.Items) && !menu.Items[itemIndex].IsSeparator {
				mb.activeItem = itemIndex
				return true, mb.SelectItem()
			}
		}

		// Click outside dropdown - close it
		mb.CloseMenu()
		return true, ActionNone
	}

	return false, ActionNone
}

// getDropdownWidth calculates the width of the dropdown menu.
func (mb *MenuBar) getDropdownWidth(menuIndex int) int {
	if menuIndex < 0 || menuIndex >= len(mb.menus) {
		return 0
	}
	menu := &mb.menus[menuIndex]
	maxWidth := 0
	for _, item := range menu.Items {
		width := len(item.Label)
		if item.Shortcut != "" {
			width += 4 + len(item.Shortcut) // 4 spaces between label and shortcut
		}
		if width > maxWidth {
			maxWidth = width
		}
	}
	return maxWidth + 4 // Add padding
}

// GetDropdownWidth returns the width of the dropdown for the active menu.
func (mb *MenuBar) GetDropdownWidth() int {
	return mb.getDropdownWidth(mb.activeMenu)
}

// FindMenuByKey finds a menu by its Alt+Key shortcut and opens it.
// Returns true if a menu was found and opened.
func (mb *MenuBar) FindMenuByKey(key rune) bool {
	// Convert to uppercase for comparison
	upperKey := key
	if key >= 'a' && key <= 'z' {
		upperKey = key - 32
	}

	for i := range mb.menus {
		menuKey := mb.menus[i].Key
		if menuKey >= 'a' && menuKey <= 'z' {
			menuKey = menuKey - 32
		}
		if menuKey == upperKey {
			mb.OpenMenu(i)
			return true
		}
	}
	return false
}

// GetMenuLabels returns the labels of all menus for rendering.
func (mb *MenuBar) GetMenuLabels() string {
	var labels []string
	for _, menu := range mb.menus {
		labels = append(labels, menu.Label)
	}
	// Join with double spaces for visual separation
	result := ""
	for i, label := range labels {
		if i > 0 {
			result += "  "
		}
		result += label
	}
	return result
}
