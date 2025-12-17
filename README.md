# ted - Terminal EDitor

[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/AndrewDonelson/ted.svg)](https://pkg.go.dev/github.com/AndrewDonelson/ted)
[![Go Report Card](https://goreportcard.com/badge/github.com/AndrewDonelson/ted?style=flat-square)](https://goreportcard.com/report/github.com/AndrewDonelson/ted)
[![GitHub stars](https://img.shields.io/github/stars/AndrewDonelson/ted.svg?style=flat-square&label=stars)](https://github.com/AndrewDonelson/ted/stargazers)
[![GitHub release](https://img.shields.io/github/release/AndrewDonelson/ted.svg?style=flat-square)](https://github.com/AndrewDonelson/ted/releases)

[![Buy Me A Coffee](https://img.shields.io/badge/Buy%20Me%20A%20Coffee-FFDD00?style=for-the-badge&logo=buy-me-a-coffee&logoColor=black)](https://buymeacoffee.com/andrewdonelson)

**ted** is a modern, cross-platform command-line text editor written in Go that uses familiar Windows-style keyboard shortcuts and intuitive arrow key navigation.

## What is ted?

ted (Terminal EDitor) is a friendly, approachable terminal text editor designed for developers who want a powerful editing experience without the learning curve of modal editors like vim. It works in any terminal on Linux, macOS, and Windows, providing a consistent experience across all platforms.

## Why ted?

Most terminal editors fall into two categories: simple but limited (like nano) or powerful but complex (like vim/emacs). ted bridges this gap by offering:

- **Familiar shortcuts** - Uses the same keyboard shortcuts you already know from modern editors (Ctrl+S to save, Ctrl+C to copy, etc.)
- **No modes** - Start typing immediately, no need to learn insert/command modes
- **Fast and lightweight** - Written in Go for speed with minimal dependencies
- **Cross-platform** - Identical experience whether you're on Linux, macOS, or Windows
- **Terminal native** - Works in any terminal, no GUI required

## Features

### Essential Editing
- Cut, Copy, Paste (Ctrl+X, Ctrl+C, Ctrl+V)
- Undo/Redo (Ctrl+Z, Ctrl+Y)
- Select all (Ctrl+A)
- Delete entire line (Ctrl+Shift+K)
- Duplicate line (Ctrl+D)
- Move line up/down (Alt+Up/Down)
- Insert line above/below (Ctrl+Shift+Enter, Ctrl+Enter)

### Navigation
- Arrow keys for cursor movement
- Word navigation (Ctrl+Left/Right)
- Jump to line start/end (Home/End)
- Jump to document start/end (Ctrl+Home/End)
- Page Up/Down navigation
- Go to line number (Ctrl+G)

### Selection
- Shift+Arrow for character selection
- Shift+Home/End for line selection
- Ctrl+Shift+Arrow for word selection
- Double-click to select word (mouse support)

### Search & Replace
- Find (Ctrl+F)
- Find next/previous (F3/Shift+F3)
- Replace (Ctrl+H)
- Replace all (Ctrl+Shift+H)
- Case-sensitive and whole word options
- Regular expression support

### Code Editing
- Syntax highlighting for multiple languages
- Auto-indentation
- Comment/uncomment (Ctrl+/)
- Indent/unindent (Tab/Shift+Tab)
- Jump to matching bracket (Ctrl+B)
- Show whitespace toggle (Ctrl+Shift+I)

### Display & UI
- Toggleable line numbers (Ctrl+L)
- Current line highlighting
- Status bar with mode, encoding, and position
- Info bar with file details (filename, size, type, settings)
- Word wrap toggle (Ctrl+Shift+W)
- Responsive layout that adapts to terminal size

### File Operations
- Open files from command line
- Save (Ctrl+S) and Save As (Ctrl+Shift+S)
- New file (Ctrl+N)
- Close file (Ctrl+W)
- Quit (Ctrl+Q)
- Unsaved changes prompt on exit

## Installation

### From Source

```bash
go install github.com/AndrewDonelson/ted@latest
```

### Binary Downloads

Download pre-built binaries from the [releases page](https://github.com/AndrewDonelson/ted/releases) for your platform.

## Usage

### Basic Usage

Open a file:
```bash
ted filename.txt
```

Create a new file:
```bash
ted
```

### Keyboard Shortcuts

#### File Operations
- **Ctrl+N** - New file
- **Ctrl+O** - Open file
- **Ctrl+S** - Save
- **Ctrl+Shift+S** - Save As
- **Ctrl+W** - Close current file
- **Ctrl+Q** - Quit editor

#### Essential Editing
- **Ctrl+Z** - Undo
- **Ctrl+Y** - Redo
- **Ctrl+X** - Cut
- **Ctrl+C** - Copy
- **Ctrl+V** - Paste
- **Ctrl+A** - Select all
- **Delete** - Delete character forward
- **Backspace** - Delete character backward

#### Line Operations
- **Ctrl+Shift+K** - Delete entire line
- **Ctrl+D** - Duplicate current line
- **Alt+Up/Down** - Move line up/down
- **Ctrl+Enter** - Insert line below
- **Ctrl+Shift+Enter** - Insert line above

#### Navigation
- **Arrow Keys** - Move cursor
- **Home/End** - Line start/end
- **Ctrl+Home/End** - Document start/end
- **Ctrl+Left/Right** - Move by word
- **Page Up/Down** - Move by page
- **Ctrl+G** - Go to line number

#### Search & Replace
- **Ctrl+F** - Find
- **F3** - Find next
- **Shift+F3** - Find previous
- **Ctrl+H** - Replace
- **Ctrl+Shift+H** - Replace all
- **Esc** - Close find/replace dialog

#### Code Editing
- **Ctrl+/** - Toggle line comment
- **Tab** - Indent
- **Shift+Tab** - Unindent
- **Ctrl+B** - Jump to matching bracket

#### Display
- **Ctrl+L** - Toggle line numbers
- **Ctrl+Shift+W** - Toggle word wrap
- **Ctrl+Shift+I** - Toggle show whitespace
- **F10** or **Alt+Key** - Activate menu bar

### Menu System

Access menus using **Alt+Key** (e.g., Alt+F for File menu) or **F10** to activate the first menu. Navigate with arrow keys and press Enter to select. Press **Esc** to close menus.

## Configuration

### Default Settings

- **Tab size:** 4 spaces
- **Use spaces:** Yes (not tabs)
- **Line numbers:** Off (toggle with Ctrl+L)
- **Word wrap:** On (toggle with Ctrl+Shift+W)
- **Color scheme:** Dark mode
- **Encoding:** UTF-8
- **Line ending:** Auto-detect (preserves file's original)

### Configuration File (Coming Soon)

Future versions will support a TOML configuration file at `~/.tedrc` or `~/.config/ted/config.toml` for customizing:
- Color schemes and themes
- File type associations
- Syntax highlighting preferences
- Editor settings (tab size, spaces vs tabs)
- Search defaults

**Note:** Keyboard shortcuts are fixed and cannot be customized. This ensures a consistent experience across all machines.

## Philosophy

ted is designed with a simple philosophy: **consistency and familiarity**. The keyboard shortcuts you learn on one machine work the same way everywhere. There's no configuration needed to get started, and no need to memorize complex key combinations. Just open a file and start editing.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines.

## Repository

**GitHub:** [github.com/AndrewDonelson/ted](https://github.com/AndrewDonelson/ted)
