package editor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AndrewDonelson/ted/core/buffer"
	"github.com/AndrewDonelson/ted/core/file"
)

func TestNewEditor(t *testing.T) {
	// This test may fail in non-terminal environments
	ed, err := NewEditor()
	if err != nil {
		t.Skipf("Skipping test - terminal not available: %v", err)
		return
	}
	defer ed.screen.Fini()

	if ed == nil {
		t.Fatal("NewEditor() returned nil")
	}

	if ed.buffer == nil {
		t.Error("Editor.buffer is nil")
	}

	if ed.layout == nil {
		t.Error("Editor.layout is nil")
	}

	if ed.renderer == nil {
		t.Error("Editor.renderer is nil")
	}

	if ed.menuBar == nil {
		t.Error("Editor.menuBar is nil")
	}

	if ed.mode != ModeInsert {
		t.Errorf("Editor.mode = %v, want ModeInsert", ed.mode)
	}

	if ed.isDirty {
		t.Error("Editor.isDirty = true, want false")
	}
}

func TestEditor_OpenFile(t *testing.T) {
	ed, err := NewEditor()
	if err != nil {
		t.Skipf("Skipping test - terminal not available: %v", err)
		return
	}
	defer ed.screen.Fini()

	// Create a temporary file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	content := []string{"line1", "line2", "line3"}

	if err := file.WriteFile(tmpFile, content, file.LineEndingLF); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Open the file
	if err := ed.OpenFile(tmpFile); err != nil {
		t.Errorf("OpenFile() error = %v", err)
	}

	if ed.filePath != tmpFile {
		t.Errorf("Editor.filePath = %q, want %q", ed.filePath, tmpFile)
	}

	if ed.buffer.LineCount() != len(content) {
		t.Errorf("Buffer line count = %d, want %d", ed.buffer.LineCount(), len(content))
	}

	if ed.isDirty {
		t.Error("Editor.isDirty = true after opening file, want false")
	}
}

func TestEditor_OpenFile_NonExistent(t *testing.T) {
	ed, err := NewEditor()
	if err != nil {
		t.Skipf("Skipping test - terminal not available: %v", err)
		return
	}
	defer ed.screen.Fini()

	nonExistent := "/nonexistent/file.txt"
	if err := ed.OpenFile(nonExistent); err == nil {
		t.Error("OpenFile() with non-existent file should return error")
	}
}

func TestEditor_OpenFile_EmptyFile(t *testing.T) {
	ed, err := NewEditor()
	if err != nil {
		t.Skipf("Skipping test - terminal not available: %v", err)
		return
	}
	defer ed.screen.Fini()

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "empty.txt")

	// Create empty file
	if err := os.WriteFile(tmpFile, []byte{}, 0644); err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}

	if err := ed.OpenFile(tmpFile); err != nil {
		t.Errorf("OpenFile() error = %v", err)
	}

	if ed.buffer.LineCount() == 0 {
		t.Error("Empty file should have at least one line")
	}
}

func TestEditor_SaveFile(t *testing.T) {
	ed, err := NewEditor()
	if err != nil {
		t.Skipf("Skipping test - terminal not available: %v", err)
		return
	}
	defer ed.screen.Fini()

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")

	// Set up editor with content
	ed.buffer.SetLines([]string{"line1", "line2", "modified"})
	ed.filePath = tmpFile
	ed.isDirty = true

	// Save the file
	if err := ed.SaveFile(); err != nil {
		t.Errorf("SaveFile() error = %v", err)
	}

	if ed.isDirty {
		t.Error("Editor.isDirty = true after saving, want false")
	}

	// Verify file was written
	lines, _, err := file.ReadFileWithInfo(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	if len(lines) != 3 {
		t.Errorf("Saved file has %d lines, want 3", len(lines))
	}
}

func TestEditor_SaveFile_NoPath(t *testing.T) {
	ed, err := NewEditor()
	if err != nil {
		t.Skipf("Skipping test - terminal not available: %v", err)
		return
	}
	defer ed.screen.Fini()

	ed.filePath = ""
	if err := ed.SaveFile(); err == nil {
		t.Error("SaveFile() with no path should return error")
	}
}

func TestEditor_InsertCharacter(t *testing.T) {
	ed, err := NewEditor()
	if err != nil {
		t.Skipf("Skipping test - terminal not available: %v", err)
		return
	}
	defer ed.screen.Fini()

	ed.buffer.SetLines([]string{"test"})
	ed.buffer.MoveCursor(buffer.Position{Line: 0, Col: 4})

	ed.insertCharacter('!')

	if ed.buffer.LineCount() != 1 {
		t.Errorf("Line count = %d, want 1", ed.buffer.LineCount())
	}

	line, _ := ed.buffer.GetLine(0)
	if line != "test!" {
		t.Errorf("Line = %q, want %q", line, "test!")
	}

	if !ed.isDirty {
		t.Error("Editor.isDirty = false after insert, want true")
	}
}

func TestEditor_InsertCharacter_Newline(t *testing.T) {
	ed, err := NewEditor()
	if err != nil {
		t.Skipf("Skipping test - terminal not available: %v", err)
		return
	}
	defer ed.screen.Fini()

	ed.buffer.SetLines([]string{"line1"})
	ed.buffer.MoveCursor(buffer.Position{Line: 0, Col: 5})

	ed.insertCharacter('\n')

	if ed.buffer.LineCount() != 2 {
		t.Errorf("Line count = %d, want 2", ed.buffer.LineCount())
	}
}

func TestEditor_HandleBackspace(t *testing.T) {
	ed, err := NewEditor()
	if err != nil {
		t.Skipf("Skipping test - terminal not available: %v", err)
		return
	}
	defer ed.screen.Fini()

	ed.buffer.SetLines([]string{"test"})
	ed.buffer.MoveCursor(buffer.Position{Line: 0, Col: 4})

	ed.handleBackspace()

	line, _ := ed.buffer.GetLine(0)
	if line != "tes" {
		t.Errorf("Line after backspace = %q, want %q", line, "tes")
	}

	if !ed.isDirty {
		t.Error("Editor.isDirty = false after backspace, want true")
	}
}

func TestEditor_HandleBackspace_AtLineStart(t *testing.T) {
	ed, err := NewEditor()
	if err != nil {
		t.Skipf("Skipping test - terminal not available: %v", err)
		return
	}
	defer ed.screen.Fini()

	ed.buffer.SetLines([]string{"line1", "line2"})
	ed.buffer.MoveCursor(buffer.Position{Line: 1, Col: 0})

	ed.handleBackspace()

	// Should join with previous line
	if ed.buffer.LineCount() != 1 {
		t.Errorf("Line count = %d, want 1", ed.buffer.LineCount())
	}

	line, _ := ed.buffer.GetLine(0)
	if line != "line1line2" {
		t.Errorf("Line = %q, want %q", line, "line1line2")
	}
}

func TestEditor_HandleDelete(t *testing.T) {
	ed, err := NewEditor()
	if err != nil {
		t.Skipf("Skipping test - terminal not available: %v", err)
		return
	}
	defer ed.screen.Fini()

	ed.buffer.SetLines([]string{"test"})
	ed.buffer.MoveCursor(buffer.Position{Line: 0, Col: 0})

	ed.handleDelete()

	line, _ := ed.buffer.GetLine(0)
	if line != "est" {
		t.Errorf("Line after delete = %q, want %q", line, "est")
	}

	if !ed.isDirty {
		t.Error("Editor.isDirty = false after delete, want true")
	}
}

func TestEditor_HandleDelete_AtLineEnd(t *testing.T) {
	ed, err := NewEditor()
	if err != nil {
		t.Skipf("Skipping test - terminal not available: %v", err)
		return
	}
	defer ed.screen.Fini()

	ed.buffer.SetLines([]string{"line1", "line2"})
	ed.buffer.MoveCursor(buffer.Position{Line: 0, Col: 5})

	ed.handleDelete()

	// Should join with next line
	if ed.buffer.LineCount() != 1 {
		t.Errorf("Line count = %d, want 1", ed.buffer.LineCount())
	}

	line, _ := ed.buffer.GetLine(0)
	if line != "line1line2" {
		t.Errorf("Line = %q, want %q", line, "line1line2")
	}
}

func TestEditor_GetFileName(t *testing.T) {
	ed, err := NewEditor()
	if err != nil {
		t.Skipf("Skipping test - terminal not available: %v", err)
		return
	}
	defer ed.screen.Fini()

	tests := []struct {
		name     string
		filePath string
		want     string
	}{
		{
			name:     "simple filename",
			filePath: "test.txt",
			want:     "test.txt",
		},
		{
			name:     "path with forward slash",
			filePath: "/path/to/test.txt",
			want:     "test.txt",
		},
		{
			name:     "path with backslash",
			filePath: "\\path\\to\\test.txt",
			want:     "test.txt",
		},
		{
			name:     "empty path",
			filePath: "",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed.filePath = tt.filePath
			got := ed.getFileName()
			if got != tt.want {
				t.Errorf("getFileName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEditor_DetectFileType(t *testing.T) {
	ed, err := NewEditor()
	if err != nil {
		t.Skipf("Skipping test - terminal not available: %v", err)
		return
	}
	defer ed.screen.Fini()

	tests := []struct {
		name     string
		filePath string
		want     string
	}{
		{
			name:     "Go file",
			filePath: "test.go",
			want:     "Go",
		},
		{
			name:     "JavaScript file",
			filePath: "test.js",
			want:     "JavaScript",
		},
		{
			name:     "TypeScript file",
			filePath: "test.ts",
			want:     "TypeScript",
		},
		{
			name:     "Python file",
			filePath: "test.py",
			want:     "Python",
		},
		{
			name:     "Markdown file",
			filePath: "test.md",
			want:     "Markdown",
		},
		{
			name:     "text file",
			filePath: "test.txt",
			want:     "Plain Text",
		},
		{
			name:     "unknown extension",
			filePath: "test.xyz",
			want:     "Plain Text",
		},
		{
			name:     "no extension",
			filePath: "test",
			want:     "Plain Text",
		},
		{
			name:     "empty path",
			filePath: "",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed.filePath = tt.filePath
			got := ed.detectFileType()
			if got != tt.want {
				t.Errorf("detectFileType() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEditor_BuildFileInfo(t *testing.T) {
	ed, err := NewEditor()
	if err != nil {
		t.Skipf("Skipping test - terminal not available: %v", err)
		return
	}
	defer ed.screen.Fini()

	ed.buffer.SetLines([]string{"line1", "line2"})
	ed.filePath = "test.txt"
	// Modify the buffer to set the modified flag
	ed.buffer.Insert(buffer.Position{Line: 0, Col: 0}, "x")

	fileInfo := ed.buildFileInfo()

	if fileInfo == nil {
		t.Fatal("buildFileInfo() returned nil")
	}

	if fileInfo.Name != "test.txt" {
		t.Errorf("FileInfo.Name = %q, want %q", fileInfo.Name, "test.txt")
	}

	if fileInfo.TotalLines != 2 {
		t.Errorf("FileInfo.TotalLines = %d, want 2", fileInfo.TotalLines)
	}

	if !fileInfo.IsModified {
		t.Error("FileInfo.IsModified = false, want true")
	}
}

func TestEditor_BuildFileInfo_NoFile(t *testing.T) {
	ed, err := NewEditor()
	if err != nil {
		t.Skipf("Skipping test - terminal not available: %v", err)
		return
	}
	defer ed.screen.Fini()

	ed.filePath = ""
	fileInfo := ed.buildFileInfo()

	if fileInfo == nil {
		t.Fatal("buildFileInfo() returned nil")
	}

	if fileInfo.Name != "" {
		t.Errorf("FileInfo.Name = %q, want empty", fileInfo.Name)
	}
}
