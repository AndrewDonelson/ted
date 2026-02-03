// Package search implements search and replace functionality for the editor.
//
// It provides text searching with support for case sensitivity, whole word matching,
// and regular expressions. The finder maintains search state including current
// match position and search history.
package search

import (
	"regexp"
	"strings"

	"github.com/AndrewDonelson/ted/core/buffer"
)

// Match represents a single search match location.
type Match struct {
	StartLine int    // Line number (0-indexed)
	StartCol  int    // Column (0-indexed, byte offset)
	EndLine   int    // End line number
	EndCol    int    // End column
	Text      string // The matched text
}

// Options controls search behavior.
type Options struct {
	CaseSensitive bool // Match case exactly
	WholeWord     bool // Match whole words only
	UseRegex      bool // Treat pattern as regular expression
	WrapAround    bool // Wrap to start when reaching end
}

// DefaultOptions returns the default search options.
func DefaultOptions() Options {
	return Options{
		CaseSensitive: false,
		WholeWord:     false,
		UseRegex:      false,
		WrapAround:    true,
	}
}

// Finder manages search operations and state.
type Finder struct {
	pattern      string   // Current search pattern
	options      Options  // Search options
	matches      []Match  // All matches in current search
	currentIndex int      // Index of current match
	history      []string // Search history
	historyIndex int      // Current position in history
	maxHistory   int      // Maximum history entries
}

// NewFinder creates a new search finder.
func NewFinder() *Finder {
	return &Finder{
		matches:    make([]Match, 0),
		history:    make([]string, 0, 20),
		maxHistory: 20,
		options:    DefaultOptions(),
	}
}

// SetPattern sets the search pattern and clears previous matches.
func (f *Finder) SetPattern(pattern string) {
	if pattern == "" {
		f.pattern = ""
		f.matches = f.matches[:0]
		return
	}

	// Only add to history if different from last search
	if pattern != f.pattern && pattern != "" {
		f.addToHistory(pattern)
	}

	f.pattern = pattern
	f.matches = f.matches[:0]
	f.currentIndex = -1
}

// GetPattern returns the current search pattern.
func (f *Finder) GetPattern() string {
	return f.pattern
}

// SetOptions sets the search options.
func (f *Finder) SetOptions(options Options) {
	f.options = options
	// Clear matches since options changed
	f.matches = f.matches[:0]
	f.currentIndex = -1
}

// GetOptions returns the current search options.
func (f *Finder) GetOptions() Options {
	return f.options
}

// addToHistory adds a pattern to the search history.
func (f *Finder) addToHistory(pattern string) {
	// Check if pattern is already at the end of history
	if len(f.history) > 0 && f.history[len(f.history)-1] == pattern {
		return
	}

	// Remove oldest if at max capacity
	if len(f.history) >= f.maxHistory {
		f.history = append(f.history[:0], f.history[1:]...)
	}

	f.history = append(f.history, pattern)
	f.historyIndex = len(f.history) - 1
}

// GetHistory returns the search history.
func (f *Finder) GetHistory() []string {
	result := make([]string, len(f.history))
	copy(result, f.history)
	return result
}

// GetHistoryItem returns a specific history item by index.
func (f *Finder) GetHistoryItem(index int) (string, bool) {
	if index < 0 || index >= len(f.history) {
		return "", false
	}
	return f.history[index], true
}

// PreviousHistory moves to the previous history entry.
func (f *Finder) PreviousHistory() (string, bool) {
	if f.historyIndex > 0 {
		f.historyIndex--
		return f.history[f.historyIndex], true
	}
	return "", false
}

// NextHistory moves to the next history entry.
func (f *Finder) NextHistory() (string, bool) {
	if f.historyIndex < len(f.history)-1 {
		f.historyIndex++
		return f.history[f.historyIndex], true
	}
	return "", false
}

// FindAll finds all matches in the buffer.
func (f *Finder) FindAll(buf *buffer.Buffer) []Match {
	if f.pattern == "" {
		return nil
	}

	f.matches = f.matches[:0]
	lines := buf.GetAllLines()

	if f.options.UseRegex {
		f.findAllRegex(lines)
	} else {
		f.findAllLiteral(lines)
	}

	return f.matches
}

// findAllLiteral finds all literal pattern matches.
func (f *Finder) findAllLiteral(lines []string) {
	pattern := f.pattern
	if !f.options.CaseSensitive {
		pattern = strings.ToLower(pattern)
	}

	for lineNum, line := range lines {
		searchLine := line
		if !f.options.CaseSensitive {
			searchLine = strings.ToLower(line)
		}

		startCol := 0
		for {
			idx := strings.Index(searchLine[startCol:], pattern)
			if idx == -1 {
				break
			}

			actualIdx := startCol + idx

			// Check whole word constraint
			if f.options.WholeWord && !f.isWholeWordMatch(line, actualIdx, len(pattern)) {
				startCol = actualIdx + 1
				continue
			}

			match := Match{
				StartLine: lineNum,
				StartCol:  actualIdx,
				EndLine:   lineNum,
				EndCol:    actualIdx + len(pattern),
				Text:      line[actualIdx : actualIdx+len(pattern)],
			}
			f.matches = append(f.matches, match)

			startCol = actualIdx + 1
		}
	}
}

// findAllRegex finds all regex pattern matches.
func (f *Finder) findAllRegex(lines []string) {
	var re *regexp.Regexp
	var err error

	if f.options.CaseSensitive {
		re, err = regexp.Compile(f.pattern)
	} else {
		re, err = regexp.Compile("(?i)" + f.pattern)
	}

	if err != nil {
		// Invalid regex, no matches
		return
	}

	for lineNum, line := range lines {
		matches := re.FindAllStringIndex(line, -1)
		for _, m := range matches {
			// m[0] is start index, m[1] is end index
			match := Match{
				StartLine: lineNum,
				StartCol:  m[0],
				EndLine:   lineNum,
				EndCol:    m[1],
				Text:      line[m[0]:m[1]],
			}
			f.matches = append(f.matches, match)
		}
	}
}

// isWholeWordMatch checks if a match is a whole word.
func (f *Finder) isWholeWordMatch(line string, start, length int) bool {
	// Check character before
	if start > 0 && isWordChar(line[start-1]) {
		return false
	}

	// Check character after
	end := start + length
	if end < len(line) && isWordChar(line[end]) {
		return false
	}

	return true
}

// isWordChar returns true if the byte is a word character.
func isWordChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

// FindNext finds the next match from the current position.
// Returns the match and true if found, otherwise returns false.
func (f *Finder) FindNext(buf *buffer.Buffer, fromPos buffer.Position) (Match, bool) {
	if f.pattern == "" {
		return Match{}, false
	}

	// Refresh matches if needed
	if len(f.matches) == 0 {
		f.FindAll(buf)
	}

	if len(f.matches) == 0 {
		return Match{}, false
	}

	// Find first match after fromPos
	for i, match := range f.matches {
		if match.StartLine > fromPos.Line ||
			(match.StartLine == fromPos.Line && match.StartCol > fromPos.Col) {
			f.currentIndex = i
			return match, true
		}
	}

	// Wrap around if enabled
	if f.options.WrapAround && len(f.matches) > 0 {
		f.currentIndex = 0
		return f.matches[0], true
	}

	return Match{}, false
}

// FindPrevious finds the previous match from the current position.
// Returns the match and true if found, otherwise returns false.
func (f *Finder) FindPrevious(buf *buffer.Buffer, fromPos buffer.Position) (Match, bool) {
	if f.pattern == "" {
		return Match{}, false
	}

	// Refresh matches if needed
	if len(f.matches) == 0 {
		f.FindAll(buf)
	}

	if len(f.matches) == 0 {
		return Match{}, false
	}

	// Find last match before fromPos
	for i := len(f.matches) - 1; i >= 0; i-- {
		match := f.matches[i]
		if match.StartLine < fromPos.Line ||
			(match.StartLine == fromPos.Line && match.StartCol < fromPos.Col) {
			f.currentIndex = i
			return match, true
		}
	}

	// Wrap around if enabled
	if f.options.WrapAround && len(f.matches) > 0 {
		f.currentIndex = len(f.matches) - 1
		return f.matches[len(f.matches)-1], true
	}

	return Match{}, false
}

// GetCurrentMatch returns the current match if any.
func (f *Finder) GetCurrentMatch() (Match, bool) {
	if f.currentIndex < 0 || f.currentIndex >= len(f.matches) {
		return Match{}, false
	}
	return f.matches[f.currentIndex], true
}

// SetCurrentMatch sets the current match by index.
func (f *Finder) SetCurrentMatch(index int) bool {
	if index < 0 || index >= len(f.matches) {
		return false
	}
	f.currentIndex = index
	return true
}

// GetMatchCount returns the total number of matches.
func (f *Finder) GetMatchCount() int {
	return len(f.matches)
}

// GetCurrentMatchIndex returns the index of the current match.
func (f *Finder) GetCurrentMatchIndex() int {
	return f.currentIndex
}

// Clear clears all matches and resets the finder state.
func (f *Finder) Clear() {
	f.matches = f.matches[:0]
	f.currentIndex = -1
}

// Reset clears the finder completely including pattern and history.
func (f *Finder) Reset() {
	f.pattern = ""
	f.matches = f.matches[:0]
	f.currentIndex = -1
	f.options = DefaultOptions()
}
