# ted - Terminal EDitor

## Project Overview
**ted** is a modern, cross-platform command-line text editor written in Go that uses familiar Windows-style keyboard shortcuts and intuitive arrow key navigation.

**Target Platforms:** Linux, macOS, Windows  
**Language:** Go  
**Philosophy:** Easy to learn, familiar shortcuts, no modal editing

**Name:** ted = Terminal EDitor - friendly, approachable, easy to remember

---

## Core Design Principles

1. **Familiar Shortcuts** - Windows/modern editor keybindings (Ctrl+S, Ctrl+C, Ctrl+V, etc.)
2. **No Modes** - Direct editing, no insert/command mode like vim
3. **Free Navigation** - Arrow keys work as expected everywhere
4. **Cross-Platform** - Identical experience on Linux, macOS, Windows
5. **Fast & Lightweight** - Written in Go for speed and minimal dependencies
6. **Terminal Native** - Works in any terminal, no GUI required

---

## Core Features

### Essential Editing Operations
- [x] Cut, Copy, Paste
- [x] Undo, Redo (with history depth)
- [x] Select from cursor to line start/end
- [x] Insert new line above/below current line
- [x] Delete entire line (fast)
- [ ] Duplicate current line
- [ ] Move line up/down
- [ ] Delete word forward/backward
- [ ] Select current word (double-click equivalent)

### Selection & Navigation
- [x] Free arrow key movement
- [ ] Shift+Arrow for selection
- [ ] Ctrl+Arrow to jump by word
- [ ] Ctrl+Home/End for document start/end
- [ ] Home/End for line start/end
- [ ] Page Up/Down navigation
- [ ] Go to line number (Ctrl+G)
- [ ] Jump to matching bracket/parenthesis

### File Operations
- [ ] Open file (single file initially)
- [ ] Save (Ctrl+S)
- [ ] Save As
- [ ] New file
- [ ] Close/Exit with unsaved changes prompt
- [ ] Recent files list
- [ ] Auto-save/backup options

### Search & Replace
- [ ] Find (Ctrl+F)
- [ ] Find next/previous
- [ ] Replace (Ctrl+H)
- [ ] Replace all
- [ ] Case-sensitive toggle
- [ ] Whole word toggle
- [ ] Regex support (optional)

### Code-Friendly Features
- [ ] Syntax highlighting (configurable)
- [ ] Auto-indentation
- [ ] Tab/Spaces conversion
- [ ] Comment/Uncomment line or block (Ctrl+/)
- [ ] Trim trailing whitespace
- [ ] Show whitespace characters (toggle)
- [ ] Multiple indentation levels (Shift+Tab)

### Display & UI
- [ ] Line numbers (toggleable)
- [ ] Current line highlighting
- [ ] Status bar (line:col, file info, mode indicators)
- [ ] Ruler/column indicator
- [ ] Word wrap (soft wrap toggle)
- [ ] Split view (horizontal/vertical - future)

### Advanced (Consider for MVP or later)
- [ ] Multiple files/tabs
- [ ] Block/column selection (Alt+Shift+Arrow)
- [ ] Bookmarks
- [ ] Macros/command recording
- [ ] Read-only mode
- [ ] File encoding detection/conversion
- [ ] Large file handling (lazy loading)

---

## MVP Feature Priority

### ✅ Aligned with Phased Development Plan

The features are now organized into **6 development phases** (see Phased Development Plan above):

- **Phase 0** (Week 1): Foundation - Basic text editing
- **Phase 1** (Week 2): Essential editing - Undo/redo, clipboard, line ops
- **Phase 2** (Week 3): Navigation - Search, line numbers, fast navigation
- **Phase 3** (Week 4): Code editing - Syntax highlighting, comments
- **Phase 4** (Week 5): Advanced search - Replace, regex
- **Phase 5** (Week 6): Configuration - Customization, themes
- **Phase 6** (Week 7+): Multiple files - Tabs/buffers

Each phase builds on previous phases and adds a coherent set of related features.

---

## Technology Decisions

### Terminal Library: tcell
**Chosen:** `github.com/gdamore/tcell/v2`

**Rationale:**
- Cross-platform (Linux, macOS, Windows)
- Low-level control (raw terminal access)
- Active maintenance
- Good performance
- Well-documented
- Event-driven model

**Alternatives Considered:**
- `termbox-go` - Less maintained, fewer features
- `bubbletea` - Too high-level, opinionated framework (good for TUIs, not editors)
- Raw termios - Too low-level, platform-specific pain

### Clipboard Library
**Chosen:** `github.com/atotto/clipboard`

**Rationale:**
- Simple, cross-platform clipboard access
- No external dependencies
- Works on Linux, macOS, Windows

### Syntax Highlighting
**Phase 3 Decision:**
- Option A: `github.com/alecthomas/chroma` - Full-featured, many languages
- Option B: Custom implementation - Lighter weight, more control
- **Recommend:** Start with Chroma, can optimize later if needed

### Configuration Format
**Chosen:** TOML (`.tedrc` or `ted.toml`)

**Rationale:**
- Human-friendly and easy to read
- Good for configuration files
- Well-supported in Go ecosystem
- Clear syntax (better than YAML's whitespace sensitivity)
- More human-friendly than JSON

**Location:**
- `~/.tedrc` (simple, traditional)
- `~/.config/ted/config.toml` (XDG standard)
- Support both, prefer XDG location

---

## Build & Distribution

### Build Process
```bash
# Development build
go build -o ted ./cmd/ted

# Production build (all platforms)
GOOS=linux GOARCH=amd64 go build -o ted-linux-amd64
GOOS=darwin GOARCH=amd64 go build -o ted-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o ted-darwin-arm64
GOOS=windows GOARCH=amd64 go build -o ted-windows-amd64.exe
```

### Installation
```bash
# Install to $GOPATH/bin or /usr/local/bin
go install github.com/yourusername/ted@latest

# Or use Homebrew (future)
brew install ted

# Or download binary from releases
```

### Dependencies
- Minimal external dependencies
- All dependencies vendored or go.mod managed
- No runtime dependencies (single binary)

---

### Unit Testing (Per Module)
Each module should have comprehensive unit tests:

```
core/buffer/buffer_test.go
core/buffer/cursor_test.go
core/buffer/selection_test.go
core/buffer/history_test.go
...
```

**Test Coverage Goals:**
- Core modules: 90%+ coverage
- UI modules: 70%+ coverage (harder to test)
- Integration tests for command workflows

### Integration Testing
- Test command execution end-to-end
- Test file operations (read, modify, save)
- Test undo/redo chains
- Test clipboard integration

### Manual Testing Checklist (Per Phase)
Each phase delivery should pass:
- [ ] All automated tests pass
- [ ] Manual smoke test of new features
- [ ] No regressions in previous features
- [ ] Works on Linux, macOS, Windows
- [ ] Performance acceptable (< 100ms latency for operations)

---

## Keyboard Shortcuts

*Windows/modern editor style - using Ctrl as primary modifier (Cmd on macOS)*

### File Operations
- **Ctrl+N** - New file
- **Ctrl+O** - Open file
- **Ctrl+S** - Save
- **Ctrl+Shift+S** - Save As
- **Ctrl+W** - Close current file
- **Ctrl+Q** - Quit editor
- **Ctrl+Tab** - Next file/tab (if multiple files)
- **Ctrl+Shift+Tab** - Previous file/tab

### Essential Editing
- **Ctrl+Z** - Undo
- **Ctrl+Y** or **Ctrl+Shift+Z** - Redo
- **Ctrl+X** - Cut
- **Ctrl+C** - Copy
- **Ctrl+V** - Paste
- **Ctrl+A** - Select all
- **Delete** - Delete character forward
- **Backspace** - Delete character backward

### Line Operations
- **Ctrl+Shift+K** - Delete entire line
- **Ctrl+D** - Duplicate current line
- **Alt+Up** - Move line up
- **Alt+Down** - Move line down
- **Ctrl+Enter** - Insert line below
- **Ctrl+Shift+Enter** - Insert line above
- **Ctrl+J** - Join lines (remove line break)

### Selection Operations
- **Shift+Home** - Select from cursor to line start
- **Shift+End** - Select from cursor to line end
- **Ctrl+Shift+Home** - Select from cursor to document start
- **Ctrl+Shift+End** - Select from cursor to document end
- **Shift+Arrow Keys** - Extend selection by character
- **Ctrl+Shift+Arrow** - Extend selection by word
- **Shift+Page Up/Down** - Select by page
- **Double-click** - Select current word (mouse support)
- **Alt+Shift+Arrow** - Block/column selection

### Navigation
- **Arrow Keys** - Move cursor (up, down, left, right)
- **Home** - Go to line start (first non-whitespace)
- **End** - Go to line end
- **Ctrl+Home** - Go to document start
- **Ctrl+End** - Go to document end
- **Ctrl+Left/Right** - Move by word
- **Page Up/Down** - Move by page
- **Ctrl+G** - Go to line number

### Word/Text Operations
- **Ctrl+Backspace** - Delete word backward
- **Ctrl+Delete** - Delete word forward
- **Ctrl+T** - Transpose characters (swap character with previous)
- **Alt+Shift+Up/Down** - Duplicate selection up/down

### Search & Replace
- **Ctrl+F** - Find
- **F3** or **Ctrl+G** - Find next
- **Shift+F3** or **Ctrl+Shift+G** - Find previous
- **Ctrl+H** - Replace
- **Ctrl+Shift+H** - Replace all
- **Esc** - Close find/replace dialog

### Code Editing
- **Ctrl+/** - Toggle line comment
- **Ctrl+Shift+/** - Toggle block comment
- **Tab** - Indent (or insert tab)
- **Shift+Tab** - Unindent (dedent)
- **Ctrl+]** - Indent line/selection
- **Ctrl+[** - Unindent line/selection
- **Ctrl+Shift+F** - Format/auto-indent document
- **Ctrl+Space** - Auto-complete (future feature)
- **Ctrl+B** - Jump to matching bracket

### Display & View
- **Ctrl+L** - Toggle line numbers
- **Ctrl+Shift+W** - Toggle word wrap
- **Ctrl+Shift+I** - Toggle show whitespace
- **Ctrl+Plus/Minus** - Increase/decrease font size (if supported)
- **Ctrl+0** - Reset font size
- **F11** - Toggle fullscreen (if supported)

### Utility
- **F10** - Activate menu bar (same as Alt+Key)
- **Ctrl+K** - Command palette (future: quick commands)
- **Ctrl+P** - Quick file open (future: fuzzy finder)
- **Ctrl+Shift+P** - Command palette
- **Esc** - Cancel operation / clear selection / close menu

---

---

## Technical Architecture

### Modular Component Design

**Core Principle:** Each component should be independent, testable, and replaceable

```
ted/
├── main.go                 # Entry point, wires components together
├── core/
│   ├── buffer/            # Text buffer management
│   │   ├── buffer.go      # Core buffer operations
│   │   ├── cursor.go      # Cursor position & movement
│   │   ├── selection.go   # Selection handling
│   │   └── history.go     # Undo/redo stack
│   ├── file/              # File I/O operations
│   │   ├── reader.go      # File reading
│   │   ├── writer.go      # File writing
│   │   └── watcher.go     # File change detection (future)
│   └── clipboard/         # System clipboard integration
│       └── clipboard.go   # Cross-platform clipboard
├── editor/
│   ├── editor.go          # Main editor controller
│   ├── commands.go        # Command dispatcher
│   └── config.go          # Configuration management
├── ui/
│   ├── terminal/          # Terminal handling
│   │   ├── screen.go      # Screen management & rendering
│   │   ├── input.go       # Keyboard input handling
│   │   ├── events.go      # Event system
│   │   └── resize.go      # Window resize handler
│   ├── menu/              # Menu system
│   │   ├── menubar.go     # Top menu bar
│   │   ├── menu.go        # Individual menu
│   │   └── menuitem.go    # Menu item
│   ├── renderer/          # Display rendering
│   │   ├── text.go        # Text rendering
│   │   ├── statusbar.go   # Top status bar (right side)
│   │   ├── infobar.go     # Bottom info bar
│   │   ├── linenumbers.go # Line numbers
│   │   ├── scrollbar.go   # Scroll indicators
│   │   └── highlight.go   # Syntax highlighting
│   ├── layout/            # Layout management
│   │   ├── viewport.go    # Viewport calculations
│   │   └── dimensions.go  # Screen dimensions
│   └── theme/             # Color schemes
│       └── theme.go       # Theme management
├── search/
│   ├── finder.go          # Search implementation
│   └── replacer.go        # Replace functionality
├── syntax/
│   ├── highlighter.go     # Syntax highlighting engine
│   ├── languages/         # Language definitions
│   │   ├── go.go
│   │   ├── javascript.go
│   │   └── ...
│   └── parser.go          # Token parser
└── utils/
    ├── keys.go            # Keyboard shortcut definitions
    └── platform.go        # OS-specific utilities
```

### Module Responsibilities

**core/buffer** - Text Buffer Management
- Stores document text in memory
- Manages cursor position(s)
- Handles text selections
- Maintains undo/redo history
- Provides insert, delete, replace operations
- **No UI dependencies** - pure data structure

**core/file** - File Operations
- Read/write files from disk
- Handle different encodings (UTF-8, etc.)
- Manage file metadata
- **No UI dependencies** - pure I/O

**core/clipboard** - System Clipboard
- Cross-platform clipboard integration
- Cut, copy, paste operations
- **OS abstraction layer**

**editor/editor** - Editor Controller
- Coordinates all components
- Dispatches commands to appropriate modules
- Manages editor state
- Handles mode switching (insert/overwrite)
- **Central orchestrator**

**editor/commands** - Command System
- Maps keyboard shortcuts to actions
- Command pattern for undo/redo
- Extensible command registry

**ui/terminal** - Terminal Interface
- Raw terminal control (tcell)
- Capture keyboard/mouse input
- Manage screen updates
- Handle terminal resize events
- Event loop

**ui/menu** - Menu System
- Top menu bar rendering
- Menu item management
- Keyboard navigation (Alt+key)
- Mouse click handling
- Dropdown menus with shortcuts

**ui/renderer** - Display Rendering
- Render buffer to screen
- Draw menu bar, status indicators, info bar
- Apply syntax highlighting
- Handle scrolling viewport
- Line numbers, scrollbars

**ui/layout** - Layout Management
- Calculate viewport dimensions
- Handle responsive resizing
- Manage screen regions (menu, edit area, info bar)
- Ensure minimum sizes

**search** - Find & Replace
- Search algorithms
- Regex support
- Replace operations
- Search history

**syntax** - Syntax Highlighting
- Tokenize source code
- Apply color themes
- Language-specific rules
- Pluggable language support

---

## Phased Development Plan

### Phase 0: Foundation (Week 1)
**Goal:** Core architecture and basic text editing with full UI layout

**Modules to Build:**
- `core/buffer` - Basic buffer (insert, delete, cursor movement)
- `core/file` - Read/write single file
- `ui/terminal` - Basic screen rendering with tcell, resize handling
- `ui/layout` - Viewport calculations, screen regions
- `ui/renderer` - Simple text display with menu bar and info bar
- `ui/menu` - Static menu bar (no interaction yet)
- `editor/editor` - Wire up basic components

**Features Delivered:**
- Open a file (command line argument)
- Display text in middle area
- Static menu bar across top (File, Edit, Search, View, Help)
- Info bar at bottom (filename, size, line count)
- Arrow key navigation (up, down, left, right)
- Type/delete characters
- Save file (Ctrl+S)
- Exit (Ctrl+Q)
- **Responsive:** Adapts to terminal resize

**UI Layout:**
```
┌─────────────────────────────────────────────────────┐
│ File  Edit  Search  View  Help      [Status Area]  │ ← Menu bar
├─────────────────────────────────────────────────────┤
│ Line numbers    Editable text area                  │
│ (if enabled)                                        │
│                                                     │ ← Editing area
│                                                     │
├─────────────────────────────────────────────────────┤
│ filename.txt │ 1.2 KB │ 45 lines │ Modified        │ ← Info bar
└─────────────────────────────────────────────────────┘
```

**Success Criteria:** Can edit a simple text file with proper layout and save it

---

### Phase 1: Essential Editing (Week 2)
**Goal:** Make it actually usable for daily work with interactive menus

**Modules to Extend:**
- `core/buffer` - Add selection, clipboard integration, undo/redo
- `core/clipboard` - System clipboard support
- `editor/commands` - Command dispatcher
- `ui/menu` - Interactive menu system (keyboard + mouse)
- `ui/renderer` - Status indicators (INS/OVR mode, encoding, position)
- `ui/renderer` - Enhanced info bar, cursor highlighting

**Features Delivered:**
- **Interactive Menus:**
  - Alt+F, Alt+E, etc. to activate menus
  - Arrow keys to navigate
  - Mouse click support
  - Show keyboard shortcuts in menus
- **Editing Operations:**
  - Undo/Redo (Ctrl+Z, Ctrl+Y)
  - Cut/Copy/Paste (Ctrl+X, Ctrl+C, Ctrl+V)
  - Select with Shift+Arrow
  - Line operations (delete line, insert line above/below)
  - Duplicate line
  - Home/End keys
- **Status Indicators:**
  - Mode: INS/OVR (toggle with Insert key)
  - Encoding: UTF-8
  - Position: LN 45, COL 12
- **Enhanced Info Bar:**
  - File type detection
  - Tab size indicator
  - Line ending (LF/CRLF)
- **UI Polish:**
  - Current line highlighting
  - Selection highlighting
  - Unsaved changes indicator (*)
  - Unsaved changes prompt on exit

**Success Criteria:** Can perform standard editing tasks comfortably with menu access

---

### Phase 2: Navigation & Productivity (Week 3)
**Goal:** Fast navigation and common power features with visual feedback

**Modules to Build:**
- `search/finder` - Basic search
- `ui/renderer` - Line numbers, scrollbar indicators

**Modules to Extend:**
- `core/buffer` - Word-based operations
- `ui/terminal` - Enhanced input handling
- `ui/layout` - Scroll position management

**Features Delivered:**
- **Search:**
  - Find (Ctrl+F) with inline search box
  - Find next/previous (F3/Shift+F3)
  - Search match counter in info bar
- **Navigation:**
  - Go to line (Ctrl+G) - dialog box
  - Ctrl+Arrow (word navigation)
  - Ctrl+Home/End (document start/end)
  - Delete word (Ctrl+Backspace/Delete)
  - Move line up/down (Alt+Up/Down)
  - Page Up/Down
- **Visual Enhancements:**
  - Line numbers (toggle with Ctrl+L or View menu)
  - Scrollbar indicators (vertical, horizontal)
  - Scroll position indicator (e.g., "Top", "45%", "Bottom")
  - Search match highlighting
- **Info Bar Updates:**
  - Show search results: "Match 3 of 12"
  - Show selection count: "45 chars selected"

**Success Criteria:** Can navigate large files efficiently with clear visual feedback

---

### Phase 3: Code Editing (Week 4)
**Goal:** Syntax highlighting and code-specific features

**Modules to Build:**
- `syntax/highlighter` - Syntax highlighting engine
- `syntax/languages` - Go, JavaScript, Python, etc.
- `ui/theme` - Color scheme system

**Modules to Extend:**
- `ui/renderer` - Apply syntax highlighting
- `editor/commands` - Comment/indent commands

**Features Delivered:**
- Syntax highlighting (configurable languages)
- Auto-indentation
- Comment/uncomment (Ctrl+/)
- Tab/Shift+Tab for indent/unindent
- Show whitespace toggle
- Jump to matching bracket
- Basic color themes

**Success Criteria:** Comfortable for coding work

---

### Phase 4: Advanced Search & Replace (Week 5)
**Goal:** Powerful find/replace capabilities

**Modules to Extend:**
- `search/finder` - Regex support, case sensitivity
- `search/replacer` - Replace functionality

**Features Delivered:**
- Replace (Ctrl+H)
- Replace all
- Case-sensitive toggle
- Whole word search
- Regex support
- Search history

**Success Criteria:** Can perform complex search/replace operations

---

### Phase 5: Configuration & Polish (Week 6)
**Goal:** Make it customizable and production-ready

**Modules to Build:**
- `editor/config` - Configuration system
- `ui/theme` - Custom themes

**Modules to Extend:**
- `editor/commands` - Custom keybindings
- All modules - Apply config settings

**Features Delivered:**
- Configuration file (~/.tedrc or ted.config)
- Custom keybindings
- Custom color schemes
- Configurable tab size, spaces vs tabs
- Auto-save options
- Trim trailing whitespace
- Word wrap toggle

**Success Criteria:** Users can customize to their preferences

---

### Phase 6: Multiple Files (Week 7+)
**Goal:** Multi-file editing

**Modules to Extend:**
- `editor/editor` - Multiple buffer management
- `ui/renderer` - Tab bar or buffer list
- `editor/commands` - File switching commands

**Features Delivered:**
- Open multiple files
- Switch between files (Ctrl+Tab)
- Tab bar or buffer list
- Close individual files
- Split view (future consideration)

**Success Criteria:** Can work with multiple files simultaneously

---

### Future Phases (v1.0+)
- Block/column selection
- Macros and command recording
- Plugin system
- LSP integration (code completion, diagnostics)
- Git integration
- Remote file editing
- Fuzzy file finder
- Command palette

---

## Module Interface Contracts

### buffer.Buffer Interface
```go
type Buffer interface {
    // Content operations
    Insert(pos Position, text string) error
    Delete(start, end Position) error
    Replace(start, end Position, text string) error
    GetText(start, end Position) string
    GetLine(lineNum int) string
    
    // Cursor operations
    MoveCursor(pos Position)
    GetCursor() Position
    
    // Selection operations
    SetSelection(start, end Position)
    GetSelection() (start, end Position, exists bool)
    ClearSelection()
    
    // History operations
    Undo() bool
    Redo() bool
    
    // Metadata
    LineCount() int
    IsModified() bool
    MarkSaved()
}
```

### renderer.Renderer Interface
```go
type Renderer interface {
    // Full screen rendering
    RenderAll(editor *Editor) error
    
    // Component rendering
    RenderMenuBar(menus []Menu, active int) error
    RenderStatusBar(mode EditorMode, encoding string, pos Position) error
    RenderInfoBar(info FileInfo) error  // Uses INVERTED colors
    RenderEditArea(buffer Buffer, viewport Viewport) error
    RenderLineNumbers(startLine, endLine, currentLine int) error
    RenderScrollbars(viewport Viewport, totalLines int) error
    
    // Screen management
    Clear()
    Refresh()
    GetSize() (width, height int)
}
```

### layout.Layout Interface
```go
type Layout interface {
    // Calculate screen regions
    GetMenuBarRegion() Region
    GetEditAreaRegion() Region
    GetInfoBarRegion() Region
    GetLineNumberWidth() int
    
    // Viewport management
    CalculateViewport(cursorPos Position, totalLines int) Viewport
    AdjustForResize(newWidth, newHeight int)
    
    // Coordinate conversion
    ScreenToBuffer(screenX, screenY int) Position
    BufferToScreen(pos Position) (screenX, screenY int)
}

type Region struct {
    X, Y          int  // Top-left corner
    Width, Height int  // Dimensions
}

type Viewport struct {
    StartLine int     // First visible line
    EndLine   int     // Last visible line
    OffsetX   int     // Horizontal scroll offset
    Width     int     // Viewport width
    Height    int     // Viewport height
}
```

### menu.MenuBar Interface
```go
type MenuBar interface {
    // Rendering
    Render(activeIndex int) string
    
    // Navigation
    GetMenuCount() int
    GetMenu(index int) Menu
    FindMenuByKey(key rune) int  // Alt+Key activation
    
    // State
    IsActive() bool
    SetActive(active bool)
}

type Menu interface {
    // Properties
    GetLabel() string
    GetKey() rune  // Alt+Key shortcut
    GetItems() []MenuItem
    
    // Rendering
    Render(width int) []string
    
    // Navigation
    GetItemCount() int
    GetSelectedItem() int
    SelectNext()
    SelectPrevious()
    
    // Execution
    ExecuteSelected(editor *Editor) error
}

type MenuItem interface {
    GetLabel() string
    GetShortcut() string  // e.g., "Ctrl+S"
    IsSeparator() bool
    IsEnabled() bool
    Execute(editor *Editor) error
}
```

### commands.Command Interface
```go
type Command interface {
    Execute(editor *Editor) error
    Undo(editor *Editor) error
    Description() string
}

type CommandRegistry interface {
    Register(name string, cmd Command)
    Execute(name string, editor *Editor) error
    GetCommand(name string) Command
}
```

### editor.Editor State
```go
type Editor struct {
    // Core components
    buffer     *Buffer
    file       *File
    clipboard  *Clipboard
    
    // UI components
    layout     *Layout
    renderer   *Renderer
    menuBar    *MenuBar
    terminal   *Terminal
    
    // State
    mode       EditorMode  // Insert/Overwrite
    config     *Config
    isDirty    bool       // Unsaved changes
    
    // Event handling
    eventChan  chan Event
    quitChan   chan bool
}

type EditorMode int
const (
    ModeInsert EditorMode = iota
    ModeOverwrite
)

type FileInfo struct {
    Name        string
    Path        string
    Size        int64
    Type        string  // "Go", "JavaScript", etc.
    Encoding    string  // "UTF-8", "ASCII", etc.
    LineEnding  string  // "LF", "CRLF", "CR"
    TabSize     int
    UseSpaces   bool
    TotalLines  int
}
```

**Key Design Decisions:**
- **Interfaces over concrete types** - allows swapping implementations
- **No circular dependencies** - core modules don't depend on UI
- **Testable modules** - each can be unit tested independently
- **Clear boundaries** - each module has single responsibility
- **Event-driven** - UI events trigger commands, not direct buffer manipulation
- **Responsive layout** - Layout adapts to screen size changes
- **Stateless rendering** - Renderer receives all data it needs, doesn't hold state

---

## UI/UX Layout

### Screen Layout Specification

**Responsive Design:** Must adapt to terminal resize events dynamically

```
┌─────────────────────────────────────────────────────────────────────┐
│ File  Edit  Search  View  Help           INS │ UTF-8 │ LN 45, COL 12│ ← Top Menu Bar
├─────────────────────────────────────────────────────────────────────┤
│  1 │ package main                                                    │
│  2 │                                                                 │
│  3 │ import (                                                        │
│  4 │     "fmt"                                                       │
│  5 │ )                                                               │
│  6 │                                                                 │
│  7 │ func main() {                                                   │ ← Scrollable
│  8 │     fmt.Println("Hello, World!")█                               │   Editing
│  9 │ }                                                               │   Area
│ 10 │                                                                 │
│ 11 │                                                                 │
│    │                                                                 │
├─────────────────────────────────────────────────────────────────────┤
│▓▓main.go▓│▓245 bytes▓│▓Go▓│▓Modified▓│▓Tab Size: 4▓│▓CRLF▓▓▓▓▓▓▓▓▓▓│ ← Bottom Info Bar
└─────────────────────────────────────────────────────────────────────┘   (INVERTED TEXT)
```

**Note:** The bottom info bar uses inverted colors (▓ represents inverted text) to make it stand out distinctly from the editing area.

### Top Menu Bar (Line 1)

**Left Side - Menu Items:**
- `File` - File operations (New, Open, Save, Save As, Close, Quit)
- `Edit` - Edit operations (Undo, Redo, Cut, Copy, Paste, Select All)
- `Search` - Search operations (Find, Replace, Go to Line)
- `View` - View options (Line Numbers, Whitespace, Word Wrap)
- `Help` - Help and About

**Right Side - Status Indicators:**
- Mode indicator: `INS` (insert) or `OVR` (overwrite)
- Encoding: `UTF-8`, `ASCII`, etc.
- Position: `LN 45, COL 12` (current line and column)

**Behavior:**
- Menu items are clickable (mouse support) or keyboard activated
- `Alt+F` for File menu, `Alt+E` for Edit, etc.
- Arrow keys navigate menus when active
- `Esc` closes active menu
- Menus display keyboard shortcuts next to commands

### Middle - Scrollable Editing Area

**Layout:**
- Optional line numbers on left (toggleable)
- Vertical scrollbar indicator on right (if content exceeds viewport)
- Horizontal scrollbar indicator on bottom (if lines exceed width)
- Current line highlighting (subtle background)
- Syntax highlighting (when enabled)
- Selection highlighting (inverted colors or distinct background)
- Cursor: Block (overwrite) or Line (insert)

**Scrolling:**
- Arrow keys scroll when cursor reaches edge
- Page Up/Down for page-based scrolling
- Ctrl+Home/End for document start/end
- Mouse wheel support
- Smooth scrolling (no jumpy movements)

**Responsiveness:**
- On window resize: recalculate viewport, adjust line wrapping
- Preserve cursor position (keep visible)
- Redraw efficiently (only changed regions)

### Bottom Info Bar (Last Line)

**Visual Style:** INVERTED TEXT - Light background with dark text to stand out from editing area

**Left to Right - File Information:**
- **Filename**: `main.go` or `[No Name]` for new files
- **File size**: `245 bytes`, `2.3 KB`, `1.5 MB`
- **File type**: `Go`, `JavaScript`, `Plain Text` (auto-detected or manual)
- **Modification status**: `Modified`, `Saved`, `[Read Only]`
- **Editor settings**: `Tab Size: 4`, `Spaces` or `Tabs`
- **Line ending**: `LF` (Unix), `CRLF` (Windows), `CR` (Mac)

**Additional Info (if space allows):**
- Total lines: `234 lines`
- Selected text: `45 chars selected` (when selection active)
- Search results: `Match 3 of 12` (during search)

**Behavior:**
- Updates in real-time as user types/navigates
- Clickable items (e.g., click encoding to change) - Phase 2+
- Responsive to window width (hide less important info on narrow terminals)
- **Always uses inverted colors** to clearly separate from content area

**Color Specification:**
- Background: #d4d4d4 (light gray) - inverted from main dark background
- Text: #1e1e1e (dark gray) - inverted from main light text
- Creates strong visual distinction from editing area

---

## Visual Elements

### Line Numbers
- **Display:** Optional, toggleable via `Ctrl+L` or View menu
- **Styling:** Muted color (gray), right-aligned
- **Width:** Dynamic based on total line count (e.g., 3 digits for 100-999 lines)
- **Current line:** Highlighted or different color

### Current Line Highlighting
- **Subtle background:** Slightly different shade from main background
- **Entire line:** Full width highlighting
- **Not distracting:** Low contrast change

### Selection Highlighting
- **Visual:** Inverted colors or distinct background color
- **Persistent:** Remains visible during typing (replaced on insert)
- **Multi-line:** Supports selecting across many lines

### Cursor Style
- **Insert mode:** Vertical line (|) or thin block
- **Overwrite mode:** Full block (█)
- **Blinking:** Optional (configurable)
- **Always visible:** Ensures cursor is in viewport

### Scrollbar Indicators
- **Vertical:** Right edge, shows position in document
- **Horizontal:** Bottom edge (if needed), shows horizontal position
- **Minimal:** Single character indicators, not intrusive
- **Example:** `↑ │ ↓` for vertical, `← ─ →` for horizontal

### Color Scheme

**Phase 0-4: Dark Mode Only**
- Background: Dark gray (#1e1e1e)
- Text: Light gray (#d4d4d4)
- Menu bar: Slightly lighter background (#252525)
- **Info bar: INVERTED - Light background (#d4d4d4) with dark text (#1e1e1e)**
- Current line: Subtle highlight (#2a2a2a)
- Selection: Blue background (#264f78)
- Line numbers: Muted gray (#858585)
- Cursor: White or bright color (#ffffff)
- Comments (syntax): Green (#6a9955)
- Keywords (syntax): Blue (#569cd6)
- Strings (syntax): Orange (#ce9178)

**Info Bar Styling:**
- Uses inverted colors to make it stand out from editing area
- Background: #d4d4d4 (light gray)
- Text: #1e1e1e (dark gray)
- Clear visual separation from content

**Phase 5+: Light Mode Option**
- Will be added as configurable theme
- Inverted colors with good contrast
- Info bar will also be inverted (dark bar on light theme)
- User can switch in configuration file

---

## Menu System Design

### Menu Activation
- **Keyboard:** `Alt+Letter` (e.g., Alt+F for File) or `F10` for first menu
- **Mouse:** Click on menu item
- **Arrow keys:** Navigate between menus when active
- **Esc:** Close menu, return to editing

### File Menu
```
┌─────────────────────────────┐
│ File                        │
├─────────────────────────────┤
│ New              Ctrl+N     │
│ Open...          Ctrl+O     │
│ Save             Ctrl+S     │
│ Save As...       Ctrl+Shift+S │
│ ─────────────────           │
│ Close            Ctrl+W     │
│ Quit             Ctrl+Q     │
└─────────────────────────────┘
```

### Edit Menu
```
┌─────────────────────────────┐
│ Edit                        │
├─────────────────────────────┤
│ Undo             Ctrl+Z     │
│ Redo             Ctrl+Y     │
│ ─────────────────           │
│ Cut              Ctrl+X     │
│ Copy             Ctrl+C     │
│ Paste            Ctrl+V     │
│ Select All       Ctrl+A     │
│ ─────────────────           │
│ Delete Line      Ctrl+Shift+K │
│ Duplicate Line   Ctrl+D     │
│ Move Line Up     Alt+Up     │
│ Move Line Down   Alt+Down   │
│ ─────────────────           │
│ Comment Toggle   Ctrl+/     │
└─────────────────────────────┘
```

### Search Menu
```
┌─────────────────────────────┐
│ Search                      │
├─────────────────────────────┤
│ Find...          Ctrl+F     │
│ Replace...       Ctrl+H     │
│ Find Next        F3         │
│ Find Previous    Shift+F3   │
│ ─────────────────           │
│ Go to Line...    Ctrl+G     │
└─────────────────────────────┘
```

### View Menu
```
┌─────────────────────────────┐
│ View                        │
├─────────────────────────────┤
│ ☑ Line Numbers   Ctrl+L     │
│ ☐ Whitespace     Ctrl+Shift+I │
│ ☐ Word Wrap      Ctrl+Shift+W │
│ ─────────────────           │
│ Zoom In          Ctrl++     │
│ Zoom Out         Ctrl+-     │
│ Reset Zoom       Ctrl+0     │
└─────────────────────────────┘
```

### Help Menu
```
┌─────────────────────────────┐
│ Help                        │
├─────────────────────────────┤
│ Keyboard Shortcuts          │
│ Documentation               │
│ ─────────────────           │
│ About ted                   │
└─────────────────────────────┘
```

---

## Mouse Support

**Phase 1+** - Basic mouse interaction

### Supported Mouse Actions

**Menu Bar:**
- **Click** - Activate menu (same as Alt+Key)
- **Hover** - Highlight menu item (visual feedback)

**Editing Area:**
- **Click** - Move cursor to click position
- **Double-click** - Select word under cursor
- **Triple-click** - Select entire line
- **Click and drag** - Create selection
- **Shift+Click** - Extend selection from cursor to click position

**Scrollbars (Future):**
- **Click on scrollbar** - Jump to that position
- **Drag scrollbar** - Smooth scrolling
- **Mouse wheel** - Scroll up/down (3 lines per tick)
- **Shift+Mouse wheel** - Horizontal scroll

**Info Bar Items (Future):**
- **Click encoding** - Change encoding dialog
- **Click line ending** - Toggle LF/CRLF
- **Click tab size** - Change tab size dialog

### Mouse Configuration
- **Enabled by default** - Can be disabled in config
- **Terminal support required** - Works with most modern terminals
- **Fallback** - Full keyboard-only operation always available

---

## Responsive Behavior

### Window Resize Handling
**Detection:** Terminal resize events via tcell

**Actions on Resize:**
1. Recalculate viewport dimensions
2. Adjust line wrapping (if word wrap enabled)
3. Reposition scrollbars
4. Keep cursor visible (scroll if needed)
5. Redraw screen efficiently (only changed regions)
6. Update menu bar and info bar widths

**Minimum Size:**
- **Width:** 40 columns (enough for basic editing)
- **Height:** 10 rows (menu + few lines + info bar)
- If smaller: Display "Terminal too small - resize to continue" message

### Adaptive Layout
**On narrow terminals (< 80 cols):**
- Abbreviate menu labels: "File" → "F", etc. (or hide some menus)
- Hide less important info bar items (file size, line ending)
- Shorten status indicators: "INSERT" → "INS"
- Reduce line number column width if possible

**On wide terminals (> 120 cols):**
- Show all info bar details
- More breathing room for menus
- Full status text: "INSERT" vs "INS"
- Could support split view (future)

**On short terminals (< 20 rows):**
- Reduce visible lines
- Ensure at least 5 lines of editing space
- Menus may overlay instead of push content

### Layout Priorities (when space is limited)
1. **Must have:** Menu bar, editing area, info bar
2. **Should have:** Line numbers, status indicators
3. **Nice to have:** Scrollbar indicators, extra info items

---

## Configuration

### Default Settings (Phase 0)
**Hard-coded defaults for initial release:**
- **Tab size:** 4 spaces
- **Use spaces (not tabs):** Yes (industry standard)
- **Line numbers:** Off (toggle with Ctrl+L)
- **Word wrap:** On (toggle with Ctrl+Shift+W)
- **Color scheme:** Dark mode
- **Auto-save:** Off
- **Show whitespace:** Off
- **Encoding:** UTF-8
- **Line ending:** Auto-detect (preserve file's original)

### Configuration File (Phase 5)
**Format:** TOML (`.tedrc` or `ted.toml`)
**Location:** `~/.tedrc` or `~/.config/ted/config.toml`

**Example config file:**
```toml
[editor]
tab_size = 4
use_spaces = true
line_numbers = false
word_wrap = true
auto_save = false
show_whitespace = false

[appearance]
theme = "dark"  # "dark" or "light" or custom theme name
font_size = 12

[search]
case_sensitive = false
whole_word = false
use_regex = false

[syntax]
# Enable/disable syntax highlighting per language
go = true
javascript = true
python = true
markdown = true

[filetypes]
# File extension to language mapping
".go" = "go"
".js" = "javascript"
".jsx" = "javascript"
".ts" = "typescript"
".py" = "python"
".md" = "markdown"
```

### Configurable Items (Phase 5+)
**YES - Will be configurable:**
- ✅ Color schemes (custom themes)
- ✅ File type associations (extension to language mapping)
- ✅ Syntax highlighting rules (per language)
- ✅ Tab size
- ✅ Spaces vs tabs
- ✅ Default view options (line numbers, word wrap, etc.)
- ✅ Search defaults (case sensitivity, regex)
- ✅ Auto-save behavior
- ✅ Theme selection

**NO - Will NOT be configurable:**
- ❌ Keyboard shortcuts (fixed for consistency across all machines)
- ❌ Menu structure
- ❌ Core UI layout

**Philosophy:** Shortcuts stay the same everywhere so users can work on any machine running ted without relearning or reconfiguring keybindings.

---

## Next Steps
*What should we design first?*
1. Keyboard shortcuts mapping
2. Core features list
3. Technical architecture/library choice
4. UI layout

## Performance Considerations

### Target Performance Metrics
- **Startup time:** < 100ms for small files (< 1MB)
- **Keystroke latency:** < 16ms (60fps feel)
- **Large file handling:** Support files up to 100MB
- **Memory usage:** < 100MB for typical editing sessions
- **Syntax highlighting:** Non-blocking, progressive

### Optimization Strategies

**Phase 0-2:** Don't optimize prematurely
- Get it working first
- Profile before optimizing

**Phase 3+:** Optimize hot paths
- Lazy syntax highlighting (only visible viewport)
- Incremental re-rendering (dirty regions only)
- Rope data structure for large buffers (if needed)
- Virtual scrolling for large files

**Large File Handling (Future):**
- Memory-mapped files
- Lazy loading (load only visible portions)
- Index-based navigation
- Read-only mode for very large files

---

## Documentation Plan

### User Documentation
- README.md - Project overview, installation
- USAGE.md - Quick start guide
- SHORTCUTS.md - Complete keyboard reference
- CONFIG.md - Configuration options (Phase 5+)

### Developer Documentation
- ARCHITECTURE.md - System design, module overview
- CONTRIBUTING.md - How to contribute
- API documentation (godoc comments)
- Each module README explaining its purpose

---

## Project Milestones

### v0.1 - "It Works" (Phase 0-1 Complete)
- Can edit and save files
- Basic clipboard and undo/redo
- Usable for simple text editing

### v0.5 - "Daily Driver" (Phase 0-3 Complete)
- All essential editing features
- Syntax highlighting
- Fast navigation
- Code editing features
- Can replace nano for most users

### v1.0 - "Production Ready" (Phase 0-5 Complete)
- Fully configurable
- Advanced search/replace
- Polished UX
- Comprehensive documentation
- Can replace most terminal editors

### v2.0 - "Power User" (Phase 6+)
- Multiple file support
- Advanced features
- Plugin system (maybe)
- LSP integration (maybe)

---

## Summary & Next Steps

### Design Complete ✅
- [x] Name: **ted** (Terminal EDitor)
- [x] Repository: github.com/AndrewDonelson/ted
- [x] License: MIT
- [x] Core features defined
- [x] Keyboard shortcuts mapped (all conflicts resolved)
- [x] UI layout designed (menu bar, edit area, info bar)
- [x] Modular architecture designed
- [x] 6-phase development plan
- [x] Technology choices made (tcell, clipboard library)
- [x] Testing strategy defined
- [x] Default settings finalized
- [x] Git workflow decided (main/develop/feature branches)
- [x] ALL design decisions finalized

### Ready to Start Phase 0 🚀

**Phase 0 Checklist:**
1. ✅ Repository already initialized at ~/Development/Golang/ted
2. Set up project structure (create module directories)
3. Initialize Go modules (`go mod init github.com/AndrewDonelson/ted`)
4. Install dependencies (tcell, clipboard library)
5. Implement `core/buffer` (basic buffer operations, cursor, insert/delete)
6. Implement `core/file` (read/write files)
7. Implement `ui/layout` (calculate screen regions, viewport)
8. Implement `ui/terminal` (tcell setup, event handling, resize detection)
9. Implement `ui/menu` (static menu bar display)
10. Implement `ui/renderer` (menu bar, edit area, **inverted info bar** rendering)
11. Implement `editor/editor` (wire all components, main loop)
12. Test: Open file, display with layout, edit, save, quit
13. Test: Resize terminal - verify responsive behavior
14. Test: Verify info bar has inverted text styling

**Estimated Time:** 5-7 days

**Success Criteria:**
- ✅ Opens file from command line
- ✅ Displays menu bar (File, Edit, Search, View, Help)
- ✅ Displays editable text in middle area
- ✅ Displays **inverted** info bar with filename, size, line count
- ✅ Arrow keys move cursor
- ✅ Can type and delete characters
- ✅ Ctrl+S saves file
- ✅ Ctrl+Q quits
- ✅ Terminal resize updates layout properly
- ✅ Info bar clearly stands out with inverted colors

### When You're Ready to Code
Just say the word and we can start implementing Phase 0! We'll:
1. Set up the project structure
2. Write the core buffer module
3. Get basic file I/O working
4. Create minimal terminal UI
5. Wire it all together
6. Test end-to-end

---

## Design Decisions ✅

All design decisions have been finalized:

### 1. Keyboard Shortcuts - RESOLVED
- **Delete line:** Ctrl+Shift+K
- **Toggle line numbers:** Ctrl+L
- **Duplicate line/Select word:** Ctrl+D
- **Close file:** Ctrl+W
- **Go to line:** Ctrl+G
- **Find next:** F3
- **NO CUSTOM KEYBINDINGS** - Shortcuts are fixed to ensure consistency across all machines and users

### 2. Repository - INITIALIZED
- **GitHub:** `github.com/AndrewDonelson/ted`
- **Local Path:** `~/Development/Golang/ted`
- **Status:** Repository already created and initialized
- **License:** MIT (permissive, simple)
- **Tools:** GitHub Speckit (specify CLI) for project management

### 3. Version Control - DECIDED
- **main** branch = stable releases
- **develop** branch = integration/development
- **Feature branches** for each phase (e.g., `feature/phase-0`, `feature/phase-1`)

### 4. Development Platform - DECIDED
- **Primary:** Ubuntu (Andrew's system)
- **Test on:** macOS, Windows

### 5. Menu Bar Behavior - DECIDED
- **Display mode:** Menus overlay content (don't push down)
- **Auto-hide:** No (menu bar always visible)
- **Activation:** Alt+Key AND F10 (both methods supported)
- **Layout:** 3 rows total (menu + edit area + info bar)

### 6. Color Scheme - DECIDED
- **Phase 0:** Dark mode only
- **Phase 5+:** Light mode option in configuration

### 7. Default Settings - DECIDED
- **Tab size:** 4 spaces
- **Spaces vs Tabs:** Spaces (industry standard for Go, JS, Python, etc.)
- **Line numbers:** Off by default (toggle with Ctrl+L)
- **Word wrap:** On by default (toggle with Ctrl+Shift+W)

---

*Last Updated: December 16, 2025*
*Status: **COMPLETE & READY FOR AI CODE AGENT IMPLEMENTATION***
*Repository: ~/Development/Golang/ted (initialized)*
*Tools: GitHub Speckit (specify CLI)*

---

## Quick Reference Summary

**Project:** ted - Terminal EDitor  
**Repository:** github.com/AndrewDonelson/ted  
**Local Path:** ~/Development/Golang/ted  
**License:** MIT  
**Language:** Go  
**Primary Platform:** Ubuntu → then macOS, Windows  
**Tools:** GitHub Speckit (specify CLI)

**Key Features:**
- Windows-style keyboard shortcuts (FIXED - no custom keybindings)
- Responsive UI (menu bar, edit area, **inverted info bar**)
- No modal editing (no vim modes)
- Cross-platform (Linux, macOS, Windows)
- Modular architecture
- 6-phase development plan

**Default Settings:**
- Tab size: 4 spaces
- Use spaces (not tabs)
- Line numbers: Off
- Word wrap: On
- Dark mode (info bar inverted)
- UTF-8 encoding

**UI Layout:**
- Top: Menu bar (File, Edit, Search, View, Help) + status (INS/OVR, UTF-8, LN/COL)
- Middle: Scrollable editing area with optional line numbers
- Bottom: **Inverted info bar** (light bg, dark text) - filename, size, type, settings

**Critical Shortcuts:**
- Save: Ctrl+S
- Quit: Ctrl+Q
- Undo/Redo: Ctrl+Z / Ctrl+Y
- Find: Ctrl+F
- Delete line: Ctrl+Shift+K
- Duplicate line: Ctrl+D
- Toggle line numbers: Ctrl+L
- Menu: Alt+Key or F10

**Configuration Philosophy:**
- ✅ Configurable: Color schemes, file associations, syntax highlighting, editor settings
- ❌ NOT Configurable: Keyboard shortcuts (consistency across all machines)


## Visual Mockup Examples

### Example 1: Editing a Go File

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ File  Edit  Search  View  Help               INS │ UTF-8 │ LN 8, COL 35     │
├─────────────────────────────────────────────────────────────────────────────┤
│  1 │ package main                                                            │
│  2 │                                                                         │
│  3 │ import (                                                                │
│  4 │     "fmt"                                                               │
│  5 │ )                                                                       │
│  6 │                                                                         │
│  7 │ func main() {                                                           │
│  8 │     fmt.Println("Hello, World!")█                                       │
│  9 │ }                                                                       │
│ 10 │                                                                         │
│ 11 │                                                                         │
│ 12 │                                                                         │
│ 13 │                                                                         │
│ 14 │                                                                         │
│    │                                                                         │
│    ↓ 45%                                                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│▓▓main.go▓│▓245 bytes▓│▓Go▓│▓Modified▓│▓Tab: 4▓│▓LF▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓│
└─────────────────────────────────────────────────────────────────────────────┘
```

### Example 2: Active Menu

```
┌─────────────────────────────────────────────────────────────────────────────┐
│[File] Edit  Search  View  Help               INS │ UTF-8 │ LN 1, COL 1      │
│ ┌─────────────────────────┐                                                 │
│ │ New           Ctrl+N    │                                                 │
│ │ Open...       Ctrl+O    │                                                 │
│ │ Save          Ctrl+S    │                                                 │
│ │ Save As...    Ctrl+Shift+S │                                              │
│ │ ───────────────────────  │                                                │
│ │ Close         Ctrl+W    │                                                 │
│ │ Quit          Ctrl+Q    │                                                 │
│ └─────────────────────────┘                                                 │
│  1 │ package main                                                            │
│  2 │                                                                         │
│  3 │ import (                                                                │
├─────────────────────────────────────────────────────────────────────────────┤
│▓▓main.go▓│▓245 bytes▓│▓Go▓│▓Saved▓│▓Tab: 4▓│▓LF▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓│
└─────────────────────────────────────────────────────────────────────────────┘
```

### Example 3: Search Active

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ File  Edit  Search  View  Help               INS │ UTF-8 │ LN 4, COL 9      │
├─────────────────────────────────────────────────────────────────────────────┤
│ Find: [fmt        ] Case: [x] Whole: [ ]                          [Close]   │
├─────────────────────────────────────────────────────────────────────────────┤
│  1 │ package main                                                            │
│  2 │                                                                         │
│  3 │ import (                                                                │
│  4 │     "⟪fmt⟫"                                                             │
│  5 │ )                                                                       │
│  6 │                                                                         │
│  7 │ func main() {                                                           │
│  8 │     ⟪fmt⟫.Println("Hello, World!")                                      │
│  9 │ }                                                                       │
│ 10 │                                                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│▓▓main.go▓│▓245 bytes▓│▓Go▓│▓Modified▓│▓Tab: 4▓│▓LF▓│▓Match 1 of 2▓▓▓▓▓▓▓▓▓▓│
└─────────────────────────────────────────────────────────────────────────────┘
```

### Example 4: Narrow Terminal (Adaptive)

```
┌────────────────────────────────────────┐
│ F E S V H       INS │ LN 8, COL 12    │
├────────────────────────────────────────┤
│  1 │ package main                      │
│  2 │                                   │
│  3 │ import (                          │
│  4 │     "fmt"                         │
│  5 │ )                                 │
│  6 │                                   │
│  7 │ func main() {                     │
│  8 │     fmt.Println█                  │
│  9 │ }                                 │
│    ↓                                   │
├────────────────────────────────────────┤
│▓▓main.go▓│▓245B▓│▓Modified▓▓▓▓▓▓▓▓▓▓▓▓│
└────────────────────────────────────────┘
```

### Example 5: With Selection

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ File  Edit  Search  View  Help               INS │ UTF-8 │ LN 4, COL 10     │
├─────────────────────────────────────────────────────────────────────────────┤
│  1 │ package main                                                            │
│  2 │                                                                         │
│  3 │ import (                                                                │
│  4 │     "█████"                                                             │
│  5 │ )                                                                       │
│  6 │                                                                         │
│  7 │ func main() {                                                           │
│  8 │     fmt.Println("Hello, World!")                                        │
│  9 │ }                                                                       │
│ 10 │                                                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│▓▓main.go▓│▓245 bytes▓│▓Go▓│▓Modified▓│▓Tab: 4▓│▓LF▓│▓5 chars selected▓▓▓▓▓▓│
└─────────────────────────────────────────────────────────────────────────────┘
```

*Note: █ = cursor, ⟪⟫ = search highlights, blocks = selection, ▓ = inverted text (info bar)*


---

## Implementation Readiness Checklist

### Pre-Development ✅
- [x] Project name decided: **ted**
- [x] Repository created: github.com/AndrewDonelson/ted
- [x] Local repository initialized: ~/Development/Golang/ted
- [x] License selected: MIT
- [x] Git workflow defined: main/develop/feature branches
- [x] All keyboard shortcuts finalized (NO custom keybindings)
- [x] UI layout fully designed (with inverted info bar)
- [x] Default settings confirmed
- [x] Technology stack chosen
- [x] **Tools selected:** GitHub Speckit (specify CLI) for project management

### Development Tools
**GitHub Speckit (specify CLI):**
- Project management and issue tracking
- Feature branch management
- Specification-driven development
- Integration with GitHub repository

### Phase 0 - Ready to Begin 🎯
**Goal:** Working editor with full UI layout in 5-7 days

**What we'll build:**
1. ✅ Repository structure (already initialized)
2. Go module initialization
3. Core buffer (text storage, cursor, basic editing)
4. File I/O (read, write, encoding detection)
5. Terminal UI (tcell integration, resize handling)
6. Layout system (menu bar, edit area, **inverted info bar**)
7. Basic rendering (static menu, text display, inverted info bar)
8. Main event loop (keyboard input, save, quit)

**What you'll be able to do:**
- Open a file: `ted myfile.txt`
- See it in a professional 3-row layout
- Info bar with clear inverted text styling
- Navigate with arrow keys
- Type and delete text
- Save with Ctrl+S
- Quit with Ctrl+Q
- Resize terminal and see it adapt

**Key Visual Features:**
- Menu bar across top
- Scrollable editing area in middle
- **Inverted info bar at bottom** (light background, dark text)

**Next Command When Ready:**
```bash
cd ~/Development/Golang/ted
# Ready to start Phase 0 implementation
```

---

*Design Document Complete - All Systems Go! 🚀*

**FOR AI CODE AGENT:**
This document contains complete specifications for implementing ted (Terminal EDitor).
All design decisions are finalized. No ambiguities remain. Ready for immediate implementation.

**CRITICAL REQUIREMENTS:**
1. Info bar MUST be inverted (light bg, dark text)
2. NO custom keybindings (shortcuts are fixed)
3. Follow modular architecture exactly as specified
4. Use GitHub Speckit (specify CLI) for project management

**START COMMAND:**
```bash
cd ~/Development/Golang/ted
# Repository is already initialized
# Begin Phase 0 implementation following "Implementation Order" section
```

**ESTIMATED COMPLETION:** 5-7 days for Phase 0
**SUCCESS CRITERIA:** All checkboxes in "Phase 0 Checklist" must be validated


---

## AI Code Agent Integration Guide

### Development Workflow with GitHub Speckit

**Repository Information:**
- **Remote:** github.com/AndrewDonelson/ted
- **Local:** ~/Development/Golang/ted
- **Status:** Initialized, ready for development
- **Branch Strategy:** main (stable) / develop (integration) / feature/* (development)

**Using specify CLI:**
```bash
# Create feature branch for Phase 0
specify create-branch feature/phase-0-foundation

# Track implementation progress
specify update-status phase-0 --progress 0%

# Mark completed tasks
specify complete-task "Initialize Go modules"
specify complete-task "Implement core buffer"

# Update documentation
specify update-docs --phase phase-0
```

### Implementation Priorities

**Phase 0 - Critical Path:**
1. **Project Structure** (Day 1)
   - Create all module directories
   - Initialize go.mod
   - Set up basic package structure

2. **Core Buffer** (Day 1-2)
   - Text storage (lines as []string initially)
   - Cursor position tracking
   - Basic insert/delete at cursor
   - No undo/redo yet (Phase 1)

3. **File I/O** (Day 2)
   - Read file into buffer
   - Write buffer to file
   - UTF-8 encoding handling
   - Error handling for file operations

4. **Terminal & Layout** (Day 3)
   - Initialize tcell
   - Calculate screen regions
   - Handle resize events
   - Define viewport structure

5. **Rendering** (Day 4-5)
   - Static menu bar rendering
   - Edit area text rendering
   - **Inverted info bar** rendering (critical visual element)
   - Cursor positioning

6. **Event Loop** (Day 5-6)
   - Keyboard input handling
   - Arrow key navigation
   - Character input
   - Ctrl+S (save) and Ctrl+Q (quit)

7. **Testing & Polish** (Day 6-7)
   - Test all basic functionality
   - Verify responsive resize
   - Check info bar inversion
   - Fix any bugs

### Key Design Constraints for AI Agent

**MUST FOLLOW:**
1. **No custom keybindings** - Shortcuts are hardcoded and identical everywhere
2. **Inverted info bar** - Bottom bar MUST use light bg (#d4d4d4) with dark text (#1e1e1e)
3. **Modular architecture** - Each module independent and testable
4. **Interfaces first** - Define interfaces before implementing
5. **No external config in Phase 0** - All settings hardcoded initially

**Code Style:**
- Follow Go conventions (gofmt, golint)
- Clear variable names
- Comments for exported functions
- Error handling always explicit
- No panics in production code

**Testing Requirements:**
- Unit tests for core/buffer
- Unit tests for core/file
- Integration tests for editor workflow
- Manual testing checklist before completion

### Module Implementation Order

```
Phase 0 Implementation Order:
1. core/buffer/buffer.go       - Text storage, cursor
2. core/file/reader.go          - File reading
3. core/file/writer.go          - File writing  
4. ui/layout/viewport.go        - Screen regions
5. ui/terminal/screen.go        - tcell setup
6. ui/terminal/resize.go        - Resize handling
7. ui/renderer/menubar.go       - Menu rendering
8. ui/renderer/text.go          - Text rendering
9. ui/renderer/infobar.go       - INVERTED info bar
10. ui/menu/menubar.go          - Static menu structure
11. editor/editor.go            - Main controller
12. main.go                     - Entry point
```

### Success Validation

**Before marking Phase 0 complete, verify:**
- ✅ Can run: `./ted testfile.txt`
- ✅ Menu bar displays correctly across top
- ✅ File contents display in middle
- ✅ Info bar is INVERTED (light bg, dark text) at bottom
- ✅ Arrow keys move cursor
- ✅ Can type characters
- ✅ Can delete with backspace
- ✅ Ctrl+S saves changes
- ✅ Ctrl+Q exits
- ✅ Terminal resize adapts layout
- ✅ No crashes or panics
- ✅ Code passes `go test ./...`
- ✅ Code passes `go vet ./...`

### Git Workflow for AI Agent

```bash
# Start Phase 0
git checkout -b feature/phase-0-foundation

# Commit frequently with clear messages
git commit -m "feat(buffer): implement basic text buffer"
git commit -m "feat(file): add file read/write operations"
git commit -m "feat(ui): implement inverted info bar rendering"

# Push to remote for review
git push origin feature/phase-0-foundation

# After review/testing, merge to develop
git checkout develop
git merge feature/phase-0-foundation

# Tag the completion
git tag -a v0.1-phase0 -m "Phase 0: Foundation complete"
git push --tags
```

### Reference Files for Implementation

**Look at these design sections when implementing:**
- **Buffer:** See "Module Interface Contracts" → buffer.Buffer
- **Rendering:** See "Color Scheme" for exact hex values
- **Info Bar:** See "Bottom Info Bar" for inverted styling specs
- **Shortcuts:** See "Keyboard Shortcuts" for all mappings
- **Layout:** See "Screen Layout Specification" for ASCII mockup

### Critical Implementation Notes

**Info Bar Rendering (CRITICAL):**
```go
// In ui/renderer/infobar.go
// MUST use inverted colors
bgColor := tcell.ColorWhite    // Light background #d4d4d4
fgColor := tcell.ColorBlack    // Dark text #1e1e1e
style := tcell.StyleDefault.Background(bgColor).Foreground(fgColor)
```

**Menu Bar vs Info Bar:**
- Menu bar: Normal colors (dark bg, light text)
- Info bar: INVERTED colors (light bg, dark text)
- This visual distinction is intentional and required

**File Operations:**
- Always use UTF-8 encoding in Phase 0
- Preserve line endings from original file
- Handle files that don't exist (create new)
- Handle read-only files (show error)

**Cursor Behavior:**
- Cursor can't move past end of line
- Cursor can't move past last line
- Cursor adjusts when line is shorter
- Cursor stays visible (scroll viewport if needed)

---

*AI Code Agent: This design document is complete and ready for implementation.*
*All decisions finalized. No ambiguities. Ready to code Phase 0.*