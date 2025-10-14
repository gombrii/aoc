package cache

import (
	"fmt"
	"os"
	"path/filepath"
)

func location() (string, error) {
	osCache, _ := os.UserCacheDir()
	return filepath.Join(osCache, "aoc-cache"), nil
}

func MakePath(key string, file string) string {
	cache, _ := location()
	return filepath.Join(cache, key, file)
}

func Contains(key string, file string) (string, bool) {
	cache, _ := location()
	path := filepath.Join(cache, key, file)

	if _, err := os.Stat(path); err != nil {
		return "", false
	}

	return path, true
}

func Store(key string, fileName string, src string) (string, error) {
	cPath, _ := location()
	dPath := filepath.Join(cPath, key)
	dst := filepath.Join(dPath, fileName)

	if _, err := os.Stat(src); err != nil {
		return "", fmt.Errorf("checking existance src: %v", err)
	}

	if err := os.MkdirAll(dPath, 0755); err != nil {
		return "", fmt.Errorf("creating cache dir: %v", err)
	}

	return dst, os.Rename(src, dst)
}

func Clear() error {
	cache, _ := location()
	return os.RemoveAll(cache)
}
