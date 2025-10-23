package commands

import "github.com/gombrii/aoc/internal/cache"

func (c Commands) ClearCache() error {
	return cache.Clear()
}
