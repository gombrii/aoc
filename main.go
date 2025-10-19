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
