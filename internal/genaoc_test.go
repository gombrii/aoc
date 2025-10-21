package internal_test

import (
	"path/filepath"
	"testing"

	"github.com/gombrii/aoc/internal"
)

// TODO: Will fail depending on current go version
func TestGenAoc(t *testing.T) {
	testRoot, _, wd := prepare(t)

	if err := internal.GenAoc("senap"); err != nil {
		t.Fail()
	}

	assertDirEqual(t, filepath.Join(wd, "testdata", "newmod"), testRoot)
}

//TODO: Add tests for all failing cases and check error and wether or not they generate anything
