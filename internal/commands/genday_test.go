package commands_test

import (
	"path/filepath"
	"testing"

	"github.com/gombrii/aoc/internal/commands"
)

//TODO: Add tests to check nothing is generated if erroring out early

// TODO: Will fail depending on current date
func TestGenDay(t *testing.T) {
	testRoot, _, wd := prepare(t)

	if err := (commands.Commands{}).GenDay(2024, 1); err != nil {
		t.Errorf("calling GenDay: %v", err)
	}

	assertEqual(t, wd, testRoot, filepath.Join(wd, "testdata", "newday"))
}

func TestGenDayWithFilesPresent(t *testing.T) {
	testRoot, _, wd := prepare(t)

	if err := (commands.Commands{}).GenDay(2024, 1); err != nil {
		t.Errorf("calling GenDay: %v", err)
	}

	if err := (commands.Commands{}).GenDay(2024, 1); err != nil {
		t.Errorf("calling GenDay: %v", err)
	}

	assertEqual(t, wd, testRoot, filepath.Join(wd, "testdata", "newday"))
}
