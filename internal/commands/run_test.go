package commands_test

import (
	"path/filepath"
	"testing"

	"github.com/gombrii/aoc/internal/commands"
	"github.com/otiai10/copy"
)

//TODO: Add tests to check nothing is generated if erroring out early

func TestRun(t *testing.T) {
	testRoot, _, wd := prepare(t)

	initMod(t, wd, testRoot)
	initDay(t, wd, testRoot)

	if err := (commands.Commands{}).Run(2024, 1, 1, "input.txt"); err != nil {
		t.Errorf("calling Run: %v", err)
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
			testRoot, _, wd := prepare(t)

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
		})
	}
}
