// Package main is the entry to the aoc app.
//
// For strictly testing purposes aoc listens to one environment variable:
//   - AOC_CACHE (overrides cache catalogue)
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/gombrii/aoc/internal/app"
	"github.com/gombrii/aoc/internal/commands"
)

func main() {
	if err := app.Start(commands.Commands{}, os.Args[1:]...); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		if errors.Is(err, app.ErrInput) {
			os.Exit(2)
		}
		os.Exit(1)
	}
}
