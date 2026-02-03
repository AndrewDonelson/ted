// Package dialog implements search and replace dialogs for the editor.
package dialog

import (
	"fmt"
	"strings"

	"github.com/AndrewDonelson/ted/core/buffer"
	"github.com/AndrewDonelson/ted/search"
	"github.com/gdamore/tcell/v2"
)

// SearchDialog is a dialog for find and replace operations.
type SearchDialog struct {
	BaseDialog
	finder        *search.Finder
	replacer      *search.Replacer
	searchInput   string
	replaceInput  string
	message       string
	isReplaceMode bool
	showOptions   bool
	options       search.Options
	onFind        func()
	onFindNext    func()
	onReplace     func()
	onReplaceAll  func()
	onCancel      func()
}

// NewSearchDialog creates a new search dialog.
func NewSearchDialog(finder *search.Finder, replacer *search.Replacer, isReplace bool, onCancel func()) *SearchDialog {
	width := 50
	height := 8
	if isReplace {
		height = 10
	}

	d := &SearchDialog{
		BaseDialog: BaseDialog{
			title:  "Find",
			width:  width,
			height: height,
		},
		finder:        finder,
		replacer:      replacer,
		isReplaceMode: isReplace,
		showOptions:   true,
		onCancel:      onCancel,
	}

	if isReplace {
		d.title = "Replace"
	}

	// Load current values from finder/replacer
	if finder != nil {
		d.searchInput = finder.GetPattern()
		d.options = finder.GetOptions()
	}
	if replacer != nil {
		d.replaceInput = replacer.GetReplacement()
	}

	return d
}

// HandleInput processes keyboard input for the search dialog.
func (d *SearchDialog) HandleInput(key tcell.Key, mod tcell.ModMask, ch rune) bool {
	switch key {
	case tcell.KeyEscape:
		d.SetCancelled()
		if d.onCancel != nil {
			d.onCancel()
		}
		return true

	case tcell.KeyTab:
		// Cycle through input fields and buttons
		d.cycleFocus()
		return true

	case tcell.KeyBackspace, tcell.KeyBackspace2:
		d.handleBackspace()
		return true

	case tcell.KeyDelete:
		d.handleDelete()
		return true

	case tcell.KeyLeft:
		d.handleCursorLeft()
		return true

	case tcell.KeyRight:
		d.handleCursorRight()
		return true

	case tcell.KeyHome:
		d.handleCursorHome()
		return true

	case tcell.KeyEnd:
		d.handleCursorEnd()
		return true

	case tcell.KeyEnter:
		return d.handleEnter()

	case tcell.KeyF3:
		// Find Next shortcut
		if d.onFindNext != nil {
			d.onFindNext()
			return true
		}
		return false

	case tcell.KeyRune:
		if ch != 0 {
			d.handleCharacter(ch)
			return true
		}
	}

	return false
}

// cycleFocus moves focus to the next input element.
func (d *SearchDialog) cycleFocus() {
	maxFocus := 2 // Find Next, Replace, Replace All buttons
	if d.isReplaceMode {
		maxFocus = 4 // + search field, replace field
	}

	d.focusIndex = (d.focusIndex + 1) % (maxFocus + 1)
}

// handleBackspace handles backspace key.
func (d *SearchDialog) handleBackspace() {
	if d.focusIndex == 0 {
		if len(d.searchInput) > 0 {
			d.searchInput = d.searchInput[:len(d.searchInput)-1]
		}
	} else if d.isReplaceMode && d.focusIndex == 1 {
		if len(d.replaceInput) > 0 {
			d.replaceInput = d.replaceInput[:len(d.replaceInput)-1]
		}
	}
}

// handleDelete handles delete key.
func (d *SearchDialog) handleDelete() {
	// In a full implementation, this would delete at cursor position
	// For now, treat same as backspace at end
	d.handleBackspace()
}

// handleCursorLeft moves cursor left in current field.
func (d *SearchDialog) handleCursorLeft() {
	// Simplified: for now just cycle focus
	d.cycleFocus()
}

// handleCursorRight moves cursor right in current field.
func (d *SearchDialog) handleCursorRight() {
	// Simplified: for now just cycle focus
	d.cycleFocus()
}

// handleCursorHome moves cursor to start of field.
func (d *SearchDialog) handleCursorHome() {
	// Simplified
}

// handleCursorEnd moves cursor to end of field.
func (d *SearchDialog) handleCursorEnd() {
	// Simplified
}

// handleEnter processes the Enter key based on current focus.
func (d *SearchDialog) handleEnter() bool {
	// Update finder/replacer with current values
	if d.finder != nil {
		d.finder.SetPattern(d.searchInput)
		d.finder.SetOptions(d.options)
	}
	if d.replacer != nil {
		d.replacer.SetReplacement(d.replaceInput)
	}

	switch d.focusIndex {
	case 0, 1, 2: // Find Next button or search field
		if d.onFindNext != nil {
			d.onFindNext()
		}
		return true
	case 3: // Replace button
		if d.isReplaceMode && d.onReplace != nil {
			d.onReplace()
		}
		return true
	case 4: // Replace All button
		if d.isReplaceMode && d.onReplaceAll != nil {
			d.onReplaceAll()
		}
		return true
	}

	return false
}

// handleCharacter processes a typed character.
func (d *SearchDialog) handleCharacter(ch rune) {
	if d.focusIndex == 0 {
		d.searchInput += string(ch)
	} else if d.isReplaceMode && d.focusIndex == 1 {
		d.replaceInput += string(ch)
	}
}

// toggleOption toggles a search option.
func (d *SearchDialog) toggleOption(option string) {
	switch option {
	case "case":
		d.options.CaseSensitive = !d.options.CaseSensitive
	case "word":
		d.options.WholeWord = !d.options.WholeWord
	case "regex":
		d.options.UseRegex = !d.options.UseRegex
	}

	if d.finder != nil {
		d.finder.SetOptions(d.options)
	}
}

// Render draws the search dialog.
func (d *SearchDialog) Render(screen Screen, style tcell.Style) {
	if !d.isOpen {
		return
	}

	// Clear dialog area
	d.Clear(screen, style)

	// Draw border
	d.DrawBorder(screen, style)

	currentY := d.y + 2

	// Draw search field label and input
	d.DrawText(screen, d.x+2, currentY, "Find:", style)
	currentY++

	// Draw search input field
	searchStyle := style
	if d.focusIndex == 0 {
		searchStyle = style.Reverse(true)
	}
	d.DrawText(screen, d.x+2, currentY, d.searchInput+"█", searchStyle)
	currentY += 2

	// Draw replace field if in replace mode
	if d.isReplaceMode {
		d.DrawText(screen, d.x+2, currentY, "Replace:", style)
		currentY++

		replaceStyle := style
		if d.focusIndex == 1 {
			replaceStyle = style.Reverse(true)
		}
		d.DrawText(screen, d.x+2, currentY, d.replaceInput+"█", replaceStyle)
		currentY += 2
	}

	// Draw options
	if d.showOptions {
		optionsText := d.buildOptionsText()
		d.DrawText(screen, d.x+2, currentY, optionsText, style)
		currentY++
	}

	// Draw message if any
	if d.message != "" {
		msgStyle := style.Foreground(tcell.ColorYellow)
		d.DrawText(screen, d.x+2, currentY, d.message, msgStyle)
		currentY++
	}

	currentY++

	// Draw buttons
	buttonY := currentY
	buttonSpacing := (d.width - 40) / 4 // Distribute buttons evenly

	// Calculate button positions
	btnX := d.x + buttonSpacing

	// Find Next button
	findNextStyle := style
	if d.focusIndex == 2 {
		findNextStyle = style.Reverse(true).Bold(true)
	}
	d.DrawButton(screen, btnX, buttonY, 0, "Find Next", findNextStyle, d.focusIndex == 2)

	if d.isReplaceMode {
		btnX += buttonSpacing + 10

		// Replace button
		replaceStyle := style
		if d.focusIndex == 3 {
			replaceStyle = style.Reverse(true).Bold(true)
		}
		d.DrawButton(screen, btnX, buttonY, 0, "Replace", replaceStyle, d.focusIndex == 3)

		btnX += buttonSpacing + 10

		// Replace All button
		replaceAllStyle := style
		if d.focusIndex == 4 {
			replaceAllStyle = style.Reverse(true).Bold(true)
		}
		d.DrawButton(screen, btnX, buttonY, 0, "Replace All", replaceAllStyle, d.focusIndex == 4)
	}
}

// buildOptionsText builds the options display text.
func (d *SearchDialog) buildOptionsText() string {
	var parts []string

	if d.options.CaseSensitive {
		parts = append(parts, "[✓] Case")
	} else {
		parts = append(parts, "[ ] Case")
	}

	if d.options.WholeWord {
		parts = append(parts, "[✓] Word")
	} else {
		parts = append(parts, "[ ] Word")
	}

	if d.options.UseRegex {
		parts = append(parts, "[✓] Regex")
	} else {
		parts = append(parts, "[ ] Regex")
	}

	return strings.Join(parts, "  ")
}

// SetMessage sets the status message displayed in the dialog.
func (d *SearchDialog) SetMessage(msg string) {
	d.message = msg
}

// GetSearchInput returns the current search input.
func (d *SearchDialog) GetSearchInput() string {
	return d.searchInput
}

// SetSearchInput sets the search input.
func (d *SearchDialog) SetSearchInput(input string) {
	d.searchInput = input
}

// GetReplaceInput returns the current replace input.
func (d *SearchDialog) GetReplaceInput() string {
	return d.replaceInput
}

// SetReplaceInput sets the replace input.
func (d *SearchDialog) SetReplaceInput(input string) {
	d.replaceInput = input
}

// SetOnFindNext sets the callback for Find Next.
func (d *SearchDialog) SetOnFindNext(fn func()) {
	d.onFindNext = fn
}

// SetOnReplace sets the callback for Replace.
func (d *SearchDialog) SetOnReplace(fn func()) {
	d.onReplace = fn
}

// SetOnReplaceAll sets the callback for Replace All.
func (d *SearchDialog) SetOnReplaceAll(fn func()) {
	d.onReplaceAll = fn
}

// GetResult returns nil for search dialog (use callbacks).
func (d *SearchDialog) GetResult() interface{} {
	return nil
}

// FindDialog is a convenience wrapper for find-only mode.
type FindDialog struct {
	*SearchDialog
}

// NewFindDialog creates a new find dialog.
func NewFindDialog(finder *search.Finder, onFindNext func(), onCancel func()) *FindDialog {
	// Create a dummy replacer since we don't need it for find
	dummyReplacer := search.NewReplacer(finder)

	searchDlg := NewSearchDialog(finder, dummyReplacer, false, onCancel)
	searchDlg.SetOnFindNext(onFindNext)

	return &FindDialog{
		SearchDialog: searchDlg,
	}
}

// ReplaceDialog is a convenience wrapper for replace mode.
type ReplaceDialog struct {
	*SearchDialog
}

// NewReplaceDialog creates a new replace dialog.
func NewReplaceDialog(finder *search.Finder, replacer *search.Replacer, onReplace func(), onReplaceAll func(), onCancel func()) *ReplaceDialog {
	searchDlg := NewSearchDialog(finder, replacer, true, onCancel)
	searchDlg.SetOnReplace(onReplace)
	searchDlg.SetOnReplaceAll(onReplaceAll)

	return &ReplaceDialog{
		SearchDialog: searchDlg,
	}
}

// SearchManager manages search state and operations.
type SearchManager struct {
	finder   *search.Finder
	replacer *search.Replacer
	history  []string
}

// NewSearchManager creates a new search manager.
func NewSearchManager() *SearchManager {
	finder := search.NewFinder()
	return &SearchManager{
		finder:   finder,
		replacer: search.NewReplacer(finder),
		history:  make([]string, 0, 20),
	}
}

// GetFinder returns the search finder.
func (sm *SearchManager) GetFinder() *search.Finder {
	return sm.finder
}

// GetReplacer returns the search replacer.
func (sm *SearchManager) GetReplacer() *search.Replacer {
	return sm.replacer
}

// FindNext finds the next match and moves the cursor there.
func (sm *SearchManager) FindNext(buf *buffer.Buffer, startPos buffer.Position) (*search.Match, bool) {
	match, found := sm.finder.FindNext(buf, startPos)
	if found {
		buf.MoveCursor(buffer.Position{Line: match.StartLine, Col: match.StartCol})
	}
	return &match, found
}

// FindPrevious finds the previous match and moves the cursor there.
func (sm *SearchManager) FindPrevious(buf *buffer.Buffer, startPos buffer.Position) (*search.Match, bool) {
	match, found := sm.finder.FindPrevious(buf, startPos)
	if found {
		buf.MoveCursor(buffer.Position{Line: match.StartLine, Col: match.StartCol})
	}
	return &match, found
}

// SetPattern sets the search pattern.
func (sm *SearchManager) SetPattern(pattern string) {
	sm.finder.SetPattern(pattern)
}

// GetPattern returns the current search pattern.
func (sm *SearchManager) GetPattern() string {
	return sm.finder.GetPattern()
}

// SetReplacement sets the replacement string.
func (sm *SearchManager) SetReplacement(replacement string) {
	sm.replacer.SetReplacement(replacement)
}

// GetReplacement returns the current replacement string.
func (sm *SearchManager) GetReplacement() string {
	return sm.replacer.GetReplacement()
}

// SetOptions sets the search options.
func (sm *SearchManager) SetOptions(options search.Options) {
	sm.finder.SetOptions(options)
}

// GetOptions returns the current search options.
func (sm *SearchManager) GetOptions() search.Options {
	return sm.finder.GetOptions()
}

// GetMatchCount returns the number of matches.
func (sm *SearchManager) GetMatchCount() int {
	return sm.finder.GetMatchCount()
}

// GetCurrentMatch returns the current match index and total.
func (sm *SearchManager) GetCurrentMatch() (int, int) {
	return sm.finder.GetCurrentMatchIndex(), sm.finder.GetMatchCount()
}

// BuildStatusMessage builds a status message for the info bar.
func (sm *SearchManager) BuildStatusMessage() string {
	current, total := sm.GetCurrentMatch()
	if total == 0 {
		if sm.GetPattern() != "" {
			return "No matches"
		}
		return ""
	}
	return fmt.Sprintf("Match %d of %d", current+1, total)
}
