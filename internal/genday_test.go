package internal_test

import (
	"path/filepath"
	"testing"

	"github.com/gombrii/aoc/internal"
)

// TODO: Will fail depending on current date
func TestGenDay(t *testing.T) {
	testRoot, _, wd := prepare(t)

	if err := internal.GenDay(2024, 1); err != nil {
		t.Fail()
	}

	assertDirEqual(t, filepath.Join(wd, "testdata", "newday"), testRoot)
}
