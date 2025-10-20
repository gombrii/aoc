package internal_test

import (
	"path/filepath"
	"testing"

	"github.com/gombrii/aoc/internal"
)

func TestGenAoc(t *testing.T) {
	tmpRoot, _, wd := prepare(t)

	internal.GenAoc("senap")

	assertDirEqual(t, filepath.Join(wd, "testdata", "newmod"), tmpRoot)
}
