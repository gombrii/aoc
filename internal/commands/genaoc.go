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
	return strings.Split(strings.TrimSpace(string(data)), "\n")
}

// Parts returns data as a slice of strings corresponding to the parts of text in the input data
// separated by the given separator.
func Parts(data []byte, separator string) []string {
	return strings.Split(strings.TrimSpace(string(data)), separator)
}

// String returns data as a continuous string.
func String(data []byte) string {
	return strings.TrimSpace(string(data))
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

type G struct {
	Symbol rune
	X      int
	Y      int
}

// GString returns a string representation of a grid of width w and height h. It can be printed to
// present a visual illustration of a two dimensional space.
//
// Glyphs is a slice containing information about symbols and their locations on the grid. Bg
// dictates the "background" for empty coordinates. Any coordinate is omitted that is in glyphs but
// outside the bounds of the grid. Glyph slice renders top to bottom, so for any two glyphs
// occupying the the same cell, the earlier one will be overwritten by the latter.
//
//	grid := render.GString(width, height, '.', []render.G{
//		render.G{Symbol: '^', X: 8, Y: 6},
//		render.G{Symbol: '#', X: 8, Y: 9},
//		render.G{Symbol: '#', X: 5, Y: 1},
//		render.G{Symbol: '#', X: 7, Y: 0},
//		render.G{Symbol: '#', X: 2, Y: 4},
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
//
// Use render.MString instead to get a string representation of an actual matrix ([][]T).
func GString(w, h int, bg rune, glyphs []G) string {
	lines := make([][]rune, h)
	for y := range h {
		lines[y] = make([]rune, w)
		for x := range w {
			lines[y][x] = bg
		}
	}
	for _, g := range glyphs {
		if g.X < 0 || g.X >= w || g.Y < 0 || g.Y >= h {
			continue
		}
		lines[g.Y][g.X] = g.Symbol
	}

	var b strings.Builder
	b.Grow(h*w + (h - 1))
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
