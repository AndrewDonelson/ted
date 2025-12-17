package file

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestWriteFile(t *testing.T) {
	tests := []struct {
		name       string
		lines      []string
		lineEnding LineEnding
		want       string
		wantErr    bool
	}{
		{
			name:       "write with LF",
			lines:      []string{"line1", "line2", "line3"},
			lineEnding: LineEndingLF,
			want:       "line1\nline2\nline3",
			wantErr:    false,
		},
		{
			name:       "write with CRLF",
			lines:      []string{"line1", "line2", "line3"},
			lineEnding: LineEndingCRLF,
			want:       "line1\r\nline2\r\nline3",
			wantErr:    false,
		},
		{
			name:       "write with CR",
			lines:      []string{"line1", "line2", "line3"},
			lineEnding: LineEndingCR,
			want:       "line1\rline2\rline3",
			wantErr:    false,
		},
		{
			name:       "write single line",
			lines:      []string{"single line"},
			lineEnding: LineEndingLF,
			want:       "single line",
			wantErr:    false,
		},
		{
			name:       "write empty lines",
			lines:      []string{"", "", ""},
			lineEnding: LineEndingLF,
			want:       "\n\n",
			wantErr:    false,
		},
		{
			name:       "write with trailing newline",
			lines:      []string{"line1", "line2", ""},
			lineEnding: LineEndingLF,
			want:       "line1\nline2\n",
			wantErr:    false,
		},
		{
			name:       "write empty file",
			lines:      []string{""},
			lineEnding: LineEndingLF,
			want:       "",
			wantErr:    false,
		},
		{
			name:       "write with UTF-8",
			lines:      []string{"hello 世界", "测试", "café"},
			lineEnding: LineEndingLF,
			want:       "hello 世界\n测试\ncafé",
			wantErr:    false,
		},
		{
			name:       "write with empty path",
			lines:      []string{"line1"},
			lineEnding: LineEndingLF,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var path string
			if tt.name == "write with empty path" {
				path = ""
			} else {
				tmpfile, err := os.CreateTemp("", "test*.txt")
				if err != nil {
					t.Fatalf("CreateTemp() error = %v", err)
				}
				path = tmpfile.Name()
				tmpfile.Close()
				defer os.Remove(path)
			}

			err := WriteFile(path, tt.lines, tt.lineEnding)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Read back and verify
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("ReadFile() error = %v", err)
				}

				got := string(data)
				if got != tt.want {
					t.Errorf("WriteFile() content = %q, want %q", got, tt.want)
				}
			}
		})
	}
}

func TestWriteFilePreserveEnding(t *testing.T) {
	tests := []struct {
		name         string
		initial      string
		lines        []string
		wantEnding   LineEnding
		wantContains string
	}{
		{
			name:         "preserve CRLF",
			initial:      "line1\r\nline2\r\n",
			lines:        []string{"new1", "new2"},
			wantEnding:   LineEndingCRLF,
			wantContains: "\r\n",
		},
		{
			name:         "preserve LF",
			initial:      "line1\nline2\n",
			lines:        []string{"new1", "new2"},
			wantEnding:   LineEndingLF,
			wantContains: "\n",
		},
		{
			name:         "default to LF for new file",
			initial:      "",
			lines:        []string{"new1", "new2"},
			wantEnding:   LineEndingLF,
			wantContains: "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfile, err := os.CreateTemp("", "test*.txt")
			if err != nil {
				t.Fatalf("CreateTemp() error = %v", err)
			}
			path := tmpfile.Name()
			tmpfile.Close()
			defer os.Remove(path)

			// Write initial content if provided
			if tt.initial != "" {
				if err := os.WriteFile(path, []byte(tt.initial), 0644); err != nil {
					t.Fatalf("WriteFile() error = %v", err)
				}
			}

			// Write new content
			if err := WriteFilePreserveEnding(path, tt.lines); err != nil {
				t.Fatalf("WriteFilePreserveEnding() error = %v", err)
			}

			// Verify line ending
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("ReadFile() error = %v", err)
			}

			content := string(data)
			if !strings.Contains(content, tt.wantContains) {
				t.Errorf("WriteFilePreserveEnding() content = %q, want to contain %q", content, tt.wantContains)
			}

			// Verify detection
			detected := detectLineEnding(path)
			if detected != tt.wantEnding {
				t.Errorf("WriteFilePreserveEnding() detected ending = %v, want %v", detected, tt.wantEnding)
			}
		})
	}
}

func TestWriteFile_RoundTrip(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		lineEnding LineEnding
	}{
		{
			name:       "round-trip LF",
			content:    "line1\nline2\nline3",
			lineEnding: LineEndingLF,
		},
		{
			name:       "round-trip CRLF",
			content:    "line1\r\nline2\r\nline3",
			lineEnding: LineEndingCRLF,
		},
		{
			name:       "round-trip CR",
			content:    "line1\rline2\rline3",
			lineEnding: LineEndingCR,
		},
		{
			name:       "round-trip with UTF-8",
			content:    "hello 世界\n测试\ncafé",
			lineEnding: LineEndingLF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfile, err := os.CreateTemp("", "test*.txt")
			if err != nil {
				t.Fatalf("CreateTemp() error = %v", err)
			}
			path := tmpfile.Name()
			tmpfile.Close()
			defer os.Remove(path)

			// Write initial content
			if err := os.WriteFile(path, []byte(tt.content), 0644); err != nil {
				t.Fatalf("WriteFile() error = %v", err)
			}

			// Read file
			lines, err := ReadFile(path)
			if err != nil {
				t.Fatalf("ReadFile() error = %v", err)
			}

			// Write back with same line ending
			if err := WriteFile(path, lines, tt.lineEnding); err != nil {
				t.Fatalf("WriteFile() error = %v", err)
			}

			// Read back and verify
			readLines, err := ReadFile(path)
			if err != nil {
				t.Fatalf("ReadFile() error = %v", err)
			}

			if !reflect.DeepEqual(readLines, lines) {
				t.Errorf("Round-trip failed: got %v, want %v", readLines, lines)
			}
		})
	}
}

func TestAtomicWrite(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatalf("CreateTemp() error = %v", err)
	}
	path := tmpfile.Name()
	tmpfile.Close()
	defer os.Remove(path)

	data := []byte("test content")
	if err := atomicWrite(path, data); err != nil {
		t.Fatalf("atomicWrite() error = %v", err)
	}

	// Verify file exists and has correct content
	readData, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if string(readData) != string(data) {
		t.Errorf("atomicWrite() content = %q, want %q", string(readData), string(data))
	}

	// Verify no temp files remain
	dir := filepath.Dir(path)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir() error = %v", err)
	}

	for _, entry := range entries {
		if strings.Contains(entry.Name(), ".tmp.") {
			t.Errorf("atomicWrite() left temp file: %q", entry.Name())
		}
	}
}

func TestAtomicWrite_CreateDirectory(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("MkdirTemp() error = %v", err)
	}
	defer os.RemoveAll(tmpdir)

	// Write to non-existent subdirectory (atomicWrite requires directory to exist)
	subdir := filepath.Join(tmpdir, "subdir")
	if err := os.MkdirAll(subdir, 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	path := filepath.Join(subdir, "test.txt")
	data := []byte("test content")

	if err := atomicWrite(path, data); err != nil {
		t.Fatalf("atomicWrite() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(path); err != nil {
		t.Errorf("atomicWrite() file does not exist: %v", err)
	}
}

func TestLineEndingToString(t *testing.T) {
	tests := []struct {
		name   string
		ending LineEnding
		want   string
	}{
		{
			name:   "LF to string",
			ending: LineEndingLF,
			want:   "\n",
		},
		{
			name:   "CRLF to string",
			ending: LineEndingCRLF,
			want:   "\r\n",
		},
		{
			name:   "CR to string",
			ending: LineEndingCR,
			want:   "\r",
		},
		{
			name:   "unknown defaults to LF",
			ending: LineEndingUnknown,
			want:   "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lineEndingToString(tt.ending)
			if got != tt.want {
				t.Errorf("lineEndingToString() = %q, want %q", got, tt.want)
			}
		})
	}
}
