# Implementation Phases: ted - Terminal EDitor

**Project Type**: Terminal Text Editor
**Language**: Go 1.25.5
**Current State**: Phase 0-1 Complete (Foundation + Essential Editing)
**Target**: Production v1.0 (Phases 2-5)
**Estimated Total**: 25-30 hours

---

## Phase 1: Line Operations & Navigation (Complete Phase 1)
**Type**: Core Editing
**Estimated**: 3-4 hours
**Files**: 
- `core/buffer/operations.go` (new)
- `editor/editor.go` (add handlers)
- `ui/terminal/input.go` (add key bindings)
- `core/buffer/operations_test.go` (new)

**Tasks**:
- [ ] Implement delete line (Ctrl+Shift+K) in buffer
- [ ] Implement duplicate line (Ctrl+D) in buffer
- [ ] Implement move line up/down (Alt+Up/Down) in buffer
- [ ] Implement insert line above/below (Ctrl+Enter, Ctrl+Shift+Enter)
- [ ] Add word navigation (Ctrl+Left/Right) with word boundary detection
- [ ] Implement Page Up/Down navigation
- [ ] Add key bindings in input handler
- [ ] Write comprehensive tests for all operations
- [ ] Update menu items to show these shortcuts

**Verification Criteria**:
- [ ] Delete line removes current line and joins with next
- [ ] Duplicate line creates exact copy below current
- [ ] Move line swaps with adjacent line
- [ ] Insert line creates empty line without moving cursor
- [ ] Word navigation jumps to next/prev word boundary
- [ ] Page Up/Down moves by viewport height
- [ ] All operations maintain cursor position appropriately
- [ ] Tests pass with >80% coverage

**Exit Criteria**: All line operations and navigation work correctly with tests

---

## Phase 2: Dialog System Framework
**Type**: UI Infrastructure
**Estimated**: 4-5 hours
**Files**:
- `ui/dialog/dialog.go` (new - base dialog interface)
- `ui/dialog/input.go` (new - text input dialog)
- `ui/dialog/confirm.go` (new - confirmation dialog)
- `ui/renderer/dialog.go` (new - dialog rendering)
- `editor/dialogs.go` (new - dialog implementations)

**Tasks**:
- [x] Create Dialog interface (Show, Hide, HandleInput, Render)
- [x] Implement InputDialog with text field and buttons
- [x] Implement ConfirmDialog for yes/no prompts
- [x] Add dialog rendering to renderer (centered, modal overlay)
- [x] Create dialog manager in editor (stack-based for nested dialogs)
- [x] Implement "Go to Line" dialog (Ctrl+G)
- [x] Implement "Open File" dialog (Ctrl+O)
- [x] Implement "Save As" dialog (Ctrl+Shift+S)
- [x] Add dialog key bindings (Tab to switch buttons, Enter to confirm, Esc to cancel)

**Verification Criteria**:
- [x] Dialog appears centered on screen with border
- [x] Text input captures keyboard input correctly
- [x] Tab cycles between input field and buttons
- [x] Enter confirms, Esc cancels
- [x] Go to Line accepts line number and jumps there
- [x] Open File shows path input and opens file
- [x] Save As shows path input and saves file
- [x] Dialog closes properly returning to editor

**Exit Criteria**: ✅ Dialog framework works, three basic dialogs functional

---

## Phase 3: Search & Replace System
**Type**: Search Module
**Estimated**: 5-6 hours
**Files**:
- `search/finder.go` (new - search engine)
- `search/replacer.go` (new - replace functionality)
- `ui/dialog/search.go` (new - search/replace dialog)
- `editor/search.go` (new - search integration)
- `search/finder_test.go` (new)
- `search/replacer_test.go` (new)

**Tasks**:
- [ ] Implement Finder with regex support (using standard regexp package)
- [ ] Support case-sensitive/insensitive search toggle
- [ ] Support whole word search toggle
- [ ] Implement find next/previous (F3/Shift+F3)
- [ ] Implement Replace functionality
- [ ] Implement Replace All with count
- [ ] Create SearchDialog UI with options checkboxes
- [ ] Highlight search matches in text (distinct background color)
- [ ] Show match count in info bar ("Match 3 of 12")
- [ ] Add search history (last 20 searches)
- [ ] Write comprehensive tests for finder and replacer

**Verification Criteria**:
- [ ] Ctrl+F opens search dialog
- [ ] Type query, press Enter finds first match
- [ ] F3 finds next, Shift+F3 finds previous
- [ ] Match highlighting visible in text
- [ ] Case sensitivity toggle works
- [ ] Whole word toggle works
- [ ] Replace replaces current match
- [ ] Replace All replaces all with count displayed
- [ ] No matches shows "No results" message
- [ ] Tests cover regex, case sensitivity, whole word

**Exit Criteria**: Full search/replace with UI working

---

## Phase 4: Display Features & Polish
**Type**: UI Enhancement
**Estimated**: 4-5 hours
**Files**:
- `ui/layout/viewport.go` (modify - add line numbers)
- `ui/renderer/text.go` (modify - line number rendering)
- `ui/renderer/linenumbers.go` (new)
- `editor/display.go` (new - display settings)
- `ui/dialog/unsaved.go` (new - unsaved changes dialog)

**Tasks**:
- [ ] Implement line numbers toggle (Ctrl+L)
- [ ] Render line numbers with dynamic width based on total lines
- [ ] Highlight current line number
- [ ] Implement word wrap toggle (Ctrl+Shift+W)
- [ ] Add soft wrap rendering in text area
- [ ] Show unsaved changes indicator (*) in info bar
- [ ] Implement unsaved changes prompt on exit (Ctrl+Q with modifications)
- [ ] Add horizontal scrollbar indicator when text exceeds width
- [ ] Implement show whitespace toggle (Ctrl+Shift+I) - render tabs/spaces visibly
- [ ] Add help dialog (F1 or Help menu) showing keyboard shortcuts
- [ ] Add about dialog showing version and credits

**Verification Criteria**:
- [ ] Ctrl+L toggles line numbers on/off
- [ ] Line numbers right-aligned with padding
- [ ] Current line number highlighted
- [ ] Word wrap toggle works (soft wrap at viewport edge)
- [ ] Unsaved changes show "*" in title bar/info bar
- [ ] Quitting with unsaved changes shows confirmation dialog
- [ ] Whitespace characters visible when toggled (» for tab, · for space)
- [ ] Help dialog lists all keyboard shortcuts
- [ ] About dialog shows version from build flags

**Exit Criteria**: All display features working, polished UI

---

## Phase 5: Syntax Highlighting Framework
**Type**: Language Support
**Estimated**: 5-6 hours
**Files**:
- `syntax/highlighter.go` (new - highlighting engine)
- `syntax/tokenizer.go` (new - token parser)
- `syntax/languages/go.go` (new - Go language rules)
- `syntax/languages/javascript.go` (new - JS language rules)
- `syntax/languages/python.go` (new - Python language rules)
- `ui/renderer/highlight.go` (new - apply colors)
- `ui/theme/theme.go` (new - color scheme management)

**Tasks**:
- [ ] Define Token types (Keyword, String, Comment, Number, etc.)
- [ ] Implement Tokenizer interface for different languages
- [ ] Create Go tokenizer (keywords, strings, comments, numbers)
- [ ] Create JavaScript tokenizer
- [ ] Create Python tokenizer
- [ ] Implement Highlighter that tokenizes visible lines only (lazy highlighting)
- [ ] Create Theme with color definitions
- [ ] Apply syntax colors in renderer
- [ ] Auto-detect language from file extension
- [ ] Add language indicator to info bar ("Go", "JavaScript", etc.)
- [ ] Tests for tokenizers

**Verification Criteria**:
- [ ] Go files show syntax colors (keywords blue, strings orange, comments green)
- [ ] JavaScript files highlighted correctly
- [ ] Python files highlighted correctly
- [ ] Language auto-detected from extension
- [ ] Info bar shows detected language
- [ ] Highlighting only applies to visible viewport (performance)
- [ ] No highlighting for unknown file types (plain text)
- [ ] Tests verify correct tokenization

**Exit Criteria**: Syntax highlighting works for Go/JS/Python

---

## Phase 6: Configuration System
**Type**: Configuration
**Estimated**: 3-4 hours
**Files**:
- `editor/config.go` (new - configuration management)
- `config/loader.go` (new - TOML parsing)
- `ui/theme/themes.go` (new - built-in themes)
- `ted.example.toml` (new - example config file)

**Tasks**:
- [ ] Define Config struct (tab size, spaces/tabs, line numbers, word wrap, theme)
- [ ] Implement TOML config loader using BurntSushi/toml
- [ ] Support ~/.tedrc and ~/.config/ted/config.toml locations
- [ ] Create default configuration
- [ ] Implement theme switching (dark/light)
- [ ] Add tab size configuration (2, 4, 8)
- [ ] Add spaces vs tabs configuration
- [ ] Add default view options (line numbers, word wrap)
- [ ] Create example config file with all options documented
- [ ] Apply config on editor startup
- [ ] Watch config file for changes (optional)

**Verification Criteria**:
- [ ] Config loads from ~/.tedrc
- [ ] Config loads from ~/.config/ted/config.toml (preferred)
- [ ] Tab size setting applies (2, 4, or 8 spaces)
- [ ] Spaces vs tabs setting works
- [ ] Default view options apply on startup
- [ ] Theme switching works (dark/light)
- [ ] Invalid config shows error message
- [ ] Missing config uses defaults

**Exit Criteria**: Configurable editor with TOML config file

---

## Phase 7: Code Editing Features
**Type**: Code Features
**Estimated**: 3-4 hours
**Files**:
- `core/buffer/indent.go` (new - indentation operations)
- `core/buffer/brackets.go` (new - bracket matching)
- `syntax/languages/comments.go` (new - comment detection)
- `editor/code.go` (new - code editing integration)

**Tasks**:
- [ ] Implement auto-indentation (maintain indent level on Enter)
- [ ] Implement Tab key for indent (spaces or tab based on config)
- [ ] Implement Shift+Tab for unindent
- [ ] Implement Ctrl+/ for toggle line comment (language-aware)
- [ ] Implement jump to matching bracket (Ctrl+B)
- [ ] Highlight matching bracket when cursor is on bracket
- [ ] Add language-specific comment markers (//, #, /* */)
- [ ] Test auto-indent with various scenarios
- [ ] Update menu items with code editing shortcuts

**Verification Criteria**:
- [ ] Enter maintains current indentation level
- [ ] Tab inserts configured indentation (spaces or tab)
- [ ] Shift+Tab removes one level of indentation
- [ ] Ctrl+/ comments current line (adds //, #, etc.)
- [ ] Ctrl+/ on commented line uncomments it
- [ ] Works for selection (comment/uncomment all selected lines)
- [ ] Ctrl+B jumps to matching bracket
- [ ] Matching bracket highlighted when cursor on bracket
- [ ] Works for (), {}, []

**Exit Criteria**: Code editing features complete

---

## Phase 8: Final Testing & Polish
**Type**: Testing & Polish
**Estimated**: 3-4 hours
**Files**:
- All test files (add missing tests)
- `README.md` (update)
- `CHANGELOG.md` (new)

**Tasks**:
- [ ] Run full test suite: `make test`
- [ ] Run race detection: `make test-race`
- [ ] Achieve >80% test coverage overall
- [ ] Add integration tests for key workflows
- [ ] Test on different terminal sizes
- [ ] Test with various file types
- [ ] Fix any remaining TODO comments
- [ ] Update README with all implemented features
- [ ] Create CHANGELOG.md for v1.0 release
- [ ] Final code review and cleanup
- [ ] Build release binaries: `make cross-build`
- [ ] Test release binaries on Linux/macOS/Windows

**Verification Criteria**:
- [ ] All tests pass
- [ ] No race conditions detected
- [ ] Overall coverage >80%
- [ ] Manual testing: open, edit, save workflow
- [ ] Manual testing: search/replace workflow
- [ ] Manual testing: line operations
- [ ] Manual testing: dialogs (open, save as, go to line)
- [ ] Manual testing: syntax highlighting
- [ ] Manual testing: configuration file
- [ ] Cross-platform binaries build successfully

**Exit Criteria**: Production-ready v1.0 release

---

## Implementation Order

1. **Phase 1** → Line Operations (foundational editing)
2. **Phase 2** → Dialog System (enables UI interactions)
3. **Phase 3** → Search/Replace (major feature)
4. **Phase 4** → Display Features (polish)
5. **Phase 5** → Syntax Highlighting (visual enhancement)
6. **Phase 6** → Configuration (customization)
7. **Phase 7** → Code Editing (productivity)
8. **Phase 8** → Final Testing (production ready)

---

## Dependencies to Add

```go
// For TOML configuration
require github.com/BurntSushi/toml v1.3.2
```

---

## Key Design Decisions

### Dialog System
- Stack-based dialog manager allows nested dialogs (e.g., unsaved changes inside another operation)
- Modal overlay blocks interaction with editor until dismissed
- Consistent styling: bordered box, centered, clear buttons

### Search/Replace
- Regex support using standard `regexp` package
- Case sensitivity and whole word as toggles in UI
- Match highlighting with distinct background color
- Search history stored in memory only (not persisted)

### Syntax Highlighting
- Tokenizer pattern: each language implements Tokenizer interface
- Lazy highlighting: only tokenize visible viewport lines
- Simple regex-based tokenizers (not full parsers - good enough for display)
- Color themes defined in Theme struct

### Configuration
- TOML format (human-friendly, well-supported in Go)
- XDG standard location preferred (~/.config/ted/config.toml)
- Traditional ~/.tedrc also supported
- Watch for changes and reload (optional enhancement)

---

## Risk Mitigation

**High Risk: Search/Replace Performance**
- Mitigation: Only search visible lines for highlighting, full search on Find Next
- Lazy evaluation of regex matches

**High Risk: Syntax Highlighting Speed**
- Mitigation: Only tokenize visible viewport
- Cache tokenization results for unchanged lines
- Simple regex patterns (not complex parsers)

**Medium Risk: Dialog Input Handling**
- Mitigation: Thoroughly test all dialog key bindings
- Clear separation between dialog and editor input modes

---

## Success Criteria (Production v1.0)

- [x] All features from README implemented
- [x] All TODO comments resolved
- [x] Test coverage >80%
- [x] No race conditions
- [x] Cross-platform builds working
- [x] Documentation complete (README, CHANGELOG)
- [x] All keyboard shortcuts from design document work
- [x] Configuration file support working
- [x] Syntax highlighting for major languages
- [x] Search/replace with regex support
- [x] Dialog system for all file/search operations
- [x] Production-ready stability

---

*Last Updated: 2026-02-03*
*Current Status: Ready to begin Phase 1*
