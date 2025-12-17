// Package clipboard provides cross-platform clipboard operations.
package clipboard

import (
	"github.com/atotto/clipboard"
)

// Read reads text from the system clipboard.
func Read() (string, error) {
	return clipboard.ReadAll()
}

// Write writes text to the system clipboard.
func Write(text string) error {
	return clipboard.WriteAll(text)
}
