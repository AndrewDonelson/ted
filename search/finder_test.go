package search

import (
	"testing"

	"github.com/AndrewDonelson/ted/core/buffer"
)

func TestNewFinder(t *testing.T) {
	f := NewFinder()

	if f.pattern != "" {
		t.Errorf("initial pattern = %q, want empty string", f.pattern)
	}

	if len(f.matches) != 0 {
		t.Errorf("initial matches length = %d, want 0", len(f.matches))
	}

	if f.currentIndex != 0 {
		t.Errorf("initial currentIndex = %d, want 0", f.currentIndex)
	}

	if len(f.history) != 0 {
		t.Errorf("initial history length = %d, want 0", len(f.history))
	}

	opts := f.GetOptions()
	if opts.CaseSensitive {
		t.Error("default CaseSensitive should be false")
	}
	if opts.WholeWord {
		t.Error("default WholeWord should be false")
	}
	if opts.UseRegex {
		t.Error("default UseRegex should be false")
	}
	if !opts.WrapAround {
		t.Error("default WrapAround should be true")
	}
}

func TestFinder_SetPattern(t *testing.T) {
	f := NewFinder()

	f.SetPattern("test")

	if f.pattern != "test" {
		t.Errorf("pattern = %q, want %q", f.pattern, "test")
	}

	// Pattern should be added to history
	if len(f.history) != 1 {
		t.Errorf("history length = %d, want 1", len(f.history))
	}

	if f.history[0] != "test" {
		t.Errorf("history[0] = %q, want %q", f.history[0], "test")
	}

	// Setting same pattern again should not add to history
	f.SetPattern("test")
	if len(f.history) != 1 {
		t.Errorf("history length after duplicate = %d, want 1", len(f.history))
	}
}

func TestFinder_SetOptions(t *testing.T) {
	f := NewFinder()

	opts := Options{
		CaseSensitive: true,
		WholeWord:     true,
		UseRegex:      false,
		WrapAround:    false,
	}

	f.SetOptions(opts)

	retrieved := f.GetOptions()
	if !retrieved.CaseSensitive {
		t.Error("CaseSensitive should be true")
	}
	if !retrieved.WholeWord {
		t.Error("WholeWord should be true")
	}
	if retrieved.UseRegex {
		t.Error("UseRegex should be false")
	}
	if retrieved.WrapAround {
		t.Error("WrapAround should be false")
	}
}

func TestFinder_History(t *testing.T) {
	f := NewFinder()

	// Add patterns
	f.SetPattern("first")
	f.SetPattern("second")
	f.SetPattern("third")

	history := f.GetHistory()
	if len(history) != 3 {
		t.Errorf("history length = %d, want 3", len(history))
	}

	// Test PreviousHistory
	item, ok := f.PreviousHistory()
	if !ok {
		t.Error("PreviousHistory should return true")
	}
	if item != "second" {
		t.Errorf("PreviousHistory = %q, want %q", item, "second")
	}

	// Test NextHistory
	item, ok = f.NextHistory()
	if !ok {
		t.Error("NextHistory should return true")
	}
	if item != "third" {
		t.Errorf("NextHistory = %q, want %q", item, "third")
	}

	// Test GetHistoryItem
	item, ok = f.GetHistoryItem(0)
	if !ok {
		t.Error("GetHistoryItem(0) should return true")
	}
	if item != "first" {
		t.Errorf("GetHistoryItem(0) = %q, want %q", item, "first")
	}

	// Invalid index
	_, ok = f.GetHistoryItem(10)
	if ok {
		t.Error("GetHistoryItem(10) should return false")
	}
}

func TestFinder_FindAll_Literal(t *testing.T) {
	f := NewFinder()
	f.SetPattern("hello")

	buf := buffer.NewBuffer()
	buf.SetLines([]string{
		"hello world",
		"hello there",
		"goodbye",
	})

	matches := f.FindAll(buf)

	if len(matches) != 2 {
		t.Errorf("found %d matches, want 2", len(matches))
	}

	if len(matches) >= 1 {
		if matches[0].StartLine != 0 {
			t.Errorf("match[0].StartLine = %d, want 0", matches[0].StartLine)
		}
		if matches[0].StartCol != 0 {
			t.Errorf("match[0].StartCol = %d, want 0", matches[0].StartCol)
		}
		if matches[0].Text != "hello" {
			t.Errorf("match[0].Text = %q, want %q", matches[0].Text, "hello")
		}
	}

	if len(matches) >= 2 {
		if matches[1].StartLine != 1 {
			t.Errorf("match[1].StartLine = %d, want 1", matches[1].StartLine)
		}
	}
}

func TestFinder_FindAll_CaseSensitive(t *testing.T) {
	f := NewFinder()
	f.SetPattern("Hello")
	opts := f.GetOptions()
	opts.CaseSensitive = true
	f.SetOptions(opts)

	buf := buffer.NewBuffer()
	buf.SetLines([]string{
		"hello world",
		"Hello there",
	})

	matches := f.FindAll(buf)

	if len(matches) != 1 {
		t.Errorf("found %d matches, want 1", len(matches))
	}

	if len(matches) >= 1 {
		if matches[0].Text != "Hello" {
			t.Errorf("match.Text = %q, want %q", matches[0].Text, "Hello")
		}
	}
}

func TestFinder_FindAll_CaseInsensitive(t *testing.T) {
	f := NewFinder()
	f.SetPattern("hello")
	opts := f.GetOptions()
	opts.CaseSensitive = false
	f.SetOptions(opts)

	buf := buffer.NewBuffer()
	buf.SetLines([]string{
		"hello world",
		"Hello there",
		"HELLO everyone",
	})

	matches := f.FindAll(buf)

	if len(matches) != 3 {
		t.Errorf("found %d matches, want 3", len(matches))
	}
}

func TestFinder_FindAll_WholeWord(t *testing.T) {
	f := NewFinder()
	f.SetPattern("test")
	opts := f.GetOptions()
	opts.WholeWord = true
	f.SetOptions(opts)

	buf := buffer.NewBuffer()
	buf.SetLines([]string{
		"this is a test",
		"testing is fun",
		"best test here",
	})

	matches := f.FindAll(buf)

	// Should match "test" but not "testing" or "best"
	if len(matches) != 2 {
		t.Errorf("found %d matches, want 2", len(matches))
	}
}

func TestFinder_FindAll_Regex(t *testing.T) {
	f := NewFinder()
	f.SetPattern("h.llo")
	opts := f.GetOptions()
	opts.UseRegex = true
	f.SetOptions(opts)

	buf := buffer.NewBuffer()
	buf.SetLines([]string{
		"hello world",
		"hallo there",
		"hillo everyone",
	})

	matches := f.FindAll(buf)

	if len(matches) != 3 {
		t.Errorf("found %d matches, want 3", len(matches))
	}
}

func TestFinder_FindNext(t *testing.T) {
	f := NewFinder()
	f.SetPattern("test")

	buf := buffer.NewBuffer()
	buf.SetLines([]string{
		"test one",
		"test two",
		"test three",
	})

	// Find from beginning - FindNext finds matches AFTER the position, so from {0,0} it finds line 1
	match, found := f.FindNext(buf, buffer.Position{Line: 0, Col: 0})
	if !found {
		t.Error("FindNext should find first match")
	}
	if match.StartLine != 1 {
		t.Errorf("first match.StartLine = %d, want 1", match.StartLine)
	}

	// Find from after first match
	match, found = f.FindNext(buf, buffer.Position{Line: 0, Col: 5})
	if !found {
		t.Error("FindNext should find second match")
	}
	if match.StartLine != 1 {
		t.Errorf("second match.StartLine = %d, want 1", match.StartLine)
	}
}

func TestFinder_FindNext_WrapAround(t *testing.T) {
	f := NewFinder()
	f.SetPattern("test")

	buf := buffer.NewBuffer()
	buf.SetLines([]string{
		"test one",
		"test two",
	})

	// Find from after all matches with wrap
	match, found := f.FindNext(buf, buffer.Position{Line: 1, Col: 10})
	if !found {
		t.Error("FindNext should find match with wrap")
	}
	if match.StartLine != 0 {
		t.Errorf("wrapped match.StartLine = %d, want 0", match.StartLine)
	}

	// Test without wrap
	opts := f.GetOptions()
	opts.WrapAround = false
	f.SetOptions(opts)
	f.Clear()

	_, found = f.FindNext(buf, buffer.Position{Line: 1, Col: 10})
	if found {
		t.Error("FindNext should not find match without wrap")
	}
}

func TestFinder_FindPrevious(t *testing.T) {
	f := NewFinder()
	f.SetPattern("test")

	buf := buffer.NewBuffer()
	buf.SetLines([]string{
		"test one",
		"test two",
		"test three",
	})

	// Find from end
	match, found := f.FindPrevious(buf, buffer.Position{Line: 2, Col: 10})
	if !found {
		t.Error("FindPrevious should find last match")
	}
	if match.StartLine != 2 {
		t.Errorf("last match.StartLine = %d, want 2", match.StartLine)
	}

	// Find from before last match
	match, found = f.FindPrevious(buf, buffer.Position{Line: 2, Col: 0})
	if !found {
		t.Error("FindPrevious should find second match")
	}
	if match.StartLine != 1 {
		t.Errorf("second match.StartLine = %d, want 1", match.StartLine)
	}
}

func TestFinder_GetCurrentMatch(t *testing.T) {
	f := NewFinder()
	f.SetPattern("test")

	buf := buffer.NewBuffer()
	buf.SetLines([]string{
		"test one",
		"test two",
	})

	f.FindAll(buf)

	// No current match set yet
	_, ok := f.GetCurrentMatch()
	if ok {
		t.Error("GetCurrentMatch should return false when no match is current")
	}

	// Set current match
	f.SetCurrentMatch(0)
	match, ok := f.GetCurrentMatch()
	if !ok {
		t.Error("GetCurrentMatch should return true after SetCurrentMatch")
	}
	if match.StartLine != 0 {
		t.Errorf("current match.StartLine = %d, want 0", match.StartLine)
	}

	// Invalid index
	if f.SetCurrentMatch(100) {
		t.Error("SetCurrentMatch with invalid index should return false")
	}
}

func TestFinder_GetMatchCount(t *testing.T) {
	f := NewFinder()
	f.SetPattern("test")

	buf := buffer.NewBuffer()
	buf.SetLines([]string{
		"test one",
		"test two",
		"no match here",
	})

	f.FindAll(buf)

	if f.GetMatchCount() != 2 {
		t.Errorf("GetMatchCount = %d, want 2", f.GetMatchCount())
	}
}

func TestFinder_Clear(t *testing.T) {
	f := NewFinder()
	f.SetPattern("test")

	buf := buffer.NewBuffer()
	buf.SetLines([]string{"test"})
	f.FindAll(buf)
	f.SetCurrentMatch(0)

	f.Clear()

	if len(f.matches) != 0 {
		t.Errorf("matches length after Clear = %d, want 0", len(f.matches))
	}

	if f.currentIndex != -1 {
		t.Errorf("currentIndex after Clear = %d, want -1", f.currentIndex)
	}

	// Pattern should still be set
	if f.pattern != "test" {
		t.Errorf("pattern after Clear = %q, want %q", f.pattern, "test")
	}
}

func TestFinder_Reset(t *testing.T) {
	f := NewFinder()
	f.SetPattern("test")

	buf := buffer.NewBuffer()
	buf.SetLines([]string{"test"})
	f.FindAll(buf)

	f.Reset()

	if f.pattern != "" {
		t.Errorf("pattern after Reset = %q, want empty", f.pattern)
	}

	if len(f.matches) != 0 {
		t.Errorf("matches length after Reset = %d, want 0", len(f.matches))
	}

	opts := f.GetOptions()
	if opts.CaseSensitive {
		t.Error("options should be reset to defaults")
	}
}

func TestEscapeLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"test.exe", `test\.exe`},
		{"a+b", `a\+b`},
		{"path/to/file", `path/to/file`}, // / is not a special regex char
		{"[test]", `\[test\]`},
	}

	for _, tt := range tests {
		result := EscapeLiteral(tt.input)
		if result != tt.expected {
			t.Errorf("EscapeLiteral(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestIsWordChar(t *testing.T) {
	tests := []struct {
		char     byte
		expected bool
	}{
		{'a', true},
		{'Z', true},
		{'5', true},
		{'_', true},
		{' ', false},
		{'.', false},
		{'-', false},
		{'!', false},
	}

	for _, tt := range tests {
		result := isWordChar(tt.char)
		if result != tt.expected {
			t.Errorf("isWordChar(%q) = %v, want %v", tt.char, result, tt.expected)
		}
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()

	if opts.CaseSensitive {
		t.Error("CaseSensitive should default to false")
	}
	if opts.WholeWord {
		t.Error("WholeWord should default to false")
	}
	if opts.UseRegex {
		t.Error("UseRegex should default to false")
	}
	if !opts.WrapAround {
		t.Error("WrapAround should default to true")
	}
}

func TestFinder_EmptyPattern(t *testing.T) {
	f := NewFinder()
	f.SetPattern("")

	buf := buffer.NewBuffer()
	buf.SetLines([]string{"hello world"})

	matches := f.FindAll(buf)
	if len(matches) != 0 {
		t.Errorf("empty pattern should return 0 matches, got %d", len(matches))
	}

	match, found := f.FindNext(buf, buffer.Position{Line: 0, Col: 0})
	if found {
		t.Error("empty pattern should not find any matches")
	}
	if match.Text != "" {
		t.Error("empty match should have empty text")
	}
}

func TestFinder_RegexInvalidPattern(t *testing.T) {
	f := NewFinder()
	f.SetPattern("[invalid")
	opts := f.GetOptions()
	opts.UseRegex = true
	f.SetOptions(opts)

	buf := buffer.NewBuffer()
	buf.SetLines([]string{"test"})

	// Should not panic, just return no matches
	matches := f.FindAll(buf)
	if len(matches) != 0 {
		t.Errorf("invalid regex should return 0 matches, got %d", len(matches))
	}
}

func TestFinder_MultipleMatchesOnSameLine(t *testing.T) {
	f := NewFinder()
	f.SetPattern("test")

	buf := buffer.NewBuffer()
	buf.SetLines([]string{
		"test test test",
	})

	matches := f.FindAll(buf)
	if len(matches) != 3 {
		t.Errorf("found %d matches, want 3", len(matches))
	}

	// Verify positions
	for i, match := range matches {
		expectedCol := i * 5 // "test " is 5 chars
		if match.StartCol != expectedCol {
			t.Errorf("match[%d].StartCol = %d, want %d", i, match.StartCol, expectedCol)
		}
	}
}

func TestFinder_MatchAcrossLines(t *testing.T) {
	f := NewFinder()
	f.SetPattern("test")

	buf := buffer.NewBuffer()
	buf.SetLines([]string{
		"first test",
		"second test",
		"third test",
	})

	matches := f.FindAll(buf)

	for i, match := range matches {
		if match.StartLine != i {
			t.Errorf("match[%d].StartLine = %d, want %d", i, match.StartLine, i)
		}
		if match.EndLine != i {
			t.Errorf("match[%d].EndLine = %d, want %d", i, match.EndLine, i)
		}
	}
}
