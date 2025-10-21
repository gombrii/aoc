package internal_test

import (
	"path/filepath"
	"testing"

	"github.com/gombrii/aoc/internal"
	"github.com/otiai10/copy"
)

func TestRun(t *testing.T) {
	testRoot, _, wd := prepare(t)

	if err := copy.Copy(filepath.Join(wd, "testdata", "newmod"), testRoot); err != nil {
		t.Fatal(err)
	}

	if err := copy.Copy(filepath.Join(wd, "testdata", "newday", "2024"), filepath.Join(testRoot, "2024")); err != nil {
		t.Fatal(err)
	}

	if err := internal.Run(2024, 1, 1, "input.txt"); err != nil {
		t.Error(err)
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

			if err := copy.Copy(filepath.Join(wd, "testdata", "newmod"), testRoot); err != nil {
				t.Fatal(err)
			}

			if err := copy.Copy(filepath.Join(wd, "testdata", "newday", "2024"), filepath.Join(testRoot, "2024")); err != nil {
				t.Fatal(err)
			}

			srcPath := filepath.Join(wd, "testdata", "puzzlefiles", params.puzzleFile)
			dstPath := filepath.Join(testRoot, "2024", "solutions", "day1", "part1.go")

			if err := copy.Copy(srcPath, dstPath); err != nil {
				t.Fatal(err)
			}

			if err := internal.Run(2024, 1, 1, "input.txt"); err != nil {
				t.Error(err)
			}
		})
	}
}

//TODO: Add tests for all failing cases and check error and wether or not they generate anything
//replaceInFile(filepath.Join(testRoot, "2024", "solutions", "day1", "part1.go"), `"NOT IMPLEMENTED!"`)
