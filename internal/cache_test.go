package internal_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gombrii/aoc/internal"
)

func TestClearCache(t *testing.T) {
	_, testCache, _ := prepare(t)
	os.Mkdir(filepath.Join(testCache, "cached-dir"), 0755)

	if err := internal.ClearCache(); err != nil {
		t.Fail()
	}

	if _, err := os.Stat(testCache); err == nil {
		t.Fail()
	}
}

func TestClearNonexistentCache(t *testing.T) {
	_, testCache, _ := prepare(t)

	if err := internal.ClearCache(); err != nil {
		t.Fail()
	}

	if err := internal.ClearCache(); err != nil {
		t.Fail()
	}

	if _, err := os.Stat(testCache); err == nil {
		t.Fail()
	}
}
