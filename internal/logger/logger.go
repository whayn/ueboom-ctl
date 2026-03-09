package logger

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var (
	infoColor    = color.New(color.FgCyan)
	warnColor    = color.New(color.FgYellow)
	errorColor   = color.New(color.FgRed, color.Bold)
	successColor = color.New(color.FgGreen)
	dimColor     = color.New(color.FgHiBlack) // Dim/Gray for Metadata
	accentColor  = color.New(color.FgCyan, color.Bold)
)

// Highlight returns a string formatted with the accent color.
func Highlight(s string) string {
	return accentColor.Sprint(s)
}

// Info prints a formatted info message to stderr.
func Info(format string, a ...interface{}) {
	infoColor.Fprint(os.Stderr, "INFO ")
	fmt.Fprintf(os.Stderr, format+"\n", a...)
}

// Warn prints a formatted warning message to stderr.
func Warn(format string, a ...interface{}) {
	warnColor.Fprint(os.Stderr, "WARN ")
	fmt.Fprintf(os.Stderr, format+"\n", a...)
}

// Error prints a formatted error message to stderr.
func Error(format string, a ...interface{}) {
	errorColor.Fprint(os.Stderr, "FAIL ")
	fmt.Fprintf(os.Stderr, format+"\n", a...)
}

// Success prints a success message to stderr.
func Success(format string, a ...interface{}) {
	successColor.Fprint(os.Stderr, "DONE ")
	fmt.Fprintf(os.Stderr, format+"\n", a...)
}

// Metadata prints a dimmed message to stderr.
func Metadata(format string, a ...interface{}) {
	dimColor.Fprintf(os.Stderr, format+"\n", a...)
}

// Section prints a bold header for grouping information.
func Section(title string) {
	fmt.Fprintln(os.Stderr)
	color.New(color.Bold, color.Underline).Fprintln(os.Stderr, title)
}

// List prints an item in a list.
func List(index int, label string, metadata string) {
	fmt.Fprintf(os.Stderr, " [%d] ", index)
	fmt.Fprint(os.Stderr, label)
	if metadata != "" {
		fmt.Fprint(os.Stderr, " ")
		dimColor.Fprintf(os.Stderr, "(%s)", metadata)
	}
	fmt.Fprintln(os.Stderr)
}
