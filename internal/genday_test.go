package internal_test

import (
	"path/filepath"
	"testing"

	"github.com/gombrii/aoc/internal"
)

func TestGenDay(t *testing.T) {
	tmpRoot, _, wd := prepare(t)

	internal.GenDay(2024, 1)

	assertDirEqual(t, filepath.Join(wd, "testdata", "newday"), tmpRoot)
}
