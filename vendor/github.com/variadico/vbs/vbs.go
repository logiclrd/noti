// Package vbs prints text to stdout if Verbose is set to true. This package
// isn't threadsafe. Use a mutex if you're going to be changing Verbose in
// goroutines.
package vbs

import (
	"fmt"
	"io"
	"os"
)

// Verbose indicates whether or not something should be printed.
var Verbose bool

// Println prints to stdout if Verbose is true.
func Println(a ...interface{}) {
	if Verbose {
		fmt.Println(a...)
	}
}

// Printf prints to stdout if Verbose is true.
func Printf(format string, a ...interface{}) {
	if Verbose {
		fmt.Printf(format, a...)
	}
}

// Printer is a conditional printer that writes to Output.
type Printer struct {
	Verbose bool
	Output  io.Writer
}

// New returns a new verbose Printer.
func New() Printer {
	return Printer{
		Verbose: false,
		Output:  os.Stdout,
	}
}

// Println prints to stdout if Verbose is true.
func (p Printer) Println(a ...interface{}) {
	if p.Verbose {
		fmt.Fprintln(p.Output, a...)
	}
}

// Printf prints to stdout if Verbose is true.
func (p Printer) Printf(format string, a ...interface{}) {
	if p.Verbose {
		fmt.Fprintf(p.Output, format, a...)
	}
}
