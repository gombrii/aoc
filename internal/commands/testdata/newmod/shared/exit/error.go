// Package exit does, in the spirit of Advent of Code, provide quick and dirty ways to fail.
package exit

import (
	"fmt"
	"os"
	"runtime"
)

// If prints err and exits if err != nil.
func If(err error) {
	if err != nil {
		if _, _, line, ok := runtime.Caller(1); ok {
			fmt.Printf("Error:%d: %v\n", line, err)
		} else {
			fmt.Printf("Error: %v\n", err)
		}
		os.Exit(1)
	}
}

// PanicIf panics if err != nil.
func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}
