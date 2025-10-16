package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gombrii/aoc/internal/files"
)

const exit = `// Package exit does, in the spirit of Advent of Code, provide quick and dirty ways to fail.
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
}`

const input = `// Package input does, in the spirit of Advent of Code, provide some common ways to interpret the 
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
}`

func GenAoc(module string) error {
	if _, err := os.Stat("go.mod"); err != nil {
		err := exec.Command("go", "mod", "init", module).Run()
		if err != nil {
			return fmt.Errorf("initiating go.mod: %v", err)
		}
	}

	if err := files.Gen(map[string]string{
		filepath.Join("shared", "parse", "input.go"): input,
		filepath.Join("shared", "exit", "error.go"):  exit,
	}, nil); err != nil {
		return fmt.Errorf("generating files: %v", err)
	}

	fmt.Printf("Initiated Advent of Code module %s\n", module)

	return nil
}
