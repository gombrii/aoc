package cache

import (
	"fmt"
	"iter"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type key string

func Key(year, day, part int, input string) key {
	return key(fmt.Sprintf("%d-day%d-part%d-%s", year, day, part, strings.Replace(input, ".txt", "", 1)))
}

func MakePath(key key, file string) string {
	cache := location()
	return filepath.Join(cache, string(key), file)
}

func ContainsKey(key key) (string, bool) {
	cache := location()
	path := filepath.Join(cache, string(key))

	if _, err := os.Stat(path); err != nil {
		return "", false
	}

	return path, true
}

func Contains(key key, file string) (string, bool) {
	cache := location()
	path := filepath.Join(cache, string(key), file)

	if _, err := os.Stat(path); err != nil {
		return "", false
	}

	return path, true
}

func Store(key key, fileName string, src string) (string, error) {
	cPath := location()
	dPath := filepath.Join(cPath, string(key))
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
	cache := location()
	return os.RemoveAll(cache)
}

func All() iter.Seq2[int, string] {
	cache := location()
	entries, _ := os.ReadDir(cache)
	sort.Slice(entries, func(i, j int) bool {
		yi, di, pi := 0, 0, 0
		yj, dj, pj := 0, 0, 0
		_, _ = fmt.Sscanf(entries[i].Name(), "%d-day%d-part%d-", &yi, &di, &pi)
		_, _ = fmt.Sscanf(entries[j].Name(), "%d-day%d-part%d-", &yj, &dj, &pj)

		if yi != yj {
			return yi < yj
		}
		if di != dj {
			return di < dj
		}
		return pi < pj
	})
	return func(yield func(int, string) bool) {
		for i, e := range entries {
			if !yield(i, filepath.Join(cache, e.Name())) {
				return
			}
		}
	}
}
