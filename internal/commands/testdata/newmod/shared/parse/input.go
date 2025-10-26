// Package input does, in the spirit of Advent of Code, provide some common ways to interpret the 
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
