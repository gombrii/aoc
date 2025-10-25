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
`

const parseTmpl = `// Package input does, in the spirit of Advent of Code, provide some common ways to interpret the 
// puzzles' input data.
package parse

import (
	"strings"
)

// Lines returns data as a slice of strings corresponding to the lines of text in the input data.
func Lines(data []byte) []string {
	return strings.Split(string(data), "\n")
}

// String returns data as a continuous string.
func String(data []byte) string {
	return string(data)
}

// Matrix returns data as a matrix. The delimiter divides the data into separates columns while rows
// correspond to the lines of text in the input data.
func Matrix(data []byte, delimiter string) [][]string {
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	matrix := make([][]string, len(lines))

	for i, line := range lines {
		matrix[i] = strings.Split(line, delimiter)
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

// GString returns a string representation of a grid of height h and width w. It can be printed
// to present a visual illustration of a two dimensional space.
//
// Points is a map containing
// information about symbols and their locations on the grid. The space parameter dictates the
// "background" for empty coordinates. Any coordinate is omitted that is in the points map but
// outside the bounds of the grid.
//
//	grid := render.GString(height, width, '.', map[rune][][2]int{
//		'^': {
//			{guard.posx, guard.posy},
//		},
//		'#': {
//			{8, 9},
//			{5, 1},
//			{7, 0},
//			{2, 4},
//		},
//	})
//
// The GString above gives the following result when printed.
//
//	.......#..
//	.....#....
//	..........
//	..........
//	..#.......
//	..........
//	........^.
//	..........
//	..........
//	........#.
func GString(h, w int, space rune, points map[rune][][2]int) string {
	lines := make([][]rune, h)
	for y := range h {
		lines[y] = make([]rune, w)
		for x := range w {
			lines[y][x] = space
		}
	}
	for symbol, coords := range points {
		for _, coord := range coords {
			x := coord[0]
			y := coord[1]
			if x < 0 || x >= w || y < 0 || y >= h {
				continue
			}
			lines[y][x] = symbol
		}
	}

	var b strings.Builder
	for i, l := range lines {
		b.WriteString(string(l))
		if i < len(lines)-1 {
			b.WriteByte('\n')
		}
	}

	return b.String()
}

// MString returns a string representation of the given matrix. It can be printed
// to present a visual illustration of the two dimensional space.
func MString[T any](matrix [][]T) string {
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
