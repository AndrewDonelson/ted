# Expert Go Developer Implementation Brief - ted (Terminal EDitor)

## Project Overview

You are tasked with implementing **ted**, a modern cross-platform terminal text editor written in Go. This is a greenfield project with complete design specifications already finalized. Your implementation will follow a 6-phase development plan, starting with Phase 0 (Foundation).

**Repository:** github.com/AndrewDonelson/ted  
**Local Path:** ~/Development/Golang/ted (already initialized)  
**Complete Design Document:** See `terminal-editor-design.md`

## Your Role & Expectations

You are an **expert Go developer** who:
- Writes production-grade, idiomatic Go code following all community standards
- Designs clean, modular architectures with clear separation of concerns
- Implements comprehensive testing with edge case coverage approaching 100%
- Creates excellent documentation that serves both users and future maintainers
- Prioritizes code quality, performance, security, and maintainability
- Follows best practices from the Go community (Effective Go, Code Review Comments, etc.)

## Core Requirements

### 1. Code Quality Standards

**Idiomatic Go:**
- Follow all conventions from [Effective Go](https://golang.org/doc/effective_go.html)
- Adhere to [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` for all code formatting (no exceptions)
- Pass `go vet` with zero warnings
- Pass `golint` with zero issues
- Consider `golangci-lint` for comprehensive linting

**Code Organization:**
- Clear package structure matching the design document
- Single responsibility principle for all modules
- Dependency injection for testability
- Interfaces defined before implementations
- Exported types/functions have complete godoc comments
- Internal complexity hidden behind clean interfaces

**Naming Conventions:**
- Package names: lowercase, single word (e.g., `buffer`, `renderer`)
- Exported identifiers: PascalCase (e.g., `Buffer`, `MoveCursor`)
- Unexported identifiers: camelCase (e.g., `cursorPos`, `lineCount`)
- Constants: PascalCase for exported, camelCase for unexported
- Interface names: typically end in -er (e.g., `Renderer`, `Commander`)

### 2. Modular Architecture

**Package Structure (from design doc):**
```
ted/
├── main.go                     # Entry point, minimal logic
├── core/
│   ├── buffer/                # Text buffer - NO UI dependencies
│   │   ├── buffer.go          # Main buffer implementation
│   │   ├── buffer_test.go     # Comprehensive tests
│   │   ├── cursor.go          # Cursor logic
│   │   ├── cursor_test.go
│   │   ├── selection.go       # Selection handling
│   │   ├── selection_test.go
│   │   └── history.go         # Undo/redo (Phase 1)
│   ├── file/                  # File I/O - NO UI dependencies
│   │   ├── reader.go
│   │   ├── reader_test.go
│   │   ├── writer.go
│   │   └── writer_test.go
│   └── clipboard/             # System clipboard
│       ├── clipboard.go
│       └── clipboard_test.go
├── editor/
│   ├── editor.go              # Main controller, coordinates all modules
│   ├── editor_test.go
│   ├── commands.go            # Command dispatcher (Phase 1)
│   └── config.go              # Configuration (Phase 5)
├── ui/
│   ├── terminal/              # Terminal interface (tcell)
│   │   ├── screen.go
│   │   ├── screen_test.go
│   │   ├── input.go
│   │   ├── events.go
│   │   └── resize.go
│   ├── layout/                # Layout calculations
│   │   ├── viewport.go
│   │   ├── viewport_test.go
│   │   └── dimensions.go
│   ├── renderer/              # Rendering components
│   │   ├── renderer.go        # Main renderer interface
│   │   ├── menubar.go
│   │   ├── text.go
│   │   ├── infobar.go         # INVERTED styling (critical)
│   │   ├── linenumbers.go
│   │   └── scrollbar.go
│   ├── menu/                  # Menu system (Phase 1)
│   │   ├── menubar.go
│   │   ├── menu.go
│   │   └── menuitem.go
│   └── theme/                 # Color schemes (Phase 5)
│       └── theme.go
├── search/                    # Search/replace (Phase 2+)
│   ├── finder.go
│   └── replacer.go
├── syntax/                    # Syntax highlighting (Phase 3+)
│   ├── highlighter.go
│   └── languages/
└── utils/
    ├── keys.go                # Keyboard definitions
    └── platform.go            # OS-specific utilities
```

**Module Independence:**
- `core/*` packages MUST NOT depend on `ui/*` packages
- Each package should be independently testable
- Use dependency injection for cross-package dependencies
- Define clear interfaces at package boundaries
- Avoid circular dependencies (enforced by Go compiler)

### 3. Performance & Optimization

**Target Metrics (from design doc):**
- Startup time: < 100ms for files < 1MB
- Keystroke latency: < 16ms (60fps feel)
- Memory usage: < 100MB for typical sessions
- Support files up to 100MB

**Implementation Guidelines:**
- Profile before optimizing (use pprof)
- Efficient data structures (consider rope for large buffers in later phases)
- Minimize allocations in hot paths
- Use sync.Pool for frequently allocated objects
- Lazy initialization where appropriate
- Incremental rendering (only redraw dirty regions)

**Phase 0 Performance Notes:**
- Simple []string for line storage is acceptable
- Optimize in later phases if profiling shows issues
- Focus on correctness first, performance second

### 4. Security Considerations

**File Operations:**
- Validate all file paths (prevent directory traversal)
- Check file permissions before reading/writing
- Handle symlinks securely
- Limit file size to prevent memory exhaustion
- Sanitize file names for display

**Input Handling:**
- Validate all keyboard/mouse input
- Prevent buffer overflows
- Handle malformed UTF-8 gracefully
- Rate limit operations if needed
- No unsafe operations in production code

**Terminal Handling:**
- Properly restore terminal state on exit
- Handle terminal resize signals safely
- Prevent terminal injection attacks
- Clean up resources in defer statements

### 5. Error Handling

**Go Best Practices:**
- Return errors, don't panic (except in truly exceptional cases)
- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- Define custom error types for domain-specific errors
- Check all errors explicitly (no `_` ignoring errors)
- Log errors appropriately
- Graceful degradation where possible

**Example Custom Error:**
```go
// In core/file/errors.go
type FileError struct {
    Path string
    Op   string
    Err  error
}

func (e *FileError) Error() string {
    return fmt.Sprintf("file operation %s on %s: %v", e.Op, e.Path, e.Err)
}

func (e *FileError) Unwrap() error {
    return e.Err
}
```

### 6. Testing Requirements

**Coverage Goals:**
- `core/*` packages: 90%+ coverage (critical business logic)
- `editor/*` packages: 85%+ coverage
- `ui/*` packages: 70%+ coverage (harder to test, but still substantial)
- Overall project: 80%+ coverage

**Test Organization:**
```
buffer/
├── buffer.go
├── buffer_test.go          # Unit tests
├── buffer_benchmark_test.go # Benchmarks
└── testdata/               # Test fixtures
    ├── sample.txt
    └── large_file.txt
```

**Testing Strategies:**

**Unit Tests:**
```go
// core/buffer/buffer_test.go
func TestBuffer_Insert(t *testing.T) {
    tests := []struct {
        name    string
        initial []string
        pos     Position
        text    string
        want    []string
        wantErr bool
    }{
        {
            name:    "insert at beginning",
            initial: []string{"hello world"},
            pos:     Position{Line: 0, Col: 0},
            text:    "foo ",
            want:    []string{"foo hello world"},
            wantErr: false,
        },
        {
            name:    "insert at end",
            initial: []string{"hello"},
            pos:     Position{Line: 0, Col: 5},
            text:    " world",
            want:    []string{"hello world"},
            wantErr: false,
        },
        {
            name:    "insert in middle",
            initial: []string{"hello world"},
            pos:     Position{Line: 0, Col: 6},
            text:    "beautiful ",
            want:    []string{"hello beautiful world"},
            wantErr: false,
        },
        // Edge cases
        {
            name:    "insert newline",
            initial: []string{"hello world"},
            pos:     Position{Line: 0, Col: 5},
            text:    "\n",
            want:    []string{"hello", " world"},
            wantErr: false,
        },
        {
            name:    "insert at invalid position",
            initial: []string{"hello"},
            pos:     Position{Line: 1, Col: 0},
            text:    "foo",
            want:    nil,
            wantErr: true,
        },
        {
            name:    "insert empty string",
            initial: []string{"hello"},
            pos:     Position{Line: 0, Col: 0},
            text:    "",
            want:    []string{"hello"},
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            buf := NewBuffer()
            buf.lines = tt.initial // Assuming lines is a field
            
            err := buf.Insert(tt.pos, tt.text)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if !tt.wantErr && !reflect.DeepEqual(buf.lines, tt.want) {
                t.Errorf("Insert() got = %v, want %v", buf.lines, tt.want)
            }
        })
    }
}
```

**Table-Driven Tests:**
- Use for all functions with multiple scenarios
- Include normal cases, edge cases, and error cases
- Clear test names describing the scenario
- Use subtests with `t.Run()` for better output

**Edge Cases to Test:**
- Empty files
- Single-character files
- Very long lines (> 10,000 chars)
- Many lines (> 100,000 lines)
- Files with no newline at end
- Files with mixed line endings (CRLF, LF, CR)
- UTF-8 edge cases (multi-byte characters, emoji, RTL text)
- Cursor at boundary conditions
- Zero-width operations
- Concurrent operations (if applicable)
- Resource exhaustion scenarios

**Benchmarks:**
```go
// core/buffer/buffer_benchmark_test.go
func BenchmarkBuffer_Insert(b *testing.B) {
    buf := NewBuffer()
    pos := Position{Line: 0, Col: 0}
    text := "test"
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        buf.Insert(pos, text)
    }
}

func BenchmarkBuffer_InsertLargeLine(b *testing.B) {
    buf := NewBuffer()
    pos := Position{Line: 0, Col: 0}
    text := strings.Repeat("x", 10000)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        buf.Insert(pos, text)
    }
}
```

**Integration Tests:**
```go
// editor/editor_integration_test.go
func TestEditor_OpenEditSave(t *testing.T) {
    // Create temp file
    tmpfile, err := ioutil.TempFile("", "test*.txt")
    if err != nil {
        t.Fatal(err)
    }
    defer os.Remove(tmpfile.Name())
    
    content := "hello world\ntest file\n"
    if _, err := tmpfile.Write([]byte(content)); err != nil {
        t.Fatal(err)
    }
    tmpfile.Close()
    
    // Open file in editor
    editor := NewEditor()
    if err := editor.Open(tmpfile.Name()); err != nil {
        t.Fatalf("Open() error = %v", err)
    }
    
    // Make edits
    editor.Insert(Position{Line: 0, Col: 5}, " beautiful")
    
    // Save file
    if err := editor.Save(); err != nil {
        t.Fatalf("Save() error = %v", err)
    }
    
    // Verify saved content
    saved, err := ioutil.ReadFile(tmpfile.Name())
    if err != nil {
        t.Fatal(err)
    }
    
    expected := "hello beautiful world\ntest file\n"
    if string(saved) != expected {
        t.Errorf("Saved content = %q, want %q", string(saved), expected)
    }
}
```

**Test Helpers:**
```go
// core/buffer/testutil.go
func newTestBuffer(lines ...string) *Buffer {
    buf := NewBuffer()
    buf.lines = lines
    return buf
}

func assertBufferEqual(t *testing.T, got *Buffer, want []string) {
    t.Helper()
    if !reflect.DeepEqual(got.lines, want) {
        t.Errorf("buffer contents = %v, want %v", got.lines, want)
    }
}
```

### 7. Documentation Standards

**Package Documentation:**
```go
// Package buffer implements a text buffer for terminal text editing.
//
// The buffer stores text as a slice of lines and provides operations
// for inserting, deleting, and querying text. It maintains cursor
// position and supports text selection.
//
// Example usage:
//
//     buf := buffer.NewBuffer()
//     buf.Insert(buffer.Position{Line: 0, Col: 0}, "hello world")
//     text := buf.GetLine(0) // Returns "hello world"
//
package buffer
```

**Type Documentation:**
```go
// Buffer represents an in-memory text buffer.
// It stores text as a slice of lines and provides methods for
// editing operations. Buffer is not safe for concurrent use.
type Buffer struct {
    lines      []string
    cursor     Position
    selection  *Selection
    modified   bool
}

// Position represents a location in the buffer.
// Line and Col are zero-indexed.
type Position struct {
    Line int // Line number (0-indexed)
    Col  int // Column number (0-indexed, byte offset)
}
```

**Function Documentation:**
```go
// Insert inserts text at the specified position.
// If text contains newlines, it will be split across multiple lines.
// Returns an error if the position is invalid.
//
// Example:
//     err := buf.Insert(Position{Line: 0, Col: 5}, "world")
func (b *Buffer) Insert(pos Position, text string) error {
    // implementation
}
```

**README.md Structure:**
```markdown
# ted - Terminal EDitor

Modern, cross-platform terminal text editor with Windows-style shortcuts.

## Features
- Windows-style keyboard shortcuts (Ctrl+S, Ctrl+C, etc.)
- No modal editing (no vim modes)
- Responsive UI with menu bar
- Syntax highlighting
- Cross-platform (Linux, macOS, Windows)

## Installation
```bash
go install github.com/AndrewDonelson/ted@latest
```

## Usage
```bash
ted filename.txt
```

## Development
See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines.

## Architecture
See [ARCHITECTURE.md](docs/ARCHITECTURE.md) for system design.

## License
MIT License - see [LICENSE](LICENSE)
```

**CONTRIBUTING.md:**
- Development setup instructions
- Code style guidelines
- How to run tests
- How to submit PRs
- Issue reporting guidelines

**ARCHITECTURE.md:**
- High-level system design
- Module responsibilities
- Data flow diagrams
- Design decisions and rationale

### 8. Git Workflow

**Branch Strategy:**
- `main` - Stable releases only
- `develop` - Integration branch for development
- `feature/phase-N-description` - Feature branches for each phase

**Commit Messages (Conventional Commits):**
```
feat(buffer): implement cursor movement
fix(file): handle UTF-8 BOM correctly
test(buffer): add edge cases for insert operation
docs(readme): update installation instructions
refactor(renderer): extract info bar to separate file
perf(buffer): optimize large file handling
```

**Format:**
```
<type>(<scope>): <subject>

<body>

<footer>
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`

**Example Commit:**
```
feat(buffer): add undo/redo support

Implements history tracking for buffer operations using
a command pattern. Each operation is stored as a command
that can be undone and redone.

- Add History struct to track commands
- Implement Undo() and Redo() methods
- Add max history depth configuration
- Update Buffer to record operations

Closes #15
```

### 9. Dependencies

**Approved Dependencies:**
- `github.com/gdamore/tcell/v2` - Terminal handling (REQUIRED)
- `github.com/atotto/clipboard` - Clipboard operations (REQUIRED)
- Standard library packages (unlimited use)

**Testing Dependencies:**
- `github.com/stretchr/testify` - Optional, for assert/require helpers
- Standard library `testing` package (preferred)

**Dependency Management:**
- Use Go modules (`go.mod`, `go.sum`)
- Vendor dependencies if needed for reproducibility
- Keep dependencies minimal
- Audit security with `go list -m all | nancy sleuth`

### 10. Critical Implementation Details

**CRITICAL: Inverted Info Bar**
The bottom info bar MUST use inverted colors:
```go
// ui/renderer/infobar.go
const (
    // Info bar uses inverted colors for visual distinction
    infoBarBg = tcell.Color252  // Light gray #d4d4d4
    infoBarFg = tcell.Color235  // Dark gray #1e1e1e
)

func (r *Renderer) RenderInfoBar(info FileInfo) error {
    style := tcell.StyleDefault.
        Background(infoBarBg).
        Foreground(infoBarFg)
    
    // Render info bar with inverted style
    // ...
}
```

**CRITICAL: No Custom Keybindings**
Keyboard shortcuts are FIXED and HARDCODED:
```go
// utils/keys.go
const (
    KeySave        = tcell.KeyCtrlS
    KeyQuit        = tcell.KeyCtrlQ
    KeyUndo        = tcell.KeyCtrlZ
    KeyRedo        = tcell.KeyCtrlY
    // etc. - ALL shortcuts defined here
)

// NO configuration file for keybindings
// NO user-customizable shortcuts
// Philosophy: Consistent experience across all machines
```

**CRITICAL: UTF-8 Handling**
```go
// core/buffer/buffer.go
import "unicode/utf8"

// Position represents a position in the buffer
type Position struct {
    Line int // Line number (0-indexed)
    Col  int // BYTE offset, not rune offset
}

// Helper to convert byte offset to rune offset
func (b *Buffer) byteToRune(line int, byteOffset int) int {
    if line >= len(b.lines) {
        return 0
    }
    return utf8.RuneCountInString(b.lines[line][:byteOffset])
}
```

**CRITICAL: Error Handling Example**
```go
// core/file/reader.go
func ReadFile(path string) ([]string, error) {
    // Validate path
    cleanPath := filepath.Clean(path)
    if !filepath.IsAbs(cleanPath) {
        var err error
        cleanPath, err = filepath.Abs(cleanPath)
        if err != nil {
            return nil, fmt.Errorf("resolve path %q: %w", path, err)
        }
    }
    
    // Check if file exists and is readable
    info, err := os.Stat(cleanPath)
    if err != nil {
        return nil, fmt.Errorf("stat file %q: %w", cleanPath, err)
    }
    
    if info.IsDir() {
        return nil, fmt.Errorf("path %q is a directory", cleanPath)
    }
    
    // Read file
    data, err := os.ReadFile(cleanPath)
    if err != nil {
        return nil, fmt.Errorf("read file %q: %w", cleanPath, err)
    }
    
    // Split into lines (handle different line endings)
    lines := splitLines(string(data))
    return lines, nil
}
```

### 11. Phase 0 Implementation Checklist

Your immediate task is **Phase 0: Foundation**. Here's the implementation order:

**Day 1: Project Setup & Core Buffer**
- [ ] Initialize `go.mod` if needed
- [ ] Create package directory structure
- [ ] Implement `core/buffer/buffer.go`
  - [ ] Basic Buffer struct
  - [ ] NewBuffer() constructor
  - [ ] Insert() method
  - [ ] Delete() method
  - [ ] GetLine() method
  - [ ] LineCount() method
- [ ] Implement `core/buffer/cursor.go`
  - [ ] Cursor position tracking
  - [ ] MoveCursor() method
  - [ ] Boundary validation
- [ ] Write comprehensive tests for buffer
  - [ ] Normal operations
  - [ ] Edge cases (empty, boundaries, etc.)
  - [ ] UTF-8 handling
  - [ ] Target: 90%+ coverage

**Day 2: File I/O**
- [ ] Implement `core/file/reader.go`
  - [ ] ReadFile() function
  - [ ] UTF-8 decoding
  - [ ] Line ending detection
  - [ ] Error handling
- [ ] Implement `core/file/writer.go`
  - [ ] WriteFile() function
  - [ ] Preserve line endings
  - [ ] Atomic write (temp file + rename)
  - [ ] Error handling
- [ ] Write tests for file operations
  - [ ] Read/write round-trip
  - [ ] Different line endings
  - [ ] Permission errors
  - [ ] Large files
  - [ ] Target: 85%+ coverage

**Day 3: Terminal & Layout**
- [ ] Implement `ui/terminal/screen.go`
  - [ ] Initialize tcell
  - [ ] Screen struct
  - [ ] Clear() method
  - [ ] Refresh() method
  - [ ] Cleanup/restore terminal state
- [ ] Implement `ui/terminal/resize.go`
  - [ ] Handle resize events
  - [ ] Update dimensions
- [ ] Implement `ui/layout/viewport.go`
  - [ ] Viewport struct
  - [ ] Calculate screen regions
  - [ ] GetMenuBarRegion()
  - [ ] GetEditAreaRegion()
  - [ ] GetInfoBarRegion()
- [ ] Write tests for layout calculations
  - [ ] Different terminal sizes
  - [ ] Minimum size handling
  - [ ] Target: 80%+ coverage

**Day 4: Rendering**
- [ ] Implement `ui/renderer/renderer.go`
  - [ ] Renderer interface
  - [ ] Basic implementation
- [ ] Implement `ui/renderer/menubar.go`
  - [ ] Render static menu bar
  - [ ] "File Edit Search View Help"
  - [ ] Status indicators on right
- [ ] Implement `ui/renderer/text.go`
  - [ ] Render buffer text
  - [ ] Handle scrolling
  - [ ] Cursor positioning
- [ ] Implement `ui/renderer/infobar.go` **CRITICAL**
  - [ ] INVERTED colors (light bg, dark text)
  - [ ] Render file info
  - [ ] Render editor state
- [ ] Implement `ui/menu/menubar.go`
  - [ ] Static menu structure
  - [ ] Menu items (no interaction yet)

**Day 5: Event Loop**
- [ ] Implement `ui/terminal/input.go`
  - [ ] Keyboard event handling
  - [ ] Map tcell events to actions
- [ ] Implement `editor/editor.go`
  - [ ] Editor struct
  - [ ] Initialize all components
  - [ ] Main event loop
  - [ ] Handle arrow keys
  - [ ] Handle character input
  - [ ] Handle Ctrl+S (save)
  - [ ] Handle Ctrl+Q (quit)
- [ ] Implement `main.go`
  - [ ] Parse command-line args
  - [ ] Create editor
  - [ ] Run event loop
  - [ ] Handle errors
  - [ ] Clean exit

**Day 6-7: Testing & Polish**
- [ ] Integration tests
  - [ ] Open file → edit → save workflow
  - [ ] Terminal resize handling
  - [ ] Keyboard input handling
- [ ] Manual testing
  - [ ] Open various file types
  - [ ] Test all supported keys
  - [ ] Resize terminal
  - [ ] Verify info bar inversion
- [ ] Documentation
  - [ ] Package godocs
  - [ ] README.md
  - [ ] Usage examples
- [ ] Code review
  - [ ] Run `go vet`
  - [ ] Run `golint`
  - [ ] Run `gofmt`
  - [ ] Check test coverage
- [ ] Performance check
  - [ ] Run benchmarks
  - [ ] Profile if needed
  - [ ] Verify startup time < 100ms

### 12. Success Criteria for Phase 0

Before marking Phase 0 complete, ALL of these must be true:

**Functionality:**
- ✅ `./ted filename.txt` opens the file
- ✅ File contents display correctly
- ✅ Menu bar shows: "File Edit Search View Help"
- ✅ Status shows: mode, encoding, line/col
- ✅ Info bar shows: filename, size, type, settings
- ✅ Info bar has INVERTED colors (visually distinct)
- ✅ Arrow keys move cursor correctly
- ✅ Character input inserts at cursor
- ✅ Backspace deletes character before cursor
- ✅ Ctrl+S saves the file
- ✅ Ctrl+Q exits the program
- ✅ Terminal resize updates layout

**Code Quality:**
- ✅ All packages have clear, single responsibilities
- ✅ All exported functions have godoc comments
- ✅ No `go vet` warnings
- ✅ No `golint` issues
- ✅ Code is `gofmt`-formatted
- ✅ No TODO comments in production code
- ✅ Error handling is explicit and comprehensive

**Testing:**
- ✅ `go test ./...` passes with no failures
- ✅ Test coverage ≥ 80% overall
- ✅ `core/buffer` coverage ≥ 90%
- ✅ `core/file` coverage ≥ 85%
- ✅ All edge cases tested
- ✅ Integration tests pass
- ✅ No race conditions (`go test -race ./...`)

**Documentation:**
- ✅ README.md is complete
- ✅ All packages have package documentation
- ✅ ARCHITECTURE.md exists (optional for Phase 0)
- ✅ Code examples in documentation work
- ✅ Installation instructions are clear

**Performance:**
- ✅ Startup time < 100ms for small files
- ✅ No noticeable lag when typing
- ✅ Memory usage is reasonable
- ✅ No memory leaks (test with long-running sessions)

## Deliverables

At the end of Phase 0, you will deliver:

1. **Source Code**
   - Complete, working implementation
   - All files properly organized
   - All tests passing
   - All documentation complete

2. **Test Suite**
   - Unit tests for all packages
   - Integration tests for workflows
   - Benchmark tests for performance
   - Test coverage report (≥ 80%)

3. **Documentation**
   - README.md with usage instructions
   - Godoc comments for all exported identifiers
   - Examples in documentation
   - Optional: ARCHITECTURE.md

4. **Build Artifacts**
   - Clean `go build` with no warnings
   - All linters passing
   - Version tag: `v0.1-phase0`

## Resources

**Reference Documents:**
- `terminal-editor-design.md` - Complete design specification
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [tcell documentation](https://pkg.go.dev/github.com/gdamore/tcell/v2)

**Testing Resources:**
- [Testing in Go](https://golang.org/pkg/testing/)
- [Table-driven tests](https://github.com/golang/go/wiki/TableDrivenTests)
- [Testify documentation](https://pkg.go.dev/github.com/stretchr/testify)

**Git Workflow:**
```bash
cd ~/Development/Golang/ted
git checkout -b feature/phase-0-foundation
# Implement Phase 0
git commit -m "feat(core): implement buffer module"
# ... more commits ...
git push origin feature/phase-0-foundation
# Create PR to develop branch
```

## Questions & Clarifications

If you encounter any ambiguity or need clarification:
1. Check `terminal-editor-design.md` first (it's comprehensive)
2. Follow Go community standards when design doc is silent
3. Make reasonable decisions and document them
4. Prioritize code quality over speed

## Final Notes

You are building the foundation of a professional terminal editor. The quality of Phase 0 will determine the success of all subsequent phases. Take the time to:

- Write clean, idiomatic Go code
- Test thoroughly with comprehensive edge cases
- Document clearly for future maintainers
- Design modular interfaces for extensibility
- Follow all community standards

Your work will be used by developers who expect a reliable, fast, well-designed tool. Make it excellent.

**Good luck, and happy coding!** 🚀