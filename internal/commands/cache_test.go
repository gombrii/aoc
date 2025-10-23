package commands_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gombrii/aoc/internal/commands"
)

func TestClearCache(t *testing.T) {
	_, testCache, _ := prepare(t)
	if err := os.Mkdir(filepath.Join(testCache, "cached-dir"), 0755); err != nil {
		t.Errorf("creating dir in cache: %v", err)
	}

	if err := (commands.Commands{}).ClearCache(); err != nil {
		t.Errorf("calling ClearCache: %v", err)
	}

	if _, err := os.Stat(testCache); err == nil {
		t.Error("cache wasn't cleared")
	}
}

func TestClearNonexistentCache(t *testing.T) {
	_, testCache, _ := prepare(t)
	cmd := commands.Commands{}

	if err := cmd.ClearCache(); err != nil {
		t.Errorf("calling ClearCache: %v", err)
	}

	if err := cmd.ClearCache(); err != nil {
		t.Errorf("calling ClearCache: %v", err)
	}

	if _, err := os.Stat(testCache); err == nil {
		t.Error("cache wasn't cleared")
	}
}
