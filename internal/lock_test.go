package internal_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gombrii/aoc/internal"
	"github.com/otiai10/copy"
)

// TODO: STOP TESTING OUTPUT: TEST ONLY WHAT EFFECTS OPERATIONS HAVE ON THE FILESYSTEM, GENERATED FILES, CHANGED VALUES, etc.
func TestLock(t *testing.T) {
	_, testCache, wd := prepare(t)

	err := copy.Copy(filepath.Join(wd, "testdata", "cache", "2024-day1-part1-input"), filepath.Join(testCache, "2024-day1-part1-input"))
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(filepath.Join(testCache, "2024-day1-part1-input", "lock"))
	if strings.TrimSpace(string(data)) != "false" {
		t.Error("lock didn't start as false")
	}

	if err = internal.Lock(2024, 1, 1, "input.txt"); err != nil {
		t.Error(err)
	}

	data, _ = os.ReadFile(filepath.Join(testCache, "2024-day1-part1-input", "lock"))
	if string(data) != "true" {
		t.Error("locking didn't set lock to true")
	}

	internal.Unlock(2024, 1, 1, "input.txt")

	data, _ = os.ReadFile(filepath.Join(testCache, "2024-day1-part1-input", "lock"))
	if string(data) != "false" {
		t.Error("unlocking didn't set lock to false")
	}
}

func TestStatusNotExists(t *testing.T) {
	_, _, _ = prepare(t)
	if err := internal.Lock(2024, 1, 1, "input.txt"); err != nil {
		t.Error(err)
	}
}
