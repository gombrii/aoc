package internal

import "github.com/gombrii/aoc/internal/cache"

func ClearCache() error {
	return cache.Clear()
}