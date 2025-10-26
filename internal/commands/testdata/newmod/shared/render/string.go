// Package render provides some ways to visualize data. Useful for debugging or just for fun.
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
