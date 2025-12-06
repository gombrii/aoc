package cache

import (
	"fmt"
	"strings"
)

type Key interface {
	ID() string
	namespace() string
}

type ConfigKey struct {
	Domain string
}

func (k ConfigKey) ID() string {
	return k.Domain

}
func (k ConfigKey) namespace() string {
	return "config"
}

type PuzzleKey struct {
	Year  int
	Day   int
	Part  int
	Input string
}

func (k PuzzleKey) ID() string {
	return fmt.Sprintf("%d-day%d-part%d-%s", k.Year, k.Day, k.Part, strings.TrimSuffix(k.Input, ".txt"))
}
func (k PuzzleKey) namespace() string {
	return "puzzles"
}
