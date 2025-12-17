# ted Constitution

## Core Principles

### I. Modular Architecture
Every module must be independent, testable, and replaceable. Core modules (`core/*`) MUST NOT depend on UI modules (`ui/*`). Each package should have a single, clear responsibility. Interfaces are defined before implementations. Dependency injection is used for testability. No circular dependencies are allowed.

### II. Code Quality Standards (NON-NEGOTIABLE)
All code must follow idiomatic Go conventions. Code must pass `gofmt`, `go vet`, and `golint` with zero warnings. Follow [Effective Go](https://golang.org/doc/effective_go.html) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments). All exported types and functions must have complete godoc comments. Error handling must be explicit - no ignored errors, no panics in production code.

### III. Test-First Development
Comprehensive testing is mandatory with high coverage requirements: `core/*` packages require 90%+ coverage, `editor/*` packages require 85%+ coverage, `ui/*` packages require 70%+ coverage, with overall project coverage of 80%+. Use table-driven tests for functions with multiple scenarios. Test edge cases including empty files, boundary conditions, UTF-8 handling, and error scenarios. Integration tests required for end-to-end workflows.

### IV. Performance Requirements
Startup time must be < 100ms for files < 1MB. Keystroke latency must be < 16ms (60fps feel). Memory usage should be < 100MB for typical sessions. Support files up to 100MB. Profile before optimizing - correctness first, performance second.

### V. Security & Error Handling
All file paths must be validated (prevent directory traversal). Check file permissions before reading/writing. Handle symlinks securely. Limit file size to prevent memory exhaustion. Validate all keyboard/mouse input. Handle malformed UTF-8 gracefully. Properly restore terminal state on exit. Return errors, don't panic (except in truly exceptional cases). Wrap errors with context using `fmt.Errorf("context: %w", err)`.

### VI. Fixed User Experience
Keyboard shortcuts are FIXED and HARDCODED - NO custom keybindings allowed. This ensures consistent experience across all machines. The info bar MUST use inverted colors (light background #d4d4d4 with dark text #1e1e1e) for visual distinction. Windows-style keyboard shortcuts are used throughout. No modal editing (no vim modes).

## Additional Constraints

### Technology Stack
- **Language:** Go (latest stable version)
- **Terminal Library:** `github.com/gdamore/tcell/v2` (REQUIRED)
- **Clipboard Library:** `github.com/atotto/clipboard` (REQUIRED)
- **Configuration Format:** TOML (Phase 5+)
- **Dependency Management:** Go modules (`go.mod`, `go.sum`)
- Keep dependencies minimal and audit security

### Cross-Platform Requirements
Must work identically on Linux, macOS, and Windows. Use platform-agnostic abstractions. Test on all target platforms before release.

### UI/UX Constraints
- Info bar MUST be inverted (light bg, dark text) - this is a critical visual requirement
- Menu bar always visible (no auto-hide)
- Responsive layout that adapts to terminal resize
- Minimum terminal size: 40 columns × 10 rows
- Dark mode default (Phase 0-4), light mode option (Phase 5+)

## Development Workflow

### Phased Development
Development follows a strict 6-phase plan:
- **Phase 0:** Foundation - Basic text editing with full UI layout
- **Phase 1:** Essential editing - Undo/redo, clipboard, line operations
- **Phase 2:** Navigation - Search, line numbers, fast navigation
- **Phase 3:** Code editing - Syntax highlighting, comments
- **Phase 4:** Advanced search - Replace, regex
- **Phase 5:** Configuration - Customization, themes
- **Phase 6:** Multiple files - Tabs/buffers

Each phase must be complete and tested before moving to the next.

### Git Workflow
- `main` branch: Stable releases only
- `develop` branch: Integration branch for development
- `feature/phase-N-description`: Feature branches for each phase
- Use Conventional Commits format: `<type>(<scope>): <subject>`
- Types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`

### Code Review Requirements
- All PRs must pass `go test ./...` with no failures
- Test coverage must meet minimum requirements
- All linters must pass (`go vet`, `golint`, `gofmt`)
- No TODO comments in production code
- All exported identifiers must have godoc comments
- Manual testing checklist must be completed

### Documentation Standards
- Package-level godoc for all packages
- Type and function documentation for all exported identifiers
- README.md with installation and usage instructions
- ARCHITECTURE.md for system design (optional for Phase 0)
- CONTRIBUTING.md for development guidelines
- Examples in documentation must work

## Governance

The constitution supersedes all other practices. All PRs and reviews must verify compliance with these principles. Complexity must be justified. Amendments require documentation, approval, and a migration plan. Use the design documents (`.design/prompt.md` and `.design/complete.md`) for detailed implementation guidance.

**Version**: 1.0.0 | **Ratified**: 2025-12-16 | **Last Amended**: 2025-12-16
