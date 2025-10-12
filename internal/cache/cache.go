package cache

import (
	"fmt"
	"os"
	"path/filepath"
)

func Location() (string, error) {
	osCache, _ := os.UserCacheDir()
	return filepath.Join(osCache, "aoc-cache"), nil
}

func Contains(key string) (string, bool) {
	cache, _ := Location()
	path := filepath.Join(cache, key)

	if _, err := os.Stat(path); err != nil {
		return "", false
	}

	return filepath.Join(path, "runner.go"), true
}

func Store(key string, src string) (string, error) {
	cPath, _ := Location()
	dPath := filepath.Join(cPath, key)
	dst := filepath.Join(dPath, "runner.go")

	if _, err := os.Stat(src); err != nil {
		return "", fmt.Errorf("checking existance src: %v", err)
	}

	if err := os.MkdirAll(dPath, 0755); err != nil {
		return "", fmt.Errorf("creating cache dir: %v", err)
	}

	return dst, os.Rename(src, dst)
}

func Clear() error {
	cache, _ := Location()
	return os.RemoveAll(cache)
}
