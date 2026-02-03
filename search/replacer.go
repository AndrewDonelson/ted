// Package search implements search and replace functionality for the editor.
package search

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/AndrewDonelson/ted/core/buffer"
	"github.com/AndrewDonelson/ted/core/history"
)

// ReplaceResult represents the result of a replace operation.
type ReplaceResult struct {
	ReplacedCount int // Number of replacements made
	SkippedCount  int // Number of matches skipped (for replace current)
}

// Replacer performs replace operations using a Finder.
type Replacer struct {
	finder      *Finder
	replacement string
}

// NewReplacer creates a new replacer with the given finder.
func NewReplacer(finder *Finder) *Replacer {
	return &Replacer{
		finder: finder,
	}
}

// SetReplacement sets the replacement string.
func (r *Replacer) SetReplacement(replacement string) {
	r.replacement = replacement
}

// GetReplacement returns the current replacement string.
func (r *Replacer) GetReplacement() string {
	return r.replacement
}

// ReplaceCurrent replaces the current match and advances to the next.
// Returns true if a replacement was made.
func (r *Replacer) ReplaceCurrent(buf *buffer.Buffer, hist *history.History) (bool, error) {
	match, ok := r.finder.GetCurrentMatch()
	if !ok {
		return false, nil
	}

	// Get the replacement text
	replacement := r.getReplacementText(match)

	// Record for undo
	deletedText, _ := buf.GetText(
		buffer.Position{Line: match.StartLine, Col: match.StartCol},
		buffer.Position{Line: match.EndLine, Col: match.EndCol},
	)

	if hist != nil {
		op := &history.DeleteOperation{
			StartPos: buffer.Position{Line: match.StartLine, Col: match.StartCol},
			EndPos:   buffer.Position{Line: match.EndLine, Col: match.EndCol},
			Deleted:  deletedText,
		}
		hist.Push(op)
	}

	// Delete the match
	if err := buf.Delete(
		buffer.Position{Line: match.StartLine, Col: match.StartCol},
		buffer.Position{Line: match.EndLine, Col: match.EndCol},
	); err != nil {
		return false, fmt.Errorf("delete match: %w", err)
	}

	// Insert replacement
	if err := buf.Insert(
		buffer.Position{Line: match.StartLine, Col: match.StartCol},
		replacement,
	); err != nil {
		return false, fmt.Errorf("insert replacement: %w", err)
	}

	// Record insert for undo
	if hist != nil {
		insertOp := &history.InsertOperation{
			Pos:  buffer.Position{Line: match.StartLine, Col: match.StartCol},
			Text: replacement,
		}
		hist.Push(insertOp)
	}

	// Clear matches and refind - positions may have changed
	r.finder.Clear()

	return true, nil
}

// ReplaceAll replaces all matches in the buffer.
// Returns the number of replacements made.
func (r *Replacer) ReplaceAll(buf *buffer.Buffer, hist *history.History) (int, error) {
	if r.finder.GetPattern() == "" {
		return 0, nil
	}

	// Find all matches
	matches := r.finder.FindAll(buf)
	if len(matches) == 0 {
		return 0, nil
	}

	// Create composite operation for undo
	compOp := &history.CompositeOperation{}
	compOp.SetDescription(fmt.Sprintf("replace all '%s' with '%s'", r.finder.GetPattern(), r.replacement))

	// Replace from end to beginning to avoid position shifting
	replaceCount := 0
	for i := len(matches) - 1; i >= 0; i-- {
		match := matches[i]
		replacement := r.getReplacementText(match)

		// Get deleted text for undo
		deletedText, _ := buf.GetText(
			buffer.Position{Line: match.StartLine, Col: match.StartCol},
			buffer.Position{Line: match.EndLine, Col: match.EndCol},
		)

		// Record delete operation
		op := &history.DeleteOperation{
			StartPos: buffer.Position{Line: match.StartLine, Col: match.StartCol},
			EndPos:   buffer.Position{Line: match.EndLine, Col: match.EndCol},
			Deleted:  deletedText,
		}
		compOp.Operations = append(compOp.Operations, op)

		// Delete the match
		if err := buf.Delete(
			buffer.Position{Line: match.StartLine, Col: match.StartCol},
			buffer.Position{Line: match.EndLine, Col: match.EndCol},
		); err != nil {
			return replaceCount, fmt.Errorf("delete match: %w", err)
		}

		// Insert replacement
		if err := buf.Insert(
			buffer.Position{Line: match.StartLine, Col: match.StartCol},
			replacement,
		); err != nil {
			return replaceCount, fmt.Errorf("insert replacement: %w", err)
		}

		// Record insert operation
		insertOp := &history.InsertOperation{
			Pos:  buffer.Position{Line: match.StartLine, Col: match.StartCol},
			Text: replacement,
		}
		compOp.Operations = append([]history.Operation{insertOp}, compOp.Operations...)

		replaceCount++
	}

	// Push composite operation to history
	if len(compOp.Operations) > 0 && hist != nil {
		hist.Push(compOp)
	}

	// Clear finder state since we changed the buffer
	r.finder.Clear()

	return replaceCount, nil
}

// getReplacementText returns the actual replacement text for a match.
// If using regex, this processes capture groups.
func (r *Replacer) getReplacementText(match Match) string {
	if !r.finder.options.UseRegex {
		return r.replacement
	}

	// Process regex replacement (handle $1, $2, etc.)
	return r.processRegexReplacement(match)
}

// processRegexReplacement processes regex capture group references.
func (r *Replacer) processRegexReplacement(match Match) string {
	result := r.replacement

	// Simple implementation: replace $1, $2, etc. with capture groups
	// In a full implementation, you'd want to parse the regex and extract groups
	// For now, we just return the replacement as-is

	// TODO: Implement full capture group replacement
	// This requires parsing the regex and extracting submatches

	return result
}

// ValidateReplacement validates that the replacement string is valid.
// For regex mode, this checks that capture group references are well-formed.
func (r *Replacer) ValidateReplacement() error {
	if !r.finder.options.UseRegex {
		return nil
	}

	// Check for invalid capture group references
	// $0, $1, $2, etc. are valid
	// $$ escapes a literal $

	for i := 0; i < len(r.replacement); i++ {
		if r.replacement[i] == '$' {
			if i+1 >= len(r.replacement) {
				return fmt.Errorf("incomplete escape at end of replacement")
			}
			next := r.replacement[i+1]
			if next == '$' || next == '&' || next == '`' || next == '\'' {
				// Valid escape
				i++
				continue
			}
			if next >= '0' && next <= '9' {
				// Capture group reference - valid
				i++
				continue
			}
			if next == '{' {
				// Named group reference - simplified, just check it's closed
				j := i + 2
				for j < len(r.replacement) && r.replacement[j] != '}' {
					j++
				}
				if j >= len(r.replacement) {
					return fmt.Errorf("unclosed named group reference")
				}
				i = j
				continue
			}
			return fmt.Errorf("invalid escape sequence: $%c", next)
		}
	}

	return nil
}

// CountMatches returns the number of matches in the buffer.
func (r *Replacer) CountMatches(buf *buffer.Buffer) int {
	if r.finder.GetPattern() == "" {
		return 0
	}

	matches := r.finder.FindAll(buf)
	return len(matches)
}

// IsPatternValid checks if the current search pattern is valid.
// For regex mode, this validates the regular expression.
func (r *Replacer) IsPatternValid() bool {
	if !r.finder.options.UseRegex {
		return true
	}

	pattern := r.finder.GetPattern()
	if pattern == "" {
		return true
	}

	var err error
	if r.finder.options.CaseSensitive {
		_, err = regexp.Compile(pattern)
	} else {
		_, err = regexp.Compile("(?i)" + pattern)
	}

	return err == nil
}

// EscapeLiteral escapes special regex characters for literal search.
func EscapeLiteral(s string) string {
	specialChars := []string{"\\", ".", "*", "+", "?", "(", ")", "[", "]", "{", "}", "^", "$", "|"}
	result := s
	for _, char := range specialChars {
		result = strings.ReplaceAll(result, char, "\\"+char)
	}
	return result
}
