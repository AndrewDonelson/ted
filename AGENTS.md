# AGENTS.md - Development Guide for ted

> **ted** is a modern, cross-platform terminal text editor written in Go.
> Module: `github.com/AndrewDonelson/ted`

## Build Commands

```bash
# Build the binary (outputs to bin/ted)
make build

# Build and run
go run .

# Install to $GOPATH/bin
make install

# Clean build artifacts
make clean
```

## Test Commands

```bash
# Run all tests (clean output via custom script)
make test

# Run single test package
go test ./core/buffer -v

# Run specific test function
go test ./core/buffer -run TestBuffer_Insert -v

# Verbose test output
make test-verbose

# Tests with coverage report (outputs coverage.html)
make test-coverage

# Race condition detection
make test-race

# Benchmark tests
make test-bench
```

## Lint/Format Commands

```bash
# Run all linters (requires golangci-lint)
make lint

# Run go vet
make vet

# Format code with gofmt
make fmt

# Check if code is formatted (fails if not)
make fmt-check

# Run all checks (format, vet, tests)
make check

# Full pipeline: clean, format, vet, test, build
make all
```

## Code Style Guidelines

### Imports
- Group imports: stdlib first, then external packages, then internal packages
- Separate groups with blank lines
- Use goimports or gofmt for formatting

```go
import (
    "fmt"
    "os"

    "github.com/gdamore/tcell/v2"

    "github.com/AndrewDonelson/ted/core/buffer"
)
```

### Formatting
- Use `gofmt` with `-s` flag for simplification
- 4-space indentation (via gofmt)
- No trailing whitespace
- End files with a newline

### Naming Conventions
- **Packages**: lowercase, no underscores (e.g., `buffer`, `file`, `renderer`)
- **Exported**: PascalCase (e.g., `NewBuffer`, `Insert`)
- **Unexported**: camelCase (e.g., `validatePosition`, `updateCursor`)
- **Types**: PascalCase, descriptive (e.g., `LineEnding`, `EditorMode`)
- **Constants**: PascalCase for exported (e.g., `LineEndingLF`)
- **Test files**: `*_test.go`, test functions: `Test<Name>` or `Test<Struct>_<Method>`

### Documentation
- Every package must have a package comment at the top
- Exported types, functions, and constants must have documentation comments
- Use complete sentences starting with the name being documented
- Include examples in doc comments where helpful

```go
// Package buffer implements a text buffer for terminal text editing.
//
// The buffer stores text as a slice of lines and provides operations
// for inserting, deleting, and querying text.
package buffer

// Insert inserts text at the specified position.
// If text contains newlines, it will be split across multiple lines.
// Returns an error if the position is invalid.
func (b *Buffer) Insert(pos Position, text string) error
```

### Error Handling
- Always check errors and handle them appropriately
- Wrap errors with context using `fmt.Errorf("...: %w", err)`
- Don't ignore errors with `_` unless explicitly documented why
- Return errors rather than logging them in library code

```go
if err := b.validatePosition(pos); err != nil {
    return fmt.Errorf("validate position: %w", err)
}
```

### Testing
- Use table-driven tests for multiple test cases
- Use `t.Fatalf` for fatal errors, `t.Errorf` for assertions
- Test file naming: `<package>_test.go` or `<file>_test.go`
- Test package name: `<package>_test` for blackbox tests, `<package>` for whitebox

```go
func TestBuffer_Insert(t *testing.T) {
    tests := []struct {
        name    string
        initial []string
        pos     Position
        text    string
        want    []string
        wantErr bool
    }{
        {name: "insert at beginning", ...},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

### Project Structure
```
ted/
├── core/          # Core editing logic (buffer, file, history, clipboard)
├── editor/        # Main editor controller
├── ui/            # User interface (renderer, layout, menus, terminal)
├── search/        # Search and replace functionality
├── syntax/        # Syntax highlighting
├── utils/         # Utility functions
├── main.go        # Entry point
└── bin/           # Build output
```

### Dependencies
- Keep dependencies minimal
- Use standard library when possible
- Main external deps: tcell (terminal UI), clipboard

### Development Workflow
1. Run `make fmt` before committing
2. Run `make check` to ensure everything passes
3. Terminal may need reset after tests: `make reset-terminal`
4. Use `make dev` for auto-rebuild development (requires `air` tool)

## Additional Notes

- **No AI Rules Found**: No Cursor rules (.cursorrules, .cursor/rules/) or Copilot instructions (.github/copilot-instructions.md) exist in this repo.
- Go version: 1.25.5
- Module path must be `github.com/AndrewDonelson/ted`
