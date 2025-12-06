package commands_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gombrii/aoc/internal/commands"
	"github.com/otiai10/copy"
)

func TestRun(t *testing.T) {
	testRoot, testCache, wd := prepare(t)

	initMod(t, wd, testRoot)
	initDay(t, wd, testRoot)

	if _, err := os.Stat(filepath.Join(testCache, "puzzles", "2024-day1-part1-input")); err == nil {
		t.Error("Cache already exists")
	}

	if err := (commands.Commands{}).Run(2024, 1, 1, "input.txt"); err != nil {
		t.Errorf("calling Run: %v", err)
	}

	if _, err := os.Stat(filepath.Join(testCache, "puzzles", "2024-day1-part1-input")); err != nil {
		t.Error("Run wasn't cached")
	}
}

func TestRunError(t *testing.T) {
	for name, params := range map[string]struct {
		puzzleFile string
	}{
		"exit": {
			puzzleFile: "exit.go",
		},
		"panic": {
			puzzleFile: "panic.go",
		},
		"compilation failure": {
			puzzleFile: "compilefail.go",
		},
	} {
		t.Run(name, func(t *testing.T) {
			testRoot, testCache, wd := prepare(t)

			initMod(t, wd, testRoot)
			initDay(t, wd, testRoot)

			srcPath := filepath.Join(wd, "testdata", "puzzlefiles", params.puzzleFile)
			dstPath := filepath.Join(testRoot, "2024", "solutions", "day1", "part1.go")

			if err := copy.Copy(srcPath, dstPath); err != nil {
				t.Fatalf("replacing part1.go in testDir: %v", err)
			}

			if err := (commands.Commands{}).Run(2024, 1, 1, "input.txt"); err != nil {
				t.Errorf("calling Run: %v", err)
			}

			if _, err := os.Stat(filepath.Join(testCache, "puzzles", "2024-day1-part1-input")); err != nil {
				t.Error("Run wasn't cached")
			}
		})
	}
}

func TestRunNoMod(t *testing.T) {
	testRoot, testCache, wd := prepare(t)

	initDay(t, wd, testRoot)

	if err := (commands.Commands{}).Run(2024, 1, 1, "input.txt"); err == nil {
		t.Error("Calling Run outside module did not return an error")
	}

	if _, err := os.Stat(filepath.Join(testCache, "2024-day1-part1-input")); err == nil {
		t.Error("Cache was created despite Run returning error")
	}
}

func TestRunNoDay(t *testing.T) {
	testRoot, testCache, wd := prepare(t)

	initMod(t, wd, testRoot)

	if err := (commands.Commands{}).Run(2024, 1, 1, "input.txt"); err == nil {
		t.Error("Calling Run without day target did not return an error")
	}

	if _, err := os.Stat(filepath.Join(testCache, "2024-day1-part1-input")); err == nil {
		t.Error("Cache was created despite Run returning error")
	}
}

func TestRunNeitherModNorDay(t *testing.T) {
	_, testCache, _ := prepare(t)

	if err := (commands.Commands{}).Run(2024, 1, 1, "input.txt"); err == nil {
		t.Error("Calling Run outside mod and without day target did not return an error")
	}

	if _, err := os.Stat(filepath.Join(testCache, "2024-day1-part1-input")); err == nil {
		t.Error("Cache was created despite Run returning error")
	}
}
