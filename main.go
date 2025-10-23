// Package main is the entry to the aoc app.
//
// For strictly testing purposes aoc listens to one environment variable:
//   - AOC_CACHE (overrides cache catalogue)
package main

import (
	"fmt"
	"os"

	"github.com/gombrii/aoc/internal/app"
	"github.com/gombrii/aoc/internal/commands"
)

//TODO: Clean up the command tests. Keep principles:
// - Don't test output
// - test that no errors despite erroneous circumstances
// - test correctly created and modified files from running commands
//
// And finally test files, exec and cache packages as well as of course app.

func main() {
	if err := app.Start(commands.Commands{}, os.Args[1:]...); err != nil {
		fmt.Println("Error:", err)
		os.Exit(2)
	}
}
