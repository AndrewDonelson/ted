// Package main is the entry point for ted (Terminal EDitor).
//
// ted is a modern, cross-platform command-line text editor written in Go
// that uses familiar Windows-style keyboard shortcuts and intuitive arrow key navigation.
package main

import (
	"fmt"
	"os"

	"github.com/AndrewDonelson/ted/editor"
)

func main() {
	// Parse command-line arguments
	var filePath string
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}

	// Create editor
	ed, err := editor.NewEditor()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing editor: %v\n", err)
		os.Exit(1)
	}

	// Open file if provided
	if filePath != "" {
		if err := ed.OpenFile(filePath); err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file %q: %v\n", filePath, err)
			// Continue anyway - allow editing new file
		}
	}

	// Run editor
	if err := ed.Run(); err != nil {
		if err == editor.ErrQuit {
			// Normal quit
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "Error running editor: %v\n", err)
		os.Exit(1)
	}
}
