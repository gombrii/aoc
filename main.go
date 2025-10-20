// Package main is the entry to the aoc app
//
// For strictly testing purposes aoc listenes to two environment variables
// - AOC_HOME (overrides "your" catalogue)
// - AOC_CACHE (overrides cache catalogue)
package main

import (
	"fmt"
	"os"

	"github.com/gombrii/aoc/internal"
)

func main() {
	if err := start(commands{
		Run:        internal.Run,
		Status:     internal.Status,
		Lock:       internal.Lock,
		Unlock:     internal.Unlock,
		GenDay:     internal.GenDay,
		GenAoc:     internal.GenAoc,
		Check:      internal.Check,
		ClearCache: internal.ClearCache,
	}, os.Args[1:]...); err != nil {
		fmt.Println("Error:", err)
		os.Exit(2)
	}
}
