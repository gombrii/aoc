// Package render provides some ways to visualize data. Useful for debugging or just for fun.
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