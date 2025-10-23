package commands_test

import (
	"path/filepath"
	"testing"

	"github.com/gombrii/aoc/internal/commands"
)

// TODO: Will fail depending on current go version
func TestGenAoc(t *testing.T) {
	testRoot, _, wd := prepare(t)

	if err := (commands.Commands{}).GenAoc("senap"); err != nil {
		t.Errorf("calling GenAoc: %v", err)
	}

	assertEqual(t, wd, testRoot, filepath.Join(wd, "testdata", "newmod"))
}

func TestGenAocWithFilesPresent(t *testing.T) {
	testRoot, _, wd := prepare(t)

	if err := (commands.Commands{}).GenAoc("senap"); err != nil {
		t.Errorf("calling GenAoc: %v", err)
	}

	if err := (commands.Commands{}).GenAoc("senap"); err != nil {
		t.Errorf("calling GenAoc: %v", err)
	}

	assertEqual(t, wd, testRoot, filepath.Join(wd, "testdata", "newmod"))
}
