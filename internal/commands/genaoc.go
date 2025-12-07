package commands

import (
	"fmt"
	"path/filepath"

	"github.com/gombrii/aoc/internal/exec"
	"github.com/gombrii/aoc/internal/files"
)

const exitTmpl = `// Package exit does, in the spirit of Advent of Code, provide quick and dirty ways to fail.
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
			fmt.Printf("Error:%d: %v", line, err)
		} else {
			fmt.Printf("Error: %v", err)
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
`

const parseTmpl = `// Package input does, in the spirit of Advent of Code, provide some common ways to interpret the
// puzzles' input data.
package parse

import (
	"strings"
)

// String returns data as a continuous string.
func String(data []byte) string {
	return strings.TrimSpace(string(data))
}

// Lines returns data as a slice of strings corresponding to the lines of text in the input data.
func Lines(data []byte) []string {
	return strings.Split(strings.TrimSpace(string(data)), "\n")
}

// Parts returns data as a slice of strings corresponding to the parts of text in the input data
// separated by the given delimiter.
func Parts(data []byte, delimiter string) []string {
	return strings.Split(strings.TrimSpace(string(data)), delimiter)
}

// Matrix returns data as a matrix. Rows correspond to the lines of text in the input data.
func Matrix(data []byte) [][]string {
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	matrix := make([][]string, len(lines))

	for i, line := range lines {
		matrix[i] = strings.Split(line, "")
	}

	return matrix
}
`

const stringTmpl = `// Package render provides some ways to visualize data. Useful for debugging or just for fun.
package render

import (
	"fmt"
	"strings"
)

// Matrix returns a string representation of the given matrix. It can be printed
// to present a visual illustration of the two dimensional space.
func Matrix[T any](matrix [][]T) string {
	var b strings.Builder
	for y := range matrix {
		for x := range matrix[0] {
			b.WriteString(fmt.Sprint(matrix[y][x]))
		}
		if y < len(matrix)-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

// MapSlice returns the string representation of a slice while allowing an operation to be performed
// on each element of it.
//
// It is handy if you have a slice of structs and want to print a certain property of each item.
// The func extracts the property that's printed.
func MapSlice[T comparable](s []T, f func(T) any) string {
	var b strings.Builder
	b.WriteString("[")
	for _, v := range s {
		b.WriteString(fmt.Sprintf("%v ", f(v)))
	}
	return b.String()[:b.Len()-1] + "]"
}
`

const printTmpl = `// Package render provides some ways to visualize data. Useful for debugging or just for fun.
package render

import (
	"fmt"
	"strings"
	"time"
)

// Animate:
//  1. Prints the given string
//  2. Sleeps the given number of milliseconds
//  3. Retraces and removes any printed lines
//
// When called in a loop and provided a string representing a continuously changing state, the state
// changes will animate on the spot. The typical use case is probably animating the changing state
// of a 2D grid.
//
// Setting different values of delayMs will speed up or slow down the animation.
func Animate(delayMs int, frame string) {
	nLines := strings.Count(frame, "\n")
	fmt.Print(frame)
	time.Sleep(time.Duration(delayMs) * time.Millisecond)
	for range nLines {
		fmt.Print("\033[1A")
		fmt.Print("\033[2K\r")
	}
}
`

func (c Commands) GenAoc(module string) error {
	if files.Exists("go.mod") {
		fmt.Println("skipping go.mod, already exists")
	} else {
		fmt.Println("creating go.mod")
		if _, err := exec.CommandAndCapture("go", "mod", "init", module); err != nil {
			return fmt.Errorf("initiating go.mod: %v", err)
		}
	}

	if err := files.Gen(map[string]string{
		filepath.Join("shared", "parse", "input.go"):   parseTmpl,
		filepath.Join("shared", "exit", "error.go"):    exitTmpl,
		filepath.Join("shared", "render", "string.go"): stringTmpl,
		filepath.Join("shared", "render", "print.go"):  printTmpl,
	}, nil); err != nil {
		return fmt.Errorf("generating files: %v", err)
	}

	return nil
}
