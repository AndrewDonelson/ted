package file

import (
	"os"
	"reflect"
	"testing"
)

func TestReadFile(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		want     []string
		wantErr  bool
		setup    func(string) error
		teardown func(string) error
	}{
		{
			name:    "read simple file",
			content: "line1\nline2\nline3",
			want:    []string{"line1", "line2", "line3"},
			wantErr: false,
		},
		{
			name:    "read file with LF endings",
			content: "line1\nline2\nline3\n",
			want:    []string{"line1", "line2", "line3", ""},
			wantErr: false,
		},
		{
			name:    "read file with CRLF endings",
			content: "line1\r\nline2\r\nline3",
			want:    []string{"line1", "line2", "line3"},
			wantErr: false,
		},
		{
			name:    "read file with CR endings",
			content: "line1\rline2\rline3",
			want:    []string{"line1", "line2", "line3"},
			wantErr: false,
		},
		{
			name:    "read empty file",
			content: "",
			want:    []string{""},
			wantErr: false,
		},
		{
			name:    "read single line file",
			content: "single line",
			want:    []string{"single line"},
			wantErr: false,
		},
		{
			name:    "read file with mixed content",
			content: "line1\nline2\r\nline3",
			want:    []string{"line1", "line2", "line3"},
			wantErr: false,
		},
		{
			name:    "read file with UTF-8 characters",
			content: "hello 世界\n测试\ncafé",
			want:    []string{"hello 世界", "测试", "café"},
			wantErr: false,
		},
		{
			name:    "read file with trailing newline",
			content: "line1\nline2\n",
			want:    []string{"line1", "line2", ""},
			wantErr: false,
		},
		{
			name:    "read file without trailing newline",
			content: "line1\nline2",
			want:    []string{"line1", "line2"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file
			tmpfile, err := os.CreateTemp("", "test*.txt")
			if err != nil {
				t.Fatalf("CreateTemp() error = %v", err)
			}
			defer os.Remove(tmpfile.Name())

			// Write content
			if _, err := tmpfile.WriteString(tt.content); err != nil {
				t.Fatalf("WriteString() error = %v", err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatalf("Close() error = %v", err)
			}

			// Read file
			got, err := ReadFile(tmpfile.Name())

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadFile_Errors(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
		setup   func() (string, func())
	}{
		{
			name:    "read non-existent file",
			path:    "/nonexistent/file.txt",
			wantErr: true,
		},
		{
			name:    "read directory",
			wantErr: true,
			setup: func() (string, func()) {
				tmpdir, err := os.MkdirTemp("", "testdir")
				if err != nil {
					t.Fatalf("MkdirTemp() error = %v", err)
				}
				return tmpdir, func() { os.RemoveAll(tmpdir) }
			},
		},
		{
			name:    "read with empty path",
			path:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var path string
			var cleanup func()

			if tt.setup != nil {
				path, cleanup = tt.setup()
				defer cleanup()
			} else {
				path = tt.path
			}

			_, err := ReadFile(path)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReadFileWithInfo(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatalf("CreateTemp() error = %v", err)
	}
	defer os.Remove(tmpfile.Name())

	content := "line1\r\nline2\r\nline3"
	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatalf("WriteString() error = %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	lines, info, err := ReadFileWithInfo(tmpfile.Name())
	if err != nil {
		t.Fatalf("ReadFileWithInfo() error = %v", err)
	}

	if len(lines) != 3 {
		t.Errorf("ReadFileWithInfo() lines = %d, want 3", len(lines))
	}

	if info == nil {
		t.Fatal("ReadFileWithInfo() info = nil, want non-nil")
	}

	if info.Path != tmpfile.Name() {
		t.Errorf("ReadFileWithInfo() info.Path = %q, want %q", info.Path, tmpfile.Name())
	}

	if info.Encoding != "UTF-8" {
		t.Errorf("ReadFileWithInfo() info.Encoding = %q, want UTF-8", info.Encoding)
	}

	if info.LineEnding != LineEndingCRLF {
		t.Errorf("ReadFileWithInfo() info.LineEnding = %q, want CRLF", info.LineEnding)
	}
}

func TestDetectLineEnding(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    LineEnding
	}{
		{
			name:    "detect LF",
			content: "line1\nline2\nline3",
			want:    LineEndingLF,
		},
		{
			name:    "detect CRLF",
			content: "line1\r\nline2\r\nline3",
			want:    LineEndingCRLF,
		},
		{
			name:    "detect CR",
			content: "line1\rline2\rline3",
			want:    LineEndingCR,
		},
		{
			name:    "detect LF for single line",
			content: "single line",
			want:    LineEndingLF,
		},
		{
			name:    "detect CRLF when mixed",
			content: "line1\nline2\r\nline3",
			want:    LineEndingCRLF, // CRLF takes precedence
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfile, err := os.CreateTemp("", "test*.txt")
			if err != nil {
				t.Fatalf("CreateTemp() error = %v", err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.WriteString(tt.content); err != nil {
				t.Fatalf("WriteString() error = %v", err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatalf("Close() error = %v", err)
			}

			got := detectLineEnding(tmpfile.Name())
			if got != tt.want {
				t.Errorf("detectLineEnding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidatePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "valid relative path",
			path:    "test.txt",
			wantErr: false,
		},
		{
			name:    "valid absolute path",
			path:    "/tmp/test.txt",
			wantErr: false,
		},
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
		},
		{
			name:    "path with ..",
			path:    "../test.txt",
			wantErr: false, // Should be cleaned, not error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validatePath(tt.path)

			if (err != nil) != tt.wantErr {
				t.Errorf("validatePath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSplitLines(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want []string
	}{
		{
			name: "split LF",
			data: []byte("line1\nline2\nline3"),
			want: []string{"line1", "line2", "line3"},
		},
		{
			name: "split CRLF",
			data: []byte("line1\r\nline2\r\nline3"),
			want: []string{"line1", "line2", "line3"},
		},
		{
			name: "split CR",
			data: []byte("line1\rline2\rline3"),
			want: []string{"line1", "line2", "line3"},
		},
		{
			name: "split with trailing newline",
			data: []byte("line1\nline2\n"),
			want: []string{"line1", "line2", ""},
		},
		{
			name: "split empty",
			data: []byte(""),
			want: []string{""},
		},
		{
			name: "split single line",
			data: []byte("single"),
			want: []string{"single"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitLines(tt.data)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitLines() = %v, want %v", got, tt.want)
			}
		})
	}
}
