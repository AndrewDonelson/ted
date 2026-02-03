package search

import (
	"testing"

	"github.com/AndrewDonelson/ted/core/buffer"
	"github.com/AndrewDonelson/ted/core/history"
)

func TestNewReplacer(t *testing.T) {
	finder := NewFinder()
	r := NewReplacer(finder)

	if r.finder != finder {
		t.Error("replacer should use the provided finder")
	}

	if r.replacement != "" {
		t.Errorf("initial replacement = %q, want empty string", r.replacement)
	}
}

func TestReplacer_SetReplacement(t *testing.T) {
	finder := NewFinder()
	r := NewReplacer(finder)

	r.SetReplacement("new text")

	if r.GetReplacement() != "new text" {
		t.Errorf("replacement = %q, want %q", r.GetReplacement(), "new text")
	}
}

func TestReplacer_CountMatches(t *testing.T) {
	finder := NewFinder()
	finder.SetPattern("test")

	r := NewReplacer(finder)

	buf := buffer.NewBuffer()
	buf.SetLines([]string{
		"test one",
		"test two",
		"no match",
	})

	count := r.CountMatches(buf)

	if count != 2 {
		t.Errorf("CountMatches = %d, want 2", count)
	}
}

func TestReplacer_CountMatches_EmptyPattern(t *testing.T) {
	finder := NewFinder()
	r := NewReplacer(finder)

	buf := buffer.NewBuffer()
	buf.SetLines([]string{"test"})

	count := r.CountMatches(buf)

	if count != 0 {
		t.Errorf("CountMatches with empty pattern = %d, want 0", count)
	}
}

func TestReplacer_IsPatternValid_Literal(t *testing.T) {
	finder := NewFinder()
	finder.SetPattern("test")

	r := NewReplacer(finder)

	if !r.IsPatternValid() {
		t.Error("literal pattern should always be valid")
	}
}

func TestReplacer_IsPatternValid_RegexValid(t *testing.T) {
	finder := NewFinder()
	finder.SetPattern("h.llo")
	opts := finder.GetOptions()
	opts.UseRegex = true
	finder.SetOptions(opts)

	r := NewReplacer(finder)

	if !r.IsPatternValid() {
		t.Error("valid regex pattern should be valid")
	}
}

func TestReplacer_IsPatternValid_RegexInvalid(t *testing.T) {
	finder := NewFinder()
	finder.SetPattern("[invalid")
	opts := finder.GetOptions()
	opts.UseRegex = true
	finder.SetOptions(opts)

	r := NewReplacer(finder)

	if r.IsPatternValid() {
		t.Error("invalid regex pattern should not be valid")
	}
}

func TestReplacer_ReplaceAll(t *testing.T) {
	finder := NewFinder()
	finder.SetPattern("test")

	r := NewReplacer(finder)
	r.SetReplacement("replaced")

	buf := buffer.NewBuffer()
	buf.SetLines([]string{
		"test one",
		"test two",
	})

	hist := history.NewHistory(100)
	count, err := r.ReplaceAll(buf, hist)

	if err != nil {
		t.Errorf("ReplaceAll error: %v", err)
	}

	if count != 2 {
		t.Errorf("ReplaceAll count = %d, want 2", count)
	}

	// Verify replacements
	lines := buf.GetAllLines()
	if len(lines) != 2 {
		t.Fatalf("buffer has %d lines, want 2", len(lines))
	}

	if lines[0] != "replaced one" {
		t.Errorf("line 0 = %q, want %q", lines[0], "replaced one")
	}

	if lines[1] != "replaced two" {
		t.Errorf("line 1 = %q, want %q", lines[1], "replaced two")
	}

	// Verify undo history was created
	if !hist.CanUndo() {
		t.Error("history should have undo available after ReplaceAll")
	}
}

func TestReplacer_ReplaceAll_EmptyPattern(t *testing.T) {
	finder := NewFinder()
	r := NewReplacer(finder)
	r.SetReplacement("new")

	buf := buffer.NewBuffer()
	buf.SetLines([]string{"test"})

	hist := history.NewHistory(100)
	count, err := r.ReplaceAll(buf, hist)

	if err != nil {
		t.Errorf("ReplaceAll error: %v", err)
	}

	if count != 0 {
		t.Errorf("ReplaceAll with empty pattern count = %d, want 0", count)
	}
}

func TestReplacer_ReplaceAll_NoMatches(t *testing.T) {
	finder := NewFinder()
	finder.SetPattern("notfound")

	r := NewReplacer(finder)
	r.SetReplacement("new")

	buf := buffer.NewBuffer()
	buf.SetLines([]string{"test"})

	hist := history.NewHistory(100)
	count, err := r.ReplaceAll(buf, hist)

	if err != nil {
		t.Errorf("ReplaceAll error: %v", err)
	}

	if count != 0 {
		t.Errorf("ReplaceAll with no matches count = %d, want 0", count)
	}
}

func TestReplacer_ReplaceAll_NilHistory(t *testing.T) {
	finder := NewFinder()
	finder.SetPattern("test")

	r := NewReplacer(finder)
	r.SetReplacement("replaced")

	buf := buffer.NewBuffer()
	buf.SetLines([]string{"test"})

	// Should not panic with nil history
	count, err := r.ReplaceAll(buf, nil)

	if err != nil {
		t.Errorf("ReplaceAll error: %v", err)
	}

	if count != 1 {
		t.Errorf("ReplaceAll count = %d, want 1", count)
	}
}

func TestReplacer_ValidateReplacement_Literal(t *testing.T) {
	finder := NewFinder()
	r := NewReplacer(finder)
	r.SetReplacement("hello world")

	if err := r.ValidateReplacement(); err != nil {
		t.Errorf("ValidateReplacement error: %v", err)
	}
}

func TestReplacer_ValidateReplacement_RegexValid(t *testing.T) {
	finder := NewFinder()
	opts := finder.GetOptions()
	opts.UseRegex = true
	finder.SetOptions(opts)

	r := NewReplacer(finder)
	r.SetReplacement("$1 is great") // Valid capture group reference

	if err := r.ValidateReplacement(); err != nil {
		t.Errorf("ValidateReplacement error: %v", err)
	}
}

func TestReplacer_ValidateReplacement_RegexInvalid(t *testing.T) {
	finder := NewFinder()
	opts := finder.GetOptions()
	opts.UseRegex = true
	finder.SetOptions(opts)

	r := NewReplacer(finder)
	r.SetReplacement("$a") // Invalid capture group reference

	if err := r.ValidateReplacement(); err == nil {
		t.Error("ValidateReplacement should error for invalid replacement")
	}
}

func TestReplacer_ValidateReplacement_RegexUnclosedGroup(t *testing.T) {
	finder := NewFinder()
	opts := finder.GetOptions()
	opts.UseRegex = true
	finder.SetOptions(opts)

	r := NewReplacer(finder)
	r.SetReplacement("${name") // Unclosed named group

	if err := r.ValidateReplacement(); err == nil {
		t.Error("ValidateReplacement should error for unclosed group reference")
	}
}

func TestReplacer_ReplaceCurrent(t *testing.T) {
	finder := NewFinder()
	finder.SetPattern("test")

	r := NewReplacer(finder)
	r.SetReplacement("replaced")

	buf := buffer.NewBuffer()
	buf.SetLines([]string{
		"test one",
		"test two",
	})

	// Find first match - from {-1, -1} to get the match at line 0
	finder.FindNext(buf, buffer.Position{Line: -1, Col: -1})

	hist := history.NewHistory(100)
	replaced, err := r.ReplaceCurrent(buf, hist)

	if err != nil {
		t.Errorf("ReplaceCurrent error: %v", err)
	}

	if !replaced {
		t.Error("ReplaceCurrent should have replaced")
	}

	// Verify first line was replaced
	lines := buf.GetAllLines()
	if lines[0] != "replaced one" {
		t.Errorf("line 0 = %q, want %q", lines[0], "replaced one")
	}

	// Second line should still have original text
	if lines[1] != "test two" {
		t.Errorf("line 1 = %q, want %q", lines[1], "test two")
	}
}

func TestReplacer_ReplaceCurrent_NoCurrentMatch(t *testing.T) {
	finder := NewFinder()
	finder.SetPattern("test")

	r := NewReplacer(finder)
	r.SetReplacement("replaced")

	buf := buffer.NewBuffer()
	buf.SetLines([]string{"test"})

	hist := history.NewHistory(100)
	replaced, err := r.ReplaceCurrent(buf, hist)

	if err != nil {
		t.Errorf("ReplaceCurrent error: %v", err)
	}

	if replaced {
		t.Error("ReplaceCurrent should not replace when no current match")
	}
}

func TestReplacer_ReplaceCurrent_NilHistory(t *testing.T) {
	finder := NewFinder()
	finder.SetPattern("test")

	r := NewReplacer(finder)
	r.SetReplacement("replaced")

	buf := buffer.NewBuffer()
	buf.SetLines([]string{"test"})

	// Find first match
	finder.FindNext(buf, buffer.Position{Line: 0, Col: 0})

	// Should not panic with nil history
	replaced, err := r.ReplaceCurrent(buf, nil)

	if err != nil {
		t.Errorf("ReplaceCurrent error: %v", err)
	}

	if !replaced {
		t.Error("ReplaceCurrent should have replaced")
	}
}

func TestReplacer_getReplacementText(t *testing.T) {
	finder := NewFinder()
	finder.SetPattern("test")

	r := NewReplacer(finder)
	r.SetReplacement("new text")

	match := Match{
		StartLine: 0,
		StartCol:  0,
		EndLine:   0,
		EndCol:    4,
		Text:      "test",
	}

	result := r.getReplacementText(match)

	if result != "new text" {
		t.Errorf("getReplacementText = %q, want %q", result, "new text")
	}
}

func TestReplacer_getReplacementText_Regex(t *testing.T) {
	finder := NewFinder()
	finder.SetPattern("h.llo")
	opts := finder.GetOptions()
	opts.UseRegex = true
	finder.SetOptions(opts)

	r := NewReplacer(finder)
	r.SetReplacement("greeting")

	match := Match{
		StartLine: 0,
		StartCol:  0,
		EndLine:   0,
		EndCol:    5,
		Text:      "hello",
	}

	result := r.getReplacementText(match)

	if result != "greeting" {
		t.Errorf("getReplacementText = %q, want %q", result, "greeting")
	}
}

func TestReplacer_ReplaceAll_MultipleLines(t *testing.T) {
	finder := NewFinder()
	finder.SetPattern("old")

	r := NewReplacer(finder)
	r.SetReplacement("new")

	buf := buffer.NewBuffer()
	buf.SetLines([]string{
		"old text on line 1",
		"old text on line 2",
		"old text on line 3",
	})

	hist := history.NewHistory(100)
	count, err := r.ReplaceAll(buf, hist)

	if err != nil {
		t.Errorf("ReplaceAll error: %v", err)
	}

	if count != 3 {
		t.Errorf("ReplaceAll count = %d, want 3", count)
	}

	// Verify all lines were replaced
	lines := buf.GetAllLines()
	for i, line := range lines {
		expected := "new text on line " + string('1'+byte(i))
		if line != expected {
			t.Errorf("line %d = %q, want %q", i, line, expected)
		}
	}
}

func TestReplacer_ReplaceAll_CaseSensitive(t *testing.T) {
	finder := NewFinder()
	finder.SetPattern("test")
	opts := finder.GetOptions()
	opts.CaseSensitive = true
	finder.SetOptions(opts)

	r := NewReplacer(finder)
	r.SetReplacement("replaced")

	buf := buffer.NewBuffer()
	buf.SetLines([]string{
		"test lowercase",
		"Test capitalized",
		"TEST uppercase",
	})

	hist := history.NewHistory(100)
	count, err := r.ReplaceAll(buf, hist)

	if err != nil {
		t.Errorf("ReplaceAll error: %v", err)
	}

	if count != 1 {
		t.Errorf("case-sensitive ReplaceAll count = %d, want 1", count)
	}

	lines := buf.GetAllLines()
	if lines[0] != "replaced lowercase" {
		t.Errorf("line 0 = %q, want %q", lines[0], "replaced lowercase")
	}
	// Other lines should be unchanged
	if lines[1] != "Test capitalized" {
		t.Errorf("line 1 = %q, want unchanged", lines[1])
	}
	if lines[2] != "TEST uppercase" {
		t.Errorf("line 2 = %q, want unchanged", lines[2])
	}
}
