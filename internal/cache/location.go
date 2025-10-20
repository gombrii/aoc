package cache

import (
	"os"
	"path/filepath"
)

const defaultName = "aoc-cache"

// location allows overriding cache location, for testing purposes
func location() string {
	osCache, _ := os.UserCacheDir()
	override := os.Getenv("AOC_CACHE")
	if override != "" {
		return override
	}
	return filepath.Join(osCache, defaultName)
}
