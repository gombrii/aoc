// Package exit does, in the spirit of Advent of Code, provide quick and dirty ways to fail.
package exit

import (
	"fmt"
	"os"
)

// If prints err and exits if err != nil.
func If(err error) {
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

// PanicIf panics if err != nil.
func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}
