// Package menu implements the menu system for the editor.
//
// Phase 0: Static menu structure only (no interaction yet).
package menu

// MenuBar represents the top menu bar.
type MenuBar struct {
	menus []Menu
}

// Menu represents a single menu (File, Edit, etc.).
type Menu struct {
	Label string
	Key   rune // Alt+Key activation
	Items []MenuItem
}

// MenuItem represents a single menu item.
type MenuItem struct {
	Label    string
	Shortcut string // e.g., "Ctrl+S"
	Action   string // Action identifier
}

// NewMenuBar creates a new menu bar with the default menus.
func NewMenuBar() *MenuBar {
	return &MenuBar{
		menus: []Menu{
			{
				Label: "File",
				Key:   'F',
				Items: []MenuItem{
					{Label: "New", Shortcut: "Ctrl+N", Action: "file.new"},
					{Label: "Open...", Shortcut: "Ctrl+O", Action: "file.open"},
					{Label: "Save", Shortcut: "Ctrl+S", Action: "file.save"},
					{Label: "Save As...", Shortcut: "Ctrl+Shift+S", Action: "file.saveas"},
					{Label: "Close", Shortcut: "Ctrl+W", Action: "file.close"},
					{Label: "Quit", Shortcut: "Ctrl+Q", Action: "file.quit"},
				},
			},
			{
				Label: "Edit",
				Key:   'E',
				Items: []MenuItem{
					{Label: "Undo", Shortcut: "Ctrl+Z", Action: "edit.undo"},
					{Label: "Redo", Shortcut: "Ctrl+Y", Action: "edit.redo"},
					{Label: "Cut", Shortcut: "Ctrl+X", Action: "edit.cut"},
					{Label: "Copy", Shortcut: "Ctrl+C", Action: "edit.copy"},
					{Label: "Paste", Shortcut: "Ctrl+V", Action: "edit.paste"},
					{Label: "Select All", Shortcut: "Ctrl+A", Action: "edit.selectall"},
				},
			},
			{
				Label: "Search",
				Key:   'S',
				Items: []MenuItem{
					{Label: "Find...", Shortcut: "Ctrl+F", Action: "search.find"},
					{Label: "Replace...", Shortcut: "Ctrl+H", Action: "search.replace"},
					{Label: "Go to Line...", Shortcut: "Ctrl+G", Action: "search.gotoline"},
				},
			},
			{
				Label: "View",
				Key:   'V',
				Items: []MenuItem{
					{Label: "Line Numbers", Shortcut: "Ctrl+L", Action: "view.linenumbers"},
					{Label: "Word Wrap", Shortcut: "Ctrl+Shift+W", Action: "view.wordwrap"},
				},
			},
			{
				Label: "Help",
				Key:   'H',
				Items: []MenuItem{
					{Label: "Keyboard Shortcuts", Action: "help.shortcuts"},
					{Label: "About", Action: "help.about"},
				},
			},
		},
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

// FindMenuByKey finds a menu by its Alt+Key shortcut.
func (mb *MenuBar) FindMenuByKey(key rune) *Menu {
	for i := range mb.menus {
		if mb.menus[i].Key == key {
			return &mb.menus[i]
		}
	}
	return nil
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
