// Package input does, in the spirit of Advent of Code, provide some common ways to interpret the
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
