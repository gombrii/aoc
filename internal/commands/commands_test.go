package commands_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otiai10/copy"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/fs"
)

func prepare(t *testing.T) (testRoot, testCache, wd string) {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("getting working directory: %v", err)
	}
	testRoot = t.TempDir()
	testCache = t.TempDir()
	t.Chdir(testRoot)
	t.Setenv("AOC_CACHE", testCache)

	return testRoot, testCache, wd
}

// initMod stores in testRoot the structure
//
//	testRoot/
//	├── go.mod
//	├── shared/
//	│   ├── exit/error.go
//	│   └── parse/input.go
func initMod(t *testing.T, wd, testRoot string) {
	if err := copy.Copy(filepath.Join(wd, "testdata", "newmod"), testRoot); err != nil {
		t.Fatalf("could not init mod: %v", err)
	}
}

// initDay stores in testRoot the structure
//
//	testRoot/
//	├── 2024/
//	│   ├── input/
//	│   │   └── day1/
//	│   │       ├── input.txt
//	│   │       └── test.txt
//	│   └── solutions/
//	│       └── day1/
//	│           ├── part1.go
//	│           ├── part2.go
//	│           └── common.go
func initDay(t *testing.T, wd, testRoot string) {
	if err := copy.Copy(filepath.Join(wd, "testdata", "newday", "2024"), filepath.Join(testRoot, "2024")); err != nil {
		t.Fatalf("could not init day: %v", err)
	}
}

// initCache stores in testCache the structure
//
//	testCache/
//	├── 2024-day1-part1-input/
//	│   ├── lock
//	│   ├── res
//	│   ├── dur
//	│   └── runner.go
func initCache(t *testing.T, wd, testCache string) {
	if err := copy.Copy(filepath.Join(wd, "testdata", "cache", "2024-day1-part1-input"), filepath.Join(testCache, "2024-day1-part1-input")); err != nil {
		t.Fatalf("could not init cache: %v", err)
	}
}

func assertEqual(t *testing.T, wd, expectedDir, actualDir string) {
	expected := fs.ManifestFromDir(t, expectedDir)
	assert.Assert(t, fs.Equal(actualDir, expected))
}
