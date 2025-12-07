// Package render provides some ways to visualize data. Useful for debugging or just for fun.
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
