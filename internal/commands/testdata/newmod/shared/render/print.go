// Package render provides some ways to visualize data. Useful for debugging or just for fun.
package render

import (
	"fmt"
	"strings"
	"time"
)

// Animate:
//  1. Prints the given string
//  2. Sleeps the given number of milliseconds
//  3. Retraces and removes any printed lines
//
// When called in a loop and provided a string representing a continuously changing state, the state
// changes will animate on the spot. The typical use case is probably animating the changing state
// of a 2D grid.
//
// Setting different values of delayMs will speed up or slow down the animation.
func Animate(delayMs int, frame string) {
	nLines := strings.Count(frame, "\n")
	fmt.Print(frame)
	time.Sleep(time.Duration(delayMs) * time.Millisecond)
	for range nLines {
		fmt.Print("\033[1A")
		fmt.Print("\033[2K\r")
	}
}
